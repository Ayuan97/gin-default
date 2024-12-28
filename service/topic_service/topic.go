package topic_service

import (
	"justus/dao"
	"justus/models"
	"justus/service/user_service"
)

type List struct {
	List []Item `json:"list"`
}

type Item struct {
	Id        int    `json:"id"`
	TopicName string `json:"topic_name"`
	HotNum    int    `json:"hot_num"`
}

func GetHotList(lange string, pageNum int, pageSize int, keyword string) (List, error) {
	var result List
	topicModel := models.Topics{}
	list, err := topicModel.GetHotList(pageNum, pageSize, lange, keyword)
	if err != nil {
		return result, err
	}
	if len(list) > 0 {
		var topicIds []int
		for _, v := range list {
			topicIds = append(topicIds, v.ID)
		}
		//获取曝光数
		topicStatisticsModel := models.TopicStatistics{}
		keyList, err2 := topicStatisticsModel.GetImpressionKeyList(topicIds)
		if err2 != nil {
			result.List = make([]Item, 0)
			return result, err2
		}
		for _, v := range list {
			tmpTop := 0
			if _, ok := keyList[v.ID]; ok {
				tmpTop = keyList[v.ID]
			}
			result.List = append(result.List, Item{v.ID, v.TopicName, tmpTop})
		}
	} else {
		result.List = make([]Item, 0)
	}
	return result, nil

}
func GetHotListV2(lange string, pageNum int, pageSize int, keyword string) (List, error) {
	var result List
	topicModel := models.Topics{}
	list, err := topicModel.GethotlistV2(pageNum, pageSize, lange, keyword)
	if err != nil {
		return result, err
	}
	if len(list) > 0 {
		for _, v := range list {
			result.List = append(result.List, Item{v.ID, v.TopicName, 0})
		}
	} else {
		result.List = make([]Item, 0)
	}
	return result, nil

}

type searchList struct {
	List []searchItem `json:"list"`
}

type searchItem struct {
	Id        int    `json:"id"`
	TopicName string `json:"topic_name"`
}

//搜索
func GetSearchList(lange string, pageNum int, pageSize int, keyword string) (searchList, error) {
	var result searchList
	topicModel := models.Topics{}
	list, err := topicModel.GetHotList(pageNum, pageSize, lange, keyword)
	if err != nil {
		return result, err
	}
	if len(list) > 0 {
		for _, v := range list {
			result.List = append(result.List, searchItem{v.ID, v.TopicName})
		}
	} else {
		result.List = make([]searchItem, 0)
	}
	return result, nil

}

func GetWebHotList(lange string, pageNum int, pageSize int) (searchList, error) {
	var result searchList
	topicModel := models.Topics{}
	list, err := topicModel.GetHotList(pageNum, pageSize, lange, "")
	if err != nil {
		return result, err
	}
	if len(list) > 0 {
		for _, v := range list {
			result.List = append(result.List, searchItem{v.ID, v.TopicName})
		}
	} else {
		result.List = make([]searchItem, 0)
	}
	return result, nil

}

//话题详情
type Detail struct {
	Id               int    `json:"id"`
	TopicName        string `json:"topic_name"`
	TopicPicture     string `json:"topic_picture"`
	HotNum           int    `json:"hot_num"`
	FollowPhotoNum   int    `json:"follow_photo_num"`
	Uid              int    `json:"uid"`
	FirstName        string `json:"first_name"`
	UserExist        int `json:"user_exist"`
	CollectStatus    int    `json:"collect_status"`
	FollowUserStatus int    `json:"follow_user_status"`
}

//通过话题ID获取话题详情
func GetDetail(uid int, topicId int) (*Detail, error) {
	var result Detail
	topicModel := models.Topics{}
	topicModel.ID = topicId
	topic, err := topicModel.GetOne()
	if err != nil {
		return nil, err
	}
	if topic.ID == 0 {
		return &result, nil
	}
	result.Id = topic.ID
	result.TopicName = topic.TopicName
	result.TopicPicture = topic.TopicPicture
	topicStatisticsModel := models.TopicStatistics{}
	topicStatisticsModel.TopicId = topicId
	statistics, _ := topicStatisticsModel.GetStatisticsItem()
	result.HotNum = statistics.ImpressionNum
	result.FollowPhotoNum = statistics.FollowPhotoNum
	result.Uid = topic.Uid
	result.UserExist = 0
	user, err := dao.GetUserInfo(topic.Uid)
	if user.Uid > 0 {
		result.FirstName = user.FirstName
		result.UserExist = 1
	}
	result.CollectStatus = GetCollectTopicStatus(uid, topicId)
	result.FollowUserStatus = user_service.GetFollowUserStatus(uid, topic.Uid)
	return &result, nil
}

//通话话题名称 获取 话题信息
func GetDetailByName(uid int, topicName string) (*Detail, error) {
	var result Detail
	topicModel := models.Topics{}
	topicModel.TopicName = topicName
	topic, err := topicModel.GetOneName()
	if err != nil {
		return nil, err
	}
	if topic.ID == 0 {
		return &result, nil
	}
	result.Id = topic.ID
	result.TopicName = topic.TopicName
	result.TopicPicture = topic.TopicPicture
	topicStatisticsModel := models.TopicStatistics{}
	topicStatisticsModel.TopicId = topic.ID
	statistics, _ := topicStatisticsModel.GetStatisticsItem()
	result.HotNum = statistics.ImpressionNum
	result.FollowPhotoNum = statistics.FollowPhotoNum
	result.Uid = topic.Uid

	user, err := dao.GetUserInfo(topic.Uid)
	if user.Uid > 0 {
		result.FirstName = user.FirstName
		result.UserExist = 1
	}
	
	result.CollectStatus = GetCollectTopicStatus(uid, topic.ID)
	result.FollowUserStatus = user_service.GetFollowUserStatus(uid, topic.Uid)
	return &result, nil
}
