package repository

import (
	"justus/internal/container"
	"justus/internal/models"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	logger container.Logger
	cache  container.Cache
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(logger container.Logger, cache container.Cache) container.UserRepository {
	return &UserRepositoryImpl{
		logger: logger,
		cache:  cache,
	}
}

// GetByID 根据ID获取用户信息
func (r *UserRepositoryImpl) GetByID(id int) (*models.User, error) {
	r.logger.Infof("Getting user by ID: %d", id)

	user := models.User{
		ID: uint(id),
	}
	result, err := user.GetUserInfo()

	if err != nil {
		r.logger.Errorf("Failed to get user by ID %d: %v", id, err)
	} else {
		r.logger.Debugf("Successfully retrieved user: %d", id)
	}

	return result, err
}

// GetByIDs 批量获取用户信息
func (r *UserRepositoryImpl) GetByIDs(ids []int) ([]*models.User, error) {
	r.logger.Infof("Getting users by IDs: %v", ids)

	user := models.User{}
	result, err := user.GetUsersByIDs(ids)

	if err != nil {
		r.logger.Errorf("Failed to get users by IDs %v: %v", ids, err)
	} else {
		r.logger.Debugf("Successfully retrieved %d users", len(result))
	}

	return result, err
}

// GetUsers 获取用户列表
func (r *UserRepositoryImpl) GetUsers(page, limit int, keyword, status string) ([]*models.User, int64, error) {
	r.logger.Infof("Getting users list - page: %d, limit: %d, keyword: %s, status: %s", page, limit, keyword, status)

	result, total, err := models.GetUsers(page, limit, keyword, status)

	if err != nil {
		r.logger.Errorf("Failed to get users list: %v", err)
	} else {
		r.logger.Debugf("Successfully retrieved %d users out of %d total", len(result), total)
	}

	return result, total, err
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(user *models.User) error {
	r.logger.Infof("Creating user: %s %s", user.FirstName, user.LastName)

	err := user.CreateUser()

	if err != nil {
		r.logger.Errorf("Failed to create user: %v", err)
	} else {
		r.logger.Infof("Successfully created user with ID: %d", user.ID)
		// 可以在这里清理相关缓存
		// r.cache.Del(fmt.Sprintf("user:%d", user.ID))
	}

	return err
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(user *models.User) error {
	r.logger.Infof("Updating user ID: %d", user.ID)

	err := user.UpdateUser()

	if err != nil {
		r.logger.Errorf("Failed to update user ID %d: %v", user.ID, err)
	} else {
		r.logger.Infof("Successfully updated user ID: %d", user.ID)
		// 可以在这里清理相关缓存
		// r.cache.Del(fmt.Sprintf("user:%d", user.ID))
	}

	return err
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(id int) error {
	r.logger.Infof("Deleting user ID: %d", id)

	user := models.User{ID: uint(id)}
	err := user.DeleteUser()

	if err != nil {
		r.logger.Errorf("Failed to delete user ID %d: %v", id, err)
	} else {
		r.logger.Infof("Successfully deleted user ID: %d", id)
		// 可以在这里清理相关缓存
		// r.cache.Del(fmt.Sprintf("user:%d", id))
	}

	return err
}
