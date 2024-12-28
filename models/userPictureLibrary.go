package models

import (
	"fmt"
	"strconv"
	"strings"
)

type UserPictureLibrary struct {
	Model
	PUid       int `json:"p_uid"`
	PId        int `json:"p_id"`
	Uid        int `json:"uid"`
	Source     int `json:"source"`
	PCreatedAt int `json:"p_created_at"`
}

func GetLatestPicture(uid int, friendUids []int) (*UserPictureLibrary, error) {
	var userPictureLibrary UserPictureLibrary
	var friendUidStrs []string
	for _, value := range friendUids {
		friendUidStrs = append(friendUidStrs, strconv.Itoa(value))
	}
	var dbStr = " AND p_uid IN (%s)"
	dbStr = fmt.Sprintf(dbStr, strings.Join(friendUidStrs, ","))
	err := db.Where("uid = ? "+dbStr+" AND source IN (0,4,5)", uid).Order("p_created_at desc").First(&userPictureLibrary).Error
	if err != nil {
		return nil, err
	}
	return &userPictureLibrary, nil
}

func (u UserPictureLibrary) GetImgList(PUid string, offset int, limit int) ([]*UserPictureLibrary, error) {
	var userPictureLibrary []*UserPictureLibrary
	var dbStr = ""
	if PUid != "" {
		dbStr = " AND p_uid IN (%s)"
		dbStr = fmt.Sprintf(PUid)
	} else {
		dbStr = ""
	}
	err := db.Where("uid = ? "+dbStr+" AND source IN (0,2,3,4,5,6)", u.Uid).Order("p_created_at desc").Offset(offset).Limit(limit).Find(&userPictureLibrary).Error
	if err != nil {
		return nil, err
	}
	return userPictureLibrary, nil
}
