package models

// Tenant 租户模型（共享表模式）
type Tenant struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:租户ID，主键"`
	Code        string    `json:"code" gorm:"size:50;uniqueIndex:uk_tenant_code;not null;comment:租户编码，唯一"`
	Name        string    `json:"name" gorm:"size:100;not null;comment:租户名称"`
	Status      int       `json:"status" gorm:"default:1;comment:状态：1-启用，0-禁用;index:idx_tenant_status"`
	Plan        string    `json:"plan" gorm:"size:50;default:'';comment:套餐/版本"`
	OwnerUserID uint      `json:"owner_user_id" gorm:"default:0;comment:拥有者管理员ID"`
	CreatedAt   GormTime  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   GormTime  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt   *GormTime `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// TenantPermission 租户权限白名单
type TenantPermission struct {
	ID           uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	TenantID     uint     `json:"tenant_id" gorm:"not null;index:idx_tenant_id;uniqueIndex:uk_tenant_permission,priority:1"`
	PermissionID uint     `json:"permission_id" gorm:"not null;index:idx_permission_id;uniqueIndex:uk_tenant_permission,priority:2"`
	CreatedAt    GormTime `json:"created_at" gorm:"autoCreateTime"`
}

// TableName 映射物理表
func (Tenant) TableName() string           { return "ay_tenants" }
func (TenantPermission) TableName() string { return "ay_tenant_permissions" }

// GetTenantPermissionIDs 获取租户的白名单权限ID集合
func GetTenantPermissionIDs(tenantID uint) ([]uint, error) {
	var ids []uint
	err := db.Model(&TenantPermission{}).
		Where("tenant_id = ?", tenantID).
		Pluck("permission_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetTenantByID 获取租户信息
func GetTenantByID(tenantID uint) (*Tenant, error) {
	var t Tenant
	if err := db.Where("id = ?", tenantID).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
