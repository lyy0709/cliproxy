<template>
  <div class="usage-query">
    <!-- 查询区域 -->
    <div class="query-section">
      <h1 class="large-title">用量查询</h1>
      <p class="subtitle">输入 API Key 查看使用统计</p>

      <div class="input-group">
        <input
          v-model="apiKey"
          type="password"
          class="input-field"
          placeholder="请输入 API Key"
          @keyup.enter="handleQuery"
        />
        <button
          class="btn-primary"
          :disabled="loading"
          @click="handleQuery"
        >
          <span v-if="loading" class="loading-spinner"></span>
          <span v-else>查询</span>
        </button>
      </div>

      <!-- 错误提示 -->
      <div v-if="errorMsg" class="alert-inline">
        <span class="alert-icon">!</span>
        <span>{{ errorMsg }}</span>
      </div>
    </div>

    <!-- 结果区域 -->
    <div v-if="usageData" class="result-section">
      <!-- Key 信息头 -->
      <div class="key-header">
        <div class="key-info">
          <h2 class="title2">{{ usageData.key_name || 'API Key' }}</h2>
          <span v-if="usageData.expires_at" :class="['status-badge', isExpired ? 'expired' : 'active']">
            {{ isExpired ? '已过期' : '有效' }}
          </span>
        </div>
      </div>

      <!-- 总统计卡片 -->
      <div class="card-section">
        <h3 class="section-header">总计</h3>
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
      <div class="card-section">
        <h3 class="section-header">今日</h3>
        <div class="stats-grid">
          <div class="stat-card compact">
            <span class="stat-value-sm">{{ formatNumber(usageData.today_requests) }}</span>
            <span class="stat-label">请求</span>
          </div>
          <div class="stat-card compact">
            <span class="stat-value-sm">{{ formatTokens(usageData.today_tokens) }}</span>
            <span class="stat-label">Token</span>
          </div>
          <div class="stat-card compact">
            <span class="stat-value-sm">${{ usageData.today_cost?.toFixed(4) || '0.00' }}</span>
            <span class="stat-label">费用</span>
          </div>
        </div>
      </div>

      <!-- 配额信息 -->
      <div v-if="usageData.daily_limit || usageData.monthly_quota || usageData.expires_at" class="card-section">
        <h3 class="section-header">配额</h3>
        <div class="list-group">
          <div v-if="usageData.daily_limit" class="list-item">
            <span class="list-label">每日限额</span>
            <span class="list-value">{{ usageData.daily_limit }} 次</span>
          </div>
          <div v-if="usageData.monthly_quota" class="list-item">
            <span class="list-label">月度配额</span>
            <span class="list-value">${{ usageData.monthly_quota }}</span>
          </div>
          <div v-if="usageData.expires_at" class="list-item">
            <span class="list-label">过期时间</span>
            <span class="list-value">{{ formatDate(usageData.expires_at) }}</span>
          </div>
        </div>
      </div>

      <!-- 模型使用分布 -->
      <div v-if="usageData.model_stats?.length" class="card-section">
        <h3 class="section-header">模型分布</h3>
        <div class="table-container">
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
                <td class="text-left model-name">{{ model.model }}</td>
                <td class="text-right">{{ formatNumber(model.request_count) }}</td>
                <td class="text-right">{{ formatTokens(model.total_tokens) }}</td>
                <td class="text-right">${{ model.total_cost?.toFixed(4) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 每日趋势 -->
      <div class="card-section">
        <div class="section-header-row">
          <h3 class="section-header">每日趋势</h3>
          <button class="btn-text" @click="loadDailyStats" :disabled="loadingDaily">
            {{ loadingDaily ? '刷新中...' : '刷新' }}
          </button>
        </div>
        <div v-if="dailyStats.length" class="table-container">
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
        <div v-else class="empty-state">
          <p>暂无数据</p>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!loading && !errorMsg" class="empty-hero">
      <div class="empty-icon">
        <svg viewBox="0 0 80 80" fill="none">
          <rect x="10" y="20" width="60" height="40" rx="4" stroke="currentColor" stroke-width="2"/>
          <path d="M10 30h60M25 40h30M25 50h20" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
      </div>
      <p class="empty-text">输入 API Key 开始查询</p>
    </div>
  </div>
</template>

<script setup>
/**
 * Apple HIG Compliant Usage Query Page
 *
 * Design Decisions:
 * - Uses SF Pro system font stack for native feel
 * - 44pt minimum touch targets for accessibility
 * - Grouped list style for data presentation
 * - Supports prefers-color-scheme for auto theming
 * - 8pt grid spacing system
 * - Concentric border radius relationships
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
/* ========== Apple HIG Design System ========== */

/* CSS Custom Properties - Apple System Colors */
:root {
  --font-system: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', Arial, sans-serif;

  /* Light Mode Colors */
  --system-blue: #007AFF;
  --system-green: #34C759;
  --system-orange: #FF9500;
  --system-red: #FF3B30;
  --system-purple: #AF52DE;
  --system-teal: #5AC8FA;

  --label-primary: #000000;
  --label-secondary: rgba(60, 60, 67, 0.6);
  --label-tertiary: rgba(60, 60, 67, 0.3);

  --bg-primary: #FFFFFF;
  --bg-secondary: #F2F2F7;
  --bg-tertiary: #FFFFFF;

  --separator: rgba(60, 60, 67, 0.12);

  /* Spacing - 8pt Grid */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;

  /* Border Radius */
  --radius-sm: 8px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;

  /* Animation */
  --ease-default: cubic-bezier(0.25, 0.1, 0.25, 1);
  --duration-fast: 200ms;
  --duration-normal: 300ms;
}

/* Dark Mode */
@media (prefers-color-scheme: dark) {
  :root {
    --system-blue: #0A84FF;
    --system-green: #30D158;
    --system-orange: #FF9F0A;
    --system-red: #FF453A;
    --system-purple: #BF5AF2;
    --system-teal: #64D2FF;

    --label-primary: #FFFFFF;
    --label-secondary: rgba(235, 235, 245, 0.6);
    --label-tertiary: rgba(235, 235, 245, 0.3);

    --bg-primary: #000000;
    --bg-secondary: #1C1C1E;
    --bg-tertiary: #2C2C2E;

    --separator: rgba(84, 84, 88, 0.36);
  }
}

/* ========== Layout ========== */

.usage-query {
  min-height: 100vh;
  background: var(--bg-secondary);
  font-family: var(--font-system);
  padding: var(--space-4);
  padding-bottom: calc(var(--space-8) + env(safe-area-inset-bottom));
}

/* ========== Query Section ========== */

.query-section {
  max-width: 600px;
  margin: 0 auto;
  padding: var(--space-8) 0 var(--space-6);
  text-align: center;
}

.large-title {
  font-size: 34px;
  font-weight: 700;
  color: var(--label-primary);
  margin: 0 0 var(--space-2);
  letter-spacing: -0.5px;
}

.subtitle {
  font-size: 17px;
  color: var(--label-secondary);
  margin: 0 0 var(--space-6);
}

.input-group {
  display: flex;
  gap: var(--space-3);
  max-width: 480px;
  margin: 0 auto;
}

.input-field {
  flex: 1;
  min-height: 50px;
  padding: 0 var(--space-4);
  font-family: var(--font-system);
  font-size: 17px;
  color: var(--label-primary);
  background: var(--bg-tertiary);
  border: none;
  border-radius: var(--radius-md);
  outline: none;
  transition: box-shadow var(--duration-fast) var(--ease-default);
}

.input-field:focus {
  box-shadow: 0 0 0 4px rgba(0, 122, 255, 0.25);
}

.input-field::placeholder {
  color: var(--label-tertiary);
}

.btn-primary {
  min-width: 88px;
  min-height: 50px;
  padding: 0 var(--space-5);
  font-family: var(--font-system);
  font-size: 17px;
  font-weight: 600;
  color: #FFFFFF;
  background: var(--system-blue);
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: opacity var(--duration-fast), transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-primary:hover {
  opacity: 0.9;
}

.btn-primary:active {
  transform: scale(0.98);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #FFFFFF;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.alert-inline {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  margin-top: var(--space-4);
  padding: var(--space-3) var(--space-4);
  background: rgba(255, 59, 48, 0.1);
  border-radius: var(--radius-sm);
  color: var(--system-red);
  font-size: 15px;
}

.alert-icon {
  width: 20px;
  height: 20px;
  background: var(--system-red);
  color: #FFFFFF;
  border-radius: 50%;
  font-size: 12px;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ========== Result Section ========== */

.result-section {
  max-width: 600px;
  margin: 0 auto;
}

.key-header {
  padding: var(--space-4) 0;
}

.key-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.title2 {
  font-size: 22px;
  font-weight: 700;
  color: var(--label-primary);
  margin: 0;
}

.status-badge {
  padding: var(--space-1) var(--space-3);
  font-size: 13px;
  font-weight: 600;
  border-radius: var(--radius-sm);
}

.status-badge.active {
  background: rgba(52, 199, 89, 0.15);
  color: var(--system-green);
}

.status-badge.expired {
  background: rgba(255, 59, 48, 0.15);
  color: var(--system-red);
}

/* ========== Card Section ========== */

.card-section {
  margin-bottom: var(--space-6);
}

.section-header {
  font-size: 13px;
  font-weight: 400;
  color: var(--label-secondary);
  text-transform: uppercase;
  letter-spacing: 0.02em;
  margin: 0 0 var(--space-2);
  padding: 0 var(--space-4);
}

.section-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 var(--space-4);
  margin-bottom: var(--space-2);
}

.section-header-row .section-header {
  margin: 0;
  padding: 0;
}

.btn-text {
  font-family: var(--font-system);
  font-size: 15px;
  color: var(--system-blue);
  background: none;
  border: none;
  cursor: pointer;
  padding: var(--space-1) var(--space-2);
}

.btn-text:disabled {
  opacity: 0.5;
}

/* ========== Stats Grid ========== */

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-3);
}

.stat-card {
  background: var(--bg-tertiary);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.stat-card.compact {
  flex-direction: column;
  text-align: center;
  padding: var(--space-4) var(--space-3);
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 22px;
  height: 22px;
}

.stat-icon.requests {
  background: rgba(0, 122, 255, 0.12);
  color: var(--system-blue);
}

.stat-icon.tokens {
  background: rgba(52, 199, 89, 0.12);
  color: var(--system-green);
}

.stat-icon.cost {
  background: rgba(255, 149, 0, 0.12);
  color: var(--system-orange);
}

.stat-content {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.stat-value {
  font-size: 22px;
  font-weight: 700;
  color: var(--label-primary);
  line-height: 1.2;
}

.stat-value-sm {
  font-size: 20px;
  font-weight: 600;
  color: var(--label-primary);
}

.stat-label {
  font-size: 13px;
  color: var(--label-secondary);
  margin-top: 2px;
}

/* ========== List Group ========== */

.list-group {
  background: var(--bg-tertiary);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  min-height: 44px;
  padding: var(--space-3) var(--space-4);
  border-bottom: 0.5px solid var(--separator);
}

.list-item:last-child {
  border-bottom: none;
}

.list-label {
  font-size: 17px;
  color: var(--label-primary);
}

.list-value {
  font-size: 17px;
  color: var(--label-secondary);
}

/* ========== Data Table ========== */

.table-container {
  background: var(--bg-tertiary);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 15px;
}

.data-table th {
  font-weight: 600;
  color: var(--label-secondary);
  padding: var(--space-3) var(--space-4);
  border-bottom: 0.5px solid var(--separator);
  background: var(--bg-secondary);
}

.data-table td {
  padding: var(--space-3) var(--space-4);
  color: var(--label-primary);
  border-bottom: 0.5px solid var(--separator);
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table .text-left {
  text-align: left;
}

.data-table .text-right {
  text-align: right;
}

.model-name {
  font-weight: 500;
  max-width: 180px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ========== Empty States ========== */

.empty-state {
  background: var(--bg-tertiary);
  border-radius: var(--radius-lg);
  padding: var(--space-8);
  text-align: center;
  color: var(--label-tertiary);
}

.empty-hero {
  max-width: 600px;
  margin: var(--space-8) auto;
  text-align: center;
}

.empty-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto var(--space-4);
  color: var(--label-tertiary);
}

.empty-text {
  font-size: 17px;
  color: var(--label-tertiary);
  margin: 0;
}

/* ========== Responsive ========== */

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .stat-card {
    flex-direction: row;
  }

  .stat-card.compact {
    flex-direction: row;
    text-align: left;
  }

  .input-group {
    flex-direction: column;
  }

  .btn-primary {
    width: 100%;
  }

  .large-title {
    font-size: 28px;
  }
}

/* ========== Reduced Motion ========== */

@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
</style>
