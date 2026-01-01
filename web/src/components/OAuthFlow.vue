<!--
 * 文件作用：OAuth授权流程组件 - Apple HIG 风格
 * 负责功能：
 *   - Claude/Gemini/OpenAI OAuth流程
 *   - SessionKey自动授权（支持批量）
 *   - 授权链接生成和Code交换
 *   - 代理配置传递
 * 重要程度：⭐⭐⭐⭐ 重要（OAuth认证）
-->
<template>
  <div class="oauth-flow">
    <!-- Claude OAuth流程 -->
    <div v-if="isClaudePlatform" class="oauth-section claude">
      <div class="oauth-header">
        <div class="platform-icon claude">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2a10 10 0 1010 10A10 10 0 0012 2zm0 18a8 8 0 118-8 8 8 0 01-8 8z"/>
            <circle cx="12" cy="12" r="3"/>
          </svg>
        </div>
        <div class="platform-info">
          <h4>Claude 账户授权</h4>
        </div>
      </div>

      <!-- 授权方式选择 -->
      <div class="auth-method-select">
        <label class="method-label">选择授权方式</label>
        <div class="method-options">
          <label :class="['method-option', { active: authMethod === 'manual' }]">
            <input type="radio" v-model="authMethod" value="manual" @change="onAuthMethodChange" />
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
            <span>手动授权</span>
          </label>
          <label :class="['method-option', { active: authMethod === 'cookie' }]">
            <input type="radio" v-model="authMethod" value="cookie" @change="onAuthMethodChange" />
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
              <circle cx="9" cy="9" r="1" fill="currentColor"/>
              <circle cx="15" cy="9" r="1" fill="currentColor"/>
            </svg>
            <span>SessionKey 自动授权</span>
          </label>
        </div>
      </div>

      <!-- SessionKey 自动授权 -->
      <div v-if="authMethod === 'cookie'" class="cookie-auth">
        <div class="info-banner">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="12" y1="16" x2="12" y2="12"/>
            <line x1="12" y1="8" x2="12.01" y2="8"/>
          </svg>
          <span>使用 claude.ai 的 sessionKey 自动完成 OAuth 授权流程</span>
        </div>

        <div class="form-group">
          <div class="form-label-row">
            <label class="form-label">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z"/>
                <path d="M8 14s1.5 2 4 2 4-2 4-2"/>
                <circle cx="9" cy="9" r="1" fill="currentColor"/>
                <circle cx="15" cy="9" r="1" fill="currentColor"/>
              </svg>
              sessionKey
            </label>
            <span v-if="parsedSessionKeyCount > 1" class="count-badge">{{ parsedSessionKeyCount }} 个</span>
            <button type="button" class="help-btn" @click="showSessionKeyHelp = !showSessionKeyHelp">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <path d="M9.09 9a3 3 0 015.83 1c0 2-3 3-3 3"/>
                <line x1="12" y1="17" x2="12.01" y2="17"/>
              </svg>
            </button>
          </div>
          <textarea
            v-model="sessionKey"
            class="form-textarea"
            rows="3"
            placeholder="每行一个 sessionKey，例如：&#10;sk-ant-sid01-xxxxx...&#10;sk-ant-sid01-yyyyy..."
          ></textarea>
          <div v-if="parsedSessionKeyCount > 1" class="session-key-tip">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="16" x2="12" y2="12"/>
              <line x1="12" y1="8" x2="12.01" y2="8"/>
            </svg>
            将批量创建 {{ parsedSessionKeyCount }} 个账户
          </div>
        </div>

        <!-- SessionKey 帮助说明 -->
        <transition name="slide">
          <div v-if="showSessionKeyHelp" class="help-section">
            <h5>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="9" y1="18" x2="15" y2="18"/>
                <line x1="10" y1="22" x2="14" y2="22"/>
                <path d="M15.09 14c.18-.98.65-1.74 1.41-2.5A4.65 4.65 0 0018 8 6 6 0 006 8c0 1 .23 2.23 1.5 3.5A4.61 4.61 0 018.91 14"/>
              </svg>
              如何获取 sessionKey
            </h5>
            <ol>
              <li>在浏览器中登录 <strong>claude.ai</strong></li>
              <li>按 <kbd>F12</kbd> 打开开发者工具</li>
              <li>切换到 <strong>Application</strong>（应用）标签页</li>
              <li>在左侧找到 <strong>Cookies</strong> → <strong>https://claude.ai</strong></li>
              <li>找到键为 <strong>sessionKey</strong> 的那一行</li>
              <li>复制其 <strong>Value</strong>（值）列的内容</li>
            </ol>
            <p class="help-tip">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="16" x2="12" y2="12"/>
                <line x1="12" y1="8" x2="12.01" y2="8"/>
              </svg>
              sessionKey 通常以 <code>sk-ant-sid01-</code> 开头
            </p>
          </div>
        </transition>

        <!-- 错误信息 -->
        <div v-if="cookieAuthError" class="error-banner">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <line x1="15" y1="9" x2="9" y2="15"/>
            <line x1="9" y1="9" x2="15" y2="15"/>
          </svg>
          <span>{{ cookieAuthError }}</span>
        </div>

        <!-- 授权按钮 -->
        <button
          class="btn btn-primary btn-block"
          :disabled="!sessionKey.trim() || cookieAuthLoading"
          @click="handleCookieAuth"
        >
          <template v-if="cookieAuthLoading && batchProgress.total > 1">
            <span class="btn-spinner"></span>
            正在授权 {{ batchProgress.current }}/{{ batchProgress.total }}...
          </template>
          <template v-else-if="cookieAuthLoading">
            <span class="btn-spinner"></span>
            正在授权...
          </template>
          <template v-else>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
            </svg>
            开始自动授权
          </template>
        </button>
      </div>

      <!-- 手动授权流程 -->
      <div v-else class="manual-auth">
        <p class="auth-desc">请按照以下步骤完成 Claude 账户的授权：</p>

        <div class="auth-steps">
          <!-- 步骤1: 生成授权链接 -->
          <div class="step-card">
            <div class="step-header">
              <div class="step-num">1</div>
              <div class="step-title">点击下方按钮生成授权链接</div>
            </div>
            <div class="step-content">
              <button v-if="!authUrl" class="btn btn-primary" :disabled="loading" @click="generateAuthUrl">
                <svg v-if="!loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71"/>
                  <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71"/>
                </svg>
                <span v-if="loading" class="btn-spinner"></span>
                {{ loading ? '生成中...' : '生成授权链接' }}
              </button>
              <div v-else class="auth-url-container">
                <div class="url-input-group">
                  <input type="text" :value="authUrl" readonly class="url-input" />
                  <button class="copy-btn" @click="copyAuthUrl">
                    <svg v-if="copied" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="check">
                      <polyline points="20 6 9 17 4 12"/>
                    </svg>
                    <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                      <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                    </svg>
                  </button>
                </div>
                <button class="link-btn" @click="regenerateAuthUrl">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="23 4 23 10 17 10"/>
                    <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
                  </svg>
                  重新生成
                </button>
              </div>
            </div>
          </div>

          <!-- 步骤2: 打开链接授权 -->
          <div class="step-card">
            <div class="step-header">
              <div class="step-num">2</div>
              <div class="step-title">在浏览器中打开链接并完成授权</div>
            </div>
            <div class="step-content">
              <p>请在新标签页中打开授权链接，登录您的 Claude 账户并授权。</p>
              <div class="warning-banner">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
                  <line x1="12" y1="9" x2="12" y2="13"/>
                  <line x1="12" y1="17" x2="12.01" y2="17"/>
                </svg>
                <span><strong>注意：</strong>如果您设置了代理，请确保浏览器也使用相同的代理访问授权页面。</span>
              </div>
            </div>
          </div>

          <!-- 步骤3: 输入授权码 -->
          <div class="step-card">
            <div class="step-header">
              <div class="step-num">3</div>
              <div class="step-title">输入 Authorization Code</div>
            </div>
            <div class="step-content">
              <p class="step-desc">授权完成后，页面会显示一个 <strong>Authorization Code</strong>，请将其复制并粘贴到下方输入框：</p>
              <textarea
                v-model="authCode"
                class="form-textarea"
                rows="3"
                placeholder="粘贴从 Claude 页面获取的 Authorization Code..."
              ></textarea>
              <p class="input-tip">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="12" y1="16" x2="12" y2="12"/>
                  <line x1="12" y1="8" x2="12.01" y2="8"/>
                </svg>
                请粘贴从 Claude 页面复制的 Authorization Code
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Gemini OAuth流程 -->
    <div v-else-if="platform === 'gemini'" class="oauth-section gemini">
      <div class="oauth-header">
        <div class="platform-icon gemini">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5z"/>
            <path d="M2 17l10 5 10-5"/>
            <path d="M2 12l10 5 10-5"/>
          </svg>
        </div>
        <div class="platform-info">
          <h4>Gemini 账户授权</h4>
        </div>
      </div>

      <p class="auth-desc">请按照以下步骤完成 Gemini 账户的授权：</p>

      <div class="auth-steps">
        <div class="step-card">
          <div class="step-header">
            <div class="step-num">1</div>
            <div class="step-title">点击下方按钮生成授权链接</div>
          </div>
          <div class="step-content">
            <button v-if="!authUrl" class="btn btn-primary" :disabled="loading" @click="generateAuthUrl">
              <svg v-if="!loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71"/>
                <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71"/>
              </svg>
              <span v-if="loading" class="btn-spinner"></span>
              {{ loading ? '生成中...' : '生成授权链接' }}
            </button>
            <div v-else class="auth-url-container">
              <div class="url-input-group">
                <input type="text" :value="authUrl" readonly class="url-input" />
                <button class="copy-btn" @click="copyAuthUrl">
                  <svg v-if="copied" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="check">
                    <polyline points="20 6 9 17 4 12"/>
                  </svg>
                  <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                  </svg>
                </button>
              </div>
              <button class="link-btn" @click="regenerateAuthUrl">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"/>
                  <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
                </svg>
                重新生成
              </button>
            </div>
          </div>
        </div>

        <div class="step-card">
          <div class="step-header">
            <div class="step-num">2</div>
            <div class="step-title">在浏览器中打开链接并完成授权</div>
          </div>
          <div class="step-content">
            <p>请在新标签页中打开授权链接，登录您的 Google 账户并授权。</p>
          </div>
        </div>

        <div class="step-card">
          <div class="step-header">
            <div class="step-num">3</div>
            <div class="step-title">输入 Authorization Code</div>
          </div>
          <div class="step-content">
            <textarea
              v-model="authCode"
              class="form-textarea"
              rows="3"
              placeholder="粘贴从 Gemini 页面获取的 Authorization Code..."
            ></textarea>
          </div>
        </div>
      </div>
    </div>

    <!-- OpenAI OAuth流程 -->
    <div v-else-if="platform === 'openai' || platform === 'openai-responses'" class="oauth-section openai">
      <div class="oauth-header">
        <div class="platform-icon openai">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2a10 10 0 1010 10A10 10 0 0012 2zm0 18a8 8 0 118-8 8 8 0 01-8 8z"/>
            <path d="M12 6v6l4 2"/>
          </svg>
        </div>
        <div class="platform-info">
          <h4>OpenAI 账户授权</h4>
        </div>
      </div>

      <p class="auth-desc">请按照以下步骤完成 OpenAI 账户的授权：</p>

      <div class="auth-steps">
        <div class="step-card">
          <div class="step-header">
            <div class="step-num">1</div>
            <div class="step-title">点击下方按钮生成授权链接</div>
          </div>
          <div class="step-content">
            <button v-if="!authUrl" class="btn btn-primary" :disabled="loading" @click="generateAuthUrl">
              <svg v-if="!loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10 13a5 5 0 007.54.54l3-3a5 5 0 00-7.07-7.07l-1.72 1.71"/>
                <path d="M14 11a5 5 0 00-7.54-.54l-3 3a5 5 0 007.07 7.07l1.71-1.71"/>
              </svg>
              <span v-if="loading" class="btn-spinner"></span>
              {{ loading ? '生成中...' : '生成授权链接' }}
            </button>
            <div v-else class="auth-url-container">
              <div class="url-input-group">
                <input type="text" :value="authUrl" readonly class="url-input" />
                <button class="copy-btn" @click="copyAuthUrl">
                  <svg v-if="copied" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="check">
                    <polyline points="20 6 9 17 4 12"/>
                  </svg>
                  <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                    <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                  </svg>
                </button>
              </div>
              <button class="link-btn" @click="regenerateAuthUrl">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"/>
                  <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
                </svg>
                重新生成
              </button>
            </div>
          </div>
        </div>

        <div class="step-card">
          <div class="step-header">
            <div class="step-num">2</div>
            <div class="step-title">在浏览器中打开链接并完成授权</div>
          </div>
          <div class="step-content">
            <p>请在新标签页中打开授权链接，登录您的 OpenAI 账户并授权。</p>
            <div class="info-banner compact">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <polyline points="12 6 12 12 16 14"/>
              </svg>
              <div>
                <strong>重要提示：</strong>授权后页面可能会加载较长时间，请耐心等待。<br/>
                当浏览器地址栏变为 <strong>http://localhost:1455/...</strong> 开头时，表示授权已完成。
              </div>
            </div>
          </div>
        </div>

        <div class="step-card">
          <div class="step-header">
            <div class="step-num">3</div>
            <div class="step-title">输入授权链接或 Code</div>
          </div>
          <div class="step-content">
            <textarea
              v-model="authCode"
              class="form-textarea"
              rows="3"
              placeholder="方式1：复制完整的链接（http://localhost:1455/auth/callback?code=...）&#10;方式2：仅复制 code 参数的值&#10;系统会自动识别并提取所需信息"
            ></textarea>
            <div class="info-banner compact">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="9" y1="18" x2="15" y2="18"/>
                <line x1="10" y1="22" x2="14" y2="22"/>
                <path d="M15.09 14c.18-.98.65-1.74 1.41-2.5A4.65 4.65 0 0018 8 6 6 0 006 8c0 1 .23 2.23 1.5 3.5A4.61 4.61 0 018.91 14"/>
              </svg>
              <span>您可以直接复制整个链接或仅复制 code 参数值，系统会自动识别。</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="oauth-actions">
      <button class="btn btn-outline" @click="$emit('back')">上一步</button>
      <button
        v-if="!(isClaudePlatform && authMethod === 'cookie')"
        class="btn btn-primary"
        :disabled="!canExchange || exchanging"
        @click="exchangeCode"
      >
        <span v-if="exchanging" class="btn-spinner"></span>
        {{ exchanging ? '验证中...' : '完成授权' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import api from '@/api'
import { ElMessage } from '@/utils/toast'

const props = defineProps({
  platform: {
    type: String,
    required: true
  },
  proxy: {
    type: Object,
    default: null
  },
  authMode: {
    type: String,
    default: 'oauth',
    validator: (value) => ['oauth', 'cookie'].includes(value)
  }
})

const emit = defineEmits(['success', 'back'])

// 状态
const loading = ref(false)
const exchanging = ref(false)
const authUrl = ref('')
const authCode = ref('')
const copied = ref(false)
const sessionId = ref('')

// Cookie自动授权相关状态
const authMethod = ref('manual')
const sessionKey = ref('')
const cookieAuthLoading = ref(false)
const cookieAuthError = ref('')
const showSessionKeyHelp = ref(false)
const batchProgress = ref({ current: 0, total: 0 })

// 根据 authMode 属性设置初始授权方式
watch(() => props.authMode, (mode) => {
  if (mode === 'cookie') {
    authMethod.value = 'cookie'
  } else {
    authMethod.value = 'manual'
  }
}, { immediate: true })

// 判断是否是 Claude 平台
const isClaudePlatform = computed(() => {
  return props.platform === 'claude' || props.platform === 'claude-official'
})

// 解析后的 sessionKey 数量
const parsedSessionKeyCount = computed(() => {
  return sessionKey.value
    .split('\n')
    .map(s => s.trim())
    .filter(s => s.length > 0).length
})

// 计算是否可以交换code
const canExchange = computed(() => {
  return authUrl.value && authCode.value.trim()
})

// 监听授权码输入，自动提取URL中的code参数
watch(authCode, (newValue) => {
  if (!newValue || typeof newValue !== 'string') return
  const trimmedValue = newValue.trim()
  if (!trimmedValue) return

  const isUrl = trimmedValue.startsWith('http://') || trimmedValue.startsWith('https://')
  if (isUrl) {
    if (trimmedValue.startsWith('http://localhost:45462') || trimmedValue.startsWith('http://localhost:1455')) {
      try {
        const url = new URL(trimmedValue)
        const code = url.searchParams.get('code')
        if (code) {
          authCode.value = code
          ElMessage.success('成功提取授权码！')
        }
      } catch (e) {
        console.error('Failed to parse URL:', e)
      }
    }
  }
})

// 获取 OAuth API 使用的 platform（openai-responses 映射为 openai）
function getOAuthPlatform() {
  if (props.platform === 'openai-responses') {
    return 'openai'
  }
  return props.platform
}

// 生成授权URL
async function generateAuthUrl() {
  loading.value = true
  authUrl.value = ''
  authCode.value = ''
  sessionId.value = ''

  try {
    const proxyConfig = props.proxy?.enabled ? {
      type: props.proxy.type,
      host: props.proxy.host,
      port: parseInt(props.proxy.port),
      username: props.proxy.username || null,
      password: props.proxy.password || null
    } : null

    const res = await api.generateOAuthUrl(getOAuthPlatform(), proxyConfig)
    authUrl.value = res.data.auth_url
    sessionId.value = res.data.session_id

    ElMessage.success('授权链接已生成')
  } catch (e) {
    ElMessage.error(e.response?.data?.message || e.message || '生成授权链接失败')
  } finally {
    loading.value = false
  }
}

// 重新生成授权URL
function regenerateAuthUrl() {
  authUrl.value = ''
  authCode.value = ''
  sessionId.value = ''
  generateAuthUrl()
}

// 复制授权URL
async function copyAuthUrl() {
  if (!authUrl.value) {
    ElMessage.error('请先生成授权链接')
    return
  }

  try {
    await navigator.clipboard.writeText(authUrl.value)
    copied.value = true
    ElMessage.success('链接已复制')
    setTimeout(() => { copied.value = false }, 2000)
  } catch (e) {
    const input = document.createElement('input')
    input.value = authUrl.value
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    copied.value = true
    ElMessage.success('链接已复制')
    setTimeout(() => { copied.value = false }, 2000)
  }
}

// 交换授权码
async function exchangeCode() {
  if (!canExchange.value) return

  exchanging.value = true
  try {
    const proxyConfig = props.proxy?.enabled ? {
      type: props.proxy.type,
      host: props.proxy.host,
      port: parseInt(props.proxy.port),
      username: props.proxy.username || null,
      password: props.proxy.password || null
    } : null

    const res = await api.exchangeOAuthCode(getOAuthPlatform(), authCode.value, sessionId.value, proxyConfig)
    const tokenInfo = res.data

    emit('success', tokenInfo)
  } catch (e) {
    ElMessage.error(e.response?.data?.message || e.message || '授权失败，请检查授权码是否正确')
  } finally {
    exchanging.value = false
  }
}

// Cookie自动授权处理（支持批量）
async function handleCookieAuth() {
  const sessionKeys = sessionKey.value
    .split('\n')
    .map(s => s.trim())
    .filter(s => s.length > 0)

  if (sessionKeys.length === 0) {
    cookieAuthError.value = '请输入至少一个 sessionKey'
    return
  }

  cookieAuthLoading.value = true
  cookieAuthError.value = ''
  batchProgress.value = { current: 0, total: sessionKeys.length }

  const proxyConfig = props.proxy?.enabled ? {
    type: props.proxy.type,
    host: props.proxy.host,
    port: parseInt(props.proxy.port),
    username: props.proxy.username || null,
    password: props.proxy.password || null
  } : null

  const results = []
  const errors = []

  for (let i = 0; i < sessionKeys.length; i++) {
    batchProgress.value.current = i + 1
    try {
      const res = await api.oauthByCookie(props.platform, sessionKeys[i], proxyConfig)
      results.push({
        ...res.data,
        session_key: sessionKeys[i]
      })
    } catch (e) {
      errors.push({
        index: i + 1,
        key: sessionKeys[i].substring(0, 20) + '...',
        error: e.response?.data?.message || e.message
      })
    }
  }

  batchProgress.value = { current: 0, total: 0 }
  cookieAuthLoading.value = false

  if (results.length > 0) {
    emit('success', results)
  }

  if (errors.length > 0 && results.length === 0) {
    cookieAuthError.value = '全部授权失败，请检查 sessionKey 是否有效'
  } else if (errors.length > 0) {
    cookieAuthError.value = `${errors.length} 个授权失败`
  }
}

// 切换授权方式时重置状态
function onAuthMethodChange() {
  sessionKey.value = ''
  cookieAuthError.value = ''
  cookieAuthLoading.value = false
  batchProgress.value = { current: 0, total: 0 }
  authUrl.value = ''
  authCode.value = ''
  sessionId.value = ''
}
</script>

<style scoped>
.oauth-flow {
  padding: 16px 0;
}

.oauth-section {
  background: var(--apple-fill-quaternary, #f5f5f7);
  border-radius: var(--apple-radius-lg, 12px);
  padding: 24px;
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
}

.oauth-section.claude {
  border-color: rgba(102, 126, 234, 0.3);
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.08) 0%, rgba(118, 75, 162, 0.08) 100%);
}

.oauth-section.gemini {
  border-color: rgba(79, 172, 254, 0.3);
  background: linear-gradient(135deg, rgba(79, 172, 254, 0.08) 0%, rgba(0, 242, 254, 0.08) 100%);
}

.oauth-section.openai {
  border-color: rgba(17, 153, 142, 0.3);
  background: linear-gradient(135deg, rgba(17, 153, 142, 0.08) 0%, rgba(56, 239, 125, 0.08) 100%);
}

.oauth-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 24px;
}

.platform-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--apple-radius-md, 12px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.platform-icon svg {
  width: 24px;
  height: 24px;
}

.platform-icon.claude {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.platform-icon.gemini {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.platform-icon.openai {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.platform-info h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
}

/* 授权方式选择 */
.auth-method-select {
  margin-bottom: 20px;
  padding: 20px;
  background: var(--apple-bg-primary, #ffffff);
  border-radius: var(--apple-radius-md, 12px);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
}

.method-label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
  margin-bottom: 14px;
}

.method-options {
  display: flex;
  gap: 12px;
}

.method-option {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  background: var(--apple-fill-quaternary, #f5f5f7);
  border: 2px solid transparent;
  border-radius: var(--apple-radius-md, 8px);
  cursor: pointer;
  transition: all 0.2s ease;
}

.method-option input {
  display: none;
}

.method-option svg {
  width: 20px;
  height: 20px;
  color: var(--apple-text-secondary, #86868b);
}

.method-option span {
  font-size: 14px;
  font-weight: 500;
  color: var(--apple-text-primary, #1d1d1f);
}

.method-option:hover {
  background: var(--apple-fill-tertiary, rgba(120, 120, 128, 0.12));
}

.method-option.active {
  border-color: var(--apple-blue, #007aff);
  background: var(--apple-blue-light, rgba(0, 122, 255, 0.1));
}

.method-option.active svg {
  color: var(--apple-blue, #007aff);
}

/* Cookie 授权 */
.cookie-auth {
  padding: 20px;
  background: var(--apple-bg-primary, #ffffff);
  border-radius: var(--apple-radius-md, 12px);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
}

/* 信息横幅 */
.info-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  background: var(--apple-blue-light, rgba(0, 122, 255, 0.1));
  border-radius: var(--apple-radius-md, 8px);
  margin-bottom: 20px;
  font-size: 14px;
  color: var(--apple-blue, #007aff);
}

.info-banner svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  margin-top: 2px;
}

.info-banner.compact {
  margin-top: 12px;
  margin-bottom: 0;
  padding: 12px 14px;
  font-size: 13px;
}

/* 警告横幅 */
.warning-banner {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  background: var(--apple-orange-light, rgba(255, 149, 0, 0.1));
  border-radius: var(--apple-radius-md, 8px);
  margin-top: 12px;
  font-size: 13px;
  color: var(--apple-orange, #ff9500);
}

.warning-banner svg {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  margin-top: 2px;
}

/* 错误横幅 */
.error-banner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: var(--apple-red-light, rgba(255, 59, 48, 0.1));
  border-radius: var(--apple-radius-md, 8px);
  margin-bottom: 16px;
  font-size: 14px;
  color: var(--apple-red, #ff3b30);
}

.error-banner svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

/* 表单元素 */
.form-group {
  margin-bottom: 16px;
}

.form-label-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 500;
  color: var(--apple-text-primary, #1d1d1f);
}

.form-label svg {
  width: 16px;
  height: 16px;
  color: var(--apple-blue, #007aff);
}

.count-badge {
  padding: 2px 8px;
  background: var(--apple-blue, #007aff);
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: var(--apple-radius-full, 100px);
}

.help-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: var(--apple-text-tertiary, #c7c7cc);
  cursor: pointer;
  transition: color 0.2s ease;
}

.help-btn:hover {
  color: var(--apple-blue, #007aff);
}

.help-btn svg {
  width: 18px;
  height: 18px;
}

.form-textarea {
  width: 100%;
  padding: 12px 14px;
  font-size: 14px;
  font-family: 'SF Mono', Monaco, monospace;
  color: var(--apple-text-primary, #1d1d1f);
  background: var(--apple-fill-quaternary, #f5f5f7);
  border: 1px solid transparent;
  border-radius: var(--apple-radius-md, 8px);
  resize: vertical;
  transition: all 0.2s ease;
}

.form-textarea::placeholder {
  color: var(--apple-text-placeholder, #c7c7cc);
  font-family: inherit;
}

.form-textarea:focus {
  background: var(--apple-bg-primary, #ffffff);
  border-color: var(--apple-blue, #007aff);
  box-shadow: 0 0 0 3px var(--apple-blue-light, rgba(0, 122, 255, 0.15));
  outline: none;
}

.session-key-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--apple-blue, #007aff);
}

.session-key-tip svg {
  width: 14px;
  height: 14px;
}

/* 帮助部分 */
.help-section {
  margin-top: 16px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(255, 204, 0, 0.1) 0%, rgba(255, 149, 0, 0.1) 100%);
  border: 1px solid rgba(255, 149, 0, 0.3);
  border-radius: var(--apple-radius-md, 8px);
}

.help-section h5 {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 600;
  color: var(--apple-orange, #ff9500);
}

.help-section h5 svg {
  width: 18px;
  height: 18px;
}

.help-section ol {
  margin: 0;
  padding-left: 20px;
  font-size: 13px;
  line-height: 1.8;
  color: var(--apple-text-primary, #1d1d1f);
}

.help-section kbd {
  background: var(--apple-fill-secondary, rgba(120, 120, 128, 0.16));
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
}

.help-section code {
  background: rgba(255, 149, 0, 0.2);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'SF Mono', Monaco, monospace;
  font-size: 12px;
}

.help-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  font-size: 12px;
  color: var(--apple-text-secondary, #86868b);
}

.help-tip svg {
  width: 14px;
  height: 14px;
}

/* 手动授权 */
.manual-auth {
  padding: 20px;
  background: var(--apple-bg-primary, #ffffff);
  border-radius: var(--apple-radius-md, 12px);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
}

.auth-desc {
  font-size: 14px;
  color: var(--apple-text-secondary, #86868b);
  margin: 0 0 20px;
}

.auth-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step-card {
  background: var(--apple-fill-quaternary, #f5f5f7);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-md, 8px);
  padding: 16px;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.step-num {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--apple-blue, #007aff);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.step-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
}

.step-content {
  padding-left: 40px;
}

.step-content p {
  font-size: 13px;
  color: var(--apple-text-secondary, #86868b);
  margin: 0;
  line-height: 1.6;
}

.step-desc {
  margin-bottom: 12px !important;
}

.input-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--apple-text-tertiary, #c7c7cc);
}

.input-tip svg {
  width: 14px;
  height: 14px;
}

/* URL 输入组 */
.auth-url-container {
  margin-top: 8px;
}

.url-input-group {
  display: flex;
  gap: 8px;
}

.url-input {
  flex: 1;
  padding: 10px 14px;
  font-size: 13px;
  font-family: 'SF Mono', Monaco, monospace;
  color: var(--apple-text-primary, #1d1d1f);
  background: var(--apple-bg-primary, #ffffff);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-md, 8px);
}

.copy-btn {
  width: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary, #f5f5f7);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-md, 8px);
  color: var(--apple-text-secondary, #86868b);
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  background: var(--apple-fill-tertiary, rgba(120, 120, 128, 0.12));
  color: var(--apple-text-primary, #1d1d1f);
}

.copy-btn svg {
  width: 18px;
  height: 18px;
}

.copy-btn svg.check {
  color: var(--apple-green, #34c759);
}

.link-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-top: 10px;
  padding: 0;
  background: transparent;
  border: none;
  font-size: 13px;
  color: var(--apple-blue, #007aff);
  cursor: pointer;
  transition: opacity 0.2s ease;
}

.link-btn:hover {
  opacity: 0.7;
}

.link-btn svg {
  width: 14px;
  height: 14px;
}

/* 按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  font-size: 15px;
  font-weight: 500;
  border: none;
  border-radius: var(--apple-radius-md, 8px);
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn svg {
  width: 18px;
  height: 18px;
}

.btn-primary {
  background: var(--apple-blue, #007aff);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--apple-blue-hover, #0066d6);
}

.btn-outline {
  background: var(--apple-fill-quaternary, #f5f5f7);
  color: var(--apple-text-primary, #1d1d1f);
  border: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-fill-tertiary, rgba(120, 120, 128, 0.12));
}

.btn-block {
  width: 100%;
}

.btn-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* 操作按钮 */
.oauth-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid var(--apple-separator, rgba(0, 0, 0, 0.1));
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
  background: var(--apple-text-primary, #1d1d1f);
  color: white;
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--apple-radius-lg, 12px);
  box-shadow: var(--apple-shadow-lg, 0 8px 32px rgba(0, 0, 0, 0.2));
  z-index: 2000;
  animation: slideUp 0.3s ease;
}

.toast svg {
  width: 18px;
  height: 18px;
}

.toast.success {
  background: var(--apple-green, #34c759);
}

.toast.error {
  background: var(--apple-red, #ff3b30);
}

/* 动画 */
@keyframes spin {
  to { transform: rotate(360deg); }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateX(-50%) translateY(20px); }
  to { opacity: 1; transform: translateX(-50%) translateY(0); }
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
  margin-top: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.slide-enter-to,
.slide-leave-from {
  opacity: 1;
  max-height: 500px;
}

/* 响应式 */
@media (max-width: 640px) {
  .method-options {
    flex-direction: column;
  }
}
</style>
