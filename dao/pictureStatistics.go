package dao

import (
	"justus/models"
)

// GetPictureStatisticsKeyList 根据ID获取统计信息
func GetPictureStatisticsKeyList(pIds []int) (map[int]models.PictureStatistics, error) {
	PictureStatistics := models.PictureStatistics{}
	return PictureStatistics.GetStatisticsKeyList(pIds)
}