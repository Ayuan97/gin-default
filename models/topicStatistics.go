package models

import (
	"fmt"
	"gorm.io/gorm"
	"justus/global"
)

type TopicStatistics struct {
	Model
	TopicId            int    `json:"topic_id"`
	ImpressionNum            int    `json:"impression_num"`
	CurImpressionNum            int    `json:"cur_impression_num"`
	ClickNum            int    `json:"click_num"`
	LikeNum            int    `json:"like_num"`
	CurLikeNum            int    `json:"cur_like_num"`
	CommentNum            int    `json:"comment_num"`
	CurCommentNum            int    `json:"cur_comment_num"`
	ShareNum            int    `json:"share_num"`
	FollowPhotoNum            int    `json:"follow_photo_num"`
	CurFollowPhotoNum            int    `json:"cur_follow_photo_num"`
}

func (m *TopicStatistics)GetImpressionKeyList(topicIds []int) (map[int]int,error) {
	var topicStatistics []TopicStatistics
	result := make(map[int]int)
	if len(topicIds) > 0 {
		err := db.Select("topic_id,impression_num").Where("topic_id in (?)", topicIds).Find(&topicStatistics).Error
		if err != nil {
			fmt.Println("getImpressionKeyList err:", err)
		}
		for _, v := range topicStatistics {
			result[v.TopicId] = v.ImpressionNum
		}
		return result,nil
	}
	return result,nil
}


//自增处理

func (m *TopicStatistics)IncData(column map[string]int) bool {
	if m.TopicId <= 0{
		return false
	}
	tx := db.Model(&TopicStatistics{}).Where("topic_id = ? ", m.TopicId)
	var updateGorm = make(map[string]interface{})
	for k,v := range column{
		switch k {
			case "impression":
				updateGorm["impression_num"] = gorm.Expr("impression_num + ?",v)
				updateGorm["cur_impression_num"] = gorm.Expr("cur_impression_num + ?",v)
			case "click":
				updateGorm["click_num"] = gorm.Expr("click_num + ?",v)
			case "like":
				updateGorm["like_num"] = gorm.Expr("like_num + ?",v)
				updateGorm["cur_like_num"] = gorm.Expr("cur_like_num + ?",v)
			case "comment":
				updateGorm["comment_num"] = gorm.Expr("comment_num + ?",v)
				updateGorm["cur_comment_num"] = gorm.Expr("cur_comment_num + ?",v)
			case "share":
				updateGorm["share_num"] = gorm.Expr("share_num + ?",v)
			case "follow_photo":
				updateGorm["follow_photo_num"] = gorm.Expr("follow_photo_num + ?",v)
				updateGorm["cur_follow_photo_num"] = gorm.Expr("cur_follow_photo_num + ?",v)
			default:
		}
	}
	if len(updateGorm) > 0{
		if err := tx.Updates(updateGorm).Error;err != nil{
			global.Logger.Error("TopicStatisticsIncData"+err.Error())
			return false
		}
	}
	return true
}


//更新权重为初值值
func (m *TopicStatistics)InitWeightData(ids []int) error {
	if len(ids) > 0{
		if err := db.Model(&TopicStatistics{}).Where("topic_id IN ? ", ids).Updates(TopicStatistics{CurImpressionNum: 1,CurLikeNum: 1,CurCommentNum: 1,CurFollowPhotoNum: 1}).Error; err != nil {
			global.Logger.Error("updateWeight"+err.Error())
			return err
		}
	}else{
		var err error
		return err
	}

	return nil
}

func (m *TopicStatistics)GetStatisticsKeyList(pIds []int) (map[int]TopicStatistics,error) {
	var topicStatistics []TopicStatistics
	result := make(map[int]TopicStatistics)
	if len(pIds) > 0 {
		err := db.Where("topic_id in (?)", pIds).Find(&topicStatistics).Error
		if err != nil {
			fmt.Println("topicStatistics,GetStatisticsKeyList err:", err)
		}
		for _, v := range topicStatistics {
			result[v.TopicId] = v
		}
		return result,nil
	}
	return result,nil
}


func (m *TopicStatistics)GetStatisticsItem() (*TopicStatistics,error) {
	var topicStatistics TopicStatistics
	err := db.Where("topic_id=?", m.TopicId).First(&topicStatistics).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("topicStatistics,GetStatisticsItem err:", err)
		return nil, err
	}
	return &topicStatistics,nil
}
