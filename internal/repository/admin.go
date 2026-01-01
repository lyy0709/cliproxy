/*
 * 文件作用：管理员配置仓库层
 * 负责功能：
 *   - 管理员配置的 CRUD 操作
 *   - 初始化默认管理员配置
 * 重要程度：⭐⭐⭐⭐ 重要（数据访问层）
 */
package repository

import (
	"cli-proxy/internal/model"

	"gorm.io/gorm"
)

// AdminConfigRepository 管理员配置仓库
type AdminConfigRepository struct {
	db *gorm.DB
}

// NewAdminConfigRepository 创建管理员配置仓库实例
func NewAdminConfigRepository(db *gorm.DB) *AdminConfigRepository {
	return &AdminConfigRepository{db: db}
}

// GetByUsername 根据用户名获取管理员配置
func (r *AdminConfigRepository) GetByUsername(username string) (*model.AdminConfig, error) {
	var config model.AdminConfig
	err := r.db.Where("username = ?", username).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetFirst 获取第一个管理员配置（系统只有一个管理员）
func (r *AdminConfigRepository) GetFirst() (*model.AdminConfig, error) {
	var config model.AdminConfig
	err := r.db.First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// Create 创建管理员配置
func (r *AdminConfigRepository) Create(config *model.AdminConfig) error {
	return r.db.Create(config).Error
}

// Update 更新管理员配置
func (r *AdminConfigRepository) Update(config *model.AdminConfig) error {
	return r.db.Save(config).Error
}

// UpdatePassword 更新密码
func (r *AdminConfigRepository) UpdatePassword(id uint, passwordHash string, mustChangePassword bool) error {
	return r.db.Model(&model.AdminConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password_hash":        passwordHash,
		"must_change_password": mustChangePassword,
		"password_changed_at":  gorm.Expr("NOW()"),
	}).Error
}

// Exists 检查管理员配置是否存在
func (r *AdminConfigRepository) Exists() (bool, error) {
	var count int64
	err := r.db.Model(&model.AdminConfig{}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InitializeDefaultAdmin 初始化默认管理员（密码: admin123，首次登录需修改）
func (r *AdminConfigRepository) InitializeDefaultAdmin() error {
	exists, err := r.Exists()
	if err != nil {
		return err
	}
	if exists {
		return nil // 已存在，不需要初始化
	}

	admin := &model.AdminConfig{
		Username:           "admin",
		MustChangePassword: true,
	}
	// 设置默认密码 admin123
	if err := admin.SetPassword("admin123"); err != nil {
		return err
	}

	return r.Create(admin)
}
