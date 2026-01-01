<!--
 * 文件作用：使用统计页面 - Apple HIG 风格
 * 负责功能：
 *   - 总览统计卡片（请求数、Token、费用）
 *   - 每日使用统计表格
 *   - 按模型分类统计
 *   - 日期范围筛选
 * 重要程度：⭐⭐⭐ 一般（统计展示）
-->
<template>
  <div class="usage-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">使用统计</h1>
        <p class="page-subtitle">查看 API 使用情况和费用统计</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon requests">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">总请求数</span>
          <span class="stat-value">{{ formatNumber(summary.totalRequests) }}</span>
          <span class="stat-unit">次</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon tokens">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polygon points="16.24,7.76 14.12,14.12 7.76,16.24 9.88,9.88"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">总 Token 消耗</span>
          <span class="stat-value">{{ formatNumber(summary.totalTokens) }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon today">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">今日消费</span>
          <span class="stat-value">${{ summary.todayCost.toFixed(4) }}</span>
        </div>
      </div>

      <div class="stat-card highlight">
        <div class="stat-icon total">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="1" x2="12" y2="23"/>
            <path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">总消费</span>
          <span class="stat-value">${{ summary.totalCost.toFixed(4) }}</span>
        </div>
      </div>
    </div>

    <!-- 每日使用统计 -->
    <div class="table-card">
      <div class="card-header">
        <h3>每日使用统计</h3>
        <div class="date-picker">
          <div class="date-input-group">
            <input
              type="date"
              v-model="startDate"
              class="date-input"
              @change="fetchDailyUsage"
            />
            <span class="date-separator">至</span>
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
              class="shortcut-btn"
              @click="applyShortcut(shortcut)"
            >
              {{ shortcut.text }}
            </button>
          </div>
        </div>
      </div>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>日期</th>
              <th class="col-right">请求数</th>
              <th class="col-right">总 Token</th>
              <th class="col-right">费用 ($)</th>
            </tr>
          </thead>
          <tbody v-if="!dailyLoading">
            <tr v-for="row in dailyUsage" :key="row.date">
              <td>{{ row.date }}</td>
              <td class="col-right">{{ formatNumber(row.request_count) }}</td>
              <td class="col-right">{{ formatNumber(row.total_tokens) }}</td>
              <td class="col-right cost-cell">${{ formatCost(row.total_cost) }}</td>
            </tr>
          </tbody>
        </table>

        <div v-if="dailyLoading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <div v-if="!dailyLoading && dailyUsage.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M3 3v18h18"/>
            <path d="M18 9l-5 5-4-4-3 3"/>
          </svg>
          <span>暂无使用记录</span>
        </div>
      </div>
    </div>

    <!-- 按模型统计 -->
    <div class="table-card">
      <div class="card-header">
        <h3>按模型统计</h3>
      </div>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>模型</th>
              <th class="col-right">请求数</th>
              <th class="col-right">总 Token</th>
              <th class="col-right">费用 ($)</th>
              <th>占比</th>
            </tr>
          </thead>
          <tbody v-if="!modelLoading">
            <tr v-for="row in modelUsage" :key="row.model">
              <td>
                <span class="model-name">{{ row.model }}</span>
              </td>
              <td class="col-right">{{ formatNumber(row.request_count) }}</td>
              <td class="col-right">{{ formatNumber(row.total_tokens) }}</td>
              <td class="col-right cost-cell">${{ formatCost(row.total_cost) }}</td>
              <td>
                <div class="progress-wrap">
                  <div class="progress-bar">
                    <div class="progress-fill" :style="{ width: getModelPercentage(row.total_cost) + '%' }"></div>
                  </div>
                  <span class="progress-text">{{ getModelPercentage(row.total_cost) }}%</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="modelLoading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <div v-if="!modelLoading && modelUsage.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <polygon points="16.24,7.76 14.12,14.12 7.76,16.24 9.88,9.88"/>
          </svg>
          <span>暂无模型使用记录</span>
        </div>
      </div>
    </div>
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
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - shortcut.days * 24 * 60 * 60 * 1000)
  startDate.value = formatDateForInput(start)
  endDate.value = formatDateForInput(end)
  fetchDailyUsage()
}

const totalModelCost = computed(() => {
  return modelUsage.value.reduce((sum, m) => sum + (m.cost || 0), 0)
})

function getModelPercentage(cost) {
  if (!totalModelCost.value || !cost) return 0
  return Math.round((cost / totalModelCost.value) * 100)
}

function formatNumber(num) {
  if (!num) return '0'
  return num.toLocaleString()
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
}

/* 页面标题 */
.page-header {
  margin-bottom: var(--apple-spacing-xl);
}

.page-title {
  font-size: var(--apple-text-3xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-xs) 0;
}

.page-subtitle {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  margin: 0;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--apple-spacing-md);
  margin-bottom: var(--apple-spacing-xl);
}

.stat-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-lg);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
}

.stat-card.highlight {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-card.highlight .stat-label,
.stat-card.highlight .stat-value,
.stat-card.highlight .stat-unit {
  color: white;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 24px;
  height: 24px;
  color: white;
}

.stat-icon.requests { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.stat-icon.tokens { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.stat-icon.today { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.stat-icon.total { background: rgba(255, 255, 255, 0.2); }

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  display: block;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-bottom: 4px;
}

.stat-value {
  font-size: var(--apple-text-xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
}

.stat-unit {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-left: 4px;
}

/* 表格卡片 */
.table-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  margin-bottom: var(--apple-spacing-xl);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
  flex-wrap: wrap;
  gap: var(--apple-spacing-md);
}

.card-header h3 {
  margin: 0;
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

/* 日期选择器 */
.date-picker {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  flex-wrap: wrap;
}

.date-input-group {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.date-input {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
}

.date-input:focus {
  outline: none;
  border-color: var(--apple-blue);
}

.date-separator {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

.date-shortcuts {
  display: flex;
  gap: var(--apple-spacing-xs);
}

.shortcut-btn {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
  color: var(--apple-blue);
  background: var(--apple-blue-light);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
}

.shortcut-btn:hover {
  background: var(--apple-blue);
  color: white;
}

/* 表格 */
.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--apple-text-sm);
}

.data-table th,
.data-table td {
  padding: var(--apple-spacing-md);
  text-align: left;
  border-bottom: 1px solid var(--apple-separator);
}

.data-table th {
  background: var(--apple-bg-secondary);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  white-space: nowrap;
}

.data-table tbody tr:hover {
  background: var(--apple-bg-secondary);
}

.col-right {
  text-align: right;
}

.cost-cell {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-green);
}

.model-name {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-primary);
}

/* 进度条 */
.progress-wrap {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: var(--apple-fill-tertiary);
  border-radius: 4px;
  overflow: hidden;
  min-width: 80px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
  transition: width var(--apple-duration-normal);
}

.progress-text {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
  min-width: 36px;
}

/* 加载和空状态 */
.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-xxxl);
  color: var(--apple-text-tertiary);
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--apple-fill-tertiary);
  border-top-color: var(--apple-blue);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: var(--apple-spacing-md);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: var(--apple-spacing-md);
  opacity: 0.5;
}

/* 响应式 */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .date-picker {
    width: 100%;
    flex-direction: column;
    align-items: flex-start;
  }

  .date-input-group {
    width: 100%;
  }

  .date-input {
    flex: 1;
  }
}
</style>
