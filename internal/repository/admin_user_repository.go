package repository

import (
	"justus/internal/container"
	"justus/internal/models"
)

// AdminUserRepositoryImpl 管理员用户仓储实现
type AdminUserRepositoryImpl struct {
	logger container.Logger
	cache  container.Cache
}

// NewAdminUserRepository 创建管理员用户仓储实例
func NewAdminUserRepository(logger container.Logger, cache container.Cache) container.AdminUserRepository {
	return &AdminUserRepositoryImpl{
		logger: logger,
		cache:  cache,
	}
}

// GetByID 根据ID获取管理员用户信息
func (r *AdminUserRepositoryImpl) GetByID(id int) (*models.AdminUser, error) {
	r.logger.Infof("Getting admin user by ID: %d", id)

	adminUser := models.AdminUser{
		ID: uint(id),
	}
	result, err := adminUser.GetAdminUserInfo()

	if err != nil {
		r.logger.Errorf("Failed to get admin user by ID %d: %v", id, err)
	} else {
		r.logger.Debugf("Successfully retrieved admin user: %d", id)
	}

	return result, err
}

// GetByUsername 根据用户名获取管理员用户信息
func (r *AdminUserRepositoryImpl) GetByUsername(username string) (*models.AdminUser, error) {
	r.logger.Infof("Getting admin user by username: %s", username)

	adminUser := models.AdminUser{
		Username: username,
	}
	result, err := adminUser.GetAdminUserByUsername()

	if err != nil {
		r.logger.Errorf("Failed to get admin user by username %s: %v", username, err)
	} else {
		r.logger.Debugf("Successfully retrieved admin user by username: %s", username)
	}

	return result, err
}

// Create 创建管理员用户
func (r *AdminUserRepositoryImpl) Create(user *models.AdminUser) error {
	r.logger.Infof("Creating admin user: %s", user.Username)

	err := user.CreateAdminUser()

	if err != nil {
		r.logger.Errorf("Failed to create admin user: %v", err)
	} else {
		r.logger.Infof("Successfully created admin user with ID: %d", user.ID)
	}

	return err
}

// Update 更新管理员用户
func (r *AdminUserRepositoryImpl) Update(user *models.AdminUser) error {
	r.logger.Infof("Updating admin user ID: %d", user.ID)

	err := user.UpdateAdminUser()

	if err != nil {
		r.logger.Errorf("Failed to update admin user ID %d: %v", user.ID, err)
	} else {
		r.logger.Infof("Successfully updated admin user ID: %d", user.ID)
	}

	return err
}

// Delete 删除管理员用户
func (r *AdminUserRepositoryImpl) Delete(id int) error {
	r.logger.Infof("Deleting admin user ID: %d", id)

	adminUser := models.AdminUser{ID: uint(id)}
	err := adminUser.DeleteAdminUser()

	if err != nil {
		r.logger.Errorf("Failed to delete admin user ID %d: %v", id, err)
	} else {
		r.logger.Infof("Successfully deleted admin user ID: %d", id)
	}

	return err
}
