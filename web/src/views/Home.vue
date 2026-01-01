<!--
 * 文件作用：系统首页 - Apple HIG 风格
 * 负责功能：
 *   - API 接口配置说明
 *   - 登录后台入口
 * 重要程度：⭐⭐⭐⭐⭐ 核心
-->
<template>
  <div class="home-page">
    <!-- 导航栏 -->
    <header class="navbar">
      <div class="navbar-content">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
          <span class="logo-text">Cli-Proxy</span>
        </div>
        <div class="nav-actions">
          <router-link to="/register" class="register-btn">
            <span>注册</span>
          </router-link>
          <router-link to="/login" class="login-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 3h4a2 2 0 012 2v14a2 2 0 01-2 2h-4"/>
              <polyline points="10,17 15,12 10,7"/>
              <line x1="15" y1="12" x2="3" y2="12"/>
            </svg>
            <span>登录</span>
          </router-link>
        </div>
      </div>
    </header>

    <!-- 主内容 -->
    <main class="main-content">
      <!-- Hero 区域 -->
      <section class="hero">
        <h1 class="hero-title">统一的 AI API 接入</h1>
        <p class="hero-subtitle">支持 Claude、OpenAI、Gemini 多平台，一个端点搞定所有</p>
        <div class="hero-actions">
          <router-link to="/register" class="hero-btn primary">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
              <circle cx="8.5" cy="7" r="4"/>
              <line x1="20" y1="8" x2="20" y2="14"/>
              <line x1="23" y1="11" x2="17" y2="11"/>
            </svg>
            免费注册
          </router-link>
          <router-link to="/login" class="hero-btn secondary">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 3h4a2 2 0 012 2v14a2 2 0 01-2 2h-4"/>
              <polyline points="10,17 15,12 10,7"/>
              <line x1="15" y1="12" x2="3" y2="12"/>
            </svg>
            登录后台
          </router-link>
        </div>
      </section>

      <!-- API 端点卡片 -->
      <section class="endpoints-section">
        <h2 class="section-title">服务端点</h2>
        <div class="endpoints-grid">
          <!-- Claude -->
          <div class="endpoint-card">
            <div class="endpoint-header">
              <div class="platform-badge claude">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M8 12h8M12 8v8"/>
                </svg>
              </div>
              <span class="platform-name">Claude</span>
            </div>
            <div class="endpoint-url">
              <code>{{ baseUrl }}/claude</code>
              <button class="copy-btn" @click="copyText(baseUrl + '/claude')" :class="{ copied: copiedUrl === baseUrl + '/claude' }">
                <svg v-if="copiedUrl !== baseUrl + '/claude'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20,6 9,17 4,12"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- OpenAI -->
          <div class="endpoint-card">
            <div class="endpoint-header">
              <div class="platform-badge openai">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M12 6v12M6 12h12"/>
                </svg>
              </div>
              <span class="platform-name">OpenAI</span>
            </div>
            <div class="endpoint-url">
              <code>{{ baseUrl }}/openai</code>
              <button class="copy-btn" @click="copyText(baseUrl + '/openai')" :class="{ copied: copiedUrl === baseUrl + '/openai' }">
                <svg v-if="copiedUrl !== baseUrl + '/openai'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20,6 9,17 4,12"/>
                </svg>
              </button>
            </div>
          </div>

          <!-- Gemini -->
          <div class="endpoint-card">
            <div class="endpoint-header">
              <div class="platform-badge gemini">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <polygon points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"/>
                </svg>
              </div>
              <span class="platform-name">Gemini</span>
            </div>
            <div class="endpoint-url">
              <code>{{ baseUrl }}/gemini</code>
              <button class="copy-btn" @click="copyText(baseUrl + '/gemini')" :class="{ copied: copiedUrl === baseUrl + '/gemini' }">
                <svg v-if="copiedUrl !== baseUrl + '/gemini'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20,6 9,17 4,12"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- 配置示例 -->
      <section class="config-section">
        <h2 class="section-title">快速配置</h2>
        <div class="config-card">
          <div class="config-tabs">
            <button
              v-for="tab in configTabs"
              :key="tab.id"
              :class="['config-tab', { active: activeTab === tab.id }]"
              @click="activeTab = tab.id"
            >
              {{ tab.name }}
            </button>
          </div>
          <div class="config-content">
            <pre><code>{{ configCode }}</code></pre>
            <button class="copy-config-btn" @click="copyText(configCode)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
              </svg>
              复制代码
            </button>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from '@/utils/toast'

const baseUrl = ref(window.location.origin)
const copiedUrl = ref('')
const activeTab = ref('claude')

const configTabs = [
  { id: 'claude', name: 'Claude' },
  { id: 'openai', name: 'OpenAI' },
  { id: 'bash', name: 'Bash' }
]

const configCode = computed(() => {
  const configs = {
    claude: `# Claude 配置
export ANTHROPIC_BASE_URL="${baseUrl.value}/claude"
export ANTHROPIC_API_KEY="your-api-key"`,
    openai: `# OpenAI 配置
export OPENAI_API_BASE="${baseUrl.value}/openai/v1"
export OPENAI_API_KEY="your-api-key"`,
    bash: `# 测试连接
curl ${baseUrl.value}/claude/v1/messages \\
  -H "Content-Type: application/json" \\
  -H "x-api-key: your-api-key" \\
  -d '{"model":"claude-3-5-sonnet-20241022","max_tokens":100,"messages":[{"role":"user","content":"Hello"}]}'`
  }
  return configs[activeTab.value]
})

const copyText = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    copiedUrl.value = text
    ElMessage.success('已复制到剪贴板')
    setTimeout(() => {
      copiedUrl.value = ''
    }, 2000)
  } catch {
    ElMessage.error('复制失败')
  }
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  background: var(--apple-bg-secondary);
}

/* 导航栏 */
.navbar {
  position: sticky;
  top: 0;
  background: var(--apple-bg-blur);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid var(--apple-separator);
  z-index: 100;
}

.navbar-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--apple-spacing-xl);
  height: 56px;
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.register-btn {
  padding: var(--apple-spacing-xs) var(--apple-spacing-md);
  color: var(--apple-blue);
  border-radius: var(--apple-radius-full);
  text-decoration: none;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.register-btn:hover {
  background: var(--apple-blue-light);
}

.login-btn {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-xs) var(--apple-spacing-md);
  background: var(--apple-blue);
  color: white;
  border-radius: var(--apple-radius-full);
  text-decoration: none;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.login-btn svg {
  width: 16px;
  height: 16px;
}

.login-btn:hover {
  background: var(--apple-blue-hover);
  transform: scale(1.02);
}

/* 主内容 */
.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: var(--apple-spacing-xxl) var(--apple-spacing-xl);
}

/* Hero */
.hero {
  text-align: center;
  padding: var(--apple-spacing-xxxl) 0;
}

.hero-title {
  font-size: var(--apple-text-5xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  letter-spacing: var(--apple-tracking-tight);
  margin-bottom: var(--apple-spacing-md);
}

.hero-subtitle {
  font-size: var(--apple-text-xl);
  color: var(--apple-text-secondary);
  max-width: 600px;
  margin: 0 auto var(--apple-spacing-xl);
}

.hero-actions {
  display: flex;
  justify-content: center;
  gap: var(--apple-spacing-md);
  flex-wrap: wrap;
}

.hero-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-md) var(--apple-spacing-xl);
  border-radius: var(--apple-radius-full);
  text-decoration: none;
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-semibold);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.hero-btn svg {
  width: 20px;
  height: 20px;
}

.hero-btn.primary {
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  color: white;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
}

.hero-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(102, 126, 234, 0.4);
}

.hero-btn.secondary {
  background: var(--apple-bg-primary);
  color: var(--apple-text-primary);
  box-shadow: var(--apple-shadow-card);
}

.hero-btn.secondary:hover {
  transform: translateY(-2px);
  box-shadow: var(--apple-shadow-card-hover);
}

/* 端点区域 */
.endpoints-section {
  margin-bottom: var(--apple-spacing-xxxl);
}

.section-title {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xl);
}

.endpoints-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: var(--apple-spacing-lg);
}

.endpoint-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  padding: var(--apple-spacing-xl);
  box-shadow: var(--apple-shadow-card);
  transition: all var(--apple-duration-normal) var(--apple-ease-default);
}

.endpoint-card:hover {
  box-shadow: var(--apple-shadow-card-hover);
  transform: translateY(-2px);
}

.endpoint-header {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  margin-bottom: var(--apple-spacing-lg);
}

.platform-badge {
  width: 44px;
  height: 44px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
}

.platform-badge svg {
  width: 24px;
  height: 24px;
  color: white;
}

.platform-badge.claude {
  background: linear-gradient(135deg, #d97706 0%, #ea580c 100%);
}

.platform-badge.openai {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.platform-badge.gemini {
  background: linear-gradient(135deg, #3b82f6 0%, #6366f1 100%);
}

.platform-name {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.endpoint-url {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  background: var(--apple-bg-tertiary);
  border-radius: var(--apple-radius-md);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
}

.endpoint-url code {
  flex: 1;
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.copy-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--apple-radius-sm);
  color: var(--apple-text-secondary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  flex-shrink: 0;
}

.copy-btn svg {
  width: 16px;
  height: 16px;
}

.copy-btn:hover {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
}

.copy-btn.copied {
  color: var(--apple-green);
}

/* 配置区域 */
.config-section {
  margin-bottom: var(--apple-spacing-xxxl);
}

.config-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  overflow: hidden;
  box-shadow: var(--apple-shadow-card);
}

.config-tabs {
  display: flex;
  padding: var(--apple-spacing-xs);
  background: var(--apple-bg-tertiary);
  gap: var(--apple-spacing-xxs);
}

.config-tab {
  flex: 1;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.config-tab:hover {
  color: var(--apple-text-primary);
}

.config-tab.active {
  background: var(--apple-bg-primary);
  color: var(--apple-text-primary);
  box-shadow: var(--apple-shadow-sm);
}

.config-content {
  position: relative;
  padding: var(--apple-spacing-xl);
}

.config-content pre {
  margin: 0;
  background: var(--apple-bg-tertiary);
  border-radius: var(--apple-radius-md);
  padding: var(--apple-spacing-lg);
  overflow-x: auto;
}

.config-content code {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  line-height: 1.6;
}

.copy-config-btn {
  position: absolute;
  top: var(--apple-spacing-xl);
  right: var(--apple-spacing-xl);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
  box-shadow: var(--apple-shadow-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.copy-config-btn svg {
  width: 14px;
  height: 14px;
}

.copy-config-btn:hover {
  color: var(--apple-text-primary);
  box-shadow: var(--apple-shadow-md);
}

/* 响应式 */
@media (max-width: 768px) {
  .nav-actions {
    gap: var(--apple-spacing-xs);
  }

  .register-btn {
    padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  }

  .hero-title {
    font-size: var(--apple-text-3xl);
  }

  .hero-subtitle {
    font-size: var(--apple-text-md);
  }

  .hero-actions {
    flex-direction: column;
    align-items: center;
  }

  .hero-btn {
    width: 100%;
    max-width: 280px;
    justify-content: center;
  }

  .endpoints-grid {
    grid-template-columns: 1fr;
  }
}
</style>
