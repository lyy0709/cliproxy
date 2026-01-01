/*
 * 文件作用：管理员处理器，处理管理员认证相关请求
 * 负责功能：
 *   - 管理员登录
 *   - JWT Token 生成
 *   - 密码修改
 *   - 首次登录强制修改密码
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

	configService := service.GetConfigService()
	clientIP := c.ClientIP()

	// 登录频率限制
	if configService.GetLoginRateLimitEnabled() {
		loginLimiter := service.GetLoginRateLimiter()
		allowed, waitSeconds := loginLimiter.Check("login:"+clientIP, configService.GetLoginRateLimitCount(), configService.GetLoginRateLimitWindow())
		if !allowed {
			response.TooManyRequests(c, service.GetRateLimitError(waitSeconds))
			return
		}
	}

	// 验证验证码（启用时强制）
	if configService.GetCaptchaEnabled() {
		if req.CaptchaId == "" || req.CaptchaCode == "" {
			response.BadRequest(c, "请输入验证码")
			return
		}
		captchaService := service.GetCaptchaService()
		if !captchaService.Verify(req.CaptchaId, req.CaptchaCode) {
			response.BadRequest(c, "验证码错误或已过期")
			return
		}
	}

	result, err := h.adminService.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	// 登录成功后重置频率限制计数
	if configService.GetLoginRateLimitEnabled() {
		service.GetLoginRateLimiter().Reset("login:" + clientIP)
	}

	response.Success(c, result)
}

// GetStatus 获取管理员状态（是否需要修改密码）
func (h *AdminHandler) GetStatus(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	result, err := h.adminService.GetStatus(username)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// ChangePassword 修改密码（需要验证旧密码）
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效: "+err.Error())
		return
	}

	if err := h.adminService.ChangePassword(username, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "密码修改成功"})
}

// ForceChangePassword 首次登录强制修改密码
func (h *AdminHandler) ForceChangePassword(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req service.ForceChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数无效: "+err.Error())
		return
	}

	if err := h.adminService.ForceChangePassword(username, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "密码设置成功"})
}
