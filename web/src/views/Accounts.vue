<!--
 * 文件作用：账户管理页面 - Apple HIG 风格
 * 负责功能：
 *   - 多平台账户展示（Claude/OpenAI/Gemini）
 *   - 账户状态管理（启用/禁用/健康检测）
 *   - 用量统计展示（5H/7D进度条）
 *   - 账户CRUD操作
 *   - Token刷新和强制恢复
 * 重要程度：⭐⭐⭐⭐⭐ 核心（账户管理）
-->
<template>
  <div class="accounts-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">账户管理</h1>
        <p class="page-subtitle">管理 Claude、Gemini、OpenAI 等账户与代理配置</p>
      </div>
      <div class="header-actions">
        <button class="btn btn-outline" @click="loadAccounts" :disabled="loading">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: loading }">
            <polyline points="23,4 23,10 17,10"/>
            <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
          </svg>
          刷新
        </button>
        <button class="btn btn-primary" @click="showFormDialog = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          添加账户
        </button>
      </div>
    </div>

    <!-- 平台统计卡片 -->
    <div class="platform-stats">
      <div
        v-for="platform in platformStats"
        :key="platform.key"
        :class="['stat-card', { active: filters.platform === platform.key }]"
        @click="filterByPlatform(platform.key)"
      >
        <div class="stat-icon" :style="{ background: platform.gradient }">
          <svg v-if="platform.key === 'claude'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2a9 9 0 019 9c0 3.9-2.5 7.2-6 8.4V22l-3-2-3 2v-2.6c-3.5-1.2-6-4.5-6-8.4a9 9 0 019-9z"/>
          </svg>
          <svg v-else-if="platform.key === 'openai'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="3"/>
            <path d="M12 1v4M12 19v4M4.22 4.22l2.83 2.83M16.95 16.95l2.83 2.83M1 12h4M19 12h4M4.22 19.78l2.83-2.83M16.95 7.05l2.83-2.83"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="12,2 2,7 12,12 22,7"/>
            <polyline points="2,17 12,22 22,17"/>
            <polyline points="2,12 12,17 22,12"/>
          </svg>
        </div>
        <div class="stat-info">
          <h3>{{ platform.name }}</h3>
          <div class="stat-numbers">
            <span class="total">{{ platform.count }} 个账户</span>
            <span class="valid">
              <span class="valid-dot"></span>
              {{ platform.validCount }} 可用
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 筛选和搜索 -->
    <div class="filter-bar">
      <div class="filter-left">
        <select v-model="filters.status" @change="loadAccounts" class="filter-select">
          <option value="">全部状态</option>
          <option value="valid">正常</option>
          <option value="invalid">无效</option>
          <option value="rate_limited">限流中</option>
          <option value="token_expired">Token过期</option>
          <option value="suspended">疑似封号</option>
          <option value="banned">已封号</option>
          <option value="disabled">已禁用</option>
        </select>
        <div class="search-box">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            v-model="filters.search"
            placeholder="搜索账户名称..."
            @input="handleSearch"
          />
        </div>
      </div>
      <div class="filter-right">
        <span v-if="filters.platform" class="filter-tag">
          {{ getPlatformName(filters.platform) }}
          <button @click="filters.platform = ''; loadAccounts()">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </span>
      </div>
    </div>

    <!-- 账户列表 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th class="col-checkbox">
                <input type="checkbox" :checked="selectedAccounts.length === accounts.length && accounts.length > 0" @change="toggleSelectAll" />
              </th>
              <th>#</th>
              <th class="col-account">账户</th>
              <th>平台/类型</th>
              <th>状态</th>
              <th>用量</th>
              <th class="col-center">启用</th>
              <th class="col-center">优先级</th>
              <th class="col-center">并发</th>
              <th class="col-right">请求数</th>
              <th class="col-right">总费用</th>
              <th>最后使用</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="(account, index) in accounts" :key="account.id">
              <td class="col-checkbox">
                <input type="checkbox" :checked="isSelected(account.id)" @change="toggleSelect(account.id)" />
              </td>
              <td class="row-index">{{ (pagination.page - 1) * pagination.pageSize + index + 1 }}</td>
              <td class="col-account">
                <div class="account-cell">
                  <div class="account-avatar" :style="{ background: getTypeColor(account.type) }">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                    </svg>
                  </div>
                  <div class="account-info">
                    <span class="account-name">{{ account.name }}</span>
                    <span class="account-type">{{ getTypeLabel(account.type, account) }}</span>
                  </div>
                </div>
              </td>
              <td>
                <span :class="['platform-badge', getPlatformClass(account.type)]">
                  {{ getPlatformLabel(account.type) }}
                </span>
              </td>
              <td>
                <div class="status-cell">
                  <span :class="['status-badge', account.status]">
                    <span class="status-dot"></span>
                    {{ getStatusLabel(account.status) }}
                  </span>
                  <div v-if="account.status === 'rate_limited' && account.rate_limit_reset_at" class="status-detail">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="12" cy="12" r="10"/>
                      <polyline points="12,6 12,12 16,14"/>
                    </svg>
                    {{ formatResetTime(account.rate_limit_reset_at) }}
                  </div>
                  <div v-if="account.status === 'suspended' && account.suspended_count > 0" class="status-detail warning">
                    连续失败 {{ account.suspended_count }} 次
                  </div>
                </div>
              </td>
              <td>
                <!-- Claude Official: 显示 5H/7D 进度条 -->
                <div v-if="account.type === 'claude-official' && hasUsageData(account)" class="usage-bars">
                  <div class="usage-bar-item" v-if="account.five_hour_utilization != null">
                    <div class="usage-bar-label">
                      <span class="label-text">5H</span>
                      <span class="label-value">{{ account.five_hour_utilization.toFixed(1) }}%</span>
                    </div>
                    <div class="usage-bar-track">
                      <div
                        :class="['usage-bar-fill', getUsageBarClass(account.five_hour_utilization)]"
                        :style="{ width: Math.min(account.five_hour_utilization, 100) + '%' }"
                      ></div>
                    </div>
                  </div>
                  <div class="usage-bar-item" v-if="account.seven_day_utilization != null">
                    <div class="usage-bar-label">
                      <span class="label-text">7D</span>
                      <span class="label-value">{{ account.seven_day_utilization.toFixed(1) }}%</span>
                    </div>
                    <div class="usage-bar-track">
                      <div
                        :class="['usage-bar-fill', getUsageBarClass(account.seven_day_utilization)]"
                        :style="{ width: Math.min(account.seven_day_utilization, 100) + '%' }"
                      ></div>
                    </div>
                  </div>
                </div>
                <!-- 其他类型: 显示预算进度条或今日统计 -->
                <div v-else-if="account.daily_budget > 0" class="usage-bars">
                  <div class="usage-bar-item">
                    <div class="usage-bar-label">
                      <span class="label-text">今日</span>
                      <span class="label-value">${{ formatCost(account.today_cost || 0) }} / ${{ formatCost(account.daily_budget) }}</span>
                    </div>
                    <div class="usage-bar-track">
                      <div
                        :class="['usage-bar-fill', getUsageBarClass(account.budget_utilization || 0)]"
                        :style="{ width: Math.min(account.budget_utilization || 0, 100) + '%' }"
                      ></div>
                    </div>
                  </div>
                </div>
                <div v-else-if="account.today_tokens > 0 || account.today_count > 0" class="usage-stats">
                  <span class="usage-stat">${{ formatCost(account.today_cost || 0) }}</span>
                  <span class="usage-stat">{{ formatTokens(account.today_tokens) }} tokens</span>
                </div>
                <span v-else class="no-data">-</span>
              </td>
              <td class="col-center">
                <label class="toggle-switch">
                  <input type="checkbox" v-model="account.enabled" @change="handleToggleEnabled(account)" />
                  <span class="toggle-slider"></span>
                </label>
              </td>
              <td class="col-center">
                <span class="badge badge-info">{{ account.priority }}</span>
              </td>
              <td class="col-center">
                <span :class="['concurrency-badge', getConcurrencyClass(account)]">
                  {{ account.current_concurrency || 0 }} / {{ account.max_concurrency || 5 }}
                </span>
              </td>
              <td class="col-right">
                <span class="mono">{{ formatNumber(account.request_count) }}</span>
              </td>
              <td class="col-right">
                <span v-if="account.total_cost > 0" class="cost">${{ formatCost(account.total_cost) }}</span>
                <span v-else class="no-data">-</span>
              </td>
              <td>
                <span v-if="account.last_used_at" class="time-ago">{{ formatRelativeTime(account.last_used_at) }}</span>
                <span v-else class="no-data">-</span>
              </td>
              <td>
                <div class="action-group">
                  <button class="action-btn" @click="handleEdit(account)" title="编辑">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                  </button>
                  <button v-if="canHealthCheck(account)" class="action-btn info" @click="handleHealthCheck(account)" :disabled="healthCheckingIds.includes(account.id)" title="检测">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
                    </svg>
                  </button>
                  <button v-if="canRecover(account)" class="action-btn success" @click="handleForceRecover(account)" :disabled="recoveringIds.includes(account.id)" title="恢复">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="23,4 23,10 17,10"/>
                      <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
                    </svg>
                  </button>
                  <button v-if="canRefreshToken(account)" class="action-btn warning" @click="handleRefreshToken(account)" :disabled="refreshingIds.includes(account.id)" title="刷新Token">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                    </svg>
                  </button>
                  <button class="action-btn danger" @click="confirmDelete(account)" title="删除">
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

        <div v-if="!loading && accounts.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
          </svg>
          <span>暂无账户</span>
        </div>
      </div>

      <!-- 分页 -->
      <div class="table-footer">
        <div class="selection-info" v-if="selectedAccounts.length > 0">
          已选择 {{ selectedAccounts.length }} 项
          <button class="link-btn danger" @click="handleBatchDelete">批量删除</button>
        </div>
        <div class="pagination-wrap">
          <div class="pagination-info">共 {{ pagination.total }} 条</div>
          <div class="pagination-controls">
            <select v-model="pagination.pageSize" @change="loadAccounts" class="page-size-select">
              <option :value="10">10 条/页</option>
              <option :value="20">20 条/页</option>
              <option :value="50">50 条/页</option>
            </select>
            <div class="page-btns">
              <button class="page-btn" :disabled="pagination.page <= 1" @click="pagination.page--; loadAccounts()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="15,18 9,12 15,6"/>
                </svg>
              </button>
              <span class="page-current">{{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) || 1 }}</span>
              <button class="page-btn" :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)" @click="pagination.page++; loadAccounts()">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="9,18 15,12 9,6"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

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
            <p class="delete-message">确定要删除账户 "{{ deleteTarget?.name }}" 吗？此操作不可撤销。</p>
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

    <!-- 添加/编辑弹窗 -->
    <AccountForm
      v-model="showFormDialog"
      :edit-data="editingAccount"
      @success="handleFormSuccess"
      @update:modelValue="handleDialogClose"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from '@/utils/toast'
import AccountForm from '@/components/AccountForm.vue'
import api from '@/api'

const loading = ref(false)
const accounts = ref([])
const selectedAccounts = ref([])
const showFormDialog = ref(false)
const editingAccount = ref(null)

// 删除相关
const deleteDialogVisible = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

// 操作加载状态
const healthCheckingIds = ref([])
const recoveringIds = ref([])
const refreshingIds = ref([])

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const filters = reactive({
  platform: '',
  status: '',
  search: ''
})

// 平台分组定义
const platformGroups = [
  {
    key: 'claude',
    name: 'Claude',
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    types: ['claude-official', 'claude-console', 'bedrock']
  },
  {
    key: 'openai',
    name: 'OpenAI',
    gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)',
    types: ['openai', 'openai-responses', 'azure-openai']
  },
  {
    key: 'gemini',
    name: 'Gemini',
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    types: ['gemini']
  }
]

// 子类型定义
const subtypeMap = {
  'claude-official': { label: 'Claude Official', color: '#667eea', platform: 'Claude' },
  'claude-console': { label: 'Claude Console', color: '#764ba2', platform: 'Claude' },
  'bedrock': { label: 'AWS Bedrock', color: '#ff9900', platform: 'Claude' },
  'openai': { label: 'OpenAI API', color: '#11998e', platform: 'OpenAI' },
  'openai-responses': { label: 'ChatGPT', color: '#38ef7d', platform: 'OpenAI' },
  'azure-openai': { label: 'Azure OpenAI', color: '#0078d4', platform: 'OpenAI' },
  'gemini': { label: 'Gemini', color: '#4facfe', platform: 'Gemini' }
}

// 平台统计
const platformStats = computed(() => {
  return platformGroups.map(group => {
    const platformAccounts = accounts.value.filter(a => group.types.includes(a.type))
    return {
      ...group,
      count: platformAccounts.length,
      validCount: platformAccounts.filter(a => a.status === 'valid' && a.enabled).length
    }
  })
})

function isSelected(id) {
  return selectedAccounts.value.includes(id)
}

function toggleSelect(id) {
  const index = selectedAccounts.value.indexOf(id)
  if (index >= 0) {
    selectedAccounts.value.splice(index, 1)
  } else {
    selectedAccounts.value.push(id)
  }
}

function toggleSelectAll() {
  if (selectedAccounts.value.length === accounts.value.length) {
    selectedAccounts.value = []
  } else {
    selectedAccounts.value = accounts.value.map(a => a.id)
  }
}

function getTypeLabel(type, row = null) {
  const baseLabel = subtypeMap[type]?.label || type
  if (!row) return baseLabel

  if (type === 'claude-official') {
    if (row.access_token && row.session_key) return 'OAuth + SK'
    if (row.access_token) return 'OAuth'
    if (row.session_key) return 'SessionKey'
    if (row.api_key) return 'API Key'
    return 'Claude Official'
  }

  if (type === 'openai') {
    if (row.base_url) {
      try {
        const url = new URL(row.base_url)
        const host = url.hostname.replace('www.', '')
        if (host.includes('openrouter')) return 'OpenRouter'
        if (host.includes('together')) return 'Together AI'
        if (host.includes('groq')) return 'Groq'
        if (host.includes('deepseek')) return 'DeepSeek'
        return host.split('.')[0]
      } catch {
        return '三方 API'
      }
    }
    return 'OpenAI'
  }

  return baseLabel
}

function getTypeColor(type) {
  return subtypeMap[type]?.color || '#999'
}

function getPlatformLabel(type) {
  return subtypeMap[type]?.platform || 'Unknown'
}

function getPlatformClass(type) {
  const platform = subtypeMap[type]?.platform?.toLowerCase()
  return platform || 'unknown'
}

function getPlatformName(key) {
  return platformGroups.find(g => g.key === key)?.name || key
}

function getStatusLabel(status) {
  const map = {
    valid: '正常',
    invalid: '无效',
    rate_limited: '限流中',
    overloaded: '过载',
    token_expired: 'Token过期',
    suspended: '疑似封号',
    banned: '已封号',
    disabled: '已禁用'
  }
  return map[status] || status
}

function hasUsageData(row) {
  return row.five_hour_utilization != null || row.seven_day_utilization != null
}

function getUsageBarClass(utilization) {
  if (utilization >= 90) return 'danger'
  if (utilization >= 70) return 'warning'
  return 'normal'
}

function getConcurrencyClass(row) {
  const current = row.current_concurrency || 0
  const max = row.max_concurrency || 5
  if (current >= max) return 'danger'
  if (current >= max * 0.8) return 'warning'
  return 'normal'
}

function formatNumber(num) {
  if (!num) return '0'
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toString()
}

function formatTokens(tokens) {
  if (!tokens) return '0'
  if (tokens >= 1000000) return (tokens / 1000000).toFixed(1) + 'M'
  if (tokens >= 1000) return (tokens / 1000).toFixed(1) + 'K'
  return tokens.toString()
}

function formatCost(cost) {
  if (!cost) return '0.00'
  if (cost >= 1000) return (cost / 1000).toFixed(2) + 'K'
  if (cost >= 1) return cost.toFixed(2)
  if (cost >= 0.01) return cost.toFixed(3)
  return cost.toFixed(4)
}

function formatRelativeTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now - date) / 1000)

  if (diff < 60) return '刚刚'
  if (diff < 3600) return Math.floor(diff / 60) + ' 分钟前'
  if (diff < 86400) return Math.floor(diff / 3600) + ' 小时前'
  if (diff < 604800) return Math.floor(diff / 86400) + ' 天前'
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}

function formatResetTime(dateStr) {
  if (!dateStr) return ''
  const resetTime = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((resetTime - now) / 1000)

  if (diff <= 0) return '即将恢复'
  if (diff < 60) return diff + ' 秒后恢复'
  if (diff < 3600) return Math.floor(diff / 60) + ' 分钟后恢复'
  const hours = Math.floor(diff / 3600)
  const mins = Math.floor((diff % 3600) / 60)
  return hours + '时' + mins + '分后恢复'
}

function canHealthCheck(row) {
  const oauthTypes = ['claude-official', 'openai-responses', 'gemini']
  return oauthTypes.includes(row.type)
}

function canRecover(row) {
  const recoverableStatuses = ['invalid', 'rate_limited', 'token_expired', 'suspended', 'banned', 'disabled']
  return recoverableStatuses.includes(row.status)
}

function canRefreshToken(row) {
  return row.type === 'claude-official' && row.session_key
}

function filterByPlatform(key) {
  if (filters.platform === key) {
    filters.platform = ''
  } else {
    filters.platform = key
  }
  loadAccounts()
}

let searchTimer = null
function handleSearch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    loadAccounts()
  }, 300)
}

onUnmounted(() => {
  if (searchTimer) {
    clearTimeout(searchTimer)
    searchTimer = null
  }
})

async function loadAccounts() {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      status: filters.status,
      search: filters.search
    }

    if (filters.platform) {
      const group = platformGroups.find(g => g.key === filters.platform)
      if (group) {
        params.types = group.types.join(',')
      }
    }

    const res = await api.getAccounts(params)
    accounts.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (e) {
    ElMessage.error('加载账户列表失败')
  } finally {
    loading.value = false
  }
}

async function handleToggleEnabled(row) {
  try {
    await api.updateAccount(row.id, { enabled: row.enabled })
    ElMessage.success('更新成功')
  } catch (e) {
    row.enabled = !row.enabled
    ElMessage.error('更新失败')
  }
}

function handleEdit(row) {
  editingAccount.value = { ...row }
  showFormDialog.value = true
}

function confirmDelete(account) {
  deleteTarget.value = account
  deleteDialogVisible.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.deleteAccount(deleteTarget.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    deleteTarget.value = null
    loadAccounts()
  } catch (e) {
    ElMessage.error('删除失败')
  } finally {
    deleting.value = false
  }
}

async function handleHealthCheck(row) {
  healthCheckingIds.value.push(row.id)
  try {
    const res = await api.checkAccountHealth(row.id)
    const data = res.data || res
    if (data.healthy) {
      ElMessage.success(`[${row.name}] ${data.message}`)
    } else {
      ElMessage.warning(`[${row.name}] ${data.message}`)
    }
    loadAccounts()
  } catch (e) {
    ElMessage.error('检测失败')
  } finally {
    healthCheckingIds.value = healthCheckingIds.value.filter(id => id !== row.id)
  }
}

async function handleForceRecover(row) {
  recoveringIds.value.push(row.id)
  try {
    await api.recoverAccount(row.id)
    ElMessage.success(`[${row.name}] 账号已强制恢复`)
    loadAccounts()
  } catch (e) {
    ElMessage.error('恢复失败')
  } finally {
    recoveringIds.value = recoveringIds.value.filter(id => id !== row.id)
  }
}

async function handleRefreshToken(row) {
  refreshingIds.value.push(row.id)
  try {
    await api.refreshAccountToken(row.id)
    ElMessage.success(`[${row.name}] Token 刷新成功`)
    loadAccounts()
  } catch (e) {
    ElMessage.error('Token 刷新失败')
  } finally {
    refreshingIds.value = refreshingIds.value.filter(id => id !== row.id)
  }
}

async function handleBatchDelete() {
  if (selectedAccounts.value.length === 0) return

  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedAccounts.value.length} 个账户吗？`,
      '批量删除',
      { type: 'warning' }
    )

    for (const id of selectedAccounts.value) {
      await api.deleteAccount(id)
    }
    ElMessage.success('批量删除成功')
    selectedAccounts.value = []
    loadAccounts()
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('批量删除失败')
    }
  }
}

function handleFormSuccess() {
  showFormDialog.value = false
  editingAccount.value = null
  loadAccounts()
}

function handleDialogClose(val) {
  if (!val) {
    editingAccount.value = null
  }
}

onMounted(() => {
  loadAccounts()
})
</script>

<style scoped>
.accounts-page {
  max-width: 1600px;
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

/* 平台统计卡片 */
.platform-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: var(--apple-spacing-md);
  margin-bottom: var(--apple-spacing-xl);
}

.stat-card {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-lg);
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  border: 2px solid transparent;
  box-shadow: var(--apple-shadow-card);
  cursor: pointer;
  transition: all var(--apple-duration-normal) var(--apple-ease-default);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--apple-shadow-card-hover);
}

.stat-card.active {
  border-color: var(--apple-blue);
  background: var(--apple-blue-light);
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

.stat-info h3 {
  margin: 0 0 var(--apple-spacing-xs);
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.stat-numbers {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
}

.stat-numbers .total {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.stat-numbers .valid {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
  font-size: var(--apple-text-sm);
  color: var(--apple-green);
}

.valid-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--apple-green);
}

/* 筛选栏 */
.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--apple-spacing-lg);
}

.filter-left {
  display: flex;
  gap: var(--apple-spacing-sm);
}

.filter-select {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
}

.search-box {
  display: flex;
  align-items: center;
  position: relative;
  width: 200px;
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
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
}

.filter-tag {
  display: inline-flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  background: var(--apple-blue-light);
  color: var(--apple-blue);
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.filter-tag button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: rgba(0, 122, 255, 0.2);
}

.filter-tag button svg {
  width: 10px;
  height: 10px;
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

.data-table tbody tr:hover {
  background: var(--apple-bg-secondary);
}

.col-checkbox { width: 40px; }
.col-account { min-width: 200px; }
.col-center { text-align: center; }
.col-right { text-align: right; }

.row-index {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

/* 账户单元格 */
.account-cell {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.account-avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.account-avatar svg {
  width: 18px;
  height: 18px;
  color: white;
}

.account-info {
  display: flex;
  flex-direction: column;
}

.account-name {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.account-type {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

/* 平台徽章 */
.platform-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
  padding: 2px 8px;
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.platform-badge.claude {
  background: #eef2ff;
  color: #667eea;
}

.platform-badge.openai {
  background: #ecfdf5;
  color: #059669;
}

.platform-badge.gemini {
  background: #eff6ff;
  color: #3b82f6;
}

/* 状态单元格 */
.status-cell {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xxs);
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.status-badge.valid { background: var(--apple-green-light); color: var(--apple-green); }
.status-badge.valid .status-dot { background: var(--apple-green); }

.status-badge.invalid { background: var(--apple-red-light); color: var(--apple-red); }
.status-badge.invalid .status-dot { background: var(--apple-red); }

.status-badge.rate_limited { background: var(--apple-orange-light); color: var(--apple-orange); }
.status-badge.rate_limited .status-dot { background: var(--apple-orange); }

.status-badge.token_expired { background: var(--apple-orange-light); color: #b45309; }
.status-badge.token_expired .status-dot { background: #f59e0b; }

.status-badge.suspended { background: #fed7aa; color: #c2410c; }
.status-badge.suspended .status-dot { background: #ea580c; }

.status-badge.banned { background: #fecaca; color: #991b1b; }
.status-badge.banned .status-dot { background: #b91c1c; }

.status-badge.disabled { background: var(--apple-fill-tertiary); color: var(--apple-text-tertiary); }
.status-badge.disabled .status-dot { background: var(--apple-text-tertiary); }

.status-detail {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
  font-size: 10px;
  color: var(--apple-text-tertiary);
}

.status-detail svg {
  width: 10px;
  height: 10px;
}

.status-detail.warning {
  color: #c2410c;
}

/* 用量进度条 */
.usage-bars {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
  min-width: 120px;
}

.usage-bar-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.usage-bar-label {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
}

.usage-bar-label .label-text {
  color: var(--apple-text-tertiary);
}

.usage-bar-label .label-value {
  font-family: var(--apple-font-mono);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.usage-bar-track {
  height: 4px;
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-full);
  overflow: hidden;
}

.usage-bar-fill {
  height: 100%;
  border-radius: var(--apple-radius-full);
  transition: width var(--apple-duration-normal);
}

.usage-bar-fill.normal { background: var(--apple-green); }
.usage-bar-fill.warning { background: var(--apple-orange); }
.usage-bar-fill.danger { background: var(--apple-red); }

.usage-stats {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.usage-stat {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.no-data {
  color: var(--apple-text-quaternary);
}

/* 开关样式 */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 36px;
  height: 20px;
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

/* 徽章 */
.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.badge-info {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-secondary);
}

.concurrency-badge {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
}

.concurrency-badge.normal { color: var(--apple-green); }
.concurrency-badge.warning { color: var(--apple-orange); }
.concurrency-badge.danger { color: var(--apple-red); }

.mono {
  font-family: var(--apple-font-mono);
  color: var(--apple-text-secondary);
}

.cost {
  font-family: var(--apple-font-mono);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-green);
}

.time-ago {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

/* 操作按钮 */
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

.action-btn:hover { background: var(--apple-blue); color: white; }
.action-btn.info:hover { background: var(--apple-teal); }
.action-btn.success:hover { background: var(--apple-green); }
.action-btn.warning:hover { background: var(--apple-orange); }
.action-btn.danger:hover { background: var(--apple-red); }

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 表格底部 */
.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
}

.selection-info {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.link-btn {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-blue);
  background: none;
  border: none;
  cursor: pointer;
}

.link-btn.danger {
  color: var(--apple-red);
}

.pagination-wrap {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
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
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }

  .platform-stats {
    grid-template-columns: 1fr;
  }

  .filter-bar {
    flex-direction: column;
    align-items: stretch;
    gap: var(--apple-spacing-sm);
  }

  .filter-left {
    flex-wrap: wrap;
  }
}
</style>
