package dao

import (
	"justus/models"
)

// GetPictureStatisticsKeyList 根据ID获取统计信息
func GetTopicStatisticsKeyList(topicIds []int) (map[int]models.TopicStatistics, error) {
	TopicStatistics := models.TopicStatistics{}
	return TopicStatistics.GetStatisticsKeyList(topicIds)
}