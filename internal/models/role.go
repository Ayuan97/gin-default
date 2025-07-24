package models

import (
	"justus/internal/global"
)

// Role 角色模型
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:角色ID，主键"`
	Name        string    `json:"name" gorm:"uniqueIndex:uk_name;not null;size:50;comment:角色标识名，英文，如admin、editor"`
	DisplayName string    `json:"display_name" gorm:"not null;size:100;default:'';comment:角色显示名称，中文，如管理员、编辑员"`
	Description string    `json:"description" gorm:"size:500;default:'';comment:角色描述信息"`
	Level       int       `json:"level" gorm:"default:1;comment:角色等级，数字越大权限越高;index:idx_level"`
	Status      int       `json:"status" gorm:"default:1;comment:角色状态：1-启用，0-禁用;index:idx_status"`
	IsSystem    bool      `json:"is_system" gorm:"default:false;comment:是否系统角色：true-系统内置不可删除，false-普通角色"`
	SortOrder   int       `json:"sort_order" gorm:"default:0;comment:排序字段，数字越小越靠前;index:idx_sort_order"`
	CreatedBy   uint      `json:"created_by" gorm:"default:0;comment:创建者ID"`
	CreatedAt   GormTime  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   GormTime  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt   *GormTime `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// Permission 权限模型
type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:权限ID，主键"`
	Name        string    `json:"name" gorm:"uniqueIndex:uk_name;not null;size:100;comment:权限标识名，格式：module.action.resource"`
	DisplayName string    `json:"display_name" gorm:"not null;size:100;default:'';comment:权限显示名称，中文描述"`
	Description string    `json:"description" gorm:"size:500;default:'';comment:权限详细描述"`
	Module      string    `json:"module" gorm:"not null;size:50;default:'';comment:所属模块：admin、api、system等;index:idx_module"`
	Action      string    `json:"action" gorm:"not null;size:50;default:'';comment:操作类型：read、write、delete、create、update等;index:idx_action"`
	Resource    string    `json:"resource" gorm:"not null;size:50;default:'';comment:资源类型：user、role、permission、system等;index:idx_resource"`
	Route       string    `json:"route" gorm:"size:200;default:'';comment:对应的路由规则，支持通配符"`
	Method      string    `json:"method" gorm:"size:20;default:'';comment:HTTP方法：GET、POST、PUT、DELETE、*"`
	ParentID    uint      `json:"parent_id" gorm:"default:0;comment:父权限ID，支持权限树结构;index:idx_parent_id"`
	Level       int       `json:"level" gorm:"default:1;comment:权限层级，根权限为1"`
	SortOrder   int       `json:"sort_order" gorm:"default:0;comment:排序字段，数字越小越靠前;index:idx_sort_order"`
	IsMenu      bool      `json:"is_menu" gorm:"default:false;comment:是否为菜单权限：true-是菜单，false-非菜单;index:idx_is_menu"`
	MenuIcon    string    `json:"menu_icon" gorm:"size:100;default:'';comment:菜单图标class或路径"`
	IsSystem    bool      `json:"is_system" gorm:"default:false;comment:是否系统权限：true-系统内置不可删除，false-普通权限"`
	CreatedAt   GormTime  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   GormTime  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt   *GormTime `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           uint     `json:"id" gorm:"primaryKey;autoIncrement;comment:关联ID，主键"`
	RoleID       uint     `json:"role_id" gorm:"not null;comment:角色ID，外键关联ay_roles.id;index:idx_role_id;uniqueIndex:uk_role_permission,priority:1"`
	PermissionID uint     `json:"permission_id" gorm:"not null;comment:权限ID，外键关联ay_permissions.id;index:idx_permission_id;uniqueIndex:uk_role_permission,priority:2"`
	GrantedBy    uint     `json:"granted_by" gorm:"default:0;comment:授权者ID，记录是谁给这个角色分配的权限;index:idx_granted_by"`
	CreatedAt    GormTime `json:"created_at" gorm:"autoCreateTime;comment:授权时间"`

	// 外键关联
	Role       Role       `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:关联ID，主键"`
	AdminUserID uint      `json:"admin_user_id" gorm:"not null;comment:管理员ID，外键关联ay_admin_users.id;index:idx_admin_user_id;uniqueIndex:uk_admin_role,priority:1"`
	RoleID      uint      `json:"role_id" gorm:"not null;comment:角色ID，外键关联ay_roles.id;index:idx_role_id;uniqueIndex:uk_admin_role,priority:2"`
	AssignedBy  uint      `json:"assigned_by" gorm:"default:0;comment:分配者ID，记录是谁给这个用户分配的角色;index:idx_assigned_by"`
	ExpiresAt   *GormTime `json:"expires_at" gorm:"comment:角色过期时间，NULL表示永不过期;index:idx_expires_at"`
	CreatedAt   GormTime  `json:"created_at" gorm:"autoCreateTime;comment:分配时间"`

	// 外键关联
	AdminUser AdminUser `json:"admin_user,omitempty" gorm:"foreignKey:AdminUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role      Role      `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
