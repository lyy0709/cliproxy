/*
 * 文件作用：网关数据模型，定义 xyrt 网关的数据结构
 * 负责功能：
 *   - 网关基础信息（名称、URL、状态）
 *   - 测试状态追踪（延迟、错误信息）
 *   - 默认网关设置
 * 重要程度：⭐⭐⭐ 重要（xyrt 授权核心）
 * 依赖模块：gorm
 */
package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Gateway xyrt 网关配置
type Gateway struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`        // 网关名称
	URL         string         `gorm:"size:500;not null" json:"url"`         // 网关地址
	Enabled     bool           `gorm:"default:true" json:"enabled"`          // 是否启用
	IsDefault   bool           `gorm:"default:false" json:"is_default"`      // 是否为默认网关
	TestStatus  string         `gorm:"size:20" json:"test_status"`           // 测试状态: success, failed, pending
	TestLatency int            `gorm:"default:0" json:"test_latency"`        // 测试延迟(ms)
	TestError   string         `gorm:"size:500" json:"test_error,omitempty"` // 测试错误信息
	LastTestAt  *time.Time     `json:"last_test_at,omitempty"`               // 最后测试时间
	Remark      string         `gorm:"size:500" json:"remark,omitempty"`     // 备注
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (g *Gateway) TableName() string {
	return "gateways"
}

// GetBaseURL 获取规范化的基础 URL（去除末尾斜杠）
func (g *Gateway) GetBaseURL() string {
	return strings.TrimSuffix(g.URL, "/")
}

// GetCodexURL 获取 Codex API 完整 URL
func (g *Gateway) GetCodexURL() string {
	return g.GetBaseURL() + "/backend-api/codex"
}

// GetRefreshURL 获取 Token 刷新 URL
func (g *Gateway) GetRefreshURL() string {
	return g.GetBaseURL() + "/auth/refresh"
}
