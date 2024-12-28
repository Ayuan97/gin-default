package topic_service

import (
	"justus/dao"
	"justus/models"
	"justus/service"
	"justus/service/picture_library_service"
)

//获取话题下的图片
func GetTopicPictureList(topicId int, userId int, pageNum int, pageSize int) (*service.ImgList, error) {
	var data service.ImgList

	imgList := make([]*models.PictureLibraryList, 0)
	data.ImgList = imgList

	topicPictureModel := models.TopicPicture{}
	topicPictureModel.TopicId = topicId
	pid := topicPictureModel.GetPublicPidByTopicId((pageNum-1)*pageSize, pageSize)
	if len(pid) > 0 {
		imgData, err := dao.GetPictureLibrarys(pid)
		if err != nil {
			return &data, err
		}
		imgList, err = picture_library_service.PictureLibraryList(imgData, userId, 0)
		if err != nil {
			return &data, err
		}
		data.ImgList = imgList
	}
	return &data, nil

}
