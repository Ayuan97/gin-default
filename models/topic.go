package models

import (
	"fmt"
	"gorm.io/gorm"
	"justus/global"
	"justus/pkg/util"
	"strings"
)

type Topics struct {
	Model
	Uid          int     `json:"uid"`
	TopicName    string  `json:"topic_name"`
	TopicPicture string  `json:"topic_picture"`
	Lange        string  `json:"lange"`
	CollectNum   int     `json:"collect_num"`
	Weight       float64 `json:"weight"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

//收藏主题
func (m *Topics) GetOne() (*Topics, error) {
	var topic *Topics
	err := db.Where("id = ?", m.ID).First(&topic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("err", err)
		return nil, err
	}
	topic = commonHandle(topic)
	return topic, nil
}

//收藏主题通过名称查找
func (m *Topics) GetOneName() (*Topics, error) {
	var topic *Topics
	err := db.Where("topic_name = ?", m.TopicName).First(&topic).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	topic = commonHandle(topic)
	return topic, nil
}

/**
 * 话题最热排序
 */
func (m *Topics) GetHotList(pageNum int, pageSize int, lange string, keyword string) ([]*Topics, error) {
	var topics []*Topics
	var result []*Topics
	fields := []string{"lange = ?"}
	values := []interface{}{lange}
	if keyword != "" {
		fields = append(fields, "topic_name like ?")
		values = append(values, keyword+"%")
	}

	err := db.Select("id,topic_name").Where(strings.Join(fields, " AND "), values...).Order("weight desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topics).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(topics) > 0 {
		for _, v := range topics {
			result = append(result, commonHandle(v))
		}
	}
	return result, nil
}

/**
 * 话题最热排序 不带#
 */
func (m *Topics) GethotlistV2(pageNum int, pageSize int, lange string, keyword string) ([]*Topics, error) {
	var topics []*Topics
	fields := []string{"lange = ?"}
	values := []interface{}{lange}
	if keyword != "" {
		fields = append(fields, "topic_name like ?")
		values = append(values, keyword+"%")
	}

	err := db.Select("id,topic_name").Where(strings.Join(fields, " AND "), values...).Order("weight desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&topics).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return topics, nil
}

//自增收藏数量
func (m *Topics) IncTopicCollectNum(tx *gorm.DB) error {
	if err := tx.Model(Topics{}).Where("id = ? ", m.ID).UpdateColumn("collect_num", gorm.Expr("collect_num + ?", 1)).Error; err != nil {
		global.Logger.Error("IncTopicCollectNum" + err.Error())
		return err
	}
	return nil
}

//减收藏数量
func (m *Topics) DecTopicCollectNum(tx *gorm.DB) error {
	if err := tx.Model(Topics{}).Where("id = ? ", m.ID).UpdateColumn("collect_num", gorm.Expr("collect_num - ?", 1)).Error; err != nil {
		global.Logger.Error("DecTopicCollectNum" + err.Error())
		return err
	}
	return nil
}

//更新权重
func (m *Topics) UpdateTopicWeight() error {
	if err := db.Model(Topics{}).Where("id = ? ", m.ID).Update("weight", m.Weight).Error; err != nil {
		global.Logger.Error("updateWeight" + err.Error())
		return err
	}
	return nil
}

//公共处理
func commonHandle(topic *Topics) *Topics {
	topic.TopicName = "#" + topic.TopicName
	topic.TopicPicture = util.GetImageUrl(topic.TopicPicture)
	return topic
}
