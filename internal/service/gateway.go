/*
 * 文件作用：网关业务逻辑层
 * 负责功能：
 *   - 网关 CRUD 业务逻辑
 *   - 网关连通性测试
 *   - 默认网关管理
 * 重要程度：⭐⭐⭐ 重要（xyrt 授权核心）
 * 依赖模块：repository, model
 */
package service

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/repository"
)

// GatewayService 网关服务
type GatewayService struct {
	repo *repository.GatewayRepository
}

// NewGatewayService 创建网关服务实例
func NewGatewayService() *GatewayService {
	return &GatewayService{
		repo: repository.NewGatewayRepository(),
	}
}

// CreateGatewayRequest 创建网关请求
type CreateGatewayRequest struct {
	Name    string `json:"name" binding:"required"`
	URL     string `json:"url" binding:"required"`
	Enabled bool   `json:"enabled"`
	Remark  string `json:"remark"`
}

// UpdateGatewayRequest 更新网关请求
type UpdateGatewayRequest struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Enabled *bool  `json:"enabled"`
	Remark  string `json:"remark"`
}

// Create 创建网关
func (s *GatewayService) Create(req *CreateGatewayRequest) (*model.Gateway, error) {
	// 规范化 URL（去除末尾斜杠）
	url := strings.TrimSuffix(strings.TrimSpace(req.URL), "/")
	if url == "" {
		return nil, errors.New("网关地址不能为空")
	}

	gateway := &model.Gateway{
		Name:    strings.TrimSpace(req.Name),
		URL:     url,
		Enabled: req.Enabled,
		Remark:  strings.TrimSpace(req.Remark),
	}

	if err := s.repo.Create(gateway); err != nil {
		return nil, err
	}

	return gateway, nil
}

// GetByID 根据 ID 获取网关
func (s *GatewayService) GetByID(id uint) (*model.Gateway, error) {
	return s.repo.GetByID(id)
}

// Update 更新网关
func (s *GatewayService) Update(id uint, req *UpdateGatewayRequest) (*model.Gateway, error) {
	gateway, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		gateway.Name = strings.TrimSpace(req.Name)
	}
	if req.URL != "" {
		gateway.URL = strings.TrimSuffix(strings.TrimSpace(req.URL), "/")
	}
	if req.Enabled != nil {
		gateway.Enabled = *req.Enabled
	}
	if req.Remark != "" {
		gateway.Remark = strings.TrimSpace(req.Remark)
	}

	if err := s.repo.Update(gateway); err != nil {
		return nil, err
	}

	return gateway, nil
}

// Delete 删除网关
func (s *GatewayService) Delete(id uint) error {
	// 检查是否有账户在使用
	count, err := s.repo.CountAccountsUsingGateway(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该网关正在被账户使用，无法删除")
	}

	return s.repo.Delete(id)
}

// List 分页获取网关列表
func (s *GatewayService) List(page, pageSize int, search string) ([]model.Gateway, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.List(page, pageSize, search)
}

// GetEnabled 获取所有启用的网关
func (s *GatewayService) GetEnabled() ([]model.Gateway, error) {
	return s.repo.GetEnabled()
}

// GetDefault 获取默认网关
func (s *GatewayService) GetDefault() (*model.Gateway, error) {
	return s.repo.GetDefault()
}

// SetDefault 设置默认网关
func (s *GatewayService) SetDefault(id uint) error {
	// 验证网关存在且启用
	gateway, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if !gateway.Enabled {
		return errors.New("只能将启用的网关设为默认")
	}

	return s.repo.SetDefault(id)
}

// ClearDefault 清除默认网关
func (s *GatewayService) ClearDefault() error {
	return s.repo.ClearDefault()
}

// ToggleEnabled 切换启用状态
func (s *GatewayService) ToggleEnabled(id uint) (*model.Gateway, error) {
	if err := s.repo.ToggleEnabled(id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(id)
}

// TestConnectivity 测试网关连通性
func (s *GatewayService) TestConnectivity(url string) (int, error) {
	// 规范化 URL
	url = strings.TrimSuffix(strings.TrimSpace(url), "/")
	if url == "" {
		return 0, errors.New("网关地址不能为空")
	}

	// 创建 HTTP 客户端，设置超时
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 测试连通性（简单 HEAD 请求）
	start := time.Now()
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	latency := int(time.Since(start).Milliseconds())

	// 只要能连接就认为成功（不检查具体状态码）
	return latency, nil
}

// TestAndUpdate 测试网关并更新状态
func (s *GatewayService) TestAndUpdate(id uint) (*model.Gateway, error) {
	gateway, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	latency, testErr := s.TestConnectivity(gateway.URL)

	var status string
	var errMsg string
	if testErr != nil {
		status = "failed"
		errMsg = testErr.Error()
	} else {
		status = "success"
		errMsg = ""
	}

	if err := s.repo.UpdateTestStatus(id, status, latency, errMsg); err != nil {
		return nil, err
	}

	// 重新获取更新后的网关
	return s.repo.GetByID(id)
}
