package wire

import (
	"justus/internal/container"
	"justus/internal/controllers/admin"
	"justus/internal/controllers/api"
	"justus/internal/controllers/common"
	"justus/internal/infrastructure"
	"justus/internal/repository"
	"justus/internal/service"
)

// WireApp 组装应用程序的所有依赖
func WireApp() (*AppContext, error) {
	// 初始化依赖注入容器
	container.InitContainer()

	// 创建基础设施层
	logger := infrastructure.NewLogger()
	cache := infrastructure.NewCache()

	// 创建 Repository 层
	userRepo := repository.NewUserRepository(logger, cache)
	adminUserRepo := repository.NewAdminUserRepository(logger, cache)

	// 创建 Service 层
	userService := service.NewUserService(userRepo, logger, cache)
	adminUserService := service.NewAdminUserService(adminUserRepo, logger, cache)

	// 将服务注册到容器中
	container.GlobalContainer.Logger = logger
	container.GlobalContainer.Cache = cache
	container.GlobalContainer.UserRepo = userRepo
	container.GlobalContainer.AdminUserRepo = adminUserRepo
	container.GlobalContainer.UserService = userService
	container.GlobalContainer.AdminUserService = adminUserService

	// 创建 API 控制器
	userController := api.NewUserController(userService, logger, cache)

	// 创建 Admin 控制器
	userManagementController := admin.NewUserManagementController(userService, adminUserService, logger)
	systemController := admin.NewSystemController(logger, cache)
	roleController := admin.NewRoleController(logger, cache)

	// 创建公共控制器
	healthController := common.NewHealthController(logger, cache)
	testController := common.NewTestController(logger)

	// 创建应用上下文
	app := &AppContext{
		Container: container.GlobalContainer,

		// API 控制器
		UserController: userController,

		// Admin 控制器
		UserManagementController: userManagementController,
		SystemController:         systemController,
		RoleController:           roleController,

		// 公共控制器
		HealthController: healthController,
		TestController:   testController,
	}

	return app, nil
}

// AppContext 应用程序上下文，包含所有注入的依赖
type AppContext struct {
	Container *container.Container

	// API 控制器
	UserController *api.UserController

	// Admin 控制器
	UserManagementController *admin.UserManagementController
	SystemController         *admin.SystemController
	RoleController           *admin.RoleController

	// 公共控制器
	HealthController *common.HealthController
	TestController   *common.TestController
}
