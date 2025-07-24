package dao

import (
	"justus/internal/models"
	"justus/pkg/util"
)

func GetUserInfo(id int) (*models.User, error) {
	user := models.User{
		ID: uint(id),
	}
	return user.GetUserInfo()

}

func GetUserInfoIDKey(userIDs []int) (map[int]interface{}, error) {
	users, _ := GetUsersByIDs(userIDs)
	userIDKey := make(map[int]interface{})
	for _, v := range users {
		userIDKey[int(v.ID)] = models.User{
			ID:        v.ID,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Avatar:    util.GetImageUrl(v.Avatar),
		}
	}
	return userIDKey, nil

}

func GetUsersByIDs(ids []int) ([]*models.User, error) {
	user := models.User{}
	return user.GetUsersByIDs(ids)
}
