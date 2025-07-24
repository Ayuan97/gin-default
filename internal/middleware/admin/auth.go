package admin

import (
	"justus/internal/models"
	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

// Auth 管理员权限验证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		// 从JWT中获取用户ID
		userId, exists := c.Get("userId")
		if !exists {
			appG.Unauthorized(e.ERROR_AUTH)
			c.Abort()
			return
		}

		uid := userId.(int)

		// 检查用户是否为管理员
		isAdmin, err := models.IsAdmin(uid)
		if err != nil {
			appG.Error(e.ERROR_DATABASE_QUERY)
			c.Abort()
			return
		}

		if !isAdmin {
			appG.Error(e.ERROR_PERMISSION_DENIED)
			c.Abort()
			return
		}

		// 将管理员信息存储到上下文中
		c.Set("isAdmin", true)
		c.Set("userRole", "admin")

		c.Next()
	}
}

// RequirePermission 检查特定权限的中间件
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		// 从JWT中获取用户ID
		userId, exists := c.Get("userId")
		if !exists {
			appG.Unauthorized(e.ERROR_AUTH)
			c.Abort()
			return
		}

		uid := userId.(int)

		// 检查用户是否有指定权限
		hasPermission, err := models.HasPermission(uid, permission)
		if err != nil {
			appG.Error(e.ERROR_DATABASE_QUERY)
			c.Abort()
			return
		}

		if !hasPermission {
			appG.Error(e.ERROR_INSUFFICIENT_PERMISSION)
			c.Abort()
			return
		}

		c.Next()
	}
}
