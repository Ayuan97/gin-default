package models

import (
	"fmt"
	"gorm.io/gorm"
	"justus/global"
)

type UserPage struct {
	*Model
	Uid            int `json:"uid" gorm:"column:uid"`
	HotNum         int `json:"hot_num" gorm:"column:hot_num"`
	FollowPhotoNum int `json:"follow_photo_num" gorm:"column:follow_photo_num"`
	FansNum        int `json:"fans_num" gorm:"column:fans_num"`
}

func (up *UserPage) GetUserPage() (*UserPage, error) {
	var userPage UserPage
	if up.Uid > 0 {
		err := db.Where("uid = ?", up.Uid).First(&userPage).Error
		if err != nil {
			return &userPage, err
		}
	} else {
		return nil, fmt.Errorf("user id is empty")
	}
	return &userPage, nil
}



//自增粉丝数量
func (m *UserPage)IncUserFansNum(tx *gorm.DB) error {
	if err := tx.Model(UserPage{}).Where("uid = ? ", m.Uid).UpdateColumn("fans_num", gorm.Expr("fans_num + ?", 1)).Error; err != nil {
		global.Logger.Error("IncUserFansNum,"+err.Error())
		return err
	}
	return nil
}


//减收藏数量
func (m *UserPage)DecUserFansNum(tx *gorm.DB) error {
	if err := tx.Model(UserPage{}).Where("uid = ? ", m.Uid).UpdateColumn("fans_num", gorm.Expr("fans_num - ?", 1)).Error; err != nil {
		global.Logger.Error("DecUserFansNum,"+err.Error())
		return err
	}
	return nil
}

//自增 曝光数
func (m *UserPage)IncUserHotNum(num int) error {
	if err := db.Model(UserPage{}).Where("uid = ? ", m.Uid).UpdateColumn("hot_num", gorm.Expr("hot_num + ?", num)).Error; err != nil {
		global.Logger.Error("IncUserHotNum,"+err.Error())
		return err
	}
	return nil
}