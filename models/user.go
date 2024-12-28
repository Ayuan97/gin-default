package models

import (
	"fmt"
	"gin-default/global"
	"gin-default/pkg/setting"
	"strings"
)

type User struct {
	Uid             int    `json:"uid"`
	UniqueId        int    `json:"unique_id"`
	Avatar          string `json:"avatar"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Lang            string `json:"lang"`
	OsType          int    `json:"os_type"`
	Phone           string `json:"phone"`
	Pin             string `json:"pin"`
	ComponentNumber int    `json:"component_number"`
	FriendNumber    int    `json:"friend_number"`
	CreatedAt       string `json:"created_at"`
}

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
	FansNum         int    `json:"fans_num"`         //粉丝数量w
}

// 图片地址拼接
func (u *User) getUrl() string {
	if strings.Contains(u.Avatar, "http") {
		return u.Avatar
	} else {
		return setting.AppSetting.ImageUrl + "/" + u.Avatar
	}
}
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
	}
}

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

// GetUsersByIDs Get users by IDs
func (u *User) GetUsersByIDs(Ids []int) ([]*User, error) {
	var users []*User
	if len(Ids) > 0 {
		db.Where("uid in (?)", Ids).Find(&users)
	}
	return users, nil

}
