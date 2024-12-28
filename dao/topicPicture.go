package dao

import "justus/models"

func GetTopicPidKey(pid []int) map[int][]models.TopicPicture {
	topicsPicture := models.TopicPicture{}
	return topicsPicture.GetTopicInfoByPictureId(pid)
}
