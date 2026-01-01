/**
 * 文件作用：Vue应用入口，初始化Vue应用和全局插件
 * 负责功能：
 *   - Vue应用创建
 *   - 路由和状态管理初始化
 *   - Apple HIG 设计系统样式加载
 * 重要程度：⭐⭐⭐⭐ 重要（前端入口）
 * 依赖模块：vue, pinia, router
 */
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './styles/index.css'

import App from './App.vue'
import router from './router'

const app = createApp(App)

app.use(createPinia())
app.use(router)

try {
  app.mount('#app')
  window.__APP_BOOTED__ = true
} catch (e) {
  // eslint-disable-next-line no-console
  console.error('前端启动失败：', e)
  window.__APP_BOOT_ERROR__ = String((e && e.message) || e)
}
