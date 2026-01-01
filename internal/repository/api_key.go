/*
 * 文件作用：API Key数据仓库，提供API Key的数据库操作
 * 负责功能：
 *   - API Key CRUD操作
 *   - 按哈希查询
 *   - 使用量统计更新
 *   - 使用日志查询
 * 重要程度：⭐⭐⭐⭐ 重要（API Key核心仓库）
 * 依赖模块：model, gorm
 */
package repository

import (
	"cli-proxy/internal/model"

	"gorm.io/gorm"
)

type APIKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository() *APIKeyRepository {
	return &APIKeyRepository{db: DB}
}

// Create 创建 API Key
func (r *APIKeyRepository) Create(key *model.APIKey) error {
	return r.db.Select("*").Create(key).Error
}

// GetByID 根据 ID 获取 API Key
func (r *APIKeyRepository) GetByID(id uint) (*model.APIKey, error) {
	var key model.APIKey
	err := r.db.First(&key, id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetByIDs 批量获取 API Key（按需字段，用于前端快速映射显示）
func (r *APIKeyRepository) GetByIDs(ids []uint) ([]model.APIKey, error) {
	if len(ids) == 0 {
		return []model.APIKey{}, nil
	}

	var keys []model.APIKey
	err := r.db.Select("id", "name", "key_prefix").
		Where("id IN ?", ids).
		Find(&keys).Error
	return keys, err
}

// GetByHash 根据哈希获取 API Key
func (r *APIKeyRepository) GetByHash(hash string) (*model.APIKey, error) {
	var key model.APIKey
	err := r.db.Where("key_hash = ?", hash).First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// List 获取所有 API Key（分页）
func (r *APIKeyRepository) List(page, pageSize int) ([]model.APIKey, int64, error) {
	var keys []model.APIKey
	var total int64

	err := r.db.Model(&model.APIKey{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&keys).Error
	return keys, total, err
}

// ListAll 获取所有 API Key（不分页）
func (r *APIKeyRepository) ListAll() ([]model.APIKey, error) {
	var keys []model.APIKey
	err := r.db.Order("created_at DESC").Find(&keys).Error
	return keys, err
}

// Update 更新 API Key
func (r *APIKeyRepository) Update(key *model.APIKey) error {
	return r.db.Save(key).Error
}

// Delete 删除 API Key
func (r *APIKeyRepository) Delete(id uint) error {
	return r.db.Delete(&model.APIKey{}, id).Error
}

// IncrementUsage 增加使用统计
func (r *APIKeyRepository) IncrementUsage(id uint, tokens int64, cost float64) error {
	return r.db.Model(&model.APIKey{}).Where("id = ?", id).Updates(map[string]interface{}{
		"request_count": gorm.Expr("request_count + 1"),
		"tokens_used":   gorm.Expr("tokens_used + ?", tokens),
		"cost_used":     gorm.Expr("cost_used + ?", cost),
		"last_used_at":  gorm.Expr("NOW()"),
	}).Error
}

// Count 统计 API Key 数量
func (r *APIKeyRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.APIKey{}).Count(&count).Error
	return count, err
}

// GetAPIKeyLogs 从 MySQL 获取 API Key 的使用日志
func (r *APIKeyRepository) GetAPIKeyLogs(keyID uint, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	r.db.Model(&model.UsageRecord{}).Where("api_key_id = ?", keyID).Count(&total)

	offset := (page - 1) * pageSize
	var records []model.UsageRecord
	err := r.db.Where("api_key_id = ?", keyID).
		Order("request_time DESC").
		Offset(offset).Limit(pageSize).
		Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	logs := make([]map[string]interface{}, 0, len(records))
	for _, rec := range records {
		logs = append(logs, map[string]interface{}{
			"id":                          rec.ID,
			"model":                       rec.Model,
			"platform":                    rec.Platform,
			"request_ip":                  rec.RequestIP,
			"input_tokens":                rec.InputTokens,
			"output_tokens":               rec.OutputTokens,
			"cache_creation_input_tokens": rec.CacheCreationInputTokens,
			"cache_read_input_tokens":     rec.CacheReadInputTokens,
			"total_tokens":                rec.TotalTokens,
			"total_cost":                  rec.TotalCost,
			"timestamp":                   rec.RequestTime,
		})
	}

	return logs, total, nil
}
