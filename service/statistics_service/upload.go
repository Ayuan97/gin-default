package statistics_service

import (
	"encoding/json"
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"strconv"
	"strings"
)

const FromTopic = 1   //来源话题
const FromPicture = 2 //来源图片

const Impression = "impression"		//曝光
const Click = "click"				//点击
const Like = "like"					//点赞
const Comment = "comment"			//评论
const Share = "share"				//分享
const FollowPhoto = "follow_photo"	//跟拍


type Data struct {
	DataType string `json:"data_type"`
	Ids string  `json:"ids"`
	Uids string  `json:"uids"`
	Inc int `json:"inc"`
}

//上传数据到集合
//ids 操作的ID集合 多个逗号隔开
//from 来源 话题 还是 图片
//dataType 数据类型 曝光、点击、评论
func UploadStatistics(idString string,uids string,fromType int,dataType string)bool  {
	if dataType != "" && idString != ""{
		return uploadData(idString,uids,fromType,dataType,1)
	}
	return false
}


func UploadStatisticsByIds(ids []int,fromType int,dataType string)bool  {
	if dataType != "" && len(ids) > 0{
		var idString []string
		for _,v := range ids{
			idString = append(idString,strconv.Itoa(v))
		}
		return uploadData(strings.Join(idString,","),"",fromType,dataType,1)

	}
	return false
}

func UploadStatisticsById(id int,fromType int,dataType string)bool  {
	if dataType != "" && id > 0{
		idString := strconv.Itoa(id)
		if fromType == FromPicture && (dataType == Like || dataType == Comment){
			topicPictureModel := models.TopicPicture{}
			topicPictureModel.PId = id
			topicId := topicPictureModel.GetTopicIdsByPId()
			UploadStatisticsByIds(topicId,FromTopic,dataType)
		}
		return uploadData(idString,"",fromType,dataType,1)
	}
	return false
}

//减1
func UploadStatisticsDecById(id int,fromType int,dataType string)bool{
	if dataType != "" && id > 0{
		idString := strconv.Itoa(id)
		if fromType == FromPicture && (dataType == Like || dataType == Comment){
			topicPictureModel := models.TopicPicture{}
			topicPictureModel.PId = id
			topicId := topicPictureModel.GetTopicIdsByPId()
			var topicIdString []string
			for _,v := range topicId{
				topicIdString = append(topicIdString,strconv.Itoa(v))
			}
			uploadData(strings.Join(topicIdString,","),"",FromTopic,dataType,-1)
		}
		return uploadData(idString,"",fromType,dataType,-1)
	}
	return false
}

func uploadData(idString string,uids string,fromType int,dataType string,inc int)bool{
	marshal, err := json.Marshal(Data{dataType, idString,uids,inc})
	if err != nil {
		return  false
	}
	key := ""
	if fromType == FromTopic {
		key = rediskey.StatisticsTopicList
	}
	if fromType == FromPicture {
		key = rediskey.StatisticsPictureList
	}
	if key == ""{
		return false
	}
	_, _ = gredis.LPush(key, marshal)
	return true
}
