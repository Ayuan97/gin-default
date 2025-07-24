package api

import (
	"strconv"

	"justus/internal/models"
	"justus/internal/service"
	"justus/pkg/app"
	"justus/pkg/e"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Lang      string `json:"lang"`
	Avatar    string `json:"avatar"`
}

// GetUser 获取用户信息
func GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	userService := &service.UserService{}
	user, err := userService.GetUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	appG.Success(gin.H{
		"user": user.Format(),
	})
}

// GetUsers 获取普通用户列表
func GetUsers(c *gin.Context) {
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
	users, total, err := models.GetUsers(page, limit, keyword, status)
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
func CreateUser(c *gin.Context) {
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

	err := user.CreateUser()
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
func UpdateUser(c *gin.Context) {
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
	userService := &service.UserService{}
	user, err := userService.GetUserInfo(id)
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

	err = user.UpdateUser()
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
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		appG.InvalidParams()
		return
	}

	// 检查用户是否存在
	userService := &service.UserService{}
	user, err := userService.GetUserInfo(id)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	err = user.DeleteUser()
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
func GetProfile(c *gin.Context) {
	appG := app.Gin{C: c}

	// 从JWT中获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		appG.Unauthorized(e.ERROR_AUTH)
		return
	}

	uid := userId.(int)
	userService := &service.UserService{}
	user, err := userService.GetUserInfo(uid)
	if err != nil {
		appG.Error(e.ERROR_USER_NOT_FOUND)
		return
	}

	appG.Success(gin.H{
		"user": user.Format(),
	})
}

// UpdateProfile 更新当前用户个人信息
func UpdateProfile(c *gin.Context) {
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
	userService := &service.UserService{}
	user, err := userService.GetUserInfo(uid)
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

	err = user.UpdateUser()
	if err != nil {
		appG.Error(e.ERROR_DATABASE_UPDATE)
		return
	}

	appG.Success(gin.H{
		"message": "个人信息更新成功",
		"user":    user.Format(),
	})
}
