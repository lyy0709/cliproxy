/*
 * 文件作用：网关数据访问层，提供 Gateway 的 CRUD 操作
 * 负责功能：
 *   - 网关创建、查询、更新、删除
 *   - 获取启用的网关列表
 *   - 默认网关管理
 *   - 测试状态更新
 * 重要程度：⭐⭐⭐ 重要（xyrt 授权核心）
 * 依赖模块：model, gorm
 */
package repository

import (
	"time"

	"cli-proxy/internal/model"

	"gorm.io/gorm"
)

// GatewayRepository 网关数据访问层
type GatewayRepository struct {
	db *gorm.DB
}

// NewGatewayRepository 创建网关仓库实例
func NewGatewayRepository() *GatewayRepository {
	return &GatewayRepository{db: DB}
}

// Create 创建网关
func (r *GatewayRepository) Create(gateway *model.Gateway) error {
	return r.db.Create(gateway).Error
}

// GetByID 根据 ID 获取网关
func (r *GatewayRepository) GetByID(id uint) (*model.Gateway, error) {
	var gateway model.Gateway
	err := r.db.First(&gateway, id).Error
	if err != nil {
		return nil, err
	}
	return &gateway, nil
}

// Update 更新网关
func (r *GatewayRepository) Update(gateway *model.Gateway) error {
	return r.db.Save(gateway).Error
}

// Delete 删除网关（软删除）
func (r *GatewayRepository) Delete(id uint) error {
	return r.db.Delete(&model.Gateway{}, id).Error
}

// List 分页获取网关列表
func (r *GatewayRepository) List(page, pageSize int, search string) ([]model.Gateway, int64, error) {
	var gateways []model.Gateway
	var total int64

	query := r.db.Model(&model.Gateway{})

	// 搜索条件
	if search != "" {
		query = query.Where("name LIKE ? OR url LIKE ? OR remark LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&gateways).Error; err != nil {
		return nil, 0, err
	}

	return gateways, total, nil
}

// GetEnabled 获取所有启用的网关
func (r *GatewayRepository) GetEnabled() ([]model.Gateway, error) {
	var gateways []model.Gateway
	err := r.db.Where("enabled = ?", true).Order("is_default DESC, name ASC").Find(&gateways).Error
	return gateways, err
}

// GetDefault 获取默认网关
func (r *GatewayRepository) GetDefault() (*model.Gateway, error) {
	var gateway model.Gateway
	err := r.db.Where("is_default = ? AND enabled = ?", true, true).First(&gateway).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &gateway, nil
}

// SetDefault 设置默认网关
func (r *GatewayRepository) SetDefault(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先清除所有默认标记
		if err := tx.Model(&model.Gateway{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置新的默认
		return tx.Model(&model.Gateway{}).Where("id = ?", id).Update("is_default", true).Error
	})
}

// ClearDefault 清除默认网关
func (r *GatewayRepository) ClearDefault() error {
	return r.db.Model(&model.Gateway{}).Where("is_default = ?", true).Update("is_default", false).Error
}

// ToggleEnabled 切换启用状态
func (r *GatewayRepository) ToggleEnabled(id uint) error {
	return r.db.Model(&model.Gateway{}).Where("id = ?", id).
		Update("enabled", gorm.Expr("NOT enabled")).Error
}

// UpdateTestStatus 更新测试状态
func (r *GatewayRepository) UpdateTestStatus(id uint, status string, latency int, errMsg string) error {
	now := time.Now()
	return r.db.Model(&model.Gateway{}).Where("id = ?", id).Updates(map[string]interface{}{
		"test_status":  status,
		"test_latency": latency,
		"test_error":   errMsg,
		"last_test_at": &now,
	}).Error
}

// CountByIDs 统计使用指定网关的账户数量
func (r *GatewayRepository) CountAccountsUsingGateway(gatewayID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Account{}).Where("gateway_id = ?", gatewayID).Count(&count).Error
	return count, err
}

// GetAll 获取所有网关（不分页）
func (r *GatewayRepository) GetAll() ([]model.Gateway, error) {
	var gateways []model.Gateway
	err := r.db.Order("created_at DESC").Find(&gateways).Error
	return gateways, err
}
