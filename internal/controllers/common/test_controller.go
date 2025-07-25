package common

import (
	"justus/internal/container"
	"justus/pkg/app"

	"github.com/gin-gonic/gin"
)

// TestController 测试控制器
type TestController struct {
	logger container.Logger
}

// NewTestController 创建测试控制器实例
func NewTestController(logger container.Logger) *TestController {
	return &TestController{
		logger: logger,
	}
}

// Test 测试接口
func (tc *TestController) Test(c *gin.Context) {
	appG := app.Gin{C: c}

	tc.logger.Info("API test endpoint accessed")

	appG.Success(gin.H{
		"message": "test ok",
		"version": "v1.0.0",
	})
}
