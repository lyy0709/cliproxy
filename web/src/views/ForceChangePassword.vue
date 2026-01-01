<!--
 * 文件作用：首次登录强制修改密码页面 - Apple HIG 风格
 * 负责功能：
 *   - 首次登录时强制修改默认密码
 *   - 密码强度验证（至少8位，包含字母和数字）
 *   - 密码确认
 * 重要程度：⭐⭐⭐⭐⭐ 核心（安全核心功能）
-->
<template>
  <div class="force-change-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="bg-circle bg-circle-1"></div>
      <div class="bg-circle bg-circle-2"></div>
      <div class="bg-circle bg-circle-3"></div>
    </div>

    <!-- 修改密码卡片 -->
    <div class="change-card">
      <!-- 警告图标 -->
      <div class="card-header">
        <div class="warning-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
          </svg>
        </div>
        <h1 class="card-title">安全提醒</h1>
        <p class="card-subtitle">首次登录，请修改默认密码以确保账户安全</p>
      </div>

      <!-- 密码表单 -->
      <form class="password-form" @submit.prevent="handleSubmit">
        <!-- 新密码 -->
        <div class="form-group">
          <label class="form-label">新密码</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
            <input
              v-model="form.newPassword"
              :type="showNewPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请输入新密码"
              autocomplete="new-password"
            />
            <button type="button" class="password-toggle" @click="showNewPassword = !showNewPassword">
              <svg v-if="!showNewPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            </button>
          </div>
          <span v-if="errors.newPassword" class="form-error">{{ errors.newPassword }}</span>
        </div>

        <!-- 密码强度指示 -->
        <div class="password-strength" v-if="form.newPassword">
          <div class="strength-bars">
            <div :class="['strength-bar', strengthLevel >= 1 ? strengthClass : '']"></div>
            <div :class="['strength-bar', strengthLevel >= 2 ? strengthClass : '']"></div>
            <div :class="['strength-bar', strengthLevel >= 3 ? strengthClass : '']"></div>
          </div>
          <span :class="['strength-text', strengthClass]">{{ strengthText }}</span>
        </div>

        <!-- 密码要求提示 -->
        <div class="password-requirements">
          <div :class="['requirement', hasMinLength ? 'valid' : '']">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline v-if="hasMinLength" points="20,6 9,17 4,12"/>
              <circle v-else cx="12" cy="12" r="10"/>
            </svg>
            <span>至少 8 个字符</span>
          </div>
          <div :class="['requirement', hasLetter ? 'valid' : '']">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline v-if="hasLetter" points="20,6 9,17 4,12"/>
              <circle v-else cx="12" cy="12" r="10"/>
            </svg>
            <span>包含字母</span>
          </div>
          <div :class="['requirement', hasDigit ? 'valid' : '']">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline v-if="hasDigit" points="20,6 9,17 4,12"/>
              <circle v-else cx="12" cy="12" r="10"/>
            </svg>
            <span>包含数字</span>
          </div>
        </div>

        <!-- 确认密码 -->
        <div class="form-group">
          <label class="form-label">确认密码</label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
            <input
              v-model="form.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请再次输入新密码"
              autocomplete="new-password"
            />
            <button type="button" class="password-toggle" @click="showConfirmPassword = !showConfirmPassword">
              <svg v-if="!showConfirmPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            </button>
          </div>
          <span v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</span>
        </div>

        <!-- 提交按钮 -->
        <button type="submit" class="submit-btn" :disabled="loading || !isValid">
          <svg v-if="loading" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
            <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
          </svg>
          <span v-else>确认修改</span>
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

const form = reactive({
  newPassword: '',
  confirmPassword: ''
})

const errors = reactive({
  newPassword: '',
  confirmPassword: ''
})

// 密码强度检测
const hasMinLength = computed(() => form.newPassword.length >= 8)
const hasLetter = computed(() => /[a-zA-Z]/.test(form.newPassword))
const hasDigit = computed(() => /\d/.test(form.newPassword))
const hasSpecial = computed(() => /[!@#$%^&*(),.?":{}|<>]/.test(form.newPassword))

const strengthLevel = computed(() => {
  let level = 0
  if (hasMinLength.value) level++
  if (hasLetter.value && hasDigit.value) level++
  if (hasSpecial.value && form.newPassword.length >= 12) level++
  return level
})

const strengthClass = computed(() => {
  if (strengthLevel.value === 1) return 'weak'
  if (strengthLevel.value === 2) return 'medium'
  if (strengthLevel.value === 3) return 'strong'
  return ''
})

const strengthText = computed(() => {
  if (strengthLevel.value === 1) return '弱'
  if (strengthLevel.value === 2) return '中'
  if (strengthLevel.value === 3) return '强'
  return ''
})

// 验证表单是否有效
const isValid = computed(() => {
  return hasMinLength.value && hasLetter.value && hasDigit.value &&
    form.confirmPassword === form.newPassword && form.confirmPassword.length > 0
})

function validate() {
  let valid = true
  errors.newPassword = ''
  errors.confirmPassword = ''

  if (!form.newPassword) {
    errors.newPassword = '请输入新密码'
    valid = false
  } else if (!hasMinLength.value || !hasLetter.value || !hasDigit.value) {
    errors.newPassword = '密码不符合要求'
    valid = false
  }

  if (!form.confirmPassword) {
    errors.confirmPassword = '请确认新密码'
    valid = false
  } else if (form.confirmPassword !== form.newPassword) {
    errors.confirmPassword = '两次输入的密码不一致'
    valid = false
  }

  return valid
}

async function handleSubmit() {
  if (!validate()) return

  loading.value = true
  try {
    await api.forceChangePassword({
      new_password: form.newPassword,
      confirm_password: form.confirmPassword
    })

    // 更新 store 状态
    userStore.setMustChangePassword(false)

    ElMessage.success('密码修改成功')
    router.push('/admin/system-monitor')
  } catch {
    // 错误已在 API 层处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.force-change-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
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

/* 卡片 */
.change-card {
  width: 100%;
  max-width: 460px;
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xxl);
  box-shadow: var(--apple-shadow-xl);
  padding: var(--apple-spacing-xxl);
  position: relative;
  z-index: 1;
  animation: apple-scale-in var(--apple-duration-normal) var(--apple-ease-spring);
}

/* 头部 */
.card-header {
  text-align: center;
  margin-bottom: var(--apple-spacing-xl);
}

.warning-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto var(--apple-spacing-lg);
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  border-radius: var(--apple-radius-xl);
  display: flex;
  align-items: center;
  justify-content: center;
}

.warning-icon svg {
  width: 40px;
  height: 40px;
  color: white;
}

.card-title {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.card-subtitle {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

/* 表单 */
.password-form {
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
  padding-right: 44px;
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

/* 密码强度 */
.password-strength {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.strength-bars {
  display: flex;
  gap: 4px;
  flex: 1;
}

.strength-bar {
  height: 4px;
  flex: 1;
  background: var(--apple-fill-tertiary);
  border-radius: 2px;
  transition: background var(--apple-duration-fast) var(--apple-ease-default);
}

.strength-bar.weak {
  background: var(--apple-red);
}

.strength-bar.medium {
  background: var(--apple-orange);
}

.strength-bar.strong {
  background: var(--apple-green);
}

.strength-text {
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  min-width: 24px;
}

.strength-text.weak {
  color: var(--apple-red);
}

.strength-text.medium {
  color: var(--apple-orange);
}

.strength-text.strong {
  color: var(--apple-green);
}

/* 密码要求 */
.password-requirements {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-md);
  background: var(--apple-fill-quaternary);
  border-radius: var(--apple-radius-md);
}

.requirement {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  transition: color var(--apple-duration-fast) var(--apple-ease-default);
}

.requirement svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.requirement.valid {
  color: var(--apple-green);
}

.requirement.valid svg {
  color: var(--apple-green);
}

/* 提交按钮 */
.submit-btn {
  width: 100%;
  padding: var(--apple-spacing-md);
  background: linear-gradient(135deg, #ff6b6b 0%, #ee5a24 100%);
  color: white;
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-semibold);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  margin-top: var(--apple-spacing-sm);
}

.submit-btn:hover:not(:disabled) {
  transform: scale(1.02);
  box-shadow: 0 4px 16px rgba(238, 90, 36, 0.4);
}

.submit-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.submit-btn:disabled {
  opacity: 0.5;
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

/* 响应式 */
@media (max-width: 480px) {
  .change-card {
    padding: var(--apple-spacing-xl);
  }
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .force-change-page {
    background: linear-gradient(135deg, #c0392b 0%, #e74c3c 100%);
  }
}
</style>
