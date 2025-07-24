package models

import (
	"justus/internal/global"
)

// Role 角色模型
type Role struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"unique;not null;size:50"`
	DisplayName string `json:"display_name" gorm:"size:100"`
	Description string `json:"description" gorm:"size:255"`
	Status      int    `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Permission 权限模型
type Permission struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"unique;not null;size:100"`
	DisplayName string `json:"display_name" gorm:"size:100"`
	Description string `json:"description" gorm:"size:255"`
	Module      string `json:"module" gorm:"size:50"`   // 模块名称：admin, api
	Action      string `json:"action" gorm:"size:50"`   // 操作类型：read, write, delete
	Resource    string `json:"resource" gorm:"size:50"` // 资源类型：user, role, system
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID       int    `json:"role_id" gorm:"not null"`
	PermissionID int    `json:"permission_id" gorm:"not null"`
	CreatedAt    string `json:"created_at"`
}

// UserRole 用户角色关联模型（已废弃，保留为兼容性）
type UserRole struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int    `json:"user_id" gorm:"not null"`
	RoleID    int    `json:"role_id" gorm:"not null"`
	CreatedAt string `json:"created_at"`
}

// AdminUserRole 管理员用户角色关联模型
type AdminUserRole struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminUserID int    `json:"admin_user_id" gorm:"not null"`
	RoleID      int    `json:"role_id" gorm:"not null"`
	CreatedAt   string `json:"created_at"`
}

// GetRoleInfo 获取角色信息
func (r *Role) GetRoleInfo() (*Role, error) {
	var role Role
	err := db.Where("id = ?", r.ID).First(&role).Error
	if err != nil {
		global.Logger.Errorf("GetRoleInfo error: %v", err)
		return nil, err
	}
	return &role, nil
}

// GetRoleByName 根据名称获取角色
func (r *Role) GetRoleByName() (*Role, error) {
	var role Role
	err := db.Where("name = ?", r.Name).First(&role).Error
	if err != nil {
		global.Logger.Errorf("GetRoleByName error: %v", err)
		return nil, err
	}
	return &role, nil
}

// GetAllRoles 获取所有角色
func (r *Role) GetAllRoles() ([]*Role, error) {
	var roles []*Role
	err := db.Where("status = ?", 1).Find(&roles).Error
	if err != nil {
		global.Logger.Errorf("GetAllRoles error: %v", err)
		return nil, err
	}
	return roles, nil
}

// GetUserRoles 获取普通用户的角色列表（已废弃，保留为兼容性）
func GetUserRoles(userID int) ([]*Role, error) {
	var roles []*Role
	err := db.Table("roles r").
		Select("r.*").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.status = 1", userID).
		Find(&roles).Error

	if err != nil {
		global.Logger.Errorf("GetUserRoles error: %v", err)
		return nil, err
	}
	return roles, nil
}

// GetRolePermissions 获取角色的权限列表
func GetRolePermissions(roleID int) ([]*Permission, error) {
	var permissions []*Permission
	err := db.Table("permissions p").
		Select("p.*").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error

	if err != nil {
		global.Logger.Errorf("GetRolePermissions error: %v", err)
		return nil, err
	}
	return permissions, nil
}

// GetUserPermissions 获取普通用户的所有权限（已废弃，保留为兼容性）
func GetUserPermissions(userID int) ([]*Permission, error) {
	var permissions []*Permission
	err := db.Table("permissions p").
		Select("DISTINCT p.*").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ?", userID).
		Find(&permissions).Error

	if err != nil {
		global.Logger.Errorf("GetUserPermissions error: %v", err)
		return nil, err
	}
	return permissions, nil
}

// GetAdminUserPermissions 获取管理员用户的所有权限
func GetAdminUserPermissions(adminUserID int) ([]*Permission, error) {
	var permissions []*Permission
	err := db.Table("permissions p").
		Select("DISTINCT p.*").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ?", adminUserID).
		Find(&permissions).Error

	if err != nil {
		global.Logger.Errorf("GetAdminUserPermissions error: %v", err)
		return nil, err
	}
	return permissions, nil
}

// HasPermission 检查普通用户是否有指定权限（已废弃，保留为兼容性）
func HasPermission(userID int, permissionName string) (bool, error) {
	var count int64
	err := db.Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN user_roles ur ON rp.role_id = ur.role_id").
		Where("ur.user_id = ? AND p.name = ?", userID, permissionName).
		Count(&count).Error

	if err != nil {
		global.Logger.Errorf("HasPermission error: %v", err)
		return false, err
	}
	return count > 0, nil
}

// HasAdminPermission 检查管理员用户是否有指定权限
func HasAdminPermission(adminUserID int, permissionName string) (bool, error) {
	var count int64
	err := db.Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ? AND p.name = ?", adminUserID, permissionName).
		Count(&count).Error

	if err != nil {
		global.Logger.Errorf("HasAdminPermission error: %v", err)
		return false, err
	}
	return count > 0, nil
}

// IsAdmin 检查普通用户是否为管理员（已废弃，保留为兼容性）
func IsAdmin(userID int) (bool, error) {
	var count int64
	err := db.Table("roles r").
		Joins("JOIN user_roles ur ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.name = 'admin' AND r.status = 1", userID).
		Count(&count).Error

	if err != nil {
		global.Logger.Errorf("IsAdmin error: %v", err)
		return false, err
	}
	return count > 0, nil
}

// IsAdminUser 检查是否为管理员用户
func IsAdminUser(adminUserID int) (bool, error) {
	var count int64
	err := db.Table("roles r").
		Joins("JOIN admin_user_roles aur ON r.id = aur.role_id").
		Where("aur.admin_user_id = ? AND r.name = 'admin' AND r.status = 1", adminUserID).
		Count(&count).Error

	if err != nil {
		global.Logger.Errorf("IsAdminUser error: %v", err)
		return false, err
	}
	return count > 0, nil
}
