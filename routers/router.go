package routers

import (
	"justus/internal/middleware/admin"
	"justus/internal/middleware/api_require"
	"justus/internal/middleware/bodyLog"
	"justus/internal/middleware/jwt"
	"justus/internal/middleware/recovers"
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
	r.Use(gin.Logger(), recovers.Recover(), bodyLog.GinBodyLogMiddleware())

	// 健康检查接口（无需认证）
	r.GET("/health", app.HealthController.Health)
	r.GET("/healthz", app.HealthController.Health)
	r.GET("/ready", app.HealthController.Readiness)
	r.GET("/live", app.HealthController.Liveness)

	// API模块路由组 - 面向普通用户
	apiGroup := r.Group("/api/v1")
	// apiGroup.Use(api_require.Common())
	// apiGroup.Use(sign.VerifySignature())
	apiGroup.Use(jwt.JWT())
	{
		// 测试接口
		apiGroup.Any("/test", app.TestController.Test)

		// 用户相关接口
		apiGroup.GET("/users", app.UserController.GetUsers)          // 获取用户列表
		apiGroup.GET("/users/:id", app.UserController.GetUser)       // 获取单个用户
		apiGroup.POST("/users", app.UserController.CreateUser)       // 创建用户
		apiGroup.PUT("/users/:id", app.UserController.UpdateUser)    // 更新用户
		apiGroup.DELETE("/users/:id", app.UserController.DeleteUser) // 删除用户

		// 用户个人相关接口
		apiGroup.GET("/profile", app.UserController.GetProfile)    // 获取个人信息
		apiGroup.PUT("/profile", app.UserController.UpdateProfile) // 更新个人信息
	}

	// Admin模块路由组 - 面向管理员
	adminGroup := r.Group("/admin/v1")
	adminGroup.Use(api_require.Common())
	adminGroup.Use(jwt.JWT())
	adminGroup.Use(admin.Auth()) // 管理员权限验证中间件
	{
		// 用户管理
		userMgmt := adminGroup.Group("/users")
		{
			userMgmt.GET("", app.UserManagementController.GetUsers)                    // 获取用户列表
			userMgmt.GET("/:id", app.UserManagementController.GetUser)                 // 获取单个用户详情
			userMgmt.POST("", app.UserManagementController.CreateUser)                 // 创建用户
			userMgmt.PUT("/:id", app.UserManagementController.UpdateUser)              // 更新用户
			userMgmt.DELETE("/:id", app.UserManagementController.DeleteUser)           // 删除用户
			userMgmt.PUT("/:id/status", app.UserManagementController.UpdateUserStatus) // 更新用户状态
		}

		// 系统管理
		systemMgmt := adminGroup.Group("/system")
		{
			systemMgmt.GET("/info", app.SystemController.GetSystemInfo)                       // 获取系统信息
			systemMgmt.GET("/stats", app.SystemController.GetSystemStats)                     // 获取系统统计
			systemMgmt.GET("/logs", app.SystemController.GetSystemLogs)                       // 获取系统日志
			systemMgmt.GET("/health", app.SystemController.GetHealthStatus)                   // 获取健康状态详情
			systemMgmt.POST("/cache/clear", app.SystemController.ClearCache)                  // 清理缓存
			systemMgmt.POST("/service/:service/restart", app.SystemController.RestartService) // 重启服务
		}

		// 权限管理
		roleMgmt := adminGroup.Group("/roles")
		{
			roleMgmt.GET("", app.RoleController.GetRoles)          // 获取角色列表
			roleMgmt.GET("/:id", app.RoleController.GetRole)       // 获取角色详情
			roleMgmt.POST("", app.RoleController.CreateRole)       // 创建角色
			roleMgmt.PUT("/:id", app.RoleController.UpdateRole)    // 更新角色
			roleMgmt.DELETE("/:id", app.RoleController.DeleteRole) // 删除角色

			// 权限相关
			roleMgmt.GET("/permissions", app.RoleController.GetPermissions) // 获取权限列表
			roleMgmt.POST("/assign", app.RoleController.AssignRole)         // 分配角色
		}
	}

	return r, nil
}
