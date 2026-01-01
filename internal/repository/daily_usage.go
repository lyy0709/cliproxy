/*
 * 文件作用：每日使用汇总数据仓库，提供使用统计的数据库操作
 * 负责功能：
 *   - 增量更新每日使用统计（UPSERT）
 *   - API Key使用统计查询（按日/月/总计）
 *   - 模型使用统计汇总
 *   - 管理员全局统计查询
 * 重要程度：⭐⭐⭐⭐ 重要（使用统计核心仓库）
 * 依赖模块：model, gorm
 */
package repository

import (
	"time"

	"cli-proxy/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// DailyUsageRepository 每日使用汇总仓库
type DailyUsageRepository struct {
	db *gorm.DB
}

// NewDailyUsageRepository 创建每日使用汇总仓库
func NewDailyUsageRepository() *DailyUsageRepository {
	return &DailyUsageRepository{db: GetDB()}
}

// IncrementUsage 增量更新每日使用统计（使用 UPSERT）
func (r *DailyUsageRepository) IncrementUsage(apiKeyID uint, modelName string, usage *model.DailyUsage) error {
	today := time.Now().Format("2006-01-02")

	// 使用 ON DUPLICATE KEY UPDATE 实现增量更新
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "api_key_id"},
			{Name: "date"},
			{Name: "model"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"request_count":               gorm.Expr("request_count + ?", usage.RequestCount),
			"input_tokens":                gorm.Expr("input_tokens + ?", usage.InputTokens),
			"output_tokens":               gorm.Expr("output_tokens + ?", usage.OutputTokens),
			"cache_creation_input_tokens": gorm.Expr("cache_creation_input_tokens + ?", usage.CacheCreationInputTokens),
			"cache_read_input_tokens":     gorm.Expr("cache_read_input_tokens + ?", usage.CacheReadInputTokens),
			"total_tokens":                gorm.Expr("total_tokens + ?", usage.TotalTokens),
			"input_cost":                  gorm.Expr("input_cost + ?", usage.InputCost),
			"output_cost":                 gorm.Expr("output_cost + ?", usage.OutputCost),
			"cache_create_cost":           gorm.Expr("cache_create_cost + ?", usage.CacheCreateCost),
			"cache_read_cost":             gorm.Expr("cache_read_cost + ?", usage.CacheReadCost),
			"total_cost":                  gorm.Expr("total_cost + ?", usage.TotalCost),
			"updated_at":                  time.Now(),
		}),
	}).Create(&model.DailyUsage{
		APIKeyID:                 apiKeyID,
		Date:                     today,
		Model:                    modelName,
		RequestCount:             usage.RequestCount,
		InputTokens:              usage.InputTokens,
		OutputTokens:             usage.OutputTokens,
		CacheCreationInputTokens: usage.CacheCreationInputTokens,
		CacheReadInputTokens:     usage.CacheReadInputTokens,
		TotalTokens:              usage.TotalTokens,
		InputCost:                usage.InputCost,
		OutputCost:               usage.OutputCost,
		CacheCreateCost:          usage.CacheCreateCost,
		CacheReadCost:            usage.CacheReadCost,
		TotalCost:                usage.TotalCost,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}).Error
}

// GetAPIKeyDailyUsage 获取 API Key 某日的使用统计
func (r *DailyUsageRepository) GetAPIKeyDailyUsage(apiKeyID uint, date string) ([]model.DailyUsage, error) {
	var usages []model.DailyUsage
	err := r.db.Where("api_key_id = ? AND date = ?", apiKeyID, date).Find(&usages).Error
	return usages, err
}

// GetAPIKeyDateRangeUsage 获取 API Key 日期范围内的使用统计
func (r *DailyUsageRepository) GetAPIKeyDateRangeUsage(apiKeyID uint, startDate, endDate string) ([]model.DailyUsage, error) {
	var usages []model.DailyUsage
	err := r.db.Where("api_key_id = ? AND date >= ? AND date <= ?", apiKeyID, startDate, endDate).
		Order("date DESC, model ASC").
		Find(&usages).Error
	return usages, err
}

// GetAPIKeyDailySummary 获取 API Key 每日汇总（不分模型）
func (r *DailyUsageRepository) GetAPIKeyDailySummary(apiKeyID uint, startDate, endDate string) ([]model.DailyUsageSummary, error) {
	var summaries []model.DailyUsageSummary
	err := r.db.Model(&model.DailyUsage{}).
		Select("date, SUM(request_count) as request_count, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Where("api_key_id = ? AND date >= ? AND date <= ?", apiKeyID, startDate, endDate).
		Group("date").
		Order("date DESC").
		Scan(&summaries).Error
	return summaries, err
}

// GetAPIKeyModelSummary 获取 API Key 模型使用汇总
func (r *DailyUsageRepository) GetAPIKeyModelSummary(apiKeyID uint, startDate, endDate string) ([]model.ModelUsageSummary, error) {
	var summaries []model.ModelUsageSummary
	err := r.db.Model(&model.DailyUsage{}).
		Select("model, SUM(request_count) as request_count, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Where("api_key_id = ? AND date >= ? AND date <= ?", apiKeyID, startDate, endDate).
		Group("model").
		Order("total_cost DESC").
		Scan(&summaries).Error
	return summaries, err
}

// GetAPIKeyTotalUsage 获取 API Key 总使用汇总
func (r *DailyUsageRepository) GetAPIKeyTotalUsage(apiKeyID uint) (*model.APIKeyUsageSummary, error) {
	var summary model.APIKeyUsageSummary
	err := r.db.Model(&model.DailyUsage{}).
		Select("api_key_id, SUM(request_count) as total_requests, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Where("api_key_id = ?", apiKeyID).
		Group("api_key_id").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	summary.APIKeyID = apiKeyID
	return &summary, nil
}

// GetAPIKeyMonthlyUsage 获取 API Key 某月使用汇总
func (r *DailyUsageRepository) GetAPIKeyMonthlyUsage(apiKeyID uint, year, month int) (*model.APIKeyUsageSummary, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	endDate := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02")

	var summary model.APIKeyUsageSummary
	err := r.db.Model(&model.DailyUsage{}).
		Select("api_key_id, SUM(request_count) as total_requests, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Where("api_key_id = ? AND date >= ? AND date <= ?", apiKeyID, startDate, endDate).
		Group("api_key_id").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	summary.APIKeyID = apiKeyID
	return &summary, nil
}

// GetTodayUsage 获取 API Key 今日使用汇总
func (r *DailyUsageRepository) GetTodayUsage(apiKeyID uint) (*model.APIKeyUsageSummary, error) {
	today := time.Now().Format("2006-01-02")
	var summary model.APIKeyUsageSummary
	err := r.db.Model(&model.DailyUsage{}).
		Select("api_key_id, SUM(request_count) as total_requests, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Where("api_key_id = ? AND date = ?", apiKeyID, today).
		Group("api_key_id").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}
	summary.APIKeyID = apiKeyID
	return &summary, nil
}

// GetAllAPIKeysDailySummary 获取所有 API Key 某日汇总（管理员用）
func (r *DailyUsageRepository) GetAllAPIKeysDailySummary(date string) ([]model.APIKeyUsageSummary, error) {
	var summaries []model.APIKeyUsageSummary
	err := r.db.Model(&model.DailyUsage{}).
		Select("api_key_id, SUM(request_count) as total_requests, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Where("date = ?", date).
		Group("api_key_id").
		Order("total_cost DESC").
		Scan(&summaries).Error
	return summaries, err
}

// GetAllDailyUsageList 获取每日使用记录列表（管理员用，分页）
func (r *DailyUsageRepository) GetAllDailyUsageList(startDate, endDate string, offset, limit int64) ([]model.DailyUsage, int64, error) {
	var usages []model.DailyUsage
	var total int64

	query := r.db.Model(&model.DailyUsage{})
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("date DESC, total_cost DESC").
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&usages).Error

	return usages, total, err
}

// GetTotalSummary 获取总体汇总（管理员用）
func (r *DailyUsageRepository) GetTotalSummary(startDate, endDate string) (*model.APIKeyUsageSummary, error) {
	var summary model.APIKeyUsageSummary

	query := r.db.Model(&model.DailyUsage{})
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	err := query.Select("SUM(request_count) as total_requests, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Scan(&summary).Error

	return &summary, err
}

// GetModelSummary 获取模型使用汇总（管理员用）
func (r *DailyUsageRepository) GetModelSummary(startDate, endDate string) ([]model.ModelUsageSummary, error) {
	var summaries []model.ModelUsageSummary

	query := r.db.Model(&model.DailyUsage{})
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	err := query.Select("model, SUM(request_count) as request_count, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Group("model").
		Order("total_cost DESC").
		Scan(&summaries).Error

	return summaries, err
}

// GetDailySummaryAll 获取每日汇总（管理员用，所有 API Key 合计）
func (r *DailyUsageRepository) GetDailySummaryAll(startDate, endDate string) ([]model.DailyUsageSummary, error) {
	var summaries []model.DailyUsageSummary

	query := r.db.Model(&model.DailyUsage{})
	if startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	err := query.Select("date, SUM(request_count) as request_count, SUM(total_tokens) as total_tokens, SUM(total_cost) as total_cost").
		Group("date").
		Order("date DESC").
		Scan(&summaries).Error

	return summaries, err
}
