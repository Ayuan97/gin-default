package admin

import (
	"strconv"

	"justus/internal/models"
	"justus/internal/service"
	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

type AdminUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name" binding:"required"`
	Role     string `json:"role"`   // 用户角色
	Status   int    `json:"status"` // 用户状态 1:正常 0:禁用
}

// GetUsers 管理员获取管理员用户列表
func GetUsers(c *gin.Context) {
	appG := app.Gin{C: c}

	// 分页参数
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		appG.InvalidParams()
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		appG.InvalidParams()
		return
	}

	// 搜索条件
	keyword := c.Query("keyword")
	status := c.Query("status")
	role := c.Query("role")

	// 调用模型层获取管理员用户列表
	adminUserService := &service.AdminUserService{}
	adminUsers, total, err := adminUserService.GetAdminUsers(page, limit, keyword, status, role)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}

	// 格式化用户信息
	var userList []interface{}
	for _, adminUser := range adminUsers {
		userList = append(userList, adminUser.Format())
	}

	appG.Success(gin.H{
		"users": userList,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
		"filters": gin.H{
			"keyword": keyword,
			"status":  status,
			"role":    role,
		},
	})
}

// GetUser 管理员获取单个管理员用户详情
func GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	adminUserService := &service.AdminUserService{}
	adminUser, err := adminUserService.GetAdminUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	appG.Success(gin.H{
		"user": adminUser.Format(),
	})
}

// CreateUser 管理员创建管理员用户
func CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var req AdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 创建管理员用户
	adminUser := &models.AdminUser{
		Username: req.Username,
		Password: req.Password, // 实际项目中需要加密
		Email:    req.Email,
		Phone:    req.Phone,
		RealName: req.RealName,
		Status:   req.Status,
	}

	adminUserService := &service.AdminUserService{}
	err := adminUserService.CreateAdminUser(adminUser)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_INSERT)
		return
	}

	appG.Success(gin.H{
		"message": "管理员用户创建成功",
		"user_id": adminUser.ID,
		"role":    req.Role,
		"status":  req.Status,
	})
}

// UpdateUser 管理员更新管理员用户信息
func UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	var req AdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 检查管理员用户是否存在
	adminUserService := &service.AdminUserService{}
	adminUser, err := adminUserService.GetAdminUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	// 更新管理员用户信息
	adminUser.Username = req.Username
	if req.Password != "" {
		adminUser.Password = req.Password // 实际项目中需要加密
	}
	adminUser.Email = req.Email
	adminUser.Phone = req.Phone
	adminUser.RealName = req.RealName
	adminUser.Status = req.Status

	err = adminUserService.UpdateAdminUser(adminUser)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_UPDATE)
		return
	}

	appG.Success(gin.H{
		"message": "管理员用户信息更新成功",
		"id":      id,
	})
}

// DeleteUser 管理员删除管理员用户
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	// 检查管理员用户是否存在
	adminUserService := &service.AdminUserService{}
	adminUser, err := adminUserService.GetAdminUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	// 不能删除管理员自己
	currentUserId, _ := c.Get("userId")
	if currentUserId.(int) == id {
		appG.Error(e.ERROR_ADMIN_SELF_OPERATION)
		return
	}

	err = adminUserService.DeleteAdminUser(adminUser)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_DELETE)
		return
	}

	appG.Success(gin.H{
		"message": "管理员用户删除成功",
		"id":      id,
	})
}

// UpdateUserStatus 管理员更新管理员用户状态
func UpdateUserStatus(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 检查管理员用户是否存在
	adminUserService := &service.AdminUserService{}
	adminUser, err := adminUserService.GetAdminUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	// 不能禁用管理员自己
	currentUserId, _ := c.Get("userId")
	if currentUserId.(int) == id && req.Status == 0 {
		appG.Error(e.ERROR_ADMIN_SELF_OPERATION)
		return
	}

	// 更新状态
	adminUser.Status = req.Status
	err = adminUserService.UpdateAdminUser(adminUser)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_UPDATE)
		return
	}

	statusText := "正常"
	if req.Status == 0 {
		statusText = "禁用"
	}

	appG.Success(gin.H{
		"message": "管理员用户状态更新成功",
		"id":      id,
		"status":  statusText,
	})
}
