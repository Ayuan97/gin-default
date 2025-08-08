package admin

import (
	"strconv"

	"justus/internal/container"
	"justus/internal/models"
	"justus/pkg/app"

	"github.com/gin-gonic/gin"
)

// RoleController 角色管理控制器
type RoleController struct {
	logger container.Logger
	cache  container.Cache
}

// NewRoleController 创建角色管理控制器实例
func NewRoleController(logger container.Logger, cache container.Cache) *RoleController {
	return &RoleController{
		logger: logger,
		cache:  cache,
	}
}

// RoleRequest 角色请求结构体
type RoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	Status      int      `json:"status"`
}

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

// GetRoles 获取角色列表
func (rc *RoleController) GetRoles(c *gin.Context) {
	appG := app.Gin{C: c}

	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		appG.InvalidParams()
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		appG.InvalidParams()
		return
	}

	keyword := c.Query("keyword")
	status := c.Query("status")

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	rc.logger.Infof("Admin getting roles list: tenant_id=%d, page=%d, limit=%d, keyword=%s, status=%s", tenantID, page, limit, keyword, status)

	roles, total, err := models.ListTenantRoles(tenantID, keyword, status, page, limit)
	if err != nil {
		appG.Error(50000)
		return
	}

	appG.Success(gin.H{
		"roles": roles,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
		"filters": gin.H{
			"keyword": keyword,
			"status":  status,
		},
	})
}

// GetRole 获取单个角色详情
func (rc *RoleController) GetRole(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	rc.logger.Infof("Admin getting role details: tenant_id=%d, id=%d", tenantID, id)

	role, err := models.GetRoleByIDForTenant(uint(id), tenantID)
	if err != nil {
		appG.Error(50000)
		return
	}
	permIDs, err := models.GetPermissionIDsOfRole(role.ID)
	if err != nil {
		appG.Error(50000)
		return
	}
	appG.Success(gin.H{"role": role, "permission_ids": permIDs})
}

// CreateRole 创建角色
func (rc *RoleController) CreateRole(c *gin.Context) {
	appG := app.Gin{C: c}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rc.logger.Errorf("Invalid role creation request: %v", err)
		appG.InvalidParams()
		return
	}

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	rc.logger.Infof("Admin creating role: tenant_id=%d, name=%s", tenantID, req.Name)

	// 创建角色
	role, err := models.CreateRoleForTenant(tenantID, req.Name, req.Name, req.Description, req.Status)
	if err != nil {
		appG.Error(50000)
		return
	}
	appG.Success(gin.H{"message": "角色创建成功", "role_id": role.ID})
}

// UpdateRole 更新角色
func (rc *RoleController) UpdateRole(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rc.logger.Errorf("Invalid role update request: %v", err)
		appG.InvalidParams()
		return
	}

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	rc.logger.Infof("Admin updating role: tenant_id=%d, id=%d, name=%s", tenantID, id, req.Name)

	if err := models.UpdateRoleForTenant(uint(id), tenantID, req.Name, req.Description, req.Status); err != nil {
		appG.Error(50000)
		return
	}
	appG.Success(gin.H{"message": "角色更新成功", "role_id": id})
}

// UpdateRolePermissions 覆盖式更新角色权限
func (rc *RoleController) UpdateRolePermissions(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		appG.InvalidParams()
		return
	}

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	// 校验角色归属
	role, err := models.GetRoleByIDForTenant(uint(id), tenantID)
	if err != nil || role == nil {
		appG.Error(50000)
		return
	}

	var req UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	if err := models.ReplaceRolePermissions(uint(id), req.PermissionIDs); err != nil {
		appG.Error(50000)
		return
	}
	appG.Success(gin.H{"message": "权限更新成功", "role_id": id})
}

// DeleteRole 删除角色
func (rc *RoleController) DeleteRole(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	rc.logger.Infof("Admin deleting role: tenant_id=%d, id=%d", tenantID, id)

	if err := models.DeleteRoleForTenant(uint(id), tenantID); err != nil {
		appG.Error(50000)
		return
	}
	appG.Success(gin.H{"message": "角色删除成功", "role_id": id})
}

// GetPermissions 获取所有可用权限列表
func (rc *RoleController) GetPermissions(c *gin.Context) {
	appG := app.Gin{C: c}

	rc.logger.Info("Admin getting permissions list")

	// TODO: 实现获取权限列表逻辑
	permissions := []gin.H{
		{
			"group":      "user",
			"group_name": "用户管理",
			"permissions": []gin.H{
				{"name": "user.read", "display_name": "查看用户", "description": "查看用户列表和详情"},
				{"name": "user.write", "display_name": "编辑用户", "description": "创建、更新用户信息"},
				{"name": "user.delete", "display_name": "删除用户", "description": "删除用户账户"},
				{"name": "user.status", "display_name": "管理用户状态", "description": "启用/禁用用户账户"},
			},
		},
		{
			"group":      "system",
			"group_name": "系统管理",
			"permissions": []gin.H{
				{"name": "system.read", "display_name": "查看系统信息", "description": "查看系统状态和统计"},
				{"name": "system.write", "display_name": "系统配置", "description": "修改系统配置"},
				{"name": "system.logs", "display_name": "查看日志", "description": "查看系统日志"},
				{"name": "system.cache", "display_name": "缓存管理", "description": "清理系统缓存"},
			},
		},
		{
			"group":      "role",
			"group_name": "角色权限",
			"permissions": []gin.H{
				{"name": "role.read", "display_name": "查看角色", "description": "查看角色列表和详情"},
				{"name": "role.write", "display_name": "编辑角色", "description": "创建、更新角色"},
				{"name": "role.delete", "display_name": "删除角色", "description": "删除角色"},
				{"name": "role.assign", "display_name": "分配角色", "description": "为用户分配角色"},
			},
		},
	}

	appG.Success(gin.H{
		"permissions": permissions,
	})
}

// AssignRole 为用户分配角色
func (rc *RoleController) AssignRole(c *gin.Context) {
	appG := app.Gin{C: c}

	var req struct {
		AdminUserID int   `json:"admin_user_id" binding:"required"`
		RoleIDs     []int `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		rc.logger.Errorf("Invalid role assignment request: %v", err)
		appG.InvalidParams()
		return
	}

	// 读取当前租户
	tenantVal, ok := c.Get("tenantId")
	if !ok {
		appG.InvalidParams()
		return
	}
	tenantID := uint(tenantVal.(int))

	rc.logger.Infof("Admin assigning roles: tenant_id=%d, admin_user_id=%d, role_ids=%v", tenantID, req.AdminUserID, req.RoleIDs)

	// 校验角色是否属于该租户或系统级
	roles, err := models.GetRolesByIDsAndTenant(req.RoleIDs, tenantID)
	if err != nil {
		appG.Error(50000)
		return
	}
	if len(roles) != len(req.RoleIDs) {
		appG.InvalidParams()
		return
	}

	// 覆盖式写入管理员在该租户的角色
	roleIDsUint := make([]uint, 0, len(req.RoleIDs))
	for _, rid := range req.RoleIDs {
		roleIDsUint = append(roleIDsUint, uint(rid))
	}
	if err := models.AssignRolesToAdminInTenant(uint(req.AdminUserID), tenantID, roleIDsUint); err != nil {
		appG.Error(50000)
		return
	}

	rc.logger.Infof("Roles assigned successfully: admin_user_id=%d", req.AdminUserID)

	appG.Success(gin.H{
		"message":       "角色分配成功",
		"admin_user_id": req.AdminUserID,
		"role_ids":      req.RoleIDs,
	})
}
