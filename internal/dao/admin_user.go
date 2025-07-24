package dao

import (
	"justus/internal/models"
)

// GetAdminUserInfo 获取管理员用户信息
func GetAdminUserInfo(id int) (*models.AdminUser, error) {
	adminUser := models.AdminUser{
		ID: uint(id),
	}
	return adminUser.GetAdminUserInfo()
}

// GetAdminUserByUsername 根据用户名获取管理员用户
func GetAdminUserByUsername(username string) (*models.AdminUser, error) {
	adminUser := models.AdminUser{
		Username: username,
	}
	return adminUser.GetAdminUserByUsername()
}

// CreateAdminUser 创建管理员用户
func CreateAdminUser(adminUser *models.AdminUser) error {
	return adminUser.CreateAdminUser()
}

// UpdateAdminUser 更新管理员用户
func UpdateAdminUser(adminUser *models.AdminUser) error {
	return adminUser.UpdateAdminUser()
}

// DeleteAdminUser 删除管理员用户
func DeleteAdminUser(adminUser *models.AdminUser) error {
	return adminUser.DeleteAdminUser()
}

// GetAdminUsers 获取管理员用户列表
func GetAdminUsers(page, limit int, keyword, status, role string) ([]*models.AdminUser, int64, error) {
	return models.GetAdminUsers(page, limit, keyword, status, role)
}
