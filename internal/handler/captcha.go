/*
 * 文件作用：验证码处理器，处理验证码相关请求
 * 负责功能：
 *   - 生成验证码图片
 * 重要程度：⭐⭐⭐⭐ 重要（安全相关）
 * 依赖模块：service, response
 */
package handler

import (
	"cli-proxy/internal/service"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	captchaService *service.CaptchaService
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler() *CaptchaHandler {
	return &CaptchaHandler{
		captchaService: service.GetCaptchaService(),
	}
}

// Generate 生成验证码
func (h *CaptchaHandler) Generate(c *gin.Context) {
	configService := service.GetConfigService()
	if !configService.GetCaptchaEnabled() {
		response.Success(c, &service.CaptchaResponse{
			CaptchaId: "",
			Image:     "",
			Enabled:   false,
		})
		return
	}

	limit := configService.GetCaptchaRateLimit()
	if limit > 0 {
		allowed, waitSeconds := service.GetCaptchaRateLimiter().Check("captcha:"+c.ClientIP(), limit, 1)
		if !allowed {
			response.TooManyRequests(c, service.GetRateLimitError(waitSeconds))
			return
		}
	}

	result, err := h.captchaService.Generate()
	if err != nil {
		response.InternalError(c, "生成验证码失败")
		return
	}

	response.Success(c, result)
}
