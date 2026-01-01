/*
 * 文件作用：管理员处理器，处理管理员认证相关请求
 * 负责功能：
 *   - 管理员登录
 *   - JWT Token 生成
 * 重要程度：⭐⭐⭐⭐⭐ 核心（认证核心）
 * 依赖模块：service, response
 */
package handler

import (
	"cli-proxy/internal/service"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员处理器
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		adminService: service.NewAdminService(),
	}
}

// LoginRequest 登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 管理员登录
func (h *AdminHandler) Login(c *gin.Context) {
	var req service.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效: "+err.Error())
		return
	}

	result, err := h.adminService.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.Success(c, result)
}
