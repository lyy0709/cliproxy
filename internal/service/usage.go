/*
 * 文件作用：使用统计服务，处理Token使用量和费用统计
 * 负责功能：
 *   - API Key 使用量统计
 *   - 每日/月度使用汇总
 *   - 按模型使用统计
 *   - 使用记录写入
 *   - 账户费用统计
 * 重要程度：⭐⭐⭐⭐ 重要（计费统计核心）
 * 依赖模块：repository, model
 */
package service

import (
	"context"
	"fmt"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/repository"
)

// UsageService 使用统计服务（直接写入 MySQL）
type UsageService struct {
	dailyUsageRepo  *repository.DailyUsageRepository
	usageRecordRepo *repository.UsageRecordRepository
	accountRepo     *repository.AccountRepository
}

func NewUsageService() *UsageService {
	return &UsageService{
		dailyUsageRepo:  repository.NewDailyUsageRepository(),
		usageRecordRepo: repository.NewUsageRecordRepository(),
		accountRepo:     repository.NewAccountRepository(),
	}
}

// UsageData Token 使用数据
type UsageData struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
	TotalTokens              int64 `json:"total_tokens"`
	Requests                 int64 `json:"requests"`
}

// CostData 费用数据
type CostData struct {
	InputCost       float64 `json:"input_cost"`
	OutputCost      float64 `json:"output_cost"`
	CacheCreateCost float64 `json:"cache_create_cost"`
	CacheReadCost   float64 `json:"cache_read_cost"`
	TotalCost       float64 `json:"total_cost"`
}

// UsageRecord 单次请求的使用记录（用于 API 返回）
type UsageRecord struct {
	ID                       string    `json:"id"`
	Model                    string    `json:"model"`
	RequestIP                string    `json:"request_ip"`
	InputTokens              int       `json:"input_tokens"`
	OutputTokens             int       `json:"output_tokens"`
	CacheCreationInputTokens int       `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int       `json:"cache_read_input_tokens"`
	TotalTokens              int       `json:"total_tokens"`
	TotalCost                float64   `json:"total_cost"`
	Timestamp                time.Time `json:"timestamp"`
}

// RecordRequest 记录一次请求（综合方法，同时更新使用量、费用和记录）
// 直接写入 MySQL，不再使用 Redis
func (s *UsageService) RecordRequest(ctx context.Context, apiKeyID uint, log *model.RequestLog, priceRate float64) error {
	now := time.Now()

	// 1. 更新每日使用统计（UPSERT 到 daily_usage 表）
	dailyUsage := &model.DailyUsage{
		RequestCount:             1,
		InputTokens:              int64(log.InputTokens),
		OutputTokens:             int64(log.OutputTokens),
		CacheCreationInputTokens: int64(log.CacheCreationInputTokens),
		CacheReadInputTokens:     int64(log.CacheReadInputTokens),
		TotalTokens:              int64(log.TotalTokens),
		InputCost:                log.InputCost,
		OutputCost:               log.OutputCost,
		CacheCreateCost:          log.CacheCreateCost,
		CacheReadCost:            log.CacheReadCost,
		TotalCost:                log.TotalCost,
	}

	if err := s.dailyUsageRepo.IncrementUsage(apiKeyID, log.Model, dailyUsage); err != nil {
		return err
	}

	// 2. 创建使用记录（INSERT 到 usage_records 表）
	record := &model.UsageRecord{
		APIKeyID:                 apiKeyID,
		Model:                    log.Model,
		Platform:                 log.Platform,
		RequestIP:                log.RequestIP,
		InputTokens:              log.InputTokens,
		OutputTokens:             log.OutputTokens,
		CacheCreationInputTokens: log.CacheCreationInputTokens,
		CacheReadInputTokens:     log.CacheReadInputTokens,
		TotalTokens:              log.TotalTokens,
		TotalCost:                log.TotalCost,
		RequestTime:              now,
	}

	if err := s.usageRecordRepo.Create(record); err != nil {
		return err
	}

	return nil
}

// IncrementAccountCost 增加账户费用（直接更新 MySQL accounts 表）
func (s *UsageService) IncrementAccountCost(ctx context.Context, accountID uint, cost float64) error {
	if accountID == 0 {
		return nil
	}
	return s.accountRepo.IncrementTotalCost(accountID, cost)
}

// GetAccountCost 获取账户总费用
func (s *UsageService) GetAccountCost(ctx context.Context, accountID uint) (float64, error) {
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return 0, err
	}
	return account.TotalCost, nil
}

// GetAccountsCost 批量获取账户总费用
func (s *UsageService) GetAccountsCost(ctx context.Context, accountIDs []uint) (map[uint]float64, error) {
	return s.accountRepo.GetTotalCostByIDs(accountIDs)
}

// GetAccountDailyCost 获取账户某天的费用（从 request_logs 聚合）
func (s *UsageService) GetAccountDailyCost(ctx context.Context, accountID uint, date string) (float64, error) {
	startTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, err
	}
	endTime := startTime.Add(24 * time.Hour)

	var result struct {
		TotalCost float64
	}
	err = repository.DB.Model(&model.RequestLog{}).
		Select("COALESCE(SUM(total_cost), 0) as total_cost").
		Where("account_id = ? AND created_at >= ? AND created_at < ?", accountID, startTime, endTime).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.TotalCost, nil
}

// ==================== API Key 使用统计 ====================

// GetAPIKeyUsage 获取 API Key 的总使用量
func (s *UsageService) GetAPIKeyUsage(ctx context.Context, apiKeyID uint) (*UsageData, error) {
	summary, err := s.dailyUsageRepo.GetAPIKeyTotalUsage(apiKeyID)
	if err != nil {
		return nil, err
	}
	return &UsageData{
		Requests:    summary.TotalRequests,
		TotalTokens: summary.TotalTokens,
	}, nil
}

// GetAPIKeyDailyUsage 获取 API Key 某天的使用量
func (s *UsageService) GetAPIKeyDailyUsage(ctx context.Context, apiKeyID uint, date string) (*UsageData, error) {
	usages, err := s.dailyUsageRepo.GetAPIKeyDailyUsage(apiKeyID, date)
	if err != nil {
		return nil, err
	}

	data := &UsageData{}
	for _, u := range usages {
		data.InputTokens += u.InputTokens
		data.OutputTokens += u.OutputTokens
		data.CacheCreationInputTokens += u.CacheCreationInputTokens
		data.CacheReadInputTokens += u.CacheReadInputTokens
		data.TotalTokens += u.TotalTokens
		data.Requests += u.RequestCount
	}
	return data, nil
}

// GetAPIKeyMonthlyUsage 获取 API Key 某月的使用量
func (s *UsageService) GetAPIKeyMonthlyUsage(ctx context.Context, apiKeyID uint, month string) (*UsageData, error) {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, err
	}

	summary, err := s.dailyUsageRepo.GetAPIKeyMonthlyUsage(apiKeyID, t.Year(), int(t.Month()))
	if err != nil {
		return nil, err
	}

	return &UsageData{
		Requests:    summary.TotalRequests,
		TotalTokens: summary.TotalTokens,
	}, nil
}

// GetAPIKeyCost 获取 API Key 的总费用
func (s *UsageService) GetAPIKeyCost(ctx context.Context, apiKeyID uint) (*CostData, error) {
	summary, err := s.dailyUsageRepo.GetAPIKeyTotalUsage(apiKeyID)
	if err != nil {
		return nil, err
	}
	return &CostData{
		TotalCost: summary.TotalCost,
	}, nil
}

// GetAPIKeyDailyCost 获取 API Key 某天的费用
func (s *UsageService) GetAPIKeyDailyCost(ctx context.Context, apiKeyID uint, date string) (*CostData, error) {
	usages, err := s.dailyUsageRepo.GetAPIKeyDailyUsage(apiKeyID, date)
	if err != nil {
		return nil, err
	}

	data := &CostData{}
	for _, u := range usages {
		data.InputCost += u.InputCost
		data.OutputCost += u.OutputCost
		data.CacheCreateCost += u.CacheCreateCost
		data.CacheReadCost += u.CacheReadCost
		data.TotalCost += u.TotalCost
	}
	return data, nil
}

// GetAPIKeyMonthlyCost 获取 API Key 某月的费用
func (s *UsageService) GetAPIKeyMonthlyCost(ctx context.Context, apiKeyID uint, month string) (*CostData, error) {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, err
	}

	summary, err := s.dailyUsageRepo.GetAPIKeyMonthlyUsage(apiKeyID, t.Year(), int(t.Month()))
	if err != nil {
		return nil, err
	}

	return &CostData{
		TotalCost: summary.TotalCost,
	}, nil
}

// GetAPIKeyRecords 获取 API Key 使用记录
func (s *UsageService) GetAPIKeyRecords(ctx context.Context, apiKeyID uint, offset, limit int64) ([]UsageRecord, error) {
	records, _, err := s.usageRecordRepo.GetByAPIKeyID(apiKeyID, int(offset), int(limit))
	if err != nil {
		return nil, err
	}
	return s.convertRecords(records), nil
}

// GetAPIKeyRecordsWithFilters 获取 API Key 使用记录（带筛选）
func (s *UsageService) GetAPIKeyRecordsWithFilters(ctx context.Context, apiKeyID uint, offset, limit int, startDate, endDate, modelFilter string) ([]UsageRecord, int64, error) {
	records, total, err := s.usageRecordRepo.GetByAPIKeyIDWithFilters(apiKeyID, offset, limit, startDate, endDate, modelFilter)
	if err != nil {
		return nil, 0, err
	}
	return s.convertRecords(records), total, nil
}

// convertRecords 转换记录格式
func (s *UsageService) convertRecords(records []model.UsageRecord) []UsageRecord {
	result := make([]UsageRecord, len(records))
	for i, r := range records {
		result[i] = UsageRecord{
			ID:                       fmt.Sprintf("%d", r.ID),
			Model:                    r.Model,
			RequestIP:                r.RequestIP,
			InputTokens:              r.InputTokens,
			OutputTokens:             r.OutputTokens,
			CacheCreationInputTokens: r.CacheCreationInputTokens,
			CacheReadInputTokens:     r.CacheReadInputTokens,
			TotalTokens:              r.TotalTokens,
			TotalCost:                r.TotalCost,
			Timestamp:                r.RequestTime,
		}
	}
	return result
}

// GetAPIKeyDailyStats 获取 API Key 指定日期范围的每日统计
func (s *UsageService) GetAPIKeyDailyStats(ctx context.Context, apiKeyID uint, startDate, endDate string) ([]model.DailyUsageStats, error) {
	summaries, err := s.dailyUsageRepo.GetAPIKeyDailySummary(apiKeyID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	stats := make([]model.DailyUsageStats, len(summaries))
	for i, s := range summaries {
		stats[i] = model.DailyUsageStats{
			Date:         s.Date,
			RequestCount: s.RequestCount,
			TotalTokens:  s.TotalTokens,
			TotalCost:    s.TotalCost,
		}
	}
	return stats, nil
}

// GetAPIKeyModelStats 获取 API Key 按模型的统计
func (s *UsageService) GetAPIKeyModelStats(ctx context.Context, apiKeyID uint) ([]model.ModelUsageStats, error) {
	summaries, err := s.dailyUsageRepo.GetAPIKeyModelSummary(apiKeyID, "2000-01-01", "2100-12-31")
	if err != nil {
		return nil, err
	}

	stats := make([]model.ModelUsageStats, len(summaries))
	for i, s := range summaries {
		stats[i] = model.ModelUsageStats{
			Model:        s.Model,
			RequestCount: s.RequestCount,
			TotalTokens:  s.TotalTokens,
			TotalCost:    s.TotalCost,
		}
	}
	return stats, nil
}

// GetAPIKeyTodayUsage 获取 API Key 今日使用汇总
func (s *UsageService) GetAPIKeyTodayUsage(ctx context.Context, apiKeyID uint) (*UsageData, error) {
	summary, err := s.dailyUsageRepo.GetTodayUsage(apiKeyID)
	if err != nil {
		return nil, err
	}
	return &UsageData{
		Requests:    summary.TotalRequests,
		TotalTokens: summary.TotalTokens,
	}, nil
}

// ==================== 管理员统计方法 ====================

// GetAllAPIKeysDailySummary 获取所有 API Key 某日汇总（管理员用）
func (s *UsageService) GetAllAPIKeysDailySummary(ctx context.Context, date string) ([]model.APIKeyUsageSummary, error) {
	return s.dailyUsageRepo.GetAllAPIKeysDailySummary(date)
}

// GetTotalSummary 获取总体汇总（管理员用）
func (s *UsageService) GetTotalSummary(ctx context.Context, startDate, endDate string) (*model.APIKeyUsageSummary, error) {
	return s.dailyUsageRepo.GetTotalSummary(startDate, endDate)
}

// GetModelSummary 获取模型使用汇总（管理员用）
func (s *UsageService) GetModelSummary(ctx context.Context, startDate, endDate string) ([]model.ModelUsageSummary, error) {
	return s.dailyUsageRepo.GetModelSummary(startDate, endDate)
}

// GetDailySummaryAll 获取每日汇总（管理员用）
func (s *UsageService) GetDailySummaryAll(ctx context.Context, startDate, endDate string) ([]model.DailyUsageSummary, error) {
	return s.dailyUsageRepo.GetDailySummaryAll(startDate, endDate)
}

// GetAllRecords 获取所有使用记录（管理员用）
func (s *UsageService) GetAllRecords(ctx context.Context, offset, limit int, startDate, endDate, modelFilter string) ([]UsageRecord, int64, error) {
	records, total, err := s.usageRecordRepo.GetAllRecords(offset, limit, startDate, endDate, modelFilter)
	if err != nil {
		return nil, 0, err
	}
	return s.convertRecords(records), total, nil
}
