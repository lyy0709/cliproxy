<!--
 * 文件作用：主布局组件 - Apple HIG 风格
 * 负责功能：
 *   - 毛玻璃侧边栏导航
 *   - 顶部导航栏
 *   - 内容区路由出口
 *   - 菜单折叠控制
 * 重要程度：⭐⭐⭐⭐ 重要（主布局框架）
-->
<template>
  <div class="app-layout">
    <!-- 侧边栏 -->
    <aside :class="['app-sidebar', { 'is-collapsed': isCollapse, 'is-mobile-open': mobileMenuOpen }]">
      <!-- 移动端遮罩 -->
      <div v-if="mobileMenuOpen" class="sidebar-backdrop" @click="mobileMenuOpen = false"></div>

      <!-- Logo -->
      <div class="sidebar-header">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
          <span v-if="!isCollapse" class="logo-text">Cli-Proxy</span>
        </div>
        <!-- 移动端关闭按钮 -->
        <button class="mobile-close-btn" @click="mobileMenuOpen = false">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"/>
            <line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>

      <!-- 导航菜单 -->
      <nav class="sidebar-nav">
        <div class="nav-group">
          <div v-if="!isCollapse" class="nav-group-title">监控</div>
          <router-link
            to="/admin/system-monitor"
            :class="['nav-item', { active: isActive('/admin/system-monitor') }]"
            @mouseenter="prefetchFor('/admin/system-monitor')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="2" y="3" width="20" height="14" rx="2"/>
                <path d="M8 21h8M12 17v4"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">系统监控</span>
          </router-link>

          <router-link
            to="/admin/account-load"
            :class="['nav-item', { active: isActive('/admin/account-load') }]"
            @mouseenter="prefetchFor('/admin/account-load')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M3 3v18h18"/>
                <path d="M18 9l-5 5-4-4-3 3"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">账户负载</span>
          </router-link>
        </div>

        <div class="nav-group">
          <div v-if="!isCollapse" class="nav-group-title">管理</div>
          <router-link
            to="/admin/accounts"
            :class="['nav-item', { active: isActive('/admin/accounts') }]"
            @mouseenter="prefetchFor('/admin/accounts')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M15 3h4a2 2 0 012 2v14a2 2 0 01-2 2h-4"/>
                <polyline points="10,17 15,12 10,7"/>
                <line x1="15" y1="12" x2="3" y2="12"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">账户管理</span>
          </router-link>

          <router-link
            to="/admin/proxies"
            :class="['nav-item', { active: isActive('/admin/proxies') }]"
            @mouseenter="prefetchFor('/admin/proxies')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <circle cx="12" cy="12" r="4"/>
                <line x1="21.17" y1="8" x2="12" y2="8"/>
                <line x1="3.95" y1="6.06" x2="8.54" y2="14"/>
                <line x1="10.88" y1="21.94" x2="15.46" y2="14"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">代理管理</span>
          </router-link>

          <router-link
            to="/admin/models"
            :class="['nav-item', { active: isActive('/admin/models') }]"
            @mouseenter="prefetchFor('/admin/models')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="4" y="4" width="16" height="16" rx="2"/>
                <rect x="9" y="9" width="6" height="6"/>
                <line x1="9" y1="1" x2="9" y2="4"/>
                <line x1="15" y1="1" x2="15" y2="4"/>
                <line x1="9" y1="20" x2="9" y2="23"/>
                <line x1="15" y1="20" x2="15" y2="23"/>
                <line x1="20" y1="9" x2="23" y2="9"/>
                <line x1="20" y1="14" x2="23" y2="14"/>
                <line x1="1" y1="9" x2="4" y2="9"/>
                <line x1="1" y1="14" x2="4" y2="14"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">模型管理</span>
          </router-link>

          <router-link
            to="/admin/api-keys"
            :class="['nav-item', { active: isActive('/admin/api-keys') }]"
            @mouseenter="prefetchFor('/admin/api-keys')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">API Key</span>
          </router-link>

          <router-link
            to="/admin/usage"
            :class="['nav-item', { active: isActive('/admin/usage') }]"
            @mouseenter="prefetchFor('/admin/usage')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M3 3v18h18"/>
                <path d="M18 9l-5 5-4-4-3 3"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">使用统计</span>
          </router-link>
        </div>

        <div class="nav-group">
          <div v-if="!isCollapse" class="nav-group-title">系统</div>
          <router-link
            to="/admin/cache"
            :class="['nav-item', { active: isActive('/admin/cache') }]"
            @mouseenter="prefetchFor('/admin/cache')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <ellipse cx="12" cy="5" rx="9" ry="3"/>
                <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/>
                <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">缓存管理</span>
          </router-link>

          <router-link
            to="/admin/settings"
            :class="['nav-item', { active: isActive('/admin/settings') }]"
            @mouseenter="prefetchFor('/admin/settings')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="3"/>
                <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">系统设置</span>
          </router-link>

          <router-link
            to="/admin/client-filter"
            :class="['nav-item', { active: isActive('/admin/client-filter') }]"
            @mouseenter="prefetchFor('/admin/client-filter')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <polygon points="22,3 2,3 10,12.46 10,19 14,21 14,12.46"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">客户端过滤</span>
          </router-link>
        </div>

        <div class="nav-group">
          <div v-if="!isCollapse" class="nav-group-title">日志</div>
          <router-link
            to="/admin/request-logs"
            :class="['nav-item', { active: isActive('/admin/request-logs') }]"
            @mouseenter="prefetchFor('/admin/request-logs')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10,9 9,9 8,9"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">请求日志</span>
          </router-link>

          <router-link
            to="/admin/error-messages"
            :class="['nav-item', { active: isActive('/admin/error-messages') }]"
            @mouseenter="prefetchFor('/admin/error-messages')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">错误消息</span>
          </router-link>

          <router-link
            to="/admin/operation-logs"
            :class="['nav-item', { active: isActive('/admin/operation-logs') }]"
            @mouseenter="prefetchFor('/admin/operation-logs')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M12 20h9"/>
                <path d="M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4L16.5 3.5z"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">操作日志</span>
          </router-link>

          <router-link
            to="/admin/system-logs"
            :class="['nav-item', { active: isActive('/admin/system-logs') }]"
            @mouseenter="prefetchFor('/admin/system-logs')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <polyline points="4,17 10,11 4,5"/>
                <line x1="12" y1="19" x2="20" y2="19"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">系统日志</span>
          </router-link>
        </div>
      </nav>

      <!-- 底部切换 -->
      <div class="sidebar-footer">
        <button class="collapse-toggle" @click="isCollapse = !isCollapse">
          <svg v-if="isCollapse" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polyline points="9,18 15,12 9,6"/>
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polyline points="15,18 9,12 15,6"/>
          </svg>
        </button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="app-main">
      <!-- 顶部导航栏 -->
      <header class="app-navbar">
        <!-- 移动端汉堡菜单 -->
        <button class="mobile-menu-btn" @click="mobileMenuOpen = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="12" x2="21" y2="12"/>
            <line x1="3" y1="6" x2="21" y2="6"/>
            <line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
        </button>
        <div class="navbar-title">
          {{ currentPageTitle }}
        </div>
        <div class="navbar-actions">
          <div class="user-dropdown" @click="showDropdown = !showDropdown" v-click-outside="closeDropdown">
            <div class="user-avatar">
              {{ userStore.user?.username?.charAt(0)?.toUpperCase() || 'A' }}
            </div>
            <span class="user-name">{{ userStore.user?.username }}</span>
            <svg class="dropdown-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6,9 12,15 18,9"/>
            </svg>

            <!-- 下拉菜单 -->
            <div v-if="showDropdown" class="dropdown-menu">
              <div class="dropdown-item" @click="openPasswordDialog">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0110 0v4"/>
                </svg>
                修改密码
              </div>
              <div class="dropdown-divider"></div>
              <div class="dropdown-item danger" @click="handleLogout">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/>
                  <polyline points="16,17 21,12 16,7"/>
                  <line x1="21" y1="12" x2="9" y2="12"/>
                </svg>
                退出登录
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- 内容区 -->
      <main class="app-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <keep-alive :max="5" :include="['Profile']">
              <component :is="Component" />
            </keep-alive>
          </transition>
        </router-view>
      </main>
    </div>

    <!-- 修改密码弹窗 -->
    <Teleport to="body">
      <div v-if="passwordDialogVisible" class="modal-overlay" @click.self="closePasswordDialog">
        <div class="modal modal-sm">
          <div class="modal-header">
            <h2>修改密码</h2>
            <button class="modal-close" @click="closePasswordDialog">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <form class="password-form" @submit.prevent="handleChangePassword">
              <div class="form-group">
                <label class="form-label">当前密码</label>
                <input v-model="passwordForm.oldPassword" type="password" class="form-input" placeholder="请输入当前密码" autocomplete="current-password" />
                <span v-if="passwordErrors.oldPassword" class="form-error">{{ passwordErrors.oldPassword }}</span>
              </div>
              <div class="form-group">
                <label class="form-label">新密码</label>
                <input v-model="passwordForm.newPassword" type="password" class="form-input" placeholder="至少8位，包含字母和数字" autocomplete="new-password" />
                <span v-if="passwordErrors.newPassword" class="form-error">{{ passwordErrors.newPassword }}</span>
              </div>
              <div class="form-group">
                <label class="form-label">确认密码</label>
                <input v-model="passwordForm.confirmPassword" type="password" class="form-input" placeholder="请再次输入新密码" autocomplete="new-password" />
                <span v-if="passwordErrors.confirmPassword" class="form-error">{{ passwordErrors.confirmPassword }}</span>
              </div>
              <div class="password-hint">密码需至少 8 位，且包含字母与数字</div>
            </form>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="closePasswordDialog">取消</button>
            <button class="btn btn-primary" :disabled="changingPassword || !isPasswordFormValid" @click="handleChangePassword">
              <span v-if="changingPassword" class="btn-loading"></span>
              {{ changingPassword ? '修改中...' : '确认修改' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { prefetchChunk } from '@/prefetch'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const showDropdown = ref(false)
const mobileMenuOpen = ref(false)
const passwordDialogVisible = ref(false)
const changingPassword = ref(false)

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordErrors = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const hasMinLength = computed(() => passwordForm.newPassword.length >= 8)
const hasLetter = computed(() => /[a-zA-Z]/.test(passwordForm.newPassword))
const hasDigit = computed(() => /\d/.test(passwordForm.newPassword))

const isPasswordFormValid = computed(() => {
  return passwordForm.oldPassword &&
    hasMinLength.value &&
    hasLetter.value &&
    hasDigit.value &&
    passwordForm.confirmPassword &&
    passwordForm.confirmPassword === passwordForm.newPassword
})

// 页面标题映射
const pageTitles = {
  '/admin/system-monitor': '系统监控',
  '/admin/accounts': '账户管理',
  '/admin/proxies': '代理管理',
  '/admin/models': '模型管理',
  '/admin/request-logs': '请求日志',
  '/admin/account-load': '账户负载',
  '/admin/cache': '缓存管理',
  '/admin/api-keys': 'API Key 管理',
  '/admin/usage': '使用统计',
  '/admin/settings': '系统设置',
  '/admin/error-messages': '错误消息',
  '/admin/operation-logs': '操作日志',
  '/admin/system-logs': '系统日志',
  '/admin/client-filter': '客户端过滤'
}

const currentPageTitle = computed(() => pageTitles[route.path] || '管理后台')

function isActive(path) {
  return route.path === path
}

function prefetchFor(path) {
  const loaders = {
    '/admin/system-monitor': () => import('@/views/SystemMonitor.vue'),
    '/admin/accounts': () => import('@/views/Accounts.vue'),
    '/admin/proxies': () => import('@/views/Proxies.vue'),
    '/admin/models': () => import('@/views/Models.vue'),
    '/admin/request-logs': () => import('@/views/RequestLogs.vue'),
    '/admin/account-load': () => import('@/views/AccountLoad.vue'),
    '/admin/cache': () => import('@/views/Cache.vue'),
    '/admin/api-keys': () => import('@/views/APIKeys.vue'),
    '/admin/usage': () => import('@/views/Usage.vue'),
    '/admin/settings': () => import('@/views/Settings.vue'),
    '/admin/error-messages': () => import('@/views/ErrorMessages.vue'),
    '/admin/operation-logs': () => import('@/views/OperationLogs.vue'),
    '/admin/system-logs': () => import('@/views/SystemLogs.vue'),
    '/admin/client-filter': () => import('@/views/ClientFilter.vue')
  }
  const loader = loaders[path]
  if (!loader) return
  prefetchChunk(path, loader)
}

function closeDropdown() {
  showDropdown.value = false
}

function handleLogout() {
  showDropdown.value = false
  userStore.logout()
  router.push('/login')
}

function openPasswordDialog() {
  showDropdown.value = false
  passwordDialogVisible.value = true
}

function closePasswordDialog() {
  passwordDialogVisible.value = false
  passwordForm.oldPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  passwordErrors.oldPassword = ''
  passwordErrors.newPassword = ''
  passwordErrors.confirmPassword = ''
}

function validatePasswordForm() {
  passwordErrors.oldPassword = ''
  passwordErrors.newPassword = ''
  passwordErrors.confirmPassword = ''

  let valid = true
  if (!passwordForm.oldPassword) {
    passwordErrors.oldPassword = '请输入当前密码'
    valid = false
  }
  if (!passwordForm.newPassword) {
    passwordErrors.newPassword = '请输入新密码'
    valid = false
  } else if (!hasMinLength.value || !hasLetter.value || !hasDigit.value) {
    passwordErrors.newPassword = '密码强度不足'
    valid = false
  }
  if (!passwordForm.confirmPassword) {
    passwordErrors.confirmPassword = '请确认新密码'
    valid = false
  } else if (passwordForm.confirmPassword !== passwordForm.newPassword) {
    passwordErrors.confirmPassword = '两次输入的密码不一致'
    valid = false
  }
  return valid
}

async function handleChangePassword() {
  if (!validatePasswordForm()) return
  changingPassword.value = true
  try {
    await api.changePassword({
      old_password: passwordForm.oldPassword,
      new_password: passwordForm.newPassword
    })
    ElMessage.success('密码修改成功')
    closePasswordDialog()
  } catch {
    // 错误已在 API 层处理
  } finally {
    changingPassword.value = false
  }
}

// 点击外部关闭下拉菜单指令
const vClickOutside = {
  mounted(el, binding) {
    el._clickOutside = (e) => {
      if (!el.contains(e.target)) {
        binding.value()
      }
    }
    document.addEventListener('click', el._clickOutside)
  },
  unmounted(el) {
    document.removeEventListener('click', el._clickOutside)
  }
}
</script>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  background: var(--apple-bg-secondary);
}

/* 侧边栏 */
.app-sidebar {
  width: var(--apple-sidebar-width);
  height: 100%;
  background: var(--apple-bg-blur);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-right: 1px solid var(--apple-separator);
  display: flex;
  flex-direction: column;
  transition: width var(--apple-duration-normal) var(--apple-ease-default);
  z-index: 100;
}

.app-sidebar.is-collapsed {
  width: var(--apple-sidebar-collapsed);
}

.sidebar-header {
  height: 56px;
  padding: 0 var(--apple-spacing-lg);
  display: flex;
  align-items: center;
  border-bottom: 1px solid var(--apple-separator);
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.logo-icon {
  width: 28px;
  height: 28px;
  color: var(--apple-blue);
}

.logo-text {
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  white-space: nowrap;
}

/* 导航 */
.sidebar-nav {
  flex: 1;
  padding: var(--apple-spacing-sm);
  overflow-y: auto;
}

.nav-group {
  margin-bottom: var(--apple-spacing-md);
}

.nav-group-title {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  text-decoration: none;
  margin-bottom: 2px;
}

.nav-item:hover {
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-primary);
}

.nav-item.active {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.nav-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-icon svg {
  width: 18px;
  height: 18px;
}

.nav-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 侧边栏底部 */
.sidebar-footer {
  padding: var(--apple-spacing-sm);
  border-top: 1px solid var(--apple-separator);
}

.user-switch {
  margin-bottom: var(--apple-spacing-sm);
}

.collapse-toggle {
  width: 100%;
  padding: var(--apple-spacing-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.collapse-toggle:hover {
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-primary);
}

.collapse-toggle svg {
  width: 18px;
  height: 18px;
}

/* 主内容区 */
.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 顶部导航栏 */
.app-navbar {
  height: var(--apple-navbar-height);
  background: var(--apple-bg-blur);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid var(--apple-separator);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--apple-spacing-xl);
  position: sticky;
  top: 0;
  z-index: 50;
}

.navbar-title {
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.navbar-actions {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
}

/* 用户下拉菜单 */
.user-dropdown {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: background var(--apple-duration-fast) var(--apple-ease-default);
}

.user-dropdown:hover {
  background: var(--apple-fill-quaternary);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-full);
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-semibold);
}

.user-name {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
}

.dropdown-arrow {
  width: 16px;
  height: 16px;
  color: var(--apple-text-tertiary);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 180px;
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-md);
  box-shadow: var(--apple-shadow-lg);
  border: 1px solid var(--apple-separator);
  padding: var(--apple-spacing-xs);
  animation: apple-fade-in-down var(--apple-duration-fast) var(--apple-ease-default);
  z-index: 100;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  border-radius: var(--apple-radius-sm);
  cursor: pointer;
  transition: background var(--apple-duration-fast) var(--apple-ease-default);
}

.dropdown-item:hover {
  background: var(--apple-fill-quaternary);
}

.dropdown-item svg {
  width: 16px;
  height: 16px;
  color: var(--apple-text-secondary);
}

.dropdown-item.danger {
  color: var(--apple-red);
}

.dropdown-item.danger svg {
  color: var(--apple-red);
}

.dropdown-divider {
  height: 1px;
  background: var(--apple-separator);
  margin: var(--apple-spacing-xs) 0;
}

/* 弹窗 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 200;
}

.modal {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-lg);
  border: 1px solid var(--apple-separator);
  width: 420px;
  max-width: calc(100vw - 32px);
  overflow: hidden;
}

.modal-sm {
  width: 420px;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.modal-header h2 {
  margin: 0;
  font-size: var(--apple-text-lg);
}

.modal-close {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--apple-text-tertiary);
}

.modal-close svg {
  width: 18px;
  height: 18px;
}

.modal-body {
  padding: var(--apple-spacing-lg);
}

.modal-footer {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
  display: flex;
  justify-content: flex-end;
  gap: var(--apple-spacing-sm);
}

.form-group {
  margin-bottom: var(--apple-spacing-md);
}

.form-label {
  display: block;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-bottom: var(--apple-spacing-xs);
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border-radius: var(--apple-radius-md);
  border: 1px solid var(--apple-separator);
  background: var(--apple-bg-secondary);
  color: var(--apple-text-primary);
  outline: none;
}

.form-input:focus {
  border-color: var(--apple-blue);
  background: var(--apple-bg-primary);
}

.form-error {
  display: block;
  margin-top: 6px;
  font-size: var(--apple-text-xs);
  color: var(--apple-red);
}

.password-hint {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

/* 内容区 */
.app-content {
  flex: 1;
  padding: var(--apple-spacing-xl);
  overflow-y: auto;
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--apple-duration-fast) var(--apple-ease-default);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 折叠状态下的样式调整 */
.is-collapsed .nav-group-title {
  display: none;
}

.is-collapsed .nav-item {
  justify-content: center;
  padding: var(--apple-spacing-sm);
}

.is-collapsed .logo {
  justify-content: center;
}

/* 移动端菜单按钮 - 默认隐藏 */
.mobile-menu-btn {
  display: none;
  width: 40px;
  height: 40px;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-primary);
  border-radius: var(--apple-radius-sm);
  transition: background var(--apple-duration-fast) var(--apple-ease-default);
}

.mobile-menu-btn:hover {
  background: var(--apple-fill-quaternary);
}

.mobile-menu-btn svg {
  width: 22px;
  height: 22px;
}

/* 移动端关闭按钮 - 默认隐藏 */
.mobile-close-btn {
  display: none;
  width: 36px;
  height: 36px;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.mobile-close-btn:hover {
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-primary);
}

.mobile-close-btn svg {
  width: 20px;
  height: 20px;
}

/* 侧边栏遮罩 - 默认隐藏 */
.sidebar-backdrop {
  display: none;
}

/* ========================================
   移动端响应式 - Apple HIG
   ======================================== */
@media (max-width: 767px) {
  /* 侧边栏变为抽屉式 */
  .app-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: 280px;
    z-index: 200;
    transform: translateX(-100%);
    transition: transform var(--apple-duration-normal) var(--apple-ease-default);
  }

  .app-sidebar.is-mobile-open {
    transform: translateX(0);
  }

  .app-sidebar.is-collapsed {
    width: 280px;
    transform: translateX(-100%);
  }

  .app-sidebar.is-collapsed.is-mobile-open {
    transform: translateX(0);
  }

  /* 侧边栏遮罩 */
  .sidebar-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.35);
    z-index: -1;
  }

  /* 显示移动端按钮 */
  .mobile-menu-btn {
    display: flex;
  }

  .mobile-close-btn {
    display: flex;
  }

  /* 隐藏桌面端折叠按钮 */
  .sidebar-footer {
    display: none;
  }

  /* 侧边栏头部调整 */
  .sidebar-header {
    justify-content: space-between;
    padding: 0 var(--apple-spacing-md);
    height: 56px;
  }

  /* 导航项目触控优化 */
  .nav-item {
    min-height: var(--apple-touch-target);
    padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  }

  /* 折叠状态下仍显示标签 */
  .is-collapsed .nav-group-title,
  .is-collapsed .nav-label {
    display: block;
  }

  .is-collapsed .nav-item {
    justify-content: flex-start;
    padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  }

  .is-collapsed .logo {
    justify-content: flex-start;
  }

  /* 顶部导航栏 */
  .app-navbar {
    padding: 0 var(--apple-spacing-md);
    height: 56px;
  }

  .navbar-title {
    font-size: var(--apple-text-base);
    flex: 1;
    text-align: center;
  }

  /* 用户下拉 */
  .user-name {
    display: none;
  }

  .dropdown-arrow {
    display: none;
  }

  .user-dropdown {
    padding: var(--apple-spacing-xxs);
  }

  .dropdown-menu {
    right: -8px;
  }

  /* 内容区 */
  .app-content {
    padding: var(--apple-spacing-md);
  }

  /* 模态框 */
  .modal-overlay {
    padding: var(--apple-spacing-md);
    align-items: flex-end;
  }

  .modal,
  .modal.modal-sm {
    max-width: 100%;
    border-radius: var(--apple-radius-xl) var(--apple-radius-xl) 0 0;
  }

  .modal-header {
    padding: var(--apple-spacing-lg) var(--apple-spacing-md);
  }

  .modal-body {
    padding: var(--apple-spacing-md);
  }

  .modal-footer {
    padding: var(--apple-spacing-md);
    flex-direction: column;
    gap: var(--apple-spacing-sm);
  }

  .modal-footer .btn {
    width: 100%;
    min-height: var(--apple-touch-target);
  }

  .form-input {
    min-height: var(--apple-touch-target);
    font-size: 16px; /* 防止 iOS 缩放 */
  }
}

/* 平板端适配 */
@media (min-width: 768px) and (max-width: 1023px) {
  .app-sidebar {
    width: var(--apple-sidebar-collapsed);
  }

  .app-sidebar .nav-group-title,
  .app-sidebar .nav-label,
  .app-sidebar .logo-text {
    display: none;
  }

  .app-sidebar .nav-item {
    justify-content: center;
    padding: var(--apple-spacing-sm);
  }

  .app-sidebar .logo {
    justify-content: center;
  }

  .sidebar-footer {
    display: none;
  }
}
</style>
