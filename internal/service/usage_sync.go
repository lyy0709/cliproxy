/*
 * 文件作用：用量同步服务，查询上游平台的账户使用情况
 * 负责功能：
 *   - Claude OAuth Usage API 查询
 *   - OpenAI Codex 响应头用量提取
 *   - Gemini Code Assist API 查询
 *   - 批量并发查询
 *   - 定时自动同步
 *   - 缓存控制（避免频繁查询）
 * 重要程度：⭐⭐⭐⭐ 重要（用量监控核心）
 * 依赖模块：repository, model, adapter, health_check
 */
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/proxy/adapter"
	"cli-proxy/internal/proxy/scheduler"
	"cli-proxy/internal/repository"
	"cli-proxy/pkg/logger"
)

// UsageSyncService 用量同步服务（单例）
type UsageSyncService struct {
	accountRepo *repository.AccountRepository
	healthCheck *AccountHealthCheckService
	log         *logger.Logger

	// 缓存控制：避免 5 分钟内重复查询同一账户
	cache    sync.Map      // key: accountID, value: time.Time
	cacheTTL time.Duration // 默认 5 分钟

	// 定时器
	ticker   *time.Ticker
	stopChan chan struct{}
	running  bool
	mu       sync.Mutex
}

var usageSyncService *UsageSyncService
var usageSyncOnce sync.Once

// GetUsageSyncService 获取用量同步服务单例
func GetUsageSyncService() *UsageSyncService {
	usageSyncOnce.Do(func() {
		usageSyncService = &UsageSyncService{
			accountRepo: repository.NewAccountRepository(),
			healthCheck: GetAccountHealthCheckService(),
			log:         logger.GetLogger("usage_sync"),
			cacheTTL:    5 * time.Minute,
			stopChan:    make(chan struct{}),
		}
	})
	return usageSyncService
}

// ==================== Claude OAuth Usage 查询 ====================

// ClaudeUsageResponse Claude OAuth Usage API 响应结构
type ClaudeUsageResponse struct {
	FiveHour struct {
		Utilization float64 `json:"utilization"` // 0-1 的小数
		ResetsAt    string  `json:"resets_at"`   // RFC3339 格式
	} `json:"five_hour"`
	SevenDay struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"seven_day"`
	SevenDaySonnet struct {
		Utilization float64 `json:"utilization"`
		ResetsAt    string  `json:"resets_at"`
	} `json:"seven_day_sonnet"`
}

// FetchClaudeOAuthUsage 查询单个 Claude OAuth 账户的用量
func (s *UsageSyncService) FetchClaudeOAuthUsage(accountID uint) (*ClaudeUsageResponse, error) {
	// 检查缓存
	if lastTime, ok := s.cache.Load(fmt.Sprintf("claude_%d", accountID)); ok {
		if time.Since(lastTime.(time.Time)) < s.cacheTTL {
			s.log.Debug("Claude 账户 %d 用量在缓存有效期内，跳过查询", accountID)
			return nil, nil
		}
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("获取账户失败: %w", err)
	}

	// 检查是否是 OAuth 账户
	if account.AccessToken == "" {
		s.log.Debug("账户 %d (%s) 非 OAuth 账户，跳过", accountID, account.Name)
		return nil, nil
	}

	// Token 已过期则先刷新
	if account.TokenExpiry != nil && account.TokenExpiry.Before(time.Now()) {
		s.log.Warn("账户 %d (%s) Token 已过期，尝试刷新", accountID, account.Name)
		if err := s.refreshAccountToken(account); err != nil {
			return nil, fmt.Errorf("Token 刷新失败: %w", err)
		}
		account, _ = s.accountRepo.GetByID(accountID)
	}

	// 构建请求
	client := adapter.GetHTTPClient(account)
	req, err := http.NewRequest("GET", "https://api.anthropic.com/api/oauth/usage", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+account.AccessToken)
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")
	req.Header.Set("User-Agent", "claude-cli/2.0.53 (external, cli)")
	req.Header.Set("Accept", "application/json")

	s.log.Debug("查询 Claude 账户 %d (%s) 用量", accountID, account.Name)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 错误处理
	if resp.StatusCode == 403 {
		// Setup Token 账户，静默跳过
		s.log.Debug("账户 %d (%s) 返回 403，可能是 Setup Token 而非 OAuth", accountID, account.Name)
		return nil, nil
	}

	if resp.StatusCode == 401 {
		// Token 过期，尝试刷新
		s.log.Warn("账户 %d (%s) Token 过期，尝试刷新", accountID, account.Name)
		if err := s.refreshAccountToken(account); err != nil {
			return nil, fmt.Errorf("Token 刷新失败: %w", err)
		}

		// 重新获取账户信息并重试
		account, _ = s.accountRepo.GetByID(accountID)
		req.Header.Set("Authorization", "Bearer "+account.AccessToken)

		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("重试请求失败: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回错误: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应
	var usage ClaudeUsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 转换并更新数据库（使用 repository.ClaudeUsageData）
	usageData := &repository.ClaudeUsageData{
		FiveHour: struct {
			Utilization *float64 `json:"utilization"`
			ResetsAt    string   `json:"resets_at"`
		}{
			Utilization: floatPtr(usage.FiveHour.Utilization * 100), // 转换为百分比
			ResetsAt:    usage.FiveHour.ResetsAt,
		},
		SevenDay: struct {
			Utilization *float64 `json:"utilization"`
			ResetsAt    string   `json:"resets_at"`
		}{
			Utilization: floatPtr(usage.SevenDay.Utilization * 100),
			ResetsAt:    usage.SevenDay.ResetsAt,
		},
		SevenDaySonnet: struct {
			Utilization *float64 `json:"utilization"`
			ResetsAt    string   `json:"resets_at"`
		}{
			Utilization: floatPtr(usage.SevenDaySonnet.Utilization * 100),
			ResetsAt:    usage.SevenDaySonnet.ResetsAt,
		},
	}

	if err := s.accountRepo.UpdateClaudeUsage(accountID, usageData); err != nil {
		return nil, fmt.Errorf("更新数据库失败: %w", err)
	}

	// 更新缓存
	s.cache.Store(fmt.Sprintf("claude_%d", accountID), time.Now())

	s.log.Info("Claude 账户 %d (%s) 用量更新成功: 5H=%.1f%%, 7D=%.1f%%",
		accountID, account.Name, usage.FiveHour.Utilization*100, usage.SevenDay.Utilization*100)

	return &usage, nil
}

// BatchFetchClaudeUsage 批量并发查询 Claude 账户用量
func (s *UsageSyncService) BatchFetchClaudeUsage(accountIDs []uint, maxConcurrency int) (map[uint]*ClaudeUsageResponse, map[uint]error) {
	if maxConcurrency <= 0 {
		maxConcurrency = 5
	}

	sem := make(chan struct{}, maxConcurrency) // 信号量控制并发数
	var wg sync.WaitGroup

	results := &sync.Map{}
	errors := &sync.Map{}

	for _, id := range accountIDs {
		wg.Add(1)
		go func(accountID uint) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量
			defer func() { <-sem }() // 释放信号量

			usage, err := s.FetchClaudeOAuthUsage(accountID)
			if err != nil {
				errors.Store(accountID, err)
			} else if usage != nil {
				results.Store(accountID, usage)
			}
		}(id)
	}

	wg.Wait()

	// 转换 sync.Map 到普通 map
	resultMap := make(map[uint]*ClaudeUsageResponse)
	errorMap := make(map[uint]error)

	results.Range(func(key, value interface{}) bool {
		resultMap[key.(uint)] = value.(*ClaudeUsageResponse)
		return true
	})

	errors.Range(func(key, value interface{}) bool {
		errorMap[key.(uint)] = value.(error)
		return true
	})

	return resultMap, errorMap
}

// SyncAllClaudeAccounts 同步所有 Claude OAuth 账户用量
func (s *UsageSyncService) SyncAllClaudeAccounts() (synced int, failed int, err error) {
	// 获取所有启用的 Claude Official 账户
	accounts, err := s.accountRepo.GetEnabledByType(model.AccountTypeClaudeOfficial)
	if err != nil {
		return 0, 0, fmt.Errorf("获取账户列表失败: %w", err)
	}

	// 过滤出 OAuth 账户（有 AccessToken）
	var oauthAccountIDs []uint
	for _, acc := range accounts {
		if acc.AccessToken != "" {
			oauthAccountIDs = append(oauthAccountIDs, acc.ID)
		}
	}

	if len(oauthAccountIDs) == 0 {
		s.log.Info("没有找到 Claude OAuth 账户")
		return 0, 0, nil
	}

	s.log.Info("开始同步 %d 个 Claude OAuth 账户用量", len(oauthAccountIDs))

	// 批量并发查询
	results, errors := s.BatchFetchClaudeUsage(oauthAccountIDs, 5)

	synced = len(results)
	failed = len(errors)

	for accID, err := range errors {
		s.log.Error("账户 %d 用量查询失败: %v", accID, err)
	}

	s.log.Info("Claude 用量同步完成: 成功 %d, 失败 %d", synced, failed)

	return synced, failed, nil
}

// ==================== OpenAI/Codex Usage 查询 ====================

// FetchOpenAICodexUsage 查询 OpenAI Codex 账户用量
// 通过发送测试请求获取响应头中的用量信息
func (s *UsageSyncService) FetchOpenAICodexUsage(accountID uint) (*repository.OpenAICodexUsageResponse, error) {
	// 检查缓存
	if lastTime, ok := s.cache.Load(fmt.Sprintf("openai_%d", accountID)); ok {
		if time.Since(lastTime.(time.Time)) < s.cacheTTL {
			s.log.Debug("OpenAI 账户 %d 用量在缓存有效期内，跳过查询", accountID)
			return nil, nil
		}
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("获取账户失败: %w", err)
	}

	// 检查是否有 AccessToken
	if account.AccessToken == "" {
		s.log.Debug("账户 %d (%s) 无 AccessToken，跳过", accountID, account.Name)
		return nil, nil
	}

	// Token 已过期则先刷新
	if account.TokenExpiry != nil && account.TokenExpiry.Before(time.Now()) {
		s.log.Warn("账户 %d (%s) Token 已过期，尝试刷新", accountID, account.Name)
		if err := s.refreshAccountToken(account); err != nil {
			return nil, fmt.Errorf("Token 刷新失败: %w", err)
		}
		account, _ = s.accountRepo.GetByID(accountID)
	}

	baseURL := account.BaseURL
	if baseURL == "" {
		baseURL = "https://chatgpt.com/backend-api/codex"
	}
	baseURL = strings.TrimSuffix(baseURL, "/")
	targetURL := baseURL + "/responses"

	// 构建最小请求
	reqBody := map[string]interface{}{
		"model":             "gpt-5",
		"input":             "hi",
		"max_output_tokens": 1,
		"stream":            false,
		"store":             false,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	createRequest := func(accessToken string) (*http.Request, error) {
		req, err := http.NewRequest("POST", targetURL, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("openai-beta", "responses=experimental")
		if account.OrganizationID != "" {
			req.Header.Set("chatgpt-account-id", account.OrganizationID)
		}
		return req, nil
	}

	req, err := createRequest(account.AccessToken)
	if err != nil {
		return nil, err
	}

	client := adapter.GetSmartHTTPClient(account, baseURL)

	s.log.Debug("查询 OpenAI 账户 %d (%s) Codex 用量", accountID, account.Name)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		s.log.Warn("账户 %d (%s) Token 过期，尝试刷新", accountID, account.Name)
		if err := s.refreshAccountToken(account); err != nil {
			return nil, fmt.Errorf("Token 刷新失败: %w", err)
		}
		account, _ = s.accountRepo.GetByID(accountID)

		req, err = createRequest(account.AccessToken)
		if err != nil {
			return nil, err
		}
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("重试请求失败: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 返回错误: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// 提取响应头
	usage := s.extractCodexUsageFromHeaders(resp.Header)

	// 检查是否提取到用量信息
	if usage.PrimaryUsedPercent == nil && usage.SecondaryUsedPercent == nil {
		s.log.Debug("账户 %d (%s) 响应头中未找到 Codex 用量信息", accountID, account.Name)
		return nil, nil
	}

	// 更新数据库
	if err := s.accountRepo.UpdateCodexUsage(accountID, usage); err != nil {
		return nil, fmt.Errorf("更新数据库失败: %w", err)
	}

	// 更新缓存
	s.cache.Store(fmt.Sprintf("openai_%d", accountID), time.Now())

	s.log.Info("OpenAI 账户 %d (%s) Codex 用量更新成功", accountID, account.Name)

	return usage, nil
}

// extractCodexUsageFromHeaders 从响应头提取 Codex 用量
func (s *UsageSyncService) extractCodexUsageFromHeaders(headers http.Header) *repository.OpenAICodexUsageResponse {
	getFloat := func(key string) *float64 {
		if val := headers.Get(key); val != "" {
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				return &f
			}
		}
		return nil
	}

	getInt := func(key string) *int64 {
		if val := headers.Get(key); val != "" {
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				return &i
			}
		}
		return nil
	}

	return &repository.OpenAICodexUsageResponse{
		PrimaryUsedPercent:         getFloat("x-codex-primary-used-percent"),
		PrimaryResetAfterSeconds:   getInt("x-codex-primary-reset-after-seconds"),
		PrimaryWindowMinutes:       getInt("x-codex-primary-window-minutes"),
		SecondaryUsedPercent:       getFloat("x-codex-secondary-used-percent"),
		SecondaryResetAfterSeconds: getInt("x-codex-secondary-reset-after-seconds"),
		SecondaryWindowMinutes:     getInt("x-codex-secondary-window-minutes"),
	}
}

// SyncAllOpenAIAccounts 同步所有 OpenAI 账户用量
func (s *UsageSyncService) SyncAllOpenAIAccounts() (synced int, failed int, err error) {
	accounts, err := s.accountRepo.GetEnabledByType(model.AccountTypeOpenAI)
	if err != nil {
		return 0, 0, fmt.Errorf("获取账户列表失败: %w", err)
	}

	if len(accounts) == 0 {
		s.log.Info("没有找到 OpenAI 账户")
		return 0, 0, nil
	}

	s.log.Info("开始同步 %d 个 OpenAI 账户 Codex 用量", len(accounts))

	for _, acc := range accounts {
		_, err := s.FetchOpenAICodexUsage(acc.ID)
		if err != nil {
			s.log.Error("账户 %d 用量查询失败: %v", acc.ID, err)
			failed++
		} else {
			synced++
		}
	}

	s.log.Info("OpenAI Codex 用量同步完成: 成功 %d, 失败 %d", synced, failed)

	return synced, failed, nil
}

// ==================== Gemini Usage 查询 ====================

// FetchGeminiUsage 查询 Gemini 账户用量
// 调用 Google Code Assist API 获取项目配置
func (s *UsageSyncService) FetchGeminiUsage(accountID uint) (*repository.GeminiUsageData, error) {
	// 检查缓存
	if lastTime, ok := s.cache.Load(fmt.Sprintf("gemini_%d", accountID)); ok {
		if time.Since(lastTime.(time.Time)) < s.cacheTTL {
			s.log.Debug("Gemini 账户 %d 用量在缓存有效期内，跳过查询", accountID)
			return nil, nil
		}
	}

	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("获取账户失败: %w", err)
	}

	// 检查是否有 AccessToken
	if account.AccessToken == "" {
		s.log.Debug("账户 %d (%s) 无 AccessToken，跳过", accountID, account.Name)
		return nil, nil
	}

	// 检查并刷新 Token（如果需要）
	if account.TokenExpiry != nil && account.TokenExpiry.Before(time.Now()) {
		s.log.Warn("账户 %d (%s) Token 过期，尝试刷新", accountID, account.Name)
		if err := s.refreshAccountToken(account); err != nil {
			return nil, fmt.Errorf("Token 刷新失败: %w", err)
		}
		account, _ = s.accountRepo.GetByID(accountID) // 重新获取
	}

	client := adapter.GetHTTPClient(account)

	// 无固定项目时，预热 tokeninfo/userinfo 以提高项目 ID 获取成功率
	if account.GeminiProjectID == "" {
		tokenInfoBody := url.Values{"access_token": []string{account.AccessToken}}.Encode()
		tokenInfoReq, err := http.NewRequest("POST", "https://oauth2.googleapis.com/tokeninfo", bytes.NewBufferString(tokenInfoBody))
		if err == nil {
			tokenInfoReq.Header.Set("Authorization", "Bearer "+account.AccessToken)
			tokenInfoReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			tokenInfoReq.Header.Set("Accept", "application/json")
			if resp, err := client.Do(tokenInfoReq); err != nil {
				s.log.Warn("tokeninfo 调用失败: %v", err)
			} else {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		}

		userInfoReq, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
		if err == nil {
			userInfoReq.Header.Set("Authorization", "Bearer "+account.AccessToken)
			userInfoReq.Header.Set("Accept", "application/json")
			if resp, err := client.Do(userInfoReq); err != nil {
				s.log.Warn("userinfo 调用失败: %v", err)
			} else {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		}
	}

	// 构建请求
	reqBody := map[string]interface{}{
		"metadata": map[string]string{
			"ideType":    "IDE_UNSPECIFIED",
			"platform":   "PLATFORM_UNSPECIFIED",
			"pluginType": "GEMINI",
		},
	}

	// 如果有固定项目 ID，添加到请求
	if account.GeminiProjectID != "" {
		reqBody["cloudaicompanionProject"] = account.GeminiProjectID
	}

	bodyBytes, _ := json.Marshal(reqBody)
	createRequest := func(accessToken string) (*http.Request, error) {
		req, err := http.NewRequest("POST",
			"https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist",
			bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	}

	req, err := createRequest(account.AccessToken)
	if err != nil {
		return nil, err
	}

	s.log.Debug("查询 Gemini 账户 %d (%s) 用量", accountID, account.Name)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		s.log.Warn("账户 %d (%s) Token 过期，尝试刷新", accountID, account.Name)
		if err := s.refreshAccountToken(account); err != nil {
			return nil, fmt.Errorf("Token 刷新失败: %w", err)
		}
		account, _ = s.accountRepo.GetByID(accountID)
		req, err = createRequest(account.AccessToken)
		if err != nil {
			return nil, err
		}
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("重试请求失败: %w", err)
		}
		defer resp.Body.Close()
	}

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Gemini API 返回错误: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		CloudaicompanionProject string `json:"cloudaicompanionProject"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	// 更新临时项目 ID（仅当没有固定项目 ID 时）
	if account.GeminiProjectID == "" && result.CloudaicompanionProject != "" {
		if err := s.accountRepo.UpdateGeminiUsage(accountID, &repository.GeminiUsageData{
			TempProjectID: result.CloudaicompanionProject,
		}); err != nil {
			return nil, fmt.Errorf("更新数据库失败: %w", err)
		}
	}

	// 更新缓存
	s.cache.Store(fmt.Sprintf("gemini_%d", accountID), time.Now())

	s.log.Info("Gemini 账户 %d (%s) 用量更新成功，项目 ID: %s", accountID, account.Name, result.CloudaicompanionProject)

	return &repository.GeminiUsageData{
		ProjectID:     account.GeminiProjectID,
		TempProjectID: result.CloudaicompanionProject,
	}, nil
}

// SyncAllGeminiAccounts 同步所有 Gemini 账户用量
func (s *UsageSyncService) SyncAllGeminiAccounts() (synced int, failed int, err error) {
	accounts, err := s.accountRepo.GetEnabledByType(model.AccountTypeGemini)
	if err != nil {
		return 0, 0, fmt.Errorf("获取账户列表失败: %w", err)
	}

	if len(accounts) == 0 {
		s.log.Info("没有找到 Gemini 账户")
		return 0, 0, nil
	}

	s.log.Info("开始同步 %d 个 Gemini 账户用量", len(accounts))

	for _, acc := range accounts {
		_, err := s.FetchGeminiUsage(acc.ID)
		if err != nil {
			s.log.Error("账户 %d 用量查询失败: %v", acc.ID, err)
			failed++
		} else {
			synced++
		}
	}

	s.log.Info("Gemini 用量同步完成: 成功 %d, 失败 %d", synced, failed)

	return synced, failed, nil
}

// ==================== 定时同步 ====================

// StartAutoSync 启动自动同步（可配置间隔）
func (s *UsageSyncService) StartAutoSync(interval time.Duration) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		s.log.Warn("用量同步定时器已在运行")
		return
	}
	s.running = true
	s.ticker = time.NewTicker(interval)
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	go func() {
		s.log.Info("用量同步定时器已启动，间隔: %v", interval)
		for {
			select {
			case <-s.ticker.C:
				s.log.Info("执行定时用量同步")

				// Claude 同步
				claudeSynced, claudeFailed, _ := s.SyncAllClaudeAccounts()
				s.log.Info("Claude: 成功 %d, 失败 %d", claudeSynced, claudeFailed)

				// OpenAI 同步
				openaiSynced, openaiFailed, _ := s.SyncAllOpenAIAccounts()
				s.log.Info("OpenAI: 成功 %d, 失败 %d", openaiSynced, openaiFailed)

				// Gemini 同步
				geminiSynced, geminiFailed, _ := s.SyncAllGeminiAccounts()
				s.log.Info("Gemini: 成功 %d, 失败 %d", geminiSynced, geminiFailed)

			case <-s.stopChan:
				s.ticker.Stop()
				s.log.Info("用量同步定时器已停止")
				return
			}
		}
	}()
}

// StopAutoSync 停止自动同步
func (s *UsageSyncService) StopAutoSync() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stopChan)
	s.running = false
}

// ==================== 工具函数 ====================

// floatPtr 返回 float64 指针
func floatPtr(f float64) *float64 {
	return &f
}

// refreshAccountToken 统一刷新 OAuth Token
func (s *UsageSyncService) refreshAccountToken(account *model.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := scheduler.GetTokenManager().ForceRefresh(ctx, account.ID); err == nil {
		return nil
	} else if account.Platform == model.PlatformClaude && account.SessionKey != "" {
		return s.healthCheck.RefreshToken(account.ID)
	} else {
		return err
	}
}
