package models

import (
	"fmt"
	"justus/internal/global"
	"justus/pkg/setting"
	"strings"
)

// User 普通用户模型
type User struct {
	Uid             int    `json:"uid" gorm:"primaryKey;autoIncrement"`
	UniqueId        int    `json:"unique_id" gorm:"unique"`
	Avatar          string `json:"avatar"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Lang            string `json:"lang"`
	OsType          int    `json:"os_type"`
	Phone           string `json:"phone"`
	Pin             string `json:"pin"`
	ComponentNumber int    `json:"component_number"`
	FriendNumber    int    `json:"friend_number"`
	Status          int    `json:"status" gorm:"default:1"` // 用户状态 1:正常 0:禁用
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// UserInfo 普通用户信息结构体
type UserInfo struct {
	Uid             int    `json:"uid"`
	UniqueId        int    `json:"unique_id"`
	Avatar          string `json:"avatar"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Lang            string `json:"lang"`
	Phone           string `json:"phone"`
	ComponentNumber int    `json:"component_number"`
	FriendNumber    int    `json:"friend_number"`
	HotNum          int    `json:"hot_num"`          //热度
	FollowPhotoNum  int    `json:"follow_photo_num"` //跟拍
	FansNum         int    `json:"fans_num"`         //粉丝数量
	Status          int    `json:"status"`           // 用户状态
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
	if u.Uid <= 0 {
		return nil
	}
	return &UserInfo{
		Uid:             u.Uid,
		UniqueId:        u.UniqueId,
		Avatar:          u.getUrl(),
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Lang:            u.Lang,
		Phone:           u.Phone,
		ComponentNumber: u.ComponentNumber,
		FriendNumber:    u.FriendNumber,
		Status:          u.Status,
	}
}

// GetUserInfo 获取用户信息
func (u *User) GetUserInfo() (*User, error) {
	var user User
	if u.Uid > 0 {
		err := db.Where("uid = ?", u.Uid).First(&user).Error
		if err != nil {
			global.Logger.Errorf("GetUserInfo error: %v db: %v", err, db)
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
		db.Where("uid in (?)", Ids).Find(&users)
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
