<!--
 * 文件作用：系统设置页面 - Apple HIG 风格
 * 负责功能：
 *   - 安全配置（验证码、登录限制）
 *   - 记录配置（保留天数、价格倍率）
 *   - 账号健康检查配置
 *   - 分级检测策略配置
 * 重要程度：⭐⭐⭐⭐ 重要（系统配置）
-->
<template>
  <div class="settings-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">系统设置</h1>
        <p class="page-subtitle">配置系统参数、安全策略和账号健康检查</p>
      </div>
    </div>

    <!-- 设置网格 -->
    <div class="settings-grid">
      <!-- 安全配置 -->
      <div class="settings-card">
        <div class="card-header">
          <div class="card-icon security">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
            </svg>
          </div>
          <h3>安全配置</h3>
        </div>
        <div class="card-body" :class="{ loading }">
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">启用登录验证码</span>
              <span class="setting-desc">开启后登录需要输入图片验证码</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="captchaEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">验证码频率限制</span>
              <span class="setting-desc">每个 IP 每分钟最多获取验证码次数</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.captcha_rate_limit" min="1" max="100" class="form-input" @change="markDirty" />
              <span class="unit">次/分钟</span>
            </div>
          </div>

          <div class="setting-divider"></div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">启用登录频率限制</span>
              <span class="setting-desc">防止暴力破解，限制登录尝试次数</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="loginRateLimitEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">登录限制次数</span>
              <span class="setting-desc">时间窗口内允许的最大登录尝试次数</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.login_rate_limit_count" min="1" max="100" class="form-input" :disabled="!loginRateLimitEnabled" @change="markDirty" />
              <span class="unit">次</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">登录限制窗口</span>
              <span class="setting-desc">频率限制的时间窗口</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.login_rate_limit_window" min="1" max="60" class="form-input" :disabled="!loginRateLimitEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 记录配置 -->
      <div class="settings-card">
        <div class="card-header">
          <div class="card-icon records">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
              <polyline points="14,2 14,8 20,8"/>
              <line x1="16" y1="13" x2="8" y2="13"/>
              <line x1="16" y1="17" x2="8" y2="17"/>
            </svg>
          </div>
          <h3>记录配置</h3>
        </div>
        <div class="card-body" :class="{ loading }">
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">记录保留天数</span>
              <span class="setting-desc">使用记录的保留时间</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.record_retention_days" min="1" max="365" class="form-input" @change="markDirty" />
              <span class="unit">天</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">最大记录数</span>
              <span class="setting-desc">每个用户保留的最大记录数</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.record_max_count" min="100" max="10000" step="100" class="form-input" @change="markDirty" />
              <span class="unit">条/用户</span>
            </div>
          </div>

          <div class="setting-divider"></div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">全局价格倍率</span>
              <span class="setting-desc">1=原价，0=免费，2=2倍。优先级：全局 → 用户</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.global_price_rate" min="0" max="10" step="0.1" class="form-input" @change="markDirty" />
              <span class="unit">倍</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 账号健康检查 -->
      <div class="settings-card">
        <div class="card-header">
          <div class="card-icon health">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
            </svg>
          </div>
          <div class="card-title-row">
            <h3>账号健康检查</h3>
            <div class="header-badges">
              <span :class="['status-badge', healthCheckEnabled ? 'success' : 'muted']">
                {{ healthCheckEnabled ? '已启用' : '已禁用' }}
              </span>
              <button class="btn btn-sm btn-outline" @click="viewHealthCheckStatus" :disabled="!healthCheckEnabled">
                查看状态
              </button>
            </div>
          </div>
        </div>
        <div class="card-body" :class="{ loading }">
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">启用健康检查</span>
              <span class="setting-desc">定期检查 OAuth 账号的有效性</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="healthCheckEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">检查间隔</span>
              <span class="setting-desc">正常账号的定期检查间隔</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.account_health_check_interval" min="1" max="60" class="form-input" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">连续错误阈值</span>
              <span class="setting-desc">连续检查失败达到此次数后标记为疑似封号</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.account_error_threshold" min="1" max="100" class="form-input" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="unit">次</span>
            </div>
          </div>

          <div class="setting-section-title">恢复策略</div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">自动恢复</span>
              <span class="setting-desc">检测通过时自动恢复账号为正常状态</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="healthCheckAutoRecovery" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">自动刷新 Token</span>
              <span class="setting-desc">Token 过期时自动尝试刷新</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="healthCheckAutoTokenRefresh" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-section-title">OAuth 配置</div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">OAuth 自动重授权</span>
              <span class="setting-desc">OAuth Token 失效时，自动使用 SessionKey 重新获取 Token</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="oauthAutoReauthorizeEnabled" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">重授权冷却时间</span>
              <span class="setting-desc">重新授权失败后的冷却时间</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.oauth_reauthorize_cooldown" min="5" max="1440" class="form-input" :disabled="!healthCheckEnabled || !oauthAutoReauthorizeEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 分级检测策略 -->
      <div class="settings-card">
        <div class="card-header">
          <div class="card-icon strategy">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="4" y1="21" x2="4" y2="14"/>
              <line x1="4" y1="10" x2="4" y2="3"/>
              <line x1="12" y1="21" x2="12" y2="12"/>
              <line x1="12" y1="8" x2="12" y2="3"/>
              <line x1="20" y1="21" x2="20" y2="16"/>
              <line x1="20" y1="12" x2="20" y2="3"/>
              <line x1="1" y1="14" x2="7" y2="14"/>
              <line x1="9" y1="8" x2="15" y2="8"/>
              <line x1="17" y1="16" x2="23" y2="16"/>
            </svg>
          </div>
          <h3>分级检测策略</h3>
        </div>
        <div class="card-body" :class="{ loading }">
          <div class="setting-section-title">限流账号探测</div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">启用主动探测</span>
              <span class="setting-desc">限流账号不等待官方 reset 时间，主动探测恢复</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="rateLimitedProbeEnabled" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">初始探测间隔</span>
              <span class="setting-desc">限流后首次探测的等待时间</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.rate_limited_probe_init_interval" min="1" max="60" class="form-input" :disabled="!healthCheckEnabled || !rateLimitedProbeEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">最大探测间隔</span>
              <span class="setting-desc">探测间隔递增的上限</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.rate_limited_probe_max_interval" min="10" max="120" class="form-input" :disabled="!healthCheckEnabled || !rateLimitedProbeEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">间隔递增因子</span>
              <span class="setting-desc">每次失败后间隔乘以此因子</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.rate_limited_probe_backoff" min="1" max="3" step="0.1" class="form-input" :disabled="!healthCheckEnabled || !rateLimitedProbeEnabled" @change="markDirty" />
              <span class="unit">倍</span>
            </div>
          </div>

          <div class="setting-section-title">疑似封号检测</div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">探测间隔</span>
              <span class="setting-desc">疑似封号账号的探测间隔</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.suspended_probe_interval" min="1" max="60" class="form-input" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">确认封号阈值</span>
              <span class="setting-desc">连续失败此次数后确认为封号状态</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.suspended_confirm_threshold" min="1" max="10" class="form-input" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="unit">次</span>
            </div>
          </div>

          <div class="setting-section-title">封号账号复活检测</div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">启用复活检测</span>
              <span class="setting-desc">定期检测已封号账号是否恢复</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="bannedProbeEnabled" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">复活探测间隔</span>
              <span class="setting-desc">封号账号的复活检测间隔</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.banned_probe_interval" min="1" max="24" class="form-input" :disabled="!healthCheckEnabled || !bannedProbeEnabled" @change="markDirty" />
              <span class="unit">小时</span>
            </div>
          </div>

          <div class="setting-section-title">Token 刷新</div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">刷新冷却时间</span>
              <span class="setting-desc">Token 刷新失败后的冷却时间</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.token_refresh_cooldown" min="5" max="120" class="form-input" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">最大重试次数</span>
              <span class="setting-desc">Token 刷新的最大重试次数</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.token_refresh_max_retries" min="1" max="10" class="form-input" :disabled="!healthCheckEnabled" @change="markDirty" />
              <span class="unit">次</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 用量同步配置 -->
      <div class="settings-card">
        <div class="card-header">
          <div class="card-icon sync">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12a9 9 0 11-6.219-8.56"/>
              <polyline points="21,4 21,12 13,12"/>
            </svg>
          </div>
          <div class="card-title-row">
            <h3>用量同步配置</h3>
            <div class="header-badges">
              <span :class="['status-badge', usageSyncEnabled ? 'success' : 'muted']">
                {{ usageSyncEnabled ? '已启用' : '已禁用' }}
              </span>
            </div>
          </div>
        </div>
        <div class="card-body" :class="{ loading }">
          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">启用自动同步</span>
              <span class="setting-desc">定期自动同步账户用量信息（Claude/OpenAI/Gemini）</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="usageSyncEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">同步间隔</span>
              <span class="setting-desc">自动同步用量的时间间隔</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.usage_sync_interval" min="5" max="1440" class="form-input" :disabled="!usageSyncEnabled" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>

          <div class="setting-divider"></div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">缓存有效期</span>
              <span class="setting-desc">避免频繁查询同一账户的缓存时间</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.usage_sync_cache_ttl" min="1" max="60" class="form-input" @change="markDirty" />
              <span class="unit">分钟</span>
            </div>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">并发查询数</span>
              <span class="setting-desc">批量同步时的最大并发请求数</span>
            </div>
            <div class="input-with-unit">
              <input type="number" v-model.number="configs.usage_sync_concurrency" min="1" max="20" class="form-input" @change="markDirty" />
              <span class="unit">个</span>
            </div>
          </div>

          <div class="setting-section-title">同步范围</div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">同步 Claude 账户</span>
              <span class="setting-desc">同步 Claude Official OAuth 账户的 5H/7D 用量</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="syncClaudeEnabled" :disabled="!usageSyncEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">同步 OpenAI 账户</span>
              <span class="setting-desc">同步 OpenAI Codex 用量（从响应头提取）</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="syncOpenAIEnabled" :disabled="!usageSyncEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div class="setting-item">
            <div class="setting-info">
              <span class="setting-label">同步 Gemini 账户</span>
              <span class="setting-desc">同步 Gemini Code Assist 项目 ID</span>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="syncGeminiEnabled" :disabled="!usageSyncEnabled" @change="markDirty" />
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
      </div>

    </div>

    <!-- 配置项列表 -->
    <div class="config-list-card">
      <div class="card-header">
        <h3>配置项列表</h3>
      </div>
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>配置项</th>
              <th>当前值</th>
              <th>说明</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="cfg in configList" :key="cfg.key">
              <td><code class="config-key">{{ cfg.key }}</code></td>
              <td><span class="config-value">{{ cfg.value }}</span></td>
              <td class="config-desc">{{ cfg.desc }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 保存按钮 -->
    <div class="action-bar">
      <button class="btn btn-secondary" @click="loadConfigs">重置</button>
      <button class="btn btn-primary" @click="saveConfigs" :disabled="saving">
        <span v-if="saving" class="btn-loading"></span>
        {{ saving ? '保存中...' : '保存配置' }}
      </button>
    </div>

    <!-- 健康检测状态弹窗 -->
    <Teleport to="body">
      <div v-if="healthStatusVisible" class="modal-overlay" @click.self="healthStatusVisible = false">
        <div class="modal">
          <div class="modal-header">
            <h2>健康检测服务状态</h2>
            <button class="modal-close" @click="healthStatusVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div v-if="healthStatus" class="health-status-grid">
              <div class="status-item">
                <span class="status-label">服务状态</span>
                <span :class="['status-badge', healthStatus.running ? 'success' : 'danger']">
                  {{ healthStatus.running ? '运行中' : '已停止' }}
                </span>
              </div>
              <div class="status-item">
                <span class="status-label">检测间隔</span>
                <span class="status-value">{{ healthStatus.interval }} 分钟</span>
              </div>
              <div class="status-item">
                <span class="status-label">上次检测</span>
                <span class="status-value">{{ healthStatus.last_check ? formatDate(healthStatus.last_check) : '暂无' }}</span>
              </div>
              <div class="status-item">
                <span class="status-label">问题账号数</span>
                <span class="status-badge warning">{{ healthStatus.problem_account_count || 0 }}</span>
              </div>
              <div class="status-item">
                <span class="status-label">已检测账号数</span>
                <span class="status-value">{{ healthStatus.checked_count || 0 }}</span>
              </div>
              <div class="status-item">
                <span class="status-label">失败账号数</span>
                <span :class="['status-badge', healthStatus.failed_count > 0 ? 'danger' : 'muted']">
                  {{ healthStatus.failed_count || 0 }}
                </span>
              </div>
              <div class="status-item">
                <span class="status-label">错误阈值</span>
                <span class="status-value">{{ healthStatus.threshold }} 次</span>
              </div>
              <div class="status-item full-width">
                <span class="status-label">最后错误</span>
                <span class="status-error">{{ healthStatus.last_error || '无' }}</span>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="refreshHealthStatus">刷新状态</button>
            <button class="btn btn-primary" @click="triggerHealthCheck" :disabled="healthChecking">
              {{ healthChecking ? '触发中...' : '手动触发检测' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const saving = ref(false)
const isDirty = ref(false)
const healthStatusVisible = ref(false)
const healthChecking = ref(false)
const healthStatus = ref(null)

const configs = reactive({
  record_retention_days: 30,
  record_max_count: 1000,
  global_price_rate: 1,
  captcha_enabled: 'true',
  captcha_rate_limit: 10,
  login_rate_limit_enable: 'true',
  login_rate_limit_count: 3,
  login_rate_limit_window: 5,
  account_health_check_enabled: 'false',
  account_health_check_interval: 5,
  account_error_threshold: 5,
  oauth_auto_reauthorize_enabled: 'true',
  oauth_reauthorize_cooldown: 30,
  health_check_auto_recovery: 'true',
  health_check_auto_token_refresh: 'true',
  rate_limited_probe_enabled: 'true',
  rate_limited_probe_init_interval: 10,
  rate_limited_probe_max_interval: 30,
  rate_limited_probe_backoff: 1.5,
  suspended_probe_interval: 5,
  suspended_confirm_threshold: 3,
  banned_probe_enabled: 'false',
  banned_probe_interval: 1,
  token_refresh_cooldown: 30,
  token_refresh_max_retries: 3,
  // 用量同步配置
  usage_sync_enabled: 'false',
  usage_sync_interval: 60,
  usage_sync_cache_ttl: 5,
  usage_sync_concurrency: 5,
  sync_claude_enabled: 'true',
  sync_openai_enabled: 'true',
  sync_gemini_enabled: 'true',
  // 邮件配置已移除
})

const configList = ref([])

function markDirty() {
  isDirty.value = true
}

const captchaEnabled = computed({
  get: () => configs.captcha_enabled === 'true',
  set: (val) => { configs.captcha_enabled = val ? 'true' : 'false' }
})

const loginRateLimitEnabled = computed({
  get: () => configs.login_rate_limit_enable === 'true',
  set: (val) => { configs.login_rate_limit_enable = val ? 'true' : 'false' }
})

const healthCheckEnabled = computed({
  get: () => configs.account_health_check_enabled === 'true',
  set: (val) => { configs.account_health_check_enabled = val ? 'true' : 'false' }
})

const oauthAutoReauthorizeEnabled = computed({
  get: () => configs.oauth_auto_reauthorize_enabled === 'true',
  set: (val) => { configs.oauth_auto_reauthorize_enabled = val ? 'true' : 'false' }
})

const healthCheckAutoRecovery = computed({
  get: () => configs.health_check_auto_recovery === 'true',
  set: (val) => { configs.health_check_auto_recovery = val ? 'true' : 'false' }
})

const healthCheckAutoTokenRefresh = computed({
  get: () => configs.health_check_auto_token_refresh === 'true',
  set: (val) => { configs.health_check_auto_token_refresh = val ? 'true' : 'false' }
})

const rateLimitedProbeEnabled = computed({
  get: () => configs.rate_limited_probe_enabled === 'true',
  set: (val) => { configs.rate_limited_probe_enabled = val ? 'true' : 'false' }
})

const bannedProbeEnabled = computed({
  get: () => configs.banned_probe_enabled === 'true',
  set: (val) => { configs.banned_probe_enabled = val ? 'true' : 'false' }
})

const usageSyncEnabled = computed({
  get: () => configs.usage_sync_enabled === 'true',
  set: (val) => { configs.usage_sync_enabled = val ? 'true' : 'false' }
})

const syncClaudeEnabled = computed({
  get: () => configs.sync_claude_enabled === 'true',
  set: (val) => { configs.sync_claude_enabled = val ? 'true' : 'false' }
})

const syncOpenAIEnabled = computed({
  get: () => configs.sync_openai_enabled === 'true',
  set: (val) => { configs.sync_openai_enabled = val ? 'true' : 'false' }
})

const syncGeminiEnabled = computed({
  get: () => configs.sync_gemini_enabled === 'true',
  set: (val) => { configs.sync_gemini_enabled = val ? 'true' : 'false' }
})

function formatDate(str) {
  if (!str) return ''
  return new Date(str).toLocaleString('zh-CN')
}

async function loadConfigs() {
  loading.value = true
  try {
    const res = await api.getSystemConfigs()
    const items = res.items || []
    configList.value = items.filter(cfg => cfg.category !== 'email')

    for (const cfg of configList.value) {
      if (cfg.key in configs) {
        if (cfg.type === 'int') {
          configs[cfg.key] = parseInt(cfg.value) || 0
        } else if (cfg.type === 'float') {
          configs[cfg.key] = parseFloat(cfg.value) || 0
        } else {
          configs[cfg.key] = cfg.value
        }
      }
    }
    isDirty.value = false
  } catch (e) {
    console.error('Failed to load configs:', e)
  } finally {
    loading.value = false
  }
}

async function saveConfigs() {
  saving.value = true
  try {
    const toSave = {
      record_retention_days: String(configs.record_retention_days),
      record_max_count: String(configs.record_max_count),
      global_price_rate: String(configs.global_price_rate),
      captcha_enabled: configs.captcha_enabled,
      captcha_rate_limit: String(configs.captcha_rate_limit),
      login_rate_limit_enable: configs.login_rate_limit_enable,
      login_rate_limit_count: String(configs.login_rate_limit_count),
      login_rate_limit_window: String(configs.login_rate_limit_window),
      account_health_check_enabled: configs.account_health_check_enabled,
      account_health_check_interval: String(configs.account_health_check_interval),
      account_error_threshold: String(configs.account_error_threshold),
      oauth_auto_reauthorize_enabled: configs.oauth_auto_reauthorize_enabled,
      oauth_reauthorize_cooldown: String(configs.oauth_reauthorize_cooldown),
      health_check_auto_recovery: configs.health_check_auto_recovery,
      health_check_auto_token_refresh: configs.health_check_auto_token_refresh,
      rate_limited_probe_enabled: configs.rate_limited_probe_enabled,
      rate_limited_probe_init_interval: String(configs.rate_limited_probe_init_interval),
      rate_limited_probe_max_interval: String(configs.rate_limited_probe_max_interval),
      rate_limited_probe_backoff: String(configs.rate_limited_probe_backoff),
      suspended_probe_interval: String(configs.suspended_probe_interval),
      suspended_confirm_threshold: String(configs.suspended_confirm_threshold),
      banned_probe_enabled: configs.banned_probe_enabled,
      banned_probe_interval: String(configs.banned_probe_interval),
      token_refresh_cooldown: String(configs.token_refresh_cooldown),
      token_refresh_max_retries: String(configs.token_refresh_max_retries)
    }
    await api.updateSystemConfigs(toSave)
    ElMessage.success('配置保存成功')
    isDirty.value = false
    await loadConfigs()
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

async function viewHealthCheckStatus() {
  healthStatusVisible.value = true
  await refreshHealthStatus()
}

async function refreshHealthStatus() {
  try {
    const res = await api.getHealthCheckStatus()
    healthStatus.value = res
  } catch (e) {
    console.error('Failed to get health status:', e)
  }
}

let healthStatusTimer = null
async function triggerHealthCheck() {
  healthChecking.value = true
  try {
    await api.triggerHealthCheck()
    ElMessage.success('健康检测已触发')
    if (healthStatusTimer) clearTimeout(healthStatusTimer)
    healthStatusTimer = setTimeout(refreshHealthStatus, 2000)
  } catch (e) {
    ElMessage.error('触发失败')
  } finally {
    healthChecking.value = false
  }
}

onUnmounted(() => {
  if (healthStatusTimer) {
    clearTimeout(healthStatusTimer)
    healthStatusTimer = null
  }
})

onMounted(() => {
  loadConfigs()
})
</script>

<style scoped>
.settings-page {
  max-width: 1200px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  margin-bottom: var(--apple-spacing-xl);
}

.page-title {
  font-size: var(--apple-text-3xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-xs) 0;
}

.page-subtitle {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  margin: 0;
}

/* 设置网格 */
.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-xl);
}

/* 设置卡片 */
.settings-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.card-header h3 {
  margin: 0;
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.card-title-row {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-badges {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.card-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.card-icon svg {
  width: 20px;
  height: 20px;
  color: white;
}

.card-icon.security { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.card-icon.records { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.card-icon.health { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.card-icon.strategy { background: linear-gradient(135deg, #fc4a1a 0%, #f7b733 100%); }
.card-icon.sync { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }

.card-body {
  padding: var(--apple-spacing-lg);
  transition: opacity var(--apple-duration-fast);
}

.card-body.loading {
  opacity: 0.5;
  pointer-events: none;
}

/* 设置项 */
.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-sm) 0;
}

.setting-info {
  flex: 1;
  margin-right: var(--apple-spacing-md);
}

.setting-label {
  display: block;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.setting-desc {
  display: block;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: 2px;
}

.setting-divider {
  height: 1px;
  background: var(--apple-separator);
  margin: var(--apple-spacing-md) 0;
}

.setting-section-title {
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: var(--apple-spacing-lg) 0 var(--apple-spacing-sm) 0;
  padding-top: var(--apple-spacing-md);
  border-top: 1px solid var(--apple-separator);
}

/* 第一个 section title 不需要上边框 */
.card-body > .setting-section-title:first-child {
  margin-top: 0;
  padding-top: 0;
  border-top: none;
}

/* 全宽设置项 */
.setting-item.full-width {
  flex-direction: column;
  align-items: flex-start;
  gap: var(--apple-spacing-xs);
}

/* 输入框带单位 */
.input-with-unit {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.input-with-unit .form-input {
  width: 80px;
  text-align: center;
}

.unit {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
  white-space: nowrap;
}

/* 开关 */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
  flex-shrink: 0;
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
  background-color: var(--apple-fill-tertiary);
  transition: 0.3s;
  border-radius: 24px;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
  box-shadow: 0 1px 3px rgba(0,0,0,0.2);
}

.toggle-switch input:checked + .toggle-slider {
  background-color: var(--apple-green);
}

.toggle-switch input:checked + .toggle-slider:before {
  transform: translateX(20px);
}

.toggle-switch input:disabled + .toggle-slider {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 配置列表卡片 */
.config-list-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
  margin-bottom: var(--apple-spacing-xl);
}

.config-list-card .card-header {
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
}

.config-list-card .card-header h3 {
  font-size: var(--apple-text-sm);
}

.table-container {
  max-height: 300px;
  overflow-y: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--apple-text-sm);
}

.data-table th,
.data-table td {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  text-align: left;
  border-bottom: 1px solid var(--apple-separator);
}

.data-table th {
  background: var(--apple-bg-secondary);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  position: sticky;
  top: 0;
}

.config-key {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-blue);
  background: var(--apple-blue-light);
  padding: 2px 6px;
  border-radius: var(--apple-radius-xs);
}

.config-value {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-primary);
}

.config-desc {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  max-width: 300px;
}

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-xl);
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-lg);
}

/* 状态徽章 */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.status-badge.success { background: var(--apple-green-light); color: var(--apple-green); }
.status-badge.warning { background: var(--apple-orange-light); color: var(--apple-orange); }
.status-badge.danger { background: var(--apple-red-light); color: var(--apple-red); }
.status-badge.muted { background: var(--apple-fill-tertiary); color: var(--apple-text-tertiary); }

/* 按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  cursor: pointer;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-sm {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
}

.btn-primary {
  background: var(--apple-blue);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--apple-blue-hover);
}

.btn-secondary {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--apple-fill-secondary);
}

.btn-outline {
  background: transparent;
  color: var(--apple-blue);
  border: 1px solid var(--apple-blue);
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-blue-light);
}

.btn-loading {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 表单输入 */
.form-input {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-input:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.form-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 下拉选择框 */
.form-select {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  cursor: pointer;
  min-width: 150px;
}

.form-select:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--apple-z-modal);
  padding: var(--apple-spacing-xl);
}

.modal {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  box-shadow: var(--apple-shadow-modal);
  width: 100%;
  max-width: 560px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: var(--apple-spacing-xl);
  border-bottom: 1px solid var(--apple-separator);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header h2 {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.modal-close svg {
  width: 18px;
  height: 18px;
}

.modal-close:hover {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
}

.modal-body {
  padding: var(--apple-spacing-xl);
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: var(--apple-spacing-lg) var(--apple-spacing-xl);
  border-top: 1px solid var(--apple-separator);
  display: flex;
  justify-content: flex-end;
  gap: var(--apple-spacing-sm);
}

/* 健康状态网格 */
.health-status-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-md);
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xxs);
  padding: var(--apple-spacing-sm);
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-sm);
}

.status-item.full-width {
  grid-column: span 2;
}

.status-label {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.status-value {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.status-error {
  font-size: var(--apple-text-sm);
  color: var(--apple-red);
}

/* 响应式 - 移动端适配 */
@media (max-width: 1024px) {
  .settings-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 767px) {
  .settings-page {
    padding: 0 var(--apple-spacing-md);
  }

  .page-title {
    font-size: var(--apple-text-2xl);
  }

  .page-subtitle {
    font-size: var(--apple-text-sm);
  }

  /* 设置卡片 */
  .settings-card {
    margin: 0 calc(var(--apple-spacing-md) * -1);
    border-radius: 0;
  }

  .card-header {
    padding: var(--apple-spacing-md);
    flex-wrap: wrap;
    gap: var(--apple-spacing-sm);
  }

  .card-header h3 {
    font-size: var(--apple-text-base);
  }

  .card-title-row {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--apple-spacing-xs);
  }

  .header-badges {
    flex-wrap: wrap;
  }

  .card-body {
    padding: var(--apple-spacing-md);
  }

  .card-icon {
    width: 36px;
    height: 36px;
  }

  .card-icon svg {
    width: 18px;
    height: 18px;
  }

  /* 设置项 */
  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--apple-spacing-sm);
  }

  .setting-info {
    margin-right: 0;
    width: 100%;
  }

  .toggle-switch {
    width: 50px;
    height: 28px;
  }

  .toggle-slider:before {
    width: 24px;
    height: 24px;
  }

  .toggle-switch input:checked + .toggle-slider:before {
    transform: translateX(22px);
  }

  .input-with-unit {
    width: 100%;
    justify-content: flex-start;
  }

  .input-with-unit .form-input {
    width: 100px;
    min-height: var(--apple-touch-target);
    font-size: 16px;
  }

  /* 配置列表卡片 */
  .config-list-card {
    margin: 0 calc(var(--apple-spacing-md) * -1);
    border-radius: 0;
  }

  .table-container {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  .data-table {
    min-width: 500px;
    font-size: var(--apple-text-xs);
  }

  .data-table th,
  .data-table td {
    padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  }

  .config-desc {
    max-width: 150px;
  }

  /* 操作栏 */
  .action-bar {
    flex-direction: column;
    gap: var(--apple-spacing-sm);
    margin: 0 calc(var(--apple-spacing-md) * -1);
    border-radius: 0;
    padding: var(--apple-spacing-lg) var(--apple-spacing-md);
  }

  .action-bar .btn {
    width: 100%;
    min-height: var(--apple-touch-target);
  }

  /* 健康状态网格 */
  .health-status-grid {
    grid-template-columns: 1fr;
  }

  .status-item.full-width {
    grid-column: span 1;
  }

  /* 模态框 */
  .modal-overlay {
    padding: var(--apple-spacing-md);
    align-items: flex-end;
  }

  .modal {
    max-height: 85vh;
    border-radius: var(--apple-radius-xl) var(--apple-radius-xl) 0 0;
    max-width: 100%;
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
}

@media (max-width: 375px) {
  .page-title {
    font-size: var(--apple-text-xl);
  }

  .setting-label {
    font-size: var(--apple-text-xs);
  }

  .config-key {
    font-size: 10px;
  }
}
</style>
