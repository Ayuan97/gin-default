package dao

import (
	"justus/models"
)

// 批量获取点赞信息
func GetBatchLikeInfo(piDs []int, uid int) []int {
	like := models.Like{
		Uid: uid,
	}
	data := like.GetBatchLikeInfo(piDs)
	var likeIds []int
	for _, v := range data {
		likeIds = append(likeIds, v.Pid)
	}

	return likeIds
}

// 批量获取点赞信息
func GetBatchLikeList(piDs []int, uid int) map[int][]*models.LikeFormated {
	like := models.Like{
		Uid: uid,
	}
	result, _ := like.GetLikeList(piDs)
	var uiDs []int
	for _, v := range result {
		uiDs = append(uiDs, v.Uid)
	}
	users, _ := GetUserInfoUidKey(uiDs)

	pidLikeKey := make(map[int][]*models.LikeFormated)
	for _, v := range result {
		info := v.Format()
		if _, ok := users[v.Uid]; ok {
			info.FirstName = users[v.Uid].(models.User).FirstName
			info.LastName = users[v.Uid].(models.User).LastName
			info.Avatar = users[v.Uid].(models.User).Avatar
		}
		var arr []*models.LikeFormated
		arr = append(arr, info)
		pidLikeKey[v.Pid] = arr
	}
	return pidLikeKey
}

// 批量获取好友的点赞信息
func GetFriendLikeInfo(piDs []int, uid int, friendIDs []int) map[int][]*models.LikeFormated {
	like := models.Like{
		Uid: uid,
	}
	result, _ := like.GetFriendLikeList(piDs, friendIDs)
	var uiDs []int
	for _, v := range result {
		uiDs = append(uiDs, v.Uid)
	}
	users, _ := GetUserInfoUidKey(uiDs)

	pidLikeKey := make(map[int][]*models.LikeFormated)
	for _, v := range result {
		info := v.Format()
		if _, ok := users[v.Uid]; ok {
			info.FirstName = users[v.Uid].(models.User).FirstName
			info.LastName = users[v.Uid].(models.User).LastName
			info.Avatar = users[v.Uid].(models.User).Avatar
		}
		var arr []*models.LikeFormated
		arr = append(arr, info)
		pidLikeKey[v.Pid] = arr
	}
	return pidLikeKey
}
