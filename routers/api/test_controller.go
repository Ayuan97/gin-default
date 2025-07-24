package api

import (
	"justus/pkg/app"

	"github.com/gin-gonic/gin"
)

// Test 测试接口
func Test(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Success(gin.H{
		"message": "test ok",
		"version": "v1.0.0",
	})
}
