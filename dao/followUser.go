package dao

import (
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
)

func GetBatchFollowInfo(followUid []int, uid int) map[int][]*models.FollowUser {
	userFollow := models.FollowUser{
		Uid: uid,
	}
	data := userFollow.IsFollow(followUid)
	result := make(map[int][]*models.FollowUser)
	for _, v := range data {
		result[v.Uid] = append(result[v.FollowUid], v)
	}
	return result
}

//DelFollowListMap 删除关注列表map
func DelFollowListMap(uid int) int64 {
	key := rediskey.GetFollowListMapKey(uid)
	value, _ := gredis.Del(key)
	return value
}
