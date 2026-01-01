/*
 * 文件作用：API Key管理处理器，处理API密钥的生成和管理
 * 负责功能：
 *   - API Key 列表查询
 *   - API Key 创建（管理员）
 *   - API Key 删除/禁用
 *   - API Key 使用量统计
 * 重要程度：⭐⭐⭐⭐ 重要（API Key管理核心）
 * 依赖模块：service
 */
package handler

import (
	"net/http"
	"strconv"
	"strings"

	"cli-proxy/internal/middleware"
	"cli-proxy/internal/service"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

type APIKeyHandler struct {
	service *service.APIKeyService
}

func NewAPIKeyHandler() *APIKeyHandler {
	return &APIKeyHandler{
		service: service.NewAPIKeyService(),
	}
}

// ========== 管理员接口 ==========

// AdminListAll 管理员获取所有 API Key（分页）
func (h *APIKeyHandler) AdminListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	keys, total, err := h.service.List(page, pageSize)
	if err != nil {
		response.InternalError(c, "获取 API Key 列表失败")
		return
	}

	response.Success(c, gin.H{
		"items": keys,
		"total": total,
		"page":  page,
	})
}

// AdminCreate 管理员创建 API Key
func (h *APIKeyHandler) AdminCreate(c *gin.Context) {
	var req service.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求数据")
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, result)
}

// AdminGet 管理员获取单个 API Key
func (h *APIKeyHandler) AdminGet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	key, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, key)
}

// AdminUpdate 管理员更新 API Key
func (h *APIKeyHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	var req service.UpdateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "无效的请求数据")
		return
	}

	key, err := h.service.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, key)
}

// AdminDelete 管理员删除 API Key
func (h *APIKeyHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 API Key ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// AdminToggleStatus 管理员切换 API Key 状态
func (h *APIKeyHandler) AdminToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 API Key ID")
		return
	}

	key, err := h.service.ToggleStatus(uint(id))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"status": key.Status})
}

// AdminLookup 管理员按 ID 批量查询 API Key（用于前端显示 sk- 前缀/完整 Key）
// GET /api/api-keys/lookup?ids=1,2,3
func (h *APIKeyHandler) AdminLookup(c *gin.Context) {
	idsStr := strings.TrimSpace(c.Query("ids"))
	if idsStr == "" {
		response.Success(c, gin.H{"items": []gin.H{}})
		return
	}

	parts := strings.Split(idsStr, ",")
	if len(parts) > 200 {
		response.BadRequest(c, "ids too many (max 200)")
		return
	}

	ids := make([]uint, 0, len(parts))
	seen := make(map[uint]struct{}, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		u64, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			response.BadRequest(c, "invalid id: "+p)
			return
		}
		id := uint(u64)
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	keys, err := h.service.Lookup(ids)
	if err != nil {
		response.InternalError(c, "获取 API Key 信息失败")
		return
	}

	// 仅返回前端映射需要的字段，减少 payload 和渲染压力
	items := make([]gin.H, 0, len(keys))
	for _, k := range keys {
		items = append(items, gin.H{
			"id":         k.ID,
			"key_prefix": k.KeyPrefix,
		})
	}

	response.Success(c, gin.H{"items": items})
}

// AdminGetAPIKeyLogs 管理员获取 API Key 的使用日志
func (h *APIKeyHandler) AdminGetAPIKeyLogs(c *gin.Context) {
	keyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 API Key ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.service.GetAPIKeyLogs(uint(keyID), page, pageSize)
	if err != nil {
		response.InternalError(c, "获取使用日志失败")
		return
	}

	response.Success(c, gin.H{
		"items": logs,
		"total": total,
		"page":  page,
	})
}

// GetMyInfo API Key 持有者获取当前 Key 信息
func (h *APIKeyHandler) GetMyInfo(c *gin.Context) {
	key := middleware.GetAPIKey(c)
	if key == nil {
		response.Unauthorized(c, "API Key not found in context")
		return
	}

	response.Success(c, gin.H{
		"id":                key.ID,
		"name":              key.Name,
		"description":       key.Description,
		"key_prefix":        key.KeyPrefix,
		"status":            key.Status,
		"allowed_platforms": key.AllowedPlatforms,
		"allowed_models":    key.AllowedModels,
		"blocked_models":    key.BlockedModels,
		"allowed_clients":   key.AllowedClients,
		"rate_limit":        key.RateLimit,
		"daily_limit":       key.DailyLimit,
		"monthly_quota":     key.MonthlyQuota,
		"expires_at":        key.ExpiresAt,
		"created_at":        key.CreatedAt,
	})
}

// ========== API Key 验证接口（供代理服务使用）==========

// Validate 验证 API Key (供代理服务使用)
func (h *APIKeyHandler) Validate(c *gin.Context) {
	apiKey := c.GetHeader("Authorization")
	if apiKey == "" {
		apiKey = c.GetHeader("X-API-Key")
	}
	if apiKey == "" {
		response.Unauthorized(c, "缺少 API Key")
		return
	}

	// 移除 Bearer 前缀
	if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
		apiKey = apiKey[7:]
	}

	key, err := h.service.ValidateKey(apiKey)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, gin.H{
		"valid":             true,
		"key_id":            key.ID,
		"allowed_platforms": key.AllowedPlatforms,
		"allowed_models":    key.AllowedModels,
		"rate_limit":        key.RateLimit,
	})
}
