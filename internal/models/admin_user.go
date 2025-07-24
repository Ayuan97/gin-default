package models

import (
	"fmt"
	"justus/internal/global"
	"justus/pkg/setting"
	"strings"
)

// AdminUser 后台管理用户模型
type AdminUser struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"unique;not null;size:50"`
	Password  string `json:"-" gorm:"not null;size:255"`
	Email     string `json:"email" gorm:"unique;size:100"`
	Phone     string `json:"phone" gorm:"size:20"`
	Avatar    string `json:"avatar" gorm:"size:255"`
	RealName  string `json:"real_name" gorm:"size:50"`
	Status    int    `json:"status" gorm:"default:1"` // 1:正常 0:禁用
	LastLogin string `json:"last_login"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// AdminUserDetail 管理员用户详细信息结构体
type AdminUserDetail struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	RealName  string `json:"real_name"`
	Status    int    `json:"status"`
	Role      string `json:"role,omitempty"` // 用户角色
	LastLogin string `json:"last_login"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 获取头像完整URL
func (au *AdminUser) getAvatarUrl() string {
	if strings.Contains(au.Avatar, "http") {
		return au.Avatar
	} else if au.Avatar != "" {
		return setting.AppSetting.ImageUrl + "/" + au.Avatar
	}
	return ""
}

// Format 格式化管理员用户信息
func (au *AdminUser) Format() *AdminUserDetail {
	if au.ID <= 0 {
		return nil
	}

	// 获取用户角色
	roles, _ := GetAdminUserRoles(au.ID)
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}
	roleName := ""
	if len(roleNames) > 0 {
		roleName = strings.Join(roleNames, ",")
	}

	return &AdminUserDetail{
		ID:        au.ID,
		Username:  au.Username,
		Email:     au.Email,
		Phone:     au.Phone,
		Avatar:    au.getAvatarUrl(),
		RealName:  au.RealName,
		Status:    au.Status,
		Role:      roleName,
		LastLogin: au.LastLogin,
		CreatedAt: au.CreatedAt,
		UpdatedAt: au.UpdatedAt,
	}
}

// GetAdminUserRoles 获取管理员用户的角色列表
func GetAdminUserRoles(adminUserID int) ([]*Role, error) {
	var roles []*Role
	err := db.Table("roles r").
		Select("r.*").
		Joins("JOIN admin_user_roles aur ON r.id = aur.role_id").
		Where("aur.admin_user_id = ? AND r.status = 1", adminUserID).
		Find(&roles).Error

	if err != nil {
		global.Logger.Errorf("GetAdminUserRoles error: %v", err)
		return nil, err
	}
	return roles, nil
}

// GetAdminUserInfo 获取管理员用户信息
func (au *AdminUser) GetAdminUserInfo() (*AdminUser, error) {
	var adminUser AdminUser
	if au.ID > 0 {
		err := db.Where("id = ?", au.ID).First(&adminUser).Error
		if err != nil {
			global.Logger.Errorf("GetAdminUserInfo error: %v", err)
			return &adminUser, err
		}
	} else {
		return nil, fmt.Errorf("admin user id is empty")
	}
	return &adminUser, nil
}

// GetAdminUserByUsername 根据用户名获取管理员用户
func (au *AdminUser) GetAdminUserByUsername() (*AdminUser, error) {
	var adminUser AdminUser
	err := db.Where("username = ?", au.Username).First(&adminUser).Error
	if err != nil {
		global.Logger.Errorf("GetAdminUserByUsername error: %v", err)
		return nil, err
	}
	return &adminUser, nil
}

// GetAdminUsers 获取管理员用户列表
func GetAdminUsers(page, limit int, keyword, status, role string) ([]*AdminUser, int64, error) {
	var adminUsers []*AdminUser
	var total int64

	query := db.Model(&AdminUser{})

	// 添加搜索条件
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 如果有角色过滤条件
	if role != "" {
		query = query.Joins("JOIN admin_user_roles aur ON admin_users.id = aur.admin_user_id").
			Joins("JOIN roles r ON aur.role_id = r.id").
			Where("r.name = ?", role)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&adminUsers).Error; err != nil {
		return nil, 0, err
	}

	return adminUsers, total, nil
}

// CreateAdminUser 创建管理员用户
func (au *AdminUser) CreateAdminUser() error {
	err := db.Create(au).Error
	if err != nil {
		global.Logger.Errorf("CreateAdminUser error: %v", err)
		return err
	}
	return nil
}

// UpdateAdminUser 更新管理员用户
func (au *AdminUser) UpdateAdminUser() error {
	err := db.Save(au).Error
	if err != nil {
		global.Logger.Errorf("UpdateAdminUser error: %v", err)
		return err
	}
	return nil
}

// DeleteAdminUser 删除管理员用户
func (au *AdminUser) DeleteAdminUser() error {
	err := db.Delete(au).Error
	if err != nil {
		global.Logger.Errorf("DeleteAdminUser error: %v", err)
		return err
	}
	return nil
}
