<!--
 * 文件作用：缓存管理页面 - Apple HIG 风格
 * 负责功能：
 *   - 缓存配置管理（TTL设置）
 *   - 会话列表（每个会话一行，显示到TTL结束）
 *   - 会话绑定详情查看
 *   - 不可用账号管理
 *   - 并发统计
 * 重要程度：⭐⭐⭐ 一般（缓存管理）
-->
<template>
  <div class="cache-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">缓存管理</h1>
        <p class="page-subtitle">管理会话缓存、不可用账号和并发控制</p>
      </div>
      <button class="btn btn-outline" @click="refreshAll">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23,4 23,10 17,10"/>
          <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
        </svg>
        刷新
      </button>
    </div>

    <!-- 缓存配置卡片 -->
    <div class="config-card">
      <div class="card-header">
        <h3>缓存配置</h3>
      </div>
      <div class="config-form" :class="{ loading: loadingConfig }">
        <div class="config-row">
          <div class="config-item">
            <label>粘性会话 TTL</label>
            <div class="input-group">
              <input v-model.number="configForm.session_ttl" type="number" min="1" max="1440" class="form-input" />
              <span class="input-suffix">分钟</span>
            </div>
          </div>
          <div class="config-item">
            <label>会话续期阈值</label>
            <div class="input-group">
              <input v-model.number="configForm.session_renewal_ttl" type="number" min="1" max="60" class="form-input" />
              <span class="input-suffix">分钟</span>
            </div>
          </div>
          <div class="config-item">
            <label>不可用标记 TTL</label>
            <div class="input-group">
              <input v-model.number="configForm.unavailable_ttl" type="number" min="1" max="60" class="form-input" />
              <span class="input-suffix">分钟</span>
            </div>
          </div>
          <div class="config-item">
            <label>并发计数 TTL</label>
            <div class="input-group">
              <input v-model.number="configForm.concurrency_ttl" type="number" min="1" max="60" class="form-input" />
              <span class="input-suffix">分钟</span>
            </div>
          </div>
          <div class="config-item action">
            <button class="btn btn-primary" @click="saveConfig" :disabled="savingConfig">
              {{ savingConfig ? '保存中...' : '保存配置' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tab 切换 -->
    <div class="tabs-container">
      <div class="tabs">
        <button
          :class="['tab', { active: activeTab === 'sessions' }]"
          @click="handleTabChange('sessions')"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 00-3-3.87"/>
            <path d="M16 3.13a4 4 0 010 7.75"/>
          </svg>
          会话列表 ({{ sessionTotal }})
        </button>
        <button
          :class="['tab', { active: activeTab === 'unavailable' }]"
          @click="handleTabChange('unavailable')"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
          </svg>
          不可用账号 ({{ unavailableList.length }})
        </button>
        <button
          :class="['tab', { active: activeTab === 'concurrency' }]"
          @click="handleTabChange('concurrency')"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 3v18h18"/>
            <path d="M18 9l-5 5-4-4-3 3"/>
          </svg>
          并发统计
        </button>
      </div>
    </div>

    <!-- 会话列表 -->
    <div v-show="activeTab === 'sessions'" class="tab-content">
      <div class="table-card">
        <div class="card-toolbar">
          <div class="search-box">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input
              v-model="sessionSearch"
              type="text"
              placeholder="搜索会话/账号/API Key..."
              @input="filterSessions"
            />
          </div>
          <button class="btn btn-danger btn-sm" @click="confirmClearAll" :disabled="sessions.length === 0">
            清除全部
          </button>
        </div>

        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>会话ID</th>
                <th>账号</th>
                <th>API Key</th>
                <th>平台</th>
                <th>模型</th>
                <th>绑定时间</th>
                <th>最后使用</th>
                <th>剩余时间</th>
                <th>客户端IP</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody v-if="!loadingSessions">
              <tr v-for="session in filteredSessions" :key="session.session_id">
                <td><code class="session-id">{{ parseSessionId(session.session_id) }}</code></td>
                <td>
                  <span class="badge badge-warning">{{ getAccountName(session.account_id) }}</span>
                </td>
                <td>
                  <code v-if="session.api_key_id" class="api-key-code">{{ getAPIKeyLabel(session.api_key_id) }}</code>
                  <span v-else class="text-muted">-</span>
                </td>
                <td><span class="badge badge-info">{{ session.platform || '-' }}</span></td>
                <td class="model-cell">{{ session.model || '-' }}</td>
                <td>{{ formatTime(session.bound_at) }}</td>
                <td>{{ formatTime(session.last_used_at) }}</td>
                <td>
                  <span :class="['ttl-badge', getTTLClass(session.remaining_ttl)]">
                    {{ formatTTL(session.remaining_ttl) }}
                  </span>
                </td>
                <td>{{ session.client_ip || '-' }}</td>
                <td>
                  <button class="action-btn danger" @click="removeSession(session.session_id)" title="移除">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <line x1="18" y1="6" x2="6" y2="18"/>
                      <line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="loadingSessions" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>

          <div v-if="!loadingSessions && sessions.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
            </svg>
            <span>暂无活跃会话</span>
          </div>
        </div>

        <div class="pagination-wrap" v-if="sessionTotal > 0">
          <div class="pagination-info">共 {{ sessionTotal }} 条记录</div>
          <div class="pagination-controls">
            <select v-model="sessionPageSize" @change="loadSessions" class="page-size-select">
              <option :value="20">20 条/页</option>
              <option :value="50">50 条/页</option>
              <option :value="100">100 条/页</option>
            </select>
            <div class="page-btns">
              <button class="page-btn" :disabled="sessionPage <= 1" @click="sessionPage--; loadSessions()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="15,18 9,12 15,6"/>
                </svg>
              </button>
              <span class="page-current">{{ sessionPage }} / {{ Math.ceil(sessionTotal / sessionPageSize) || 1 }}</span>
              <button class="page-btn" :disabled="sessionPage >= Math.ceil(sessionTotal / sessionPageSize)" @click="sessionPage++; loadSessions()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9,18 15,12 9,6"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 不可用账号 -->
    <div v-show="activeTab === 'unavailable'" class="tab-content">
      <div class="table-card">
        <div class="card-toolbar">
          <span class="toolbar-title">临时不可用账号</span>
          <button class="btn btn-outline btn-sm" @click="loadUnavailable">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23,4 23,10 17,10"/>
              <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
            </svg>
            刷新
          </button>
        </div>

        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>账号</th>
                <th>原因</th>
                <th>剩余时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody v-if="!loadingUnavailable">
              <tr v-for="item in unavailableList" :key="item.account_id">
                <td>
                  <span class="badge badge-warning">{{ getAccountName(item.account_id) }}</span>
                </td>
                <td class="reason-cell">{{ item.reason }}</td>
                <td>
                  <span class="ttl-badge danger">{{ formatTTL(item.remaining_ttl) }}</span>
                </td>
                <td>
                  <button class="btn btn-success btn-sm" @click="clearUnavailable(item.account_id)">
                    恢复
                  </button>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="loadingUnavailable" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>

          <div v-if="!loadingUnavailable && unavailableList.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M22 11.08V12a10 10 0 11-5.93-9.14"/>
              <polyline points="22,4 12,14.01 9,11.01"/>
            </svg>
            <span>暂无不可用账号</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 并发统计 -->
    <div v-show="activeTab === 'concurrency'" class="tab-content">
      <div class="concurrency-grid">
        <!-- 账号并发 -->
        <div class="table-card">
          <div class="card-toolbar">
            <span class="toolbar-title">账号并发 ({{ accountConcurrencyList.length }})</span>
          </div>
          <div class="table-container compact">
            <table class="data-table">
              <thead>
                <tr>
                  <th>账号</th>
                  <th>并发</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in accountConcurrencyList" :key="item.account_id">
                  <td>{{ getAccountName(item.account_id) }}</td>
                  <td>
                    <span :class="getConcurrencyClass(item.concurrency, item.max_concurrency)">
                      {{ item.concurrency }}
                    </span> / {{ item.max_concurrency }}
                  </td>
                  <td>
                    <button class="btn btn-danger btn-xs" @click="resetAccountConcurrency(item.account_id)">
                      重置
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
            <div v-if="accountConcurrencyList.length === 0" class="empty-state small">
              <span>暂无数据</span>
            </div>
          </div>
        </div>

      </div>
    </div>

    <!-- 清除全部确认弹窗 -->
    <Teleport to="body">
      <div v-if="clearAllDialogVisible" class="modal-overlay" @click.self="clearAllDialogVisible = false">
        <div class="modal modal-sm">
          <div class="modal-header danger">
            <div class="danger-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
                <line x1="12" y1="9" x2="12" y2="13"/>
                <line x1="12" y1="17" x2="12.01" y2="17"/>
              </svg>
            </div>
            <h2>确认清除</h2>
          </div>
          <div class="modal-body">
            <p class="delete-message">确定要清除所有会话缓存吗？此操作不可撤销。</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="clearAllDialogVisible = false">取消</button>
            <button class="btn btn-danger" @click="clearAllSessions">确认清除</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

// 配置
const loadingConfig = ref(false)
const savingConfig = ref(false)
const configForm = reactive({
  session_ttl: 60,
  session_renewal_ttl: 14,
  unavailable_ttl: 5,
  concurrency_ttl: 5
})

// Tab
const activeTab = ref('sessions')

// 会话列表
const loadingSessions = ref(false)
const sessions = ref([])
const sessionTotal = ref(0)
const sessionPage = ref(1)
const sessionPageSize = ref(20)
const sessionSearch = ref('')

// 账号和 API Key 名称缓存
const accountNames = ref({})
const apiKeyLabels = ref({})

// 不可用账号
const loadingUnavailable = ref(false)
const unavailableList = ref([])

// 清除确认
const clearAllDialogVisible = ref(false)

// 并发统计（从会话列表中计算）
const accountConcurrencyList = computed(() => {
  const map = new Map()
  sessions.value.forEach(s => {
    if (!map.has(s.account_id)) {
      map.set(s.account_id, { account_id: s.account_id, concurrency: 0, max_concurrency: 5 })
    }
    map.get(s.account_id).concurrency++
  })
  return Array.from(map.values())
})

// 过滤后的会话
const filteredSessions = computed(() => {
  if (!sessionSearch.value) return sessions.value
  const search = sessionSearch.value.toLowerCase()
  return sessions.value.filter(s => {
    return s.session_id?.toLowerCase().includes(search) ||
           getAccountName(s.account_id).toLowerCase().includes(search) ||
           getAPIKeyLabel(s.api_key_id).toLowerCase().includes(search)
  })
})

// 自动刷新定时器
let refreshTimer = null

function formatTime(time) {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}

function getConcurrencyClass(current, limit) {
  if (current >= limit) return 'concurrency-danger'
  if (current >= limit * 0.8) return 'concurrency-warning'
  return 'concurrency-success'
}

function parseSessionId(sessionId) {
  if (!sessionId) return '-'
  if (sessionId.startsWith('apikey:')) {
    const parts = sessionId.split(':')
    if (parts.length >= 3) {
      return parts.slice(2).join(':')
    }
    return `Key#${parts[1]}`
  }
  return sessionId
}

function formatTTL(seconds) {
  if (!seconds || seconds <= 0) return '已过期'
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分${seconds % 60}秒`
  return `${Math.floor(seconds / 3600)}时${Math.floor((seconds % 3600) / 60)}分`
}

function getTTLClass(seconds) {
  if (!seconds || seconds <= 0) return 'danger'
  if (seconds < 300) return 'warning'
  return 'success'
}

function getAccountName(accountId) {
  if (!accountId) return '-'
  return accountNames.value[accountId] || `#${accountId}`
}

function getAPIKeyLabel(apiKeyId) {
  if (!apiKeyId) return '-'
  return apiKeyLabels.value[apiKeyId] || 'sk-...'
}

function filterSessions() {
  // 搜索在 computed 中实现
}

async function fetchConfig() {
  loadingConfig.value = true
  try {
    const res = await api.getCacheConfig()
    const data = res.data || {}
    configForm.session_ttl = data.session_ttl || 60
    configForm.session_renewal_ttl = data.session_renewal_ttl || 14
    configForm.unavailable_ttl = data.unavailable_ttl || 5
    configForm.concurrency_ttl = data.concurrency_ttl || 5
  } catch (e) {
    console.error('Failed to fetch config:', e)
  } finally {
    loadingConfig.value = false
  }
}

async function saveConfig() {
  savingConfig.value = true
  try {
    await api.updateCacheConfig(configForm)
    ElMessage.success('配置已保存')
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    savingConfig.value = false
  }
}

async function loadSessions() {
  loadingSessions.value = true
  try {
    const offset = (sessionPage.value - 1) * sessionPageSize.value
    const res = await api.getCacheSessions({ offset, limit: sessionPageSize.value })
    sessions.value = res.data?.sessions || []
    sessionTotal.value = res.data?.total || 0

    const accountIds = [...new Set(sessions.value.map(s => s.account_id).filter(Boolean))]
    const apiKeyIds = [...new Set(sessions.value.map(s => s.api_key_id).filter(Boolean))]

    await loadAccountNames(accountIds)
    await loadAPIKeyLabels(apiKeyIds)
  } catch (e) {
    console.error('Failed to load sessions:', e)
  } finally {
    loadingSessions.value = false
  }
}

async function loadAccountNames(accountIds) {
  if (!accountIds.length) return
  try {
    const missing = accountIds.filter(id => id && !accountNames.value[id])
    if (!missing.length) return

    const res = await api.getAccounts({ page: 1, page_size: 1000 })
    const items = res.data?.items || res.data?.list || []
    items.forEach(acc => {
      if (acc?.id) {
        accountNames.value[acc.id] = acc.name || `账号${acc.id}`
      }
    })

    for (const id of missing) {
      if (!accountNames.value[id]) {
        try {
          const r = await api.getAccount(id)
          if (r.data?.id) {
            accountNames.value[r.data.id] = r.data.name || `账号${r.data.id}`
          }
        } catch {
          // ignore
        }
      }
    }
  } catch (e) {
    console.error('Failed to load account names:', e)
  }
}

async function loadAPIKeyLabels(apiKeyIds) {
  if (!apiKeyIds.length) return
  const missing = apiKeyIds.filter(id => id && !apiKeyLabels.value[id])
  if (!missing.length) return

  try {
    const res = await api.adminLookupAPIKeys(missing)
    const items = res.data?.items || []
    items.forEach(k => {
      if (!k?.id) return
      if (k.key_prefix) {
        apiKeyLabels.value[k.id] = k.key_prefix
      }
    })
  } catch (e) {
    console.error('Failed to load API key labels:', e)
  }
}

async function removeSession(sessionId) {
  try {
    await api.removeCacheSession(sessionId)
    ElMessage.success('会话已移除')
    loadSessions()
  } catch (e) {
    ElMessage.error('移除失败')
  }
}

function confirmClearAll() {
  clearAllDialogVisible.value = true
}

async function clearAllSessions() {
  try {
    await api.clearCache('sessions')
    ElMessage.success('所有会话已清除')
    sessions.value = []
    sessionTotal.value = 0
    clearAllDialogVisible.value = false
  } catch (e) {
    ElMessage.error('清除失败')
  }
}

async function loadUnavailable() {
  loadingUnavailable.value = true
  try {
    const res = await api.getUnavailableAccounts()
    unavailableList.value = res.data?.accounts || []
  } catch (e) {
    console.error('Failed to load unavailable accounts:', e)
  } finally {
    loadingUnavailable.value = false
  }
}

async function clearUnavailable(accountId) {
  try {
    await api.clearAccountUnavailable(accountId)
    ElMessage.success('已恢复账号')
    loadUnavailable()
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

async function resetAccountConcurrency(accountId) {
  try {
    await api.resetAccountConcurrency(accountId)
    ElMessage.success('已重置')
    loadSessions()
  } catch (e) {
    ElMessage.error('重置失败')
  }
}

function handleTabChange(tab) {
  activeTab.value = tab
  if (tab === 'sessions') {
    loadSessions()
  } else if (tab === 'unavailable') {
    loadUnavailable()
  } else if (tab === 'concurrency') {
    loadSessions()
  }
}

function refreshAll() {
  fetchConfig()
  loadSessions()
  if (activeTab.value === 'unavailable') {
    loadUnavailable()
  }
}

function startAutoRefresh() {
  refreshTimer = setInterval(() => {
    if (activeTab.value === 'sessions') {
      loadSessions()
    }
  }, 30000)
}

function stopAutoRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  fetchConfig()
  loadSessions()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.cache-page {
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

/* 配置卡片 */
.config-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  margin-bottom: var(--apple-spacing-xl);
  overflow: hidden;
}

.card-header {
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.card-header h3 {
  margin: 0;
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.config-form {
  padding: var(--apple-spacing-lg);
  transition: opacity var(--apple-duration-fast);
}

.config-form.loading {
  opacity: 0.5;
  pointer-events: none;
}

.config-row {
  display: flex;
  flex-wrap: wrap;
  gap: var(--apple-spacing-lg);
  align-items: flex-end;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.config-item label {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.config-item.action {
  flex-direction: row;
  align-items: flex-end;
}

.input-group {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.input-group .form-input {
  width: 100px;
}

.input-suffix {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

/* Tabs */
.tabs-container {
  margin-bottom: var(--apple-spacing-lg);
}

.tabs {
  display: flex;
  gap: var(--apple-spacing-xs);
  background: var(--apple-bg-secondary);
  padding: var(--apple-spacing-xs);
  border-radius: var(--apple-radius-md);
}

.tab {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.tab svg {
  width: 16px;
  height: 16px;
}

.tab:hover {
  color: var(--apple-text-primary);
  background: var(--apple-bg-tertiary);
}

.tab.active {
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  box-shadow: var(--apple-shadow-xs);
}

/* 表格卡片 */
.table-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.card-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.toolbar-title {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.search-box {
  display: flex;
  align-items: center;
  max-width: 300px;
  position: relative;
}

.search-box svg {
  position: absolute;
  left: var(--apple-spacing-sm);
  width: 16px;
  height: 16px;
  color: var(--apple-text-tertiary);
  pointer-events: none;
}

.search-box input {
  width: 100%;
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  padding-left: 36px;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-secondary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
}

.search-box input:focus {
  outline: none;
  border-color: var(--apple-blue);
}

.table-container {
  overflow-x: auto;
}

.table-container.compact {
  max-height: 300px;
  overflow-y: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--apple-text-sm);
}

.data-table th,
.data-table td {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  text-align: left;
  border-bottom: 1px solid var(--apple-separator);
}

.data-table th {
  background: var(--apple-bg-secondary);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  white-space: nowrap;
  position: sticky;
  top: 0;
}

.data-table tbody tr:hover {
  background: var(--apple-bg-secondary);
}

.session-id {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.api-key-code {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
  background: var(--apple-bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--apple-radius-xs);
}

.model-cell {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.reason-cell {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.text-muted {
  color: var(--apple-text-tertiary);
}

/* 徽章 */
.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.badge-primary { background: var(--apple-blue-light); color: var(--apple-blue); }
.badge-success { background: var(--apple-green-light); color: var(--apple-green); }
.badge-warning { background: var(--apple-orange-light); color: var(--apple-orange); }
.badge-danger { background: var(--apple-red-light); color: var(--apple-red); }
.badge-info { background: var(--apple-fill-tertiary); color: var(--apple-text-secondary); }

.ttl-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.ttl-badge.success { background: var(--apple-green-light); color: var(--apple-green); }
.ttl-badge.warning { background: var(--apple-orange-light); color: var(--apple-orange); }
.ttl-badge.danger { background: var(--apple-red-light); color: var(--apple-red); }

.concurrency-success { color: var(--apple-green); font-weight: var(--apple-font-semibold); }
.concurrency-warning { color: var(--apple-orange); font-weight: var(--apple-font-semibold); }
.concurrency-danger { color: var(--apple-red); font-weight: var(--apple-font-semibold); }

/* 操作按钮 */
.action-btn {
  width: 28px;
  height: 28px;
  border-radius: var(--apple-radius-xs);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-secondary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.action-btn svg {
  width: 14px;
  height: 14px;
}

.action-btn:hover {
  background: var(--apple-blue);
  color: white;
}

.action-btn.danger:hover {
  background: var(--apple-red);
}

/* 分页 */
.pagination-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
}

.pagination-info {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
}

.page-size-select {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  background: var(--apple-bg-primary);
}

.page-btns {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.page-btn {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-secondary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.page-btn svg {
  width: 16px;
  height: 16px;
}

.page-btn:hover:not(:disabled) {
  background: var(--apple-blue);
  color: white;
}

.page-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.page-current {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  min-width: 60px;
  text-align: center;
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

.loading-state.small, .empty-state.small {
  padding: var(--apple-spacing-xl);
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

/* 并发网格 */
.concurrency-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-lg);
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
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  cursor: pointer;
}

.btn svg {
  width: 16px;
  height: 16px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-sm {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
}

.btn-xs {
  padding: 2px var(--apple-spacing-xs);
  font-size: var(--apple-text-xs);
}

.btn-primary {
  background: var(--apple-blue);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--apple-blue-hover);
}

.btn-secondary {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--apple-fill-secondary);
}

.btn-danger {
  background: var(--apple-red);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #e6362d;
}

.btn-success {
  background: var(--apple-green);
  color: white;
}

.btn-success:hover:not(:disabled) {
  background: #2db84d;
}

.btn-outline {
  background: transparent;
  color: var(--apple-blue);
  border: 1px solid var(--apple-blue);
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-blue-light);
}

/* 表单 */
.form-input {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-input:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--apple-z-modal);
  padding: var(--apple-spacing-xl);
}

.modal {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  box-shadow: var(--apple-shadow-modal);
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal.modal-sm { max-width: 400px; }

.modal-header {
  padding: var(--apple-spacing-xl);
  border-bottom: 1px solid var(--apple-separator);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header.danger {
  flex-direction: column;
  text-align: center;
  gap: var(--apple-spacing-md);
}

.danger-icon {
  width: 56px;
  height: 56px;
  background: var(--apple-orange-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.danger-icon svg {
  width: 28px;
  height: 28px;
  color: var(--apple-orange);
}

.modal-header h2 {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0;
}

.modal-body {
  padding: var(--apple-spacing-xl);
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: var(--apple-spacing-lg) var(--apple-spacing-xl);
  border-top: 1px solid var(--apple-separator);
  display: flex;
  justify-content: flex-end;
  gap: var(--apple-spacing-sm);
}

.delete-message {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  text-align: center;
  margin: 0;
}

/* 响应式 */
@media (max-width: 1024px) {
  .config-row {
    flex-direction: column;
  }

  .concurrency-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }

  .tabs {
    flex-direction: column;
  }

  .card-toolbar {
    flex-direction: column;
    gap: var(--apple-spacing-sm);
  }

  .search-box {
    max-width: 100%;
    width: 100%;
  }
}
</style>
