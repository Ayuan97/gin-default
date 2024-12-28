package models

import (
	"gorm.io/gorm"
)

type SmallComponentList struct {
	Model
	GroupFriendListId string `json:"group_friend_list_id"`
}

// GetArticle Get a single article based on ID
func GetComponent(id int) (*SmallComponentList, error) {
	var component SmallComponentList
	err := db.Where("id = ?", id).First(&component).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &component, nil
}
