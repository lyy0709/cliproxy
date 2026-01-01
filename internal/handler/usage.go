/*
 * 文件作用：使用统计处理器，处理管理员的使用量统计查询
 * 负责功能：
 *   - 总使用量汇总
 *   - 每日使用统计
 *   - 按模型使用统计
 *   - 使用记录列表查询
 *   - API Key 使用统计
 * 重要程度：⭐⭐⭐⭐ 重要（数据统计核心）
 * 依赖模块：service, repository
 */
package handler

import (
	"strconv"
	"time"

	"cli-proxy/internal/repository"
	"cli-proxy/internal/service"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

// UsageHandler 使用统计处理器
type UsageHandler struct {
	usageService    *service.UsageService
	pricingService  *service.PricingService
	dailyUsageRepo  *repository.DailyUsageRepository
	usageRecordRepo *repository.UsageRecordRepository
}

// NewUsageHandler 创建使用统计处理器
func NewUsageHandler() *UsageHandler {
	return &UsageHandler{
		usageService:    service.NewUsageService(),
		pricingService:  service.NewPricingService(),
		dailyUsageRepo:  repository.NewDailyUsageRepository(),
		usageRecordRepo: repository.NewUsageRecordRepository(),
	}
}

// ========== 管理员接口 ==========

// GetTotalUsageSummary 获取总使用量汇总
func (h *UsageHandler) GetTotalUsageSummary(c *gin.Context) {
	ctx := c.Request.Context()

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = "2000-01-01"
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// 获取总汇总
	summary, err := h.usageService.GetTotalSummary(ctx, startDate, endDate)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取今日汇总
	today := time.Now().Format("2006-01-02")
	todaySummary, err := h.dailyUsageRepo.GetTotalSummary(today, today)
	todayTokens := int64(0)
	todayRequests := int64(0)
	todayCost := float64(0)
	if err == nil && todaySummary != nil {
		todayTokens = todaySummary.TotalTokens
		todayRequests = todaySummary.TotalRequests
		todayCost = todaySummary.TotalCost
	}

	response.Success(c, gin.H{
		"total_requests": summary.TotalRequests,
		"total_tokens":   summary.TotalTokens,
		"total_cost":     summary.TotalCost,
		"today_requests": todayRequests,
		"today_tokens":   todayTokens,
		"today_cost":     todayCost,
	})
}

// GetDailySummary 获取每日使用汇总
func (h *UsageHandler) GetDailySummary(c *gin.Context) {
	ctx := c.Request.Context()

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		// 默认最近 30 天
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	summaries, err := h.usageService.GetDailySummaryAll(ctx, startDate, endDate)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"start_date": startDate,
		"end_date":   endDate,
		"daily":      summaries,
	})
}

// GetModelSummary 获取模型使用汇总
func (h *UsageHandler) GetModelSummary(c *gin.Context) {
	ctx := c.Request.Context()

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" {
		startDate = "2000-01-01"
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	summaries, err := h.usageService.GetModelSummary(ctx, startDate, endDate)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"models": summaries,
	})
}

// GetAPIKeyUsageSummary 获取各 API Key 的使用汇总
func (h *UsageHandler) GetAPIKeyUsageSummary(c *gin.Context) {
	ctx := c.Request.Context()

	date := c.DefaultQuery("date", time.Now().Format("2006-01-02"))

	summaries, err := h.usageService.GetAllAPIKeysDailySummary(ctx, date)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"date":     date,
		"api_keys": summaries,
	})
}

// GetAllUsageRecords 获取所有使用记录
func (h *UsageHandler) GetAllUsageRecords(c *gin.Context) {
	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	modelFilter := c.Query("model")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	records, total, err := h.usageService.GetAllRecords(ctx, offset, pageSize, startDate, endDate, modelFilter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items":     records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAPIKeyUsage 获取指定 API Key 的使用量
func (h *UsageHandler) GetAPIKeyUsage(c *gin.Context) {
	keyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	ctx := c.Request.Context()

	// 获取使用量
	usage, err := h.usageService.GetAPIKeyUsage(ctx, uint(keyID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取费用
	cost, err := h.usageService.GetAPIKeyCost(ctx, uint(keyID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取今日使用量
	today := time.Now().Format("2006-01-02")
	todayUsage, _ := h.usageService.GetAPIKeyDailyUsage(ctx, uint(keyID), today)
	todayCost, _ := h.usageService.GetAPIKeyDailyCost(ctx, uint(keyID), today)

	// 获取模型统计
	modelStats, _ := h.usageService.GetAPIKeyModelStats(ctx, uint(keyID))

	response.Success(c, gin.H{
		"api_key_id":                  keyID,
		"total_tokens":                usage.TotalTokens,
		"input_tokens":                usage.InputTokens,
		"output_tokens":               usage.OutputTokens,
		"cache_creation_input_tokens": usage.CacheCreationInputTokens,
		"cache_read_input_tokens":     usage.CacheReadInputTokens,
		"total_requests":              usage.Requests,
		"total_cost":                  cost.TotalCost,
		"today_tokens":                todayUsage.TotalTokens,
		"today_requests":              todayUsage.Requests,
		"today_cost":                  todayCost.TotalCost,
		"model_count":                 len(modelStats),
		"model_stats":                 modelStats,
	})
}

// GetModels 获取所有模型定价（用于展示）
func (h *UsageHandler) GetModels(c *gin.Context) {
	ctx := c.Request.Context()

	models, err := h.pricingService.GetAllModels(ctx)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"models": models,
	})
}

// ========== API Key 统计详情接口 ==========

// GetAPIKeyDailyStats 获取 API Key 的每日统计
func (h *UsageHandler) GetAPIKeyDailyStats(c *gin.Context) {
	keyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	ctx := c.Request.Context()

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	stats, err := h.usageService.GetAPIKeyDailyStats(ctx, uint(keyID), startDate, endDate)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"api_key_id": keyID,
		"start_date": startDate,
		"end_date":   endDate,
		"stats":      stats,
	})
}

// GetAPIKeyRecords 获取 API Key 的使用记录
func (h *UsageHandler) GetAPIKeyRecords(c *gin.Context) {
	keyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	modelFilter := c.Query("model")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	records, total, err := h.usageService.GetAPIKeyRecordsWithFilters(ctx, uint(keyID), offset, pageSize, startDate, endDate, modelFilter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"api_key_id": keyID,
		"items":      records,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
	})
}

// ========== API Key 持有者自查接口（使用 x-api-key 认证） ==========

// GetMyUsage 获取当前 API Key 的使用量（API Key 持有者自查）
func (h *UsageHandler) GetMyUsage(c *gin.Context) {
	// 从上下文获取 API Key ID（由 APIKeyAuth 中间件设置）
	apiKeyIDRaw, exists := c.Get("api_key_id")
	if !exists {
		response.Unauthorized(c, "API Key not found in context")
		return
	}
	keyID := apiKeyIDRaw.(uint)

	ctx := c.Request.Context()

	// 获取使用量
	usage, err := h.usageService.GetAPIKeyUsage(ctx, keyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取费用
	cost, err := h.usageService.GetAPIKeyCost(ctx, keyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 获取今日使用量
	today := time.Now().Format("2006-01-02")
	todayUsage, _ := h.usageService.GetAPIKeyDailyUsage(ctx, keyID, today)
	todayCost, _ := h.usageService.GetAPIKeyDailyCost(ctx, keyID, today)

	// 获取模型统计
	modelStats, _ := h.usageService.GetAPIKeyModelStats(ctx, keyID)

	// 获取 API Key 信息（名称、限额等）
	apiKeyRepo := repository.NewAPIKeyRepository()
	apiKey, _ := apiKeyRepo.GetByID(keyID)

	result := gin.H{
		"total_tokens":                usage.TotalTokens,
		"input_tokens":                usage.InputTokens,
		"output_tokens":               usage.OutputTokens,
		"cache_creation_input_tokens": usage.CacheCreationInputTokens,
		"cache_read_input_tokens":     usage.CacheReadInputTokens,
		"total_requests":              usage.Requests,
		"total_cost":                  cost.TotalCost,
		"today_tokens":                todayUsage.TotalTokens,
		"today_requests":              todayUsage.Requests,
		"today_cost":                  todayCost.TotalCost,
		"model_count":                 len(modelStats),
		"model_stats":                 modelStats,
	}

	// 添加限额信息（如果存在）
	if apiKey != nil {
		result["key_name"] = apiKey.Name
		result["daily_limit"] = apiKey.DailyLimit
		result["monthly_quota"] = apiKey.MonthlyQuota
		if apiKey.ExpiresAt != nil {
			result["expires_at"] = apiKey.ExpiresAt
		}
	}

	response.Success(c, result)
}

// GetMyDailyStats 获取当前 API Key 的每日统计（API Key 持有者自查）
func (h *UsageHandler) GetMyDailyStats(c *gin.Context) {
	apiKeyIDRaw, exists := c.Get("api_key_id")
	if !exists {
		response.Unauthorized(c, "API Key not found in context")
		return
	}
	keyID := apiKeyIDRaw.(uint)

	ctx := c.Request.Context()

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		endDate = time.Now().Format("2006-01-02")
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}

	stats, err := h.usageService.GetAPIKeyDailyStats(ctx, keyID, startDate, endDate)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"start_date": startDate,
		"end_date":   endDate,
		"stats":      stats,
	})
}

// GetMyRecords 获取当前 API Key 的使用记录（API Key 持有者自查）
func (h *UsageHandler) GetMyRecords(c *gin.Context) {
	apiKeyIDRaw, exists := c.Get("api_key_id")
	if !exists {
		response.Unauthorized(c, "API Key not found in context")
		return
	}
	keyID := apiKeyIDRaw.(uint)

	ctx := c.Request.Context()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	modelFilter := c.Query("model")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	records, total, err := h.usageService.GetAPIKeyRecordsWithFilters(ctx, keyID, offset, pageSize, startDate, endDate, modelFilter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items":     records,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
