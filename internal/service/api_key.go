/*
 * 文件作用：API Key管理服务，处理API密钥的业务逻辑
 * 负责功能：
 *   - API Key CRUD操作
 *   - API Key 验证
 *   - API Key 状态管理
 *   - 使用量统计
 * 重要程度：⭐⭐⭐⭐ 重要（API Key管理核心）
 * 依赖模块：repository, model
 */
package service

import (
	"errors"
	"sync"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/repository"
	"cli-proxy/pkg/logger"
)

var (
	apiKeyLog     *logger.Logger
	apiKeyLogOnce sync.Once
)

func getAPIKeyLog() *logger.Logger {
	apiKeyLogOnce.Do(func() {
		apiKeyLog = logger.GetLogger("main")
	})
	return apiKeyLog
}

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService() *APIKeyService {
	return &APIKeyService{
		repo: repository.NewAPIKeyRepository(),
	}
}

// CreateAPIKeyRequest 创建 API Key 请求
type CreateAPIKeyRequest struct {
	Name             string     `json:"name" binding:"required"`
	Description      string     `json:"description"`       // 描述用途
	AllowedPlatforms string     `json:"allowed_platforms"` // 允许的平台
	AllowedModels    string     `json:"allowed_models"`    // 允许的模型
	BlockedModels    string     `json:"blocked_models"`    // 禁止的模型
	AllowedClients   string     `json:"allowed_clients"`   // 允许的客户端
	RateLimit        int        `json:"rate_limit"`        // 每分钟请求限制
	DailyLimit       int        `json:"daily_limit"`       // 每日请求限制
	MonthlyQuota     float64    `json:"monthly_quota"`     // 月额度
	ExpiresAt        *time.Time `json:"expires_at"`        // 过期时间
}

// CreateAPIKeyResponse 创建 API Key 响应 (只在创建时返回完整 key)
type CreateAPIKeyResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`        // 只在创建时返回
	KeyPrefix string `json:"key_prefix"` // 用于显示
}

// Create 创建新的 API Key
func (s *APIKeyService) Create(req *CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	getAPIKeyLog().Info("[apikey] 创建 API Key 请求 | Name: %s", req.Name)

	// 生成新的 API Key
	key, hash, prefix, err := model.GenerateAPIKey()
	if err != nil {
		getAPIKeyLog().Error("[apikey] 创建 API Key 失败 | 原因: 生成 Key 失败: %v", err)
		return nil, errors.New("生成 API Key 失败")
	}

	// 设置默认值
	rateLimit := 60
	if req.RateLimit > 0 {
		rateLimit = req.RateLimit
	}

	allowedPlatforms := "all"
	if req.AllowedPlatforms != "" {
		allowedPlatforms = req.AllowedPlatforms
	}

	apiKey := &model.APIKey{
		Name:             req.Name,
		Description:      req.Description,
		KeyHash:          hash,
		KeyPrefix:        prefix,
		Status:           "active",
		AllowedPlatforms: allowedPlatforms,
		AllowedModels:    req.AllowedModels,
		BlockedModels:    req.BlockedModels,
		AllowedClients:   req.AllowedClients,
		RateLimit:        rateLimit,
		DailyLimit:       req.DailyLimit,
		MonthlyQuota:     req.MonthlyQuota,
		ExpiresAt:        req.ExpiresAt,
	}

	if err := s.repo.Create(apiKey); err != nil {
		getAPIKeyLog().Error("[apikey] 创建 API Key 失败 | 原因: 数据库错误: %v", err)
		return nil, err
	}

	getAPIKeyLog().Info("[apikey] 创建 API Key 成功 | KeyID: %d | Name: %s", apiKey.ID, apiKey.Name)

	return &CreateAPIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Key:       key, // 只在创建时返回完整 key
		KeyPrefix: prefix,
	}, nil
}

// List 获取所有 API Key（分页）
func (s *APIKeyService) List(page, pageSize int) ([]model.APIKey, int64, error) {
	return s.repo.List(page, pageSize)
}

// ListAll 获取所有 API Key（不分页）
func (s *APIKeyService) ListAll() ([]model.APIKey, error) {
	return s.repo.ListAll()
}

// GetByID 根据 ID 获取 API Key
func (s *APIKeyService) GetByID(id uint) (*model.APIKey, error) {
	return s.repo.GetByID(id)
}

// Delete 删除 API Key
func (s *APIKeyService) Delete(id uint) error {
	getAPIKeyLog().Info("[apikey] 删除 API Key 请求 | KeyID: %d", id)
	if err := s.repo.Delete(id); err != nil {
		getAPIKeyLog().Error("[apikey] 删除 API Key 失败 | KeyID: %d | 原因: %v", id, err)
		return err
	}
	getAPIKeyLog().Info("[apikey] 删除 API Key 成功 | KeyID: %d", id)
	return nil
}

// ToggleStatus 切换 API Key 状态
func (s *APIKeyService) ToggleStatus(id uint) (*model.APIKey, error) {
	getAPIKeyLog().Info("[apikey] 切换 API Key 状态请求 | KeyID: %d", id)
	key, err := s.repo.GetByID(id)
	if err != nil {
		getAPIKeyLog().Error("[apikey] 切换 API Key 状态失败 | KeyID: %d | 原因: 查询失败: %v", id, err)
		return nil, err
	}

	oldStatus := key.Status
	if key.Status == "active" {
		key.Status = "disabled"
	} else {
		key.Status = "active"
	}

	if err := s.repo.Update(key); err != nil {
		getAPIKeyLog().Error("[apikey] 切换 API Key 状态失败 | KeyID: %d | 原因: 更新失败: %v", id, err)
		return nil, err
	}

	getAPIKeyLog().Info("[apikey] 切换 API Key 状态成功 | KeyID: %d | %s -> %s", id, oldStatus, key.Status)
	return key, nil
}

// ValidateKey 验证 API Key 并返回 Key 信息
func (s *APIKeyService) ValidateKey(keyStr string) (*model.APIKey, error) {
	hash := model.HashAPIKey(keyStr)
	key, err := s.repo.GetByHash(hash)
	if err != nil {
		return nil, errors.New("无效的 API Key")
	}

	if !key.IsActive() {
		if key.Status == "disabled" {
			return nil, errors.New("API Key 已被禁用")
		}
		if key.IsExpired() {
			return nil, errors.New("API Key 已过期")
		}
		return nil, errors.New("API Key 不可用")
	}

	return key, nil
}

// IncrementUsage 记录使用统计
func (s *APIKeyService) IncrementUsage(id uint, tokens int64, cost float64) error {
	return s.repo.IncrementUsage(id, tokens, cost)
}

// UpdateAPIKeyRequest 更新 API Key 请求
type UpdateAPIKeyRequest struct {
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	AllowedPlatforms string     `json:"allowed_platforms"`
	AllowedModels    string     `json:"allowed_models"`
	BlockedModels    string     `json:"blocked_models"`
	AllowedClients   string     `json:"allowed_clients"`
	RateLimit        int        `json:"rate_limit"`
	DailyLimit       int        `json:"daily_limit"`
	MonthlyQuota     float64    `json:"monthly_quota"`
	ExpiresAt        *time.Time `json:"expires_at"`
	Status           string     `json:"status"`
}

// Update 更新 API Key
func (s *APIKeyService) Update(id uint, req *UpdateAPIKeyRequest) (*model.APIKey, error) {
	key, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		key.Name = req.Name
	}
	if req.Description != "" {
		key.Description = req.Description
	}
	if req.AllowedPlatforms != "" {
		key.AllowedPlatforms = req.AllowedPlatforms
	}
	// 允许清空以下字段
	key.AllowedModels = req.AllowedModels
	key.BlockedModels = req.BlockedModels
	key.AllowedClients = req.AllowedClients

	if req.RateLimit > 0 {
		key.RateLimit = req.RateLimit
	}
	key.DailyLimit = req.DailyLimit
	key.MonthlyQuota = req.MonthlyQuota
	key.ExpiresAt = req.ExpiresAt

	if req.Status != "" {
		key.Status = req.Status
	}

	if err := s.repo.Update(key); err != nil {
		return nil, err
	}

	return key, nil
}

// Lookup 按 ID 批量查询 API Key（用于前端显示映射）
func (s *APIKeyService) Lookup(ids []uint) ([]model.APIKey, error) {
	return s.repo.GetByIDs(ids)
}

// GetAPIKeyLogs 获取 API Key 的使用日志
func (s *APIKeyService) GetAPIKeyLogs(keyID uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	return s.repo.GetAPIKeyLogs(keyID, page, pageSize)
}

// Count 统计 API Key 数量
func (s *APIKeyService) Count() (int64, error) {
	return s.repo.Count()
}
