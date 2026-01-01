/*
 * 文件作用：API Key 权限校验辅助方法
 * 负责功能：
 *   - 平台访问限制校验
 *   - 模型允许/禁止列表校验
 * 重要程度：⭐⭐⭐ 一般（权限辅助）
 */
package handler

import (
	"strings"

	"cli-proxy/internal/middleware"
	"cli-proxy/internal/model"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

func enforceAPIKeyAccess(c *gin.Context, platform, modelName string) bool {
	apiKey := middleware.GetAPIKey(c)
	if apiKey == nil {
		response.CustomUnauthorizedAbort(c, model.ErrorTypeAuthFailed, "缺少 API Key")
		return false
	}

	if platform != "" && !isPlatformAllowed(apiKey.AllowedPlatforms, platform) {
		response.CustomForbiddenAbort(c, model.ErrorTypePlatformForbid, "平台访问受限")
		return false
	}

	if modelName != "" {
		if isModelBlocked(apiKey.BlockedModels, modelName) {
			response.CustomForbiddenAbort(c, model.ErrorTypeModelForbidden, "模型访问受限")
			return false
		}
		if !isModelAllowed(apiKey.AllowedModels, modelName) {
			response.CustomForbiddenAbort(c, model.ErrorTypeModelForbidden, "模型访问受限")
			return false
		}
	}

	return true
}

func isPlatformAllowed(allowedPlatforms, platform string) bool {
	allowed := strings.TrimSpace(strings.ToLower(allowedPlatforms))
	if allowed == "" || allowed == "all" {
		return true
	}
	target := strings.TrimSpace(strings.ToLower(platform))
	for _, p := range strings.Split(allowed, ",") {
		p = strings.TrimSpace(strings.ToLower(p))
		if p == "" {
			continue
		}
		if p == target {
			return true
		}
	}
	return false
}

func isModelAllowed(allowedModels, modelName string) bool {
	allowed := strings.TrimSpace(strings.ToLower(allowedModels))
	if allowed == "" || allowed == "all" {
		return true
	}
	return matchModelList(allowedModels, modelName)
}

func isModelBlocked(blockedModels, modelName string) bool {
	blocked := strings.TrimSpace(strings.ToLower(blockedModels))
	if blocked == "" {
		return false
	}
	if blocked == "all" {
		return true
	}
	return matchModelList(blockedModels, modelName)
}

func matchModelList(list, modelName string) bool {
	modelLower := strings.TrimSpace(strings.ToLower(modelName))
	if modelLower == "" {
		return false
	}
	for _, item := range strings.Split(list, ",") {
		item = strings.TrimSpace(strings.ToLower(item))
		if item == "" {
			continue
		}
		if modelLower == item || strings.HasPrefix(modelLower, item) {
			return true
		}
	}
	return false
}
