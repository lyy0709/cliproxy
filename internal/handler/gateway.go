/*
 * 文件作用：网关管理 API 处理器
 * 负责功能：
 *   - 网关 CRUD 接口
 *   - 网关测试接口
 *   - 默认网关管理接口
 * 重要程度：⭐⭐⭐ 重要（xyrt 授权核心）
 * 依赖模块：service, response
 */
package handler

import (
	"strconv"

	"cli-proxy/internal/service"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

// GatewayHandler 网关管理处理器
type GatewayHandler struct {
	service *service.GatewayService
}

// NewGatewayHandler 创建网关处理器实例
func NewGatewayHandler() *GatewayHandler {
	return &GatewayHandler{
		service: service.NewGatewayService(),
	}
}

// List 获取网关列表
// GET /admin/gateways
func (h *GatewayHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	gateways, total, err := h.service.List(page, pageSize, search)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, gateways, total, page, pageSize)
}

// GetEnabled 获取启用的网关列表（供下拉选择）
// GET /admin/gateways/enabled
func (h *GatewayHandler) GetEnabled(c *gin.Context) {
	gateways, err := h.service.GetEnabled()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items": gateways,
		"total": len(gateways),
	})
}

// Get 获取单个网关
// GET /admin/gateways/:id
func (h *GatewayHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的网关 ID")
		return
	}

	gateway, err := h.service.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "网关不存在")
		return
	}

	response.Success(c, gateway)
}

// Create 创建网关
// POST /admin/gateways
func (h *GatewayHandler) Create(c *gin.Context) {
	var req service.CreateGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	gateway, err := h.service.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, gateway)
}

// Update 更新网关
// PUT /admin/gateways/:id
func (h *GatewayHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的网关 ID")
		return
	}

	var req service.UpdateGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	gateway, err := h.service.Update(uint(id), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gateway)
}

// Delete 删除网关
// DELETE /admin/gateways/:id
func (h *GatewayHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的网关 ID")
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ToggleEnabled 切换启用状态
// PUT /admin/gateways/:id/toggle
func (h *GatewayHandler) ToggleEnabled(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的网关 ID")
		return
	}

	gateway, err := h.service.ToggleEnabled(uint(id))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gateway)
}

// GetDefault 获取默认网关
// GET /admin/gateways/default
func (h *GatewayHandler) GetDefault(c *gin.Context) {
	gateway, err := h.service.GetDefault()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gateway)
}

// SetDefault 设置默认网关
// PUT /admin/gateways/:id/default
func (h *GatewayHandler) SetDefault(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的网关 ID")
		return
	}

	if err := h.service.SetDefault(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ClearDefault 清除默认网关
// DELETE /admin/gateways/default
func (h *GatewayHandler) ClearDefault(c *gin.Context) {
	if err := h.service.ClearDefault(); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Test 测试网关连通性
// POST /admin/gateways/test
func (h *GatewayHandler) Test(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	latency, err := h.service.TestConnectivity(req.URL)
	if err != nil {
		response.Success(c, gin.H{
			"status":  "failed",
			"latency": 0,
			"error":   err.Error(),
		})
		return
	}

	response.Success(c, gin.H{
		"status":  "success",
		"latency": latency,
		"error":   "",
	})
}

// TestByID 测试指定网关并更新状态
// POST /admin/gateways/:id/test
func (h *GatewayHandler) TestByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的网关 ID")
		return
	}

	gateway, err := h.service.TestAndUpdate(uint(id))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gateway)
}
