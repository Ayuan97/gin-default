package recommend_service

import (
	"fmt"
	"justus/dao"
	"justus/models"
	"justus/pkg/setting"
	"justus/service"
	"justus/service/picture_library_service"
)

func RecommendList(userId int, page int, lange string) *service.ImgList {
	var data service.ImgList
	var offset int
	if page > 0 {
		offset = (page - 1) * setting.AppSetting.PageSize
	} else {
		offset = 0
	}
	imgList := make([]*models.PictureLibraryList, 0)
	data.ImgList = imgList
	ReadRecordStr := dao.GetRecommendListReadRecord(userId)
	//获取推荐列表
	recommendData, err := dao.GetRecommendList(ReadRecordStr, offset, setting.AppSetting.PageSize, lange)
	if err != nil {
		fmt.Println("err:", err)
		return &data
	}
	if len(recommendData) > 0 {
		var pid []int
		for _, v := range recommendData {
			pid = append(pid, v.PId)
		}
		imgData, err := dao.GetPictureLibrarys(pid)
		if err != nil {
			fmt.Println("err:", err)
			return &data
		}
		imgList, err = picture_library_service.PictureLibraryList(imgData, userId, 0)
		if err != nil {
			fmt.Println("err:", err)
			return &data
		}
		data.ImgList = imgList
	}

	return &data

}
