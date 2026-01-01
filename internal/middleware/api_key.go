/*
 * 文件作用：API Key认证中间件，验证代理请求的API密钥
 * 负责功能：
 *   - API Key 解析（支持多种Header格式）
 *   - API Key 有效性验证
 *   - API Key 信息注入上下文
 *   - 请求日志记录
 * 重要程度：⭐⭐⭐⭐⭐ 核心（代理认证核心）
 * 依赖模块：service, model
 */
package middleware

import (
	"strconv"
	"strings"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/service"
	"cli-proxy/pkg/logger"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth API Key 认证中间件
func APIKeyAuth() gin.HandlerFunc {
	apiKeyService := service.NewAPIKeyService()
	configService := service.GetConfigService()
	usageService := service.NewUsageService()
	apiKeyRateLimiter := service.GetAPIKeyRateLimiter()
	log := logger.GetLogger("auth")

	return func(c *gin.Context) {
		// 从 Header 获取 API Key（支持多种格式）
		apiKey := c.GetHeader("Authorization")
		if apiKey == "" {
			apiKey = c.GetHeader("X-API-Key")
		}
		if apiKey == "" {
			apiKey = c.GetHeader("x-api-key") // Claude 标准格式
		}

		if apiKey == "" {
			log.Debug("API Key 认证失败 | IP: %s | 原因: 缺少API Key", c.ClientIP())
			response.CustomUnauthorizedAbort(c, model.ErrorTypeAuthFailed, "缺少 API Key，请在 Authorization 或 x-api-key header 中提供")
			return
		}

		// 移除 Bearer 前缀
		if strings.HasPrefix(apiKey, "Bearer ") {
			apiKey = apiKey[7:]
		}

		// 验证 API Key
		key, err := apiKeyService.ValidateKey(apiKey)
		if err != nil {
			log.Debug("API Key 认证失败 | IP: %s | Key: %s... | 原因: %v", c.ClientIP(), maskAPIKey(apiKey), err)
			// 根据错误内容确定错误类型
			errorType := getAPIKeyErrorType(err.Error())
			response.CustomUnauthorizedAbort(c, errorType, err.Error())
			return
		}

		// 限流/额度检查（跳过自助查询等非代理接口）
		if shouldEnforceAPIKeyLimits(c) {
			// 频率限制（每分钟请求数）
			if key.RateLimit > 0 {
				allowed, waitSeconds := apiKeyRateLimiter.Check("apikey:"+strconv.FormatUint(uint64(key.ID), 10), key.RateLimit, 1)
				if !allowed {
					if waitSeconds > 0 {
						c.Header("Retry-After", strconv.Itoa(waitSeconds))
					}
					response.CustomTooManyRequestsAbort(c, model.ErrorTypeRateLimit, service.GetRateLimitError(waitSeconds))
					return
				}
			}

			// 每日请求限制
			if key.DailyLimit > 0 {
				todayUsage, err := usageService.GetAPIKeyTodayUsage(c.Request.Context(), key.ID)
				if err != nil {
					log.Warn("API Key 日限额检查失败 | KeyID: %d | Err: %v", key.ID, err)
				} else if todayUsage.Requests >= int64(key.DailyLimit) {
					response.CustomForbiddenAbort(c, model.ErrorTypeDailyLimit, "今日额度已用完")
					return
				}
			}

			// 月度配额限制（按费用）
			if key.MonthlyQuota > 0 {
				month := time.Now().Format("2006-01")
				monthlyCost, err := usageService.GetAPIKeyMonthlyCost(c.Request.Context(), key.ID, month)
				if err != nil {
					log.Warn("API Key 月配额检查失败 | KeyID: %d | Err: %v", key.ID, err)
				} else if monthlyCost.TotalCost >= key.MonthlyQuota {
					response.CustomForbiddenAbort(c, model.ErrorTypeMonthlyQuota, "本月额度已用完")
					return
				}
			}
		}

		log.Debug("API Key 认证成功 | IP: %s | KeyID: %d", c.ClientIP(), key.ID)

		// 将 API Key 信息存储到 Context 中
		c.Set("api_key", key)
		c.Set("api_key_id", key.ID)
		c.Set("api_key_allowed_platforms", key.AllowedPlatforms)
		c.Set("api_key_allowed_models", key.AllowedModels)
		c.Set("api_key_rate_limit", key.RateLimit)

		// 使用全局价格倍率
		priceRate := configService.GetGlobalPriceRate()
		c.Set("api_key_price_rate", priceRate)

		// 调试日志：每次请求都输出倍率信息
		log.Info("API Key 认证 | KeyID: %d | GlobalRate: %.2f | Path: %s",
			key.ID, priceRate, c.Request.URL.Path)

		c.Next()
	}
}

func shouldEnforceAPIKeyLimits(c *gin.Context) bool {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/usage/me") {
		return false
	}
	if strings.HasPrefix(path, "/api/key/") {
		return false
	}
	return true
}

// maskAPIKey 遮蔽API Key用于日志
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:8]
}

// GetAPIKey 从 Context 获取 API Key 信息
func GetAPIKey(c *gin.Context) *model.APIKey {
	if key, exists := c.Get("api_key"); exists {
		return key.(*model.APIKey)
	}
	return nil
}

// CheckPlatformAccess 检查平台访问权限
func CheckPlatformAccess(c *gin.Context, platform string) bool {
	allowed, exists := c.Get("api_key_allowed_platforms")
	if !exists {
		return true // 没有设置 API Key，允许访问
	}

	allowedStr := allowed.(string)
	if allowedStr == "all" || allowedStr == "" {
		return true
	}

	// 检查是否在允许的平台列表中
	platforms := strings.Split(allowedStr, ",")
	for _, p := range platforms {
		if strings.TrimSpace(p) == platform {
			return true
		}
	}

	return false
}

// CheckModelAccess 检查模型访问权限
func CheckModelAccess(c *gin.Context, modelName string) bool {
	allowed, exists := c.Get("api_key_allowed_models")
	if !exists {
		return true
	}

	allowedStr := allowed.(string)
	if allowedStr == "" {
		return true // 空字符串表示允许所有模型
	}

	// 检查是否在允许的模型列表中
	models := strings.Split(allowedStr, ",")
	for _, m := range models {
		if strings.TrimSpace(m) == modelName {
			return true
		}
	}

	return false
}

// getAPIKeyErrorType 根据错误信息判断错误类型
func getAPIKeyErrorType(errMsg string) string {
	switch {
	case strings.Contains(errMsg, "禁用"):
		return model.ErrorTypeKeyDisabled
	case strings.Contains(errMsg, "过期"):
		return model.ErrorTypeKeyExpired
	case strings.Contains(errMsg, "无效"):
		return model.ErrorTypeKeyInvalid
	default:
		return model.ErrorTypeAuthFailed
	}
}
