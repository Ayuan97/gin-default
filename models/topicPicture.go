package models

import (
	"fmt"
	"gorm.io/gorm"
	"justus/global"
)

type TopicPicture struct {
	Model
	TopicId   int     `json:"topic_id"`
	PId       int     `gorm:"p_id" json:"p_id"`
	TopicName string  `json:"topic_name"`
	Weight    float64 `json:"weight" gorm:"column:weight"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}
type TopicPictureFormated struct {
	TopicId   int    `json:"topic_id"`
	TopicName string `json:"topic_name"`
}

// 通过图片ID 查询话题信息
func (m *TopicPicture) GetTopicInfoByPictureId(pids []int) map[int][]TopicPicture {
	var topicPicture []TopicPicture
	result := make(map[int][]TopicPicture)
	if len(pids) > 0 {
		err := db.Where("p_id in (?)", pids).Order("created_at asc").Find(&topicPicture).Error
		if err != nil {
			fmt.Println("GetPictureTopicAll err:", err)
		}
		for _, v := range topicPicture {
			if _, ok := result[v.PId]; ok {
				v.TopicName = "#" + v.TopicName
				result[v.PId] = append(result[v.PId], v)
			} else {
				v.TopicName = "#" + v.TopicName
				result[v.PId] = append(result[v.PId], v)
			}
		}
		return result
	}
	return result
}

// 通过话题ID
func (m *TopicPicture) GetPublicPidByTopicId(offset int, pageSize int) []int {
	var pid []int
	var topicPicture []TopicPicture
	err := db.Select("p_id").Where("topic_id = ? AND is_visible=1", m.TopicId).Order("weight desc,id desc").Offset(offset).Limit(pageSize).Find(&topicPicture).Error
	if err != nil {
		fmt.Println("GetPublicPidByTopicId err:", err)
	}
	for _, v := range topicPicture {
		pid = append(pid, v.PId)
	}
	return pid

}

// 通过图片ID 获取话题ID集合
func (m *TopicPicture) GetTopicIdsByPId() []int {
	var topicId []int
	var topicPicture []TopicPicture
	err := db.Select("topic_id").Where("p_id = ?", m.PId).Find(&topicPicture).Error
	if err != nil {
		fmt.Println("GetPublicPidByTopicId err:", err)
	}
	for _, v := range topicPicture {
		topicId = append(topicId, v.TopicId)
	}
	return topicId

}

//更新权重
func (p *TopicPicture) UpdatePictureWeight() error {
	if err := db.Model(&TopicPicture{}).Where("p_id = ? ", p.PId).Update("weight", p.Weight).Error; err != nil {
		global.Logger.Error("UpdatePictureWeight_topic" + err.Error())
		return err
	}
	return nil
}
