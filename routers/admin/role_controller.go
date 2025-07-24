package admin

import (
	"strconv"

	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

type RoleRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// GetRoles 获取角色列表
func GetRoles(c *gin.Context) {
	appG := app.Gin{C: c}

	// 这里应该从数据库获取角色列表
	// 示例数据
	appG.Success(gin.H{
		"roles": []gin.H{
			{
				"id":          1,
				"name":        "admin",
				"description": "系统管理员",
				"permissions": []string{"user.read", "user.write", "user.delete", "system.read", "system.write"},
				"created_at":  "2024-01-01 00:00:00",
			},
			{
				"id":          2,
				"name":        "moderator",
				"description": "内容管理员",
				"permissions": []string{"user.read", "user.write", "content.read", "content.write"},
				"created_at":  "2024-01-01 00:00:00",
			},
			{
				"id":          3,
				"name":        "user",
				"description": "普通用户",
				"permissions": []string{"profile.read", "profile.write"},
				"created_at":  "2024-01-01 00:00:00",
			},
		},
	})
}

// CreateRole 创建角色
func CreateRole(c *gin.Context) {
	appG := app.Gin{C: c}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 检查角色名是否已存在
	// 这里应该调用service层检查

	// 这里应该调用service层创建角色
	// 示例响应
	appG.Success(gin.H{
		"message":     "角色创建成功",
		"role_id":     123,
		"name":        req.Name,
		"description": req.Description,
		"permissions": req.Permissions,
	})
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 检查角色是否存在
	// 这里应该调用service层检查

	// 不能修改admin角色的核心权限
	if id == 1 {
		appG.Error(e.ERROR_PERMISSION_DENIED)
		return
	}

	// 这里应该调用service层更新角色
	// 示例响应
	appG.Success(gin.H{
		"message": "角色更新成功",
		"role_id": id,
	})
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	// 不能删除admin角色
	if id == 1 {
		appG.Error(e.ERROR_PERMISSION_DENIED)
		return
	}

	// 检查是否有用户使用该角色
	// 这里应该调用service层检查

	// 这里应该调用service层删除角色
	// 示例响应
	appG.Success(gin.H{
		"message": "角色删除成功",
		"role_id": id,
	})
}
