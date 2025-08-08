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

		// 优先允许超级管理员
		if isSuper, _ := c.Get("isSuper"); isSuper == true {
			c.Set("isAdmin", true)
			c.Set("userRole", "super_admin")
			c.Next()
			return
		}

		// 校验是否管理员用户
		isAdminUser, err := models.IsAdminUser(uid)
		if err != nil {
			appG.Error(e.ERROR_DATABASE_QUERY)
			c.Abort()
			return
		}
		if !isAdminUser {
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

		// 超级管理员直接放行
		if isSuper, _ := c.Get("isSuper"); isSuper == true {
			c.Next()
			return
		}

		// 读取上下文租户ID
		tenantVal, ok := c.Get("tenantId")
		if !ok {
			appG.Error(e.INVALID_PARAMS)
			c.Abort()
			return
		}
		tenantID := uint(tenantVal.(int))

		// 基于租户白名单的权限校验
		hasPermission, err := models.HasAdminPermissionInTenant(uid, permission, tenantID)
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
