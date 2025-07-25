package api

import (
	"strconv"

	"justus/internal/container"
	"justus/internal/models"
	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService container.UserService
	logger      container.Logger
	cache       container.Cache
}

// NewUserController 创建用户控制器实例
func NewUserController(userService container.UserService, logger container.Logger, cache container.Cache) *UserController {
	return &UserController{
		userService: userService,
		logger:      logger,
		cache:       cache,
	}
}

type UserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Lang      string `json:"lang"`
	Avatar    string `json:"avatar"`
}

// GetUser 获取用户信息
func (uc *UserController) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	user, err := uc.userService.GetUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	appG.Success(gin.H{
		"user": user.Format(),
	})
}

// GetUsers 获取普通用户列表
func (uc *UserController) GetUsers(c *gin.Context) {
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

	// 获取普通用户列表
	users, total, err := uc.userService.GetUsers(page, limit, keyword, status)
	if err != nil {
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

// CreateUser 创建普通用户
func (uc *UserController) CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 创建普通用户
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Lang:      req.Lang,
		Avatar:    req.Avatar,
		Status:    1, // 默认正常状态
	}

	err := uc.userService.CreateUser(user)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_INSERT)
		return
	}

	appG.Success(gin.H{
		"message": "用户创建成功",
		"user_id": user.ID,
	})
}

// UpdateUser 更新用户信息
func (uc *UserController) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 检查用户是否存在
	user, err := uc.userService.GetUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	// 更新用户信息
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.Lang = req.Lang
	user.Avatar = req.Avatar

	err = uc.userService.UpdateUser(user)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_UPDATE)
		return
	}

	appG.Success(gin.H{
		"message": "用户更新成功",
		"user":    user.Format(),
	})
}

// DeleteUser 删除用户
func (uc *UserController) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	// 检查用户是否存在
	_, err = uc.userService.GetUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	err = uc.userService.DeleteUser(id)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_DELETE)
		return
	}

	appG.Success(gin.H{
		"message": "用户删除成功",
		"id":      id,
	})
}

// GetProfile 获取当前用户个人信息
func (uc *UserController) GetProfile(c *gin.Context) {
	appG := app.Gin{C: c}

	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		appG.Unauthorized(e.ERROR_AUTH)
		return
	}

	uid := userId.(int)
	user, err := uc.userService.GetUserInfo(uid)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	appG.Success(gin.H{
		"user": user.Format(),
	})
}

// UpdateProfile 更新当前用户个人信息
func (uc *UserController) UpdateProfile(c *gin.Context) {
	appG := app.Gin{C: c}

	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		appG.Unauthorized(e.ERROR_AUTH)
		return
	}

	uid := userId.(int)

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		appG.InvalidParams()
		return
	}

	// 检查用户是否存在
	user, err := uc.userService.GetUserInfo(uid)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	// 更新个人信息
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.Lang = req.Lang
	user.Avatar = req.Avatar

	err = uc.userService.UpdateUser(user)
	if err != nil {
		appG.Error(e.ERROR_DATABASE_UPDATE)
		return
	}

	appG.Success(gin.H{
		"message": "个人信息更新成功",
		"user":    user.Format(),
	})
}
