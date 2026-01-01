<!--
 * 文件作用：代理管理页面 - Apple HIG 风格
 * 负责功能：
 *   - 代理列表和CRUD
 *   - 代理连通性测试
 *   - 设置默认代理
 *   - 支持URL快速导入
 * 重要程度：⭐⭐⭐ 一般（代理配置）
-->
<template>
  <div class="proxies-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">代理管理</h1>
        <p class="page-subtitle">管理代理服务器配置，用于账户连接</p>
      </div>
      <div class="header-actions">
        <button class="btn btn-outline" @click="loadProxies" :disabled="loading">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: loading }">
            <polyline points="23,4 23,10 17,10"/>
            <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
          </svg>
          刷新
        </button>
        <button class="btn btn-primary" @click="handleAdd">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加代理
        </button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="filter-bar">
      <div class="search-box">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input
          v-model="keyword"
          placeholder="搜索代理名称、地址..."
          @input="handleSearch"
        />
      </div>
    </div>

    <!-- 代理列表 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>代理名称</th>
              <th>类型</th>
              <th>地址</th>
              <th class="col-center">连接状态</th>
              <th class="col-center">启用</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="proxy in proxies" :key="proxy.id">
              <td>
                <div class="proxy-cell">
                  <div class="proxy-icon" :class="proxy.type">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                    </svg>
                  </div>
                  <div class="proxy-info">
                    <div class="proxy-name-row">
                      <span class="proxy-name">{{ proxy.name }}</span>
                      <span v-if="proxy.is_default" class="badge badge-warning">
                        <svg viewBox="0 0 24 24" fill="currentColor" width="10" height="10">
                          <polygon points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"/>
                        </svg>
                        默认
                      </span>
                    </div>
                    <span v-if="proxy.remark" class="proxy-remark">{{ proxy.remark }}</span>
                  </div>
                </div>
              </td>
              <td>
                <span :class="['type-badge', proxy.type]">
                  {{ proxy.type.toUpperCase() }}
                </span>
              </td>
              <td>
                <code class="proxy-address">{{ proxy.host }}:{{ proxy.port }}</code>
              </td>
              <td class="col-center">
                <div class="status-cell">
                  <span v-if="proxy.testStatus === 'testing'" class="test-badge testing">
                    <span class="test-spinner"></span>
                    测试中
                  </span>
                  <span v-else-if="proxy.testStatus === 'success'" class="test-badge success">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="20,6 9,17 4,12"/>
                    </svg>
                    {{ proxy.latency }}ms
                  </span>
                  <span v-else-if="proxy.testStatus === 'failed'" class="test-badge failed">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <line x1="18" y1="6" x2="6" y2="18"/>
                      <line x1="6" y1="6" x2="18" y2="18"/>
                    </svg>
                    失败
                  </span>
                  <span v-else class="no-test">未测试</span>
                </div>
              </td>
              <td class="col-center">
                <label class="toggle-switch">
                  <input type="checkbox" v-model="proxy.enabled" @change="handleToggleEnabled(proxy)" />
                  <span class="toggle-slider"></span>
                </label>
              </td>
              <td>
                <div class="action-group">
                  <button
                    class="action-btn warning"
                    @click="handleSetDefault(proxy)"
                    :disabled="proxy.is_default || !proxy.enabled"
                    title="设为默认"
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polygon points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"/>
                    </svg>
                  </button>
                  <button
                    class="action-btn success"
                    @click="handleTestRow(proxy)"
                    :disabled="proxy.testStatus === 'testing'"
                    title="测试连接"
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M22 11.08V12a10 10 0 11-5.93-9.14"/>
                      <polyline points="22,4 12,14.01 9,11.01"/>
                    </svg>
                  </button>
                  <button class="action-btn" @click="handleEdit(proxy)" title="编辑">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                  </button>
                  <button class="action-btn danger" @click="confirmDelete(proxy)" title="删除">
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

        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <div v-if="!loading && proxies.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
          </svg>
          <span>暂无代理配置</span>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="pagination.total > 0" class="table-footer">
        <div class="pagination-wrap">
          <div class="pagination-info">共 {{ pagination.total }} 条</div>
          <div class="pagination-controls">
            <select v-model="pagination.pageSize" @change="loadProxies" class="page-size-select">
              <option :value="10">10 条/页</option>
              <option :value="20">20 条/页</option>
              <option :value="50">50 条/页</option>
            </select>
            <div class="page-btns">
              <button class="page-btn" :disabled="pagination.page <= 1" @click="pagination.page--; loadProxies()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="15,18 9,12 15,6"/>
                </svg>
              </button>
              <span class="page-current">{{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) || 1 }}</span>
              <button class="page-btn" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)" @click="pagination.page++; loadProxies()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9,18 15,12 9,6"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑弹窗 -->
    <Teleport to="body">
      <div v-if="dialogVisible" class="modal-overlay" @click.self="dialogVisible = false">
        <div class="modal">
          <div class="modal-header">
            <h2>{{ isEdit ? '编辑代理' : '添加代理' }}</h2>
            <button class="modal-close" @click="dialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <!-- 快速导入 -->
            <div class="form-group">
              <label class="form-label">快速导入</label>
              <div class="import-box">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/>
                  <polyline points="17,8 12,3 7,8"/>
                  <line x1="12" y1="3" x2="12" y2="15"/>
                </svg>
                <input
                  v-model="proxyUrl"
                  type="text"
                  class="form-input"
                  placeholder="粘贴代理URL，如: http://user:pass@host:port"
                  @input="parseProxyUrl"
                />
              </div>
              <p class="form-tip">支持格式: http://, https://, socks5:// 开头的代理地址</p>
            </div>

            <div class="form-divider">
              <span>代理配置</span>
            </div>

            <div class="form-group">
              <label class="form-label">代理名称 <span class="required">*</span></label>
              <input v-model="form.name" type="text" class="form-input" placeholder="为代理设置一个易识别的名称" />
            </div>

            <div class="form-group">
              <label class="form-label">代理类型 <span class="required">*</span></label>
              <div class="type-selector">
                <label :class="['type-option', { active: form.type === 'http' }]">
                  <input type="radio" v-model="form.type" value="http" />
                  <span>HTTP</span>
                </label>
                <label :class="['type-option', { active: form.type === 'https' }]">
                  <input type="radio" v-model="form.type" value="https" />
                  <span>HTTPS</span>
                </label>
                <label :class="['type-option', { active: form.type === 'socks5' }]">
                  <input type="radio" v-model="form.type" value="socks5" />
                  <span>SOCKS5</span>
                </label>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group flex-2">
                <label class="form-label">主机地址 <span class="required">*</span></label>
                <input v-model="form.host" type="text" class="form-input" placeholder="例如: 127.0.0.1 或 proxy.example.com" />
              </div>
              <div class="form-group flex-1">
                <label class="form-label">端口 <span class="required">*</span></label>
                <input v-model.number="form.port" type="number" class="form-input" min="1" max="65535" />
              </div>
            </div>

            <div class="form-divider">
              <span>认证配置（可选）</span>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">用户名</label>
                <input v-model="form.username" type="text" class="form-input" placeholder="代理认证用户名" />
              </div>
              <div class="form-group">
                <label class="form-label">密码</label>
                <input v-model="form.password" type="password" class="form-input" placeholder="代理认证密码" />
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">备注</label>
              <textarea v-model="form.remark" class="form-textarea" rows="2" placeholder="备注信息（可选）"></textarea>
            </div>

            <div class="form-group">
              <label class="form-label">启用</label>
              <label class="toggle-switch large">
                <input type="checkbox" v-model="form.enabled" />
                <span class="toggle-slider"></span>
              </label>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="dialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="submitting" @click="handleSubmit">
              <span v-if="submitting" class="btn-loading"></span>
              {{ submitting ? '保存中...' : (isEdit ? '保存修改' : '添加代理') }}
            </button>
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
            <p class="delete-message">确定要删除代理 "{{ deleteTarget?.name }}" 吗？</p>
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
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const proxies = ref([])
const keyword = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const proxyUrl = ref('')

// 删除相关
const deleteDialogVisible = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const defaultForm = {
  name: '',
  type: 'http',
  host: '',
  port: 7890,
  username: '',
  password: '',
  remark: '',
  enabled: true
}

const form = reactive({ ...defaultForm })

// 加载代理列表
async function loadProxies() {
  loading.value = true
  try {
    const res = await api.getProxyConfigs({
      page: pagination.page,
      page_size: pagination.pageSize,
      keyword: keyword.value
    })
    proxies.value = (res.items || []).map(p => ({
      ...p,
      testStatus: p.test_status || null,
      latency: p.test_latency || null
    }))
    pagination.total = res.total || 0
  } catch (e) {
    console.error('Failed to load proxies:', e)
    ElMessage.error('加载代理列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
let searchTimer = null
function handleSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.page = 1
    loadProxies()
  }, 300)
}

onUnmounted(() => {
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
})

// 解析代理URL
function parseProxyUrl() {
  const url = proxyUrl.value.trim()
  if (!url) return

  try {
    const regex = /^(https?|socks5):\/\/(?:([^:@]+):([^@]+)@)?([^:\/]+):(\d+)\/?$/i
    const match = url.match(regex)

    if (match) {
      const [, protocol, username, password, host, port] = match
      form.type = protocol.toLowerCase()
      form.host = host
      form.port = parseInt(port, 10)
      form.username = username || ''
      form.password = password || ''

      if (!form.name) {
        form.name = `${protocol.toUpperCase()} - ${host}:${port}`
      }

      ElMessage.success('代理地址解析成功')
    } else {
      const simpleMatch = url.match(/^([^:\/]+):(\d+)$/)
      if (simpleMatch) {
        form.host = simpleMatch[1]
        form.port = parseInt(simpleMatch[2], 10)
        if (!form.name) {
          form.name = `HTTP - ${form.host}:${form.port}`
        }
        ElMessage.success('代理地址解析成功')
      }
    }
  } catch (e) {
    console.error('解析代理URL失败:', e)
  }
}

// 测试代理连通性
async function testProxy(proxyData) {
  try {
    const res = await api.testProxyConnectivity({
      id: proxyData.id || 0,
      type: proxyData.type,
      host: proxyData.host,
      port: proxyData.port,
      username: proxyData.username || '',
      password: proxyData.password || ''
    })
    return res
  } catch (e) {
    return { success: false, error: e.message || '测试请求失败' }
  }
}

// 测试列表中的代理
async function handleTestRow(row) {
  row.testStatus = 'testing'
  row.latency = null

  const res = await testProxy(row)

  if (res.success) {
    row.testStatus = 'success'
    row.latency = res.latency
    ElMessage.success(`代理连接成功，延迟: ${res.latency}ms`)
  } else {
    row.testStatus = 'failed'
    ElMessage.error(res.error || '代理连接失败')
  }
}

// 添加
function handleAdd() {
  Object.assign(form, { ...defaultForm })
  proxyUrl.value = ''
  isEdit.value = false
  dialogVisible.value = true
}

// 编辑
function handleEdit(row) {
  Object.assign(form, { ...row })
  proxyUrl.value = ''
  isEdit.value = true
  dialogVisible.value = true
}

// 切换启用状态
async function handleToggleEnabled(row) {
  try {
    await api.toggleProxyConfig(row.id)
    ElMessage.success('更新成功')
    if (!row.enabled && row.is_default) {
      row.is_default = false
    }
  } catch (e) {
    row.enabled = !row.enabled
    ElMessage.error('更新失败')
  }
}

// 设置默认代理
async function handleSetDefault(row) {
  try {
    await api.setDefaultProxyConfig(row.id)
    proxies.value.forEach(p => {
      p.is_default = p.id === row.id
    })
    ElMessage.success(`已将 "${row.name}" 设为默认代理`)
  } catch (e) {
    ElMessage.error(e.message || '设置失败')
  }
}

// 删除确认
function confirmDelete(proxy) {
  deleteTarget.value = proxy
  deleteDialogVisible.value = true
}

// 删除
async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.deleteProxyConfig(deleteTarget.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    deleteTarget.value = null
    loadProxies()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  } finally {
    deleting.value = false
  }
}

// 提交表单
async function handleSubmit() {
  if (!form.name || !form.host || !form.port) {
    ElMessage.warning('请填写必填项')
    return
  }

  submitting.value = true

  try {
    ElMessage.info('正在测试代理连通性...')
    const testRes = await testProxy(form)

    if (!testRes.success) {
      ElMessage.warning(`代理连接测试失败: ${testRes.error}，仍将保存配置`)
    } else {
      ElMessage.success(`代理连接测试成功，延迟: ${testRes.latency}ms`)
    }

    if (isEdit.value) {
      await api.updateProxyConfig(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await api.createProxyConfig(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadProxies()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadProxies()
})
</script>

<style scoped>
.proxies-page {
  max-width: 1200px;
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

.header-actions {
  display: flex;
  gap: var(--apple-spacing-sm);
}

/* 筛选栏 */
.filter-bar {
  margin-bottom: var(--apple-spacing-lg);
}

.search-box {
  display: flex;
  align-items: center;
  position: relative;
  width: 300px;
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
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  padding-left: 36px;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
}

.search-box input:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
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

.col-center { text-align: center; }

/* 代理单元格 */
.proxy-cell {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.proxy-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.proxy-icon.https {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.proxy-icon.socks5 {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.proxy-icon svg {
  width: 20px;
  height: 20px;
  color: white;
}

.proxy-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.proxy-name-row {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.proxy-name {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.proxy-remark {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.proxy-address {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  background: var(--apple-fill-quaternary);
  padding: 2px 8px;
  border-radius: var(--apple-radius-xs);
}

/* 类型徽章 */
.type-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
  background: #eef2ff;
  color: #667eea;
}

.type-badge.https {
  background: #ecfdf5;
  color: #059669;
}

.type-badge.socks5 {
  background: #fef2f2;
  color: #dc2626;
}

/* 测试状态 */
.status-cell {
  display: flex;
  justify-content: center;
}

.test-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
  padding: 4px 10px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.test-badge svg {
  width: 12px;
  height: 12px;
}

.test-badge.testing {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-secondary);
}

.test-badge.success {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.test-badge.failed {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

.test-spinner {
  width: 12px;
  height: 12px;
  border: 2px solid var(--apple-fill-tertiary);
  border-top-color: var(--apple-text-secondary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.no-test {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-quaternary);
}

/* 徽章 */
.badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.badge-warning {
  background: var(--apple-orange-light);
  color: var(--apple-orange);
}

.badge-warning svg {
  width: 10px;
  height: 10px;
}

/* 开关 */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 36px;
  height: 20px;
}

.toggle-switch.large {
  margin-top: var(--apple-spacing-xs);
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--apple-fill-tertiary);
  transition: 0.3s;
  border-radius: 20px;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
}

.toggle-switch input:checked + .toggle-slider {
  background-color: var(--apple-green);
}

.toggle-switch input:checked + .toggle-slider:before {
  transform: translateX(16px);
}

/* 操作按钮 */
.action-group {
  display: flex;
  gap: var(--apple-spacing-xxs);
}

.action-btn {
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

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn:hover:not(:disabled) { background: var(--apple-blue); color: white; }
.action-btn.success:hover:not(:disabled) { background: var(--apple-green); }
.action-btn.warning:hover:not(:disabled) { background: var(--apple-orange); }
.action-btn.danger:hover:not(:disabled) { background: var(--apple-red); }

.action-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* 分页 */
.table-footer {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
}

.pagination-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-info {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
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

.page-btn svg {
  width: 14px;
  height: 14px;
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

.spinning {
  animation: spin 1s linear infinite;
}

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
  max-width: 520px;
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

/* 表单 */
.form-group {
  margin-bottom: var(--apple-spacing-lg);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-md);
}

.form-row .flex-1 { grid-column: span 1; }
.form-row .flex-2 { grid-column: span 1; }

.form-label {
  display: block;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.form-label .required {
  color: var(--apple-red);
}

.form-input, .form-textarea {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-textarea {
  resize: vertical;
  min-height: 60px;
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.form-tip {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: var(--apple-spacing-xxs);
}

.form-divider {
  display: flex;
  align-items: center;
  margin: var(--apple-spacing-lg) 0;
  color: var(--apple-text-tertiary);
  font-size: var(--apple-text-xs);
}

.form-divider::before,
.form-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--apple-separator);
}

.form-divider span {
  padding: 0 var(--apple-spacing-sm);
}

/* 快速导入 */
.import-box {
  position: relative;
}

.import-box svg {
  position: absolute;
  left: var(--apple-spacing-sm);
  top: 50%;
  transform: translateY(-50%);
  width: 16px;
  height: 16px;
  color: var(--apple-text-tertiary);
}

.import-box .form-input {
  padding-left: 36px;
}

/* 类型选择器 */
.type-selector {
  display: flex;
  gap: var(--apple-spacing-xs);
}

.type-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  cursor: pointer;
  transition: all var(--apple-duration-fast);
}

.type-option input {
  display: none;
}

.type-option:hover {
  border-color: var(--apple-blue);
  color: var(--apple-blue);
}

.type-option.active {
  background: var(--apple-blue);
  border-color: var(--apple-blue);
  color: white;
}

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

  .search-box {
    width: 100%;
  }

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
