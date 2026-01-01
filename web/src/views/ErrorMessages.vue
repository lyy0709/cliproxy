<!--
 * 文件作用：错误配置页面 - Apple HIG 风格
 * 负责功能：
 *   - 错误消息CRUD和状态管理
 *   - 错误规则CRUD和优先级配置
 *   - 缓存刷新和批量操作
 *   - 自动账户禁用/限流规则
 * 重要程度：⭐⭐⭐ 一般（错误处理配置）
-->
<template>
  <div class="error-config-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">错误配置</h1>
        <p class="page-subtitle">管理错误消息和自动处理规则</p>
      </div>
    </div>

    <!-- Tab 切换 -->
    <div class="tab-container">
      <div class="tab-nav">
        <button :class="['tab-btn', { active: activeTab === 'messages' }]" @click="activeTab = 'messages'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="8" x2="12" y2="12"/>
            <line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          错误消息
        </button>
        <button :class="['tab-btn', { active: activeTab === 'rules' }]" @click="activeTab = 'rules'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="12 2 2 7 12 12 22 7 12 2"/>
            <polyline points="2 17 12 22 22 17"/>
            <polyline points="2 12 12 17 22 12"/>
          </svg>
          错误规则
        </button>
      </div>
    </div>

    <!-- 错误消息面板 -->
    <div v-show="activeTab === 'messages'" class="tab-panel">
      <!-- 操作栏 -->
      <div class="action-bar">
        <div class="action-group">
          <button class="btn btn-success" @click="enableAllMessages" :disabled="enablingAll">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: enablingAll }">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            全部启用
          </button>
          <button class="btn btn-warning" @click="disableAllMessages" :disabled="disablingAll">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: disablingAll }">
              <circle cx="12" cy="12" r="10"/>
              <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
            </svg>
            全部禁用
          </button>
        </div>
        <div class="action-group">
          <button class="btn btn-primary" @click="showCreateMessageDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            新增配置
          </button>
          <button class="btn btn-outline" @click="initMessageDefaults" :disabled="initializingMessages">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: initializingMessages }">
              <path d="M21 12a9 9 0 11-2-5.65"/>
              <path d="M21 3v6h-6"/>
            </svg>
            初始化默认
          </button>
          <button class="btn btn-outline" @click="refreshMessageCache" :disabled="refreshingMessages">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: refreshingMessages }">
              <polyline points="23,4 23,10 17,10"/>
              <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
            </svg>
            刷新缓存
          </button>
        </div>
      </div>

      <!-- 消息表格 -->
      <div class="table-card">
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th width="90">状态码</th>
                <th width="180">错误类型</th>
                <th>原始信息</th>
                <th>自定义消息</th>
                <th width="150">说明</th>
                <th width="80">状态</th>
                <th width="100">操作</th>
              </tr>
            </thead>
            <tbody v-if="!loadingMessages">
              <tr v-for="row in messages" :key="row.id">
                <td>
                  <span :class="['code-badge', getCodeClass(row.code)]">{{ row.code }}</span>
                </td>
                <td>
                  <span class="error-type">{{ row.error_type }}</span>
                </td>
                <td>
                  <span class="muted-text">{{ row.original_message || '-' }}</span>
                </td>
                <td>
                  <span class="custom-message">{{ row.custom_message }}</span>
                </td>
                <td>
                  <span class="description">{{ row.description || '-' }}</span>
                </td>
                <td>
                  <button :class="['toggle-switch', { active: row.enabled }]" @click="toggleMessageEnabled(row)">
                    <span class="toggle-thumb"></span>
                  </button>
                </td>
                <td>
                  <div class="action-btns">
                    <button class="action-btn edit" @click="showEditMessageDialog(row)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn delete" @click="confirmDeleteMessage(row)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="loadingMessages" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>

          <div v-if="!loadingMessages && messages.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <span>暂无错误消息配置</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误规则面板 -->
    <div v-show="activeTab === 'rules'" class="tab-panel">
      <!-- 规则说明 -->
      <div class="guide-card">
        <div class="guide-header">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3"/>
            <line x1="12" y1="17" x2="12.01" y2="17"/>
          </svg>
          <span>规则说明</span>
        </div>
        <div class="guide-content">
          <div class="guide-item">
            <span class="guide-badge danger">禁用账户</span>
            <span class="guide-text">账户被封号、认证失效时，自动禁用该账户</span>
          </div>
          <div class="guide-item">
            <span class="guide-badge warning">临时限流</span>
            <span class="guide-text">遇到限流或临时错误，切换其他账户重试（1小时后恢复）</span>
          </div>
          <div class="guide-item">
            <span class="guide-badge info">过载</span>
            <span class="guide-text">服务过载，临时切换账户</span>
          </div>
        </div>
      </div>

      <!-- 操作栏 -->
      <div class="action-bar">
        <div class="action-group">
          <button class="btn btn-success" @click="enableAllRules" :disabled="enablingAllRules">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: enablingAllRules }">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            全部启用
          </button>
          <button class="btn btn-warning" @click="disableAllRules" :disabled="disablingAllRules">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: disablingAllRules }">
              <circle cx="12" cy="12" r="10"/>
              <line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
            </svg>
            全部禁用
          </button>
        </div>
        <div class="action-group">
          <button class="btn btn-primary" @click="showCreateRuleDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            新增规则
          </button>
          <button class="btn btn-outline" @click="resetRulesToDefault" :disabled="resettingRules">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: resettingRules }">
              <path d="M21 12a9 9 0 11-2-5.65"/>
              <path d="M21 3v6h-6"/>
            </svg>
            重置为默认
          </button>
          <button class="btn btn-outline" @click="refreshRuleCache" :disabled="refreshingRules">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: refreshingRules }">
              <polyline points="23,4 23,10 17,10"/>
              <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
            </svg>
            刷新缓存
          </button>
        </div>
      </div>

      <!-- 规则表格 -->
      <div class="table-card">
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th width="120">处理方式</th>
                <th>匹配条件</th>
                <th width="200">说明</th>
                <th width="80">优先级</th>
                <th width="70">启用</th>
                <th width="100">操作</th>
              </tr>
            </thead>
            <tbody v-if="!loadingRules">
              <tr v-for="row in rules" :key="row.id">
                <td>
                  <span :class="['status-badge', getTargetStatusClass(row.target_status)]">
                    {{ getTargetStatusLabel(row.target_status) }}
                  </span>
                </td>
                <td>
                  <div class="match-condition">
                    <span v-if="row.http_status_code > 0" :class="['code-badge', getCodeClass(row.http_status_code)]">
                      HTTP {{ row.http_status_code }}
                    </span>
                    <span v-if="row.http_status_code > 0 && row.keyword" class="condition-sep">+</span>
                    <code v-if="row.keyword" class="keyword-code">{{ row.keyword }}</code>
                    <span v-if="!row.http_status_code && !row.keyword" class="muted-text">匹配所有错误</span>
                  </div>
                </td>
                <td>
                  <span class="description">{{ row.description || '-' }}</span>
                </td>
                <td>
                  <span class="priority-num">{{ row.priority }}</span>
                </td>
                <td>
                  <button :class="['toggle-switch', { active: row.enabled }]" @click="toggleRuleEnabled(row)">
                    <span class="toggle-thumb"></span>
                  </button>
                </td>
                <td>
                  <div class="action-btns">
                    <button class="action-btn edit" @click="showEditRuleDialog(row)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn delete" @click="confirmDeleteRule(row)" title="删除">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <div v-if="loadingRules" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>

          <div v-if="!loadingRules && rules.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <polygon points="12 2 2 7 12 12 22 7 12 2"/>
              <polyline points="2 17 12 22 22 17"/>
              <polyline points="2 12 12 17 22 12"/>
            </svg>
            <span>暂无错误规则配置</span>
          </div>
        </div>

        <div v-if="rulePagination.total > 0" class="table-footer">
          <div class="pagination-info">共 {{ rulePagination.total }} 条</div>
          <div class="pagination-controls">
            <select v-model="rulePagination.pageSize" @change="loadRules" class="page-size-select">
              <option :value="20">20 条/页</option>
              <option :value="50">50 条/页</option>
              <option :value="100">100 条/页</option>
            </select>
            <div class="page-btns">
              <button class="page-btn" :disabled="rulePagination.page <= 1" @click="rulePagination.page--; loadRules()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
              </button>
              <span class="page-info">{{ rulePagination.page }} / {{ Math.ceil(rulePagination.total / rulePagination.pageSize) || 1 }}</span>
              <button class="page-btn" :disabled="rulePagination.page >= Math.ceil(rulePagination.total / rulePagination.pageSize)" @click="rulePagination.page++; loadRules()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="9 18 15 12 9 6"/></svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误消息对话框 -->
    <div v-if="messageDialogVisible" class="modal-overlay" @click.self="messageDialogVisible = false">
      <div class="modal">
        <div class="modal-header">
          <h3 class="modal-title">{{ isEditMessage ? '编辑错误消息' : '新增错误消息' }}</h3>
          <button class="modal-close" @click="messageDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">状态码 <span class="required">*</span></label>
            <select v-model="messageForm.code" class="form-select" :disabled="isEditMessage">
              <option :value="400">400 - Bad Request</option>
              <option :value="401">401 - Unauthorized</option>
              <option :value="403">403 - Forbidden</option>
              <option :value="429">429 - Too Many Requests</option>
              <option :value="500">500 - Internal Server Error</option>
              <option :value="502">502 - Bad Gateway</option>
              <option :value="503">503 - Service Unavailable</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">错误类型 <span class="required">*</span></label>
            <input v-model="messageForm.error_type" class="form-input" :disabled="isEditMessage" placeholder="如: auth_failed" />
          </div>
          <div class="form-group">
            <label class="form-label">自定义消息 <span class="required">*</span></label>
            <textarea v-model="messageForm.custom_message" class="form-textarea" rows="3" placeholder="返回给用户的错误消息"></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">说明</label>
            <input v-model="messageForm.description" class="form-input" placeholder="管理员备注（可选）" />
          </div>
          <div class="form-group form-inline">
            <label class="form-label">启用</label>
            <button :class="['toggle-switch', { active: messageForm.enabled }]" @click="messageForm.enabled = !messageForm.enabled">
              <span class="toggle-thumb"></span>
            </button>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="messageDialogVisible = false">取消</button>
          <button class="btn btn-primary" @click="submitMessage" :disabled="submittingMessage">
            <span v-if="submittingMessage" class="btn-spinner"></span>
            确定
          </button>
        </div>
      </div>
    </div>

    <!-- 错误规则对话框 -->
    <div v-if="ruleDialogVisible" class="modal-overlay" @click.self="ruleDialogVisible = false">
      <div class="modal modal-lg">
        <div class="modal-header">
          <h3 class="modal-title">{{ isEditRule ? '编辑规则' : '新增规则' }}</h3>
          <button class="modal-close" @click="ruleDialogVisible = false">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">HTTP状态码</label>
            <div class="form-row">
              <input type="number" v-model.number="ruleForm.http_status_code" class="form-input" min="0" max="599" placeholder="0" style="width: 120px" />
              <span class="form-hint">0 表示匹配任意状态码</span>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">错误关键词</label>
            <input v-model="ruleForm.keyword" class="form-input" placeholder="匹配错误消息中的关键词（不区分大小写，留空表示任意）" />
          </div>
          <div class="form-group">
            <label class="form-label">目标状态 <span class="required">*</span></label>
            <select v-model="ruleForm.target_status" class="form-select">
              <option value="valid">valid - 不修改（仅记录错误）</option>
              <option value="invalid">invalid - 无效（禁用账户）</option>
              <option value="rate_limited">rate_limited - 限流（1小时后恢复）</option>
              <option value="overloaded">overloaded - 过载（临时不可用）</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">优先级</label>
            <div class="form-row">
              <input type="number" v-model.number="ruleForm.priority" class="form-input" min="0" max="1000" style="width: 120px" />
              <span class="form-hint">数值越大优先级越高</span>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">描述</label>
            <input v-model="ruleForm.description" class="form-input" placeholder="规则说明（可选）" />
          </div>
          <div class="form-group form-inline">
            <label class="form-label">启用</label>
            <button :class="['toggle-switch', { active: ruleForm.enabled }]" @click="ruleForm.enabled = !ruleForm.enabled">
              <span class="toggle-thumb"></span>
            </button>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="ruleDialogVisible = false">取消</button>
          <button class="btn btn-primary" @click="submitRule" :disabled="submittingRule">
            <span v-if="submittingRule" class="btn-spinner"></span>
            确定
          </button>
        </div>
      </div>
    </div>

    <!-- 确认对话框 -->
    <div v-if="confirmDialogVisible" class="modal-overlay" @click.self="confirmDialogVisible = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3 class="modal-title">{{ confirmTitle }}</h3>
        </div>
        <div class="modal-body">
          <p class="confirm-message">{{ confirmMessage }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline" @click="confirmDialogVisible = false">取消</button>
          <button class="btn btn-danger" @click="confirmAction">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import api from '@/api'
import { ElMessage, ElMessageBox } from '@/utils/toast'

const activeTab = ref('messages')

// 确认对话框
const confirmDialogVisible = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
let confirmCallback = null

function showConfirm(title, message, callback) {
  confirmTitle.value = title
  confirmMessage.value = message
  confirmCallback = callback
  confirmDialogVisible.value = true
}

function confirmAction() {
  confirmDialogVisible.value = false
  if (confirmCallback) confirmCallback()
}

// ==================== 错误消息相关 ====================
const loadingMessages = ref(false)
const initializingMessages = ref(false)
const refreshingMessages = ref(false)
const submittingMessage = ref(false)
const enablingAll = ref(false)
const disablingAll = ref(false)
const messageDialogVisible = ref(false)
const isEditMessage = ref(false)
const messages = ref([])

const messageForm = reactive({
  id: null,
  code: 400,
  error_type: '',
  custom_message: '',
  description: '',
  enabled: true
})

async function loadMessages() {
  loadingMessages.value = true
  try {
    const res = await api.getErrorMessages()
    messages.value = res.data || res || []
  } catch (e) {
    ElMessage.error('加载失败')
  }
  finally { loadingMessages.value = false }
}

function showCreateMessageDialog() {
  isEditMessage.value = false
  Object.assign(messageForm, { id: null, code: 400, error_type: '', custom_message: '', description: '', enabled: true })
  messageDialogVisible.value = true
}

function showEditMessageDialog(row) {
  isEditMessage.value = true
  Object.assign(messageForm, { id: row.id, code: row.code, error_type: row.error_type, custom_message: row.custom_message, description: row.description || '', enabled: row.enabled })
  messageDialogVisible.value = true
}

async function submitMessage() {
  if (!messageForm.code || !messageForm.error_type || !messageForm.custom_message) {
    ElMessage.error('请填写必填项')
    return
  }
  submittingMessage.value = true
  try {
    if (isEditMessage.value) {
      await api.updateErrorMessage(messageForm.id, { custom_message: messageForm.custom_message, enabled: messageForm.enabled, description: messageForm.description })
      ElMessage.success('更新成功')
    } else {
      await api.createErrorMessage({ code: messageForm.code, error_type: messageForm.error_type, custom_message: messageForm.custom_message, enabled: messageForm.enabled, description: messageForm.description })
      ElMessage.success('创建成功')
    }
    messageDialogVisible.value = false
    loadMessages()
  } catch (e) {
    ElMessage.error('操作失败')
  }
  finally { submittingMessage.value = false }
}

async function toggleMessageEnabled(row) {
  try {
    await api.toggleErrorMessage(row.id)
    row.enabled = !row.enabled
    ElMessage.success('状态已更新')
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

function confirmDeleteMessage(row) {
  showConfirm('确认删除', '确定删除该错误消息配置?', async () => {
    try {
      await api.deleteErrorMessage(row.id)
      ElMessage.success('删除成功')
      loadMessages()
    } catch (e) {
      ElMessage.error('删除失败')
    }
  })
}

async function initMessageDefaults() {
  initializingMessages.value = true
  try {
    await api.initErrorMessages()
    ElMessage.success('默认配置已初始化')
    loadMessages()
  } catch (e) {
    ElMessage.error('初始化失败')
  }
  finally { initializingMessages.value = false }
}

async function refreshMessageCache() {
  refreshingMessages.value = true
  try {
    await api.refreshErrorMessages()
    ElMessage.success('缓存已刷新')
  } catch (e) {
    ElMessage.error('刷新失败')
  }
  finally { refreshingMessages.value = false }
}

async function enableAllMessages() {
  enablingAll.value = true
  try {
    const res = await api.enableAllErrorMessages()
    ElMessage.success(`已启用所有配置（影响 ${res.data?.affected || 0} 条）`)
    loadMessages()
  } catch (e) {
    ElMessage.error('操作失败')
  }
  finally { enablingAll.value = false }
}

async function disableAllMessages() {
  disablingAll.value = true
  try {
    const res = await api.disableAllErrorMessages()
    ElMessage.success(`已禁用所有配置（影响 ${res.data?.affected || 0} 条）`)
    loadMessages()
  } catch (e) {
    ElMessage.error('操作失败')
  }
  finally { disablingAll.value = false }
}

// ==================== 错误规则相关 ====================
const loadingRules = ref(false)
const resettingRules = ref(false)
const refreshingRules = ref(false)
const submittingRule = ref(false)
const enablingAllRules = ref(false)
const disablingAllRules = ref(false)
const ruleDialogVisible = ref(false)
const isEditRule = ref(false)
const rules = ref([])

const rulePagination = reactive({ page: 1, pageSize: 50, total: 0 })

const ruleForm = reactive({
  id: null,
  http_status_code: 0,
  keyword: '',
  target_status: 'valid',
  priority: 50,
  description: '',
  enabled: true
})

async function loadRules() {
  loadingRules.value = true
  try {
    const res = await api.getErrorRules({ page: rulePagination.page, page_size: rulePagination.pageSize })
    rules.value = res.data?.items || []
    rulePagination.total = res.data?.total || 0
  } catch (e) {
    ElMessage.error('加载失败')
  }
  finally { loadingRules.value = false }
}

function showCreateRuleDialog() {
  isEditRule.value = false
  Object.assign(ruleForm, { id: null, http_status_code: 0, keyword: '', target_status: 'valid', priority: 50, description: '', enabled: true })
  ruleDialogVisible.value = true
}

function showEditRuleDialog(row) {
  isEditRule.value = true
  Object.assign(ruleForm, { id: row.id, http_status_code: row.http_status_code, keyword: row.keyword, target_status: row.target_status, priority: row.priority, description: row.description || '', enabled: row.enabled })
  ruleDialogVisible.value = true
}

async function submitRule() {
  if (!ruleForm.target_status) {
    ElMessage.error('请选择目标状态')
    return
  }
  submittingRule.value = true
  try {
    const data = { http_status_code: ruleForm.http_status_code || 0, keyword: ruleForm.keyword || '', target_status: ruleForm.target_status, priority: ruleForm.priority, description: ruleForm.description, enabled: ruleForm.enabled }
    if (isEditRule.value) {
      await api.updateErrorRule(ruleForm.id, data)
      ElMessage.success('更新成功')
    } else {
      await api.createErrorRule(data)
      ElMessage.success('创建成功')
    }
    ruleDialogVisible.value = false
    loadRules()
  } catch (e) {
    ElMessage.error('操作失败')
  }
  finally { submittingRule.value = false }
}

async function toggleRuleEnabled(row) {
  try {
    await api.updateErrorRule(row.id, { enabled: !row.enabled })
    row.enabled = !row.enabled
    ElMessage.success('状态已更新')
  } catch (e) {
    ElMessage.error('操作失败')
  }
}

function confirmDeleteRule(row) {
  showConfirm('确认删除', '确定删除该规则?', async () => {
    try {
      await api.deleteErrorRule(row.id)
      ElMessage.success('删除成功')
      loadRules()
    } catch (e) {
      ElMessage.error('删除失败')
    }
  })
}

async function resetRulesToDefault() {
  showConfirm('确认重置', '此操作将删除所有现有规则并恢复为默认规则，是否继续？', async () => {
    resettingRules.value = true
    try {
      await api.resetErrorRules()
      ElMessage.success('已重置为默认规则')
      loadRules()
    } catch (e) {
      ElMessage.error('重置失败')
    }
    finally { resettingRules.value = false }
  })
}

async function refreshRuleCache() {
  refreshingRules.value = true
  try {
    const res = await api.refreshErrorRulesCache()
    ElMessage.success(`缓存已刷新，共 ${res.data?.rule_count || 0} 条规则`)
  } catch (e) {
    ElMessage.error('刷新失败')
  }
  finally { refreshingRules.value = false }
}

async function enableAllRules() {
  enablingAllRules.value = true
  try {
    const res = await api.enableAllErrorRules()
    ElMessage.success(`已启用所有规则（影响 ${res.data?.affected || 0} 条）`)
    loadRules()
  } catch (e) {
    ElMessage.error('操作失败')
  }
  finally { enablingAllRules.value = false }
}

async function disableAllRules() {
  disablingAllRules.value = true
  try {
    const res = await api.disableAllErrorRules()
    ElMessage.success(`已禁用所有规则（影响 ${res.data?.affected || 0} 条）`)
    loadRules()
  } catch (e) {
    ElMessage.error('操作失败')
  }
  finally { disablingAllRules.value = false }
}

// ==================== 通用方法 ====================
function getCodeClass(code) {
  if (code >= 500) return 'danger'
  if (code === 429) return 'warning'
  if (code >= 400) return 'info'
  return ''
}

function getTargetStatusClass(status) {
  switch (status) {
    case 'valid': return 'success'
    case 'invalid': return 'danger'
    case 'rate_limited': return 'warning'
    case 'overloaded': return 'info'
    default: return ''
  }
}

function getTargetStatusLabel(status) {
  switch (status) {
    case 'valid': return '忽略'
    case 'invalid': return '禁用账户'
    case 'rate_limited': return '临时限流'
    case 'overloaded': return '过载'
    default: return status
  }
}

// 切换 Tab 时加载数据
watch(activeTab, (val) => {
  if (val === 'messages' && messages.value.length === 0) loadMessages()
  if (val === 'rules' && rules.value.length === 0) loadRules()
})

onMounted(() => {
  loadMessages()
})
</script>

<style scoped>
.error-config-page {
  padding: var(--apple-spacing-xl);
  max-width: 1400px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  margin-bottom: var(--apple-spacing-lg);
}

.header-content {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--apple-text-primary);
  letter-spacing: -0.5px;
}

.page-subtitle {
  font-size: 15px;
  color: var(--apple-text-secondary);
}

/* Tab 切换 */
.tab-container {
  margin-bottom: var(--apple-spacing-lg);
}

.tab-nav {
  display: inline-flex;
  background: var(--apple-fill-quaternary);
  padding: 4px;
  border-radius: var(--apple-radius-md);
  gap: 4px;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
  color: var(--apple-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--apple-radius-sm);
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn svg {
  width: 18px;
  height: 18px;
}

.tab-btn:hover {
  color: var(--apple-text-primary);
}

.tab-btn.active {
  background: var(--apple-bg-primary);
  color: var(--apple-text-primary);
  box-shadow: var(--apple-shadow-sm);
}

/* 指南卡片 */
.guide-card {
  background: linear-gradient(135deg, rgba(0, 122, 255, 0.05) 0%, rgba(88, 86, 214, 0.05) 100%);
  border: 1px solid rgba(0, 122, 255, 0.15);
  border-radius: var(--apple-radius-lg);
  padding: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-lg);
}

.guide-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-md);
}

.guide-header svg {
  width: 20px;
  height: 20px;
  color: var(--apple-blue);
}

.guide-content {
  display: flex;
  flex-wrap: wrap;
  gap: var(--apple-spacing-lg);
}

.guide-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.guide-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 600;
  border-radius: var(--apple-radius-full);
}

.guide-badge.danger {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

.guide-badge.warning {
  background: var(--apple-orange-light);
  color: var(--apple-orange);
}

.guide-badge.info {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.guide-text {
  font-size: 13px;
  color: var(--apple-text-secondary);
}

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: var(--apple-spacing-md);
  margin-bottom: var(--apple-spacing-lg);
}

.action-group {
  display: flex;
  gap: var(--apple-spacing-sm);
  flex-wrap: wrap;
}

/* 按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn svg {
  width: 16px;
  height: 16px;
}

.btn-primary {
  background: var(--apple-blue);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--apple-blue-hover);
}

.btn-success {
  background: var(--apple-green);
  color: white;
}

.btn-success:hover:not(:disabled) {
  background: #2db84e;
}

.btn-warning {
  background: var(--apple-orange);
  color: white;
}

.btn-warning:hover:not(:disabled) {
  background: #e68600;
}

.btn-danger {
  background: var(--apple-red);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #e6352b;
}

.btn-outline {
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-primary);
  border: 1px solid var(--apple-separator);
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-fill-tertiary);
}

.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
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
}

.data-table th {
  padding: 14px 16px;
  font-size: 12px;
  font-weight: 600;
  color: var(--apple-text-secondary);
  text-align: left;
  background: var(--apple-bg-secondary);
  border-bottom: 1px solid var(--apple-separator);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-table td {
  padding: 14px 16px;
  font-size: 14px;
  color: var(--apple-text-primary);
  border-bottom: 1px solid var(--apple-separator);
  vertical-align: middle;
}

.data-table tr:last-child td {
  border-bottom: none;
}

.data-table tr:hover td {
  background: var(--apple-fill-quaternary);
}

/* 状态码徽章 */
.code-badge {
  display: inline-flex;
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 600;
  border-radius: var(--apple-radius-full);
  font-family: 'SF Mono', Monaco, monospace;
}

.code-badge.danger {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

.code-badge.warning {
  background: var(--apple-orange-light);
  color: var(--apple-orange);
}

.code-badge.info {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

/* 状态徽章 */
.status-badge {
  display: inline-flex;
  padding: 5px 12px;
  font-size: 12px;
  font-weight: 600;
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

.status-badge.warning {
  background: var(--apple-orange-light);
  color: var(--apple-orange);
}

.status-badge.info {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

/* 匹配条件 */
.match-condition {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.condition-sep {
  color: var(--apple-text-tertiary);
  font-size: 12px;
}

.keyword-code {
  background: var(--apple-fill-secondary);
  padding: 3px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'SF Mono', Monaco, monospace;
  color: var(--apple-orange);
}

/* 其他文本样式 */
.error-type {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 13px;
  color: var(--apple-text-primary);
}

.custom-message {
  font-size: 13px;
  color: var(--apple-text-primary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.description {
  font-size: 13px;
  color: var(--apple-text-secondary);
}

.muted-text {
  color: var(--apple-text-tertiary);
  font-size: 13px;
}

.priority-num {
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 13px;
  color: var(--apple-text-secondary);
}

/* 开关按钮 */
.toggle-switch {
  position: relative;
  width: 44px;
  height: 24px;
  background: var(--apple-fill-secondary);
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.toggle-switch.active {
  background: var(--apple-green);
}

.toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s ease;
}

.toggle-switch.active .toggle-thumb {
  transform: translateX(20px);
}

/* 操作按钮 */
.action-btns {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary);
  border: none;
  border-radius: var(--apple-radius-sm);
  color: var(--apple-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn.edit:hover {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.action-btn.delete:hover {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  gap: 16px;
  color: var(--apple-text-secondary);
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--apple-fill-tertiary);
  border-top-color: var(--apple-blue);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  gap: 16px;
  color: var(--apple-text-tertiary);
}

.empty-state svg {
  width: 48px;
  height: 48px;
}

/* 表格底部 */
.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid var(--apple-separator);
  background: var(--apple-bg-secondary);
}

.pagination-info {
  font-size: 13px;
  color: var(--apple-text-secondary);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-size-select {
  padding: 6px 28px 6px 12px;
  font-size: 13px;
  border: 1px solid var(--apple-separator);
  border-radius: var(--apple-radius-sm);
  background: var(--apple-bg-primary);
  color: var(--apple-text-primary);
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%2386868b' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 8px center;
}

.page-btns {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary);
  border: 1px solid var(--apple-separator);
  border-radius: var(--apple-radius-sm);
  color: var(--apple-text-primary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-btn:not(:disabled):hover {
  background: var(--apple-fill-tertiary);
}

.page-btn svg {
  width: 16px;
  height: 16px;
}

.page-info {
  font-size: 13px;
  color: var(--apple-text-secondary);
  min-width: 60px;
  text-align: center;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

.modal {
  width: 100%;
  max-width: 480px;
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  box-shadow: var(--apple-shadow-modal);
  animation: scaleIn 0.25s ease;
}

.modal-sm {
  max-width: 360px;
}

.modal-lg {
  max-width: 560px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--apple-separator);
}

.modal-title {
  font-size: 17px;
  font-weight: 600;
  color: var(--apple-text-primary);
}

.modal-close {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-tertiary);
  border: none;
  border-radius: 50%;
  color: var(--apple-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
}

.modal-close:hover {
  background: var(--apple-fill-secondary);
  color: var(--apple-text-primary);
}

.modal-close svg {
  width: 14px;
  height: 14px;
}

.modal-body {
  padding: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--apple-separator);
  background: var(--apple-bg-secondary);
}

/* 确认消息 */
.confirm-message {
  font-size: 14px;
  color: var(--apple-text-secondary);
  line-height: 1.6;
}

/* 表单 */
.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-inline {
  display: flex;
  align-items: center;
  gap: 16px;
}

.form-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--apple-text-primary);
  margin-bottom: 8px;
}

.form-inline .form-label {
  margin-bottom: 0;
}

.required {
  color: var(--apple-red);
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 10px 14px;
  font-size: 14px;
  color: var(--apple-text-primary);
  background: var(--apple-fill-quaternary);
  border: 1px solid transparent;
  border-radius: var(--apple-radius-md);
  transition: all 0.2s ease;
}

.form-input::placeholder,
.form-textarea::placeholder {
  color: var(--apple-text-placeholder);
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  background: var(--apple-bg-primary);
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
  outline: none;
}

.form-input:disabled,
.form-select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2386868b' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 40px;
  cursor: pointer;
}

.form-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.form-hint {
  font-size: 12px;
  color: var(--apple-text-tertiary);
}

/* Toast 通知 */
.toast {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: var(--apple-text-primary);
  color: white;
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-lg);
  z-index: 2000;
  animation: slideUp 0.3s ease;
}

.toast svg {
  width: 18px;
  height: 18px;
}

.toast.success {
  background: var(--apple-green);
}

.toast.error {
  background: var(--apple-red);
}

/* 动画 */
@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes scaleIn {
  from { opacity: 0; transform: scale(0.95); }
  to { opacity: 1; transform: scale(1); }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateX(-50%) translateY(20px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}

.spinning {
  animation: spin 1s linear infinite;
}

/* 响应式 */
@media (max-width: 768px) {
  .error-config-page {
    padding: var(--apple-spacing-md);
  }

  .page-title {
    font-size: 24px;
  }

  .action-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .action-group {
    justify-content: flex-start;
  }

  .guide-content {
    flex-direction: column;
    gap: var(--apple-spacing-sm);
  }
}
</style>
