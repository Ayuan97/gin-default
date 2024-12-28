package models

import (
	"gorm.io/gorm"
)

type FriendList struct {
	Model
	Uid       int `json:"uid"`
	FriendUid int `json:"friend_uid"`
}

func GetFriendList(uid int) ([]*FriendList, error) {
	var friends []*FriendList
	err := db.Select("friend_uid").Where("uid=?", uid).Find(&friends).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return friends, nil
}
