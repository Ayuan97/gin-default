package common

import (
	"justus/internal/container"
	"justus/internal/models"
	"justus/pkg/app"
	"justus/pkg/gredis"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct {
	logger container.Logger
	cache  container.Cache
}

// NewHealthController 创建健康检查控制器实例
func NewHealthController(logger container.Logger, cache container.Cache) *HealthController {
	return &HealthController{
		logger: logger,
		cache:  cache,
	}
}

// Health 健康检查接口
func (hc *HealthController) Health(c *gin.Context) {
	appG := app.Gin{C: c}

	// 检查各项服务状态
	status := hc.checkServicesHealth()

	if status["overall"] == "healthy" {
		c.Status(200)
		appG.Success(gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"services":  status,
		})
	} else {
		c.Status(503)
		appG.Error(503)
	}
}

// Readiness 就绪检查接口
func (hc *HealthController) Readiness(c *gin.Context) {
	appG := app.Gin{C: c}

	// 检查应用是否已准备好处理请求
	ready := hc.checkReadiness()

	if ready {
		appG.Success(gin.H{
			"status":    "ready",
			"timestamp": time.Now().Unix(),
		})
	} else {
		c.Status(503)
		appG.Error(503)
	}
}

// Liveness 存活检查接口
func (hc *HealthController) Liveness(c *gin.Context) {
	appG := app.Gin{C: c}

	// 简单的存活检查
	appG.Success(gin.H{
		"status":    "alive",
		"timestamp": time.Now().Unix(),
	})
}

// checkServicesHealth 检查各项服务健康状态
func (hc *HealthController) checkServicesHealth() map[string]string {
	status := make(map[string]string)
	overall := "healthy"

	// 检查数据库连接
	db := models.GetDb()
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			status["database"] = "unhealthy"
			overall = "unhealthy"
			hc.logger.Errorf("Database health check failed: %v", err)
		} else {
			status["database"] = "healthy"
		}
	} else {
		status["database"] = "unhealthy"
		overall = "unhealthy"
		hc.logger.Error("Database connection is nil")
	}

	// 检查Redis连接
	if gredis.Get("health_check") == "" {
		// 尝试设置一个临时键值
		err := gredis.Set("health_check", "ok", time.Minute)
		if err != nil {
			status["redis"] = "unhealthy"
			overall = "unhealthy"
			hc.logger.Errorf("Redis health check failed: %v", err)
		} else {
			status["redis"] = "healthy"
		}
	} else {
		status["redis"] = "healthy"
	}

	status["overall"] = overall
	return status
}

// checkReadiness 检查应用就绪状态
func (hc *HealthController) checkReadiness() bool {
	// 检查数据库是否可用
	db := models.GetDb()
	if db == nil {
		hc.logger.Error("Database connection is nil in readiness check")
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		hc.logger.Errorf("Failed to get SQL DB in readiness check: %v", err)
		return false
	}

	if err := sqlDB.Ping(); err != nil {
		hc.logger.Errorf("Database ping failed in readiness check: %v", err)
		return false
	}

	// 检查Redis是否可用
	if err := gredis.Set("readiness_check", "ok", time.Minute); err != nil {
		hc.logger.Errorf("Redis connection failed in readiness check: %v", err)
		return false
	}

	return true
}
