package message_service

const MessageKeyPictureLike = "picture_like"		//照片点赞
const MessageKeyPictureComment = "picture_comment"		//照片评论
const MessageKeyFollowUser = "follow_user"		//关注


//消息处理
type MessageSetting struct {
	Content string  `json:"content"`
}

var MessageSettingList = make(map[string]MessageSetting)

func GetMessageSetting(messageKey string){
	MessageSettingList[MessageKeyPictureLike] = MessageSetting{Content: ""}
}









