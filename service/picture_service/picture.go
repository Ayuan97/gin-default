package picture_service

import (
	"fmt"
	"justus/dao"
	"justus/models"
	"justus/pkg/setting"
	"justus/service"
	"justus/service/picture_library_service"
)

//const SOURCE_FRIEND = 0         //正常好友发的
//const SOURCE_JUSTUS = 1         //just us发的
const SOURCE_FRIEND_BEFORE = 2 //加好友X天前发的
//const SOURCE_GENERATE_VIDEO = 3 //生成的视频
//const SOURCE_TUNE_SUCCESS = 4   //合拍完成的
//const SOURCE_TUNE_FAIL = 5      //合拍未完成的
//const SOURCE_ANONYMOUS = 6      //匿名用户拍的照片

type PostPicture struct {
	PId int `json:"pid"`
}

type PostImgList struct {
	Limit       int    `json:"limit"`
	ScreeningId string `json:"screening_id"`
	Page        int    `json:"page"`
}

type LatestLibrary struct {
	Id     int    `json:"id"`
	Type   int    `json:"type"`
	ImgUrl string `json:"img_url"`
}

func GetLatestLibrary(uid int, friendUids []int) (*LatestLibrary, error) {
	var latestLibrary LatestLibrary
	userPicture, err := models.GetLatestPicture(uid, friendUids)
	if err != nil {
		return nil, err
	}
	data, err := models.GetPictureLibrary(userPicture.PId)
	if err != nil {
		return nil, err
	}
	latestLibrary.Id = data.ID
	latestLibrary.Type = data.Type
	latestLibrary.ImgUrl = setting.AppSetting.ImageUrl + "/" + data.ImgUrl
	return &latestLibrary, nil
}

func GetPicture(p PostPicture, uid int) (*models.PictureLibraryList, error) {
	var picture models.PictureLibrary
	picture.ID = p.PId
	img, err := picture.GetPictureOne()
	if err != nil {
		return nil, err
	}
	if img.ID > 0 {
		PictureList := make([]*models.PictureLibrary, 0)
		PictureList = append(PictureList, img)
		imgData, err := picture_library_service.PictureLibraryList(PictureList, uid, 0)
		if err != nil {
			fmt.Println(err)
		}
		var result *models.PictureLibraryList
		result = imgData[0]
		return result, nil
	}
	return nil, nil

}

func GetPictureList(post PostImgList, uid int) (*service.ImgList, error) {
	var data service.ImgList
	var offset int
	if post.Page > 0 {
		offset = (post.Page - 1) * setting.AppSetting.PageSize
	} else {
		offset = 0
	}
	userPictureLibrary, err := dao.GetUserPictureLibrary(post.ScreeningId, uid, offset, setting.AppSetting.PageSize)
	if err != nil {
		fmt.Println(err)
	}
	var piDs []int
	var userPictureSource map[int]int
	userPictureSource = make(map[int]int)
	for _, value := range userPictureLibrary {
		userPictureSource[value.PId] = value.Source
		piDs = append(piDs, value.PId)
	}
	//fmt.Println("userPictureSource",userPictureSource)
	imgListData, _ := dao.GetPictureLibrarys(piDs)
	imgList, err := picture_library_service.PictureLibraryList(imgListData, uid, 0)
	for _, v := range imgList {
		if userPictureSource[v.ID] == SOURCE_FRIEND_BEFORE {
			v.IsFuzzy = 1
		} else {
			v.IsFuzzy = 0
		}
		v.Source = userPictureSource[v.ID]
	}
	if err != nil {
		fmt.Println("err", err)
	}
	data.ImgList = imgList
	return &data, nil
}
