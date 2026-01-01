<!--
 * 文件作用：账户负载页面 - Apple HIG 风格
 * 负责功能：
 *   - 时间范围筛选
 *   - 请求数和成功率统计
 *   - Token使用和耗时统计
 *   - 负载可视化进度条
 * 重要程度：⭐⭐⭐ 一般（负载监控）
-->
<template>
  <div class="account-load-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">账户负载</h1>
        <p class="page-subtitle">查看各账户的请求负载和性能统计</p>
      </div>
      <button class="btn btn-outline" @click="fetchLoadStats" :disabled="loading">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: loading }">
          <polyline points="23,4 23,10 17,10"/>
          <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
        </svg>
        刷新
      </button>
    </div>

    <!-- 时间范围筛选 -->
    <div class="filter-card">
      <div class="filter-row">
        <div class="filter-group">
          <label class="filter-label">时间范围</label>
          <div class="time-range-selector">
            <button
              v-for="range in timeRanges"
              :key="range.value"
              :class="['range-btn', { active: timeRange === range.value }]"
              @click="timeRange = range.value; handleTimeRangeChange()"
            >
              {{ range.label }}
            </button>
          </div>
        </div>
        <div v-if="timeRange === 'custom'" class="filter-group">
          <label class="filter-label">自定义范围</label>
          <div class="date-inputs">
            <input type="datetime-local" v-model="customStart" class="date-input" @change="handleDateChange" />
            <span class="date-sep">至</span>
            <input type="datetime-local" v-model="customEnd" class="date-input" @change="handleDateChange" />
          </div>
        </div>
      </div>
    </div>

    <!-- 负载列表 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>账户名称</th>
              <th>平台</th>
              <th class="col-right">请求数</th>
              <th>成功率</th>
              <th class="col-right">Token 使用</th>
              <th class="col-right">平均耗时</th>
              <th>最后使用</th>
              <th>负载</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="row in loadStats" :key="row.account_id">
              <td><code class="id-code">#{{ row.account_id }}</code></td>
              <td><span class="account-name">{{ row.account_name }}</span></td>
              <td>
                <span :class="['platform-badge', row.platform]">{{ row.platform }}</span>
              </td>
              <td class="col-right">
                <div class="request-count">
                  <span class="total">{{ row.request_count }}</span>
                  <span class="detail">
                    (<span class="success">{{ row.success_count }}</span> /
                    <span class="error">{{ row.error_count }}</span>)
                  </span>
                </div>
              </td>
              <td>
                <div class="progress-wrap">
                  <div class="progress-bar">
                    <div
                      class="progress-fill"
                      :class="getProgressClass(row)"
                      :style="{ width: getSuccessRate(row) + '%' }"
                    ></div>
                  </div>
                  <span class="progress-text">{{ getSuccessRate(row) }}%</span>
                </div>
              </td>
              <td class="col-right">
                <span class="token-value">{{ formatTokens(row.total_tokens) }}</span>
              </td>
              <td class="col-right">
                <span class="duration-value">{{ row.avg_duration }}ms</span>
              </td>
              <td>
                <span class="time-text">{{ formatTime(row.last_used_at) }}</span>
              </td>
              <td>
                <div class="load-bar-wrap">
                  <div class="load-bar">
                    <div
                      class="load-fill"
                      :style="{ width: getLoadPercentage(row) + '%', background: getLoadColor(row) }"
                    ></div>
                  </div>
                  <span class="load-text">{{ getLoadPercentage(row) }}%</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <div v-if="!loading && loadStats.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
          <span>暂无负载数据</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'

const loading = ref(false)
const loadStats = ref([])
const timeRange = ref('24h')
const customStart = ref('')
const customEnd = ref('')
const startTime = ref('')
const endTime = ref('')

const timeRanges = [
  { label: '1小时', value: '1h' },
  { label: '24小时', value: '24h' },
  { label: '7天', value: '7d' },
  { label: '30天', value: '30d' },
  { label: '自定义', value: 'custom' }
]

const maxRequests = computed(() => {
  return Math.max(...loadStats.value.map(s => s.request_count), 1)
})

function formatTokens(tokens) {
  if (tokens >= 1000000) return (tokens / 1000000).toFixed(1) + 'M'
  if (tokens >= 1000) return (tokens / 1000).toFixed(1) + 'K'
  return tokens || 0
}

function formatTime(time) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

function getSuccessRate(row) {
  if (row.request_count === 0) return 100
  return Math.round((row.success_count / row.request_count) * 100)
}

function getProgressClass(row) {
  const rate = getSuccessRate(row)
  if (rate >= 95) return 'success'
  if (rate >= 80) return 'warning'
  return 'error'
}

function getLoadPercentage(row) {
  return Math.round((row.request_count / maxRequests.value) * 100)
}

function getLoadColor(row) {
  const percentage = getLoadPercentage(row)
  if (percentage > 80) return 'var(--apple-red)'
  if (percentage > 50) return 'var(--apple-orange)'
  return 'var(--apple-green)'
}

function handleTimeRangeChange() {
  const now = new Date()
  switch (timeRange.value) {
    case '1h':
      startTime.value = new Date(now - 60 * 60 * 1000).toISOString()
      endTime.value = now.toISOString()
      fetchLoadStats()
      break
    case '24h':
      startTime.value = new Date(now - 24 * 60 * 60 * 1000).toISOString()
      endTime.value = now.toISOString()
      fetchLoadStats()
      break
    case '7d':
      startTime.value = new Date(now - 7 * 24 * 60 * 60 * 1000).toISOString()
      endTime.value = now.toISOString()
      fetchLoadStats()
      break
    case '30d':
      startTime.value = new Date(now - 30 * 24 * 60 * 60 * 1000).toISOString()
      endTime.value = now.toISOString()
      fetchLoadStats()
      break
    case 'custom':
      break
  }
}

function handleDateChange() {
  if (customStart.value && customEnd.value) {
    startTime.value = new Date(customStart.value).toISOString()
    endTime.value = new Date(customEnd.value).toISOString()
    fetchLoadStats()
  }
}

async function fetchLoadStats() {
  loading.value = true
  try {
    const params = {}
    if (startTime.value) params.start_time = startTime.value
    if (endTime.value) params.end_time = endTime.value
    const res = await api.getAccountLoadStats(params)
    loadStats.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch load stats:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  handleTimeRangeChange()
})
</script>

<style scoped>
.account-load-page {
  max-width: 1400px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--apple-spacing-xl);
}

.header-content { flex: 1; }

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

/* 筛选卡片 */
.filter-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-xl);
}

.filter-row {
  display: flex;
  align-items: flex-end;
  gap: var(--apple-spacing-xl);
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.filter-label {
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.time-range-selector {
  display: flex;
  background: var(--apple-fill-quaternary);
  border-radius: var(--apple-radius-sm);
  padding: 2px;
}

.range-btn {
  padding: var(--apple-spacing-xs) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  border-radius: var(--apple-radius-xs);
  transition: all var(--apple-duration-fast);
}

.range-btn:hover {
  color: var(--apple-text-primary);
}

.range-btn.active {
  background: var(--apple-bg-primary);
  color: var(--apple-blue);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.date-inputs {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.date-input {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  background: var(--apple-bg-primary);
}

.date-input:focus {
  outline: none;
  border-color: var(--apple-blue);
}

.date-sep {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

/* 表格卡片 */
.table-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

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

.col-right { text-align: right; }

.id-code {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.account-name {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

/* 平台徽章 */
.platform-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
  text-transform: capitalize;
}

.platform-badge.claude { background: #f3e8ff; color: #9333ea; }
.platform-badge.openai { background: #dcfce7; color: #16a34a; }
.platform-badge.gemini { background: #dbeafe; color: #2563eb; }

/* 请求计数 */
.request-count {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.request-count .total {
  font-weight: var(--apple-font-bold);
  font-size: var(--apple-text-md);
  color: var(--apple-text-primary);
}

.request-count .detail {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.request-count .success { color: var(--apple-green); }
.request-count .error { color: var(--apple-red); }

/* 进度条 */
.progress-wrap, .load-bar-wrap {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.progress-bar, .load-bar {
  flex: 1;
  height: 8px;
  background: var(--apple-fill-tertiary);
  border-radius: 4px;
  overflow: hidden;
  min-width: 80px;
}

.progress-fill, .load-fill {
  height: 100%;
  border-radius: 4px;
  transition: width var(--apple-duration-normal);
}

.progress-fill.success { background: var(--apple-green); }
.progress-fill.warning { background: var(--apple-orange); }
.progress-fill.error { background: var(--apple-red); }

.progress-text, .load-text {
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  min-width: 36px;
  text-align: right;
}

.token-value, .duration-value {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
}

.time-text {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
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

@keyframes spin { to { transform: rotate(360deg); } }
.spinning { animation: spin 1s linear infinite; }

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: var(--apple-spacing-md);
  opacity: 0.5;
}

/* 按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
  cursor: pointer;
}

.btn svg { width: 16px; height: 16px; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }

.btn-outline {
  background: transparent;
  color: var(--apple-blue);
  border: 1px solid var(--apple-blue);
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-blue-light);
}

/* 响应式 */
@media (max-width: 1024px) {
  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }

  .time-range-selector {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }
}
</style>
