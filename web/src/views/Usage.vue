<!--
 * 文件作用：使用统计页面 - Apple HIG 风格
 * 负责功能：
 *   - 总览统计卡片（请求数、Token、费用）
 *   - 每日使用统计表格
 *   - 按模型分类统计
 *   - 日期范围筛选
 * 重要程度：⭐⭐⭐ 一般（统计展示）
 *
 * Design Decisions:
 * - Uses SF Pro system font stack for native feel
 * - 44pt minimum touch targets for accessibility
 * - Card-based layout with subtle shadows
 * - Supports prefers-color-scheme for auto theming
 * - Gradient icons for visual hierarchy
-->
<template>
  <div class="usage-page">
    <!-- 页面标题 -->
    <header class="page-header">
      <div class="header-content">
        <h1 class="page-title">使用统计</h1>
        <p class="page-subtitle">查看 API 使用情况和费用统计</p>
      </div>
      <button class="refresh-btn" @click="refreshAll" :disabled="isRefreshing">
        <svg :class="['refresh-icon', { spinning: isRefreshing }]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M23 4v6h-6M1 20v-6h6"/>
          <path d="M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15"/>
        </svg>
        <span>{{ isRefreshing ? '刷新中...' : '刷新' }}</span>
      </button>
    </header>

    <!-- 统计卡片网格 -->
    <section class="stats-section">
      <div class="stats-grid">
        <!-- 总请求数 -->
        <div class="stat-card">
          <div class="stat-icon requests">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
            </svg>
          </div>
          <div class="stat-body">
            <span class="stat-label">总请求数</span>
            <div class="stat-value-row">
              <span class="stat-value">{{ formatNumber(summary.totalRequests) }}</span>
              <span class="stat-unit">次</span>
            </div>
          </div>
        </div>

        <!-- 总 Token -->
        <div class="stat-card">
          <div class="stat-icon tokens">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <polygon points="16.24,7.76 14.12,14.12 7.76,16.24 9.88,9.88"/>
            </svg>
          </div>
          <div class="stat-body">
            <span class="stat-label">总 Token 消耗</span>
            <div class="stat-value-row">
              <span class="stat-value">{{ formatTokens(summary.totalTokens) }}</span>
            </div>
          </div>
        </div>

        <!-- 今日消费 -->
        <div class="stat-card">
          <div class="stat-icon today">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
              <line x1="16" y1="2" x2="16" y2="6"/>
              <line x1="8" y1="2" x2="8" y2="6"/>
              <line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
          </div>
          <div class="stat-body">
            <span class="stat-label">今日消费</span>
            <div class="stat-value-row">
              <span class="stat-value">${{ formatCost(summary.todayCost) }}</span>
            </div>
          </div>
        </div>

        <!-- 总消费 -->
        <div class="stat-card featured">
          <div class="stat-icon total">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="1" x2="12" y2="23"/>
              <path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
            </svg>
          </div>
          <div class="stat-body">
            <span class="stat-label">总消费</span>
            <div class="stat-value-row">
              <span class="stat-value">${{ formatCost(summary.totalCost) }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 每日使用统计 -->
    <section class="table-section">
      <div class="section-card">
        <div class="section-header">
          <div class="section-title-row">
            <svg class="section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 3v18h18"/>
              <path d="M18 9l-5 5-4-4-3 3"/>
            </svg>
            <h2>每日使用统计</h2>
          </div>
          <div class="date-controls">
            <div class="date-picker">
              <input
                type="date"
                v-model="startDate"
                class="date-input"
                @change="fetchDailyUsage"
              />
              <span class="date-separator">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="5" y1="12" x2="19" y2="12"/>
                  <polyline points="12,5 19,12 12,19"/>
                </svg>
              </span>
              <input
                type="date"
                v-model="endDate"
                class="date-input"
                @change="fetchDailyUsage"
              />
            </div>
            <div class="date-shortcuts">
              <button
                v-for="shortcut in dateShortcuts"
                :key="shortcut.text"
                :class="['shortcut-btn', { active: activeShortcut === shortcut.text }]"
                @click="applyShortcut(shortcut)"
              >
                {{ shortcut.text }}
              </button>
            </div>
          </div>
        </div>

        <div class="table-wrapper">
          <table class="data-table" v-if="!dailyLoading && dailyUsage.length > 0">
            <thead>
              <tr>
                <th>日期</th>
                <th class="text-right">请求数</th>
                <th class="text-right">总 Token</th>
                <th class="text-right">费用</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in dailyUsage" :key="row.date" class="table-row">
                <td>
                  <span class="date-cell">{{ row.date }}</span>
                </td>
                <td class="text-right">
                  <span class="number-cell">{{ formatNumber(row.request_count) }}</span>
                </td>
                <td class="text-right">
                  <span class="number-cell">{{ formatNumber(row.total_tokens) }}</span>
                </td>
                <td class="text-right">
                  <span class="cost-cell">${{ formatCost(row.total_cost) }}</span>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- 加载状态 -->
          <div v-if="dailyLoading" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>

          <!-- 空状态 -->
          <div v-if="!dailyLoading && dailyUsage.length === 0" class="empty-state">
            <div class="empty-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M3 3v18h18"/>
                <path d="M18 9l-5 5-4-4-3 3"/>
              </svg>
            </div>
            <span class="empty-text">暂无使用记录</span>
            <span class="empty-hint">选择日期范围查看统计数据</span>
          </div>
        </div>
      </div>
    </section>

    <!-- 按模型统计 -->
    <section class="table-section">
      <div class="section-card">
        <div class="section-header">
          <div class="section-title-row">
            <svg class="section-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <polygon points="16.24,7.76 14.12,14.12 7.76,16.24 9.88,9.88"/>
            </svg>
            <h2>按模型统计</h2>
          </div>
        </div>

        <div class="table-wrapper">
          <table class="data-table" v-if="!modelLoading && modelUsage.length > 0">
            <thead>
              <tr>
                <th>模型</th>
                <th class="text-right">请求数</th>
                <th class="text-right">总 Token</th>
                <th class="text-right">费用</th>
                <th>占比</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in modelUsage" :key="row.model" class="table-row">
                <td>
                  <span class="model-badge">{{ row.model }}</span>
                </td>
                <td class="text-right">
                  <span class="number-cell">{{ formatNumber(row.request_count) }}</span>
                </td>
                <td class="text-right">
                  <span class="number-cell">{{ formatNumber(row.total_tokens) }}</span>
                </td>
                <td class="text-right">
                  <span class="cost-cell">${{ formatCost(row.total_cost) }}</span>
                </td>
                <td>
                  <div class="progress-cell">
                    <div class="progress-bar">
                      <div
                        class="progress-fill"
                        :style="{ width: getModelPercentage(row.total_cost) + '%' }"
                      ></div>
                    </div>
                    <span class="progress-label">{{ getModelPercentage(row.total_cost) }}%</span>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- 加载状态 -->
          <div v-if="modelLoading" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>

          <!-- 空状态 -->
          <div v-if="!modelLoading && modelUsage.length === 0" class="empty-state">
            <div class="empty-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <polygon points="16.24,7.76 14.12,14.12 7.76,16.24 9.88,9.88"/>
              </svg>
            </div>
            <span class="empty-text">暂无模型使用记录</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import api from '@/api'

const summary = reactive({
  totalRequests: 0,
  totalTokens: 0,
  todayCost: 0,
  totalCost: 0
})

const dailyUsage = ref([])
const dailyLoading = ref(false)
const modelUsage = ref([])
const modelLoading = ref(false)
const isRefreshing = ref(false)
const activeShortcut = ref('最近一周')

// 日期
const today = new Date()
const weekAgo = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000)

const startDate = ref(formatDateForInput(weekAgo))
const endDate = ref(formatDateForInput(today))

const dateShortcuts = [
  { text: '最近一周', days: 7 },
  { text: '最近一月', days: 30 },
  { text: '最近三月', days: 90 }
]

function formatDateForInput(date) {
  return date.toISOString().split('T')[0]
}

function applyShortcut(shortcut) {
  activeShortcut.value = shortcut.text
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - shortcut.days * 24 * 60 * 60 * 1000)
  startDate.value = formatDateForInput(start)
  endDate.value = formatDateForInput(end)
  fetchDailyUsage()
  fetchModelUsage()
}

const totalModelCost = computed(() => {
  return modelUsage.value.reduce((sum, m) => sum + (m.total_cost || 0), 0)
})

function getModelPercentage(cost) {
  if (!totalModelCost.value || !cost) return 0
  return Math.round((cost / totalModelCost.value) * 100)
}

function formatNumber(num) {
  if (!num) return '0'
  return num.toLocaleString()
}

function formatTokens(tokens) {
  if (!tokens) return '0'
  if (tokens >= 1000000) return (tokens / 1000000).toFixed(2) + 'M'
  if (tokens >= 1000) return (tokens / 1000).toFixed(1) + 'K'
  return tokens.toLocaleString()
}

function formatCost(cost) {
  if (cost === undefined || cost === null) return '0.0000'
  return cost.toFixed(4)
}

async function fetchSummary() {
  try {
    const res = await api.adminGetUsageSummary()
    if (res.data) {
      summary.totalRequests = res.data.total_requests || 0
      summary.totalTokens = res.data.total_tokens || 0
      summary.todayCost = res.data.today_cost || 0
      summary.totalCost = res.data.total_cost || 0
    }
  } catch (e) {
    console.error('Failed to fetch summary:', e)
  }
}

async function fetchDailyUsage() {
  if (!startDate.value || !endDate.value) return

  dailyLoading.value = true
  try {
    const res = await api.adminGetDailyUsage({
      start_date: startDate.value,
      end_date: endDate.value
    })
    dailyUsage.value = res.data?.daily || []
  } catch (e) {
    console.error('Failed to fetch daily usage:', e)
  } finally {
    dailyLoading.value = false
  }
}

async function fetchModelUsage() {
  modelLoading.value = true
  try {
    const res = await api.adminGetModelUsage({
      start_date: startDate.value,
      end_date: endDate.value
    })
    modelUsage.value = res.data?.models || []
  } catch (e) {
    console.error('Failed to fetch model usage:', e)
  } finally {
    modelLoading.value = false
  }
}

async function refreshAll() {
  isRefreshing.value = true
  await Promise.all([
    fetchSummary(),
    fetchDailyUsage(),
    fetchModelUsage()
  ])
  isRefreshing.value = false
}

onMounted(() => {
  fetchSummary()
  fetchDailyUsage()
  fetchModelUsage()
})
</script>

<style scoped>
.usage-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--apple-spacing-md);
}

/* ========== 页面标题 ========== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--apple-spacing-xl);
  gap: var(--apple-spacing-md);
}

.page-title {
  font-size: var(--apple-text-3xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-xs) 0;
  letter-spacing: var(--apple-tracking-tight);
}

.page-subtitle {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  margin: 0;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-xs) var(--apple-spacing-md);
  background: var(--apple-bg-primary);
  color: var(--apple-blue);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-full);
  box-shadow: var(--apple-shadow-card);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  white-space: nowrap;
}

.refresh-btn:hover:not(:disabled) {
  box-shadow: var(--apple-shadow-card-hover);
  transform: translateY(-1px);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-icon {
  width: 16px;
  height: 16px;
}

.refresh-icon.spinning {
  animation: spin 1s linear infinite;
}

/* ========== 统计卡片 ========== */
.stats-section {
  margin-bottom: var(--apple-spacing-xl);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--apple-spacing-md);
}

.stat-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  padding: var(--apple-spacing-lg);
  box-shadow: var(--apple-shadow-card);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  transition: all var(--apple-duration-normal) var(--apple-ease-default);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--apple-shadow-card-hover);
}

.stat-card.featured {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-card.featured .stat-label,
.stat-card.featured .stat-value,
.stat-card.featured .stat-unit {
  color: white;
}

.stat-icon {
  width: 52px;
  height: 52px;
  border-radius: var(--apple-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 26px;
  height: 26px;
  color: white;
}

.stat-icon.requests {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.tokens {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.today {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.stat-icon.total {
  background: rgba(255, 255, 255, 0.2);
}

.stat-body {
  flex: 1;
  min-width: 0;
}

.stat-label {
  display: block;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-bottom: 4px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value-row {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.stat-value {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  line-height: 1.2;
}

.stat-unit {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

/* ========== 表格区块 ========== */
.table-section {
  margin-bottom: var(--apple-spacing-xl);
}

.section-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
  flex-wrap: wrap;
  gap: var(--apple-spacing-md);
}

.section-title-row {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.section-icon {
  width: 24px;
  height: 24px;
  color: var(--apple-blue);
}

.section-header h2 {
  margin: 0;
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

/* ========== 日期选择器 ========== */
.date-controls {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  flex-wrap: wrap;
}

.date-picker {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  background: var(--apple-bg-secondary);
  padding: var(--apple-spacing-xs);
  border-radius: var(--apple-radius-md);
}

.date-input {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  font-family: var(--apple-font-family);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.date-input:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.date-separator {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  color: var(--apple-text-tertiary);
}

.date-separator svg {
  width: 16px;
  height: 16px;
}

.date-shortcuts {
  display: flex;
  gap: var(--apple-spacing-xs);
}

.shortcut-btn {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  color: var(--apple-blue);
  background: var(--apple-blue-light);
  border-radius: var(--apple-radius-full);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.shortcut-btn:hover {
  background: var(--apple-blue);
  color: white;
}

.shortcut-btn.active {
  background: var(--apple-blue);
  color: white;
}

/* ========== 数据表格 ========== */
.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--apple-text-sm);
}

.data-table th,
.data-table td {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  text-align: left;
}

.data-table th {
  background: var(--apple-bg-secondary);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  font-size: var(--apple-text-xs);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
  border-bottom: 1px solid var(--apple-separator);
}

.data-table .text-right {
  text-align: right;
}

.table-row {
  transition: background var(--apple-duration-fast) var(--apple-ease-default);
  border-bottom: 1px solid var(--apple-separator);
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background: var(--apple-bg-secondary);
}

.date-cell {
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.number-cell {
  font-family: var(--apple-font-mono);
  color: var(--apple-text-primary);
}

.cost-cell {
  font-family: var(--apple-font-mono);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-green);
}

.model-badge {
  display: inline-block;
  padding: var(--apple-spacing-xxs) var(--apple-spacing-sm);
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-blue);
  background: var(--apple-blue-light);
  border-radius: var(--apple-radius-sm);
}

/* ========== 进度条 ========== */
.progress-cell {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  min-width: 140px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: var(--apple-fill-tertiary);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
  transition: width var(--apple-duration-normal) var(--apple-ease-default);
}

.progress-label {
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  min-width: 36px;
  text-align: right;
}

/* ========== 加载和空状态 ========== */
.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-xxxl);
  color: var(--apple-text-tertiary);
}

.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid var(--apple-fill-tertiary);
  border-top-color: var(--apple-blue);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: var(--apple-spacing-md);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin-bottom: var(--apple-spacing-md);
  opacity: 0.4;
}

.empty-icon svg {
  width: 100%;
  height: 100%;
}

.empty-text {
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  margin-bottom: var(--apple-spacing-xs);
}

.empty-hint {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

/* ========== 响应式 - 移动端适配 ========== */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 767px) {
  .usage-page {
    padding: 0 var(--apple-spacing-md);
  }

  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: var(--apple-spacing-md);
  }

  .page-title {
    font-size: var(--apple-text-2xl);
  }

  .page-subtitle {
    font-size: var(--apple-text-sm);
  }

  .refresh-btn {
    width: 100%;
    justify-content: center;
    min-height: var(--apple-touch-target);
  }

  /* 统计卡片 */
  .stats-grid {
    grid-template-columns: 1fr;
    gap: var(--apple-spacing-sm);
  }

  .stat-card {
    padding: var(--apple-spacing-md);
  }

  .stat-icon {
    width: 44px;
    height: 44px;
  }

  .stat-icon svg {
    width: 22px;
    height: 22px;
  }

  .stat-value {
    font-size: var(--apple-text-xl);
  }

  /* 表格区块 */
  .table-section {
    margin-bottom: var(--apple-spacing-lg);
  }

  .section-card {
    margin: 0 calc(var(--apple-spacing-md) * -1);
    border-radius: 0;
  }

  .section-header {
    padding: var(--apple-spacing-md);
    flex-direction: column;
    align-items: stretch;
  }

  .section-title-row h2 {
    font-size: var(--apple-text-md);
  }

  /* 日期选择器 */
  .date-controls {
    width: 100%;
    flex-direction: column;
    align-items: stretch;
  }

  .date-picker {
    width: 100%;
    flex-direction: column;
    gap: var(--apple-spacing-xs);
  }

  .date-input {
    flex: 1;
    min-height: var(--apple-touch-target);
    font-size: 16px; /* 防止 iOS 缩放 */
  }

  .date-separator {
    display: none;
  }

  .date-shortcuts {
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    gap: var(--apple-spacing-xs);
  }

  .shortcut-btn {
    flex: 1;
    text-align: center;
    min-height: var(--apple-touch-target-sm);
  }

  /* 表格 */
  .data-table th,
  .data-table td {
    padding: var(--apple-spacing-sm) var(--apple-spacing-xs);
    font-size: var(--apple-text-xs);
  }

  .model-badge {
    font-size: 10px;
    padding: 2px 6px;
  }

  .progress-cell {
    min-width: 80px;
  }

  .progress-label {
    min-width: 28px;
    font-size: 10px;
  }

  /* 加载和空状态 */
  .loading-state,
  .empty-state {
    padding: var(--apple-spacing-xl);
  }

  .empty-icon {
    width: 48px;
    height: 48px;
  }
}

@media (max-width: 375px) {
  .page-title {
    font-size: var(--apple-text-xl);
  }

  .stat-value {
    font-size: var(--apple-text-lg);
  }

  .section-title-row h2 {
    font-size: var(--apple-text-base);
  }

  .data-table th,
  .data-table td {
    font-size: 10px;
    padding: var(--apple-spacing-xs);
  }
}

/* ========== 减少动画偏好 ========== */
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
