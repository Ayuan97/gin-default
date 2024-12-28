package topic_service

import (
	"justus/dao"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
)

func CollectTopic(uid int, topicId int, status int) bool {
	key := rediskey.GetUserCollectTopicKey(uid)
	if status == 1 { //收藏
		err := dao.CollectTopic(uid, topicId)
		if err != nil {
			return false
		}
		dao.DelFollowListMap(uid)
		_, _ = gredis.SAdd(key, topicId)
		return true
	} else { //取消收藏
		err := dao.UnCollectTopic(uid, topicId)
		if err != nil {
			return false
		}
		dao.DelFollowListMap(uid)
		_, _ = gredis.SRem(key, topicId)
		return true
	}
}

func GetCollectTopicStatus(uid int, topicId int) int {
	key := rediskey.GetUserCollectTopicKey(uid)
	if gredis.SIsMember(key, topicId) {
		return 1
	} else {
		return 0
	}
}
