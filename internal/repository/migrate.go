/*
 * 文件作用：数据库迁移和初始化，管理数据库表结构和默认数据
 * 负责功能：
 *   - 自动迁移所有数据表
 *   - 初始化默认配置项
 *   - 初始化默认客户端过滤规则
 *   - 初始化默认错误消息配置
 * 重要程度：⭐⭐⭐⭐ 重要（数据库初始化核心）
 * 依赖模块：model, logger, gorm
 */
package repository

import (
	"cli-proxy/internal/model"
	"cli-proxy/pkg/logger"
)

func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.Proxy{},
		&model.Account{},
		&model.AccountGroup{},
		&model.RequestLog{},
		&model.AIModel{},
		&model.APIKey{},
		&model.DailyUsage{},
		&model.SystemConfig{},
		&model.UsageRecord{},
		&model.OperationLog{},
		// 客户端过滤相关
		&model.ClientType{},
		&model.ClientFilterRule{},
		&model.ClientFilterConfig{},
		// 错误消息配置
		&model.ErrorMessage{},
		// 错误规则配置
		&model.ErrorRule{},
		// 模型映射
		&model.ModelMapping{},
	)
}

// InitDefaultConfigs 初始化默认系统配置
func InitDefaultConfigs() error {
	log := logger.GetLogger("main")

	for _, cfg := range model.DefaultConfigs {
		var existing model.SystemConfig
		result := DB.Where("config_key = ?", cfg.Key).First(&existing)
		if result.Error == nil {
			// 已存在，跳过
			continue
		}

		if err := DB.Create(&cfg).Error; err != nil {
			log.Error("创建默认配置失败: %s, %v", cfg.Key, err)
			continue
		}
		log.Info("默认配置已创建: %s = %s", cfg.Key, cfg.Value)
	}

	return nil
}

// InitDefaultClientFilters 初始化默认客户端过滤配置
func InitDefaultClientFilters() error {
	repo := NewClientFilterRepository()
	return repo.InitDefaultData()
}

// InitDefaultErrorMessages 初始化默认错误消息配置
func InitDefaultErrorMessages() error {
	repo := NewErrorMessageRepository()
	return repo.InitDefaultData()
}

// InitDefaultErrorRules 初始化默认错误规则
func InitDefaultErrorRules() error {
	repo := NewErrorRuleRepository()
	return repo.InitDefaultRules()
}
