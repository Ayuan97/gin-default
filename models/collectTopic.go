package models

import (
	"fmt"
	"gorm.io/gorm"
)

type CollectTopic struct {
	Model
	Uid int `json:"uid"`
	TopicId int `json:"topic_id"`
}

//收藏主题
func (m *CollectTopic)Create(tx *gorm.DB) error {
	if err := tx.Create(m).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return  nil
}

//删除
func (m *CollectTopic)Del(tx *gorm.DB) error {
	if err := tx.Where("uid=? AND topic_id=?",m.Uid,m.TopicId).Delete(m).Error; err != nil {
		return err
	}
	return  nil
}


