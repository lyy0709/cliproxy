/*
 * 文件作用：强制修改密码中间件
 * 负责功能：
 *   - 检查管理员是否必须修改密码
 *   - 必须修改时仅放行指定接口
 * 重要程度：⭐⭐⭐⭐ 重要（安全核心）
 */
package middleware

import (
	"net/http"

	"cli-proxy/internal/repository"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

// MustChangePasswordGuard 强制修改密码保护
func MustChangePasswordGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		username := c.GetString("username")
		if username == "" {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		repo := repository.NewAdminConfigRepository(repository.DB)
		adminConfig, err := repo.GetByUsername(username)
		if err != nil {
			response.Forbidden(c, "管理员配置不存在")
			c.Abort()
			return
		}

		if !adminConfig.MustChangePassword {
			c.Next()
			return
		}

		path := c.FullPath()
		if path == "/api/admin/status" || path == "/api/admin/force-change-password" {
			c.Next()
			return
		}

		response.Forbidden(c, "首次登录需修改密码")
		c.Abort()
	}
}
