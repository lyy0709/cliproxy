/*
 * 文件作用：AI模型数据模型，定义模型配置和定价信息
 * 负责功能：
 *   - 模型基础信息（名称、平台、提供商）
 *   - 定价配置（输入/输出/缓存价格）
 *   - 模型能力和限制
 *   - 别名和分类
 * 重要程度：⭐⭐⭐ 一般（模型数据结构）
 * 依赖模块：gorm
 */
package model

import (
	"time"

	"gorm.io/gorm"
)

// AIModel AI 模型定义
type AIModel struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	Name             string         `gorm:"size:100;not null;uniqueIndex" json:"name"`              // 模型名称，如 claude-3-5-sonnet
	DisplayName      string         `gorm:"size:100" json:"display_name"`                           // 显示名称
	Platform         string         `gorm:"size:20;not null;index" json:"platform"`                 // 平台: claude/openai/gemini
	Provider         string         `gorm:"size:50" json:"provider"`                                // 提供商: anthropic/openai/google
	Description      string         `gorm:"size:500" json:"description"`                            // 描述
	Category         string         `gorm:"size:30" json:"category"`                                // 分类: chat/completion/embedding/image
	ContextSize      int            `gorm:"default:0" json:"context_size"`                          // 上下文长度
	MaxOutput        int            `gorm:"default:0" json:"max_output"`                            // 最大输出长度
	InputPrice       float64        `gorm:"type:decimal(10,6);default:0" json:"input_price"`        // 输入价格 ($/1M tokens)
	OutputPrice      float64        `gorm:"type:decimal(10,6);default:0" json:"output_price"`       // 输出价格 ($/1M tokens)
	CacheCreatePrice float64        `gorm:"type:decimal(10,6);default:0" json:"cache_create_price"` // 缓存创建价格 ($/1M tokens)
	CacheReadPrice   float64        `gorm:"type:decimal(10,6);default:0" json:"cache_read_price"`   // 缓存读取价格 ($/1M tokens)
	Enabled          bool           `gorm:"default:true" json:"enabled"`                            // 是否启用
	IsDefault        bool           `gorm:"default:false" json:"is_default"`                        // 是否默认模型
	SortOrder        int            `gorm:"default:0" json:"sort_order"`                            // 排序
	Aliases          string         `gorm:"type:text" json:"aliases"`                               // 别名列表，逗号分隔
	Capabilities     string         `gorm:"type:text" json:"capabilities"`                          // 能力列表 JSON
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *AIModel) TableName() string {
	return "ai_models"
}

// 预定义模型数据 (2025-2026年最新模型，价格来源于官方文档)
// Claude Code: https://claude.com/pricing
// OpenAI Codex: https://developers.openai.com/codex/pricing/
// Gemini CLI: https://ai.google.dev/gemini-api/docs/pricing
var DefaultModels = []AIModel{
	// ==================== Claude Code 系列 (Anthropic) ====================
	// Claude 4.5 系列 (2025 最新) - 官方价格
	{Name: "claude-opus-4-5-20251101", DisplayName: "Claude Opus 4.5", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 32000, InputPrice: 5.0, OutputPrice: 25.0, Enabled: true, IsDefault: true, SortOrder: 1, Aliases: "claude-opus-4-5,claude-opus-4.5"},
	{Name: "claude-sonnet-4-5-20250929", DisplayName: "Claude Sonnet 4.5", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 64000, InputPrice: 3.0, OutputPrice: 15.0, Enabled: true, SortOrder: 2, Aliases: "claude-sonnet-4-5,claude-sonnet-4.5"},
	{Name: "claude-haiku-4-5-20251001", DisplayName: "Claude Haiku 4.5", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 8192, InputPrice: 1.0, OutputPrice: 5.0, Enabled: true, SortOrder: 3, Aliases: "claude-haiku-4-5,claude-haiku-4.5"},

	// Claude 4.1 系列 (2025)
	{Name: "claude-opus-4-1-20250805", DisplayName: "Claude Opus 4.1", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 32000, InputPrice: 5.0, OutputPrice: 25.0, Enabled: true, SortOrder: 4, Aliases: "claude-opus-4-1,claude-opus-4.1"},

	// Claude 4 系列 (2025)
	{Name: "claude-sonnet-4-20250514", DisplayName: "Claude Sonnet 4", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 1000000, MaxOutput: 64000, InputPrice: 3.0, OutputPrice: 15.0, Enabled: true, SortOrder: 5, Aliases: "claude-sonnet-4"},
	{Name: "claude-opus-4-20250514", DisplayName: "Claude Opus 4", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 32000, InputPrice: 5.0, OutputPrice: 25.0, Enabled: true, SortOrder: 6, Aliases: "claude-opus-4"},

	// Claude 3.7 系列 (2025)
	{Name: "claude-3-7-sonnet-20250219", DisplayName: "Claude 3.7 Sonnet", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 128000, InputPrice: 3.0, OutputPrice: 15.0, Enabled: true, SortOrder: 7, Aliases: "claude-3.7-sonnet,claude-3-7-sonnet"},

	// Claude 3.5 系列 (保留)
	{Name: "claude-3-5-sonnet-20241022", DisplayName: "Claude 3.5 Sonnet", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 8192, InputPrice: 3.0, OutputPrice: 15.0, Enabled: true, SortOrder: 8, Aliases: "claude-3-5-sonnet-latest,claude-3-5-sonnet,claude-3.5-sonnet"},
	{Name: "claude-3-5-haiku-20241022", DisplayName: "Claude 3.5 Haiku", Platform: "claude", Provider: "anthropic", Category: "chat", ContextSize: 200000, MaxOutput: 8192, InputPrice: 0.8, OutputPrice: 4.0, Enabled: true, SortOrder: 9, Aliases: "claude-3-5-haiku-latest,claude-3-5-haiku,claude-3.5-haiku"},

	// ==================== OpenAI Codex 系列 ====================
	// GPT-5 Codex 系列 (2025-2026 最新)
	{Name: "gpt-5-codex", DisplayName: "GPT-5 Codex", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 200000, MaxOutput: 100000, InputPrice: 1.25, OutputPrice: 10.0, Enabled: true, SortOrder: 10, Aliases: "gpt-5,gpt5-codex"},
	{Name: "gpt-5.2-codex", DisplayName: "GPT-5.2 Codex", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 200000, MaxOutput: 100000, InputPrice: 1.25, OutputPrice: 10.0, Enabled: true, SortOrder: 11, Aliases: "gpt-5.2,gpt52-codex"},
	{Name: "gpt-5.1-codex-mini", DisplayName: "GPT-5.1 Codex Mini", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 200000, MaxOutput: 100000, InputPrice: 0.5, OutputPrice: 2.0, Enabled: true, SortOrder: 12, Aliases: "gpt-5.1-mini,gpt51-codex-mini"},
	{Name: "codex-mini-latest", DisplayName: "Codex Mini Latest", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 128000, MaxOutput: 65536, InputPrice: 1.5, OutputPrice: 6.0, Enabled: true, SortOrder: 13, Aliases: "codex-mini"},

	// GPT-4.1 系列 (2025)
	{Name: "gpt-4.1", DisplayName: "GPT-4.1", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 1047576, MaxOutput: 32768, InputPrice: 2.0, OutputPrice: 8.0, Enabled: true, SortOrder: 14, Aliases: "gpt-4.1-2025-04-14"},
	{Name: "gpt-4.1-mini", DisplayName: "GPT-4.1 Mini", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 1047576, MaxOutput: 32768, InputPrice: 0.4, OutputPrice: 1.6, Enabled: true, SortOrder: 15, Aliases: "gpt-4.1-mini-2025-04-14"},
	{Name: "gpt-4.1-nano", DisplayName: "GPT-4.1 Nano", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 1047576, MaxOutput: 32768, InputPrice: 0.1, OutputPrice: 0.4, Enabled: true, SortOrder: 16, Aliases: "gpt-4.1-nano-2025-04-14"},

	// OpenAI o1/o3 系列 (2024-2025)
	{Name: "o1", DisplayName: "o1", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 200000, MaxOutput: 100000, InputPrice: 15.0, OutputPrice: 60.0, Enabled: true, SortOrder: 17, Aliases: "o1-2024-12-17"},
	{Name: "o1-mini", DisplayName: "o1 Mini", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 128000, MaxOutput: 65536, InputPrice: 1.1, OutputPrice: 4.4, Enabled: true, SortOrder: 18},
	{Name: "o3-mini", DisplayName: "o3 Mini", Platform: "openai", Provider: "openai", Category: "chat", ContextSize: 200000, MaxOutput: 100000, InputPrice: 1.1, OutputPrice: 4.4, Enabled: true, SortOrder: 19, Aliases: "o3-mini-2025-01-31"},

	// ==================== Gemini CLI 系列 (Google) ====================
	// Gemini 3 系列 (2025-2026 最新)
	{Name: "gemini-3-flash", DisplayName: "Gemini 3 Flash", Platform: "gemini", Provider: "google", Category: "chat", ContextSize: 1048576, MaxOutput: 65535, InputPrice: 0.5, OutputPrice: 3.0, Enabled: true, SortOrder: 20, Aliases: "gemini-3.0-flash"},
	{Name: "gemini-3-pro", DisplayName: "Gemini 3 Pro", Platform: "gemini", Provider: "google", Category: "chat", ContextSize: 1048576, MaxOutput: 65535, InputPrice: 2.0, OutputPrice: 12.0, Enabled: true, SortOrder: 21, Aliases: "gemini-3.0-pro"},

	// Gemini 2.5 系列 (2025)
	{Name: "gemini-2.5-pro", DisplayName: "Gemini 2.5 Pro", Platform: "gemini", Provider: "google", Category: "chat", ContextSize: 1048576, MaxOutput: 65535, InputPrice: 1.25, OutputPrice: 10.0, Enabled: true, SortOrder: 22, Aliases: "gemini-2.5-pro-exp-03-25"},
	{Name: "gemini-2.5-flash", DisplayName: "Gemini 2.5 Flash", Platform: "gemini", Provider: "google", Category: "chat", ContextSize: 1048576, MaxOutput: 65535, InputPrice: 0.3, OutputPrice: 2.5, Enabled: true, SortOrder: 23},
	{Name: "gemini-2.5-flash-lite", DisplayName: "Gemini 2.5 Flash Lite", Platform: "gemini", Provider: "google", Category: "chat", ContextSize: 1048576, MaxOutput: 65535, InputPrice: 0.1, OutputPrice: 0.4, Enabled: true, SortOrder: 24},

	// Gemini 2.0 系列 (2024-2025)
	{Name: "gemini-2.0-flash", DisplayName: "Gemini 2.0 Flash", Platform: "gemini", Provider: "google", Category: "chat", ContextSize: 1048576, MaxOutput: 8192, InputPrice: 0.1, OutputPrice: 0.4, Enabled: true, SortOrder: 25, Aliases: "gemini-2.0-flash-exp"},
}
