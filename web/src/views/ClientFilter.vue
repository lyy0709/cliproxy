<!--
 * 文件作用：客户端过滤页面 - Apple HIG 风格
 * 负责功能：
 *   - 过滤功能开关
 *   - 允许客户端类型选择
 *   - 验证模式切换（简单/严格）
 *   - 请求验证测试
 * 重要程度：⭐⭐⭐ 一般（安全配置）
-->
<template>
  <div class="client-filter-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">客户端过滤</h1>
        <p class="page-subtitle">限制只允许指定客户端访问 API</p>
      </div>
    </div>

    <!-- 主开关卡片 -->
    <div class="switch-card">
      <div class="switch-content">
        <div class="switch-icon" :class="{ active: config.filter_enabled }">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
          </svg>
        </div>
        <div class="switch-info">
          <h2>客户端过滤</h2>
          <p>{{ config.filter_enabled ? '过滤已启用，只有允许的客户端可以访问' : '过滤已关闭，所有请求都将被允许' }}</p>
        </div>
        <label class="toggle-switch large">
          <input type="checkbox" v-model="config.filter_enabled" @change="saveConfig" :disabled="configLoading" />
          <span class="toggle-slider"></span>
        </label>
      </div>
    </div>

    <!-- 过滤关闭提示 -->
    <div v-if="!config.filter_enabled" class="info-banner">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"/>
        <line x1="12" y1="16" x2="12" y2="12"/>
        <line x1="12" y1="8" x2="12.01" y2="8"/>
      </svg>
      <span>过滤功能已关闭，所有请求都将被允许通过</span>
    </div>

    <template v-if="config.filter_enabled">
      <!-- 允许的客户端 -->
      <div class="section-card">
        <div class="section-header">
          <h3>允许的客户端</h3>
          <span class="section-badge">{{ clientTypes.filter(c => c.enabled).length }}/{{ clientTypes.length }} 已启用</span>
        </div>
        <div class="client-grid">
          <div
            v-for="ct in clientTypes"
            :key="ct.id"
            :class="['client-card', { enabled: ct.enabled }]"
            @click="toggleClient(ct)"
          >
            <div class="client-icon">{{ ct.icon }}</div>
            <div class="client-info">
              <span class="client-name">{{ ct.name }}</span>
              <span class="client-desc">{{ ct.description }}</span>
            </div>
            <div v-if="ct.enabled" class="check-mark">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <polyline points="20,6 9,17 4,12"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- 验证模式 -->
      <div class="section-card">
        <div class="section-header">
          <h3>验证模式</h3>
        </div>
        <div class="mode-selector">
          <label :class="['mode-option', { active: config.filter_mode === 'simple' }]" @click="config.filter_mode = 'simple'; saveConfig()">
            <div class="mode-radio">
              <span class="radio-dot"></span>
            </div>
            <div class="mode-content">
              <span class="mode-title">简单模式</span>
              <span class="mode-desc">宽松验证，User-Agent 匹配即可</span>
            </div>
          </label>
          <label :class="['mode-option', { active: config.filter_mode === 'strict' }]" @click="config.filter_mode = 'strict'; saveConfig()">
            <div class="mode-radio">
              <span class="radio-dot"></span>
            </div>
            <div class="mode-content">
              <span class="mode-title">严格模式</span>
              <span class="mode-desc">完整验证，包括请求头和请求体</span>
            </div>
          </label>
        </div>
        <div class="mode-detail">
          <template v-if="config.filter_mode === 'simple'">
            <p><strong>简单模式</strong>：User-Agent 以 <code>claude-cli/版本号</code> 开头即可通过验证</p>
          </template>
          <template v-else>
            <p><strong>严格模式</strong>：User-Agent 必须是完整格式 <code>claude-cli/版本 (external, cli)</code>，并验证请求头</p>
          </template>
        </div>
      </div>

      <!-- 验证测试 -->
      <div class="section-card">
        <div class="section-header collapsible" @click="testExpanded = !testExpanded">
          <h3>验证测试</h3>
          <svg :class="['expand-icon', { expanded: testExpanded }]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6,9 12,15 18,9"/>
          </svg>
        </div>
        <div v-if="testExpanded" class="test-section">
          <div class="form-group">
            <label class="form-label">User-Agent</label>
            <input
              v-model="testData.user_agent"
              type="text"
              class="form-input"
              placeholder="claude-cli/1.0.15 (external, cli)"
            />
          </div>
          <div class="form-group">
            <label class="form-label">请求头 (JSON)</label>
            <textarea
              v-model="testData.headersJson"
              class="form-textarea"
              rows="2"
              placeholder='{"x-app": "claude-code", "anthropic-version": "2023-06-01"}'
            ></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">请求体 (JSON)</label>
            <textarea
              v-model="testData.bodyJson"
              class="form-textarea"
              rows="3"
              placeholder='{"system": [{"text": "You are Claude Code..."}]}'
            ></textarea>
          </div>
          <div class="test-actions">
            <button class="btn btn-secondary" @click="fillTestExample">填入示例</button>
            <button class="btn btn-primary" @click="runTest" :disabled="testing">
              <span v-if="testing" class="btn-loading"></span>
              {{ testing ? '测试中...' : '测试验证' }}
            </button>
          </div>

          <div v-if="testResult" :class="['test-result', testResult.allowed ? 'success' : 'error']">
            <div class="result-icon">
              <svg v-if="testResult.allowed" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 11.08V12a10 10 0 11-5.93-9.14"/>
                <polyline points="22,4 12,14.01 9,11.01"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
            </div>
            <div class="result-content">
              <span class="result-title">{{ testResult.allowed ? '验证通过' : '验证失败' }}</span>
              <span class="result-client">客户端: {{ testResult.client_name || '未识别' }}</span>
              <span v-if="testResult.details?.reason" class="result-reason">原因: {{ testResult.details.reason }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const config = reactive({
  filter_enabled: false,
  filter_mode: 'simple'
})
const configLoading = ref(false)
const clientTypes = ref([])
const testExpanded = ref(false)
const testing = ref(false)

const testData = reactive({
  user_agent: '',
  headersJson: '',
  bodyJson: ''
})
const testResult = ref(null)

async function loadData() {
  try {
    const [configRes, typesRes] = await Promise.all([
      api.getClientFilterConfig(),
      api.getClientTypes()
    ])
    Object.assign(config, configRes.data)
    if (!config.filter_mode) {
      config.filter_mode = 'simple'
    }
    clientTypes.value = typesRes.data || []
  } catch (e) {
    console.error('Failed to load data:', e)
  }
}

async function saveConfig() {
  configLoading.value = true
  try {
    await api.updateClientFilterConfig(config)
    ElMessage.success('配置已保存')
  } catch (e) {
    await loadData()
    ElMessage.error('保存失败')
  } finally {
    configLoading.value = false
  }
}

async function toggleClient(ct) {
  ct.enabled = !ct.enabled
  try {
    await api.toggleClientType(ct.id)
    ElMessage.success(`${ct.name} 已${ct.enabled ? '启用' : '禁用'}`)
  } catch (e) {
    ct.enabled = !ct.enabled
    ElMessage.error('操作失败')
  }
}

function fillTestExample() {
  testData.user_agent = 'claude-cli/1.0.15 (external, cli)'
  testData.headersJson = JSON.stringify({
    'x-app': 'claude-code',
    'anthropic-version': '2023-06-01',
    'anthropic-beta': 'oauth-2025-04-20',
    'x-stainless-os': 'Linux'
  }, null, 2)
  testData.bodyJson = JSON.stringify({
    system: [{ text: "You are Claude Code, Anthropic's official CLI for Claude." }],
    metadata: {
      user_id: 'user_d98385411c93cd074b2cefd5c9831fe77f24a53e4ecdcd1f830bba586fe62cb9_account__session_17cf0fd3-d51b-4b59-977d-b899dafb3022'
    }
  }, null, 2)
}

async function runTest() {
  testing.value = true
  testResult.value = null
  try {
    let headers = {}, body = {}
    try {
      if (testData.headersJson) headers = JSON.parse(testData.headersJson)
    } catch {
      ElMessage.error('请求头 JSON 格式错误')
      testing.value = false
      return
    }
    try {
      if (testData.bodyJson) body = JSON.parse(testData.bodyJson)
    } catch {
      ElMessage.error('请求体 JSON 格式错误')
      testing.value = false
      return
    }

    const res = await api.testClientFilter({
      user_agent: testData.user_agent,
      path: '/v1/messages',
      headers,
      body
    })
    testResult.value = res.data
  } catch (e) {
    ElMessage.error('测试请求失败')
  } finally {
    testing.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.client-filter-page {
  max-width: 800px;
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

/* 主开关卡片 */
.switch-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-xl);
  margin-bottom: var(--apple-spacing-lg);
}

.switch-content {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-lg);
}

.switch-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--apple-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-tertiary);
  transition: all var(--apple-duration-normal);
}

.switch-icon.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.switch-icon svg {
  width: 28px;
  height: 28px;
  color: var(--apple-text-secondary);
  transition: color var(--apple-duration-normal);
}

.switch-icon.active svg {
  color: white;
}

.switch-info {
  flex: 1;
}

.switch-info h2 {
  margin: 0 0 var(--apple-spacing-xxs) 0;
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.switch-info p {
  margin: 0;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

/* 开关 */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 51px;
  height: 31px;
  flex-shrink: 0;
}

.toggle-switch.large {
  width: 51px;
  height: 31px;
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
  border-radius: 31px;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 27px;
  width: 27px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.toggle-switch input:checked + .toggle-slider {
  background-color: var(--apple-green);
}

.toggle-switch input:checked + .toggle-slider:before {
  transform: translateX(20px);
}

.toggle-switch input:disabled + .toggle-slider {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 信息横幅 */
.info-banner {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
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
  color: var(--apple-blue);
}

/* 区块卡片 */
.section-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  margin-bottom: var(--apple-spacing-lg);
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.section-header.collapsible {
  cursor: pointer;
  user-select: none;
}

.section-header.collapsible:hover {
  background: var(--apple-bg-secondary);
}

.section-header h3 {
  margin: 0;
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.section-badge {
  padding: 4px 10px;
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.expand-icon {
  width: 20px;
  height: 20px;
  color: var(--apple-text-tertiary);
  transition: transform var(--apple-duration-fast);
}

.expand-icon.expanded {
  transform: rotate(180deg);
}

/* 客户端网格 */
.client-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-lg);
}

.client-card {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-md);
  border: 2px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: all var(--apple-duration-fast);
  position: relative;
}

.client-card:hover {
  border-color: var(--apple-blue);
}

.client-card.enabled {
  border-color: var(--apple-green);
  background: var(--apple-green-light);
}

.client-icon {
  font-size: 28px;
  flex-shrink: 0;
}

.client-info {
  flex: 1;
  min-width: 0;
}

.client-name {
  display: block;
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.client-desc {
  display: block;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.check-mark {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 20px;
  height: 20px;
  background: var(--apple-green);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.check-mark svg {
  width: 12px;
  height: 12px;
  color: white;
}

/* 模式选择器 */
.mode-selector {
  display: flex;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-lg);
  padding-bottom: 0;
}

.mode-option {
  flex: 1;
  display: flex;
  align-items: flex-start;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-md);
  border: 2px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: all var(--apple-duration-fast);
}

.mode-option:hover {
  border-color: var(--apple-blue);
}

.mode-option.active {
  border-color: var(--apple-blue);
  background: var(--apple-blue-light);
}

.mode-radio {
  width: 20px;
  height: 20px;
  border: 2px solid var(--apple-separator-opaque);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.mode-option.active .mode-radio {
  border-color: var(--apple-blue);
}

.radio-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: transparent;
  transition: background var(--apple-duration-fast);
}

.mode-option.active .radio-dot {
  background: var(--apple-blue);
}

.mode-content {
  flex: 1;
}

.mode-title {
  display: block;
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.mode-desc {
  display: block;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: 2px;
}

.mode-detail {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  margin: var(--apple-spacing-md);
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-sm);
}

.mode-detail p {
  margin: 0;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.mode-detail code {
  background: var(--apple-fill-tertiary);
  padding: 2px 6px;
  border-radius: var(--apple-radius-xs);
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
}

/* 测试区域 */
.test-section {
  padding: var(--apple-spacing-lg);
}

/* 表单 */
.form-group {
  margin-bottom: var(--apple-spacing-md);
}

.form-label {
  display: block;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.form-input, .form-textarea {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
}

.form-textarea {
  resize: vertical;
  min-height: 60px;
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.test-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--apple-spacing-sm);
  margin-top: var(--apple-spacing-lg);
}

/* 测试结果 */
.test-result {
  display: flex;
  align-items: flex-start;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-md);
  border-radius: var(--apple-radius-md);
  margin-top: var(--apple-spacing-lg);
}

.test-result.success {
  background: var(--apple-green-light);
}

.test-result.error {
  background: var(--apple-red-light);
}

.result-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.test-result.success .result-icon {
  background: var(--apple-green);
}

.test-result.error .result-icon {
  background: var(--apple-red);
}

.result-icon svg {
  width: 20px;
  height: 20px;
  color: white;
}

.result-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.result-title {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.result-client, .result-reason {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
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
  transition: all var(--apple-duration-fast);
  cursor: pointer;
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

.btn-loading {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 响应式 */
@media (max-width: 768px) {
  .switch-content {
    flex-wrap: wrap;
  }

  .client-grid {
    grid-template-columns: 1fr;
  }

  .mode-selector {
    flex-direction: column;
  }
}
</style>
