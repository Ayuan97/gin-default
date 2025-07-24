package dao

import (
	"justus/internal/models"
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
