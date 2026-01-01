/**
 * 文件作用：Apple HIG 风格 Toast 消息提示工具
 * 负责功能：
 *   - 显示成功/错误/警告/信息提示
 *   - 自动消失动画
 * 重要程度：⭐⭐⭐ 基础（UI反馈）
 */

let toastContainer = null

function ensureContainer() {
  if (!toastContainer) {
    toastContainer = document.createElement('div')
    toastContainer.className = 'apple-toast-container'
    toastContainer.style.cssText = `
      position: fixed;
      top: 20px;
      left: 50%;
      transform: translateX(-50%);
      z-index: 9999;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8px;
      pointer-events: none;
    `
    document.body.appendChild(toastContainer)
  }
  return toastContainer
}

function createToast(message, type = 'info') {
  const container = ensureContainer()

  const toast = document.createElement('div')
  toast.className = `apple-toast apple-toast-${type}`

  const colors = {
    success: { bg: '#34c759', color: '#fff' },
    error: { bg: '#ff3b30', color: '#fff' },
    warning: { bg: '#ff9500', color: '#fff' },
    info: { bg: '#fff', color: '#1d1d1f' }
  }

  const { bg, color } = colors[type] || colors.info

  toast.style.cssText = `
    padding: 12px 24px;
    background: ${bg};
    color: ${color};
    border-radius: 12px;
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
    font-size: 14px;
    font-weight: 500;
    font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', sans-serif;
    opacity: 0;
    transform: translateY(-20px) scale(0.95);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    pointer-events: auto;
    max-width: 400px;
    text-align: center;
  `

  toast.textContent = message
  container.appendChild(toast)

  // 显示动画
  requestAnimationFrame(() => {
    toast.style.opacity = '1'
    toast.style.transform = 'translateY(0) scale(1)'
  })

  // 自动消失
  setTimeout(() => {
    toast.style.opacity = '0'
    toast.style.transform = 'translateY(-20px) scale(0.95)'
    setTimeout(() => {
      toast.remove()
      // 如果容器为空，移除容器
      if (container.children.length === 0) {
        container.remove()
        toastContainer = null
      }
    }, 300)
  }, 3000)

  return toast
}

export function showMessage(message, type = 'info') {
  return createToast(message, type)
}

// 便捷方法
export const toast = {
  success: (msg) => showMessage(msg, 'success'),
  error: (msg) => showMessage(msg, 'error'),
  warning: (msg) => showMessage(msg, 'warning'),
  info: (msg) => showMessage(msg, 'info')
}

// 兼容 Element Plus 的 ElMessage API
export const ElMessage = {
  success: (msg) => showMessage(typeof msg === 'string' ? msg : msg?.message || '', 'success'),
  error: (msg) => showMessage(typeof msg === 'string' ? msg : msg?.message || '', 'error'),
  warning: (msg) => showMessage(typeof msg === 'string' ? msg : msg?.message || '', 'warning'),
  info: (msg) => showMessage(typeof msg === 'string' ? msg : msg?.message || '', 'info')
}

// 兼容 Element Plus 的 ElMessageBox API
export const ElMessageBox = {
  confirm: (message, title, options = {}) => {
    return new Promise((resolve, reject) => {
      const overlay = document.createElement('div')
      overlay.style.cssText = `
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.4);
        backdrop-filter: blur(4px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 10000;
        opacity: 0;
        transition: opacity 0.3s ease;
      `

      const dialog = document.createElement('div')
      dialog.style.cssText = `
        background: #fff;
        border-radius: 16px;
        box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
        padding: 24px;
        max-width: 400px;
        width: 90%;
        transform: scale(0.95) translateY(20px);
        transition: transform 0.3s ease;
        font-family: -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', sans-serif;
      `

      const iconColor = options.type === 'warning' ? '#ff9500' : options.type === 'error' ? '#ff3b30' : '#007aff'

      dialog.innerHTML = `
        <div style="display: flex; align-items: flex-start; gap: 16px; margin-bottom: 20px;">
          <div style="width: 40px; height: 40px; border-radius: 10px; background: ${iconColor}15; display: flex; align-items: center; justify-content: center; flex-shrink: 0;">
            <svg viewBox="0 0 24 24" fill="none" stroke="${iconColor}" stroke-width="2" style="width: 20px; height: 20px;">
              ${options.type === 'warning' ? '<path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>' : '<circle cx="12" cy="12" r="10"/><line x1="12" y1="16" x2="12" y2="12"/><line x1="12" y1="8" x2="12.01" y2="8"/>'}
            </svg>
          </div>
          <div style="flex: 1;">
            <h3 style="margin: 0 0 8px; font-size: 17px; font-weight: 600; color: #1d1d1f;">${title}</h3>
            <p style="margin: 0; font-size: 14px; color: #86868b; line-height: 1.5;">${message}</p>
          </div>
        </div>
        <div style="display: flex; gap: 12px; justify-content: flex-end;">
          <button class="cancel-btn" style="padding: 10px 20px; font-size: 14px; font-weight: 500; border-radius: 8px; border: none; background: rgba(0,0,0,0.05); color: #1d1d1f; cursor: pointer; transition: background 0.2s;">取消</button>
          <button class="confirm-btn" style="padding: 10px 20px; font-size: 14px; font-weight: 500; border-radius: 8px; border: none; background: ${options.type === 'warning' ? '#ff9500' : '#ff3b30'}; color: #fff; cursor: pointer; transition: background 0.2s;">确定</button>
        </div>
      `

      overlay.appendChild(dialog)
      document.body.appendChild(overlay)

      requestAnimationFrame(() => {
        overlay.style.opacity = '1'
        dialog.style.transform = 'scale(1) translateY(0)'
      })

      const close = (confirmed) => {
        overlay.style.opacity = '0'
        dialog.style.transform = 'scale(0.95) translateY(20px)'
        setTimeout(() => {
          overlay.remove()
          if (confirmed) {
            resolve()
          } else {
            reject(new Error('cancelled'))
          }
        }, 300)
      }

      dialog.querySelector('.cancel-btn').addEventListener('click', () => close(false))
      dialog.querySelector('.confirm-btn').addEventListener('click', () => close(true))
      overlay.addEventListener('click', (e) => {
        if (e.target === overlay) close(false)
      })
    })
  }
}

export default toast
