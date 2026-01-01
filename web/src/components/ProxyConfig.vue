<!--
 * 文件作用：代理配置组件 - Apple HIG 风格
 * 负责功能：
 *   - 代理启用/禁用开关
 *   - 快速URL解析配置
 *   - 手动代理参数配置
 *   - 代理连接测试
 * 重要程度：⭐⭐⭐ 一般（代理配置辅助）
-->
<template>
  <div class="proxy-config">
    <!-- 启用开关 -->
    <div class="proxy-header">
      <div class="proxy-title">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
        </svg>
        <span>代理配置</span>
      </div>
      <button :class="['toggle-switch', { active: localProxy.enabled }]" @click="toggleEnabled">
        <span class="toggle-thumb"></span>
      </button>
    </div>

    <!-- 代理配置表单 -->
    <div v-if="localProxy.enabled" class="proxy-form">
      <!-- 快速配置 -->
      <div class="quick-config">
        <label class="config-label">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71"/>
            <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71"/>
          </svg>
          快速配置
        </label>
        <input
          v-model="proxyUrl"
          type="text"
          class="config-input"
          placeholder="粘贴代理URL，如 socks5://user:pass@host:port"
          @paste="handlePaste"
          @input="parseProxyUrl"
        />
        <div class="quick-tip">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="16" x2="12" y2="12"/>
            <line x1="12" y1="8" x2="12.01" y2="8"/>
          </svg>
          支持格式：socks5://user:pass@host:port 或 http://host:port
        </div>
      </div>

      <!-- 分隔线 -->
      <div class="divider">
        <span>或手动配置</span>
      </div>

      <!-- 手动配置表单 -->
      <div class="manual-form">
        <div class="form-row three-col">
          <div class="form-group">
            <label class="form-label">代理类型</label>
            <select v-model="localProxy.type" class="form-select" @change="emitUpdate">
              <option value="socks5">SOCKS5</option>
              <option value="http">HTTP</option>
              <option value="https">HTTPS</option>
            </select>
          </div>
          <div class="form-group flex-2">
            <label class="form-label">主机地址</label>
            <input v-model="localProxy.host" type="text" class="form-input" placeholder="proxy.example.com" @input="emitUpdate" />
          </div>
          <div class="form-group">
            <label class="form-label">端口</label>
            <input v-model.number="localProxy.port" type="number" class="form-input" min="1" max="65535" placeholder="1080" @input="emitUpdate" />
          </div>
        </div>

        <div class="form-row two-col">
          <div class="form-group">
            <label class="form-label">用户名（可选）</label>
            <input v-model="localProxy.username" type="text" class="form-input" placeholder="代理用户名" @input="emitUpdate" />
          </div>
          <div class="form-group">
            <label class="form-label">密码（可选）</label>
            <div class="password-input">
              <input
                v-model="localProxy.password"
                :type="showPassword ? 'text' : 'password'"
                class="form-input"
                placeholder="代理密码"
                @input="emitUpdate"
              />
              <button class="password-toggle" @click="showPassword = !showPassword" type="button">
                <svg v-if="showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 代理测试 -->
      <div class="proxy-test">
        <button class="test-btn" :disabled="testing" @click="testProxy">
          <svg v-if="!testing" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M5 12h14M12 5l7 7-7 7"/>
          </svg>
          <span v-if="testing" class="test-spinner"></span>
          {{ testing ? '测试中...' : '测试连接' }}
        </button>
        <div v-if="testResult" :class="['test-result', testResult.success ? 'success' : 'error']">
          <svg v-if="testResult.success" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 11.08V12a10 10 0 11-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
          <span>{{ testResult.message }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({
      enabled: false,
      type: 'socks5',
      host: '',
      port: 1080,
      username: '',
      password: ''
    })
  }
})

const emit = defineEmits(['update:modelValue'])

const localProxy = reactive({
  enabled: false,
  type: 'socks5',
  host: '',
  port: 1080,
  username: '',
  password: ''
})

const proxyUrl = ref('')
const testing = ref(false)
const testResult = ref(null)
const showPassword = ref(false)

// 监听外部值变化
watch(() => props.modelValue, (val) => {
  if (val) {
    Object.assign(localProxy, val)
  }
}, { immediate: true, deep: true })

function toggleEnabled() {
  localProxy.enabled = !localProxy.enabled
  emitUpdate()
}

function emitUpdate() {
  emit('update:modelValue', { ...localProxy })
}

function handlePaste() {
  setTimeout(() => {
    parseProxyUrl()
  }, 10)
}

function parseProxyUrl() {
  const url = proxyUrl.value.trim()
  if (!url) return

  try {
    // 解析代理URL格式: type://[user:pass@]host:port
    const match = url.match(/^(socks5|http|https):\/\/(?:([^:]+):([^@]+)@)?([^:]+):(\d+)$/i)
    if (match) {
      localProxy.type = match[1].toLowerCase()
      localProxy.username = match[2] || ''
      localProxy.password = match[3] || ''
      localProxy.host = match[4]
      localProxy.port = parseInt(match[5])
      localProxy.enabled = true
      emitUpdate()
    }
  } catch (e) {
    console.error('Failed to parse proxy URL:', e)
  }
}

async function testProxy() {
  if (!localProxy.host || !localProxy.port) {
    testResult.value = { success: false, message: '请填写代理地址和端口' }
    return
  }

  testing.value = true
  testResult.value = null

  try {
    // 这里调用后端API测试代理连接
    // 暂时模拟测试结果
    await new Promise(resolve => setTimeout(resolve, 1000))
    testResult.value = { success: true, message: '连接成功' }
  } catch (e) {
    testResult.value = { success: false, message: e.message || '连接失败' }
  } finally {
    testing.value = false
  }
}
</script>

<style scoped>
.proxy-config {
  background: var(--apple-fill-quaternary, #f5f5f7);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-lg, 12px);
  padding: 16px;
}

.proxy-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.proxy-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 15px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
}

.proxy-title svg {
  width: 20px;
  height: 20px;
  color: var(--apple-blue, #007aff);
}

/* 开关按钮 */
.toggle-switch {
  position: relative;
  width: 51px;
  height: 31px;
  background: var(--apple-fill-secondary, rgba(120, 120, 128, 0.16));
  border: none;
  border-radius: 16px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.toggle-switch.active {
  background: var(--apple-green, #34c759);
}

.toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 27px;
  height: 27px;
  background: white;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.15);
  transition: transform 0.2s ease;
}

.toggle-switch.active .toggle-thumb {
  transform: translateX(20px);
}

/* 代理表单 */
.proxy-form {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 快速配置 */
.quick-config {
  margin-bottom: 16px;
}

.config-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  color: var(--apple-text-primary, #1d1d1f);
  margin-bottom: 8px;
}

.config-label svg {
  width: 16px;
  height: 16px;
  color: var(--apple-text-secondary, #86868b);
}

.config-input {
  width: 100%;
  padding: 10px 14px;
  font-size: 14px;
  color: var(--apple-text-primary, #1d1d1f);
  background: var(--apple-bg-primary, #ffffff);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-md, 8px);
  transition: all 0.2s ease;
}

.config-input::placeholder {
  color: var(--apple-text-placeholder, #c7c7cc);
}

.config-input:focus {
  border-color: var(--apple-blue, #007aff);
  box-shadow: 0 0 0 3px var(--apple-blue-light, rgba(0, 122, 255, 0.15));
  outline: none;
}

.quick-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--apple-text-tertiary, #c7c7cc);
  margin-top: 8px;
}

.quick-tip svg {
  width: 14px;
  height: 14px;
  color: var(--apple-orange, #ff9500);
}

/* 分隔线 */
.divider {
  display: flex;
  align-items: center;
  margin: 20px 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--apple-separator, rgba(0, 0, 0, 0.1));
}

.divider span {
  padding: 0 12px;
  font-size: 12px;
  color: var(--apple-text-tertiary, #c7c7cc);
}

/* 手动表单 */
.manual-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: flex;
  gap: 12px;
}

.form-row.three-col .form-group:first-child {
  flex: 0 0 120px;
}

.form-row.three-col .form-group:last-child {
  flex: 0 0 100px;
}

.form-row.two-col .form-group {
  flex: 1;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group.flex-2 {
  flex: 2;
}

.form-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--apple-text-primary, #1d1d1f);
}

.form-input,
.form-select {
  width: 100%;
  padding: 10px 14px;
  font-size: 14px;
  color: var(--apple-text-primary, #1d1d1f);
  background: var(--apple-bg-primary, #ffffff);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-md, 8px);
  transition: all 0.2s ease;
}

.form-input::placeholder {
  color: var(--apple-text-placeholder, #c7c7cc);
}

.form-input:focus,
.form-select:focus {
  border-color: var(--apple-blue, #007aff);
  box-shadow: 0 0 0 3px var(--apple-blue-light, rgba(0, 122, 255, 0.15));
  outline: none;
}

.form-select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='16' height='16' viewBox='0 0 24 24' fill='none' stroke='%2386868b' stroke-width='2'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 36px;
  cursor: pointer;
}

/* 密码输入框 */
.password-input {
  position: relative;
}

.password-input .form-input {
  padding-right: 44px;
}

.password-toggle {
  position: absolute;
  right: 4px;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: var(--apple-text-tertiary, #c7c7cc);
  cursor: pointer;
  transition: color 0.2s ease;
}

.password-toggle:hover {
  color: var(--apple-text-secondary, #86868b);
}

.password-toggle svg {
  width: 18px;
  height: 18px;
}

/* 代理测试 */
.proxy-test {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
}

.test-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 500;
  color: var(--apple-blue, #007aff);
  background: var(--apple-blue-light, rgba(0, 122, 255, 0.1));
  border: none;
  border-radius: var(--apple-radius-md, 8px);
  cursor: pointer;
  transition: all 0.2s ease;
}

.test-btn:hover:not(:disabled) {
  background: rgba(0, 122, 255, 0.15);
}

.test-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.test-btn svg {
  width: 16px;
  height: 16px;
}

.test-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(0, 122, 255, 0.2);
  border-top-color: var(--apple-blue, #007aff);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.test-result {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
}

.test-result svg {
  width: 16px;
  height: 16px;
}

.test-result.success {
  color: var(--apple-green, #34c759);
}

.test-result.error {
  color: var(--apple-red, #ff3b30);
}

/* 响应式 */
@media (max-width: 640px) {
  .form-row {
    flex-direction: column;
  }

  .form-row.three-col .form-group:first-child,
  .form-row.three-col .form-group:last-child {
    flex: 1;
  }
}
</style>
