package dao

import (
	"justus/models"
)

//获取用户图片集合
func GetUserPictureLibrary(ScreeningId string, uid int, offset int, limit int) ([]*models.UserPictureLibrary, error) {
	UserPictureLibrary := models.UserPictureLibrary{
		Uid: uid,
	}

	return UserPictureLibrary.GetImgList(ScreeningId, offset, limit)
	//return UserPictureLibrary.GetImgList(offset, limit)
}
