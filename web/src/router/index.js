/**
 * 文件作用：Vue路由配置，定义应用页面路由和导航守卫
 * 负责功能：
 *   - 页面路由定义
 *   - 权限路由守卫
 *   - 管理员路由
 *   - 登录状态检查
 * 重要程度：⭐⭐⭐⭐ 重要（前端路由核心）
 * 依赖模块：vue-router, user store
 */
import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true }
  },
  // API Key 用量查询页面（公开，无需登录）
  {
    path: '/usage-query',
    name: 'UsageQuery',
    component: () => import('@/views/UsageQuery.vue'),
    meta: { guest: true }
  },
  // 首页（重定向到管理后台）
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { guest: true }
  },
  // 后台入口
  {
    path: '/dashboard',
    redirect: '/admin/system-monitor'
  },
  // 管理员后台路由
  {
    path: '/admin',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/admin/system-monitor'
      },
      {
        path: 'system-monitor',
        name: 'SystemMonitor',
        component: () => import('@/views/SystemMonitor.vue')
      },
      {
        path: 'accounts',
        name: 'Accounts',
        component: () => import('@/views/Accounts.vue')
      },
      {
        path: 'models',
        name: 'Models',
        component: () => import('@/views/Models.vue')
      },
      {
        path: 'api-keys',
        name: 'AdminAPIKeys',
        component: () => import('@/views/APIKeys.vue')
      },
      {
        path: 'request-logs',
        name: 'RequestLogs',
        component: () => import('@/views/RequestLogs.vue')
      },
      {
        path: 'account-load',
        name: 'AccountLoad',
        component: () => import('@/views/AccountLoad.vue')
      },
      {
        path: 'cache',
        name: 'Cache',
        component: () => import('@/views/Cache.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/Settings.vue')
      },
      {
        path: 'proxies',
        name: 'Proxies',
        component: () => import('@/views/Proxies.vue')
      },
      {
        path: 'operation-logs',
        name: 'OperationLogs',
        component: () => import('@/views/OperationLogs.vue')
      },
      {
        path: 'client-filter',
        name: 'ClientFilter',
        component: () => import('@/views/ClientFilter.vue')
      },
      {
        path: 'error-messages',
        name: 'ErrorMessages',
        component: () => import('@/views/ErrorMessages.vue')
      },
      {
        path: 'system-logs',
        name: 'SystemLogs',
        component: () => import('@/views/SystemLogs.vue')
      },
      {
        path: 'usage',
        name: 'Usage',
        component: () => import('@/views/Usage.vue')
      }
    ]
  },
  // 旧路由兼容重定向
  {
    path: '/system-monitor',
    redirect: '/admin/system-monitor'
  },
  {
    path: '/accounts',
    redirect: '/admin/accounts'
  },
  {
    path: '/models',
    redirect: '/admin/models'
  },
  {
    path: '/api-keys',
    redirect: '/admin/api-keys'
  },
  {
    path: '/request-logs',
    redirect: '/admin/request-logs'
  },
  {
    path: '/settings',
    redirect: '/admin/settings'
  },
  // 404 处理
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  // 需要登录但未登录
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
    return
  }

  // 已登录访问登录页，跳转到后台
  if (to.path === '/login' && userStore.isLoggedIn) {
    next('/admin/system-monitor')
    return
  }

  next()
})

export default router
