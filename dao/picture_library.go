package dao

import (
	"justus/models"
)

//var PictureLibraryType = map[string]int{
//	"照片": 1,
//	"视频": 2,
//	"合拍": 3,
//}
//var PictureLibraryTuneType = map[int]int{
//	2: 2,
//	3: 2,
//	4: 3,
//	5: 3,
//	6: 3,
//	7: 4,
//}

// GetPictureLibrary 获取用户自己的图片
func GetPictureLibrary(uid int, offset int, limit int) ([]*models.PictureLibrary, error) {
	PictureLibrary := models.PictureLibrary{
		Uid: uid,
	}
	return PictureLibrary.GetImgList(offset, limit)
}

// GetPictureLibrarys  获取多个用户的图片
func GetPictureLibrarys(pid []int) ([]*models.PictureLibrary, error) {
	var list []*models.PictureLibrary
	pictureLibrary := models.PictureLibrary{}
	keyList := pictureLibrary.GetPictureKeyListByIds(pid)
	if len(keyList) > 0 {
		for _, id := range pid {
			if _, ok := keyList[id]; ok {
				list = append(list, keyList[id])
			}
		}
	}

	return list, nil

}
