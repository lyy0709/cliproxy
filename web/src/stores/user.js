/**
 * 文件作用：用户状态管理，管理登录状态和用户信息
 * 负责功能：
 *   - 用户登录/登出
 *   - Token存储管理
 *   - 用户信息获取
 *   - 登录状态判断
 *   - 首次登录强制修改密码状态
 * 重要程度：⭐⭐⭐⭐ 重要（认证状态核心）
 * 依赖模块：pinia, api
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const mustChangePassword = ref(JSON.parse(localStorage.getItem('mustChangePassword') || 'false'))

  const isLoggedIn = computed(() => !!token.value)

  async function login(loginData) {
    const res = await api.login(loginData)
    token.value = res.data.token
    user.value = {
      username: res.data.username,
      role: res.data.role
    }
    // 保存是否需要强制修改密码的状态
    mustChangePassword.value = res.data.must_change_password || false

    localStorage.setItem('token', token.value)
    localStorage.setItem('user', JSON.stringify(user.value))
    localStorage.setItem('mustChangePassword', JSON.stringify(mustChangePassword.value))
    return res
  }

  function logout() {
    token.value = ''
    user.value = null
    mustChangePassword.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('mustChangePassword')
  }

  function setMustChangePassword(value) {
    mustChangePassword.value = value
    localStorage.setItem('mustChangePassword', JSON.stringify(value))
  }

  async function fetchProfile() {
    const res = await api.getProfile()
    user.value = res.data
    localStorage.setItem('user', JSON.stringify(user.value))
  }

  async function checkAdminStatus() {
    try {
      const res = await api.getAdminStatus()
      mustChangePassword.value = res.data.must_change_password || false
      localStorage.setItem('mustChangePassword', JSON.stringify(mustChangePassword.value))
      return res.data
    } catch {
      return null
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    mustChangePassword,
    login,
    logout,
    fetchProfile,
    setMustChangePassword,
    checkAdminStatus
  }
})
