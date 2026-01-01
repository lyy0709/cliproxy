<!--
 * 文件作用：管理员登录页面 - Apple HIG 风格
 * 负责功能：
 *   - 用户名密码输入
 *   - 验证码获取和验证
 *   - 登录请求和状态管理
 * 重要程度：⭐⭐⭐⭐ 重要（系统入口）
-->
<template>
  <div class="login-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="bg-circle bg-circle-1"></div>
      <div class="bg-circle bg-circle-2"></div>
      <div class="bg-circle bg-circle-3"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- Logo -->
      <div class="login-header">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
        </div>
        <h1 class="login-title">Cli-Proxy</h1>
        <p class="login-subtitle">AI API 代理管理平台</p>
      </div>

      <!-- 登录表单 -->
      <form class="login-form" @submit.prevent="handleLogin">
        <!-- 用户名 -->
        <div class="form-group">
          <label class="form-label">用户名</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
            <input
              v-model="form.username"
              type="text"
              class="form-input"
              placeholder="请输入用户名"
              autocomplete="username"
            />
          </div>
          <span v-if="errors.username" class="form-error">{{ errors.username }}</span>
        </div>

        <!-- 密码 -->
        <div class="form-group">
          <label class="form-label">密码</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请输入密码"
              autocomplete="current-password"
            />
            <button type="button" class="password-toggle" @click="showPassword = !showPassword">
              <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            </button>
          </div>
          <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
        </div>

        <!-- 验证码 -->
        <div v-if="captchaEnabled" class="form-group">
          <label class="form-label">验证码</label>
          <div class="captcha-row">
            <div class="input-wrapper captcha-input">
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
              </svg>
              <input
                v-model="form.captchaCode"
                type="text"
                class="form-input"
                placeholder="请输入验证码"
                autocomplete="off"
                @keyup.enter="handleLogin"
              />
            </div>
            <div class="captcha-image-wrapper" @click="refreshCaptcha" :title="captchaLoading ? '加载中...' : '点击刷新'">
              <img v-if="captchaImage" :src="captchaImage" alt="验证码" class="captcha-image" />
              <div v-else class="captcha-placeholder">
                <svg v-if="captchaLoading" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
                  <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
                </svg>
                <span v-else>加载中</span>
              </div>
            </div>
          </div>
          <span v-if="errors.captcha" class="form-error">{{ errors.captcha }}</span>
        </div>

        <!-- 登录按钮 -->
        <button type="submit" class="login-btn" :disabled="loading">
          <svg v-if="loading" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
            <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
          </svg>
          <span v-else>登录</span>
        </button>
      </form>

      <!-- 返回首页 -->
      <div class="login-footer">
        <router-link to="/" class="back-link">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="19" y1="12" x2="5" y2="12"/>
            <polyline points="12,19 5,12 12,5"/>
          </svg>
          返回首页
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const showPassword = ref(false)
const captchaLoading = ref(false)
const captchaImage = ref('')
const captchaId = ref('')
const captchaEnabled = ref(true)

const form = reactive({
  username: '',
  password: '',
  captchaCode: ''
})

const errors = reactive({
  username: '',
  password: '',
  captcha: ''
})

// 获取验证码
async function refreshCaptcha() {
  captchaLoading.value = true
  try {
    const res = await api.getCaptcha()
    if (res.data.enabled === false) {
      captchaEnabled.value = false
      captchaId.value = ''
      captchaImage.value = ''
      return
    }
    captchaEnabled.value = true
    captchaId.value = res.data.captcha_id || ''
    captchaImage.value = res.data.image || ''
  } catch {
    ElMessage.error('获取验证码失败')
  } finally {
    captchaLoading.value = false
  }
}

function validate() {
  let valid = true
  errors.username = ''
  errors.password = ''
  errors.captcha = ''

  if (!form.username.trim()) {
    errors.username = '请输入用户名'
    valid = false
  }

  if (!form.password) {
    errors.password = '请输入密码'
    valid = false
  }

  if (captchaEnabled.value) {
    if (!captchaId.value) {
      errors.captcha = '验证码未加载，请刷新'
      valid = false
    } else if (!form.captchaCode.trim()) {
      errors.captcha = '请输入验证码'
      valid = false
    }
  }

  return valid
}

async function handleLogin() {
  if (!validate()) return

  loading.value = true
  try {
    const loginData = {
      username: form.username,
      password: form.password
    }
    if (captchaEnabled.value) {
      loginData.captcha_id = captchaId.value
      loginData.captcha_code = form.captchaCode
    }
    const res = await userStore.login(loginData)
    ElMessage.success('登录成功')

    // 检查是否需要强制修改密码
    if (res.data.must_change_password) {
      router.push('/force-change-password')
    } else {
      router.push('/admin/system-monitor')
    }
  } catch {
    // 登录失败，刷新验证码
    refreshCaptcha()
    form.captchaCode = ''
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: var(--apple-spacing-xl);
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.bg-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.bg-circle-1 {
  width: 600px;
  height: 600px;
  top: -200px;
  right: -100px;
}

.bg-circle-2 {
  width: 400px;
  height: 400px;
  bottom: -100px;
  left: -50px;
}

.bg-circle-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 20%;
}

/* 登录卡片 */
.login-card {
  width: 100%;
  max-width: 420px;
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xxl);
  box-shadow: var(--apple-shadow-xl);
  padding: var(--apple-spacing-xxl);
  position: relative;
  z-index: 1;
  animation: apple-scale-in var(--apple-duration-normal) var(--apple-ease-spring);
}

/* 头部 */
.login-header {
  text-align: center;
  margin-bottom: var(--apple-spacing-xxl);
}

.logo {
  width: 64px;
  height: 64px;
  margin: 0 auto var(--apple-spacing-lg);
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  border-radius: var(--apple-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon {
  width: 36px;
  height: 36px;
  color: white;
}

.login-title {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.login-subtitle {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

/* 表单 */
.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-lg);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.form-label {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--apple-spacing-md);
  width: 18px;
  height: 18px;
  color: var(--apple-text-tertiary);
  pointer-events: none;
}

.form-input {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  padding-left: 44px;
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-fill-quaternary);
  border: 1px solid transparent;
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-input::placeholder {
  color: var(--apple-text-placeholder);
}

.form-input:hover {
  background: var(--apple-fill-tertiary);
}

.form-input:focus {
  background: var(--apple-bg-primary);
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.password-toggle {
  position: absolute;
  right: var(--apple-spacing-sm);
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.password-toggle:hover {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-secondary);
}

.password-toggle svg {
  width: 18px;
  height: 18px;
}

.form-error {
  font-size: var(--apple-text-xs);
  color: var(--apple-red);
}

/* 验证码 */
.captcha-row {
  display: flex;
  gap: var(--apple-spacing-sm);
}

.captcha-input {
  flex: 1;
}

.captcha-image-wrapper {
  width: 120px;
  height: 44px;
  border-radius: var(--apple-radius-sm);
  cursor: pointer;
  background: var(--apple-fill-quaternary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid var(--apple-separator);
  overflow: hidden;
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.captcha-image-wrapper:hover {
  border-color: var(--apple-blue);
}

.captcha-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.captcha-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  padding: var(--apple-spacing-md);
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  color: white;
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-semibold);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  margin-top: var(--apple-spacing-sm);
}

.login-btn:hover:not(:disabled) {
  transform: scale(1.02);
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.4);
}

.login-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* Spinner */
.spinner {
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 底部 */
.login-footer {
  margin-top: var(--apple-spacing-xl);
  text-align: center;
}

.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  text-decoration: none;
  transition: color var(--apple-duration-fast) var(--apple-ease-default);
}

.back-link:hover {
  color: var(--apple-blue);
}

.back-link svg {
  width: 16px;
  height: 16px;
}

/* 响应式 - 移动端适配 */
@media (max-width: 767px) {
  .login-page {
    padding: var(--apple-spacing-md);
    padding-top: calc(var(--apple-safe-area-top) + var(--apple-spacing-md));
    padding-bottom: calc(var(--apple-safe-area-bottom) + var(--apple-spacing-md));
  }

  .login-card {
    padding: var(--apple-spacing-lg);
    border-radius: var(--apple-radius-xl);
  }

  .login-header {
    margin-bottom: var(--apple-spacing-xl);
  }

  .logo {
    width: 56px;
    height: 56px;
    margin-bottom: var(--apple-spacing-md);
  }

  .logo-icon {
    width: 30px;
    height: 30px;
  }

  .login-title {
    font-size: var(--apple-text-xl);
  }

  .login-subtitle {
    font-size: var(--apple-text-xs);
  }

  .login-form {
    gap: var(--apple-spacing-md);
  }

  .form-input {
    padding: var(--apple-spacing-md);
    padding-left: 44px;
    min-height: var(--apple-touch-target);
    font-size: 16px; /* 防止 iOS 自动缩放 */
  }

  .password-toggle {
    width: 44px;
    height: 44px;
    right: 0;
  }

  .captcha-row {
    flex-direction: column;
    gap: var(--apple-spacing-sm);
  }

  .captcha-input {
    width: 100%;
  }

  .captcha-image-wrapper {
    width: 100%;
    height: 48px;
  }

  .login-btn {
    min-height: var(--apple-touch-target);
    font-size: var(--apple-text-md);
  }

  /* 隐藏部分装饰圆 */
  .bg-circle-3 {
    display: none;
  }
}

@media (max-width: 375px) {
  .login-card {
    padding: var(--apple-spacing-md);
  }

  .logo {
    width: 48px;
    height: 48px;
  }

  .login-title {
    font-size: var(--apple-text-lg);
  }
}
</style>
