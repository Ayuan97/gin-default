package component_service

import (
	"justus/models"
	"strconv"
	"strings"
)

func GetGroupFriendListId(uid int) ([]int,error) {
	var friendUids []int
	component, err := models.GetComponent(uid)
	if err != nil {
		return friendUids,err
	}
	friendUidArr := strings.Split(component.GroupFriendListId,`,`)
	for _ ,value := range friendUidArr {
		uid, _ := strconv.Atoi(value)
		friendUids = append(friendUids,uid)
	}
	return friendUids,nil
}