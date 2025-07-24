package service

import (
	"justus/internal/dao"
	"justus/internal/models"
)

// AdminUserService 管理员用户服务
type AdminUserService struct{}

// GetAdminUserInfo 获取管理员用户信息
func (s *AdminUserService) GetAdminUserInfo(id int) (*models.AdminUser, error) {
	return dao.GetAdminUserInfo(id)
}

// GetAdminUserByUsername 根据用户名获取管理员用户
func (s *AdminUserService) GetAdminUserByUsername(username string) (*models.AdminUser, error) {
	return dao.GetAdminUserByUsername(username)
}

// CreateAdminUser 创建管理员用户
func (s *AdminUserService) CreateAdminUser(adminUser *models.AdminUser) error {
	return dao.CreateAdminUser(adminUser)
}

// UpdateAdminUser 更新管理员用户
func (s *AdminUserService) UpdateAdminUser(adminUser *models.AdminUser) error {
	return dao.UpdateAdminUser(adminUser)
}

// DeleteAdminUser 删除管理员用户
func (s *AdminUserService) DeleteAdminUser(adminUser *models.AdminUser) error {
	return dao.DeleteAdminUser(adminUser)
}

// GetAdminUsers 获取管理员用户列表
func (s *AdminUserService) GetAdminUsers(page, limit int, keyword, status, role string) ([]*models.AdminUser, int64, error) {
	return dao.GetAdminUsers(page, limit, keyword, status, role)
}
