package models

import (
	"justus/internal/global"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:角色ID，主键"`
	TenantID    uint      `json:"tenant_id" gorm:"default:0;comment:租户ID；0表示系统级角色;index:idx_role_tenant"`
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
	Component   string    `json:"component" gorm:"size:200;default:'';comment:前端组件路径;column:component"`
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
	TenantID    uint      `json:"tenant_id" gorm:"not null;default:0;comment:租户ID;index:idx_admin_role_tenant;uniqueIndex:uk_admin_role,priority:3"`
	AssignedBy  uint      `json:"assigned_by" gorm:"default:0;comment:分配者ID，记录是谁给这个用户分配的角色;index:idx_assigned_by"`
	ExpiresAt   *GormTime `json:"expires_at" gorm:"comment:角色过期时间，NULL表示永不过期;index:idx_expires_at"`
	CreatedAt   GormTime  `json:"created_at" gorm:"autoCreateTime;comment:分配时间"`

	// 外键关联
	AdminUser AdminUser `json:"admin_user,omitempty" gorm:"foreignKey:AdminUserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role      Role      `json:"role,omitempty" gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// 表名映射（使用 ay_ 前缀）
func (Role) TableName() string           { return "ay_roles" }
func (Permission) TableName() string     { return "ay_permissions" }
func (RolePermission) TableName() string { return "ay_role_permissions" }
func (AdminUserRole) TableName() string  { return "ay_admin_user_roles" }

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

// 已废弃的方法移除：GetUserRoles

// GetRolePermissions 获取角色的权限列表
func GetRolePermissions(roleID int) ([]*Permission, error) {
	var permissions []*Permission
	err := db.Table("ay_permissions p").
		Select("p.*").
		Joins("JOIN ay_role_permissions rp ON p.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error

	if err != nil {
		global.Logger.Errorf("GetRolePermissions error: %v", err)
		return nil, err
	}
	return permissions, nil
}

// 已废弃的方法移除：GetUserPermissions

// GetAdminUserPermissions 获取管理员用户的所有权限
func GetAdminUserPermissions(adminUserID int) ([]*Permission, error) {
	var permissions []*Permission
	err := db.Table("ay_permissions p").
		Select("DISTINCT p.*").
		Joins("JOIN ay_role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN ay_admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ?", adminUserID).
		Find(&permissions).Error

	if err != nil {
		global.Logger.Errorf("GetAdminUserPermissions error: %v", err)
		return nil, err
	}
	return permissions, nil
}

// GetAdminUserPermissionIDs 获取管理员用户的所有权限ID（便于快速求交集与构建菜单）
func GetAdminUserPermissionIDs(adminUserID int) ([]uint, error) {
	var ids []uint
	err := db.Table("ay_permissions p").
		Select("DISTINCT p.id").
		Joins("JOIN ay_role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN ay_admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ?", adminUserID).
		Pluck("p.id", &ids).Error
	if err != nil {
		global.Logger.Errorf("GetAdminUserPermissionIDs error: %v", err)
		return nil, err
	}
	return ids, nil
}

// GetAdminUserPermissionIDsInTenant 获取管理员在指定租户的权限ID集合
func GetAdminUserPermissionIDsInTenant(adminUserID int, tenantID uint) ([]uint, error) {
	var ids []uint
	err := db.Table("ay_permissions p").
		Select("DISTINCT p.id").
		Joins("JOIN ay_role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN ay_admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ? AND aur.tenant_id = ?", adminUserID, tenantID).
		Pluck("p.id", &ids).Error
	if err != nil {
		global.Logger.Errorf("GetAdminUserPermissionIDsInTenant error: %v", err)
		return nil, err
	}
	return ids, nil
}

// GetAdminUserPermissionNamesInTenant 获取管理员在指定租户的权限名称集合
func GetAdminUserPermissionNamesInTenant(adminUserID int, tenantID uint) ([]string, error) {
	var names []string
	err := db.Table("ay_permissions p").
		Select("DISTINCT p.name").
		Joins("JOIN ay_role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN ay_admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ? AND aur.tenant_id = ?", adminUserID, tenantID).
		Pluck("p.name", &names).Error
	if err != nil {
		global.Logger.Errorf("GetAdminUserPermissionNamesInTenant error: %v", err)
		return nil, err
	}
	return names, nil
}

// GetRolesByIDsAndTenant 获取指定租户可用的角色（包含系统级角色 tenant_id=0）
func GetRolesByIDsAndTenant(roleIDs []int, tenantID uint) ([]Role, error) {
	if len(roleIDs) == 0 {
		return []Role{}, nil
	}
	var roles []Role
	if err := db.Table("ay_roles").
		Where("id IN ?", roleIDs).
		Where("status = 1").
		Where("tenant_id IN ?", []uint{0, tenantID}).
		Find(&roles).Error; err != nil {
		global.Logger.Errorf("GetRolesByIDsAndTenant error: %v", err)
		return nil, err
	}
	return roles, nil
}

// AssignRolesToAdminInTenant 在指定租户下为管理员设置角色（覆盖式）
func AssignRolesToAdminInTenant(adminUserID uint, tenantID uint, roleIDs []uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("ay_admin_user_roles").
			Where("admin_user_id = ? AND tenant_id = ?", adminUserID, tenantID).
			Delete(&AdminUserRole{}).Error; err != nil {
			return err
		}
		if len(roleIDs) == 0 {
			return nil
		}
		rows := make([]AdminUserRole, 0, len(roleIDs))
		for _, rid := range roleIDs {
			rows = append(rows, AdminUserRole{
				AdminUserID: adminUserID,
				RoleID:      rid,
				TenantID:    tenantID,
			})
		}
		if err := tx.Table("ay_admin_user_roles").Create(&rows).Error; err != nil {
			return err
		}
		return nil
	})
}

// ListTenantRoles 按租户分页查询角色（仅本租户角色，不含系统级）
func ListTenantRoles(tenantID uint, keyword string, status string, page, limit int) ([]Role, int64, error) {
	var (
		roles []Role
		total int64
	)
	query := db.Model(&Role{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR display_name LIKE ?", like, like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id DESC").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

// GetRoleByIDForTenant 获取本租户的角色（不包含系统级）
func GetRoleByIDForTenant(roleID uint, tenantID uint) (*Role, error) {
	var role Role
	if err := db.Where("id = ? AND tenant_id = ?", roleID, tenantID).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetPermissionIDsOfRole 获取角色绑定的权限ID集合
func GetPermissionIDsOfRole(roleID uint) ([]uint, error) {
	var ids []uint
	if err := db.Table("ay_role_permissions").Where("role_id = ?", roleID).Pluck("permission_id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

// GetPermissionsByNames 批量根据名称获取权限
func GetPermissionsByNames(names []string) ([]Permission, error) {
	var list []Permission
	if len(names) == 0 {
		return list, nil
	}
	if err := db.Table("ay_permissions").Where("name IN ?", names).Find(&list).Error; err != nil {
		global.Logger.Errorf("GetPermissionsByNames error: %v", err)
		return nil, err
	}
	return list, nil
}

// CreateRoleForTenant 在指定租户创建角色
func CreateRoleForTenant(tenantID uint, name, displayName, description string, status int) (*Role, error) {
	role := &Role{
		TenantID:    tenantID,
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Status:      status,
	}
	if err := db.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// UpdateRoleForTenant 更新本租户的角色（不允许编辑系统级角色）
func UpdateRoleForTenant(roleID uint, tenantID uint, displayName, description string, status int) error {
	// 只允许更新本租户角色
	return db.Model(&Role{}).
		Where("id = ? AND tenant_id = ?", roleID, tenantID).
		Updates(map[string]interface{}{
			"display_name": displayName,
			"description":  description,
			"status":       status,
		}).Error
}

// DeleteRoleForTenant 删除本租户角色（需无绑定）
func DeleteRoleForTenant(roleID uint, tenantID uint) error {
	// 检查绑定
	var cnt int64
	if err := db.Table("ay_admin_user_roles").Where("role_id = ? AND tenant_id = ?", roleID, tenantID).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return gorm.ErrInvalidData
	}
	// 删除关联权限
	if err := db.Table("ay_role_permissions").Where("role_id = ?", roleID).Delete(&RolePermission{}).Error; err != nil {
		return err
	}
	// 删除角色
	return db.Where("id = ? AND tenant_id = ?", roleID, tenantID).Delete(&Role{}).Error
}

// ReplaceRolePermissions 覆盖式替换角色的权限集合
func ReplaceRolePermissions(roleID uint, permissionIDs []uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("ay_role_permissions").Where("role_id = ?", roleID).Delete(&RolePermission{}).Error; err != nil {
			return err
		}
		if len(permissionIDs) == 0 {
			return nil
		}
		rows := make([]RolePermission, 0, len(permissionIDs))
		for _, pid := range permissionIDs {
			rows = append(rows, RolePermission{RoleID: roleID, PermissionID: pid})
		}
		if err := tx.Table("ay_role_permissions").Create(&rows).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetPermissionByName 根据权限名获取权限
func GetPermissionByName(name string) (*Permission, error) {
	var p Permission
	if err := db.Where("name = ?", name).First(&p).Error; err != nil {
		global.Logger.Errorf("GetPermissionByName error: %v", err)
		return nil, err
	}
	return &p, nil
}

// HasAdminPermissionInTenant 基于租户白名单的权限校验
// 语义：管理员拥有该权限（基于角色） 且 租户白名单允许该权限
func HasAdminPermissionInTenant(adminUserID int, permissionName string, tenantID uint) (bool, error) {
	// 先查权限ID
	perm, err := GetPermissionByName(permissionName)
	if err != nil {
		return false, err
	}

	// 租户白名单必须允许
	var whiteCount int64
	if err := db.Table("ay_tenant_permissions").
		Where("tenant_id = ? AND permission_id = ?", tenantID, perm.ID).
		Count(&whiteCount).Error; err != nil {
		global.Logger.Errorf("HasAdminPermissionInTenant whitelist check error: %v", err)
		return false, err
	}
	if whiteCount == 0 {
		return false, nil
	}

	// 用户是否在该租户拥有该权限（通过租户内角色）
	var count int64
	if err := db.Table("ay_permissions p").
		Joins("JOIN ay_role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN ay_admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ? AND aur.tenant_id = ? AND p.name = ?", adminUserID, tenantID, permissionName).
		Count(&count).Error; err != nil {
		global.Logger.Errorf("HasAdminPermissionInTenant role check error: %v", err)
		return false, err
	}
	return count > 0, nil
}

// 已废弃的方法移除：HasPermission

// HasAdminPermission 检查管理员用户是否有指定权限
func HasAdminPermission(adminUserID int, permissionName string) (bool, error) {
	var count int64
	err := db.Table("ay_permissions p").
		Joins("JOIN ay_role_permissions rp ON p.id = rp.permission_id").
		Joins("JOIN ay_admin_user_roles aur ON rp.role_id = aur.role_id").
		Where("aur.admin_user_id = ? AND p.name = ?", adminUserID, permissionName).
		Count(&count).Error

	if err != nil {
		global.Logger.Errorf("HasAdminPermission error: %v", err)
		return false, err
	}
	return count > 0, nil
}

// 已废弃的方法移除：IsAdmin

// IsAdminUser 检查是否为管理员用户
func IsAdminUser(adminUserID int) (bool, error) {
	var count int64
	err := db.Table("ay_roles r").
		Joins("JOIN ay_admin_user_roles aur ON r.id = aur.role_id").
		Where("aur.admin_user_id = ? AND r.name = 'admin' AND r.status = 1", adminUserID).
		Count(&count).Error

	if err != nil {
		global.Logger.Errorf("IsAdminUser error: %v", err)
		return false, err
	}
	return count > 0, nil
}
