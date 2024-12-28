package dao

import (
	"justus/global"
	"justus/models"
)


// 收藏
func CollectTopic(uid int, topicId int) error {
	tx := models.GetDb().Begin()		//开启事务
	var collectTopic = models.CollectTopic{}
	collectTopic.Uid = uid
	collectTopic.TopicId = topicId
	err := collectTopic.Create(tx)
	if err != nil {
		global.Logger.Error("create collect_topic failed")
		tx.Rollback()
		return err
	}
	var topic = models.Topics{}
	topic.ID = topicId
	err = topic.IncTopicCollectNum(tx)
	if err != nil {
		global.Logger.Error("inc topic collect num failed")
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// 取消收藏
func UnCollectTopic(uid int, topicId int) error {
	tx := models.GetDb().Begin()		//开启事务
	var collectTopic = models.CollectTopic{}
	collectTopic.TopicId = topicId
	collectTopic.Uid = uid
	err := collectTopic.Del(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	var topic = models.Topics{}
	topic.ID = topicId
	err = topic.DecTopicCollectNum(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
