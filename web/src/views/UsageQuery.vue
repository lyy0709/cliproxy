<!--
 * 文件作用：API Key 用量查询页面 - Apple HIG 风格
 * 负责功能：
 *   - 输入 API Key 查询使用统计
 *   - 展示总计/今日用量
 *   - 展示配额信息
 *   - 展示模型使用分布
 *   - 展示每日趋势
 * 重要程度：⭐⭐⭐ 一般（公开查询页面）
-->
<template>
  <div class="usage-query-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="bg-circle bg-circle-1"></div>
      <div class="bg-circle bg-circle-2"></div>
    </div>

    <!-- 查询卡片 -->
    <div class="query-card">
      <!-- 头部 -->
      <div class="query-header">
        <div class="header-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
        </div>
        <h1 class="header-title">用量查询</h1>
        <p class="header-subtitle">输入 API Key 查看使用统计</p>
      </div>

      <!-- 查询表单 -->
      <div class="query-form">
        <div class="input-group">
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
            </svg>
            <input
              v-model="apiKey"
              type="password"
              class="form-input"
              placeholder="请输入 API Key"
              @keyup.enter="handleQuery"
            />
          </div>
          <button class="btn-query" :disabled="loading" @click="handleQuery">
            <svg v-if="loading" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
              <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
            </svg>
            <span v-else>查询</span>
          </button>
        </div>

        <!-- 错误提示 -->
        <div v-if="errorMsg" class="error-alert">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          <span>{{ errorMsg }}</span>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!usageData && !loading && !errorMsg" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 80 80" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="10" y="20" width="60" height="40" rx="4"/>
            <path d="M10 30h60M25 40h30M25 50h20" stroke-linecap="round"/>
          </svg>
        </div>
        <p class="empty-text">输入 API Key 开始查询</p>
      </div>
    </div>

    <!-- 结果区域 -->
    <div v-if="usageData" class="result-container">
      <!-- Key 信息头 -->
      <div class="key-info-card">
        <div class="key-info-header">
          <div class="key-name">
            <h2>{{ usageData.key_name || 'API Key' }}</h2>
            <span :class="['status-badge', isExpired ? 'expired' : 'active']">
              {{ isExpired ? '已过期' : '有效' }}
            </span>
          </div>
        </div>
      </div>

      <!-- 统计卡片网格 -->
      <div class="stats-section">
        <div class="section-header">
          <h3>总计</h3>
        </div>
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-icon requests">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
              </svg>
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ formatNumber(usageData.total_requests) }}</span>
              <span class="stat-label">请求数</span>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon tokens">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 6v12M6 12h12"/>
              </svg>
            </div>
            <div class="stat-content">
              <span class="stat-value">{{ formatTokens(usageData.total_tokens) }}</span>
              <span class="stat-label">Token</span>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon cost">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 6v2m0 8v2M9 10c0-1.1.9-2 2-2h2a2 2 0 010 4h-2a2 2 0 000 4h2a2 2 0 002-2"/>
              </svg>
            </div>
            <div class="stat-content">
              <span class="stat-value">${{ usageData.total_cost?.toFixed(4) || '0.00' }}</span>
              <span class="stat-label">费用</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 今日统计 -->
      <div class="stats-section">
        <div class="section-header">
          <h3>今日</h3>
        </div>
        <div class="stats-grid compact">
          <div class="stat-card mini">
            <span class="stat-value-sm">{{ formatNumber(usageData.today_requests) }}</span>
            <span class="stat-label">请求</span>
          </div>
          <div class="stat-card mini">
            <span class="stat-value-sm">{{ formatTokens(usageData.today_tokens) }}</span>
            <span class="stat-label">Token</span>
          </div>
          <div class="stat-card mini">
            <span class="stat-value-sm">${{ usageData.today_cost?.toFixed(4) || '0.00' }}</span>
            <span class="stat-label">费用</span>
          </div>
        </div>
      </div>

      <!-- 配额信息 -->
      <div v-if="usageData.daily_limit || usageData.monthly_quota || usageData.expires_at" class="quota-section">
        <div class="section-header">
          <h3>配额</h3>
        </div>
        <div class="quota-card">
          <div v-if="usageData.daily_limit" class="quota-item">
            <span class="quota-label">每日限额</span>
            <span class="quota-value">{{ usageData.daily_limit }} 次</span>
          </div>
          <div v-if="usageData.monthly_quota" class="quota-item">
            <span class="quota-label">月度配额</span>
            <span class="quota-value">${{ usageData.monthly_quota }}</span>
          </div>
          <div v-if="usageData.expires_at" class="quota-item">
            <span class="quota-label">过期时间</span>
            <span class="quota-value">{{ formatDate(usageData.expires_at) }}</span>
          </div>
        </div>
      </div>

      <!-- 模型使用分布 -->
      <div v-if="usageData.model_stats?.length" class="model-section">
        <div class="section-header">
          <h3>模型分布</h3>
        </div>
        <div class="table-card">
          <table class="data-table">
            <thead>
              <tr>
                <th class="text-left">模型</th>
                <th class="text-right">请求</th>
                <th class="text-right">Token</th>
                <th class="text-right">费用</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="model in usageData.model_stats" :key="model.model">
                <td class="text-left">
                  <span class="model-name">{{ model.model }}</span>
                </td>
                <td class="text-right">{{ formatNumber(model.request_count) }}</td>
                <td class="text-right">{{ formatTokens(model.total_tokens) }}</td>
                <td class="text-right">${{ model.total_cost?.toFixed(4) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 每日趋势 -->
      <div class="trend-section">
        <div class="section-header">
          <h3>每日趋势</h3>
          <button class="btn-refresh" @click="loadDailyStats" :disabled="loadingDaily">
            <svg v-if="loadingDaily" class="spinner-sm" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
              <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
            </svg>
            <span v-else>刷新</span>
          </button>
        </div>
        <div v-if="dailyStats.length" class="table-card">
          <table class="data-table">
            <thead>
              <tr>
                <th class="text-left">日期</th>
                <th class="text-right">请求</th>
                <th class="text-right">Token</th>
                <th class="text-right">费用</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="day in dailyStats.slice(0, 14)" :key="day.date">
                <td class="text-left">{{ day.date }}</td>
                <td class="text-right">{{ formatNumber(day.request_count) }}</td>
                <td class="text-right">{{ formatTokens(day.total_tokens) }}</td>
                <td class="text-right">${{ day.total_cost?.toFixed(4) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-else class="empty-table">
          <p>暂无数据</p>
        </div>
      </div>

      <!-- 返回首页 -->
      <div class="footer-link">
        <router-link to="/" class="back-link">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="19" y1="12" x2="5" y2="12"/>
            <polyline points="12,19 5,12 12,5"/>
          </svg>
          返回首页
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * API Key Usage Query Page - Apple HIG Style
 *
 * Design Features:
 * - Uses global --apple-* CSS variables
 * - SF Pro system font stack
 * - 44pt minimum touch targets
 * - Supports prefers-color-scheme
 * - Card-based layout
 */
import { ref, computed } from 'vue'

const apiKey = ref('')
const loading = ref(false)
const loadingDaily = ref(false)
const usageData = ref(null)
const dailyStats = ref([])
const errorMsg = ref('')

const isExpired = computed(() => {
  if (!usageData.value?.expires_at) return false
  return new Date(usageData.value.expires_at) < new Date()
})

async function handleQuery() {
  if (!apiKey.value.trim()) {
    errorMsg.value = '请输入 API Key'
    return
  }

  loading.value = true
  errorMsg.value = ''
  usageData.value = null
  dailyStats.value = []

  try {
    const response = await fetch('/api/key/usage', {
      headers: { 'x-api-key': apiKey.value.trim() }
    })
    const data = await response.json()

    if (!response.ok) {
      throw new Error(data.message || data.error || '查询失败')
    }

    usageData.value = data.data
    await loadDailyStats()
  } catch (err) {
    errorMsg.value = err.message || '查询失败，请检查 API Key'
  } finally {
    loading.value = false
  }
}

async function loadDailyStats() {
  if (!apiKey.value.trim()) return

  loadingDaily.value = true
  try {
    const response = await fetch('/api/key/usage/daily', {
      headers: { 'x-api-key': apiKey.value.trim() }
    })
    const data = await response.json()

    if (response.ok && data.data?.stats) {
      dailyStats.value = [...data.data.stats].sort((a, b) =>
        new Date(b.date) - new Date(a.date)
      )
    }
  } catch (err) {
    console.error('加载每日统计失败:', err)
  } finally {
    loadingDaily.value = false
  }
}

function formatNumber(num) {
  if (!num) return '0'
  return num.toLocaleString()
}

function formatTokens(tokens) {
  if (!tokens) return '0'
  if (tokens >= 1000000) return (tokens / 1000000).toFixed(2) + 'M'
  if (tokens >= 1000) return (tokens / 1000).toFixed(1) + 'K'
  return tokens.toString()
}

function formatDate(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}
</script>

<style scoped>
.usage-query-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: var(--apple-spacing-xl);
  position: relative;
  overflow-x: hidden;
}

/* 背景装饰 */
.bg-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.bg-circle-1 {
  width: 500px;
  height: 500px;
  top: -150px;
  right: -100px;
}

.bg-circle-2 {
  width: 300px;
  height: 300px;
  bottom: -50px;
  left: -50px;
}

/* 查询卡片 */
.query-card {
  max-width: 520px;
  margin: 0 auto var(--apple-spacing-xl);
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xxl);
  box-shadow: var(--apple-shadow-xl);
  padding: var(--apple-spacing-xxl);
  position: relative;
  z-index: 1;
  animation: apple-scale-in var(--apple-duration-normal) var(--apple-ease-spring);
}

/* 头部 */
.query-header {
  text-align: center;
  margin-bottom: var(--apple-spacing-xl);
}

.header-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto var(--apple-spacing-md);
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  border-radius: var(--apple-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-icon svg {
  width: 32px;
  height: 32px;
  color: white;
}

.header-title {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.header-subtitle {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

/* 查询表单 */
.query-form {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-md);
}

.input-group {
  display: flex;
  gap: var(--apple-spacing-sm);
}

.input-wrapper {
  flex: 1;
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--apple-spacing-md);
  width: 18px;
  height: 18px;
  color: var(--apple-text-tertiary);
  pointer-events: none;
}

.form-input {
  width: 100%;
  padding: var(--apple-spacing-md);
  padding-left: 44px;
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-fill-quaternary);
  border: 1px solid transparent;
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-input::placeholder {
  color: var(--apple-text-placeholder);
}

.form-input:hover {
  background: var(--apple-fill-tertiary);
}

.form-input:focus {
  background: var(--apple-bg-primary);
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.btn-query {
  min-width: 88px;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  color: white;
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-semibold);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-query:hover:not(:disabled) {
  transform: scale(1.02);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
}

.btn-query:active:not(:disabled) {
  transform: scale(0.98);
}

.btn-query:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 错误提示 */
.error-alert {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  background: var(--apple-red-light);
  border-radius: var(--apple-radius-sm);
  color: var(--apple-red);
  font-size: var(--apple-text-sm);
}

.error-alert svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: var(--apple-spacing-xxl) 0;
}

.empty-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto var(--apple-spacing-md);
  color: var(--apple-text-quaternary);
}

.empty-text {
  font-size: var(--apple-text-base);
  color: var(--apple-text-tertiary);
}

/* 结果容器 */
.result-container {
  max-width: 680px;
  margin: 0 auto;
  position: relative;
  z-index: 1;
}

/* Key 信息卡片 */
.key-info-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-lg);
  animation: apple-fade-in-up var(--apple-duration-normal) var(--apple-ease-default);
}

.key-name {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.key-name h2 {
  font-size: var(--apple-text-xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin: 0;
}

.status-badge {
  padding: var(--apple-spacing-xxs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
  border-radius: var(--apple-radius-full);
}

.status-badge.active {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.status-badge.expired {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

/* 统计区块 */
.stats-section,
.quota-section,
.model-section,
.trend-section {
  margin-bottom: var(--apple-spacing-lg);
  animation: apple-fade-in-up var(--apple-duration-normal) var(--apple-ease-default);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--apple-spacing-sm);
  padding: 0 var(--apple-spacing-xs);
}

.section-header h3 {
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
  color: rgba(255, 255, 255, 0.7);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0;
}

/* 统计网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--apple-spacing-sm);
}

.stat-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-md);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.stat-card.mini {
  flex-direction: column;
  text-align: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-sm);
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 20px;
  height: 20px;
}

.stat-icon.requests {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.stat-icon.tokens {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.stat-icon.cost {
  background: var(--apple-orange-light);
  color: var(--apple-orange);
}

.stat-content {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.stat-value {
  font-size: var(--apple-text-xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  line-height: 1.2;
}

.stat-value-sm {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.stat-label {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
  margin-top: 2px;
}

/* 配额卡片 */
.quota-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.quota-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 44px;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  border-bottom: 1px solid var(--apple-separator);
}

.quota-item:last-child {
  border-bottom: none;
}

.quota-label {
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
}

.quota-value {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
}

/* 表格卡片 */
.table-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--apple-text-sm);
}

.data-table th {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  border-bottom: 1px solid var(--apple-separator);
  background: var(--apple-bg-secondary);
}

.data-table td {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  color: var(--apple-text-primary);
  border-bottom: 1px solid var(--apple-separator);
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table .text-left { text-align: left; }
.data-table .text-right { text-align: right; }

.model-name {
  font-weight: var(--apple-font-medium);
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
}

/* 刷新按钮 */
.btn-refresh {
  padding: var(--apple-spacing-xxs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: rgba(255, 255, 255, 0.9);
  background: rgba(255, 255, 255, 0.15);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
}

.btn-refresh:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.25);
}

.btn-refresh:disabled {
  opacity: 0.5;
}

/* 空表格 */
.empty-table {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-xxl);
  text-align: center;
  color: var(--apple-text-tertiary);
}

/* 底部链接 */
.footer-link {
  text-align: center;
  padding: var(--apple-spacing-xl) 0;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  font-size: var(--apple-text-sm);
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: color var(--apple-duration-fast) var(--apple-ease-default);
}

.back-link:hover {
  color: white;
}

.back-link svg {
  width: 16px;
  height: 16px;
}

/* Spinner */
.spinner {
  width: 20px;
  height: 20px;
  animation: apple-spin 1s linear infinite;
}

.spinner-sm {
  width: 14px;
  height: 14px;
  animation: apple-spin 1s linear infinite;
}

/* 响应式 */
@media (max-width: 640px) {
  .usage-query-page {
    padding: var(--apple-spacing-md);
  }

  .query-card {
    padding: var(--apple-spacing-xl);
  }

  .input-group {
    flex-direction: column;
  }

  .btn-query {
    width: 100%;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stat-card {
    flex-direction: row;
  }

  .stat-card.mini {
    flex-direction: row;
    text-align: left;
    justify-content: space-between;
  }

  .stat-card.mini .stat-label {
    order: -1;
    margin-top: 0;
    font-size: var(--apple-text-sm);
  }

  .header-title {
    font-size: var(--apple-text-xl);
  }
}

/* 减少动画偏好 */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
</style>
