<!--
 * 文件作用：模型管理页面 - Apple HIG 风格
 * 负责功能：
 *   - 模型列表和CRUD
 *   - 模型价格配置
 *   - 模型映射管理
 *   - 平台筛选和状态切换
 * 重要程度：⭐⭐⭐⭐ 重要（模型配置）
-->
<template>
  <div class="models-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">模型管理</h1>
        <p class="page-subtitle">管理 AI 模型配置和模型映射规则</p>
      </div>
    </div>

    <!-- Tab 切换 -->
    <div class="tabs-container">
      <div class="tabs">
        <button :class="['tab', { active: activeTab === 'models' }]" @click="activeTab = 'models'">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="12,2 2,7 12,12 22,7"/>
            <polyline points="2,17 12,22 22,17"/>
            <polyline points="2,12 12,17 22,12"/>
          </svg>
          模型列表
        </button>
        <button :class="['tab', { active: activeTab === 'mappings' }]" @click="activeTab = 'mappings'; loadMappings(); loadMappingCacheStats()">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="22" y1="2" x2="11" y2="13"/>
            <polygon points="22,2 15,22 11,13 2,9"/>
          </svg>
          模型映射
        </button>
      </div>
    </div>

    <!-- 模型列表 Tab -->
    <div v-show="activeTab === 'models'" class="tab-content">
      <div class="table-card">
        <div class="card-toolbar">
          <div class="toolbar-left">
            <select v-model="filterPlatform" @change="loadModels" class="filter-select">
              <option value="">全部平台</option>
              <option value="claude">Claude</option>
              <option value="openai">OpenAI</option>
              <option value="gemini">Gemini</option>
            </select>
          </div>
          <button class="btn btn-primary" @click="showAddDialog">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            添加模型
          </button>
        </div>

        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>模型名称</th>
                <th>平台</th>
                <th class="col-right">上下文</th>
                <th class="col-right">价格 ($/1M)</th>
                <th class="col-center">状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody v-if="!loading">
              <tr v-for="model in filteredModels" :key="model.id">
                <td>
                  <div class="model-name-cell">
                    <span class="model-name">{{ model.name }}</span>
                    <span v-if="model.is_default" class="badge badge-success">默认</span>
                  </div>
                  <div v-if="model.display_name" class="display-name">{{ model.display_name }}</div>
                </td>
                <td>
                  <span :class="['badge', getPlatformBadgeClass(model.platform)]">
                    {{ model.platform }}
                  </span>
                </td>
                <td class="col-right">
                  <span class="mono">{{ formatNumber(model.context_size) }}</span>
                </td>
                <td class="col-right">
                  <div class="price-info">
                    <span>入: ${{ model.input_price }}</span>
                    <span>出: ${{ model.output_price }}</span>
                  </div>
                </td>
                <td class="col-center">
                  <label class="toggle-switch">
                    <input type="checkbox" v-model="model.enabled" @change="toggleEnabled(model)" />
                    <span class="toggle-slider"></span>
                  </label>
                </td>
                <td>
                  <div class="action-group">
                    <button class="action-btn" @click="editModel(model)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn danger" @click="confirmDeleteModel(model)" title="删除">
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

          <div v-if="!loading && filteredModels.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <polygon points="12,2 2,7 12,12 22,7"/>
              <polyline points="2,17 12,22 22,17"/>
              <polyline points="2,12 12,17 22,12"/>
            </svg>
            <span>暂无模型</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 模型映射 Tab -->
    <div v-show="activeTab === 'mappings'" class="tab-content">
      <!-- 说明 -->
      <div class="info-banner">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="16" x2="12" y2="12"/>
          <line x1="12" y1="8" x2="12.01" y2="8"/>
        </svg>
        <span>模型映射允许将客户端请求的模型名自动转换为实际的模型名。例如: claude-3-5-sonnet → claude-sonnet-4-20250514</span>
      </div>

      <div class="table-card">
        <div class="card-toolbar">
          <span class="toolbar-title">映射规则</span>
          <div class="toolbar-actions">
            <button class="btn btn-outline" @click="refreshMappingCache">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23,4 23,10 17,10"/>
                <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
              </svg>
              刷新缓存
            </button>
            <button class="btn btn-primary" @click="showAddMappingDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="12" y1="5" x2="12" y2="19"/>
                <line x1="5" y1="12" x2="19" y2="12"/>
              </svg>
              添加映射
            </button>
          </div>
        </div>

        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>源模型</th>
                <th class="col-center"></th>
                <th>目标模型</th>
                <th class="col-center">优先级</th>
                <th>描述</th>
                <th class="col-center">状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody v-if="!mappingLoading">
              <tr v-for="mapping in mappings" :key="mapping.id">
                <td>
                  <code class="model-code">{{ mapping.source_model }}</code>
                </td>
                <td class="col-center">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="arrow-icon">
                    <line x1="5" y1="12" x2="19" y2="12"/>
                    <polyline points="12,5 19,12 12,19"/>
                  </svg>
                </td>
                <td>
                  <code class="model-code">{{ mapping.target_model }}</code>
                </td>
                <td class="col-center">
                  <span class="badge badge-info">{{ mapping.priority }}</span>
                </td>
                <td>
                  <span class="description">{{ mapping.description || '-' }}</span>
                </td>
                <td class="col-center">
                  <label class="toggle-switch">
                    <input type="checkbox" v-model="mapping.enabled" @change="toggleMappingEnabled(mapping)" />
                    <span class="toggle-slider"></span>
                  </label>
                </td>
                <td>
                  <div class="action-group">
                    <button class="action-btn" @click="editMapping(mapping)" title="编辑">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                      </svg>
                    </button>
                    <button class="action-btn danger" @click="confirmDeleteMapping(mapping)" title="删除">
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

          <div v-if="mappingLoading" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载中...</span>
          </div>

          <div v-if="!mappingLoading && mappings.length === 0" class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <line x1="22" y1="2" x2="11" y2="13"/>
              <polygon points="22,2 15,22 11,13 2,9"/>
            </svg>
            <span>暂无映射规则</span>
          </div>
        </div>

        <!-- 缓存统计 -->
        <div v-if="mappingCacheStats" class="cache-stats">
          <div class="cache-stats-header">
            <strong>当前生效的映射 ({{ mappingCacheStats.count }} 条)</strong>
          </div>
          <div v-if="mappingCacheStats.mappings && Object.keys(mappingCacheStats.mappings).length > 0" class="cache-stats-content">
            <span
              v-for="(target, source) in mappingCacheStats.mappings"
              :key="source"
              class="mapping-tag"
            >
              {{ source }} → {{ target }}
            </span>
          </div>
          <div v-else class="no-mappings">暂无生效的映射</div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑模型弹窗 -->
    <Teleport to="body">
      <div v-if="dialogVisible" class="modal-overlay" @click.self="dialogVisible = false">
        <div class="modal modal-lg">
          <div class="modal-header">
            <h2>{{ isEdit ? '编辑模型' : '添加模型' }}</h2>
            <button class="modal-close" @click="dialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">模型名称 <span class="required">*</span></label>
                <input v-model="form.name" type="text" class="form-input" placeholder="如: claude-3-5-sonnet" />
              </div>
              <div class="form-group">
                <label class="form-label">显示名称</label>
                <input v-model="form.display_name" type="text" class="form-input" placeholder="如: Claude 3.5 Sonnet" />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">平台 <span class="required">*</span></label>
                <select v-model="form.platform" class="form-select">
                  <option value="claude">Claude</option>
                  <option value="openai">OpenAI</option>
                  <option value="gemini">Gemini</option>
                </select>
                <p class="form-tip">决定使用哪个 API 处理请求</p>
              </div>
              <div class="form-group">
                <label class="form-label">提供商</label>
                <input v-model="form.provider" type="text" class="form-input" placeholder="如: anthropic" />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">上下文长度</label>
                <input v-model.number="form.context_size" type="number" class="form-input" />
              </div>
              <div class="form-group">
                <label class="form-label">最大输出</label>
                <input v-model.number="form.max_output" type="number" class="form-input" />
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">输入价格</label>
                <input v-model.number="form.input_price" type="number" step="0.0001" class="form-input" />
                <p class="form-tip">$/1M tokens</p>
              </div>
              <div class="form-group">
                <label class="form-label">输出价格</label>
                <input v-model.number="form.output_price" type="number" step="0.0001" class="form-input" />
                <p class="form-tip">$/1M tokens</p>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">缓存写入价格</label>
                <input v-model.number="form.cache_create_price" type="number" step="0.0001" class="form-input" />
                <p class="form-tip">$/1M tokens (可选)</p>
              </div>
              <div class="form-group">
                <label class="form-label">缓存读取价格</label>
                <input v-model.number="form.cache_read_price" type="number" step="0.0001" class="form-input" />
                <p class="form-tip">$/1M tokens (可选)</p>
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">别名</label>
              <input v-model="form.aliases" type="text" class="form-input" placeholder="多个别名用逗号分隔" />
            </div>
            <div class="form-group">
              <label class="form-label">描述</label>
              <textarea v-model="form.description" class="form-textarea" rows="2"></textarea>
            </div>
            <div class="form-row three">
              <div class="form-group">
                <label class="form-label">排序</label>
                <input v-model.number="form.sort_order" type="number" class="form-input" />
              </div>
              <div class="form-group">
                <label class="form-label">启用</label>
                <label class="toggle-switch large">
                  <input type="checkbox" v-model="form.enabled" />
                  <span class="toggle-slider"></span>
                </label>
              </div>
              <div class="form-group">
                <label class="form-label">默认</label>
                <label class="toggle-switch large">
                  <input type="checkbox" v-model="form.is_default" />
                  <span class="toggle-slider"></span>
                </label>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="dialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="submitting" @click="submitForm">
              <span v-if="submitting" class="btn-loading"></span>
              {{ submitting ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 添加/编辑映射弹窗 -->
    <Teleport to="body">
      <div v-if="mappingDialogVisible" class="modal-overlay" @click.self="mappingDialogVisible = false">
        <div class="modal">
          <div class="modal-header">
            <h2>{{ isMappingEdit ? '编辑映射' : '添加映射' }}</h2>
            <button class="modal-close" @click="mappingDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">源模型 <span class="required">*</span></label>
              <input v-model="mappingForm.source_model" type="text" class="form-input" placeholder="客户端请求的模型名" list="source-models" />
              <datalist id="source-models">
                <option v-for="m in models" :key="m.id" :value="m.name">{{ m.display_name }}</option>
              </datalist>
              <p class="form-tip">客户端请求时使用的模型名称（可输入自定义名称）</p>
            </div>
            <div class="form-group">
              <label class="form-label">目标模型 <span class="required">*</span></label>
              <select v-model="mappingForm.target_model" class="form-select">
                <option value="" disabled>请选择目标模型</option>
                <option v-for="m in enabledModels" :key="m.id" :value="m.name">
                  {{ m.name }} {{ m.display_name ? `(${m.display_name})` : '' }}
                </option>
              </select>
              <p class="form-tip">实际发送到上游的模型名称</p>
            </div>
            <div class="form-group">
              <label class="form-label">优先级</label>
              <input v-model.number="mappingForm.priority" type="number" min="0" max="100" class="form-input" />
              <p class="form-tip">数值越大优先级越高</p>
            </div>
            <div class="form-group">
              <label class="form-label">描述</label>
              <textarea v-model="mappingForm.description" class="form-textarea" rows="2" placeholder="映射说明（可选）"></textarea>
            </div>
            <div class="form-group">
              <label class="form-label">启用</label>
              <label class="toggle-switch large">
                <input type="checkbox" v-model="mappingForm.enabled" />
                <span class="toggle-slider"></span>
              </label>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="mappingDialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="mappingSubmitting" @click="submitMappingForm">
              <span v-if="mappingSubmitting" class="btn-loading"></span>
              {{ mappingSubmitting ? '保存中...' : '保存' }}
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
            <p class="delete-message">{{ deleteMessage }}</p>
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

// 公共状态
const activeTab = ref('models')

// 模型列表状态
const loading = ref(false)
const models = ref([])
const filterPlatform = ref('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)

const defaultForm = {
  name: '',
  display_name: '',
  platform: 'claude',
  provider: '',
  description: '',
  category: 'chat',
  context_size: 200000,
  max_output: 8192,
  input_price: 0,
  output_price: 0,
  cache_create_price: 0,
  cache_read_price: 0,
  enabled: true,
  is_default: false,
  sort_order: 0,
  aliases: ''
}

const form = reactive({ ...defaultForm })

const filteredModels = computed(() => {
  if (!filterPlatform.value) return models.value
  return models.value.filter(m => m.platform === filterPlatform.value)
})

const enabledModels = computed(() => {
  return models.value.filter(m => m.enabled)
})

// 模型映射状态
const mappingLoading = ref(false)
const mappings = ref([])
const mappingDialogVisible = ref(false)
const isMappingEdit = ref(false)
const mappingSubmitting = ref(false)
const mappingCacheStats = ref(null)

const defaultMappingForm = {
  source_model: '',
  target_model: '',
  priority: 0,
  description: '',
  enabled: true
}

const mappingForm = reactive({ ...defaultMappingForm })

// 删除相关
const deleteDialogVisible = ref(false)
const deleteMessage = ref('')
const deleteTarget = ref(null)
const deleteType = ref('') // 'model' or 'mapping'
const deleting = ref(false)

function formatNumber(num) {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(0) + 'K'
  return num
}

function getPlatformBadgeClass(platform) {
  const map = { claude: 'badge-purple', openai: 'badge-success', gemini: 'badge-info' }
  return map[platform] || 'badge-secondary'
}

// 模型相关方法
async function loadModels() {
  loading.value = true
  try {
    const res = await api.getModels()
    models.value = res.data || []
  } catch (e) {
    ElMessage.error('加载模型列表失败')
  } finally {
    loading.value = false
  }
}

function showAddDialog() {
  isEdit.value = false
  Object.assign(form, { ...defaultForm })
  dialogVisible.value = true
}

function editModel(row) {
  isEdit.value = true
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

async function submitForm() {
  if (!form.name || !form.platform) {
    ElMessage.warning('请填写必填项')
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await api.updateModel(form.id, form)
      ElMessage.success('更新成功')
    } else {
      await api.createModel(form)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    loadModels()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function toggleEnabled(row) {
  try {
    await api.toggleModel(row.id)
  } catch (e) {
    row.enabled = !row.enabled
    ElMessage.error('切换状态失败')
  }
}

function confirmDeleteModel(row) {
  deleteTarget.value = row
  deleteType.value = 'model'
  deleteMessage.value = `确定删除模型 "${row.name}" 吗？`
  deleteDialogVisible.value = true
}

// 映射相关方法
async function loadMappings() {
  mappingLoading.value = true
  try {
    const res = await api.getModelMappings()
    mappings.value = res.data?.mappings || []
  } catch (e) {
    ElMessage.error('加载映射列表失败')
  } finally {
    mappingLoading.value = false
  }
}

async function loadMappingCacheStats() {
  try {
    const res = await api.getModelMappingCacheStats()
    mappingCacheStats.value = res.data || null
  } catch (e) {
    console.error('加载缓存统计失败:', e)
  }
}

function showAddMappingDialog() {
  isMappingEdit.value = false
  Object.assign(mappingForm, { ...defaultMappingForm })
  mappingDialogVisible.value = true
}

function editMapping(row) {
  isMappingEdit.value = true
  Object.assign(mappingForm, { ...row })
  mappingDialogVisible.value = true
}

async function submitMappingForm() {
  if (!mappingForm.source_model || !mappingForm.target_model) {
    ElMessage.warning('请填写必填项')
    return
  }

  mappingSubmitting.value = true
  try {
    if (isMappingEdit.value) {
      await api.updateModelMapping(mappingForm.id, mappingForm)
      ElMessage.success('更新成功')
    } else {
      await api.createModelMapping(mappingForm)
      ElMessage.success('创建成功')
    }
    mappingDialogVisible.value = false
    loadMappings()
    loadMappingCacheStats()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    mappingSubmitting.value = false
  }
}

async function toggleMappingEnabled(row) {
  try {
    await api.toggleModelMapping(row.id)
    loadMappingCacheStats()
  } catch (e) {
    row.enabled = !row.enabled
    ElMessage.error('切换状态失败')
  }
}

function confirmDeleteMapping(row) {
  deleteTarget.value = row
  deleteType.value = 'mapping'
  deleteMessage.value = `确定删除映射 "${row.source_model}" 吗？`
  deleteDialogVisible.value = true
}

async function refreshMappingCache() {
  try {
    await api.refreshModelMappingCache()
    ElMessage.success('缓存已刷新')
    loadMappingCacheStats()
  } catch (e) {
    ElMessage.error('刷新缓存失败')
  }
}

// 删除处理
async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    if (deleteType.value === 'model') {
      await api.deleteModel(deleteTarget.value.id)
      loadModels()
    } else {
      await api.deleteModelMapping(deleteTarget.value.id)
      loadMappings()
      loadMappingCacheStats()
    }
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    deleteTarget.value = null
  } catch (e) {
    ElMessage.error('删除失败')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadModels()
  loadMappings()
  loadMappingCacheStats()
})

watch(activeTab, (newTab) => {
  if (newTab === 'mappings') {
    loadMappings()
    loadMappingCacheStats()
  }
})
</script>

<style scoped>
.models-page {
  max-width: 1200px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  margin-bottom: var(--apple-spacing-xl);
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

/* 信息横幅 */
.info-banner {
  display: flex;
  align-items: flex-start;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-md);
  background: var(--apple-blue-light);
  border-radius: var(--apple-radius-md);
  margin-bottom: var(--apple-spacing-lg);
}

.info-banner svg {
  width: 20px;
  height: 20px;
  color: var(--apple-blue);
  flex-shrink: 0;
}

.info-banner span {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  line-height: 1.5;
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

.toolbar-actions {
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

.col-center { text-align: center; }
.col-right { text-align: right; }

.model-name-cell {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.model-name {
  font-family: var(--apple-font-mono);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.display-name {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: 2px;
}

.model-code {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-primary);
  background: var(--apple-bg-tertiary);
  padding: 2px 8px;
  border-radius: var(--apple-radius-xs);
}

.arrow-icon {
  width: 16px;
  height: 16px;
  color: var(--apple-text-tertiary);
}

.description {
  color: var(--apple-text-tertiary);
  font-size: var(--apple-text-sm);
}

.mono {
  font-family: var(--apple-font-mono);
  color: var(--apple-text-secondary);
}

.price-info {
  display: flex;
  flex-direction: column;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

/* 徽章 */
.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.badge-success { background: var(--apple-green-light); color: var(--apple-green); }
.badge-purple { background: #eef2ff; color: #667eea; }
.badge-info { background: var(--apple-fill-tertiary); color: var(--apple-text-secondary); }
.badge-secondary { background: var(--apple-fill-quaternary); color: var(--apple-text-tertiary); }

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

/* 缓存统计 */
.cache-stats {
  padding: var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
}

.cache-stats-header {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-bottom: var(--apple-spacing-md);
}

.cache-stats-content {
  display: flex;
  flex-wrap: wrap;
  gap: var(--apple-spacing-xs);
}

.mapping-tag {
  padding: var(--apple-spacing-xxs) var(--apple-spacing-sm);
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-xs);
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.no-mappings {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
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
.modal.modal-lg { max-width: 650px; }

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

.form-row.three {
  grid-template-columns: repeat(3, 1fr);
}

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

.form-input, .form-select, .form-textarea {
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

.form-input:focus, .form-select:focus, .form-textarea:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.form-tip {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: var(--apple-spacing-xxs);
}

.delete-message {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  text-align: center;
  margin: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .form-row, .form-row.three {
    grid-template-columns: 1fr;
  }

  .card-toolbar {
    flex-direction: column;
    gap: var(--apple-spacing-sm);
    align-items: stretch;
  }

  .toolbar-actions {
    justify-content: flex-end;
  }
}
</style>
