/*
 * 文件作用：路由注册中心，定义所有HTTP接口路由
 * 负责功能：
 *   - 管理员登录接口
 *   - 管理后台路由（/api/admin/*）
 *   - 代理转发路由（/claude/*, /openai/*, /responses）
 *   - 中间件配置（JWT、API Key、操作日志）
 *   - 静态文件服务
 * 重要程度：⭐⭐⭐⭐⭐ 核心（所有请求的入口）
 * 依赖模块：middleware, handler
 */
package handler

import (
	"cli-proxy/internal/middleware"
	"cli-proxy/internal/repository"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 全局操作日志中间件（放在认证之后，记录所有写操作）
	r.Use(middleware.OperationLogger())

	// 初始化处理器
	adminHandler := NewAdminHandler()
	captchaHandler := NewCaptchaHandler()
	accountHandler := NewAccountHandler()
	proxyHandler := NewProxyHandler()
	requestLogHandler := NewRequestLogHandler()
	oauthHandler := NewOAuthHandler()
	usageHandler := NewUsageHandler()
	cacheHandler := NewCacheHandler()
	operationLogHandler := NewOperationLogHandler()

	// 模型管理
	modelRepo := repository.NewAIModelRepository(repository.GetDB())
	modelHandler := NewAIModelHandler(modelRepo)
	// 初始化默认模型
	modelRepo.InitDefaultModels()

	// 模型映射
	modelMappingHandler := NewModelMappingHandler()

	// OpenAI Responses API Handler
	openaiResponsesHandler := NewOpenAIResponsesHandler()

	// API Key Handler
	apiKeyHandler := NewAPIKeyHandler()

	// 管理员认证接口（公开）
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", adminHandler.Login)
		auth.GET("/captcha", captchaHandler.Generate) // 生成验证码
	}

	// 管理员状态接口（需要 JWT 认证）
	adminStatus := r.Group("/api/admin")
	adminStatus.Use(middleware.JWTAuth())
	adminStatus.Use(middleware.AdminRequired())
	adminStatus.Use(middleware.MustChangePasswordGuard())
	{
		adminStatus.GET("/status", adminHandler.GetStatus)                         // 获取管理员状态
		adminStatus.POST("/force-change-password", adminHandler.ForceChangePassword) // 首次强制修改密码
		adminStatus.PUT("/password", adminHandler.ChangePassword)                   // 修改密码
	}

	// ========== 代理转发接口 (需要 API Key 认证) ==========
	proxyGroup := r.Group("")
	proxyGroup.Use(middleware.APIKeyAuth())
	proxyGroup.Use(middleware.ClientFilter())        // 客户端过滤
	proxyGroup.Use(middleware.CheckAllowedClients()) // API Key 客户端限制检查
	{
		// ========== 按平台区分的路由 ==========
		// Claude 平台 - 使用 Claude 原生格式
		proxyGroup.POST("/claude/v1/messages", proxyHandler.ClaudeMessages)

		// OpenAI 平台 - 使用 OpenAI 原生格式
		proxyGroup.POST("/openai/v1/chat/completions", proxyHandler.OpenAIChatCompletions)

		// OpenAI Responses API (Codex CLI) - 平台路由版本
		proxyGroup.POST("/openai/responses", openaiResponsesHandler.HandleResponses)
		proxyGroup.POST("/openai/v1/responses", openaiResponsesHandler.HandleResponses)
		proxyGroup.POST("/openai/responses/compact", openaiResponsesHandler.HandleResponses)
		proxyGroup.POST("/openai/v1/responses/compact", openaiResponsesHandler.HandleResponses)

		// OpenAI Responses API (Codex CLI) - 兼容根路径（部分客户端固定使用 /responses 或 /v1/responses）
		proxyGroup.POST("/responses", openaiResponsesHandler.HandleResponses)
		proxyGroup.POST("/v1/responses", openaiResponsesHandler.HandleResponses)
		proxyGroup.POST("/responses/compact", openaiResponsesHandler.HandleResponses)
		proxyGroup.POST("/v1/responses/compact", openaiResponsesHandler.HandleResponses)

		// Gemini 平台 - 使用 Gemini 原生格式
		proxyGroup.POST("/gemini/v1/chat", proxyHandler.GeminiChat)

		// ========== API Key 持有者自查接口 ==========
		// 使用 x-api-key 认证，允许 API Key 持有者查询自己的使用情况（兼容旧路径）
		proxyGroup.GET("/usage/me", usageHandler.GetMyUsage)              // 获取使用量汇总
		proxyGroup.GET("/usage/me/daily", usageHandler.GetMyDailyStats)   // 获取每日统计
		proxyGroup.GET("/usage/me/records", usageHandler.GetMyRecords)    // 获取使用记录
	}

	// ========== API Key 自助查询接口 (需要 API Key 认证) ==========
	keyGroup := r.Group("/api/key")
	keyGroup.Use(middleware.APIKeyAuth())
	{
		keyGroup.GET("/info", apiKeyHandler.GetMyInfo)           // 查询当前 Key 信息
		keyGroup.GET("/usage", usageHandler.GetMyUsage)          // 查询当前 Key 使用量
		keyGroup.GET("/usage/daily", usageHandler.GetMyDailyStats)   // 查询每日使用量
		keyGroup.GET("/usage/records", usageHandler.GetMyRecords)    // 查询使用记录
	}

	// 需要认证的接口（管理员）
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.MustChangePasswordGuard())
	admin.Use(middleware.AdminRequired())
	{
		// API Key 管理（管理员直接管理）
		apiKeys := admin.Group("/api-keys")
		{
			apiKeys.GET("", apiKeyHandler.AdminListAll)
			apiKeys.POST("", apiKeyHandler.AdminCreate)
			apiKeys.GET("/lookup", apiKeyHandler.AdminLookup)          // 按ID批量查询
			apiKeys.GET("/:id", apiKeyHandler.AdminGet)
			apiKeys.PUT("/:id", apiKeyHandler.AdminUpdate)
			apiKeys.DELETE("/:id", apiKeyHandler.AdminDelete)
			apiKeys.PUT("/:id/toggle", apiKeyHandler.AdminToggleStatus)
			apiKeys.GET("/:id/logs", apiKeyHandler.AdminGetAPIKeyLogs)
			apiKeys.GET("/:id/usage", usageHandler.GetAPIKeyUsage)
		}

		// 使用统计
		usage := admin.Group("/usage")
		{
			usage.GET("/summary", usageHandler.GetTotalUsageSummary)    // 总使用量汇总
			usage.GET("/daily", usageHandler.GetDailySummary)           // 每日汇总
			usage.GET("/models", usageHandler.GetModelSummary)          // 模型汇总
			usage.GET("/api-keys", usageHandler.GetAPIKeyUsageSummary)  // 各 API Key 使用汇总
			usage.GET("/records", usageHandler.GetAllUsageRecords)      // 所有使用记录
		}

		// 模型价格查询
		admin.GET("/models", usageHandler.GetModels)

		// 账户管理
		accounts := admin.Group("/accounts")
		{
			accounts.GET("/types", accountHandler.GetTypes)
			accounts.GET("", accountHandler.List)
			accounts.POST("", accountHandler.Create)
			accounts.GET("/:id", accountHandler.Get)
			accounts.PUT("/:id", accountHandler.Update)
			accounts.DELETE("/:id", accountHandler.Delete)
			accounts.PUT("/:id/status", accountHandler.UpdateStatus)
			// 健康检测相关操作
			accounts.POST("/:id/health-check", accountHandler.HealthCheck)   // 手动触发单个账号健康检测
			accounts.POST("/:id/recover", accountHandler.ForceRecover)       // 强制恢复账号
			accounts.POST("/:id/refresh-token", accountHandler.RefreshToken) // 刷新 Token

			// 用量查询相关操作
			accounts.GET("/:id/usage", accountHandler.FetchUsage)         // 查询单个账户用量
			accounts.POST("/batch-usage", accountHandler.BatchFetchUsage) // 批量同步所有账户用量
		}

		// 健康检测服务管理
		healthCheck := admin.Group("/health-check")
		{
			healthCheck.GET("/status", accountHandler.GetHealthCheckStatus) // 获取健康检测服务状态
			healthCheck.POST("/trigger", accountHandler.TriggerHealthCheck) // 手动触发全局健康检测
		}

		// 账户分组管理
		groups := admin.Group("/account-groups")
		{
			groups.GET("", accountHandler.ListGroups)
			groups.GET("/all", accountHandler.GetAllGroups)
			groups.POST("", accountHandler.CreateGroup)
			groups.GET("/:id", accountHandler.GetGroup)
			groups.PUT("/:id", accountHandler.UpdateGroup)
			groups.DELETE("/:id", accountHandler.DeleteGroup)
			groups.POST("/:id/accounts", accountHandler.AddAccountToGroup)
			groups.DELETE("/:id/accounts/:accountId", accountHandler.RemoveAccountFromGroup)
		}

		// OAuth 授权
		oauth := admin.Group("/oauth")
		{
			oauth.POST("/generate-url", oauthHandler.GenerateURL)
			oauth.POST("/exchange", oauthHandler.Exchange)
			oauth.POST("/cookie-auth", oauthHandler.CookieAuth)
		}

		// 请求日志
		logs := admin.Group("/logs")
		{
			logs.GET("", requestLogHandler.List)
			logs.GET("/summary", requestLogHandler.GetSummary)
			logs.GET("/account-load", requestLogHandler.GetAccountLoadStats)
		}

		// 操作日志
		opLogs := admin.Group("/operation-logs")
		{
			opLogs.GET("", operationLogHandler.List)
			opLogs.GET("/stats", operationLogHandler.GetStats)
			opLogs.GET("/:id", operationLogHandler.Get)
			opLogs.DELETE("/cleanup", operationLogHandler.Cleanup)
		}

		// 模型管理
		models := admin.Group("/models-config")
		{
			models.GET("", modelHandler.List)
			models.GET("/platforms", modelHandler.GetPlatforms)
			models.POST("", modelHandler.Create)
			models.GET("/:id", modelHandler.Get)
			models.PUT("/:id", modelHandler.Update)
			models.DELETE("/:id", modelHandler.Delete)
			models.PUT("/:id/toggle", modelHandler.ToggleEnabled)
			models.POST("/init-defaults", modelHandler.InitDefaults)
			models.POST("/reset-defaults", modelHandler.ResetDefaults)
		}

		// 模型映射管理
		modelMappings := admin.Group("/model-mappings")
		{
			modelMappings.GET("", modelMappingHandler.List)
			modelMappings.POST("", modelMappingHandler.Create)
			modelMappings.GET("/:id", modelMappingHandler.Get)
			modelMappings.PUT("/:id", modelMappingHandler.Update)
			modelMappings.DELETE("/:id", modelMappingHandler.Delete)
			modelMappings.POST("/:id/toggle", modelMappingHandler.Toggle)
			modelMappings.POST("/refresh", modelMappingHandler.RefreshCache)
			modelMappings.GET("/cache", modelMappingHandler.GetCacheStats)
		}

		// 缓存管理
		cache := admin.Group("/cache")
		{
			cache.GET("/stats", cacheHandler.GetStats)                       // 获取缓存统计
			cache.GET("/sessions", cacheHandler.ListSessions)                // 列出所有会话
			cache.DELETE("/sessions/:sessionId", cacheHandler.RemoveSession) // 移除会话
			cache.GET("/accounts", cacheHandler.ListAccountsCache)           // 列出有缓存的账号（聚合）
			cache.GET("/unavailable", cacheHandler.ListUnavailableAccounts)  // 列出不可用账户
			cache.POST("/clear", cacheHandler.ClearCache)                    // 按类型清理缓存
			cache.DELETE("/api-keys/:id", cacheHandler.ClearAPIKeyCache)     // 清理 API Key 缓存
			cache.GET("/config", cacheHandler.GetCacheConfig)                // 获取缓存配置
			cache.PUT("/config", cacheHandler.UpdateCacheConfig)             // 更新缓存配置
		}

		// 账户缓存管理（并发控制和不可用标记）
		accountCache := admin.Group("/accounts/:id/cache")
		{
			accountCache.DELETE("/sessions", cacheHandler.ClearAccountSessions)       // 清除账户会话
			accountCache.POST("/unavailable", cacheHandler.MarkAccountUnavailable)    // 标记账户不可用
			accountCache.DELETE("/unavailable", cacheHandler.ClearAccountUnavailable) // 清除不可用标记
			accountCache.GET("/concurrency", cacheHandler.GetAccountConcurrency)      // 获取并发信息
			accountCache.PUT("/concurrency", cacheHandler.SetAccountConcurrencyLimit) // 设置并发限制
			accountCache.DELETE("/concurrency", cacheHandler.ResetAccountConcurrency) // 重置并发计数
		}

		// 系统配置管理
		configHandler := NewConfigHandler()
		configs := admin.Group("/configs")
		{
			configs.GET("", configHandler.GetAll)                           // 获取所有配置
			configs.GET("/category/:category", configHandler.GetByCategory) // 获取分类配置
			configs.PUT("", configHandler.Update)                           // 更新配置
		}

		// 代理配置管理
		proxyConfigs := admin.Group("/proxy-configs")
		{
			proxyConfigs.GET("", ListProxyConfigs)                    // 获取代理列表
			proxyConfigs.GET("/enabled", GetEnabledProxyConfigs)      // 获取启用的代理（用于下拉选择）
			proxyConfigs.GET("/default", GetDefaultProxyConfig)       // 获取默认代理
			proxyConfigs.DELETE("/default", ClearDefaultProxyConfig)  // 清除默认代理
			proxyConfigs.POST("", CreateProxyConfig)                  // 创建代理
			proxyConfigs.POST("/test", TestProxyConnectivity)         // 测试代理连通性
			proxyConfigs.GET("/:id", GetProxyConfig)                  // 获取单个代理
			proxyConfigs.PUT("/:id", UpdateProxyConfig)               // 更新代理
			proxyConfigs.DELETE("/:id", DeleteProxyConfig)            // 删除代理
			proxyConfigs.PUT("/:id/toggle", ToggleProxyConfigEnabled) // 切换启用状态
			proxyConfigs.PUT("/:id/default", SetDefaultProxyConfig)   // 设置为默认代理
		}

		// 网关配置管理（xyrt 专用）
		gatewayHandler := NewGatewayHandler()
		gateways := admin.Group("/gateways")
		{
			gateways.GET("", gatewayHandler.List)                   // 获取网关列表
			gateways.GET("/enabled", gatewayHandler.GetEnabled)     // 获取启用的网关（用于下拉选择）
			gateways.GET("/default", gatewayHandler.GetDefault)     // 获取默认网关
			gateways.DELETE("/default", gatewayHandler.ClearDefault) // 清除默认网关
			gateways.POST("", gatewayHandler.Create)                // 创建网关
			gateways.POST("/test", gatewayHandler.Test)             // 测试网关连通性
			gateways.GET("/:id", gatewayHandler.Get)                // 获取单个网关
			gateways.PUT("/:id", gatewayHandler.Update)             // 更新网关
			gateways.DELETE("/:id", gatewayHandler.Delete)          // 删除网关
			gateways.PUT("/:id/toggle", gatewayHandler.ToggleEnabled) // 切换启用状态
			gateways.PUT("/:id/default", gatewayHandler.SetDefault)   // 设置为默认网关
			gateways.POST("/:id/test", gatewayHandler.TestByID)       // 测试指定网关
		}

		// 系统监控
		monitorHandler := NewSystemMonitorHandler()
		monitor := admin.Group("/monitor")
		{
			monitor.GET("", monitorHandler.GetMonitorData)           // 获取完整监控数据
			monitor.GET("/system", monitorHandler.GetSystemStats)    // 系统资源
			monitor.GET("/cache", monitorHandler.GetCacheStats)      // 缓存统计
			monitor.GET("/mysql", monitorHandler.GetMySQLStats)      // MySQL 统计
			monitor.GET("/accounts", monitorHandler.GetAccountStats) // 账号统计
			monitor.GET("/api-keys", monitorHandler.GetAPIKeyStats)  // API Key 统计
			monitor.GET("/today", monitorHandler.GetTodayUsageStats) // 今日使用统计
		}

		// 错误消息管理
		errorMsgHandler := NewErrorMessageHandler()
		errorMessages := admin.Group("/error-messages")
		{
			errorMessages.GET("", errorMsgHandler.List)
			errorMessages.GET("/code/:code", errorMsgHandler.GetByCode)
			errorMessages.GET("/:id", errorMsgHandler.Get)
			errorMessages.POST("", errorMsgHandler.Create)
			errorMessages.PUT("/:id", errorMsgHandler.Update)
			errorMessages.DELETE("/:id", errorMsgHandler.Delete)
			errorMessages.PUT("/:id/toggle", errorMsgHandler.ToggleEnabled)
			errorMessages.POST("/init", errorMsgHandler.InitDefault)
			errorMessages.POST("/refresh", errorMsgHandler.RefreshCache)
			errorMessages.PUT("/enable-all", errorMsgHandler.EnableAll)
			errorMessages.PUT("/disable-all", errorMsgHandler.DisableAll)
		}

		// 系统日志查看
		systemLogHandler := NewSystemLogHandler()
		sysLogs := admin.Group("/system-logs")
		{
			sysLogs.GET("/files", systemLogHandler.ListFiles)       // 获取日志文件列表
			sysLogs.GET("/read", systemLogHandler.ReadFile)         // 读取日志内容
			sysLogs.GET("/tail", systemLogHandler.TailFile)         // 查看日志末尾
			sysLogs.GET("/download", systemLogHandler.DownloadFile) // 下载日志文件
			sysLogs.DELETE("/file", systemLogHandler.DeleteFile)    // 删除日志文件
		}

		// 客户端过滤管理
		clientFilterHandler := NewClientFilterHandler()
		clientFilter := admin.Group("/client-filter")
		{
			// 全局配置
			clientFilter.GET("/config", clientFilterHandler.GetConfig)
			clientFilter.PUT("/config", clientFilterHandler.UpdateConfig)
			clientFilter.POST("/reload", clientFilterHandler.ReloadCache)
			clientFilter.POST("/test", clientFilterHandler.TestValidation)

			// 客户端类型管理
			clientTypes := clientFilter.Group("/client-types")
			{
				clientTypes.GET("", clientFilterHandler.ListClientTypes)
				clientTypes.POST("", clientFilterHandler.CreateClientType)
				clientTypes.GET("/:id", clientFilterHandler.GetClientType)
				clientTypes.PUT("/:id", clientFilterHandler.UpdateClientType)
				clientTypes.DELETE("/:id", clientFilterHandler.DeleteClientType)
				clientTypes.PUT("/:id/toggle", clientFilterHandler.ToggleClientType)
			}

			// 过滤规则管理
			rules := clientFilter.Group("/rules")
			{
				rules.GET("", clientFilterHandler.ListRules)
				rules.POST("", clientFilterHandler.CreateRule)
				rules.GET("/:id", clientFilterHandler.GetRule)
				rules.PUT("/:id", clientFilterHandler.UpdateRule)
				rules.DELETE("/:id", clientFilterHandler.DeleteRule)
				rules.PUT("/:id/toggle", clientFilterHandler.ToggleRule)
			}
		}

		// 错误规则管理
		errorRuleHandler := NewErrorRuleHandler()
		errorRules := admin.Group("/error-rules")
		{
			errorRules.GET("", errorRuleHandler.List)
			errorRules.POST("", errorRuleHandler.Create)
			errorRules.GET("/:id", errorRuleHandler.Get)
			errorRules.PUT("/:id", errorRuleHandler.Update)
			errorRules.DELETE("/:id", errorRuleHandler.Delete)
			errorRules.POST("/reset", errorRuleHandler.ResetToDefault)
			errorRules.POST("/refresh", errorRuleHandler.RefreshCache)
			errorRules.PUT("/enable-all", errorRuleHandler.EnableAll)
			errorRules.PUT("/disable-all", errorRuleHandler.DisableAll)
		}
	}

	// 静态文件（前端）- 放在最后作为兜底
	RegisterStatic(r)
}
