package models

import "gorm.io/gorm"

// WithTenant 返回带有 tenant_id 过滤条件的 DB 会话
// 注意：仅对包含 tenant_id 字段的表生效，像全局的 permissions 表不应使用该作用域
func WithTenant(db *gorm.DB, tenantID uint) *gorm.DB {
	return db.Where("tenant_id = ?", tenantID)
}
