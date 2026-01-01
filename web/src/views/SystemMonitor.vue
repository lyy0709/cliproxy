<!--
 * 文件作用：系统监控页面 - Apple HIG 风格
 * 负责功能：
 *   - 今日使用概览（消费、Token、请求、用户）
 *   - Token使用详情
 *   - 账号状态统计
 *   - 系统资源监控（CPU、内存、磁盘）
 *   - 缓存和MySQL状态
 * 重要程度：⭐⭐⭐⭐ 重要（运维监控）
 * 依赖模块：api, toast
-->
<template>
  <div class="system-monitor-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">系统监控</h1>
        <p class="page-subtitle">实时查看系统运行状态和资源使用情况</p>
      </div>
      <button class="btn btn-primary" @click="fetchData" :disabled="loading">
        <svg v-if="loading" class="btn-spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
          <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23,4 23,10 17,10"/>
          <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
        </svg>
        刷新数据
      </button>
    </div>

    <!-- 今日使用概览 -->
    <section class="section">
      <h2 class="section-title">今日使用</h2>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon cost">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="1" x2="12" y2="23"/>
              <path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">${{ formatNumber(data.today_usage?.total_cost || 0, 4) }}</div>
            <div class="stat-label">今日消费</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon token">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M12 6v6l4 2"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatLargeNumber(data.today_usage?.total_tokens || 0) }}</div>
            <div class="stat-label">今日 Token</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon request">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="22,12 18,12 15,21 9,3 6,12 2,12"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatNumber(data.today_usage?.request_count || 0) }}</div>
            <div class="stat-label">今日请求</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon user">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 00-3-3.87"/>
              <path d="M16 3.13a4 4 0 010 7.75"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ data.users?.active || 0 }}</div>
            <div class="stat-label">活跃用户</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 累计使用统计 -->
    <section class="section">
      <h2 class="section-title">累计统计</h2>
      <div class="stats-grid">
        <div class="stat-card total">
          <div class="stat-icon total-cost">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="1" y="4" width="22" height="16" rx="2" ry="2"/>
              <line x1="1" y1="10" x2="23" y2="10"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">${{ formatNumber(data.total_usage?.total_cost || 0, 2) }}</div>
            <div class="stat-label">总消费</div>
          </div>
        </div>

        <div class="stat-card total">
          <div class="stat-icon total-token">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 16V8a2 2 0 00-1-1.73l-7-4a2 2 0 00-2 0l-7 4A2 2 0 003 8v8a2 2 0 001 1.73l7 4a2 2 0 002 0l7-4A2 2 0 0021 16z"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatLargeNumber(data.total_usage?.total_tokens || 0) }}</div>
            <div class="stat-label">总 Token</div>
          </div>
        </div>

        <div class="stat-card total">
          <div class="stat-icon total-request">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 11-5.93-9.14"/>
              <polyline points="22,4 12,14.01 9,11.01"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatLargeNumber(data.total_usage?.request_count || 0) }}</div>
            <div class="stat-label">总请求</div>
          </div>
        </div>

        <div class="stat-card total">
          <div class="stat-icon total-cache">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ formatLargeNumber(data.total_usage?.cache_creation_tokens || 0) }}</div>
            <div class="stat-label">缓存创建</div>
          </div>
        </div>
      </div>
    </section>

    <!-- 详情卡片 -->
    <section class="section">
      <div class="cards-grid two-col">
        <!-- Token 详情 -->
        <div class="detail-card">
          <div class="card-header">
            <h3 class="card-title">Token 使用详情</h3>
          </div>
          <div class="card-body">
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">输入 Token</span>
                <span class="detail-value">{{ formatNumber(data.today_usage?.input_tokens || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">输出 Token</span>
                <span class="detail-value">{{ formatNumber(data.today_usage?.output_tokens || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">缓存创建</span>
                <span class="detail-value">{{ formatNumber(data.today_usage?.cache_creation_tokens || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">缓存读取</span>
                <span class="detail-value">{{ formatNumber(data.today_usage?.cache_read_tokens || 0) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 用户统计 -->
        <div class="detail-card">
          <div class="card-header">
            <h3 class="card-title">用户统计</h3>
          </div>
          <div class="card-body">
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">总用户数</span>
                <span class="detail-value">{{ data.users?.total || 0 }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">今日活跃</span>
                <span class="detail-value highlight">{{ data.users?.active || 0 }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">今日新增</span>
                <span class="detail-value">
                  <span v-if="data.users?.new_today > 0" class="badge success">+{{ data.users?.new_today }}</span>
                  <span v-else>0</span>
                </span>
              </div>
              <div class="detail-item">
                <span class="detail-label">活跃率</span>
                <span class="detail-value">{{ data.users?.total > 0 ? ((data.users?.active / data.users?.total) * 100).toFixed(1) : 0 }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 账号状态 -->
    <section class="section">
      <h2 class="section-title">账号状态</h2>
      <div class="account-stats-grid">
        <div class="account-stat-card">
          <div class="account-value">{{ data.accounts?.total || 0 }}</div>
          <div class="account-label">总账号</div>
        </div>
        <div class="account-stat-card success">
          <div class="account-value">{{ data.accounts?.active || 0 }}</div>
          <div class="account-label">正常可用</div>
          <div class="account-indicator"></div>
        </div>
        <div class="account-stat-card warning">
          <div class="account-value">{{ data.accounts?.rate_limited || 0 }}</div>
          <div class="account-label">限流中</div>
          <div class="account-indicator"></div>
        </div>
        <div class="account-stat-card danger">
          <div class="account-value">{{ data.accounts?.invalid || 0 }}</div>
          <div class="account-label">无效/禁用</div>
          <div class="account-indicator"></div>
        </div>
      </div>
    </section>

    <!-- 系统资源 -->
    <section class="section">
      <h2 class="section-title">系统资源</h2>
      <div class="resource-grid">
        <!-- CPU -->
        <div class="resource-card">
          <div class="resource-header">
            <div class="resource-title">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="4" y="4" width="16" height="16" rx="2" ry="2"/>
                <rect x="9" y="9" width="6" height="6"/>
                <line x1="9" y1="1" x2="9" y2="4"/>
                <line x1="15" y1="1" x2="15" y2="4"/>
                <line x1="9" y1="20" x2="9" y2="23"/>
                <line x1="15" y1="20" x2="15" y2="23"/>
                <line x1="20" y1="9" x2="23" y2="9"/>
                <line x1="20" y1="14" x2="23" y2="14"/>
                <line x1="1" y1="9" x2="4" y2="9"/>
                <line x1="1" y1="14" x2="4" y2="14"/>
              </svg>
              <span>CPU</span>
            </div>
            <span class="resource-badge">{{ data.system?.cpu_cores || 0 }} 核</span>
          </div>
          <div class="progress-container">
            <div class="progress-bar" :style="{ width: (data.system?.cpu_usage || 0) + '%' }" :class="getProgressClass(data.system?.cpu_usage)"></div>
          </div>
          <div class="resource-percent">{{ (data.system?.cpu_usage || 0).toFixed(1) }}%</div>
        </div>

        <!-- 内存 -->
        <div class="resource-card">
          <div class="resource-header">
            <div class="resource-title">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="6" width="20" height="12" rx="2"/>
                <line x1="6" y1="10" x2="6" y2="14"/>
                <line x1="10" y1="10" x2="10" y2="14"/>
                <line x1="14" y1="10" x2="14" y2="14"/>
                <line x1="18" y1="10" x2="18" y2="14"/>
              </svg>
              <span>内存</span>
            </div>
            <span class="resource-badge">{{ formatBytes(data.system?.memory_total) }}</span>
          </div>
          <div class="progress-container">
            <div class="progress-bar" :style="{ width: (data.system?.memory_usage || 0) + '%' }" :class="getProgressClass(data.system?.memory_usage)"></div>
          </div>
          <div class="resource-percent">{{ (data.system?.memory_usage || 0).toFixed(1) }}%</div>
          <div class="resource-detail">
            已用 {{ formatBytes(data.system?.memory_used) }} / 可用 {{ formatBytes(data.system?.memory_free) }}
          </div>
        </div>

        <!-- 磁盘 -->
        <div class="resource-card">
          <div class="resource-header">
            <div class="resource-title">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <ellipse cx="12" cy="5" rx="9" ry="3"/>
                <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
                <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
              </svg>
              <span>磁盘</span>
            </div>
            <span class="resource-badge">{{ formatBytes(data.system?.disk_total) }}</span>
          </div>
          <div class="progress-container">
            <div class="progress-bar" :style="{ width: (data.system?.disk_usage || 0) + '%' }" :class="getProgressClass(data.system?.disk_usage)"></div>
          </div>
          <div class="resource-percent">{{ (data.system?.disk_usage || 0).toFixed(1) }}%</div>
          <div class="resource-detail">
            已用 {{ formatBytes(data.system?.disk_used) }} / 可用 {{ formatBytes(data.system?.disk_free) }}
          </div>
        </div>
      </div>
    </section>

    <!-- 数据库状态 -->
    <section class="section">
      <div class="cards-grid two-col">
        <!-- 内存缓存 -->
        <div class="detail-card">
          <div class="card-header">
            <h3 class="card-title">内存缓存</h3>
            <span class="status-badge success">运行中</span>
          </div>
          <div class="card-body">
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">会话绑定数</span>
                <span class="detail-value">{{ formatNumber(data.cache?.session_count || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">账户并发数</span>
                <span class="detail-value">{{ formatNumber(data.cache?.account_concurrency_count || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">用户并发数</span>
                <span class="detail-value">{{ formatNumber(data.cache?.user_concurrency_count || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">不可用标记</span>
                <span class="detail-value">{{ formatNumber(data.cache?.unavailable_count || 0) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- MySQL -->
        <div class="detail-card">
          <div class="card-header">
            <h3 class="card-title">MySQL</h3>
            <span :class="['status-badge', data.mysql?.connected ? 'success' : 'danger']">
              {{ data.mysql?.connected ? '已连接' : '未连接' }}
            </span>
          </div>
          <div class="card-body">
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">表数量</span>
                <span class="detail-value">{{ data.mysql?.table_count || 0 }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">数据大小</span>
                <span class="detail-value">{{ formatBytes(data.mysql?.data_size) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">索引大小</span>
                <span class="detail-value">{{ formatBytes(data.mysql?.index_size) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">总大小</span>
                <span class="detail-value highlight">{{ formatBytes(data.mysql?.total_size) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- 更新时间 -->
    <div class="update-time" v-if="data.updated_at">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <polyline points="12,6 12,12 16,14"/>
      </svg>
      最后更新: {{ formatTime(data.updated_at) }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from '@/utils/toast'

const loading = ref(false)
const data = ref({})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await api.getMonitorData()
    if (res.code === 0) {
      data.value = res.data
    } else {
      ElMessage.error(res.message || '获取监控数据失败')
    }
  } catch (err) {
    ElMessage.error('获取监控数据失败')
  } finally {
    loading.value = false
  }
}

const formatNumber = (num, decimals = 0) => {
  if (num === undefined || num === null) return '0'
  if (decimals > 0) {
    return num.toFixed(decimals)
  }
  return num.toLocaleString()
}

const formatLargeNumber = (num) => {
  if (num === undefined || num === null || num === 0) return '0'
  if (num >= 1000000000) return (num / 1000000000).toFixed(2) + 'B'
  if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toLocaleString()
}

const formatBytes = (bytes) => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(2) + ' ' + sizes[i]
}

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN')
}

const getProgressClass = (percentage) => {
  if (percentage < 60) return 'success'
  if (percentage < 80) return 'warning'
  return 'danger'
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.system-monitor-page {
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

.header-content {
  flex: 1;
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

/* 按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  cursor: pointer;
}

.btn svg {
  width: 16px;
  height: 16px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--apple-blue);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--apple-blue-hover);
}

.btn-spinner {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Section */
.section {
  margin-bottom: var(--apple-spacing-xl);
}

.section-title {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-md) 0;
}

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--apple-spacing-md);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-lg);
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.stat-card:hover {
  box-shadow: var(--apple-shadow-card-hover);
  transform: translateY(-2px);
}

.stat-card.total {
  border-left: 3px solid var(--apple-blue);
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

.stat-icon.cost { background: linear-gradient(135deg, #667eea, #764ba2); }
.stat-icon.token { background: linear-gradient(135deg, #f093fb, #f5576c); }
.stat-icon.request { background: linear-gradient(135deg, #4facfe, #00f2fe); }
.stat-icon.user { background: linear-gradient(135deg, #43e97b, #38f9d7); }
.stat-icon.total-cost { background: linear-gradient(135deg, #ff6b6b, #ee5a24); }
.stat-icon.total-token { background: linear-gradient(135deg, #a55eea, #8854d0); }
.stat-icon.total-request { background: linear-gradient(135deg, #20bf6b, #26de81); }
.stat-icon.total-cache { background: linear-gradient(135deg, #fa8231, #fd9644); }

.stat-info {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: var(--apple-text-xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.stat-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-top: var(--apple-spacing-xxs);
}

/* 详情卡片网格 */
.cards-grid {
  display: grid;
  gap: var(--apple-spacing-md);
}

.cards-grid.two-col {
  grid-template-columns: repeat(2, 1fr);
}

.detail-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.card-title {
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0;
}

.card-body {
  padding: var(--apple-spacing-lg);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-md);
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-sm);
  background: var(--apple-fill-quaternary);
  border-radius: var(--apple-radius-sm);
}

.detail-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.detail-value {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.detail-value.highlight {
  color: var(--apple-blue);
}

/* 状态徽章 */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: var(--apple-spacing-xxs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-full);
}

.status-badge.success {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.status-badge.danger {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

.badge {
  display: inline-flex;
  padding: 2px 8px;
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-full);
}

.badge.success {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

/* 账号状态网格 */
.account-stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--apple-spacing-md);
}

.account-stat-card {
  position: relative;
  text-align: center;
  padding: var(--apple-spacing-xl);
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.account-stat-card .account-value {
  font-size: var(--apple-text-3xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  line-height: 1;
}

.account-stat-card .account-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-top: var(--apple-spacing-sm);
}

.account-stat-card .account-indicator {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
}

.account-stat-card.success .account-value { color: var(--apple-green); }
.account-stat-card.success .account-indicator { background: var(--apple-green); }

.account-stat-card.warning .account-value { color: var(--apple-orange); }
.account-stat-card.warning .account-indicator { background: var(--apple-orange); }

.account-stat-card.danger .account-value { color: var(--apple-red); }
.account-stat-card.danger .account-indicator { background: var(--apple-red); }

/* 资源监控网格 */
.resource-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--apple-spacing-md);
}

.resource-card {
  padding: var(--apple-spacing-lg);
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
}

.resource-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--apple-spacing-md);
}

.resource-title {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.resource-title svg {
  width: 18px;
  height: 18px;
  color: var(--apple-text-secondary);
}

.resource-badge {
  font-size: var(--apple-text-xs);
  padding: var(--apple-spacing-xxs) var(--apple-spacing-sm);
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-secondary);
  border-radius: var(--apple-radius-full);
}

/* 进度条 */
.progress-container {
  height: 8px;
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-full);
  overflow: hidden;
  margin-bottom: var(--apple-spacing-sm);
}

.progress-bar {
  height: 100%;
  border-radius: var(--apple-radius-full);
  transition: width var(--apple-duration-normal) var(--apple-ease-default);
}

.progress-bar.success { background: var(--apple-green); }
.progress-bar.warning { background: var(--apple-orange); }
.progress-bar.danger { background: var(--apple-red); }

.resource-percent {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  text-align: center;
}

.resource-detail {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  text-align: center;
  margin-top: var(--apple-spacing-xs);
}

/* 更新时间 */
.update-time {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: var(--apple-spacing-xs);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
  margin-top: var(--apple-spacing-lg);
}

.update-time svg {
  width: 14px;
  height: 14px;
}

/* 响应式布局 */
@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .account-stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .resource-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-md);
  }

  .stats-grid,
  .account-stats-grid {
    grid-template-columns: 1fr;
  }

  .cards-grid.two-col {
    grid-template-columns: 1fr;
  }

  .resource-grid {
    grid-template-columns: 1fr;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>
