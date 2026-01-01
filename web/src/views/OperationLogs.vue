<!--
 * 文件作用：操作日志页面 - Apple HIG 风格
 * 负责功能：
 *   - 操作日志列表和筛选
 *   - 模块/操作类型分类
 *   - 日志详情查看
 *   - 历史日志清理
 * 重要程度：⭐⭐ 辅助（审计日志）
-->
<template>
  <div class="operation-logs-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">操作日志</h1>
        <p class="page-subtitle">记录系统所有操作行为，便于审计和追溯</p>
      </div>
      <div class="header-actions">
        <button class="btn btn-outline" @click="loadLogs" :disabled="loading">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: loading }">
            <polyline points="23,4 23,10 17,10"/>
            <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
          </svg>
          刷新
        </button>
        <button class="btn btn-warning" @click="confirmCleanup">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
          </svg>
          清理旧日志
        </button>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="filter-card">
      <div class="filter-row">
        <div class="filter-group">
          <label class="filter-label">模块</label>
          <select v-model="filters.module" @change="handleFilterChange" class="filter-select">
            <option value="">全部模块</option>
            <option v-for="(label, key) in moduleLabels" :key="key" :value="key">{{ label }}</option>
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">操作</label>
          <select v-model="filters.action" @change="handleFilterChange" class="filter-select">
            <option value="">全部操作</option>
            <option v-for="(label, key) in actionLabels" :key="key" :value="key">{{ label }}</option>
          </select>
        </div>
        <div class="filter-group">
          <label class="filter-label">时间范围</label>
          <div class="date-inputs">
            <input type="date" v-model="dateStart" class="date-input" @change="handleDateChange" />
            <span class="date-sep">至</span>
            <input type="date" v-model="dateEnd" class="date-input" @change="handleDateChange" />
          </div>
        </div>
        <div class="filter-group search-group">
          <label class="filter-label">搜索</label>
          <div class="search-box">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input v-model="filters.search" placeholder="用户/描述/路径..." @input="handleSearch" />
          </div>
        </div>
      </div>
    </div>

    <!-- 日志列表 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th width="50">#</th>
              <th width="160">时间</th>
              <th width="120">用户</th>
              <th width="90">模块</th>
              <th width="80">操作</th>
              <th>描述</th>
              <th width="180">路径</th>
              <th width="80">结果</th>
              <th width="70">耗时</th>
              <th width="60">操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="(row, index) in logs" :key="row.id">
              <td><span class="row-index">{{ (pagination.page - 1) * pagination.pageSize + index + 1 }}</span></td>
              <td><span class="log-time">{{ formatTime(row.created_at) }}</span></td>
              <td>
                <div class="user-cell">
                  <span class="username">{{ row.username || '-' }}</span>
                  <span class="user-ip">{{ row.ip }}</span>
                </div>
              </td>
              <td><span :class="['module-badge', row.module]">{{ getModuleLabel(row.module) }}</span></td>
              <td><span :class="['action-badge', row.action]">{{ getActionLabel(row.action) }}</span></td>
              <td>
                <div class="description-cell">
                  <span class="description">{{ row.description }}</span>
                  <span class="target" v-if="row.target_name">→ {{ row.target_name }}</span>
                </div>
              </td>
              <td>
                <div class="path-cell">
                  <span :class="['method-badge', row.method]">{{ row.method }}</span>
                  <span class="path">{{ row.path }}</span>
                </div>
              </td>
              <td>
                <span :class="['result-badge', row.response_code === 0 ? 'success' : 'error']">
                  {{ row.response_code === 0 ? '成功' : '失败' }}
                </span>
              </td>
              <td><span class="duration">{{ row.duration }}ms</span></td>
              <td>
                <button class="action-btn" @click="showDetail(row)" title="详情">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                    <circle cx="12" cy="12" r="3"/>
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <div v-if="!loading && logs.length === 0" class="empty-state">
          <span>暂无日志记录</span>
        </div>
      </div>

      <div v-if="pagination.total > 0" class="table-footer">
        <div class="pagination-info">共 {{ pagination.total }} 条</div>
        <div class="pagination-controls">
          <select v-model="pagination.pageSize" @change="loadLogs" class="page-size-select">
            <option :value="20">20 条/页</option>
            <option :value="50">50 条/页</option>
            <option :value="100">100 条/页</option>
          </select>
          <div class="page-btns">
            <button class="page-btn" :disabled="pagination.page <= 1" @click="pagination.page--; loadLogs()">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15,18 9,12 15,6"/></svg>
            </button>
            <span class="page-current">{{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) || 1 }}</span>
            <button class="page-btn" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)" @click="pagination.page++; loadLogs()">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9,18 15,12 9,6"/></svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 详情弹窗 -->
    <Teleport to="body">
      <div v-if="detailVisible" class="modal-overlay" @click.self="detailVisible = false">
        <div class="modal modal-lg">
          <div class="modal-header">
            <h2>操作日志详情</h2>
            <button class="modal-close" @click="detailVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body" v-if="currentLog">
            <div class="detail-grid">
              <div class="detail-item"><span class="detail-label">ID</span><span class="detail-value">{{ currentLog.id }}</span></div>
              <div class="detail-item"><span class="detail-label">时间</span><span class="detail-value">{{ formatTime(currentLog.created_at) }}</span></div>
              <div class="detail-item"><span class="detail-label">用户</span><span class="detail-value">{{ currentLog.username }} (ID: {{ currentLog.user_id }})</span></div>
              <div class="detail-item"><span class="detail-label">IP</span><span class="detail-value mono">{{ currentLog.ip }}</span></div>
              <div class="detail-item"><span class="detail-label">模块</span><span class="detail-value">{{ getModuleLabel(currentLog.module) }}</span></div>
              <div class="detail-item"><span class="detail-label">操作</span><span class="detail-value">{{ getActionLabel(currentLog.action) }}</span></div>
              <div class="detail-item"><span class="detail-label">目标ID</span><span class="detail-value">{{ currentLog.target_id || '-' }}</span></div>
              <div class="detail-item"><span class="detail-label">目标名称</span><span class="detail-value">{{ currentLog.target_name || '-' }}</span></div>
              <div class="detail-item full"><span class="detail-label">描述</span><span class="detail-value">{{ currentLog.description }}</span></div>
              <div class="detail-item"><span class="detail-label">请求方法</span><span class="detail-value mono">{{ currentLog.method }}</span></div>
              <div class="detail-item"><span class="detail-label">请求路径</span><span class="detail-value mono">{{ currentLog.path }}</span></div>
              <div class="detail-item"><span class="detail-label">响应码</span><span class="detail-value">{{ currentLog.response_code }}</span></div>
              <div class="detail-item"><span class="detail-label">响应消息</span><span class="detail-value">{{ currentLog.response_msg || '-' }}</span></div>
              <div class="detail-item"><span class="detail-label">耗时</span><span class="detail-value">{{ currentLog.duration }}ms</span></div>
              <div class="detail-item full"><span class="detail-label">User-Agent</span><span class="detail-value ua">{{ currentLog.user_agent || '-' }}</span></div>
              <div v-if="currentLog.request_body" class="detail-item full">
                <span class="detail-label">请求体</span>
                <pre class="request-body">{{ formatJSON(currentLog.request_body) }}</pre>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="detailVisible = false">关闭</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 清理确认弹窗 -->
    <Teleport to="body">
      <div v-if="cleanupVisible" class="modal-overlay" @click.self="cleanupVisible = false">
        <div class="modal modal-sm">
          <div class="modal-header danger">
            <div class="danger-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
                <line x1="12" y1="9" x2="12" y2="13"/>
                <line x1="12" y1="17" x2="12.01" y2="17"/>
              </svg>
            </div>
            <h2>确认清理</h2>
          </div>
          <div class="modal-body">
            <p class="delete-message">确定要清理 90 天前的日志吗？此操作不可恢复。</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="cleanupVisible = false">取消</button>
            <button class="btn btn-danger" :disabled="cleaning" @click="handleCleanup">
              <span v-if="cleaning" class="btn-loading"></span>
              {{ cleaning ? '清理中...' : '确认清理' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const logs = ref([])
const detailVisible = ref(false)
const currentLog = ref(null)
const cleanupVisible = ref(false)
const cleaning = ref(false)
const dateStart = ref('')
const dateEnd = ref('')

const pagination = reactive({ page: 1, pageSize: 20, total: 0 })
const filters = reactive({ module: '', action: '', search: '', start_time: '', end_time: '' })

const moduleLabels = {
  auth: '认证', user: '用户', account: '账户', apikey: 'API Key',
  model: '模型', config: '配置', cache: '缓存', proxy: '代理',
  package: '套餐', group: '分组', system: '系统'
}

const actionLabels = {
  login: '登录', create: '创建', update: '更新', delete: '删除', clear: '清除'
}

function getModuleLabel(module) { return moduleLabels[module] || module }
function getActionLabel(action) { return actionLabels[action] || action }

function formatTime(dateStr) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

function formatJSON(str) {
  try { return JSON.stringify(JSON.parse(str), null, 2) }
  catch { return str }
}

let searchTimer = null
function handleSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => { pagination.page = 1; loadLogs() }, 300)
}

onUnmounted(() => { if (searchTimer) clearTimeout(searchTimer) })

function handleDateChange() {
  filters.start_time = dateStart.value || ''
  filters.end_time = dateEnd.value || ''
  pagination.page = 1
  loadLogs()
}

function handleFilterChange() {
  pagination.page = 1
  loadLogs()
}

async function loadLogs() {
  loading.value = true
  try {
    const params = { page: pagination.page, page_size: pagination.pageSize, ...filters }
    Object.keys(params).forEach(k => { if (!params[k]) delete params[k] })
    const res = await api.getOperationLogs(params)
    logs.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (e) {
    console.error('Failed to load logs:', e)
  } finally {
    loading.value = false
  }
}

function showDetail(row) {
  currentLog.value = row
  detailVisible.value = true
}

function confirmCleanup() { cleanupVisible.value = true }

async function handleCleanup() {
  cleaning.value = true
  try {
    const res = await api.cleanupOperationLogs(90)
    ElMessage.success(`清理完成，删除了 ${res.data.deleted} 条记录`)
    cleanupVisible.value = false
    loadLogs()
  } catch (e) {
    ElMessage.error('清理失败')
  } finally {
    cleaning.value = false
  }
}

onMounted(() => { loadLogs() })
</script>

<style scoped>
.operation-logs-page { max-width: 1400px; margin: 0 auto; }

/* 页面标题 */
.page-header {
  display: flex; justify-content: space-between; align-items: flex-start;
  margin-bottom: var(--apple-spacing-xl);
}
.header-content { flex: 1; }
.page-title {
  font-size: var(--apple-text-3xl); font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary); margin: 0 0 var(--apple-spacing-xs) 0;
}
.page-subtitle { font-size: var(--apple-text-base); color: var(--apple-text-secondary); margin: 0; }
.header-actions { display: flex; gap: var(--apple-spacing-sm); }

/* 筛选卡片 */
.filter-card {
  background: var(--apple-bg-primary); border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card); padding: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-xl);
}
.filter-row { display: flex; align-items: flex-end; gap: var(--apple-spacing-md); flex-wrap: wrap; }
.filter-group { display: flex; flex-direction: column; gap: var(--apple-spacing-xxs); }
.filter-label { font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); }
.filter-select, .date-input {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm); font-size: var(--apple-text-sm);
  border: 1px solid var(--apple-separator-opaque); border-radius: var(--apple-radius-sm);
  background: var(--apple-bg-primary);
}
.filter-select:focus, .date-input:focus { outline: none; border-color: var(--apple-blue); }
.date-inputs { display: flex; align-items: center; gap: var(--apple-spacing-xs); }
.date-sep { font-size: var(--apple-text-sm); color: var(--apple-text-tertiary); }
.search-group { flex: 1; min-width: 200px; }
.search-box {
  display: flex; align-items: center; position: relative;
}
.search-box svg {
  position: absolute; left: var(--apple-spacing-sm); width: 16px; height: 16px;
  color: var(--apple-text-tertiary); pointer-events: none;
}
.search-box input {
  width: 100%; padding: var(--apple-spacing-xs) var(--apple-spacing-md); padding-left: 36px;
  font-size: var(--apple-text-sm); border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm); background: var(--apple-bg-primary);
}
.search-box input:focus { outline: none; border-color: var(--apple-blue); }

/* 表格 */
.table-card {
  background: var(--apple-bg-primary); border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card); overflow: hidden;
}
.table-container { overflow-x: auto; }
.data-table { width: 100%; border-collapse: collapse; font-size: var(--apple-text-sm); }
.data-table th, .data-table td {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md); text-align: left;
  border-bottom: 1px solid var(--apple-separator);
}
.data-table th {
  background: var(--apple-bg-secondary); font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary); white-space: nowrap;
}
.data-table tbody tr:hover { background: var(--apple-bg-secondary); }

.row-index { font-family: var(--apple-font-mono); font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); }
.log-time { font-family: var(--apple-font-mono); font-size: var(--apple-text-xs); color: var(--apple-text-secondary); }

.user-cell { display: flex; flex-direction: column; }
.username { font-weight: var(--apple-font-semibold); color: var(--apple-text-primary); }
.user-ip { font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); font-family: var(--apple-font-mono); }

/* 徽章 */
.module-badge, .action-badge, .method-badge, .result-badge {
  display: inline-block; padding: 2px 8px; border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs); font-weight: var(--apple-font-semibold);
}
.module-badge { background: var(--apple-fill-tertiary); color: var(--apple-text-secondary); }
.module-badge.auth { background: #fef3c7; color: #d97706; }
.module-badge.user { background: #dbeafe; color: #2563eb; }
.module-badge.account { background: #dcfce7; color: #16a34a; }
.module-badge.apikey { background: #e0e7ff; color: #4f46e5; }
.module-badge.config { background: #fef3c7; color: #d97706; }
.module-badge.cache { background: #cffafe; color: #0891b2; }

.action-badge.login { background: #fef3c7; color: #d97706; }
.action-badge.create { background: #dcfce7; color: #16a34a; }
.action-badge.update { background: #dbeafe; color: #2563eb; }
.action-badge.delete { background: #fee2e2; color: #dc2626; }
.action-badge.clear { background: #fef3c7; color: #d97706; }

.method-badge { background: #dcfce7; color: #16a34a; }
.method-badge.PUT { background: #fef3c7; color: #d97706; }
.method-badge.DELETE { background: #fee2e2; color: #dc2626; }
.method-badge.GET { background: #e0e7ff; color: #4f46e5; }

.result-badge.success { background: var(--apple-green-light); color: var(--apple-green); }
.result-badge.error { background: var(--apple-red-light); color: var(--apple-red); }

.description-cell { display: flex; flex-direction: column; gap: 2px; }
.description { color: var(--apple-text-primary); }
.target { font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); }

.path-cell { display: flex; align-items: center; gap: var(--apple-spacing-xs); }
.path { font-family: var(--apple-font-mono); font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 120px; }

.duration { font-family: var(--apple-font-mono); font-size: var(--apple-text-xs); color: var(--apple-text-secondary); }

.action-btn {
  width: 28px; height: 28px; border-radius: var(--apple-radius-sm);
  display: flex; align-items: center; justify-content: center;
  background: var(--apple-fill-quaternary); color: var(--apple-text-secondary);
  transition: all var(--apple-duration-fast);
}
.action-btn svg { width: 14px; height: 14px; }
.action-btn:hover { background: var(--apple-blue); color: white; }

/* 分页 */
.table-footer {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg); border-top: 1px solid var(--apple-separator);
  display: flex; justify-content: space-between; align-items: center;
}
.pagination-info { font-size: var(--apple-text-sm); color: var(--apple-text-secondary); }
.pagination-controls { display: flex; align-items: center; gap: var(--apple-spacing-sm); }
.page-size-select { padding: var(--apple-spacing-xs) var(--apple-spacing-sm); font-size: var(--apple-text-sm); border: 1px solid var(--apple-separator-opaque); border-radius: var(--apple-radius-sm); }
.page-btns { display: flex; align-items: center; gap: var(--apple-spacing-xs); }
.page-btn { width: 28px; height: 28px; border-radius: var(--apple-radius-sm); display: flex; align-items: center; justify-content: center; background: var(--apple-fill-quaternary); color: var(--apple-text-secondary); transition: all var(--apple-duration-fast); }
.page-btn svg { width: 14px; height: 14px; }
.page-btn:hover:not(:disabled) { background: var(--apple-blue); color: white; }
.page-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.page-current { font-size: var(--apple-text-sm); color: var(--apple-text-primary); min-width: 60px; text-align: center; }

/* 加载和空状态 */
.loading-state, .empty-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: var(--apple-spacing-xxl); color: var(--apple-text-tertiary); }
.loading-spinner { width: 24px; height: 24px; border: 2px solid var(--apple-fill-tertiary); border-top-color: var(--apple-blue); border-radius: 50%; animation: spin 1s linear infinite; margin-bottom: var(--apple-spacing-sm); }
@keyframes spin { to { transform: rotate(360deg); } }
.spinning { animation: spin 1s linear infinite; }

/* 按钮 */
.btn { display: inline-flex; align-items: center; justify-content: center; gap: var(--apple-spacing-xs); padding: var(--apple-spacing-sm) var(--apple-spacing-lg); font-size: var(--apple-text-sm); font-weight: var(--apple-font-medium); border-radius: var(--apple-radius-sm); transition: all var(--apple-duration-fast); cursor: pointer; }
.btn svg { width: 16px; height: 16px; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-primary { background: var(--apple-blue); color: white; }
.btn-primary:hover:not(:disabled) { background: var(--apple-blue-hover); }
.btn-secondary { background: var(--apple-fill-tertiary); color: var(--apple-text-primary); }
.btn-secondary:hover:not(:disabled) { background: var(--apple-fill-secondary); }
.btn-warning { background: var(--apple-orange); color: white; }
.btn-warning:hover:not(:disabled) { background: #e68900; }
.btn-danger { background: var(--apple-red); color: white; }
.btn-danger:hover:not(:disabled) { background: #e6362d; }
.btn-outline { background: transparent; color: var(--apple-blue); border: 1px solid var(--apple-blue); }
.btn-outline:hover:not(:disabled) { background: var(--apple-blue-light); }
.btn-loading { width: 14px; height: 14px; border: 2px solid rgba(255, 255, 255, 0.3); border-top-color: white; border-radius: 50%; animation: spin 1s linear infinite; }

/* 模态框 */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.5); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: var(--apple-z-modal); padding: var(--apple-spacing-xl); }
.modal { background: var(--apple-bg-primary); border-radius: var(--apple-radius-xl); box-shadow: var(--apple-shadow-modal); width: 100%; max-width: 500px; max-height: 90vh; overflow: hidden; display: flex; flex-direction: column; }
.modal.modal-sm { max-width: 400px; }
.modal.modal-lg { max-width: 700px; }
.modal-header { padding: var(--apple-spacing-xl); border-bottom: 1px solid var(--apple-separator); display: flex; align-items: center; justify-content: space-between; }
.modal-header.danger { flex-direction: column; text-align: center; gap: var(--apple-spacing-md); }
.danger-icon { width: 56px; height: 56px; background: var(--apple-orange-light); border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.danger-icon svg { width: 28px; height: 28px; color: var(--apple-orange); }
.modal-header h2 { font-size: var(--apple-text-lg); font-weight: var(--apple-font-semibold); color: var(--apple-text-primary); margin: 0; }
.modal-close { width: 32px; height: 32px; border-radius: var(--apple-radius-sm); display: flex; align-items: center; justify-content: center; color: var(--apple-text-tertiary); transition: all var(--apple-duration-fast); }
.modal-close svg { width: 18px; height: 18px; }
.modal-close:hover { background: var(--apple-fill-tertiary); color: var(--apple-text-primary); }
.modal-body { padding: var(--apple-spacing-xl); overflow-y: auto; flex: 1; }
.modal-footer { padding: var(--apple-spacing-lg) var(--apple-spacing-xl); border-top: 1px solid var(--apple-separator); display: flex; justify-content: flex-end; gap: var(--apple-spacing-sm); }
.delete-message { font-size: var(--apple-text-base); color: var(--apple-text-secondary); text-align: center; margin: 0; }

/* 详情网格 */
.detail-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: var(--apple-spacing-md); }
.detail-item { display: flex; flex-direction: column; gap: 4px; padding: var(--apple-spacing-sm); background: var(--apple-bg-secondary); border-radius: var(--apple-radius-sm); }
.detail-item.full { grid-column: span 2; }
.detail-label { font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); }
.detail-value { font-size: var(--apple-text-sm); color: var(--apple-text-primary); word-break: break-all; }
.detail-value.mono { font-family: var(--apple-font-mono); }
.detail-value.ua { font-size: var(--apple-text-xs); }
.request-body { background: var(--apple-fill-quaternary); padding: var(--apple-spacing-sm); border-radius: var(--apple-radius-sm); font-family: var(--apple-font-mono); font-size: var(--apple-text-xs); white-space: pre-wrap; word-break: break-all; margin: var(--apple-spacing-xs) 0 0 0; max-height: 200px; overflow: auto; }

/* 响应式 */
@media (max-width: 768px) {
  .page-header { flex-direction: column; gap: var(--apple-spacing-lg); }
  .filter-row { flex-direction: column; align-items: stretch; }
  .filter-group, .search-group { width: 100%; }
  .detail-grid { grid-template-columns: 1fr; }
  .detail-item.full { grid-column: span 1; }
}
</style>
