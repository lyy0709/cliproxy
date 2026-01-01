<!--
 * 文件作用：API Key 管理页面（管理后台）- Apple HIG 风格
 * 负责功能：
 *   - 查看所有用户的 API Key 列表
 *   - 搜索和过滤 API Key
 *   - 启用/禁用 API Key
 *   - 删除 API Key
 *   - 查看 API Key 详情和使用日志
 * 重要程度：⭐⭐⭐⭐ 重要（管理员核心功能）
-->
<template>
  <div class="apikeys-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">API Key 管理</h1>
        <p class="page-subtitle">查看和管理所有 API Key</p>
      </div>
      <button class="btn btn-outline" @click="fetchAPIKeys">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23,4 23,10 17,10"/>
          <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
        </svg>
        刷新
      </button>
    </div>

    <!-- API Key 表格 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th class="col-id">ID</th>
              <th class="col-name">名称</th>
              <th class="col-prefix">Key 前缀</th>
              <th class="col-status">状态</th>
              <th class="col-rate">限速</th>
              <th class="col-requests">请求数</th>
              <th class="col-cost">费用</th>
              <th class="col-time">创建时间</th>
              <th class="col-actions">操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="key in apiKeys" :key="key.id">
              <td class="col-id">{{ key.id }}</td>
              <td class="col-name">
                <span class="key-name">{{ key.name }}</span>
              </td>
              <td class="col-prefix">
                <code class="key-prefix">{{ key.key_prefix }}...</code>
              </td>
              <td class="col-status">
                <span :class="['badge', getStatusType(key)]">
                  {{ getStatusLabel(key) }}
                </span>
              </td>
              <td class="col-rate">
                <span class="rate-info">{{ key.rate_limit }}/分</span>
              </td>
              <td class="col-requests">{{ formatNumber(key.request_count) }}</td>
              <td class="col-cost">
                <span class="cost">${{ (key.cost_used || 0).toFixed(4) }}</span>
              </td>
              <td class="col-time">{{ formatDate(key.created_at) }}</td>
              <td class="col-actions">
                <div class="action-group">
                  <button class="action-btn" @click="viewLogs(key)" title="日志">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                      <polyline points="14,2 14,8 20,8"/>
                      <line x1="16" y1="13" x2="8" y2="13"/>
                      <line x1="16" y1="17" x2="8" y2="17"/>
                    </svg>
                  </button>
                  <button
                    :class="['action-btn', key.status === 'active' ? 'warning' : 'success']"
                    @click="handleToggle(key)"
                    :title="key.status === 'active' ? '禁用' : '启用'"
                  >
                    <svg v-if="key.status === 'active'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                      <path d="M7 11V7a5 5 0 0110 0v4"/>
                    </svg>
                    <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                      <path d="M7 11V7a5 5 0 019.9-1"/>
                    </svg>
                  </button>
                  <button class="action-btn danger" @click="confirmDelete(key)" title="删除">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3,6 5,6 21,6"/>
                      <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && apiKeys.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
          </svg>
          <span>暂无 API Key</span>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrap" v-if="pagination.total > 0">
        <div class="pagination-info">
          共 {{ pagination.total }} 条记录
        </div>
        <div class="pagination-controls">
          <select v-model="pagination.pageSize" @change="fetchAPIKeys" class="page-size-select">
            <option :value="10">10 条/页</option>
            <option :value="20">20 条/页</option>
            <option :value="50">50 条/页</option>
          </select>
          <div class="page-btns">
            <button
              class="page-btn"
              :disabled="pagination.page <= 1"
              @click="pagination.page--; fetchAPIKeys()"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15,18 9,12 15,6"/>
              </svg>
            </button>
            <span class="page-current">{{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) || 1 }}</span>
            <button
              class="page-btn"
              :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)"
              @click="pagination.page++; fetchAPIKeys()"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9,18 15,12 9,6"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 使用日志弹窗 -->
    <Teleport to="body">
      <div v-if="logDialogVisible" class="modal-overlay" @click.self="logDialogVisible = false">
        <div class="modal modal-xl">
          <div class="modal-header">
            <div class="modal-title-group">
              <h2>{{ currentKey?.key_prefix }} 使用日志</h2>
              <div class="modal-badges">
                <span class="badge badge-success">请求数: {{ currentKey?.request_count || 0 }}</span>
                <span class="badge badge-warning">费用: ${{ (currentKey?.cost_used || 0).toFixed(4) }}</span>
              </div>
            </div>
            <button class="modal-close" @click="logDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <!-- 日志表格 -->
            <div class="log-table-container">
              <table class="data-table compact">
                <thead>
                  <tr>
                    <th>时间</th>
                    <th>模型</th>
                    <th>输入 Token</th>
                    <th>输出 Token</th>
                    <th>缓存创建</th>
                    <th>缓存读取</th>
                    <th>总 Token</th>
                    <th>费用</th>
                  </tr>
                </thead>
                <tbody v-if="!logLoading">
                  <tr v-for="(log, idx) in logs" :key="idx">
                    <td>{{ formatDate(log.timestamp) }}</td>
                    <td><code class="model-name">{{ log.model }}</code></td>
                    <td>{{ formatNumber(log.input_tokens || 0) }}</td>
                    <td>{{ formatNumber(log.output_tokens || 0) }}</td>
                    <td>{{ formatNumber(log.cache_creation_input_tokens || 0) }}</td>
                    <td>{{ formatNumber(log.cache_read_input_tokens || 0) }}</td>
                    <td>{{ formatNumber(log.total_tokens || 0) }}</td>
                    <td class="cost">${{ (log.total_cost || 0).toFixed(6) }}</td>
                  </tr>
                </tbody>
              </table>

              <div v-if="logLoading" class="loading-state small">
                <div class="loading-spinner"></div>
                <span>加载中...</span>
              </div>

              <div v-if="!logLoading && logs.length === 0" class="empty-state small">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                  <polyline points="14,2 14,8 20,8"/>
                </svg>
                <span>暂无日志</span>
              </div>
            </div>

            <!-- 日志分页 -->
            <div class="log-pagination" v-if="logPagination.total > 0">
              <div class="pagination-info">
                共 {{ logPagination.total }} 条记录
              </div>
              <div class="pagination-controls">
                <select v-model="logPagination.pageSize" @change="fetchLogs" class="page-size-select">
                  <option :value="20">20 条/页</option>
                  <option :value="50">50 条/页</option>
                  <option :value="100">100 条/页</option>
                </select>
                <div class="page-btns">
                  <button
                    class="page-btn"
                    :disabled="logPagination.page <= 1"
                    @click="logPagination.page--; fetchLogs()"
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="15,18 9,12 15,6"/>
                    </svg>
                  </button>
                  <span class="page-current">{{ logPagination.page }} / {{ Math.ceil(logPagination.total / logPagination.pageSize) || 1 }}</span>
                  <button
                    class="page-btn"
                    :disabled="logPagination.page >= Math.ceil(logPagination.total / logPagination.pageSize)"
                    @click="logPagination.page++; fetchLogs()"
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="9,18 15,12 9,6"/>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="logDialogVisible = false">关闭</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 删除确认弹窗 -->
    <Teleport to="body">
      <div v-if="deleteDialogVisible" class="modal-overlay" @click.self="deleteDialogVisible = false">
        <div class="modal modal-sm">
          <div class="modal-header danger">
            <div class="danger-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
            </div>
            <h2>确认删除</h2>
          </div>
          <div class="modal-body">
            <p class="delete-message">确定要删除此 API Key 吗？此操作不可撤销。</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="deleteDialogVisible = false">取消</button>
            <button class="btn btn-danger" :disabled="deleting" @click="handleDelete">
              <span v-if="deleting" class="btn-loading"></span>
              {{ deleting ? '删除中...' : '删除' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const apiKeys = ref([])
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })

// 日志相关
const logDialogVisible = ref(false)
const logLoading = ref(false)
const currentKey = ref(null)
const logs = ref([])
const logPagination = reactive({ page: 1, pageSize: 20, total: 0 })

// 删除相关
const deleteDialogVisible = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

function formatDate(str) {
  if (!str) return '-'
  return new Date(str).toLocaleString('zh-CN')
}

function formatNumber(num) {
  if (!num) return '0'
  return num.toLocaleString()
}

// 判断 API Key 是否过期
function isExpired(row) {
  if (!row.expires_at) return false
  return new Date(row.expires_at) < new Date()
}

// 获取状态显示标签
function getStatusLabel(row) {
  if (row.status === 'disabled') return '禁用'
  if (isExpired(row)) return '已过期'
  return '正常'
}

// 获取状态标签类型
function getStatusType(row) {
  if (row.status === 'disabled') return 'badge-danger'
  if (isExpired(row)) return 'badge-warning'
  return 'badge-success'
}

async function fetchAPIKeys() {
  loading.value = true
  try {
    const res = await api.adminGetAllAPIKeys({ page: pagination.page, page_size: pagination.pageSize })
    apiKeys.value = res.data?.items || []
    pagination.total = res.data?.total || 0
  } catch (e) {
    ElMessage.error('获取 API Key 列表失败')
  } finally {
    loading.value = false
  }
}

async function handleToggle(row) {
  try {
    await api.adminToggleAPIKey(row.id)
    ElMessage.success('状态已更新')
    fetchAPIKeys()
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

function confirmDelete(key) {
  deleteTarget.value = key
  deleteDialogVisible.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.adminDeleteAPIKey(deleteTarget.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    deleteTarget.value = null
    fetchAPIKeys()
  } catch (e) {
    ElMessage.error('删除失败')
  } finally {
    deleting.value = false
  }
}

// 查看日志
function viewLogs(row) {
  currentKey.value = row
  logPagination.page = 1
  logDialogVisible.value = true
  fetchLogs()
}

async function fetchLogs() {
  if (!currentKey.value) return
  logLoading.value = true
  logs.value = []
  try {
    const res = await api.adminGetAPIKeyLogs(currentKey.value.id, {
      page: logPagination.page,
      page_size: logPagination.pageSize
    })
    logs.value = res.data?.items || []
    logPagination.total = res.data?.total || 0
  } catch (e) {
    ElMessage.error('获取日志失败')
  } finally {
    logLoading.value = false
  }
}

onMounted(() => {
  fetchAPIKeys()
})
</script>

<style scoped>
.apikeys-page {
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

.data-table.compact {
  font-size: var(--apple-text-xs);
}

.data-table th,
.data-table td {
  padding: var(--apple-spacing-md);
  text-align: left;
  border-bottom: 1px solid var(--apple-separator);
}

.data-table.compact th,
.data-table.compact td {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
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

.col-id { width: 60px; }
.col-name { min-width: 100px; }
.col-user { min-width: 120px; }
.col-prefix { width: 120px; }
.col-status { width: 80px; }
.col-rate { width: 80px; }
.col-requests { width: 80px; }
.col-cost { width: 90px; }
.col-time { width: 150px; }
.col-actions { width: 120px; }

.key-name {
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

/* 用户信息 */
.user-info {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.user-avatar {
  width: 28px;
  height: 28px;
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  color: white;
  border-radius: var(--apple-radius-xs);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--apple-font-semibold);
  font-size: var(--apple-text-xs);
}

.key-prefix {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
  background: var(--apple-bg-tertiary);
  padding: 2px 8px;
  border-radius: var(--apple-radius-xs);
}

.model-name {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-primary);
  background: var(--apple-bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--apple-radius-xs);
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

.rate-info {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.cost {
  font-weight: var(--apple-font-medium);
  color: var(--apple-blue);
}

/* 操作按钮组 */
.action-group {
  display: flex;
  gap: var(--apple-spacing-xxs);
}

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

.action-btn.success:hover { background: var(--apple-green); }
.action-btn.warning:hover { background: var(--apple-orange); }
.action-btn.danger:hover { background: var(--apple-red); }

/* 分页 */
.pagination-wrap, .log-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
}

.log-pagination {
  padding: var(--apple-spacing-md) 0 0 0;
  border-top: none;
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

.empty-state.small svg {
  width: 32px;
  height: 32px;
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

.btn-outline {
  background: transparent;
  color: var(--apple-blue);
  border: 1px solid var(--apple-blue);
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-blue-light);
}

.btn-loading {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
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
.modal.modal-lg { max-width: 700px; }
.modal.modal-xl { max-width: 1000px; }

.modal-header {
  padding: var(--apple-spacing-xl);
  border-bottom: 1px solid var(--apple-separator);
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
}

.modal-header.danger {
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--apple-spacing-md);
}

.modal-title-group {
  flex: 1;
}

.modal-title-group h2 {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-sm) 0;
}

.modal-badges {
  display: flex;
  gap: var(--apple-spacing-sm);
  flex-wrap: wrap;
}

.danger-icon {
  width: 56px;
  height: 56px;
  background: var(--apple-red-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.danger-icon svg {
  width: 28px;
  height: 28px;
  color: var(--apple-red);
}

.modal-header h2 {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.modal-close svg {
  width: 18px;
  height: 18px;
}

.modal-close:hover {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
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

/* 日志表格容器 */
.log-table-container {
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
  overflow: hidden;
  max-height: 400px;
  overflow-y: auto;
}

/* 删除确认 */
.delete-message {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  text-align: center;
  margin: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }

  .modal-badges {
    flex-direction: column;
  }
}
</style>
