package friend_service

import (
	"justus/models"
)


func GetFriendUids(uid int) ([]int,error) {
	var friend_uid []int
	list, err := models.GetFriendList(uid)
	if err != nil {
		return friend_uid,err
	}
	for _,value := range list{
		friend_uid = append(friend_uid,value.FriendUid)
	}

	return friend_uid,nil
}