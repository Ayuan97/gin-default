package routers

import (
	"justus/internal/middleware/admin"
	"justus/internal/middleware/api_require"
	"justus/internal/middleware/bodyLog"
	"justus/internal/middleware/jwt"
	"justus/internal/middleware/recovers"
	"justus/internal/middleware/sign"
	adminController "justus/routers/admin"
	"justus/routers/api"

	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), recovers.Recover(), bodyLog.GinBodyLogMiddleware())

	// 健康检查接口（无需认证）
	r.GET("/health", api.Health)
	r.GET("/healthz", api.Health)
	r.GET("/ready", api.Readiness)
	r.GET("/live", api.Liveness)

	// API模块路由组
	apiGroup := r.Group("/api/v1")
	apiGroup.Use(api_require.Common())
	apiGroup.Use(sign.VerifySignature())
	apiGroup.Use(jwt.JWT())
	{
		// 测试接口
		apiGroup.Any("/test", api.Test)

		// 用户相关接口
		apiGroup.GET("/users", api.GetUsers)          // 获取用户列表
		apiGroup.GET("/users/:id", api.GetUser)       // 获取单个用户
		apiGroup.POST("/users", api.CreateUser)       // 创建用户
		apiGroup.PUT("/users/:id", api.UpdateUser)    // 更新用户
		apiGroup.DELETE("/users/:id", api.DeleteUser) // 删除用户

		// 用户个人相关接口
		apiGroup.GET("/profile", api.GetProfile)    // 获取个人信息
		apiGroup.PUT("/profile", api.UpdateProfile) // 更新个人信息
	}

	// Admin模块路由组 - 面向管理员
	adminGroup := r.Group("/admin/v1")
	adminGroup.Use(api_require.Common())
	adminGroup.Use(jwt.JWT())
	adminGroup.Use(admin.Auth()) // 管理员权限验证中间件
	{
		// 管理员测试接口
		adminGroup.POST("/test", adminController.Test)

		// 用户管理
		adminGroup.GET("/users", adminController.GetUsers)                    // 获取用户列表
		adminGroup.GET("/users/:id", adminController.GetUser)                 // 获取单个用户详情
		adminGroup.POST("/users", adminController.CreateUser)                 // 创建用户
		adminGroup.PUT("/users/:id", adminController.UpdateUser)              // 更新用户
		adminGroup.DELETE("/users/:id", adminController.DeleteUser)           // 删除用户
		adminGroup.PUT("/users/:id/status", adminController.UpdateUserStatus) // 更新用户状态

		// 系统管理
		adminGroup.GET("/system/info", adminController.GetSystemInfo)   // 获取系统信息
		adminGroup.GET("/system/stats", adminController.GetSystemStats) // 获取系统统计

		// 权限管理
		adminGroup.GET("/roles", adminController.GetRoles)          // 获取角色列表
		adminGroup.POST("/roles", adminController.CreateRole)       // 创建角色
		adminGroup.PUT("/roles/:id", adminController.UpdateRole)    // 更新角色
		adminGroup.DELETE("/roles/:id", adminController.DeleteRole) // 删除角色
	}

	return r
}
