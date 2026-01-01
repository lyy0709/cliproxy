/*
 * 文件作用：管理员配置数据模型
 * 负责功能：
 *   - 管理员密码存储（bcrypt加密）
 *   - 首次登录强制修改密码标志
 *   - 密码修改时间记录
 * 重要程度：⭐⭐⭐⭐ 重要（安全核心模型）
 */
package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AdminConfig 管理员配置
type AdminConfig struct {
	ID                 uint       `json:"id" gorm:"primaryKey"`
	Username           string     `json:"username" gorm:"uniqueIndex;size:50;not null"`
	PasswordHash       string     `json:"-" gorm:"size:255;not null"`
	MustChangePassword bool       `json:"must_change_password" gorm:"default:true"`
	PasswordChangedAt  *time.Time `json:"password_changed_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// TableName 指定表名
func (AdminConfig) TableName() string {
	return "admin_configs"
}

// SetPassword 设置密码（bcrypt加密）
func (a *AdminConfig) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.PasswordHash = string(hash)
	return nil
}

// VerifyPassword 验证密码
func (a *AdminConfig) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	return err == nil
}

// ValidatePasswordStrength 验证密码强度（至少8位，包含字母和数字）
func ValidatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasDigit := false

	for _, c := range password {
		switch {
		case 'a' <= c && c <= 'z':
			hasLetter = true
		case 'A' <= c && c <= 'Z':
			hasLetter = true
		case '0' <= c && c <= '9':
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}
