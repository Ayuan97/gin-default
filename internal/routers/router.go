package routers

import (
	"justus/internal/middleware/admin"
	"justus/internal/middleware/api_require"
	"justus/internal/middleware/bodyLog"
	"justus/internal/middleware/cors"
	"justus/internal/middleware/jwt"
	"justus/internal/middleware/recovers"
	tenantmw "justus/internal/middleware/tenant"
	"justus/internal/wire"

	"github.com/gin-gonic/gin"
)

// InitRouterWith 使用依赖注入初始化路由
func InitRouterWith() (*gin.Engine, error) {
	// 组装依赖
	app, err := wire.WireApp()
	if err != nil {
		return nil, err
	}

	r := gin.New()
	// 中间件顺序：Logger -> Recover -> CORS -> BodyLog
	r.Use(gin.Logger(), recovers.Recover(), cors.Cors(), bodyLog.GinBodyLogMiddleware())

	// 健康检查接口（无需认证）
	r.GET("/health", app.HealthController.Health)
	r.GET("/healthz", app.HealthController.Health)
	r.GET("/ready", app.HealthController.Readiness)
	r.GET("/live", app.HealthController.Liveness)

	// API模块路由组
	apiGroup := r.Group("/api/v1")
	// apiGroup.Use(api_require.Common())
	apiGroup.Use(jwt.JWT())
	{
		apiGroup.Any("/test", app.TestController.Test)

		apiGroup.GET("/users", app.UserController.GetUsers)
		apiGroup.GET("/users/:id", app.UserController.GetUser)
		apiGroup.POST("/users", app.UserController.CreateUser)
		apiGroup.PUT("/users/:id", app.UserController.UpdateUser)
		apiGroup.DELETE("/users/:id", app.UserController.DeleteUser)

		apiGroup.GET("/profile", app.UserController.GetProfile)
		apiGroup.PUT("/profile", app.UserController.UpdateProfile)
	}

	// Admin模块路由组 - 面向管理员
	adminGroup := r.Group("/admin/v1")
	adminGroup.Use(api_require.Common())
	adminGroup.Use(jwt.JWT())
	adminGroup.Use(tenantmw.Resolve())
	adminGroup.Use(admin.Auth())
	{
		// 用户管理
		userMgmt := adminGroup.Group("/users")
		{
			userMgmt.GET("", app.UserManagementController.GetUsers)
			userMgmt.GET("/:id", app.UserManagementController.GetUser)
			userMgmt.POST("", app.UserManagementController.CreateUser)
			userMgmt.PUT("/:id", app.UserManagementController.UpdateUser)
			userMgmt.DELETE("/:id", app.UserManagementController.DeleteUser)
			userMgmt.PUT("/:id/status", app.UserManagementController.UpdateUserStatus)
		}

		// 系统管理
		systemMgmt := adminGroup.Group("/system")
		{
			systemMgmt.GET("/info", app.SystemController.GetSystemInfo)
			systemMgmt.GET("/stats", app.SystemController.GetSystemStats)
			systemMgmt.GET("/logs", app.SystemController.GetSystemLogs)
			systemMgmt.GET("/health", app.SystemController.GetHealthStatus)
			systemMgmt.POST("/cache/clear", app.SystemController.ClearCache)
			systemMgmt.POST("/service/:service/restart", app.SystemController.RestartService)
		}

		// 权限管理
		roleMgmt := adminGroup.Group("/roles")
		{
			roleMgmt.GET("", app.RoleController.GetRoles)
			roleMgmt.GET("/:id", app.RoleController.GetRole)
			roleMgmt.POST("", app.RoleController.CreateRole)
			roleMgmt.PUT("/:id", app.RoleController.UpdateRole)
			roleMgmt.DELETE("/:id", app.RoleController.DeleteRole)
			roleMgmt.PUT("/:id/permissions", app.RoleController.UpdateRolePermissions)

			roleMgmt.GET("/permissions", app.RoleController.GetPermissions)
			roleMgmt.POST("/assign", app.RoleController.AssignRole)
		}

		// 菜单相关（租户感知）
		menuMgmt := adminGroup.Group("/menus")
		{
			menuMgmt.GET("", app.MenuController.GetMyMenus)
			menuMgmt.GET("/vben", app.MenuController.GetMyMenusVben)
		}

		// 超级管理员配置租户菜单白名单
		tenantMenuMgmt := adminGroup.Group("/tenants")
		{
			tenantMenuMgmt.GET(":id/menus", app.MenuController.GetTenantMenus)
			tenantMenuMgmt.PUT(":id/menus", app.MenuController.UpdateTenantMenus)
		}

		// 权限码（按钮级）接口
		accessGroup := adminGroup.Group("/access")
		{
			accessGroup.GET("/codes", app.AccessController.GetAccessCodes)
		}

		// 认证相关
		authGroup := adminGroup.Group("/auth")
		{
			authGroup.GET("/profile", app.AuthController.Profile)
		}
	}

	// 兜底路由
	r.NoRoute(func(c *gin.Context) { c.JSON(404, gin.H{"code": 404, "msg": "not found"}) })
	r.NoMethod(func(c *gin.Context) { c.JSON(405, gin.H{"code": 405, "msg": "method not allowed"}) })

	return r, nil
}
