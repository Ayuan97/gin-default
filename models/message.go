package models

import (
	"gorm.io/gorm"
	"justus/global"
	"justus/pkg/util"
)

type Message struct {
	Model
	Uid int `json:"uid"`
	OperateUid int `json:"operate_uid"`
	MessageKey string `json:"message_key"`
	ContentId int `json:"content_id"`
	Picture string `json:"picture"`
	Content string `json:"content"`
}



//入库
func (m *Message)Create() (error) {
	if err := db.Create(m).Error; err != nil {
		return err
	}
	return  nil
}

//消息列表
func (m *Message)GetList(uid int,pageNum int,pageSize int) ([]*Message,error) {
	var message []*Message
	err := db.Select("id,uid,operate_uid,picture,message_key,content,content_id,picture,created_at").Where("uid=?",uid).Order("id desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&message).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.Logger.Error("message GetList:",err)
		return message, err
	}
	for k,v := range message{
		message[k].Picture = util.GetImageUrl(v.Picture)
	}
	return message,nil
}

//清空用户消息
func (m *Message)Clear()error{
	if err := db.Where("uid=?",m.Uid).Delete(m).Error; err != nil {
		global.Logger.Error("message clear:uid:",m.Uid,",error:",err)
		return err
	}
	return  nil
}


//删除
func (m *Message)Delete()error{
	if err := db.Where("uid=? AND id=?",m.Uid,m.ID).Delete(m).Error; err != nil {
		global.Logger.Error("message delete:uid:",m.Uid,",error:",err)
		return err
	}
	return  nil
}


