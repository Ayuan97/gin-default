package user_service

import (
	"justus/dao"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"justus/service/message_service"
	"justus/service/push_service"
)

func FollowUser(uid int, followUid int, status int) bool {
	key := rediskey.GetUserFollowUserKey(uid)
	if status == 1 { //关注
		err := dao.FollowUser(uid, followUid)
		if err != nil {
			return false
		}
		_, _ = gredis.SAdd(key, followUid)
		dao.DelFollowListMap(uid)
		go func() { //协程 添加消息
			message_service.AddFollowUserMessage(uid, followUid)
			userInfo, _ := dao.GetUserInfo(uid)
			if userInfo.Uid > 0 {
				push_service.SendPushInList(push_service.PushSendTypeFollowUser, followUid, push_service.PushMessageParam{Name: userInfo.FirstName, Uid: uid})
			}
		}()
		return true
	} else { //取消关注
		err := dao.UnFollowUser(uid, followUid)
		if err != nil {
			return false
		}
		dao.DelFollowListMap(uid)
		_, _ = gredis.SRem(key, followUid)
		return true
	}
}

func GetFollowUserStatus(uid int, followUid int) int {
	key := rediskey.GetUserFollowUserKey(uid)
	if gredis.SIsMember(key, followUid) {
		return 1
	} else {
		return 0
	}
}

func GetFriendStatus(uid int, FriendUid int) int {
	key := rediskey.GetFriendListMap(uid)
	if gredis.SIsMember(key, FriendUid) {
		return 1
	} else {
		return 0
	}
}
