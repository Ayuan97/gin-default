package service

import (
	"justus/internal/dao"
	"justus/internal/models"
)

// UserService 用户服务
type UserService struct{}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(uid int) (*models.User, error) {
	return dao.GetUserInfo(uid)
}

// GetUsersByIDs 批量获取用户信息
func (s *UserService) GetUsersByIDs(uids []int) ([]*models.User, error) {
	return dao.GetUsersByIDs(uids)
}

// TODO: 这里可以添加更多用户相关的业务逻辑
// 比如用户注册、登录、更新信息等功能
