/*
 * 文件作用：OAuth Token 管理器，处理 Claude/OpenAI/xyrt OAuth Token 刷新
 * 负责功能：
 *   - Access Token 自动刷新
 *   - Token 过期检测
 *   - 刷新锁防止并发刷新
 *   - Token 持久化更新
 *   - xyrt Token 每日定时刷新
 * 重要程度：⭐⭐⭐⭐ 重要（OAuth账户必需）
 * 依赖模块：model, repository
 */
package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/proxy/adapter"
	"cli-proxy/internal/repository"
	"cli-proxy/pkg/logger"
	"cli-proxy/pkg/utils"
)

// TokenManager OAuth Token 管理器
type TokenManager struct {
	mu   sync.RWMutex
	repo *repository.AccountRepository

	// Token 刷新配置
	refreshThreshold time.Duration // 提前刷新时间
	refreshing       map[uint]bool // 正在刷新的账户
}

var defaultTokenManager *TokenManager
var tokenManagerOnce sync.Once

// GetTokenManager 获取 Token 管理器单例
func GetTokenManager() *TokenManager {
	tokenManagerOnce.Do(func() {
		defaultTokenManager = &TokenManager{
			repo:             repository.NewAccountRepository(),
			refreshThreshold: 5 * time.Minute, // 提前5分钟刷新
			refreshing:       make(map[uint]bool),
		}
		// 启动后台刷新协程
		go defaultTokenManager.backgroundRefresh()
		// 启动 xyrt 每日刷新协程
		go defaultTokenManager.dailyXyrtRefresh()
	})
	return defaultTokenManager
}

// SetRefreshThreshold 设置提前刷新时间
func (m *TokenManager) SetRefreshThreshold(d time.Duration) {
	m.refreshThreshold = d
}

// CheckAndRefreshToken 检查并刷新 Token
func (m *TokenManager) CheckAndRefreshToken(ctx context.Context, account *model.Account) error {
	// 只处理有 Token 过期时间的账户
	if account.TokenExpiry == nil {
		return nil
	}

	// 检查是否需要刷新
	if time.Until(*account.TokenExpiry) > m.refreshThreshold {
		return nil
	}

	// 检查是否正在刷新
	m.mu.Lock()
	if m.refreshing[account.ID] {
		m.mu.Unlock()
		return nil
	}
	m.refreshing[account.ID] = true
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.refreshing, account.ID)
		m.mu.Unlock()
	}()

	// 根据账户类型刷新
	switch account.Type {
	case model.AccountTypeClaudeOfficial:
		return m.refreshClaudeOfficialToken(ctx, account)
	case model.AccountTypeOpenAI, model.AccountTypeOpenAIResponses:
		return m.refreshOpenAIToken(ctx, account)
	case model.AccountTypeGemini:
		return m.refreshGeminiToken(ctx, account)
	default:
		return nil
	}
}

// refreshClaudeOfficialToken 刷新 Claude Official Token
func (m *TokenManager) refreshClaudeOfficialToken(ctx context.Context, account *model.Account) error {
	if account.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// Claude OAuth 刷新端点 - 使用官方 Console API
	tokenURL := "https://console.anthropic.com/v1/oauth/token"

	// 构建 JSON 请求体（与 clove 项目保持一致）
	payload := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": account.RefreshToken,
		"client_id":     "9d1c250a-e61b-44d9-88ed-5944d1962f5e", // OAuth Client ID
	}
	jsonBody, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := utils.ReadAllWithLimit(resp.Body, utils.MaxResponseBodyBytes)
		return fmt.Errorf("token refresh failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	// 更新账户
	expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	account.AccessToken = tokenResp.AccessToken
	account.TokenExpiry = &expiry
	if tokenResp.RefreshToken != "" {
		account.RefreshToken = tokenResp.RefreshToken
	}

	return m.repo.UpdateToken(account.ID, account.AccessToken, account.RefreshToken, &expiry)
}

// refreshGeminiToken 刷新 Gemini Token
func (m *TokenManager) refreshGeminiToken(ctx context.Context, account *model.Account) error {
	if account.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// Google OAuth 刷新端点
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", account.RefreshToken)
	// 需要 client_id 和 client_secret，从 APIKey 和 APISecret 获取
	if account.APIKey != "" {
		data.Set("client_id", account.APIKey)
	}
	if account.APISecret != "" {
		data.Set("client_secret", account.APISecret)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := utils.ReadAllWithLimit(resp.Body, utils.MaxResponseBodyBytes)
		return fmt.Errorf("token refresh failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	// 更新账户
	expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	account.AccessToken = tokenResp.AccessToken
	account.TokenExpiry = &expiry

	return m.repo.UpdateToken(account.ID, account.AccessToken, account.RefreshToken, &expiry)
}

// refreshOpenAIToken 刷新 OpenAI OAuth Token
func (m *TokenManager) refreshOpenAIToken(ctx context.Context, account *model.Account) error {
	if account.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	const openAITokenURL = "https://auth.openai.com/oauth/token"
	const openAIClientID = "app_EMoamEEZ73f0CkXaXp7hrann"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", openAIClientID)
	data.Set("refresh_token", account.RefreshToken)

	req, err := http.NewRequestWithContext(ctx, "POST", openAITokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := adapter.GetChromeTLSClient(account)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := utils.ReadAllWithLimit(resp.Body, utils.MaxResponseBodyBytes)
		return fmt.Errorf("token refresh failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return err
	}

	expiry := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	account.AccessToken = tokenResp.AccessToken
	account.TokenExpiry = &expiry
	if tokenResp.RefreshToken != "" {
		account.RefreshToken = tokenResp.RefreshToken
	}

	return m.repo.UpdateToken(account.ID, account.AccessToken, account.RefreshToken, &expiry)
}

// backgroundRefresh 后台定期检查并刷新 Token
func (m *TokenManager) backgroundRefresh() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.refreshAllExpiring()
	}
}

// refreshAllExpiring 刷新所有即将过期的 Token
func (m *TokenManager) refreshAllExpiring() {
	// 获取所有启用的账户
	accounts, err := m.repo.GetAllEnabled()
	if err != nil {
		return
	}

	ctx := context.Background()
	for _, account := range accounts {
		if account.TokenExpiry != nil && time.Until(*account.TokenExpiry) < m.refreshThreshold {
			go m.CheckAndRefreshToken(ctx, &account)
		}
	}
}

// ForceRefresh 强制刷新指定账户的 Token
func (m *TokenManager) ForceRefresh(ctx context.Context, accountID uint) error {
	account, err := m.repo.GetByID(accountID)
	if err != nil {
		return err
	}

	switch account.Type {
	case model.AccountTypeClaudeOfficial:
		return m.refreshClaudeOfficialToken(ctx, account)
	case model.AccountTypeOpenAI, model.AccountTypeOpenAIResponses:
		// xyrt 类型使用专门的刷新方法
		if account.AuthType == "xyrt" {
			return m.refreshXyrtToken(ctx, account)
		}
		return m.refreshOpenAIToken(ctx, account)
	case model.AccountTypeGemini:
		return m.refreshGeminiToken(ctx, account)
	default:
		return fmt.Errorf("account type %s does not support token refresh", account.Type)
	}
}

// ========== xyrt Token 刷新相关 ==========

// getXyrtLog 懒加载获取 xyrt 日志器（避免包初始化顺序问题）
func getXyrtLog() *logger.Logger {
	return logger.GetLogger("xyrt")
}

// refreshXyrtToken 刷新 xyrt Token
func (m *TokenManager) refreshXyrtToken(ctx context.Context, account *model.Account) error {
	if account.XyrtRefreshToken == "" {
		return fmt.Errorf("no xyrt refresh token available")
	}

	if account.GatewayURL == "" {
		return fmt.Errorf("no gateway URL configured for xyrt")
	}

	// 构建刷新 URL
	refreshURL := strings.TrimSuffix(account.GatewayURL, "/") + "/auth/refresh"

	// 构建请求体
	payload := map[string]string{
		"refresh_token": account.XyrtRefreshToken,
	}
	jsonBody, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", refreshURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// xyrt 不使用代理，直接访问网关
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		getXyrtLog().Error("[xyrt] Token 刷新失败 | AccountID: %d | 原因: %v", account.ID, err)
		m.repo.MarkAsTokenExpired(account.ID, fmt.Sprintf("xyrt refresh failed: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := utils.ReadAllWithLimit(resp.Body, utils.MaxResponseBodyBytes)
		errMsg := fmt.Sprintf("xyrt token refresh failed: HTTP %d - %s", resp.StatusCode, string(body))
		getXyrtLog().Error("[xyrt] Token 刷新失败 | AccountID: %d | %s", account.ID, errMsg)
		m.repo.MarkAsTokenExpired(account.ID, errMsg)
		return fmt.Errorf(errMsg)
	}

	// 解析响应
	var tokenResp struct {
		AccessToken      string `json:"accessToken"`
		AccountCheckInfo struct {
			TeamIDs  []string `json:"team_ids"`
			PlanType string   `json:"plan_type"`
		} `json:"accountCheckInfo"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		errMsg := fmt.Sprintf("xyrt response parse failed: %v", err)
		getXyrtLog().Error("[xyrt] Token 刷新失败 | AccountID: %d | %s", account.ID, errMsg)
		m.repo.MarkAsTokenExpired(account.ID, errMsg)
		return err
	}

	// 更新账户
	now := time.Now()
	orgID := ""
	planType := tokenResp.AccountCheckInfo.PlanType

	// 如果是 team 或 k12，设置组织 ID
	if (planType == "team" || planType == "k12") && len(tokenResp.AccountCheckInfo.TeamIDs) > 0 {
		orgID = tokenResp.AccountCheckInfo.TeamIDs[0]
	}

	// 更新数据库
	if err := m.repo.UpdateXyrtToken(account.ID, tokenResp.AccessToken, orgID, planType, &now); err != nil {
		getXyrtLog().Error("[xyrt] Token 更新失败 | AccountID: %d | 原因: %v", account.ID, err)
		return err
	}

	getXyrtLog().Info("[xyrt] Token 刷新成功 | AccountID: %d | PlanType: %s | OrgID: %s", account.ID, planType, orgID)
	return nil
}

// dailyXyrtRefresh 每日定时刷新 xyrt Token
func (m *TokenManager) dailyXyrtRefresh() {
	// 每小时检查一次
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// 启动时立即检查一次
	m.refreshExpiredXyrtTokens()

	for range ticker.C {
		m.refreshExpiredXyrtTokens()
	}
}

// refreshExpiredXyrtTokens 刷新超过24小时未刷新的 xyrt Token
func (m *TokenManager) refreshExpiredXyrtTokens() {
	accounts, err := m.repo.GetXyrtAccountsNeedingRefresh()
	if err != nil {
		getXyrtLog().Error("[xyrt] 获取需要刷新的账户失败: %v", err)
		return
	}

	if len(accounts) == 0 {
		return
	}

	getXyrtLog().Info("[xyrt] 开始刷新 %d 个 xyrt 账户", len(accounts))

	ctx := context.Background()
	for _, account := range accounts {
		// 检查是否正在刷新
		m.mu.Lock()
		if m.refreshing[account.ID] {
			m.mu.Unlock()
			continue
		}
		m.refreshing[account.ID] = true
		m.mu.Unlock()

		go func(acc model.Account) {
			defer func() {
				m.mu.Lock()
				delete(m.refreshing, acc.ID)
				m.mu.Unlock()
			}()
			m.refreshXyrtToken(ctx, &acc)
		}(account)
	}
}
