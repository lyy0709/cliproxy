/*
 * 文件作用：API Key数据模型，定义API密钥的数据结构
 * 负责功能：
 *   - API Key基础信息（名称、状态、描述）
 *   - Key哈希存储
 *   - 权限控制（平台、模型、客户端）
 *   - 限制配置（频率、每日限制、月额度）
 *   - Key生成和验证方法
 * 重要程度：⭐⭐⭐⭐ 重要（核心数据结构）
 * 依赖模块：gorm
 */
package model

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

const (
	APIKeyPrefix = "sk-" // API Key 前缀
)

// APIKey API 密钥模型
type APIKey struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`         // 名称
	Description string         `gorm:"size:500" json:"description"`           // 描述用途
	KeyHash     string         `gorm:"size:64;uniqueIndex;not null" json:"-"` // Key 的 SHA256 哈希
	KeyPrefix   string         `gorm:"size:20" json:"key_prefix"`             // Key 前缀用于显示 (如 sk-xxx...)
	Status      string         `gorm:"size:20;default:active" json:"status"`  // 状态: active, disabled, expired

	// 权限控制
	AllowedPlatforms string `gorm:"size:100;default:all" json:"allowed_platforms"` // 允许的平台: all, claude, openai, gemini (逗号分隔)
	AllowedModels    string `gorm:"type:text" json:"allowed_models,omitempty"`     // 允许的模型列表 (逗号分隔)
	BlockedModels    string `gorm:"type:text" json:"blocked_models,omitempty"`     // 禁止的模型列表 (逗号分隔)
	AllowedClients   string `gorm:"size:200" json:"allowed_clients,omitempty"`     // 允许的客户端类型 (逗号分隔, 如: claude_code,codex_cli)

	// 限制配置
	RateLimit     int        `gorm:"default:0" json:"rate_limit"`                // 每分钟请求限制（0=不限）
	DailyLimit    int        `gorm:"default:0" json:"daily_limit"`               // 每日请求限制 (0=不限)
	MonthlyQuota  float64    `gorm:"type:decimal(10,2);default:0" json:"monthly_quota"` // 月额度 (美元，0=不限)
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`                       // 过期时间

	// 统计字段
	RequestCount   int64      `gorm:"default:0" json:"request_count"`            // 总请求次数
	TokensUsed     int64      `gorm:"default:0" json:"tokens_used"`              // 已使用 tokens
	CostUsed       float64    `gorm:"type:decimal(10,4);default:0" json:"cost_used"` // 已使用金额
	LastUsedAt     *time.Time `json:"last_used_at,omitempty"`                    // 最后使用时间

	// 时间戳
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (k *APIKey) TableName() string {
	return "api_keys"
}

// GenerateKey 生成新的 API Key
func GenerateAPIKey() (key string, hash string, prefix string, err error) {
	// 生成 32 字节的随机数据
	bytes := make([]byte, 32)
	if _, err = rand.Read(bytes); err != nil {
		return "", "", "", err
	}

	// 转换为十六进制字符串
	key = APIKeyPrefix + hex.EncodeToString(bytes)

	// 计算 SHA256 哈希
	hashBytes := sha256.Sum256([]byte(key))
	hash = hex.EncodeToString(hashBytes[:])

	// 生成前缀用于显示 (如 sk-a1b2c3...)
	prefix = key[:10] + "..."

	return key, hash, prefix, nil
}

// HashAPIKey 计算 API Key 的哈希
func HashAPIKey(key string) string {
	hashBytes := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hashBytes[:])
}

// IsExpired 检查 API Key 是否过期
func (k *APIKey) IsExpired() bool {
	if k.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*k.ExpiresAt)
}

// IsActive 检查 API Key 是否可用
func (k *APIKey) IsActive() bool {
	return k.Status == "active" && !k.IsExpired()
}
