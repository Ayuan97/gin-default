package models

import (
	"gorm.io/gorm"
	"justus/global"
)

type FollowUser struct {
	*Model
	Uid       int `json:"uid" gorm:"column:uid"`
	FollowUid int `json:"follow_uid" gorm:"column:follow_uid"`
}

func (f *FollowUser) IsFollow(followUid []int) []*FollowUser {
	var followUser []*FollowUser
	var query = db.Model(&followUser)
	if f.Uid > 0 && len(followUid) > 0 {
		query.Where("uid in (?)", followUid).Where("follow_uid = ? ", f.Uid)
	} else {
		return followUser
	}
	err := query.Find(&followUser).Error
	if err != nil {
		global.Logger.Errorf("IsFollow error: %v", err)
		return followUser
	}
	return followUser
}

//收藏主题
func (m *FollowUser) Create(tx *gorm.DB) error {
	if err := tx.Create(m).Error; err != nil {
		global.Logger.Error("FollowUser Create," + err.Error())
		return err
	}
	return nil
}

//删除
func (m *FollowUser) Del() error {
	if err := db.Where("uid=? AND follow_uid=?", m.Uid, m.FollowUid).Delete(m).Error; err != nil {
		global.Logger.Error("FollowUser Del," + err.Error())
		return err
	}
	return nil
}
