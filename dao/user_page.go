package dao

import (
	"justus/models"
)

func GetUserPage(uid int) (*models.UserPage, error) {
	userPage := models.UserPage{
		Uid: uid,
	}
	return userPage.GetUserPage()

}
