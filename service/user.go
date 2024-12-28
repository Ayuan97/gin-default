package service

import (
	"justus/dao"
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"justus/pkg/setting"
	"justus/service/picture_library_service"
	"justus/service/user_service"
)

type UserCenter struct {
	Page      int `json:"page"`                   //页码 大于0
	FollowUid int `json:"follow_uid" binding:"-"` //被关注用户id
}
type HomePage struct {
	UserInfo   *models.UserInfo `json:"user_info"`
	IsFollow   int              `json:"is_follow"`
	IsFriend   int              `json:"is_friend"`
	MessageDot int              `json:"message_dot"`
}
type ImgList struct {
	ImgList []*models.PictureLibraryList `json:"list"`
}

// GetUserCenter 获取用户个人主页数据
func GetUserCenter(uid int, followUid int) (*HomePage, error) {
	//查询用户信息
	var userInfo *models.User
	var UserId int
	if followUid > 0 {
		UserId = followUid
	} else {
		UserId = uid
	}
	userInfo, err := dao.GetUserInfo(UserId)

	if err != nil {
		return nil, err
	}

	data := HomePage{}
	data.UserInfo = userInfo.Format()

	//查询用户 跟拍 粉丝 数量
	var userPage *models.UserPage
	userPage, _ = dao.GetUserPage(UserId)
	data.UserInfo.FollowPhotoNum = userPage.FollowPhotoNum
	data.UserInfo.FansNum = userPage.FansNum
	data.UserInfo.HotNum = userPage.HotNum
	//是否是好友
	if uid == UserId {
		data.IsFriend = 1
	} else {
		data.IsFriend = user_service.GetFriendStatus(uid, followUid)
	}
	//是否被关注
	if uid == UserId {
		data.IsFollow = 0
	} else {
		data.IsFollow = user_service.GetFollowUserStatus(uid, UserId)
	}
	//获取小红点
	messageDot := 0
	messageReadstatus := gredis.Get(rediskey.GetMessageReadStatusKey(uid))
	if err != nil {
		return nil, err
	}
	if messageReadstatus == "1" {
		messageDot = 1
	}
	data.MessageDot = messageDot

	return &data, nil
}

// GetOneSelfImgList 获取指定uid 的图片列表
func GetOneSelfImgList(uid int, page int, userId int) (*ImgList, error) {
	var data ImgList

	var offset int
	if page > 0 {
		offset = (page - 1) * setting.AppSetting.PageSize
	} else {
		offset = 0
	}
	//查询用户图片
	imgListData, err := dao.GetPictureLibrary(uid, offset, setting.AppSetting.PageSize)
	if err != nil {
		return nil, err
	}
	imgList, err := picture_library_service.PictureLibraryList(imgListData, userId, 0)
	data.ImgList = imgList
	return &data, nil
}
