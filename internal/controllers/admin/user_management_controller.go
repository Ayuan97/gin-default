package admin

import (
	"strconv"

	"justus/internal/container"
	"justus/internal/models"
	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

// UserManagementController 用户管理控制器
type UserManagementController struct {
	userService      container.UserService
	adminUserService container.AdminUserService
	logger           container.Logger
}

// NewUserManagementController 创建用户管理控制器实例
func NewUserManagementController(userService container.UserService, adminUserService container.AdminUserService, logger container.Logger) *UserManagementController {
	return &UserManagementController{
		userService:      userService,
		adminUserService: adminUserService,
		logger:           logger,
	}
}

// UserRequest 用户请求结构体
type UserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Lang      string `json:"lang"`
	Avatar    string `json:"avatar"`
}

// GetUsers 获取用户列表 (管理员)
func (umc *UserManagementController) GetUsers(c *gin.Context) {
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

	// 搜索条件
	keyword := c.Query("keyword")
	status := c.Query("status")

	umc.logger.Infof("Admin getting users list: page=%d, limit=%d, keyword=%s, status=%s", page, limit, keyword, status)

	// 获取用户列表
	users, total, err := umc.userService.GetUsers(page, limit, keyword, status)
	if err != nil {
		umc.logger.Errorf("Failed to get users list: %v", err)
		appG.Error(e.ERROR_DATABASE_QUERY)
		return
	}

	// 格式化用户信息
	var userList []interface{}
	for _, user := range users {
		userList = append(userList, user.Format())
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
		},
	})
}

// GetUser 获取单个用户详情 (管理员)
func (umc *UserManagementController) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	umc.logger.Infof("Admin getting user details: id=%d", id)

	user, err := umc.userService.GetUserInfo(id)
	if err != nil {
		umc.logger.Errorf("Failed to get user details: id=%d, error=%v", id, err)
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	appG.Success(gin.H{
		"user": user.Format(),
	})
}

// CreateUser 创建用户 (管理员)
func (umc *UserManagementController) CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		umc.logger.Errorf("Invalid user creation request: %v", err)
		appG.InvalidParams()
		return
	}

	umc.logger.Infof("Admin creating user: phone=%s", req.Phone)

	// 创建用户
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Lang:      req.Lang,
		Avatar:    req.Avatar,
		Status:    1, // 默认正常状态
	}

	err := umc.userService.CreateUser(user)
	if err != nil {
		umc.logger.Errorf("Failed to create user: %v", err)
		appG.Error(e.ERROR_DATABASE_INSERT)
		return
	}

	umc.logger.Infof("User created successfully: id=%d", user.ID)

	appG.Success(gin.H{
		"message": "用户创建成功",
		"user_id": user.ID,
	})
}

// UpdateUser 更新用户 (管理员)
func (umc *UserManagementController) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		umc.logger.Errorf("Invalid user update request: %v", err)
		appG.InvalidParams()
		return
	}

	umc.logger.Infof("Admin updating user: id=%d", id)

	// 检查用户是否存在
	user, err := umc.userService.GetUserInfo(id)
	if err != nil {
		umc.logger.Errorf("User not found for update: id=%d", id)
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	// 更新用户信息
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.Lang = req.Lang
	user.Avatar = req.Avatar

	err = umc.userService.UpdateUser(user)
	if err != nil {
		umc.logger.Errorf("Failed to update user: id=%d, error=%v", id, err)
		appG.Error(e.ERROR_DATABASE_UPDATE)
		return
	}

	umc.logger.Infof("User updated successfully: id=%d", id)

	appG.Success(gin.H{
		"message": "用户更新成功",
		"user":    user.Format(),
	})
}

// DeleteUser 删除用户 (管理员)
func (umc *UserManagementController) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	umc.logger.Infof("Admin deleting user: id=%d", id)

	// 检查用户是否存在
	_, err = umc.userService.GetUserInfo(id)
	if err != nil {
		umc.logger.Errorf("User not found for deletion: id=%d", id)
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	err = umc.userService.DeleteUser(id)
	if err != nil {
		umc.logger.Errorf("Failed to delete user: id=%d, error=%v", id, err)
		appG.Error(e.ERROR_DATABASE_DELETE)
		return
	}

	umc.logger.Infof("User deleted successfully: id=%d", id)

	appG.Success(gin.H{
		"message": "用户删除成功",
		"id":      id,
	})
}

// UpdateUserStatus 更新用户状态 (管理员)
func (umc *UserManagementController) UpdateUserStatus(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		umc.logger.Errorf("Invalid status update request: %v", err)
		appG.InvalidParams()
		return
	}

	umc.logger.Infof("Admin updating user status: id=%d, status=%d", id, req.Status)

	// 检查用户是否存在
	user, err := umc.userService.GetUserInfo(id)
	if err != nil {
		umc.logger.Errorf("User not found for status update: id=%d", id)
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	// 更新用户状态
	user.Status = req.Status
	err = umc.userService.UpdateUser(user)
	if err != nil {
		umc.logger.Errorf("Failed to update user status: id=%d, error=%v", id, err)
		appG.Error(e.ERROR_DATABASE_UPDATE)
		return
	}

	umc.logger.Infof("User status updated successfully: id=%d, status=%d", id, req.Status)

	appG.Success(gin.H{
		"message": "用户状态更新成功",
		"user":    user.Format(),
	})
}
