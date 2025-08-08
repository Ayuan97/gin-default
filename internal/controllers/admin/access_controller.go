package admin

import (
	"justus/internal/container"
	"justus/internal/models"
	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

// AccessController 提供权限码等访问控制相关接口
type AccessController struct {
	logger container.Logger
	cache  container.Cache
}

func NewAccessController(logger container.Logger, cache container.Cache) *AccessController {
	return &AccessController{logger: logger, cache: cache}
}

// GetAccessCodes 返回当前管理员在当前租户下拥有的权限名称数组（用于前端按钮级控制）
func (ac *AccessController) GetAccessCodes(c *gin.Context) {
	appG := app.Gin{C: c}

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	userVal, ok := c.Get("userId")
	if !ok {
		appG.Unauthorized(e.ERROR_AUTH)
		return
	}
	adminUserID := userVal.(int)

	codes, err := models.GetAdminUserPermissionNamesInTenant(adminUserID, tenantID)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}
	appG.Success(gin.H{"codes": codes})
}
