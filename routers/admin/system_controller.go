package admin

import (
	"runtime"
	"time"

	"justus/pkg/app"

	"github.com/gin-gonic/gin"
)

// Test 管理员测试接口
func Test(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Success(gin.H{
		"message": "Admin module test successful",
		"module":  "admin",
		"version": "v1",
	})
}

// GetSystemInfo 获取系统信息
func GetSystemInfo(c *gin.Context) {
	appG := app.Gin{C: c}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	appG.Success(gin.H{
		"system": gin.H{
			"go_version":      runtime.Version(),
			"os":              runtime.GOOS,
			"arch":            runtime.GOARCH,
			"cpu_count":       runtime.NumCPU(),
			"goroutine_count": runtime.NumGoroutine(),
			"memory": gin.H{
				"alloc":      m.Alloc / 1024 / 1024,      // MB
				"total":      m.TotalAlloc / 1024 / 1024, // MB
				"sys":        m.Sys / 1024 / 1024,        // MB
				"heap_alloc": m.HeapAlloc / 1024 / 1024,  // MB
				"heap_sys":   m.HeapSys / 1024 / 1024,    // MB
			},
			"server_time": time.Now().Format("2006-01-02 15:04:05"),
		},
	})
}

// GetSystemStats 获取系统统计信息
func GetSystemStats(c *gin.Context) {
	appG := app.Gin{C: c}

	// 这里应该从数据库获取真实的统计数据
	// 示例数据
	appG.Success(gin.H{
		"stats": gin.H{
			"total_users":         1000,
			"active_users":        850,
			"inactive_users":      150,
			"total_sessions":      1250,
			"today_registrations": 25,
			"today_logins":        320,
			"database": gin.H{
				"total_size":   "125.5 MB",
				"tables_count": 15,
				"connections":  5,
			},
			"cache": gin.H{
				"hit_rate":     "95.2%",
				"memory_usage": "45.8 MB",
				"total_keys":   2580,
			},
		},
		"updated_at": time.Now().Format("2006-01-02 15:04:05"),
	})
}
