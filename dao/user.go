package dao

import (
	"justus/global"
	"justus/models"
	"justus/pkg/util"
)

func GetUserInfo(uid int) (*models.User, error) {
	user := models.User{
		Uid: uid,
	}
	return user.GetUserInfo()

}

func GetUserInfoUidKey(userIDs []int) (map[int]interface{}, error) {
	users, _ := GetUsersByIDs(userIDs)
	userUidKey := make(map[int]interface{})
	for _, v := range users {
		userUidKey[v.Uid] = models.User{
			Uid:       v.Uid,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Avatar:    util.GetImageUrl(v.Avatar),
		}
	}
	return userUidKey, nil

}

func GetUsersByIDs(uids []int) ([]*models.User, error) {
	user := models.User{}
	return user.GetUsersByIDs(uids)
}

// 关注
func FollowUser(uid int, followUid int) error {
	tx := models.GetDb().Begin() //开启事务
	var model = models.FollowUser{}
	model.Uid = uid
	model.FollowUid = followUid
	err := model.Create(tx)
	if err != nil {
		global.Logger.Error("create follow_user failed")
		tx.Rollback()
		return err
	}
	var userPageModel = models.UserPage{}
	userPageModel.Uid = followUid
	err = userPageModel.IncUserFansNum(tx)
	if err != nil {
		global.Logger.Error("inc user fans num failed")
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 取消关注
func UnFollowUser(uid int, followUid int) error {
	tx := models.GetDb().Begin() //开启事务
	var model = models.FollowUser{}
	model.Uid = uid
	model.FollowUid = followUid
	err := model.Del()
	if err != nil {
		global.Logger.Error("DEL follow_user failed")
		tx.Rollback()
		return err
	}
	var userPageModel = models.UserPage{}
	userPageModel.Uid = followUid
	err = userPageModel.DecUserFansNum(tx)
	if err != nil {
		global.Logger.Error("DEC user fans num failed")
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
