package service

import (
	"justus/internal/container"
	"justus/internal/models"
)

// AdminUserServiceImpl 管理员用户服务实现
type AdminUserServiceImpl struct {
	adminUserRepo container.AdminUserRepository
	logger        container.Logger
	cache         container.Cache
}

// NewAdminUserService 创建管理员用户服务实例
func NewAdminUserService(adminUserRepo container.AdminUserRepository, logger container.Logger, cache container.Cache) container.AdminUserService {
	return &AdminUserServiceImpl{
		adminUserRepo: adminUserRepo,
		logger:        logger,
		cache:         cache,
	}
}

// GetAdminUserInfo 获取管理员用户信息
func (s *AdminUserServiceImpl) GetAdminUserInfo(id int) (*models.AdminUser, error) {
	s.logger.Infof("AdminUserService: Getting admin user info for ID: %d", id)

	user, err := s.adminUserRepo.GetByID(id)
	if err != nil {
		s.logger.Errorf("AdminUserService: Failed to get admin user info for ID %d: %v", id, err)
		return nil, err
	}

	s.logger.Debugf("AdminUserService: Successfully retrieved admin user info for ID: %d", id)
	return user, nil
}

// GetByUsername 根据用户名获取管理员用户信息
func (s *AdminUserServiceImpl) GetByUsername(username string) (*models.AdminUser, error) {
	s.logger.Infof("AdminUserService: Getting admin user by username: %s", username)

	user, err := s.adminUserRepo.GetByUsername(username)
	if err != nil {
		s.logger.Errorf("AdminUserService: Failed to get admin user by username %s: %v", username, err)
		return nil, err
	}

	s.logger.Debugf("AdminUserService: Successfully retrieved admin user by username: %s", username)
	return user, nil
}

// CreateAdminUser 创建管理员用户
func (s *AdminUserServiceImpl) CreateAdminUser(user *models.AdminUser) error {
	s.logger.Infof("AdminUserService: Creating admin user: %s", user.Username)

	// 这里可以添加业务逻辑，例如：
	// - 验证管理员权限
	// - 密码加密
	// - 记录创建日志等

	err := s.adminUserRepo.Create(user)
	if err != nil {
		s.logger.Errorf("AdminUserService: Failed to create admin user: %v", err)
		return err
	}

	s.logger.Infof("AdminUserService: Admin user created successfully with ID: %d", user.ID)
	return nil
}

// UpdateAdminUser 更新管理员用户
func (s *AdminUserServiceImpl) UpdateAdminUser(user *models.AdminUser) error {
	s.logger.Infof("AdminUserService: Updating admin user ID: %d", user.ID)

	// 这里可以添加业务逻辑，例如：
	// - 验证更新权限
	// - 数据验证
	// - 记录更新日志等

	err := s.adminUserRepo.Update(user)
	if err != nil {
		s.logger.Errorf("AdminUserService: Failed to update admin user ID %d: %v", user.ID, err)
		return err
	}

	s.logger.Infof("AdminUserService: Admin user ID %d updated successfully", user.ID)
	return nil
}

// DeleteAdminUser 删除管理员用户
func (s *AdminUserServiceImpl) DeleteAdminUser(id int) error {
	s.logger.Infof("AdminUserService: Deleting admin user ID: %d", id)

	// 这里可以添加业务逻辑，例如：
	// - 验证删除权限
	// - 记录删除日志
	// - 清理相关数据等

	err := s.adminUserRepo.Delete(id)
	if err != nil {
		s.logger.Errorf("AdminUserService: Failed to delete admin user ID %d: %v", id, err)
		return err
	}

	s.logger.Infof("AdminUserService: Admin user ID %d deleted successfully", id)
	return nil
}
