package admin

import (
	"justus/internal/container"
	"justus/internal/models"
	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

// AuthController 提供认证相关接口
type AuthController struct {
	logger container.Logger
	cache  container.Cache
}

func NewAuthController(logger container.Logger, cache container.Cache) *AuthController {
	return &AuthController{logger: logger, cache: cache}
}

// Profile 返回当前管理员与租户信息
func (ac *AuthController) Profile(c *gin.Context) {
	appG := app.Gin{C: c}

	userVal, ok := c.Get("userId")
	if !ok {
		appG.Unauthorized(e.ERROR_AUTH)
		return
	}
	adminUserID := userVal.(int)

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	// 管理员基础信息
	au := &models.AdminUser{ID: uint(adminUserID)}
	adminInfo, err := au.GetAdminUserInfo()
	if err != nil {
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}
	profile := adminInfo.Format()

	// 租户信息
	tenant, err := models.GetTenantByID(tenantID)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}

	appG.Success(gin.H{
		"admin_user": gin.H{
			"id":         profile.ID,
			"username":   profile.Username,
			"real_name":  profile.RealName,
			"email":      profile.Email,
			"avatar":     profile.Avatar,
			"department": profile.Department,
			"position":   profile.Position,
			"is_super":   profile.IsSuper,
			"roles":      profile.Role,
			"last_login": profile.LastLoginAt,
		},
		"tenant": gin.H{
			"id":     tenant.ID,
			"code":   tenant.Code,
			"name":   tenant.Name,
			"status": tenant.Status,
			"plan":   tenant.Plan,
		},
	})
}
