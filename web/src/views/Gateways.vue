<template>
  <div class="gateways-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <svg class="title-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="2" y1="12" x2="22" y2="12"/>
            <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/>
          </svg>
          网关管理
        </h1>
        <p class="page-desc">管理 xyrt 授权专用的网关配置，用于替换 chatgpt.com 域名</p>
      </div>
      <button class="btn btn-primary" @click="showAddModal = true">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        添加网关
      </button>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <div class="search-input-wrapper">
        <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          class="search-input"
          placeholder="搜索网关名称或地址..."
          @input="debouncedSearch"
        />
      </div>
    </div>

    <!-- 网关列表 -->
    <div class="gateway-list">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <span>加载中...</span>
      </div>

      <div v-else-if="gateways.length === 0" class="empty-state">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"/>
          <line x1="2" y1="12" x2="22" y2="12"/>
          <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/>
        </svg>
        <h3>暂无网关配置</h3>
        <p>点击上方"添加网关"按钮创建第一个网关</p>
      </div>

      <div v-else class="gateway-grid">
        <div
          v-for="gateway in gateways"
          :key="gateway.id"
          class="gateway-card"
          :class="{ disabled: !gateway.enabled, default: gateway.is_default }"
        >
          <div class="card-header">
            <div class="gateway-info">
              <h3 class="gateway-name">
                {{ gateway.name }}
                <span v-if="gateway.is_default" class="default-badge">默认</span>
              </h3>
              <p class="gateway-url">{{ gateway.url }}</p>
            </div>
            <div class="gateway-status">
              <span
                class="status-badge"
                :class="gateway.test_status || 'pending'"
              >
                {{ getStatusText(gateway.test_status) }}
              </span>
              <span v-if="gateway.test_latency > 0" class="latency">
                {{ gateway.test_latency }}ms
              </span>
            </div>
          </div>

          <div class="card-body">
            <div class="info-row">
              <span class="label">Codex URL:</span>
              <code class="value">{{ gateway.url }}/backend-api/codex</code>
            </div>
            <div v-if="gateway.remark" class="info-row">
              <span class="label">备注:</span>
              <span class="value">{{ gateway.remark }}</span>
            </div>
            <div v-if="gateway.test_error" class="info-row error">
              <span class="label">错误:</span>
              <span class="value">{{ gateway.test_error }}</span>
            </div>
          </div>

          <div class="card-footer">
            <div class="toggle-wrapper">
              <label class="toggle">
                <input
                  type="checkbox"
                  :checked="gateway.enabled"
                  @change="toggleEnabled(gateway)"
                />
                <span class="toggle-slider"></span>
              </label>
              <span class="toggle-label">{{ gateway.enabled ? '已启用' : '已禁用' }}</span>
            </div>
            <div class="actions">
              <button class="btn-icon" title="测试连接" @click="testGateway(gateway)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 11-5.93-9.14"/>
                  <polyline points="22 4 12 14.01 9 11.01"/>
                </svg>
              </button>
              <button
                v-if="!gateway.is_default && gateway.enabled"
                class="btn-icon"
                title="设为默认"
                @click="setDefault(gateway)"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
                </svg>
              </button>
              <button class="btn-icon" title="编辑" @click="editGateway(gateway)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                  <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                </svg>
              </button>
              <button class="btn-icon danger" title="删除" @click="confirmDelete(gateway)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="3 6 5 6 21 6"/>
                  <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div v-if="total > pageSize" class="pagination">
      <button :disabled="page === 1" @click="page--; loadGateways()">上一页</button>
      <span>{{ page }} / {{ Math.ceil(total / pageSize) }}</span>
      <button :disabled="page >= Math.ceil(total / pageSize)" @click="page++; loadGateways()">下一页</button>
    </div>

    <!-- 添加/编辑弹窗 -->
    <div v-if="showAddModal || showEditModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ showEditModal ? '编辑网关' : '添加网关' }}</h2>
          <button class="close-btn" @click="closeModal">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label required">网关名称</label>
            <input v-model="form.name" type="text" class="form-input" placeholder="例如：主网关" />
          </div>
          <div class="form-group">
            <label class="form-label required">网关地址</label>
            <input v-model="form.url" type="text" class="form-input" placeholder="例如：https://gateway.example.com" />
            <p class="form-tip">网关地址用于替换 chatgpt.com，请求将发送到 网关/backend-api/codex</p>
          </div>
          <div class="form-group">
            <label class="form-label">备注</label>
            <textarea v-model="form.remark" class="form-textarea" placeholder="可选备注信息"></textarea>
          </div>
          <div class="form-group">
            <label class="toggle-inline">
              <input v-model="form.enabled" type="checkbox" />
              <span>启用网关</span>
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeModal">取消</button>
          <button class="btn btn-primary" :disabled="!form.name || !form.url" @click="saveGateway">
            {{ showEditModal ? '保存' : '添加' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h2>确认删除</h2>
          <button class="close-btn" @click="showDeleteModal = false">×</button>
        </div>
        <div class="modal-body">
          <p>确定要删除网关 <strong>{{ deleteTarget?.name }}</strong> 吗？</p>
          <p class="warning-text">此操作不可恢复。</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showDeleteModal = false">取消</button>
          <button class="btn btn-danger" @click="deleteGateway">确认删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

// 状态
const loading = ref(false)
const gateways = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const searchQuery = ref('')

// 弹窗状态
const showAddModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const deleteTarget = ref(null)
const editTarget = ref(null)

// 表单数据
const form = ref({
  name: '',
  url: '',
  remark: '',
  enabled: true
})

// 加载网关列表
async function loadGateways() {
  loading.value = true
  try {
    const res = await api.getGateways({
      page: page.value,
      page_size: pageSize.value,
      search: searchQuery.value
    })
    const payload = res.data || res
    gateways.value = payload.items || []
    total.value = payload.total || 0
  } catch (e) {
    console.error('加载网关失败:', e)
  } finally {
    loading.value = false
  }
}

// 搜索防抖
let searchTimeout
function debouncedSearch() {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    page.value = 1
    loadGateways()
  }, 300)
}

// 获取状态文本
function getStatusText(status) {
  const map = {
    success: '正常',
    failed: '异常',
    pending: '未测试'
  }
  return map[status] || '未测试'
}

// 切换启用状态
async function toggleEnabled(gateway) {
  try {
    await api.toggleGateway(gateway.id)
    loadGateways()
  } catch (e) {
    console.error('切换状态失败:', e)
  }
}

// 测试网关
async function testGateway(gateway) {
  try {
    await api.testGatewayById(gateway.id)
    loadGateways()
  } catch (e) {
    console.error('测试失败:', e)
  }
}

// 设为默认
async function setDefault(gateway) {
  try {
    await api.setDefaultGateway(gateway.id)
    loadGateways()
  } catch (e) {
    console.error('设置默认失败:', e)
  }
}

// 编辑网关
function editGateway(gateway) {
  editTarget.value = gateway
  form.value = {
    name: gateway.name,
    url: gateway.url,
    remark: gateway.remark || '',
    enabled: gateway.enabled
  }
  showEditModal.value = true
}

// 确认删除
function confirmDelete(gateway) {
  deleteTarget.value = gateway
  showDeleteModal.value = true
}

// 删除网关
async function deleteGateway() {
  try {
    await api.deleteGateway(deleteTarget.value.id)
    showDeleteModal.value = false
    deleteTarget.value = null
    loadGateways()
  } catch (e) {
    alert(e.message || '删除失败')
  }
}

// 保存网关
async function saveGateway() {
  try {
    if (showEditModal.value) {
      await api.updateGateway(editTarget.value.id, form.value)
    } else {
      await api.createGateway(form.value)
    }
    closeModal()
    loadGateways()
  } catch (e) {
    alert(e.message || '保存失败')
  }
}

// 关闭弹窗
function closeModal() {
  showAddModal.value = false
  showEditModal.value = false
  editTarget.value = null
  form.value = { name: '', url: '', remark: '', enabled: true }
}

onMounted(loadGateways)
</script>

<style scoped>
.gateways-view {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.title-icon {
  width: 28px;
  height: 28px;
  color: var(--primary-color, #007AFF);
}

.page-desc {
  margin: 8px 0 0;
  color: #666;
  font-size: 14px;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn svg {
  width: 18px;
  height: 18px;
}

.btn-primary {
  background: var(--primary-color, #007AFF);
  color: white;
}

.btn-primary:hover {
  background: #0056b3;
}

.btn-secondary {
  background: #f0f0f0;
  color: #333;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.search-bar {
  margin-bottom: 24px;
}

.search-input-wrapper {
  position: relative;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 18px;
  height: 18px;
  color: #999;
}

.search-input {
  width: 100%;
  padding: 10px 12px 10px 40px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 14px;
}

.gateway-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.gateway-card {
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s;
}

.gateway-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.gateway-card.disabled {
  opacity: 0.6;
}

.gateway-card.default {
  border-color: var(--primary-color, #007AFF);
}

.card-header {
  padding: 16px;
  background: #f8f9fa;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.gateway-name {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.default-badge {
  font-size: 11px;
  padding: 2px 6px;
  background: var(--primary-color, #007AFF);
  color: white;
  border-radius: 4px;
  font-weight: 500;
}

.gateway-url {
  font-size: 13px;
  color: #666;
  margin: 4px 0 0;
  word-break: break-all;
}

.gateway-status {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.status-badge {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.status-badge.success {
  background: #d4edda;
  color: #155724;
}

.status-badge.failed {
  background: #f8d7da;
  color: #721c24;
}

.status-badge.pending {
  background: #e2e3e5;
  color: #6c757d;
}

.latency {
  font-size: 12px;
  color: #666;
}

.card-body {
  padding: 16px;
}

.info-row {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 13px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-row .label {
  color: #666;
  white-space: nowrap;
}

.info-row .value {
  color: #333;
  word-break: break-all;
}

.info-row code {
  background: #f5f5f5;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
}

.info-row.error .value {
  color: #dc3545;
}

.card-footer {
  padding: 12px 16px;
  background: #f8f9fa;
  border-top: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toggle-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toggle {
  position: relative;
  width: 40px;
  height: 22px;
}

.toggle input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  inset: 0;
  background: #ccc;
  border-radius: 22px;
  transition: 0.3s;
}

.toggle-slider::before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 2px;
  bottom: 2px;
  background: white;
  border-radius: 50%;
  transition: 0.3s;
}

.toggle input:checked + .toggle-slider {
  background: var(--primary-color, #007AFF);
}

.toggle input:checked + .toggle-slider::before {
  transform: translateX(18px);
}

.toggle-label {
  font-size: 13px;
  color: #666;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-icon {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.btn-icon svg {
  width: 18px;
  height: 18px;
  color: #666;
}

.btn-icon:hover {
  background: #e9ecef;
}

.btn-icon:hover svg {
  color: #333;
}

.btn-icon.danger:hover {
  background: #fee2e2;
}

.btn-icon.danger:hover svg {
  color: #dc3545;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #666;
}

.loading-state .spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid var(--primary-color, #007AFF);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state svg {
  width: 64px;
  height: 64px;
  color: #ccc;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 18px;
  margin: 0 0 8px;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 24px;
}

.pagination button {
  padding: 8px 16px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: white;
  cursor: pointer;
}

.pagination button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 12px;
  width: 500px;
  max-width: 90vw;
  max-height: 90vh;
  overflow: auto;
}

.modal.modal-sm {
  width: 400px;
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  font-size: 18px;
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  font-size: 24px;
  cursor: pointer;
  color: #666;
}

.modal-body {
  padding: 20px;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid #e0e0e0;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 6px;
}

.form-label.required::after {
  content: ' *';
  color: #dc3545;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.form-textarea {
  min-height: 80px;
  resize: vertical;
}

.form-tip {
  font-size: 12px;
  color: #666;
  margin: 6px 0 0;
}

.toggle-inline {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
}

.warning-text {
  color: #dc3545;
  font-size: 13px;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .gateways-view {
    background: #1a1a1a;
    color: #e0e0e0;
  }

  .gateway-card {
    background: #2a2a2a;
    border-color: #3a3a3a;
  }

  .card-header,
  .card-footer {
    background: #222;
    border-color: #3a3a3a;
  }

  .search-input,
  .form-input,
  .form-textarea {
    background: #2a2a2a;
    border-color: #3a3a3a;
    color: #e0e0e0;
  }

  .modal {
    background: #2a2a2a;
  }

  .modal-header,
  .modal-footer {
    border-color: #3a3a3a;
  }
}
</style>
