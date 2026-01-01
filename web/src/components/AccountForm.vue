<!--
 * 文件作用：账户表单组件 - Apple HIG 风格
 * 负责功能：
 *   - 多平台选择和配置
 *   - OAuth/SessionKey/API Key授权方式
 *   - 基本信息和代理配置
 *   - 模型限制和映射配置
 * 重要程度：⭐⭐⭐⭐ 重要（账户管理核心）
-->
<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @click.self="handleClose">
        <div class="modal-container">
          <!-- 模态框头部 -->
          <div class="modal-header">
            <h2 class="modal-title">{{ isEdit ? '编辑账户' : '添加账户' }}</h2>
            <button class="modal-close" @click="handleClose">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>

          <!-- 步骤指示器 -->
          <div v-if="!isEdit && needsOAuth" class="steps-indicator">
            <div class="step" :class="{ active: step >= 1, done: step > 1 }">
              <div class="step-num">
                <svg v-if="step > 1" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
                <span v-else>1</span>
              </div>
              <span class="step-label">基本信息</span>
            </div>
            <div class="step-line" :class="{ done: step > 1 }"></div>
            <div class="step" :class="{ active: step >= 2 }">
              <div class="step-num">2</div>
              <span class="step-label">授权认证</span>
            </div>
          </div>

          <!-- 模态框内容 -->
          <div class="modal-body">
            <!-- 步骤1: 基本信息 -->
            <div v-if="step === 1 && !isEdit" class="form-step">
              <!-- 平台选择 -->
              <div class="form-section">
                <h4 class="section-title">选择平台</h4>
                <div class="platform-groups">
                  <div
                    v-for="group in platformGroups"
                    :key="group.key"
                    class="platform-card"
                    :class="{ selected: platformGroup === group.key }"
                    @click="selectPlatformGroup(group.key)"
                  >
                    <div class="platform-icon" :style="{ background: group.gradient }">
                      <component :is="group.iconComponent" />
                    </div>
                    <div class="platform-info">
                      <h5>{{ group.name }}</h5>
                      <p>{{ group.desc }}</p>
                    </div>
                    <div v-if="platformGroup === group.key" class="check-badge">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                        <polyline points="20 6 9 17 4 12"/>
                      </svg>
                    </div>
                  </div>
                </div>

                <!-- 子平台选择 -->
                <div v-if="platformGroup" class="subplatform-section">
                  <p class="subplatform-label">选择具体平台类型：</p>
                  <div class="subplatform-grid">
                    <label
                      v-for="sub in currentSubplatforms"
                      :key="sub.value"
                      class="subplatform-card"
                      :class="{ selected: form.type === sub.value }"
                    >
                      <input v-model="form.type" type="radio" :value="sub.value" class="sr-only" />
                      <div class="sub-icon" :style="{ background: sub.color }">
                        <component :is="sub.iconComponent" />
                      </div>
                      <div class="sub-info">
                        <span class="sub-name">{{ sub.label }}</span>
                        <span class="sub-desc">{{ sub.desc }}</span>
                      </div>
                      <div v-if="form.type === sub.value" class="check-badge small">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </div>
                    </label>
                  </div>
                </div>
              </div>

              <!-- 添加方式选择 -->
              <div v-if="showAddTypeSelection" class="form-section">
                <h4 class="section-title">添加方式</h4>
                <div class="add-type-grid">
                  <!-- Claude 平台 -->
                  <template v-if="form.type === 'claude-official'">
                    <div
                      class="type-card"
                      :class="{ selected: form.addType === 'cookie' }"
                      @click="form.addType = 'cookie'"
                    >
                      <div class="type-icon cookie">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <circle cx="8" cy="10" r="1" fill="currentColor"/>
                          <circle cx="15" cy="8" r="1" fill="currentColor"/>
                          <circle cx="10" cy="15" r="1" fill="currentColor"/>
                          <circle cx="16" cy="14" r="1" fill="currentColor"/>
                        </svg>
                      </div>
                      <div class="type-content">
                        <h5>SessionKey 授权</h5>
                        <p>使用 Claude.ai 的 sessionKey 自动完成授权</p>
                      </div>
                      <div v-if="form.addType === 'cookie'" class="check-badge small">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </div>
                    </div>
                    <div
                      class="type-card"
                      :class="{ selected: form.addType === 'oauth' }"
                      @click="form.addType = 'oauth'"
                    >
                      <div class="type-icon oauth">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                          <path d="M7 11V7a5 5 0 0110 0v4"/>
                        </svg>
                      </div>
                      <div class="type-content">
                        <h5>OAuth 手动授权</h5>
                        <p>生成授权链接，手动完成授权流程</p>
                      </div>
                      <div v-if="form.addType === 'oauth'" class="check-badge small">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </div>
                    </div>
                  </template>

                  <!-- ChatGPT 官方平台 -->
                  <template v-else-if="form.type === 'openai-responses'">
                    <div
                      class="type-card"
                      :class="{ selected: form.addType === 'oauth' }"
                      @click="form.addType = 'oauth'"
                    >
                      <div class="type-icon oauth">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                        </svg>
                      </div>
                      <div class="type-content">
                        <h5>OAuth 授权</h5>
                        <p>通过 OAuth 授权访问 ChatGPT 账户</p>
                      </div>
                      <div v-if="form.addType === 'oauth'" class="check-badge small">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </div>
                    </div>
                    <div
                      class="type-card"
                      :class="{ selected: form.addType === 'cookie' }"
                      @click="form.addType = 'cookie'"
                    >
                      <div class="type-icon cookie">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <circle cx="8" cy="10" r="1" fill="currentColor"/>
                          <circle cx="15" cy="8" r="1" fill="currentColor"/>
                        </svg>
                      </div>
                      <div class="type-content">
                        <h5>SessionKey 授权</h5>
                        <p>使用 ChatGPT 的 SessionKey 直接认证</p>
                      </div>
                      <div v-if="form.addType === 'cookie'" class="check-badge small">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </div>
                    </div>
                  </template>

                  <!-- Gemini 平台 -->
                  <template v-else-if="form.type === 'gemini'">
                    <div
                      class="type-card"
                      :class="{ selected: form.addType === 'apikey' }"
                      @click="form.addType = 'apikey'"
                    >
                      <div class="type-icon apikey">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                        </svg>
                      </div>
                      <div class="type-content">
                        <h5>API Key</h5>
                        <p>使用 Gemini API Key 直接连接</p>
                      </div>
                      <div v-if="form.addType === 'apikey'" class="check-badge small">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </div>
                    </div>
                    <div
                      class="type-card"
                      :class="{ selected: form.addType === 'oauth' }"
                      @click="form.addType = 'oauth'"
                    >
                      <div class="type-icon oauth">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                        </svg>
                      </div>
                      <div class="type-content">
                        <h5>OAuth 授权</h5>
                        <p>通过 Google OAuth 授权访问</p>
                      </div>
                      <div v-if="form.addType === 'oauth'" class="check-badge small">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                          <polyline points="20 6 9 17 4 12"/>
                        </svg>
                      </div>
                    </div>
                  </template>
                </div>
              </div>

              <!-- 基本信息表单 -->
              <div class="form-section">
                <h4 class="section-title">基本信息</h4>
                <div class="form-grid">
                  <div class="form-group full">
                    <label class="form-label required">账户名称</label>
                    <input
                      v-model="form.name"
                      type="text"
                      class="form-input"
                      placeholder="为账户设置一个易识别的名称"
                    />
                  </div>

                  <div class="form-group full">
                    <label class="form-label">描述（可选）</label>
                    <textarea
                      v-model="form.description"
                      class="form-textarea"
                      rows="2"
                      placeholder="账户用途说明..."
                    ></textarea>
                  </div>

                  <div class="form-group">
                    <label class="form-label">账户类型</label>
                    <select v-model="form.accountType" class="form-select">
                      <option value="shared">共享账户</option>
                      <option value="dedicated">专属账户</option>
                    </select>
                  </div>

                  <div class="form-group">
                    <label class="form-label">优先级</label>
                    <input
                      v-model.number="form.priority"
                      type="number"
                      min="1"
                      max="100"
                      class="form-input"
                    />
                  </div>

                  <div class="form-group">
                    <label class="form-label">最大并发</label>
                    <input
                      v-model.number="form.max_concurrency"
                      type="number"
                      min="1"
                      max="100"
                      class="form-input"
                    />
                  </div>

                  <div class="form-group">
                    <label class="form-label">每日预算($)</label>
                    <input
                      v-model.number="form.daily_budget"
                      type="number"
                      min="0"
                      max="10000"
                      step="0.01"
                      class="form-input"
                    />
                  </div>

                  <div class="form-group">
                    <label class="form-label">启用</label>
                    <label class="toggle-switch">
                      <input type="checkbox" v-model="form.enabled" />
                      <span class="toggle-slider"></span>
                    </label>
                  </div>
                </div>
              </div>

              <!-- 平台特定配置 -->
              <div v-if="showPlatformConfig" class="form-section">
                <h4 class="section-title">{{ getPlatformConfigTitle }}</h4>

                <!-- Claude Console 配置 -->
                <template v-if="form.type === 'claude-console'">
                  <div class="form-grid">
                    <div class="form-group full">
                      <label class="form-label required">API URL</label>
                      <input v-model="form.api_url" type="text" class="form-input" placeholder="例如：https://api.anthropic.com" />
                    </div>
                    <div class="form-group full">
                      <label class="form-label required">API Key</label>
                      <div class="input-with-toggle">
                        <input
                          v-model="form.api_key"
                          :type="showApiKey ? 'text' : 'password'"
                          class="form-input"
                          placeholder="sk-ant-api03-..."
                        />
                        <button type="button" class="toggle-visibility" @click="showApiKey = !showApiKey">
                          <svg v-if="showApiKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                </template>

                <!-- Bedrock 配置 -->
                <template v-if="form.type === 'bedrock'">
                  <div class="form-grid">
                    <div class="form-group full">
                      <label class="form-label required">AWS Access Key ID</label>
                      <input v-model="form.aws_access_key_id" type="text" class="form-input" placeholder="请输入 AWS Access Key ID" />
                    </div>
                    <div class="form-group full">
                      <label class="form-label required">AWS Secret Access Key</label>
                      <div class="input-with-toggle">
                        <input
                          v-model="form.aws_secret_access_key"
                          :type="showAwsSecret ? 'text' : 'password'"
                          class="form-input"
                          placeholder="请输入 AWS Secret Access Key"
                        />
                        <button type="button" class="toggle-visibility" @click="showAwsSecret = !showAwsSecret">
                          <svg v-if="showAwsSecret" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                            <line x1="1" y1="1" x2="23" y2="23"/>
                          </svg>
                          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                            <circle cx="12" cy="12" r="3"/>
                          </svg>
                        </button>
                      </div>
                    </div>
                    <div class="form-group full">
                      <label class="form-label required">AWS Region</label>
                      <select v-model="form.aws_region" class="form-select">
                        <option value="us-east-1">us-east-1 (美国东部)</option>
                        <option value="us-west-2">us-west-2 (美国西部)</option>
                        <option value="eu-west-1">eu-west-1 (欧洲爱尔兰)</option>
                        <option value="ap-northeast-1">ap-northeast-1 (东京)</option>
                        <option value="ap-southeast-1">ap-southeast-1 (新加坡)</option>
                      </select>
                    </div>
                  </div>
                </template>

                <!-- OpenAI 配置 -->
                <template v-if="form.type === 'openai' && form.addType === 'apikey'">
                  <div class="form-grid">
                    <div class="form-group full">
                      <label class="form-label required">
                        <svg class="label-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                        </svg>
                        API Key
                      </label>
                      <div class="input-with-toggle">
                        <input
                          v-model="form.api_key"
                          :type="showApiKey ? 'text' : 'password'"
                          class="form-input"
                          placeholder="sk-..."
                        />
                        <button type="button" class="toggle-visibility" @click="showApiKey = !showApiKey">
                          <svg v-if="showApiKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                            <line x1="1" y1="1" x2="23" y2="23"/>
                          </svg>
                          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                            <circle cx="12" cy="12" r="3"/>
                          </svg>
                        </button>
                      </div>
                      <p class="form-tip">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <line x1="12" y1="16" x2="12" y2="12"/>
                          <line x1="12" y1="8" x2="12.01" y2="8"/>
                        </svg>
                        在 <a href="https://platform.openai.com/api-keys" target="_blank">OpenAI Platform</a> 获取 API Key
                      </p>
                    </div>
                    <div class="form-group full">
                      <label class="form-label">组织 ID（可选）</label>
                      <input v-model="form.organization_id" type="text" class="form-input" placeholder="org-..." />
                    </div>
                    <div class="form-group full">
                      <label class="form-label">API Base URL（可选）</label>
                      <input v-model="form.api_url" type="text" class="form-input" placeholder="默认: https://api.openai.com/v1" />
                      <p class="form-tip">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <line x1="12" y1="16" x2="12" y2="12"/>
                          <line x1="12" y1="8" x2="12.01" y2="8"/>
                        </svg>
                        如果使用代理或自定义端点，请填写完整 URL
                      </p>
                    </div>
                  </div>
                </template>

                <!-- ChatGPT 官方 SessionKey 配置 -->
                <template v-if="form.type === 'openai-responses' && form.addType === 'cookie'">
                  <div class="form-grid">
                    <div class="form-group full">
                      <label class="form-label required">
                        <svg class="label-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <circle cx="8" cy="10" r="1" fill="currentColor"/>
                          <circle cx="15" cy="8" r="1" fill="currentColor"/>
                        </svg>
                        SessionKey (Access Token)
                      </label>
                      <div class="input-with-toggle">
                        <input
                          v-model="form.session_key"
                          :type="showSessionKey ? 'text' : 'password'"
                          class="form-input"
                          placeholder="eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
                        />
                        <button type="button" class="toggle-visibility" @click="showSessionKey = !showSessionKey">
                          <svg v-if="showSessionKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
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
                  <div class="help-card">
                    <div class="help-header">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"/>
                        <line x1="12" y1="16" x2="12" y2="12"/>
                        <line x1="12" y1="8" x2="12.01" y2="8"/>
                      </svg>
                      <span>如何获取 SessionKey</span>
                    </div>
                    <ol class="help-steps">
                      <li>在浏览器中登录 <a href="https://chatgpt.com" target="_blank">chatgpt.com</a></li>
                      <li>按 <kbd>F12</kbd> 打开开发者工具，切换到 <strong>Network</strong> 标签</li>
                      <li>刷新页面，找到任意 API 请求（如 <code>conversations</code>）</li>
                      <li>在请求头中找到 <strong>Authorization: Bearer xxx</strong></li>
                      <li>复制 <strong>Bearer</strong> 后面的完整 token</li>
                    </ol>
                    <p class="help-note">SessionKey 通常以 <code>eyJhbGciOiJSUzI1NiI</code> 开头，有效期约 7 天</p>
                  </div>
                </template>

                <!-- Azure OpenAI 配置 -->
                <template v-if="form.type === 'azure-openai'">
                  <div class="form-grid">
                    <div class="form-group full">
                      <label class="form-label required">Azure Endpoint</label>
                      <input v-model="form.azure_endpoint" type="text" class="form-input" placeholder="https://your-resource.openai.azure.com" />
                    </div>
                    <div class="form-group full">
                      <label class="form-label required">API Key</label>
                      <div class="input-with-toggle">
                        <input
                          v-model="form.api_key"
                          :type="showApiKey ? 'text' : 'password'"
                          class="form-input"
                          placeholder="Azure API Key"
                        />
                        <button type="button" class="toggle-visibility" @click="showApiKey = !showApiKey">
                          <svg v-if="showApiKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                            <line x1="1" y1="1" x2="23" y2="23"/>
                          </svg>
                          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                            <circle cx="12" cy="12" r="3"/>
                          </svg>
                        </button>
                      </div>
                    </div>
                    <div class="form-group">
                      <label class="form-label required">部署名称</label>
                      <input v-model="form.azure_deployment_name" type="text" class="form-input" placeholder="gpt-4" />
                    </div>
                    <div class="form-group">
                      <label class="form-label">API 版本</label>
                      <input v-model="form.azure_api_version" type="text" class="form-input" placeholder="2024-02-01" />
                    </div>
                  </div>
                </template>

                <!-- Gemini 配置 -->
                <template v-if="form.type === 'gemini' && form.addType === 'apikey'">
                  <div class="form-grid">
                    <div class="form-group full">
                      <label class="form-label required">
                        <svg class="label-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                        </svg>
                        API Key
                      </label>
                      <div class="input-with-toggle">
                        <input
                          v-model="form.api_key"
                          :type="showApiKey ? 'text' : 'password'"
                          class="form-input"
                          placeholder="Gemini API Key"
                        />
                        <button type="button" class="toggle-visibility" @click="showApiKey = !showApiKey">
                          <svg v-if="showApiKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                            <line x1="1" y1="1" x2="23" y2="23"/>
                          </svg>
                          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                            <circle cx="12" cy="12" r="3"/>
                          </svg>
                        </button>
                      </div>
                      <p class="form-tip">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <circle cx="12" cy="12" r="10"/>
                          <line x1="12" y1="16" x2="12" y2="12"/>
                          <line x1="12" y1="8" x2="12.01" y2="8"/>
                        </svg>
                        在 <a href="https://aistudio.google.com/apikey" target="_blank">Google AI Studio</a> 获取 API Key
                      </p>
                    </div>
                    <div class="form-group full">
                      <label class="form-label">API Base URL（可选）</label>
                      <input v-model="form.api_url" type="text" class="form-input" placeholder="默认: https://generativelanguage.googleapis.com/v1beta" />
                    </div>
                  </div>
                </template>
              </div>

              <!-- 模型配置 -->
              <div class="form-section">
                <h4 class="section-title">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polygon points="12,2 2,7 12,12 22,7"/>
                    <polyline points="2,17 12,22 22,17"/>
                    <polyline points="2,12 12,17 22,12"/>
                  </svg>
                  模型配置（可选）
                </h4>
                <div class="form-group full">
                  <label class="form-label">允许的模型</label>
                  <div class="tags-input" @click="focusModelInput">
                    <span v-for="model in form.allowedModelsList" :key="model" class="tag">
                      {{ model }}
                      <button type="button" class="tag-remove" @click.stop="removeModel(model)">×</button>
                    </span>
                    <input
                      ref="modelInputRef"
                      v-model="modelInputValue"
                      type="text"
                      class="tag-input"
                      placeholder="输入模型名称后按回车添加"
                      @keydown.enter.prevent="addModel"
                      @keydown.delete="handleModelBackspace"
                    />
                  </div>
                  <p class="form-tip">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="12" cy="12" r="10"/>
                      <line x1="12" y1="16" x2="12" y2="12"/>
                      <line x1="12" y1="8" x2="12.01" y2="8"/>
                    </svg>
                    留空表示允许所有模型，可手动输入自定义模型名
                  </p>
                </div>
              </div>

              <!-- 代理配置 -->
              <div class="form-section">
                <h4 class="section-title">代理配置</h4>
                <div class="form-group full">
                  <label class="form-label">选择代理</label>
                  <select v-model="form.proxy_id" class="form-select">
                    <option :value="null">直连（不使用代理）</option>
                    <option v-for="proxy in proxyList" :key="proxy.id" :value="proxy.id">
                      {{ proxy.name }} ({{ proxy.type }}://{{ proxy.host }}:{{ proxy.port }})
                    </option>
                  </select>
                  <p class="form-tip">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="12" cy="12" r="10"/>
                      <line x1="12" y1="16" x2="12" y2="12"/>
                      <line x1="12" y1="8" x2="12.01" y2="8"/>
                    </svg>
                    留空表示直连，选择代理后将通过该代理访问目标服务
                  </p>
                </div>
              </div>
            </div>

            <!-- 步骤2: OAuth 授权 -->
            <div v-if="step === 2" class="form-step">
              <OAuthFlow
                :platform="form.type"
                :proxy="selectedProxy"
                :auth-mode="form.addType"
                @success="handleOAuthSuccess"
                @back="step = 1"
              />
            </div>

            <!-- 编辑模式 -->
            <div v-if="isEdit" class="form-step">
              <div class="edit-type-banner">
                <div class="type-badge" :style="{ background: getTypeColor(form.type) }">
                  <component :is="getTypeIconComponent(form.type)" />
                </div>
                <span>{{ getTypeLabel(form.type) }}</span>
              </div>

              <div class="form-grid">
                <div class="form-group full">
                  <label class="form-label required">账户名称</label>
                  <input v-model="form.name" type="text" class="form-input" placeholder="账户名称" />
                </div>

                <div class="form-group full">
                  <label class="form-label">描述（可选）</label>
                  <textarea v-model="form.description" class="form-textarea" rows="2" placeholder="账户用途说明..."></textarea>
                </div>

                <div class="form-group">
                  <label class="form-label">优先级</label>
                  <input v-model.number="form.priority" type="number" min="1" max="100" class="form-input" />
                </div>

                <div class="form-group">
                  <label class="form-label">权重</label>
                  <input v-model.number="form.weight" type="number" min="1" max="1000" class="form-input" />
                </div>

                <div class="form-group">
                  <label class="form-label">最大并发</label>
                  <input v-model.number="form.max_concurrency" type="number" min="1" max="100" class="form-input" />
                </div>

                <div class="form-group">
                  <label class="form-label">启用</label>
                  <label class="toggle-switch">
                    <input type="checkbox" v-model="form.enabled" />
                    <span class="toggle-slider"></span>
                  </label>
                </div>
              </div>

              <!-- API 配置 -->
              <div v-if="showEditApiConfig" class="form-section">
                <h4 class="section-title">API 配置</h4>
                <div class="form-grid">
                  <div class="form-group full">
                    <label class="form-label">API URL</label>
                    <input v-model="form.api_url" type="text" class="form-input" placeholder="API Base URL" />
                    <p class="form-tip">留空使用默认地址</p>
                  </div>
                  <div class="form-group full">
                    <label class="form-label">API Key</label>
                    <div class="input-with-toggle">
                      <input
                        v-model="form.api_key"
                        :type="showApiKey ? 'text' : 'password'"
                        class="form-input"
                        placeholder="API Key"
                      />
                      <button type="button" class="toggle-visibility" @click="showApiKey = !showApiKey">
                        <svg v-if="showApiKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                          <line x1="1" y1="1" x2="23" y2="23"/>
                        </svg>
                        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                          <circle cx="12" cy="12" r="3"/>
                        </svg>
                      </button>
                    </div>
                    <p class="form-tip">留空保持原有 Key 不变</p>
                  </div>
                </div>
              </div>

              <!-- AWS Bedrock 配置 -->
              <div v-if="form.type === 'bedrock'" class="form-section">
                <h4 class="section-title">AWS Bedrock 配置</h4>
                <div class="form-grid">
                  <div class="form-group full">
                    <label class="form-label">AWS Access Key ID</label>
                    <input v-model="form.aws_access_key_id" type="text" class="form-input" placeholder="请输入 AWS Access Key ID" />
                  </div>
                  <div class="form-group full">
                    <label class="form-label">AWS Secret Access Key</label>
                    <div class="input-with-toggle">
                      <input
                        v-model="form.aws_secret_access_key"
                        :type="showAwsSecret ? 'text' : 'password'"
                        class="form-input"
                        placeholder="请输入 AWS Secret Access Key（留空保持不变）"
                      />
                      <button type="button" class="toggle-visibility" @click="showAwsSecret = !showAwsSecret">
                        <svg v-if="showAwsSecret" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                          <line x1="1" y1="1" x2="23" y2="23"/>
                        </svg>
                        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                          <circle cx="12" cy="12" r="3"/>
                        </svg>
                      </button>
                    </div>
                  </div>
                  <div class="form-group full">
                    <label class="form-label">AWS Region</label>
                    <select v-model="form.aws_region" class="form-select">
                      <option value="us-east-1">us-east-1 (美国东部)</option>
                      <option value="us-west-2">us-west-2 (美国西部)</option>
                      <option value="eu-west-1">eu-west-1 (欧洲爱尔兰)</option>
                      <option value="ap-northeast-1">ap-northeast-1 (东京)</option>
                      <option value="ap-southeast-1">ap-southeast-1 (新加坡)</option>
                    </select>
                  </div>
                </div>
              </div>

              <!-- Azure OpenAI 编辑配置 -->
              <div v-if="form.type === 'azure-openai'" class="form-section">
                <h4 class="section-title">Azure OpenAI 配置</h4>
                <div class="form-grid">
                  <div class="form-group full">
                    <label class="form-label">Azure Endpoint</label>
                    <input v-model="form.azure_endpoint" type="text" class="form-input" placeholder="https://your-resource.openai.azure.com" />
                  </div>
                  <div class="form-group full">
                    <label class="form-label">API Key</label>
                    <div class="input-with-toggle">
                      <input
                        v-model="form.api_key"
                        :type="showApiKey ? 'text' : 'password'"
                        class="form-input"
                        placeholder="Azure API Key（留空保持不变）"
                      />
                      <button type="button" class="toggle-visibility" @click="showApiKey = !showApiKey">
                        <svg v-if="showApiKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                          <line x1="1" y1="1" x2="23" y2="23"/>
                        </svg>
                        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                          <circle cx="12" cy="12" r="3"/>
                        </svg>
                      </button>
                    </div>
                  </div>
                  <div class="form-group">
                    <label class="form-label">部署名称</label>
                    <input v-model="form.azure_deployment_name" type="text" class="form-input" placeholder="gpt-4" />
                  </div>
                  <div class="form-group">
                    <label class="form-label">API 版本</label>
                    <input v-model="form.azure_api_version" type="text" class="form-input" placeholder="2024-02-01" />
                  </div>
                </div>
              </div>

              <!-- ChatGPT 官方编辑配置 -->
              <div v-if="form.type === 'openai-responses'" class="form-section">
                <h4 class="section-title">ChatGPT 官方配置</h4>
                <div class="form-grid">
                  <div class="form-group full">
                    <label class="form-label">Base URL</label>
                    <input v-model="form.api_url" type="text" class="form-input" placeholder="默认: https://chatgpt.com/backend-api/codex" />
                    <p class="form-tip">留空使用默认地址</p>
                  </div>
                  <div class="form-group full">
                    <label class="form-label">SessionKey (Access Token)</label>
                    <div class="input-with-toggle">
                      <input
                        v-model="form.session_key"
                        :type="showSessionKey ? 'text' : 'password'"
                        class="form-input"
                        placeholder="eyJhbGciOiJSUzI1NiI...（留空保持不变）"
                      />
                      <button type="button" class="toggle-visibility" @click="showSessionKey = !showSessionKey">
                        <svg v-if="showSessionKey" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94"/>
                          <line x1="1" y1="1" x2="23" y2="23"/>
                        </svg>
                        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                          <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                          <circle cx="12" cy="12" r="3"/>
                        </svg>
                      </button>
                    </div>
                    <p class="form-tip">留空保持原有 SessionKey 不变</p>
                  </div>
                  <div class="form-group full">
                    <label class="form-label">Organization ID（可选）</label>
                    <input v-model="form.organization_id" type="text" class="form-input" placeholder="org-..." />
                  </div>
                </div>
              </div>

              <!-- 模型配置 -->
              <div class="form-section">
                <h4 class="section-title">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polygon points="12,2 2,7 12,12 22,7"/>
                    <polyline points="2,17 12,22 22,17"/>
                    <polyline points="2,12 12,17 22,12"/>
                  </svg>
                  模型配置
                </h4>
                <div class="form-group full">
                  <label class="form-label">允许的模型</label>
                  <div class="tags-input" @click="focusModelInput">
                    <span v-for="model in form.allowedModelsList" :key="model" class="tag">
                      {{ model }}
                      <button type="button" class="tag-remove" @click.stop="removeModel(model)">×</button>
                    </span>
                    <input
                      ref="modelInputRef"
                      v-model="modelInputValue"
                      type="text"
                      class="tag-input"
                      placeholder="输入模型名称后按回车添加"
                      @keydown.enter.prevent="addModel"
                      @keydown.delete="handleModelBackspace"
                    />
                  </div>
                  <p class="form-tip">留空表示允许所有模型，可手动输入自定义模型名</p>
                </div>
                <div class="form-group full">
                  <label class="form-label">模型映射（可选）</label>
                  <select v-model="form.selectedMappingIds" multiple class="form-select multi">
                    <option v-for="mapping in globalMappings" :key="mapping.id" :value="mapping.id">
                      {{ mapping.source_model }} → {{ mapping.target_model }}
                    </option>
                  </select>
                  <p class="form-tip">从全局模型映射中选择要应用的规则（按住 Ctrl/Cmd 多选）</p>
                </div>
              </div>

              <!-- 代理配置 -->
              <div class="form-section">
                <h4 class="section-title">代理配置</h4>
                <div class="form-group full">
                  <label class="form-label">选择代理</label>
                  <select v-model="form.proxy_id" class="form-select">
                    <option :value="null">直连（不使用代理）</option>
                    <option v-for="proxy in proxyList" :key="proxy.id" :value="proxy.id">
                      {{ proxy.name }} ({{ proxy.type }}://{{ proxy.host }}:{{ proxy.port }})
                    </option>
                  </select>
                  <p class="form-tip">留空表示直连，选择代理后将通过该代理访问目标服务</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 模态框底部 -->
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="handleClose">取消</button>
            <template v-if="!isEdit && step === 1">
              <button
                type="button"
                class="btn btn-primary"
                :disabled="!canProceed"
                @click="handleNext"
              >
                {{ needsOAuth ? '下一步' : '创建账户' }}
                <svg v-if="needsOAuth" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="5" y1="12" x2="19" y2="12"/>
                  <polyline points="12 5 19 12 12 19"/>
                </svg>
              </button>
            </template>
            <template v-if="isEdit">
              <button type="button" class="btn btn-primary" :disabled="submitting" @click="handleSubmit">
                <span v-if="submitting" class="btn-loading"></span>
                {{ submitting ? '保存中...' : '保存修改' }}
              </button>
            </template>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed, watch, h } from 'vue'
import OAuthFlow from './OAuthFlow.vue'
import api from '@/api'

// 图标组件
const BrainIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { d: 'M12 2a10 10 0 1010 10A10 10 0 0012 2zm0 18a8 8 0 118-8 8 8 0 01-8 8z' }),
  h('circle', { cx: '12', cy: '12', r: '3' })
])

const RobotIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('rect', { x: '3', y: '8', width: '18', height: '12', rx: '2' }),
  h('path', { d: 'M12 2v4' }),
  h('circle', { cx: '8', cy: '14', r: '2' }),
  h('circle', { cx: '16', cy: '14', r: '2' })
])

const GoogleIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('circle', { cx: '12', cy: '12', r: '10' }),
  h('path', { d: 'M12 8v8M8 12h8' })
])

const KeyIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { d: 'M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4' })
])

const TerminalIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('rect', { x: '2', y: '3', width: '20', height: '18', rx: '2' }),
  h('polyline', { points: '6 9 10 13 6 17' }),
  h('line', { x1: '14', y1: '17', x2: '18', y2: '17' })
])

const AwsIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { d: 'M12 4L4 8v8l8 4 8-4V8l-8-4z' }),
  h('path', { d: 'M4 8l8 4 8-4M12 12v8' })
])

const BoltIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('polygon', { points: '13 2 3 14 12 14 11 22 21 10 12 10 13 2' })
])

const ChatIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('path', { d: 'M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z' })
])

const MicrosoftIcon = () => h('svg', { viewBox: '0 0 24 24', fill: 'none', stroke: 'currentColor', 'stroke-width': '2' }, [
  h('rect', { x: '3', y: '3', width: '8', height: '8' }),
  h('rect', { x: '13', y: '3', width: '8', height: '8' }),
  h('rect', { x: '3', y: '13', width: '8', height: '8' }),
  h('rect', { x: '13', y: '13', width: '8', height: '8' })
])

const props = defineProps({
  modelValue: Boolean,
  editData: Object
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const isEdit = computed(() => !!props.editData?.id)

const step = ref(1)
const submitting = ref(false)
const platformGroup = ref('')
const showApiKey = ref(false)
const showAwsSecret = ref(false)
const showSessionKey = ref(false)
const modelInputRef = ref(null)
const modelInputValue = ref('')

// 平台分组定义
const platformGroups = [
  {
    key: 'claude',
    name: 'Claude',
    desc: 'Anthropic AI',
    iconComponent: BrainIcon,
    gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    key: 'openai',
    name: 'OpenAI',
    desc: 'GPT 系列模型',
    iconComponent: RobotIcon,
    gradient: 'linear-gradient(135deg, #11998e 0%, #38ef7d 100%)'
  },
  {
    key: 'gemini',
    name: 'Gemini',
    desc: 'Google AI',
    iconComponent: GoogleIcon,
    gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  }
]

// 子平台定义
const subplatformMap = {
  claude: [
    { value: 'claude-official', label: 'Claude Official', desc: 'OAuth 认证', iconComponent: KeyIcon, color: '#667eea' },
    { value: 'claude-console', label: 'Claude Console', desc: 'API Key 认证', iconComponent: TerminalIcon, color: '#764ba2' },
    { value: 'bedrock', label: 'AWS Bedrock', desc: 'AWS 托管服务', iconComponent: AwsIcon, color: '#ff9900' }
  ],
  openai: [
    { value: 'openai', label: 'OpenAI 三方 API', desc: 'API Key 认证', iconComponent: BoltIcon, color: '#11998e' },
    { value: 'openai-responses', label: 'ChatGPT 官方', desc: 'OAuth / SessionKey', iconComponent: ChatIcon, color: '#38ef7d' },
    { value: 'azure-openai', label: 'Azure OpenAI', desc: 'Azure 托管', iconComponent: MicrosoftIcon, color: '#0078d4' }
  ],
  gemini: [
    { value: 'gemini', label: 'Gemini', desc: 'Google AI Studio', iconComponent: GoogleIcon, color: '#4facfe' }
  ]
}

const currentSubplatforms = computed(() => subplatformMap[platformGroup.value] || [])

// 代理列表
const proxyList = ref([])
const loadingProxies = ref(false)

// 全局模型映射列表
const globalMappings = ref([])
const loadingMappings = ref(false)

// 根据 proxy_id 获取选中的代理对象
const selectedProxy = computed(() => {
  if (!form.proxy_id) return null
  const proxy = proxyList.value.find(p => p.id === form.proxy_id)
  if (!proxy) return null
  return {
    enabled: true,
    type: proxy.type,
    host: proxy.host,
    port: proxy.port,
    username: proxy.username || '',
    password: proxy.password || ''
  }
})

// 加载代理列表
async function loadProxies() {
  loadingProxies.value = true
  try {
    const res = await api.getEnabledProxyConfigs()
    proxyList.value = res.items || []
  } catch (e) {
    console.error('Failed to load proxies:', e)
  } finally {
    loadingProxies.value = false
  }
}

// 加载全局模型映射列表
async function loadMappings() {
  loadingMappings.value = true
  try {
    const res = await api.getModelMappings()
    globalMappings.value = (res.data?.mappings || []).filter(m => m.enabled)
  } catch (e) {
    console.error('Failed to load model mappings:', e)
  } finally {
    loadingMappings.value = false
  }
}

const defaultForm = {
  name: '',
  type: '',
  description: '',
  enabled: true,
  priority: 50,
  weight: 100,
  max_concurrency: 5,
  daily_budget: 0,
  accountType: 'shared',
  addType: 'oauth',
  api_key: '',
  api_url: '',
  access_token: '',
  refresh_token: '',
  session_key: '',
  organization_id: '',
  aws_access_key_id: '',
  aws_secret_access_key: '',
  aws_region: 'us-east-1',
  azure_endpoint: '',
  azure_deployment_name: '',
  azure_api_version: '2024-02-01',
  allowed_models: '',
  allowedModelsList: [],
  model_mapping: '',
  selectedMappingIds: [],
  proxy_id: null
}

const form = reactive({ ...defaultForm })

// 是否显示添加方式选择
const showAddTypeSelection = computed(() => {
  const type = form.type
  return ['claude-official', 'openai-responses', 'gemini'].includes(type)
})

// 是否显示平台特定配置
const showPlatformConfig = computed(() => {
  const type = form.type
  if (['claude-console', 'bedrock', 'azure-openai'].includes(type)) return true
  if ((type === 'openai' || type === 'gemini') && form.addType === 'apikey') return true
  if (type === 'openai-responses' && form.addType === 'cookie') return true
  return false
})

// 获取平台配置标题
const getPlatformConfigTitle = computed(() => {
  const typeLabels = {
    'claude-console': 'Claude Console 配置',
    'bedrock': 'AWS Bedrock 配置',
    'azure-openai': 'Azure OpenAI 配置',
    'openai': 'OpenAI 三方 API 配置',
    'openai-responses': 'ChatGPT 官方配置',
    'gemini': 'Gemini 配置'
  }
  if (form.addType === 'manual') return 'Token 配置'
  return typeLabels[form.type] || '平台配置'
})

// 是否需要 OAuth
const needsOAuth = computed(() => {
  if (form.type === 'openai-responses' && form.addType === 'cookie') {
    return false
  }
  return form.addType === 'oauth' || form.addType === 'cookie'
})

// 编辑模式下是否显示 API 配置
const showEditApiConfig = computed(() => {
  const type = form.type
  return ['claude-console', 'openai', 'gemini'].includes(type)
})

// 是否可以继续
const canProceed = computed(() => {
  if (!form.type || !form.name) return false
  return true
})

// 监听编辑数据变化
watch(() => props.editData, (val) => {
  if (val) {
    Object.assign(form, { ...defaultForm, ...val })
    if (val.base_url) {
      form.api_url = val.base_url
    }
    if (val.allowed_models) {
      try {
        form.allowedModelsList = val.allowed_models.split(',').map(s => s.trim()).filter(Boolean)
      } catch {
        form.allowedModelsList = []
      }
    } else {
      form.allowedModelsList = []
    }
    if (val.model_mapping) {
      try {
        const mappingObj = JSON.parse(val.model_mapping)
        form._rawModelMapping = mappingObj
        form.selectedMappingIds = []
      } catch {
        form.selectedMappingIds = []
      }
    } else {
      form.selectedMappingIds = []
    }
    for (const [key, subs] of Object.entries(subplatformMap)) {
      if (subs.some(s => s.value === val.type)) {
        platformGroup.value = key
        break
      }
    }
  }
}, { immediate: true })

// 监听弹窗显示
watch(visible, async (val) => {
  if (val && !isEdit.value) {
    resetForm()
  }
  if (val) {
    loadProxies()
    await loadMappings()

    if (form._rawModelMapping && globalMappings.value.length > 0) {
      const matchedIds = []
      for (const [source, target] of Object.entries(form._rawModelMapping)) {
        const found = globalMappings.value.find(
          m => m.source_model === source && m.target_model === target
        )
        if (found) {
          matchedIds.push(found.id)
        }
      }
      form.selectedMappingIds = matchedIds
      delete form._rawModelMapping
    }
  }
})

function selectPlatformGroup(key) {
  platformGroup.value = key
  form.type = ''
  form.allowedModelsList = []
  if (key === 'claude') {
    form.addType = 'cookie'
  } else {
    form.addType = 'apikey'
  }
}

// 监听子平台类型变化
watch(() => form.type, (newType) => {
  if (!newType) return
  if (newType === 'claude-official') {
    form.addType = 'cookie'
  } else if (newType === 'openai-responses') {
    form.addType = 'oauth'
  } else if (['claude-console', 'bedrock', 'azure-openai'].includes(newType)) {
    form.addType = 'apikey'
  } else if (['openai', 'gemini'].includes(newType)) {
    form.addType = 'apikey'
  }
})

function resetForm() {
  Object.assign(form, { ...defaultForm })
  platformGroup.value = ''
  step.value = 1
  showApiKey.value = false
  showAwsSecret.value = false
  showSessionKey.value = false
  modelInputValue.value = ''
}

function handleClose() {
  visible.value = false
  resetForm()
}

function focusModelInput() {
  modelInputRef.value?.focus()
}

function addModel() {
  const value = modelInputValue.value.trim()
  if (value && !form.allowedModelsList.includes(value)) {
    form.allowedModelsList.push(value)
  }
  modelInputValue.value = ''
}

function removeModel(model) {
  const index = form.allowedModelsList.indexOf(model)
  if (index > -1) {
    form.allowedModelsList.splice(index, 1)
  }
}

function handleModelBackspace() {
  if (!modelInputValue.value && form.allowedModelsList.length > 0) {
    form.allowedModelsList.pop()
  }
}

function showMessage(msg, type = 'info') {
  // 简单的消息提示
  const div = document.createElement('div')
  div.className = `toast-message toast-${type}`
  div.textContent = msg
  document.body.appendChild(div)
  setTimeout(() => {
    div.classList.add('show')
  }, 10)
  setTimeout(() => {
    div.classList.remove('show')
    setTimeout(() => div.remove(), 300)
  }, 3000)
}

async function handleNext() {
  if (!form.name) {
    showMessage('请输入账户名称', 'warning')
    return
  }

  if (needsOAuth.value) {
    step.value = 2
  } else {
    await handleSubmit()
  }
}

async function handleOAuthSuccess(tokenInfo) {
  if (Array.isArray(tokenInfo)) {
    for (const token of tokenInfo) {
      form.access_token = token.access_token
      form.refresh_token = token.refresh_token
      form.session_key = token.session_key
      await createAccount()
    }
  } else {
    form.access_token = tokenInfo.access_token
    form.refresh_token = tokenInfo.refresh_token
    await createAccount()
  }
}

async function createAccount() {
  submitting.value = true
  try {
    const submitData = buildSubmitData()
    await api.createAccount(submitData)
    showMessage('创建成功', 'success')
    emit('success')
    handleClose()
  } catch (e) {
    showMessage(e.message || '创建失败', 'error')
  } finally {
    submitting.value = false
  }
}

async function handleSubmit() {
  if (!form.name) {
    showMessage('请输入账户名称', 'warning')
    return
  }

  submitting.value = true
  try {
    const submitData = buildSubmitData()

    if (isEdit.value) {
      await api.updateAccount(form.id, submitData)
      showMessage('更新成功', 'success')
    } else {
      await api.createAccount(submitData)
      showMessage('创建成功', 'success')
    }
    emit('success')
    handleClose()
  } catch (e) {
    showMessage(e.message || '操作失败', 'error')
  } finally {
    submitting.value = false
  }
}

function buildSubmitData() {
  const data = {
    name: form.name,
    type: form.type,
    description: form.description,
    enabled: form.enabled,
    priority: form.priority,
    weight: form.weight,
    max_concurrency: form.max_concurrency,
    account_type: form.accountType
  }

  if (form.api_key) data.api_key = form.api_key
  if (form.api_url) data.base_url = form.api_url
  if (form.access_token) data.access_token = form.access_token
  if (form.refresh_token) data.refresh_token = form.refresh_token
  if (form.session_key) data.session_key = form.session_key
  if (form.organization_id) data.organization_id = form.organization_id

  if (form.type === 'bedrock') {
    data.aws_access_key = form.aws_access_key_id
    data.aws_secret_key = form.aws_secret_access_key
    data.aws_region = form.aws_region
  }

  if (form.type === 'azure-openai') {
    data.azure_endpoint = form.azure_endpoint
    data.azure_deployment_name = form.azure_deployment_name
    data.azure_api_version = form.azure_api_version
  }

  if (form.proxy_id) {
    data.proxy_id = form.proxy_id
  } else if (isEdit.value && props.editData?.proxy_id) {
    data.clear_proxy = true
  }

  if (form.allowedModelsList?.length > 0) {
    data.allowed_models = form.allowedModelsList.join(',')
  } else if (isEdit.value && props.editData?.allowed_models) {
    data.clear_allowed_models = true
  }

  if (form.selectedMappingIds?.length > 0) {
    const mappingObj = {}
    form.selectedMappingIds.forEach(id => {
      const mapping = globalMappings.value.find(m => m.id === id)
      if (mapping) {
        mappingObj[mapping.source_model] = mapping.target_model
      }
    })
    if (Object.keys(mappingObj).length > 0) {
      data.model_mapping = JSON.stringify(mappingObj)
    }
  } else if (isEdit.value && props.editData?.model_mapping) {
    data.clear_model_mapping = true
  }

  return data
}

function getTypeLabel(type) {
  for (const subs of Object.values(subplatformMap)) {
    const found = subs.find(s => s.value === type)
    if (found) return found.label
  }
  return type
}

function getTypeIconComponent(type) {
  for (const subs of Object.values(subplatformMap)) {
    const found = subs.find(s => s.value === type)
    if (found) return found.iconComponent
  }
  return BrainIcon
}

function getTypeColor(type) {
  for (const subs of Object.values(subplatformMap)) {
    const found = subs.find(s => s.value === type)
    if (found) return found.color
  }
  return '#999'
}
</script>

<style scoped>
/* Modal 基础样式 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-container {
  background: var(--apple-bg-primary, #fff);
  border-radius: var(--apple-radius-xl, 20px);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  width: 100%;
  max-width: 800px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--apple-border, rgba(0, 0, 0, 0.1));
}

.modal-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: none;
  background: var(--apple-bg-tertiary, rgba(0, 0, 0, 0.05));
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.modal-close:hover {
  background: var(--apple-bg-secondary, rgba(0, 0, 0, 0.1));
}

.modal-close svg {
  width: 16px;
  height: 16px;
  color: var(--apple-text-secondary, #86868b);
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--apple-border, rgba(0, 0, 0, 0.1));
  background: var(--apple-bg-secondary, rgba(0, 0, 0, 0.02));
}

/* 步骤指示器 */
.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 16px 24px;
  background: var(--apple-bg-secondary, rgba(0, 0, 0, 0.02));
  border-bottom: 1px solid var(--apple-border, rgba(0, 0, 0, 0.1));
}

.step {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--apple-text-tertiary, #c7c7cc);
}

.step.active {
  color: var(--apple-accent, #007aff);
}

.step.done {
  color: var(--apple-success, #34c759);
}

.step-num {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--apple-bg-tertiary, rgba(0, 0, 0, 0.05));
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
}

.step.active .step-num {
  background: var(--apple-accent, #007aff);
  color: white;
}

.step.done .step-num {
  background: var(--apple-success, #34c759);
  color: white;
}

.step-num svg {
  width: 14px;
  height: 14px;
}

.step-label {
  font-size: 14px;
  font-weight: 500;
}

.step-line {
  width: 60px;
  height: 2px;
  background: var(--apple-border, rgba(0, 0, 0, 0.1));
}

.step-line.done {
  background: var(--apple-success, #34c759);
}

/* 表单区域 */
.form-step {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
  margin: 0;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--apple-border, rgba(0, 0, 0, 0.1));
}

.section-title svg {
  width: 18px;
  height: 18px;
  color: var(--apple-accent, #007aff);
}

/* 平台卡片 */
.platform-groups {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.platform-card {
  position: relative;
  padding: 16px;
  border: 2px solid var(--apple-border, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-md, 12px);
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--apple-bg-primary, #fff);
}

.platform-card:hover {
  border-color: var(--apple-accent, #007aff);
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.15);
}

.platform-card.selected {
  border-color: var(--apple-accent, #007aff);
  background: rgba(0, 122, 255, 0.05);
}

.platform-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--apple-radius-md, 12px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-bottom: 12px;
}

.platform-icon svg {
  width: 24px;
  height: 24px;
}

.platform-info h5 {
  margin: 0 0 4px;
  font-size: 15px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
}

.platform-info p {
  margin: 0;
  font-size: 12px;
  color: var(--apple-text-secondary, #86868b);
}

.check-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--apple-accent, #007aff);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
}

.check-badge svg {
  width: 12px;
  height: 12px;
}

.check-badge.small {
  width: 20px;
  height: 20px;
}

.check-badge.small svg {
  width: 10px;
  height: 10px;
}

/* 子平台 */
.subplatform-section {
  margin-top: 16px;
  padding: 16px;
  background: var(--apple-bg-secondary, rgba(0, 0, 0, 0.02));
  border-radius: var(--apple-radius-md, 12px);
  border: 1px solid var(--apple-border, rgba(0, 0, 0, 0.1));
}

.subplatform-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--apple-text-secondary, #86868b);
  margin-bottom: 12px;
}

.subplatform-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.subplatform-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 2px solid var(--apple-border, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-sm, 8px);
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--apple-bg-primary, #fff);
}

.subplatform-card:hover {
  border-color: var(--apple-accent, #007aff);
}

.subplatform-card.selected {
  border-color: var(--apple-accent, #007aff);
  background: rgba(0, 122, 255, 0.05);
}

.sub-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--apple-radius-sm, 8px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.sub-icon svg {
  width: 18px;
  height: 18px;
}

.sub-info {
  display: flex;
  flex-direction: column;
}

.sub-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
}

.sub-desc {
  font-size: 11px;
  color: var(--apple-text-secondary, #86868b);
}

.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

/* 添加方式 */
.add-type-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.type-card {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 16px;
  border: 2px solid var(--apple-border, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-md, 12px);
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--apple-bg-primary, #fff);
}

.type-card:hover {
  border-color: var(--apple-accent, #007aff);
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.1);
}

.type-card.selected {
  border-color: var(--apple-accent, #007aff);
  background: rgba(0, 122, 255, 0.05);
}

.type-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--apple-radius-sm, 8px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.type-icon.cookie {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.type-icon.oauth {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.type-icon.apikey {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.type-icon svg {
  width: 20px;
  height: 20px;
}

.type-content {
  flex: 1;
}

.type-content h5 {
  margin: 0 0 4px;
  font-size: 14px;
  font-weight: 600;
  color: var(--apple-text-primary, #1d1d1f);
}

.type-content p {
  margin: 0;
  font-size: 12px;
  color: var(--apple-text-secondary, #86868b);
  line-height: 1.4;
}

/* 表单网格 */
.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group.full {
  grid-column: 1 / -1;
}

.form-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 500;
  color: var(--apple-text-primary, #1d1d1f);
}

.form-label.required::after {
  content: '*';
  color: var(--apple-danger, #ff3b30);
}

.label-icon {
  width: 14px;
  height: 14px;
  color: var(--apple-accent, #007aff);
}

.form-input,
.form-select,
.form-textarea {
  padding: 10px 14px;
  font-size: 14px;
  border: 1px solid var(--apple-border, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-sm, 8px);
  background: var(--apple-bg-primary, #fff);
  color: var(--apple-text-primary, #1d1d1f);
  transition: all 0.2s ease;
  font-family: inherit;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--apple-accent, #007aff);
  box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 60px;
}

.form-select.multi {
  min-height: 100px;
}

.input-with-toggle {
  position: relative;
  display: flex;
}

.input-with-toggle .form-input {
  flex: 1;
  padding-right: 44px;
}

.toggle-visibility {
  position: absolute;
  right: 4px;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  border: none;
  background: transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--apple-radius-sm, 8px);
}

.toggle-visibility:hover {
  background: var(--apple-bg-tertiary, rgba(0, 0, 0, 0.05));
}

.toggle-visibility svg {
  width: 18px;
  height: 18px;
  color: var(--apple-text-secondary, #86868b);
}

.form-tip {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--apple-text-secondary, #86868b);
  margin-top: 4px;
}

.form-tip svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
}

.form-tip a {
  color: var(--apple-accent, #007aff);
  text-decoration: none;
}

.form-tip a:hover {
  text-decoration: underline;
}

/* Toggle Switch */
.toggle-switch {
  position: relative;
  display: inline-flex;
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
  background-color: var(--apple-bg-tertiary, rgba(0, 0, 0, 0.1));
  border-radius: 31px;
  transition: 0.3s;
}

.toggle-slider::before {
  position: absolute;
  content: '';
  height: 27px;
  width: 27px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  border-radius: 50%;
  transition: 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.toggle-switch input:checked + .toggle-slider {
  background-color: var(--apple-success, #34c759);
}

.toggle-switch input:checked + .toggle-slider::before {
  transform: translateX(20px);
}

/* Tags Input */
.tags-input {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px 12px;
  border: 1px solid var(--apple-border, rgba(0, 0, 0, 0.1));
  border-radius: var(--apple-radius-sm, 8px);
  background: var(--apple-bg-primary, #fff);
  min-height: 44px;
  cursor: text;
  transition: all 0.2s ease;
}

.tags-input:focus-within {
  border-color: var(--apple-accent, #007aff);
  box-shadow: 0 0 0 3px rgba(0, 122, 255, 0.1);
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: var(--apple-accent, #007aff);
  color: white;
  border-radius: 16px;
  font-size: 13px;
}

.tag-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border: none;
  background: rgba(255, 255, 255, 0.3);
  color: white;
  border-radius: 50%;
  cursor: pointer;
  font-size: 12px;
  line-height: 1;
}

.tag-remove:hover {
  background: rgba(255, 255, 255, 0.5);
}

.tag-input {
  flex: 1;
  min-width: 120px;
  border: none;
  outline: none;
  font-size: 14px;
  background: transparent;
  color: var(--apple-text-primary, #1d1d1f);
}

.tag-input::placeholder {
  color: var(--apple-text-tertiary, #c7c7cc);
}

/* Help Card */
.help-card {
  margin-top: 16px;
  padding: 16px;
  background: rgba(0, 122, 255, 0.05);
  border: 1px solid rgba(0, 122, 255, 0.2);
  border-radius: var(--apple-radius-md, 12px);
}

.help-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--apple-accent, #007aff);
  margin-bottom: 12px;
}

.help-header svg {
  width: 18px;
  height: 18px;
}

.help-steps {
  margin: 0;
  padding-left: 20px;
  font-size: 13px;
  line-height: 1.8;
  color: var(--apple-text-primary, #1d1d1f);
}

.help-steps li {
  margin-bottom: 4px;
}

.help-steps a {
  color: var(--apple-accent, #007aff);
  text-decoration: none;
}

.help-steps a:hover {
  text-decoration: underline;
}

.help-steps kbd {
  background: var(--apple-bg-tertiary, rgba(0, 0, 0, 0.05));
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
}

.help-steps code {
  background: rgba(0, 122, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
  color: var(--apple-accent, #007aff);
}

.help-note {
  margin-top: 12px;
  font-size: 12px;
  color: var(--apple-text-secondary, #86868b);
}

.help-note code {
  background: var(--apple-bg-tertiary, rgba(0, 0, 0, 0.05));
  padding: 2px 4px;
  border-radius: 3px;
  font-family: monospace;
}

/* Edit Banner */
.edit-type-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--apple-bg-secondary, rgba(0, 0, 0, 0.02));
  border-radius: var(--apple-radius-md, 12px);
  margin-bottom: 20px;
}

.type-badge {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-sm, 8px);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.type-badge svg {
  width: 16px;
  height: 16px;
}

.edit-type-banner > span {
  font-weight: 500;
  color: var(--apple-text-primary, #1d1d1f);
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 20px;
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--apple-radius-sm, 8px);
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn svg {
  width: 16px;
  height: 16px;
}

.btn-primary {
  background: var(--apple-accent, #007aff);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #0066d6;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--apple-bg-tertiary, rgba(0, 0, 0, 0.05));
  color: var(--apple-text-primary, #1d1d1f);
}

.btn-secondary:hover {
  background: var(--apple-bg-secondary, rgba(0, 0, 0, 0.1));
}

.btn-loading {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Modal Transition */
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.95) translateY(20px);
}

/* Toast Messages */
:global(.toast-message) {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%) translateY(-100px);
  padding: 12px 24px;
  background: var(--apple-bg-primary, #fff);
  border-radius: var(--apple-radius-md, 12px);
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  font-size: 14px;
  font-weight: 500;
  z-index: 9999;
  transition: transform 0.3s ease;
}

:global(.toast-message.show) {
  transform: translateX(-50%) translateY(0);
}

:global(.toast-success) {
  color: var(--apple-success, #34c759);
}

:global(.toast-error) {
  color: var(--apple-danger, #ff3b30);
}

:global(.toast-warning) {
  color: var(--apple-warning, #ff9500);
}

/* Dark Mode Support */
@media (prefers-color-scheme: dark) {
  .modal-container {
    background: #1c1c1e;
  }

  .modal-title,
  .section-title,
  .platform-info h5,
  .sub-name,
  .type-content h5,
  .form-label {
    color: #f5f5f7;
  }

  .modal-close {
    background: rgba(255, 255, 255, 0.1);
  }

  .modal-close:hover {
    background: rgba(255, 255, 255, 0.15);
  }

  .platform-card,
  .subplatform-card,
  .type-card {
    background: #2c2c2e;
    border-color: rgba(255, 255, 255, 0.1);
  }

  .form-input,
  .form-select,
  .form-textarea {
    background: #2c2c2e;
    border-color: rgba(255, 255, 255, 0.1);
    color: #f5f5f7;
  }

  .tags-input {
    background: #2c2c2e;
    border-color: rgba(255, 255, 255, 0.1);
  }

  .tag-input {
    color: #f5f5f7;
  }
}
</style>
