package models

import (
	"fmt"
	"justus/internal/global"
	"justus/pkg/setting"
	"strings"
)

// User 普通用户模型
type User struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID，主键"`
	Username      string    `json:"username" gorm:"uniqueIndex:uk_username;not null;size:50;comment:用户名，唯一标识"`
	Email         string    `json:"email" gorm:"uniqueIndex:uk_email;size:100;comment:邮箱地址，可选"`
	Phone         string    `json:"phone" gorm:"size:20;comment:手机号码"`
	Avatar        string    `json:"avatar" gorm:"size:500;default:'';comment:头像URL地址"`
	FirstName     string    `json:"first_name" gorm:"size:50;default:'';comment:名字（西方习惯）"`
	LastName      string    `json:"last_name" gorm:"size:50;default:'';comment:姓氏（西方习惯）"`
	Nickname      string    `json:"nickname" gorm:"size:50;default:'';comment:昵称"`
	Gender        int       `json:"gender" gorm:"default:0;comment:性别：0-未知，1-男，2-女"`
	Birthday      *GormDate `json:"birthday" gorm:"comment:生日"`
	Lang          string    `json:"lang" gorm:"size:10;default:'zh-Hans';comment:语言偏好：zh-Hans-简体中文，en-英文等"`
	Timezone      string    `json:"timezone" gorm:"size:50;default:'Asia/Shanghai';comment:时区设置"`
	Status        int       `json:"status" gorm:"default:1;comment:用户状态：1-正常，0-禁用，2-待激活;index:idx_status"`
	EmailVerified bool      `json:"email_verified" gorm:"default:false;comment:邮箱验证状态：false-未验证，true-已验证"`
	PhoneVerified bool      `json:"phone_verified" gorm:"default:false;comment:手机验证状态：false-未验证，true-已验证"`
	LastLoginAt   *GormTime `json:"last_login_at" gorm:"comment:最后登录时间;index:idx_last_login"`
	LastLoginIP   string    `json:"last_login_ip" gorm:"size:45;default:'';comment:最后登录IP地址"`
	LoginCount    int       `json:"login_count" gorm:"default:0;comment:登录次数统计"`
	CreatedAt     GormTime  `json:"created_at" gorm:"autoCreateTime;comment:创建时间;index:idx_created_at"`
	UpdatedAt     GormTime  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt     *GormTime `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// UserInfo 普通用户信息结构体
type UserInfo struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	FullName  string `json:"full_name"`
	Lang      string `json:"lang"`
	Status    int    `json:"status"` // 用户状态
}

// 图片地址拼接
func (u *User) getUrl() string {
	if strings.Contains(u.Avatar, "http") {
		return u.Avatar
	} else if u.Avatar != "" {
		return setting.AppSetting.ImageUrl + "/" + u.Avatar
	}
	return ""
}

// Format 格式化普通用户信息
func (u *User) Format() *UserInfo {
	if u.ID <= 0 {
		return nil
	}

	fullName := u.FirstName
	if u.LastName != "" {
		if fullName != "" {
			fullName += " " + u.LastName
		} else {
			fullName = u.LastName
		}
	}

	return &UserInfo{
		ID:        int(u.ID),
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.getUrl(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		FullName:  fullName,
		Lang:      u.Lang,
		Status:    u.Status,
	}
}

// GetUserInfo 获取用户信息
func (u *User) GetUserInfo() (*User, error) {
	var user User
	if u.ID > 0 {
		err := db.Where("id = ?", u.ID).First(&user).Error
		if err != nil {
			global.Logger.Errorf("GetUserInfo error: %v", err)
			return &user, err
		}
	} else {
		return nil, fmt.Errorf("user id is empty")
	}

	return &user, nil
}

// GetUsersByIDs 根据ID列表获取用户
func (u *User) GetUsersByIDs(Ids []int) ([]*User, error) {
	var users []*User
	if len(Ids) > 0 {
		db.Where("id in (?)", Ids).Find(&users)
	}
	return users, nil
}

// GetUsers 获取普通用户列表（API使用）
func GetUsers(page, limit int, keyword, status string) ([]*User, int64, error) {
	var users []*User
	var total int64

	query := db.Model(&User{})

	// 添加搜索条件
	if keyword != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR phone LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// CreateUser 创建普通用户
func (u *User) CreateUser() error {
	err := db.Create(u).Error
	if err != nil {
		global.Logger.Errorf("CreateUser error: %v", err)
		return err
	}
	return nil
}

// UpdateUser 更新用户信息
func (u *User) UpdateUser() error {
	err := db.Save(u).Error
	if err != nil {
		global.Logger.Errorf("UpdateUser error: %v", err)
		return err
	}
	return nil
}

// DeleteUser 删除用户
func (u *User) DeleteUser() error {
	err := db.Delete(u).Error
	if err != nil {
		global.Logger.Errorf("DeleteUser error: %v", err)
		return err
	}
	return nil
}
