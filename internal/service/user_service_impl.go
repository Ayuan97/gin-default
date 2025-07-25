package service

import (
	"justus/internal/container"
	"justus/internal/models"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	userRepo container.UserRepository
	logger   container.Logger
	cache    container.Cache
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo container.UserRepository, logger container.Logger, cache container.Cache) container.UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
		logger:   logger,
		cache:    cache,
	}
}

// GetUserInfo 获取用户信息
func (s *UserServiceImpl) GetUserInfo(id int) (*models.User, error) {
	s.logger.Infof("UserService: Getting user info for ID: %d", id)

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.Errorf("UserService: Failed to get user info for ID %d: %v", id, err)
		return nil, err
	}

	s.logger.Debugf("UserService: Successfully retrieved user info for ID: %d", id)
	return user, nil
}

// GetUsersByIDs 批量获取用户信息
func (s *UserServiceImpl) GetUsersByIDs(ids []int) ([]*models.User, error) {
	s.logger.Infof("UserService: Getting users info for IDs: %v", ids)

	users, err := s.userRepo.GetByIDs(ids)
	if err != nil {
		s.logger.Errorf("UserService: Failed to get users info for IDs %v: %v", ids, err)
		return nil, err
	}

	s.logger.Debugf("UserService: Successfully retrieved %d users", len(users))
	return users, nil
}

// GetUsers 获取用户列表
func (s *UserServiceImpl) GetUsers(page, limit int, keyword, status string) ([]*models.User, int64, error) {
	s.logger.Infof("UserService: Getting users list with filters")

	users, total, err := s.userRepo.GetUsers(page, limit, keyword, status)
	if err != nil {
		s.logger.Errorf("UserService: Failed to get users list: %v", err)
		return nil, 0, err
	}

	s.logger.Infof("UserService: Successfully retrieved users list - %d users found", len(users))
	return users, total, nil
}

// CreateUser 创建用户
func (s *UserServiceImpl) CreateUser(user *models.User) error {
	s.logger.Infof("UserService: Creating user - %s %s", user.FirstName, user.LastName)

	// 这里可以添加业务逻辑，例如：
	// - 验证用户数据
	// - 密码加密
	// - 发送欢迎邮件等

	err := s.userRepo.Create(user)
	if err != nil {
		s.logger.Errorf("UserService: Failed to create user: %v", err)
		return err
	}

	s.logger.Infof("UserService: User created successfully with ID: %d", user.ID)
	return nil
}

// UpdateUser 更新用户
func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	s.logger.Infof("UserService: Updating user ID: %d", user.ID)

	// 这里可以添加业务逻辑，例如：
	// - 验证更新权限
	// - 数据验证
	// - 缓存更新等

	err := s.userRepo.Update(user)
	if err != nil {
		s.logger.Errorf("UserService: Failed to update user ID %d: %v", user.ID, err)
		return err
	}

	s.logger.Infof("UserService: User ID %d updated successfully", user.ID)
	return nil
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(id int) error {
	s.logger.Infof("UserService: Deleting user ID: %d", id)

	// 这里可以添加业务逻辑，例如：
	// - 验证删除权限
	// - 软删除逻辑
	// - 清理相关数据等

	err := s.userRepo.Delete(id)
	if err != nil {
		s.logger.Errorf("UserService: Failed to delete user ID %d: %v", id, err)
		return err
	}

	s.logger.Infof("UserService: User ID %d deleted successfully", id)
	return nil
}
