package push_service

import (
	"encoding/json"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
)
const PushSendTypePictureLike = "picture_like"         //照片被点赞
const PushSendTypePictureComment = "picture_comment"         //照片被评论
const PushSendTypeFollowUser = "follow_user"        //我关注了你

type pushMessage struct {
	SendType string `json:"send_type"`
	UId int `json:"uid"`
	MessageParam PushMessageParam `json:"message_param"`
}

type PushMessageParam struct {
	PictureId int `json:"picture_id"`		//操作的图片ID
	Name string `json:"name"`				//用户昵称 ，一般推送内容中需要使用
	Uid int `json:"uid"`				//推送发送人的用户ID
}
//评论推送
func SendPushInList(sendType string,uid int,param PushMessageParam) bool {
	message := pushMessage{
		sendType,uid,param,
	}
	marshal, err := json.Marshal(message)
	if err != nil {
		return false
	}
	_, _ = gredis.LPush(rediskey.PushListQueue, marshal)

	return true
}