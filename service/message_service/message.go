package message_service

import (
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"time"
)

//获取图片
func addMessagePicture(message *models.Message,pictureId int) *models.Message {
	picture, err := models.GetPictureLibrary(pictureId)
	if err == nil{
		message.Picture = picture.ImgUrl
		message.Uid = picture.Uid
		return message
	}
	return message
}

func baseMessage(operateUid int,messageKey string)*models.Message  {
	var messageModel models.Message
	messageModel.MessageKey = messageKey
	messageModel.OperateUid = operateUid
	return &messageModel
}

//添加点赞消息
func AddPictureLikeMessage(operateUid int,pictureId int)bool{
	 messageModel := baseMessage(operateUid,MessageKeyPictureLike)
	 messageModel.ContentId = pictureId
	 messageModel = addMessagePicture(messageModel,pictureId)
	 if messageModel.Uid != messageModel.OperateUid{
		 return createMessage(messageModel)
	 }
	 return true
}

//评论
func AddPictureCommentMessage(operateUid int,pictureId int,commentContent string)bool{
	messageModel := baseMessage(operateUid,MessageKeyPictureComment)
	messageModel.ContentId = pictureId
	messageModel.Content = commentContent
	messageModel = addMessagePicture(messageModel,pictureId)
	if messageModel.Uid != messageModel.OperateUid{
		return createMessage(messageModel)
	}
	return true
}

//关注
func AddFollowUserMessage(operateUid int,uid int)bool{
	messageModel := baseMessage(operateUid,MessageKeyFollowUser)
	messageModel.Uid = uid
	messageModel.ContentId = operateUid
	messageModel.Content = "follow_user"
	if messageModel.Uid != messageModel.OperateUid{
		return createMessage(messageModel)
	}
	return true
}

func createMessage(messageModel *models.Message) bool{
	err := messageModel.Create()
	if err != nil {
		return false
	}
	_ = gredis.Set(rediskey.GetMessageReadStatusKey(messageModel.Uid), 1, time.Hour*24*7)
	return true
}
