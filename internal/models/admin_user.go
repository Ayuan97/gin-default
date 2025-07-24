package models

import (
	"fmt"
	"justus/internal/global"
	"justus/pkg/setting"
	"strings"
)

// AdminUser 后台管理用户模型
type AdminUser struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:管理员ID，主键"`
	Username          string    `json:"username" gorm:"uniqueIndex:uk_username;not null;size:50;comment:管理员用户名，唯一标识"`
	Password          string    `json:"-" gorm:"not null;size:255;comment:登录密码，bcrypt加密"`
	Email             string    `json:"email" gorm:"uniqueIndex:uk_email;size:100;comment:邮箱地址"`
	Phone             string    `json:"phone" gorm:"size:20;comment:手机号码"`
	Avatar            string    `json:"avatar" gorm:"size:500;default:'';comment:头像URL地址"`
	RealName          string    `json:"real_name" gorm:"size:50;default:'';comment:真实姓名"`
	Department        string    `json:"department" gorm:"size:100;default:'';comment:所属部门;index:idx_department"`
	Position          string    `json:"position" gorm:"size:50;default:'';comment:职位"`
	Status            int       `json:"status" gorm:"default:1;comment:账户状态：1-正常，0-禁用，2-锁定;index:idx_status"`
	IsSuper           bool      `json:"is_super" gorm:"default:false;comment:是否超级管理员：false-否，true-是"`
	LoginCount        int       `json:"login_count" gorm:"default:0;comment:登录次数统计"`
	LastLoginAt       *GormTime `json:"last_login_at" gorm:"comment:最后登录时间;index:idx_last_login"`
	LastLoginIP       string    `json:"last_login_ip" gorm:"size:45;default:'';comment:最后登录IP地址"`
	PasswordChangedAt *GormTime `json:"password_changed_at" gorm:"comment:密码最后修改时间"`
	FailedLoginCount  int       `json:"failed_login_count" gorm:"default:0;comment:连续登录失败次数"`
	LockedUntil       *GormTime `json:"locked_until" gorm:"comment:账户锁定到期时间"`
	CreatedBy         uint      `json:"created_by" gorm:"default:0;comment:创建者ID"`
	CreatedAt         GormTime  `json:"created_at" gorm:"autoCreateTime;comment:创建时间;index:idx_created_at"`
	UpdatedAt         GormTime  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt         *GormTime `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// AdminUserDetail 管理员用户详细信息结构体
type AdminUserDetail struct {
	ID               uint   `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Avatar           string `json:"avatar"`
	RealName         string `json:"real_name"`
	Department       string `json:"department"`
	Position         string `json:"position"`
	Status           int    `json:"status"`
	IsSuper          bool   `json:"is_super"`
	Role             string `json:"role,omitempty"` // 用户角色
	LastLoginAt      string `json:"last_login_at"`
	LoginCount       int    `json:"login_count"`
	FailedLoginCount int    `json:"failed_login_count"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
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

	// 格式化时间
	lastLoginAt := ""
	if au.LastLoginAt != nil && !au.LastLoginAt.Time.IsZero() {
		lastLoginAt = au.LastLoginAt.Time.Format("2006-01-02 15:04:05")
	}

	return &AdminUserDetail{
		ID:               au.ID,
		Username:         au.Username,
		Email:            au.Email,
		Phone:            au.Phone,
		Avatar:           au.getAvatarUrl(),
		RealName:         au.RealName,
		Department:       au.Department,
		Position:         au.Position,
		Status:           au.Status,
		IsSuper:          au.IsSuper,
		Role:             roleName,
		LastLoginAt:      lastLoginAt,
		LoginCount:       au.LoginCount,
		FailedLoginCount: au.FailedLoginCount,
		CreatedAt:        au.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:        au.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
	}
}

// GetAdminUserRoles 获取管理员用户的角色列表
func GetAdminUserRoles(adminUserID uint) ([]*Role, error) {
	var roles []*Role
	err := db.Table("ay_roles r").
		Select("r.*").
		Joins("JOIN ay_admin_user_roles aur ON r.id = aur.role_id").
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
		query = query.Where("username LIKE ? OR real_name LIKE ? OR email LIKE ? OR phone LIKE ? OR department LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 如果有角色过滤条件
	if role != "" {
		query = query.Joins("JOIN ay_admin_user_roles aur ON ay_admin_users.id = aur.admin_user_id").
			Joins("JOIN ay_roles r ON aur.role_id = r.id").
			Where("r.name = ?", role)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&adminUsers).Error; err != nil {
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
