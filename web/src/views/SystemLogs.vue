<!--
 * 文件作用：系统日志页面 - Apple HIG 风格
 * 负责功能：
 *   - 应用日志和服务器日志切换
 *   - 日志文件列表和筛选
 *   - 结构化日志解析展示
 *   - 实时日志追踪
 * 重要程度：⭐⭐ 辅助（运维日志）
-->
<template>
  <div class="system-logs-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">系统日志</h1>
        <p class="page-subtitle">查看应用日志和服务器日志</p>
      </div>
      <button class="btn btn-outline" @click="loadFiles" :disabled="loadingFiles">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: loadingFiles }">
          <polyline points="23,4 23,10 17,10"/>
          <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
        </svg>
        刷新
      </button>
    </div>

    <!-- Tab 切换 -->
    <div class="source-tabs">
      <button :class="['source-tab', { active: logSource === 'app' }]" @click="logSource = 'app'; handleSourceChange()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
          <polyline points="14,2 14,8 20,8"/>
        </svg>
        应用日志
      </button>
      <button :class="['source-tab', { active: logSource === 'server' }]" @click="logSource = 'server'; handleSourceChange()">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
          <line x1="8" y1="21" x2="16" y2="21"/>
          <line x1="12" y1="17" x2="12" y2="21"/>
        </svg>
        服务器日志
      </button>
    </div>

    <!-- 主内容区 -->
    <div class="logs-layout">
      <!-- 左侧：文件列表 -->
      <div class="file-panel">
        <div class="panel-header">
          <span class="panel-title">{{ logSource === 'app' ? '应用日志文件' : '服务器日志文件' }}</span>
          <span class="file-count">{{ files.length }} 个</span>
        </div>

        <!-- 筛选栏 -->
        <div class="filter-bar">
          <select v-model="filterCategory" @change="loadFiles" class="filter-select">
            <option value="">全部分类</option>
            <option v-for="cat in categories" :key="cat.name" :value="cat.name">{{ cat.label }} ({{ cat.count }})</option>
          </select>
          <input v-if="logSource === 'app'" type="date" v-model="filterDate" @change="loadFiles" class="filter-date" />
        </div>

        <!-- 文件列表 -->
        <div class="file-list" v-if="!loadingFiles">
          <div
            v-for="file in files"
            :key="file.name"
            :class="['file-item', { selected: selectedFile === file.name }]"
            @click="selectFile(file)"
          >
            <div class="file-info">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
              </svg>
              <span class="file-name">{{ file.name }}</span>
            </div>
            <div class="file-meta">
              <span class="file-size">{{ file.size_human }}</span>
              <div class="file-actions">
                <button class="file-action-btn" @click.stop="handleTail(file)" title="实时">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polygon points="5,3 19,12 5,21"/>
                  </svg>
                </button>
                <button class="file-action-btn" @click.stop="handleDownload(file)" title="下载">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                    <polyline points="7,10 12,15 17,10"/>
                    <line x1="12" y1="15" x2="12" y2="3"/>
                  </svg>
                </button>
                <button v-if="logSource === 'app'" class="file-action-btn danger" @click.stop="confirmDelete(file)" title="删除">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3,6 5,6 21,6"/>
                    <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
          <div v-if="files.length === 0" class="empty-files">暂无日志文件</div>
        </div>
        <div v-else class="loading-files">
          <div class="loading-spinner small"></div>
          <span>加载中...</span>
        </div>
      </div>

      <!-- 右侧：日志内容 -->
      <div class="content-panel">
        <div class="panel-header">
          <div class="content-info" v-if="selectedFile">
            <span class="selected-file">{{ selectedFile }}</span>
            <span v-if="logInfo.size_human" class="info-badge">{{ logInfo.size_human }}</span>
            <span v-if="logInfo.total_lines" class="info-badge success">{{ logInfo.total_lines }} 行</span>
          </div>
          <span v-else class="no-file">请选择日志文件</span>

          <div class="content-tools" v-if="selectedFile">
            <div class="search-box small">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="11" cy="11" r="8"/>
                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <input v-model="searchKeyword" placeholder="搜索..." @keyup.enter="loadLogContent" />
            </div>
            <label class="checkbox-label">
              <input type="checkbox" v-model="reverseOrder" @change="loadLogContent" />
              <span>倒序</span>
            </label>
            <label class="checkbox-label">
              <input type="checkbox" v-model="defaultExpanded" />
              <span>展开</span>
            </label>
            <button class="btn btn-sm btn-secondary" @click="loadLogContent" :disabled="loadingContent">查询</button>
            <button class="btn btn-sm btn-success" @click="tailLog" :disabled="loadingTail">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="5,3 19,12 5,21"/>
              </svg>
              实时
            </button>
          </div>
        </div>

        <!-- 日志内容区 -->
        <div class="log-content" :class="{ loading: loadingContent || loadingTail }">
          <div v-if="!selectedFile" class="empty-content">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
              <polyline points="14,2 14,8 20,8"/>
            </svg>
            <span>请从左侧选择日志文件</span>
          </div>

          <div v-else-if="logEntries.length === 0 && logLines.length === 0 && !loadingContent" class="empty-content">
            <span>没有日志内容</span>
          </div>

          <!-- 结构化日志 -->
          <div v-else-if="logSource === 'app' && logEntries.length > 0" class="log-entries" ref="logPreRef">
            <div
              v-for="(entry, index) in logEntries"
              :key="index"
              :class="['log-entry', getLevelClass(entry.level), { expanded: isExpanded(index) }]"
              @click="toggleExpand(index)"
            >
              <template v-if="entry.is_json">
                <div class="log-main-line">
                  <span class="log-time">{{ formatTime(entry.timestamp) }}</span>
                  <span :class="['log-level', getLevelClass(entry.level)]">{{ entry.level }}</span>
                  <span class="log-module" v-if="entry.module">[{{ entry.module }}]</span>
                  <span class="log-message">{{ entry.message }}</span>
                  <span class="inline-fields" v-if="hasInlineFields(entry)">
                    <span v-if="entry.fields?.user_id" class="inline-field user">👤{{ entry.fields.user_id }}</span>
                    <span v-if="entry.fields?.model" class="inline-field model">🤖{{ entry.fields.model }}</span>
                    <span v-if="entry.fields?.status" :class="['inline-field', getStatusClass(entry.fields.status)]">{{ entry.fields.status }}</span>
                    <span v-if="entry.fields?.latency !== undefined" class="inline-field latency">{{ formatLatency(entry.fields.latency) }}</span>
                  </span>
                  <span class="expand-icon" v-if="hasFields(entry)">{{ isExpanded(index) ? '▼' : '▶' }}</span>
                </div>
                <div v-if="isExpanded(index) && hasFields(entry)" class="log-fields">
                  <div class="fields-grid">
                    <div v-for="(value, key) in getSortedFields(entry.fields)" :key="key" class="log-field">
                      <span class="field-key">{{ key }}:</span>
                      <span class="field-value">{{ formatFieldValue(key, value) }}</span>
                    </div>
                    <div v-if="entry.caller" class="log-field full">
                      <span class="field-key">caller:</span>
                      <span class="field-value caller">{{ entry.caller }}</span>
                    </div>
                  </div>
                </div>
              </template>
              <template v-else>
                <span class="log-raw">{{ entry.raw }}</span>
              </template>
            </div>
          </div>

          <!-- 原始文本 -->
          <pre v-else-if="logLines.length > 0" class="log-pre" ref="logPreRef">{{ logLines.join('\n') }}</pre>

          <div v-if="loadingContent || loadingTail" class="loading-overlay">
            <div class="loading-spinner"></div>
          </div>
        </div>

        <!-- 分页 -->
        <div v-if="logInfo.total_pages > 1" class="content-footer">
          <div class="pagination-info">共 {{ logInfo.total_lines }} 行</div>
          <div class="pagination-controls">
            <select v-model="pageSize" @change="handlePageSizeChange" class="page-size-select">
              <option :value="100">100 行</option>
              <option :value="200">200 行</option>
              <option :value="500">500 行</option>
            </select>
            <div class="page-btns">
              <button class="page-btn" :disabled="page <= 1" @click="page--; loadLogContent()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15,18 9,12 15,6"/></svg>
              </button>
              <span class="page-current">{{ page }} / {{ logInfo.total_pages }}</span>
              <button class="page-btn" :disabled="page >= logInfo.total_pages" @click="page++; loadLogContent()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9,18 15,12 9,6"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <Teleport to="body">
      <div v-if="deleteVisible" class="modal-overlay" @click.self="deleteVisible = false">
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
            <p class="delete-message">确定要删除日志文件 "{{ deleteTarget }}" 吗？</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="deleteVisible = false">取消</button>
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
import { ref, nextTick, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const logSource = ref('app')
const files = ref([])
const categories = ref([])
const loadingFiles = ref(false)
const filterCategory = ref('')
const filterDate = ref('')

const selectedFile = ref('')
const logLines = ref([])
const logEntries = ref([])
const logInfo = ref({})
const loadingContent = ref(false)
const loadingTail = ref(false)
const page = ref(1)
const pageSize = ref(200)
const searchKeyword = ref('')
const reverseOrder = ref(false)
const logPreRef = ref(null)
const expandedRows = ref([])
const defaultExpanded = ref(true)

const deleteVisible = ref(false)
const deleteTarget = ref('')
const deleting = ref(false)

function handleSourceChange() {
  files.value = []
  categories.value = []
  selectedFile.value = ''
  logLines.value = []
  logEntries.value = []
  logInfo.value = {}
  filterCategory.value = ''
  filterDate.value = ''
  page.value = 1
  searchKeyword.value = ''
  reverseOrder.value = false
  expandedRows.value = []
  loadFiles()
}

async function loadFiles() {
  loadingFiles.value = true
  try {
    const params = { source: logSource.value }
    if (filterCategory.value) params.category = filterCategory.value
    if (filterDate.value && logSource.value === 'app') params.date = filterDate.value
    const res = await api.getSystemLogFiles(params)
    files.value = res.data?.files || []
    categories.value = res.data?.categories || []
  } catch (e) {
    console.error('Failed to load files:', e)
  } finally {
    loadingFiles.value = false
  }
}

function selectFile(file) {
  selectedFile.value = file.name
  page.value = 1
  searchKeyword.value = ''
  expandedRows.value = []
  loadLogContent()
}

async function loadLogContent() {
  if (!selectedFile.value) return
  loadingContent.value = true
  expandedRows.value = []
  try {
    const res = await api.readSystemLog({
      file: selectedFile.value,
      source: logSource.value,
      page: page.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value,
      reverse: reverseOrder.value ? 'true' : ''
    })
    logLines.value = res.data?.lines || []
    logEntries.value = res.data?.entries || []
    logInfo.value = {
      size_human: res.data?.size_human,
      total_lines: res.data?.total_lines,
      total_pages: res.data?.total_pages
    }
    nextTick(() => { if (logPreRef.value) logPreRef.value.scrollTop = 0 })
  } catch (e) {
    ElMessage.error('加载日志失败')
  } finally {
    loadingContent.value = false
  }
}

async function tailLog() {
  if (!selectedFile.value) return
  loadingTail.value = true
  expandedRows.value = []
  try {
    const res = await api.tailSystemLog({ file: selectedFile.value, source: logSource.value, lines: 200 })
    logLines.value = res.data?.lines || []
    logEntries.value = res.data?.entries || []
    logInfo.value = { size_human: res.data?.size_human, total_lines: res.data?.count, total_pages: 1 }
    page.value = 1
    nextTick(() => { if (logPreRef.value) logPreRef.value.scrollTop = logPreRef.value.scrollHeight })
    ElMessage.success('已加载最新 ' + (res.data?.count || 0) + ' 行')
  } catch (e) {
    ElMessage.error('加载实时日志失败')
  } finally {
    loadingTail.value = false
  }
}

function handlePageSizeChange() {
  page.value = 1
  loadLogContent()
}

function handleTail(file) {
  selectedFile.value = file.name
  tailLog()
}

function handleDownload(file) {
  window.open(api.downloadSystemLog(file.name, logSource.value), '_blank')
}

function confirmDelete(file) {
  deleteTarget.value = file.name
  deleteVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  try {
    await api.deleteSystemLog(deleteTarget.value, logSource.value)
    ElMessage.success('删除成功')
    if (selectedFile.value === deleteTarget.value) {
      selectedFile.value = ''
      logLines.value = []
      logEntries.value = []
      logInfo.value = {}
    }
    deleteVisible.value = false
    loadFiles()
  } catch (e) {
    ElMessage.error('删除失败')
  } finally {
    deleting.value = false
  }
}

function getLevelClass(level) {
  if (!level) return ''
  const l = level.toUpperCase()
  if (l === 'ERROR' || l === 'FATAL') return 'level-error'
  if (l === 'WARN' || l === 'WARNING') return 'level-warn'
  if (l === 'DEBUG') return 'level-debug'
  return 'level-info'
}

function formatTime(timestamp) {
  if (!timestamp) return ''
  try {
    return new Date(timestamp).toLocaleString('zh-CN', { hour12: false })
  } catch { return timestamp }
}

function hasFields(entry) {
  return (entry.fields && Object.keys(entry.fields).length > 0) || entry.caller
}

function hasInlineFields(entry) {
  if (!entry.fields) return false
  return ['user_id', 'model', 'status', 'latency'].some(k => entry.fields[k] !== undefined)
}

function isExpanded(index) {
  return defaultExpanded.value ? !expandedRows.value.includes(index) : expandedRows.value.includes(index)
}

function toggleExpand(index) {
  const idx = expandedRows.value.indexOf(index)
  if (idx > -1) expandedRows.value.splice(idx, 1)
  else expandedRows.value.push(index)
}

function getStatusClass(status) {
  if (status >= 500) return 'status-error'
  if (status >= 400) return 'status-warn'
  if (status >= 200 && status < 300) return 'status-success'
  return ''
}

function formatLatency(latency) {
  if (latency === undefined) return ''
  return latency >= 1000 ? (latency / 1000).toFixed(2) + 's' : latency + 'ms'
}

function getSortedFields(fields) {
  if (!fields) return {}
  const priority = ['user_id', 'api_key_id', 'account_id', 'model', 'client_ip', 'method', 'path', 'status', 'latency', 'input_tokens', 'output_tokens', 'total_tokens', 'total_cost', 'error']
  const sorted = {}
  for (const k of priority) if (fields[k] !== undefined) sorted[k] = fields[k]
  for (const k of Object.keys(fields)) if (sorted[k] === undefined) sorted[k] = fields[k]
  return sorted
}

function formatFieldValue(key, value) {
  if (value === null || value === undefined) return '-'
  if (['input_cost', 'output_cost', 'total_cost'].includes(key)) return '$' + Number(value).toFixed(6)
  if (['latency', 'exec_duration'].includes(key)) return typeof value === 'number' ? (value >= 1000 ? (value / 1000).toFixed(2) + 's' : value + 'ms') : value
  if (['input_tokens', 'output_tokens', 'total_tokens'].includes(key)) return Number(value).toLocaleString()
  if (typeof value === 'object') return JSON.stringify(value)
  return String(value)
}

onMounted(() => { loadFiles() })
</script>

<style scoped>
.system-logs-page { max-width: 1600px; margin: 0 auto; }

/* 页面标题 */
.page-header { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: var(--apple-spacing-lg); }
.header-content { flex: 1; }
.page-title { font-size: var(--apple-text-3xl); font-weight: var(--apple-font-bold); color: var(--apple-text-primary); margin: 0 0 var(--apple-spacing-xs) 0; }
.page-subtitle { font-size: var(--apple-text-base); color: var(--apple-text-secondary); margin: 0; }

/* Tab 切换 */
.source-tabs { display: flex; gap: var(--apple-spacing-xs); margin-bottom: var(--apple-spacing-lg); }
.source-tab {
  display: flex; align-items: center; gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  font-size: var(--apple-text-sm); font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary); background: var(--apple-fill-quaternary);
  border-radius: var(--apple-radius-sm); transition: all var(--apple-duration-fast);
}
.source-tab svg { width: 16px; height: 16px; }
.source-tab:hover { background: var(--apple-fill-tertiary); color: var(--apple-text-primary); }
.source-tab.active { background: var(--apple-blue); color: white; }

/* 主布局 */
.logs-layout { display: grid; grid-template-columns: 320px 1fr; gap: var(--apple-spacing-lg); height: calc(100vh - 240px); }

/* 文件面板 */
.file-panel {
  background: var(--apple-bg-primary); border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card); display: flex; flex-direction: column; overflow: hidden;
}
.panel-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}
.panel-title { font-weight: var(--apple-font-semibold); color: var(--apple-text-primary); }
.file-count { font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); background: var(--apple-fill-tertiary); padding: 2px 8px; border-radius: var(--apple-radius-full); }

.filter-bar { display: flex; gap: var(--apple-spacing-xs); padding: var(--apple-spacing-sm) var(--apple-spacing-md); border-bottom: 1px solid var(--apple-separator); }
.filter-select, .filter-date {
  flex: 1; padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs); border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm); background: var(--apple-bg-primary);
}

.file-list { flex: 1; overflow-y: auto; padding: var(--apple-spacing-xs); }
.file-item {
  display: flex; justify-content: space-between; align-items: center;
  padding: var(--apple-spacing-sm); border-radius: var(--apple-radius-sm);
  cursor: pointer; transition: background var(--apple-duration-fast);
}
.file-item:hover { background: var(--apple-bg-secondary); }
.file-item.selected { background: var(--apple-blue-light); }
.file-info { display: flex; align-items: center; gap: var(--apple-spacing-xs); flex: 1; min-width: 0; }
.file-info svg { width: 16px; height: 16px; color: var(--apple-text-tertiary); flex-shrink: 0; }
.file-name { font-size: var(--apple-text-sm); color: var(--apple-text-primary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-meta { display: flex; align-items: center; gap: var(--apple-spacing-xs); }
.file-size { font-size: var(--apple-text-xs); color: var(--apple-text-tertiary); }
.file-actions { display: flex; gap: 2px; opacity: 0; transition: opacity var(--apple-duration-fast); }
.file-item:hover .file-actions { opacity: 1; }
.file-action-btn {
  width: 24px; height: 24px; border-radius: var(--apple-radius-xs);
  display: flex; align-items: center; justify-content: center;
  color: var(--apple-text-tertiary); transition: all var(--apple-duration-fast);
}
.file-action-btn svg { width: 12px; height: 12px; }
.file-action-btn:hover { background: var(--apple-blue); color: white; }
.file-action-btn.danger:hover { background: var(--apple-red); }

.empty-files, .loading-files { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: var(--apple-spacing-xl); color: var(--apple-text-tertiary); font-size: var(--apple-text-sm); }

/* 内容面板 */
.content-panel {
  background: var(--apple-bg-primary); border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card); display: flex; flex-direction: column; overflow: hidden;
}
.content-info { display: flex; align-items: center; gap: var(--apple-spacing-sm); }
.selected-file { font-weight: var(--apple-font-semibold); color: var(--apple-blue); }
.info-badge { font-size: var(--apple-text-xs); padding: 2px 8px; background: var(--apple-fill-tertiary); border-radius: var(--apple-radius-full); color: var(--apple-text-secondary); }
.info-badge.success { background: var(--apple-green-light); color: var(--apple-green); }
.no-file { color: var(--apple-text-tertiary); }
.content-tools { display: flex; align-items: center; gap: var(--apple-spacing-sm); }
.search-box.small { width: 120px; }
.search-box.small input { padding: var(--apple-spacing-xxs) var(--apple-spacing-sm); padding-left: 28px; font-size: var(--apple-text-xs); }
.search-box.small svg { width: 12px; height: 12px; left: var(--apple-spacing-xs); }
.checkbox-label { display: flex; align-items: center; gap: 4px; font-size: var(--apple-text-xs); color: var(--apple-text-secondary); cursor: pointer; }
.checkbox-label input { width: 14px; height: 14px; }

/* 日志内容 */
.log-content { flex: 1; overflow: hidden; position: relative; min-height: 300px; }
.log-content.loading { pointer-events: none; }
.empty-content { display: flex; flex-direction: column; align-items: center; justify-content: center; height: 100%; color: var(--apple-text-tertiary); gap: var(--apple-spacing-md); }
.empty-content svg { width: 48px; height: 48px; opacity: 0.5; }
.loading-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; background: rgba(255,255,255,0.8); display: flex; align-items: center; justify-content: center; }

/* 结构化日志 - 暗色终端风格 */
.log-entries {
  background: #1e1e1e; height: 100%; overflow: auto; padding: var(--apple-spacing-sm);
  font-family: 'SF Mono', Consolas, Monaco, monospace; font-size: 12px;
}
.log-entry { padding: 4px 8px; border-radius: 3px; margin-bottom: 2px; cursor: pointer; line-height: 1.6; }
.log-entry:hover { background: rgba(255,255,255,0.08); }
.log-entry.expanded { background: rgba(255,255,255,0.03); }
.log-entry.level-error { background: rgba(244,135,113,0.1); border-left: 3px solid #f48771; }
.log-entry.level-warn { background: rgba(204,167,0,0.1); border-left: 3px solid #cca700; }
.log-main-line { display: flex; align-items: center; flex-wrap: wrap; gap: 4px; }
.log-time { color: #6a9955; margin-right: 8px; }
.log-level { font-weight: bold; padding: 1px 4px; border-radius: 2px; font-size: 11px; }
.log-level.level-error { background: #5a1d1d; color: #f48771; }
.log-level.level-warn { background: #5a4a1d; color: #cca700; }
.log-level.level-info { background: #1d3a5a; color: #4fc1ff; }
.log-level.level-debug { background: #2d2d2d; color: #888; }
.log-module { color: #c586c0; margin-right: 8px; }
.log-message { color: #d4d4d4; }
.log-raw { color: #d4d4d4; white-space: pre-wrap; word-break: break-all; }
.inline-fields { display: inline-flex; gap: 6px; margin-left: 10px; }
.inline-field { padding: 1px 6px; border-radius: 3px; font-size: 11px; background: rgba(255,255,255,0.1); }
.inline-field.user { background: rgba(79,193,255,0.2); color: #4fc1ff; }
.inline-field.model { background: rgba(220,220,170,0.2); color: #dcdcaa; }
.inline-field.latency { background: rgba(206,145,120,0.2); color: #ce9178; }
.inline-field.status-success { background: rgba(106,153,85,0.3); color: #6a9955; }
.inline-field.status-warn { background: rgba(204,167,0,0.3); color: #cca700; }
.inline-field.status-error { background: rgba(244,135,113,0.3); color: #f48771; }
.expand-icon { margin-left: auto; color: #666; font-size: 10px; }
.log-fields { margin-top: 8px; padding: 10px 12px; background: rgba(0,0,0,0.4); border-radius: 4px; border-left: 3px solid #4fc1ff; }
.fields-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); gap: 4px 20px; }
.log-field { display: flex; align-items: flex-start; padding: 2px 0; }
.log-field.full { grid-column: 1 / -1; }
.field-key { color: #9cdcfe; margin-right: 8px; min-width: 100px; }
.field-value { color: #ce9178; word-break: break-all; }
.field-value.caller { color: #888; font-size: 11px; }

/* 原始日志 */
.log-pre {
  margin: 0; padding: 15px; background: #1e1e1e; color: #d4d4d4; height: 100%; overflow: auto;
  font-family: 'SF Mono', Consolas, Monaco, monospace; font-size: 12px; line-height: 1.6; white-space: pre-wrap;
}

/* 分页 */
.content-footer { padding: var(--apple-spacing-sm) var(--apple-spacing-lg); border-top: 1px solid var(--apple-separator); display: flex; justify-content: space-between; align-items: center; }
.pagination-info { font-size: var(--apple-text-sm); color: var(--apple-text-secondary); }
.pagination-controls { display: flex; align-items: center; gap: var(--apple-spacing-sm); }
.page-size-select { padding: var(--apple-spacing-xxs) var(--apple-spacing-xs); font-size: var(--apple-text-xs); border: 1px solid var(--apple-separator-opaque); border-radius: var(--apple-radius-sm); }
.page-btns { display: flex; align-items: center; gap: var(--apple-spacing-xs); }
.page-btn { width: 24px; height: 24px; border-radius: var(--apple-radius-sm); display: flex; align-items: center; justify-content: center; background: var(--apple-fill-quaternary); color: var(--apple-text-secondary); }
.page-btn svg { width: 12px; height: 12px; }
.page-btn:hover:not(:disabled) { background: var(--apple-blue); color: white; }
.page-btn:disabled { opacity: 0.3; }
.page-current { font-size: var(--apple-text-xs); color: var(--apple-text-primary); min-width: 50px; text-align: center; }

/* 按钮 */
.btn { display: inline-flex; align-items: center; justify-content: center; gap: var(--apple-spacing-xs); padding: var(--apple-spacing-sm) var(--apple-spacing-lg); font-size: var(--apple-text-sm); font-weight: var(--apple-font-medium); border-radius: var(--apple-radius-sm); transition: all var(--apple-duration-fast); cursor: pointer; }
.btn svg { width: 14px; height: 14px; }
.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-sm { padding: var(--apple-spacing-xxs) var(--apple-spacing-sm); font-size: var(--apple-text-xs); }
.btn-secondary { background: var(--apple-fill-tertiary); color: var(--apple-text-primary); }
.btn-secondary:hover:not(:disabled) { background: var(--apple-fill-secondary); }
.btn-success { background: var(--apple-green); color: white; }
.btn-success:hover:not(:disabled) { background: #2db553; }
.btn-danger { background: var(--apple-red); color: white; }
.btn-danger:hover:not(:disabled) { background: #e6362d; }
.btn-outline { background: transparent; color: var(--apple-blue); border: 1px solid var(--apple-blue); }
.btn-outline:hover:not(:disabled) { background: var(--apple-blue-light); }
.btn-loading { width: 14px; height: 14px; border: 2px solid rgba(255,255,255,0.3); border-top-color: white; border-radius: 50%; animation: spin 1s linear infinite; }

.loading-spinner { width: 24px; height: 24px; border: 2px solid var(--apple-fill-tertiary); border-top-color: var(--apple-blue); border-radius: 50%; animation: spin 1s linear infinite; }
.loading-spinner.small { width: 16px; height: 16px; }
@keyframes spin { to { transform: rotate(360deg); } }
.spinning { animation: spin 1s linear infinite; }

/* 搜索框 */
.search-box { display: flex; align-items: center; position: relative; }
.search-box svg { position: absolute; left: var(--apple-spacing-sm); width: 14px; height: 14px; color: var(--apple-text-tertiary); pointer-events: none; }
.search-box input { width: 100%; padding: var(--apple-spacing-xs) var(--apple-spacing-sm); padding-left: 32px; font-size: var(--apple-text-sm); border: 1px solid var(--apple-separator-opaque); border-radius: var(--apple-radius-sm); background: var(--apple-bg-primary); }
.search-box input:focus { outline: none; border-color: var(--apple-blue); }

/* 模态框 */
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: var(--apple-z-modal); }
.modal { background: var(--apple-bg-primary); border-radius: var(--apple-radius-xl); box-shadow: var(--apple-shadow-modal); width: 100%; max-width: 400px; }
.modal.modal-sm { max-width: 360px; }
.modal-header { padding: var(--apple-spacing-xl); border-bottom: 1px solid var(--apple-separator); display: flex; align-items: center; justify-content: space-between; }
.modal-header.danger { flex-direction: column; text-align: center; gap: var(--apple-spacing-md); }
.danger-icon { width: 56px; height: 56px; background: var(--apple-red-light); border-radius: 50%; display: flex; align-items: center; justify-content: center; }
.danger-icon svg { width: 28px; height: 28px; color: var(--apple-red); }
.modal-header h2 { font-size: var(--apple-text-lg); font-weight: var(--apple-font-semibold); color: var(--apple-text-primary); margin: 0; }
.modal-body { padding: var(--apple-spacing-xl); }
.modal-footer { padding: var(--apple-spacing-lg) var(--apple-spacing-xl); border-top: 1px solid var(--apple-separator); display: flex; justify-content: flex-end; gap: var(--apple-spacing-sm); }
.delete-message { font-size: var(--apple-text-base); color: var(--apple-text-secondary); text-align: center; margin: 0; }

/* 响应式 */
@media (max-width: 1024px) {
  .logs-layout { grid-template-columns: 1fr; height: auto; }
  .file-panel { max-height: 300px; }
  .content-panel { min-height: 500px; }
}
</style>
