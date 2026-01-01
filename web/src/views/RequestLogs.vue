<!--
 * 文件作用：请求日志页面 - Apple HIG 风格
 * 负责功能：
 *   - 每日请求汇总
 *   - 按模型统计
 *   - 详细记录查询
 *   - Token和费用统计
 * 重要程度：⭐⭐⭐ 一般（日志查看）
-->
<template>
  <div class="request-logs-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">请求日志</h1>
        <p class="page-subtitle">查看 API 请求统计和详细记录</p>
      </div>
      <button class="btn btn-outline" @click="refreshAll" :disabled="loadingDaily">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: loadingDaily }">
          <polyline points="23,4 23,10 17,10"/>
          <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
        </svg>
        刷新
      </button>
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
          <span class="stat-value">{{ summary.total_requests || 0 }}</span>
          <span class="stat-label">总请求数</span>
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
          <span class="stat-value">{{ formatTokens(summary.total_tokens) }}</span>
          <span class="stat-label">总 Token</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon cost">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="1" x2="12" y2="23"/>
            <path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">${{ (summary.total_cost || 0).toFixed(4) }}</span>
          <span class="stat-label">总费用</span>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon models">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"/>
            <rect x="14" y="3" width="7" height="7"/>
            <rect x="14" y="14" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-value">{{ modelStats.length }}</span>
          <span class="stat-label">模型数</span>
        </div>
      </div>
    </div>

    <!-- Tab 切换 -->
    <div class="tabs-container">
      <div class="tabs-header">
        <button :class="['tab-btn', { active: activeTab === 'daily' }]" @click="activeTab = 'daily'">每日汇总</button>
        <button :class="['tab-btn', { active: activeTab === 'models' }]" @click="activeTab = 'models'">模型统计</button>
        <button :class="['tab-btn', { active: activeTab === 'records' }]" @click="activeTab = 'records'">详细记录</button>
      </div>

      <!-- 每日汇总 -->
      <div v-show="activeTab === 'daily'" class="tab-content">
        <div class="table-card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>日期</th>
                  <th class="col-right">请求数</th>
                  <th class="col-right">Token</th>
                  <th class="col-right">费用</th>
                </tr>
              </thead>
              <tbody v-if="!loadingDaily">
                <tr v-for="row in dailyStats" :key="row.date">
                  <td>{{ row.date }}</td>
                  <td class="col-right">{{ row.request_count }}</td>
                  <td class="col-right">{{ formatTokens(row.total_tokens) }}</td>
                  <td class="col-right cost-cell">${{ row.total_cost?.toFixed(4) || '0' }}</td>
                </tr>
              </tbody>
            </table>
            <div v-if="loadingDaily" class="loading-state">
              <div class="loading-spinner"></div>
              <span>加载中...</span>
            </div>
            <div v-if="!loadingDaily && dailyStats.length === 0" class="empty-state">
              <span>暂无数据</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 模型统计 -->
      <div v-show="activeTab === 'models'" class="tab-content">
        <div class="table-card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>模型</th>
                  <th class="col-right">请求数</th>
                  <th class="col-right">Token</th>
                  <th class="col-right">费用</th>
                </tr>
              </thead>
              <tbody v-if="!loadingModels">
                <tr v-for="row in modelStats" :key="row.model">
                  <td><code class="model-code">{{ row.model }}</code></td>
                  <td class="col-right">{{ row.request_count }}</td>
                  <td class="col-right">{{ formatTokens(row.total_tokens) }}</td>
                  <td class="col-right cost-cell">${{ row.total_cost?.toFixed(4) || '0' }}</td>
                </tr>
              </tbody>
            </table>
            <div v-if="loadingModels" class="loading-state">
              <div class="loading-spinner"></div>
              <span>加载中...</span>
            </div>
            <div v-if="!loadingModels && modelStats.length === 0" class="empty-state">
              <span>暂无数据</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 详细记录 -->
      <div v-show="activeTab === 'records'" class="tab-content">
        <!-- 筛选栏 -->
        <div class="filter-bar">
          <div class="filter-group">
            <label class="filter-label">时间范围</label>
            <div class="date-inputs">
              <input type="date" v-model="dateStart" class="date-input" />
              <span class="date-sep">至</span>
              <input type="date" v-model="dateEnd" class="date-input" />
            </div>
          </div>
          <div class="filter-group">
            <label class="filter-label">模型</label>
            <select v-model="filterModel" class="filter-select">
              <option value="">全部模型</option>
              <option v-for="m in availableModels" :key="m" :value="m">{{ m }}</option>
            </select>
          </div>
          <div class="filter-actions">
            <button class="btn btn-primary btn-sm" @click="fetchRecords">查询</button>
            <button class="btn btn-secondary btn-sm" @click="resetFilters">重置</button>
          </div>
        </div>

        <div class="table-card">
          <div class="table-container">
            <table class="data-table">
              <thead>
                <tr>
                  <th>模型</th>
                  <th>请求IP</th>
                  <th class="col-right">输入</th>
                  <th class="col-right">输出</th>
                  <th class="col-right">缓存创建</th>
                  <th class="col-right">缓存读取</th>
                  <th class="col-right">总Token</th>
                  <th class="col-right">费用</th>
                  <th>时间</th>
                </tr>
              </thead>
              <tbody v-if="!loadingRecords">
                <tr v-for="row in records" :key="row.id">
                  <td><code class="model-code">{{ row.model }}</code></td>
                  <td><span class="ip-text">{{ row.request_ip }}</span></td>
                  <td class="col-right">{{ formatTokens(row.input_tokens) }}</td>
                  <td class="col-right">{{ formatTokens(row.output_tokens) }}</td>
                  <td class="col-right">
                    <span :class="{ 'cache-highlight': row.cache_creation_input_tokens > 0 }">
                      {{ formatTokens(row.cache_creation_input_tokens) }}
                    </span>
                  </td>
                  <td class="col-right">
                    <span :class="{ 'cache-read-highlight': row.cache_read_input_tokens > 0 }">
                      {{ formatTokens(row.cache_read_input_tokens) }}
                    </span>
                  </td>
                  <td class="col-right">{{ formatTokens(row.total_tokens) }}</td>
                  <td class="col-right cost-cell">${{ (row.total_cost || 0).toFixed(4) }}</td>
                  <td><span class="time-text">{{ formatTime(row.timestamp || row.request_time) }}</span></td>
                </tr>
              </tbody>
            </table>
            <div v-if="loadingRecords" class="loading-state">
              <div class="loading-spinner"></div>
              <span>加载中...</span>
            </div>
            <div v-if="!loadingRecords && records.length === 0" class="empty-state">
              <span>暂无记录</span>
            </div>
          </div>

          <div v-if="pagination.total > 0" class="table-footer">
            <div class="pagination-info">共 {{ pagination.total }} 条</div>
            <div class="pagination-controls">
              <select v-model="pagination.pageSize" @change="fetchRecords" class="page-size-select">
                <option :value="20">20 条/页</option>
                <option :value="50">50 条/页</option>
                <option :value="100">100 条/页</option>
              </select>
              <div class="page-btns">
                <button class="page-btn" :disabled="pagination.page <= 1" @click="pagination.page--; fetchRecords()">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15,18 9,12 15,6"/></svg>
                </button>
                <span class="page-current">{{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) || 1 }}</span>
                <button class="page-btn" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)" @click="pagination.page++; fetchRecords()">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9,18 15,12 9,6"/></svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import api from '@/api'

const activeTab = ref('daily')

const dateStart = ref('')
const dateEnd = ref('')
const filterModel = ref('')

const summary = reactive({
  total_requests: 0,
  total_tokens: 0,
  total_cost: 0
})

const dailyStats = ref([])
const loadingDaily = ref(false)
const modelStats = ref([])
const loadingModels = ref(false)
const records = ref([])
const loadingRecords = ref(false)
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

const availableModels = computed(() => {
  return modelStats.value.map(m => m.model).filter(Boolean)
})

function formatTokens(tokens) {
  if (!tokens) return '0'
  if (tokens >= 1000000) return (tokens / 1000000).toFixed(1) + 'M'
  if (tokens >= 1000) return (tokens / 1000).toFixed(1) + 'K'
  return tokens
}

function formatTime(time) {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

async function fetchAllSummary() {
  loadingDaily.value = true
  loadingModels.value = true
  try {
    const summaryRes = await api.adminGetUsageSummary()
    const dailyRes = await api.adminGetDailyUsage({})
    const modelsRes = await api.adminGetModelUsage({})

    summary.total_requests = summaryRes.data?.total_requests || 0
    summary.total_tokens = summaryRes.data?.total_tokens || 0
    summary.total_cost = summaryRes.data?.total_cost || 0
    dailyStats.value = dailyRes.data?.daily || []
    modelStats.value = modelsRes.data?.models || []
  } catch (e) {
    console.error('Failed to fetch summary:', e)
  } finally {
    loadingDaily.value = false
    loadingModels.value = false
  }
}

async function fetchRecords() {
  loadingRecords.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (dateStart.value && dateEnd.value) {
      params.start_date = dateStart.value
      params.end_date = dateEnd.value
    }
    if (filterModel.value) {
      params.model = filterModel.value
    }
    const res = await api.adminGetUsageRecords(params)
    const data = res.data || {}
    records.value = data.items || []
    pagination.total = data.total || records.value.length
  } catch (e) {
    console.error('Failed to fetch records:', e)
  } finally {
    loadingRecords.value = false
  }
}

function resetFilters() {
  dateStart.value = ''
  dateEnd.value = ''
  filterModel.value = ''
  pagination.page = 1
  fetchRecords()
}

function refreshAll() {
  fetchAllSummary()
  fetchRecords()
}

onMounted(() => {
  fetchAllSummary()
})
</script>

<style scoped>
.request-logs-page {
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

/* 统计卡片 */
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

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg { width: 24px; height: 24px; color: white; }
.stat-icon.requests { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.stat-icon.tokens { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.stat-icon.cost { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.stat-icon.models { background: linear-gradient(135deg, #fc4a1a 0%, #f7b733 100%); }

.stat-content { flex: 1; }
.stat-value {
  display: block;
  font-size: var(--apple-text-xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
}
.stat-label {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

/* Tabs */
.tabs-container {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.tabs-header {
  display: flex;
  border-bottom: 1px solid var(--apple-separator);
  padding: 0 var(--apple-spacing-lg);
}

.tab-btn {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: all var(--apple-duration-fast);
}

.tab-btn:hover { color: var(--apple-text-primary); }
.tab-btn.active {
  color: var(--apple-blue);
  border-bottom-color: var(--apple-blue);
}

.tab-content { padding: var(--apple-spacing-lg); }

/* 筛选栏 */
.filter-bar {
  display: flex;
  align-items: flex-end;
  gap: var(--apple-spacing-md);
  flex-wrap: wrap;
  margin-bottom: var(--apple-spacing-lg);
  padding: var(--apple-spacing-md);
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
}

.filter-group { display: flex; flex-direction: column; gap: var(--apple-spacing-xxs); }
.filter-label {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}
.filter-select, .date-input {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  background: var(--apple-bg-primary);
}
.filter-select:focus, .date-input:focus {
  outline: none;
  border-color: var(--apple-blue);
}

.date-inputs { display: flex; align-items: center; gap: var(--apple-spacing-xs); }
.date-sep { font-size: var(--apple-text-sm); color: var(--apple-text-tertiary); }
.filter-actions { display: flex; gap: var(--apple-spacing-xs); align-items: flex-end; }

/* 信息横幅 */
.info-banner {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  background: var(--apple-blue-light);
  border-radius: var(--apple-radius-md);
  margin-bottom: var(--apple-spacing-lg);
}
.info-banner svg { width: 20px; height: 20px; color: var(--apple-blue); flex-shrink: 0; }
.info-banner span { font-size: var(--apple-text-sm); color: var(--apple-blue); }

/* 表格卡片 */
.table-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-md);
  border: 1px solid var(--apple-separator);
  overflow: hidden;
}

.table-container { overflow-x: auto; }

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--apple-text-sm);
}

.data-table th, .data-table td {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  text-align: left;
  border-bottom: 1px solid var(--apple-separator);
}

.data-table th {
  background: var(--apple-bg-secondary);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  white-space: nowrap;
}

.data-table tbody tr:hover { background: var(--apple-bg-secondary); }
.col-right { text-align: right; }
.cost-cell { font-weight: var(--apple-font-semibold); color: var(--apple-green); }

.model-code {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  background: var(--apple-fill-quaternary);
  padding: 2px 6px;
  border-radius: var(--apple-radius-xs);
}

.ip-text, .time-text {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.cache-highlight { color: var(--apple-orange); font-weight: var(--apple-font-semibold); }
.cache-read-highlight { color: var(--apple-green); font-weight: var(--apple-font-semibold); }

/* 加载和空状态 */
.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-xxl);
  color: var(--apple-text-tertiary);
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--apple-fill-tertiary);
  border-top-color: var(--apple-blue);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: var(--apple-spacing-sm);
}

@keyframes spin { to { transform: rotate(360deg); } }
.spinning { animation: spin 1s linear infinite; }

/* 分页 */
.table-footer {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-info { font-size: var(--apple-text-sm); color: var(--apple-text-secondary); }
.pagination-controls { display: flex; align-items: center; gap: var(--apple-spacing-sm); }
.page-size-select {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
}
.page-btns { display: flex; align-items: center; gap: var(--apple-spacing-xs); }
.page-btn {
  width: 28px;
  height: 28px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-secondary);
  transition: all var(--apple-duration-fast);
}
.page-btn svg { width: 14px; height: 14px; }
.page-btn:hover:not(:disabled) { background: var(--apple-blue); color: white; }
.page-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.page-current {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  min-width: 60px;
  text-align: center;
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
.btn-sm { padding: var(--apple-spacing-xs) var(--apple-spacing-md); }
.btn-primary { background: var(--apple-blue); color: white; }
.btn-primary:hover:not(:disabled) { background: var(--apple-blue-hover); }
.btn-secondary { background: var(--apple-fill-tertiary); color: var(--apple-text-primary); }
.btn-secondary:hover:not(:disabled) { background: var(--apple-fill-secondary); }
.btn-outline {
  background: transparent;
  color: var(--apple-blue);
  border: 1px solid var(--apple-blue);
}
.btn-outline:hover:not(:disabled) { background: var(--apple-blue-light); }

/* 响应式 */
@media (max-width: 1024px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
}
@media (max-width: 768px) {
  .stats-grid { grid-template-columns: 1fr; }
  .filter-bar { flex-direction: column; align-items: stretch; }
  .filter-group { width: 100%; }
  .filter-select, .date-inputs { width: 100%; }
}
</style>
