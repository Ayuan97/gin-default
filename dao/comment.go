package dao

import (
	"fmt"
	"justus/models"
	"strconv"
	"strings"
)

// CrateComment  创建评论
func CrateComment(pid int, uid int, content string) (*models.Comment, error) {
	comment := models.Comment{
		PId:     pid,
		Uid:     uid,
		Content: content,
	}
	return comment.CreatComment()
}

// DelComment 删除评论
func DelComment(id int, pid int, uid int) (bool, error) {
	comment := models.Comment{
		Model: models.Model{
			ID: id,
		},
		PId: pid,
		Uid: uid,
	}
	return comment.DelComment()
}

// GetComment 获取评论内容
func GetComment(id int) (*models.Comment, error) {
	comment := models.Comment{
		Model: models.Model{
			ID: id,
		},
	}
	return comment.GetComment()
}

// GetCommentList 获取评论列表
func GetCommentList(pid int, page int, pageSize int) ([]*models.Comment, error) {
	comment := models.Comment{
		PId: pid,
	}
	return comment.GetCommentList(page, pageSize)
}

// GetCommentFirst 批量获取评论列表第一条
func GetCommentFirst(pid []int, userId int) map[int][]*models.CommentFormated {
	comment := models.Comment{}
	data, err := comment.GetCommentFirst(pid)
	if err != nil {
		fmt.Println("GetCommentFirst error:", err)
	}
	var piDs []int
	for _, v := range data {
		//字符串分割数组
		pidArr := strings.Split(v.Id, ",")
		//字符串转int
		id := make([]int, len(pidArr))
		for i, value := range pidArr {
			id[i], _ = strconv.Atoi(value)
		}
		piDs = append(piDs, id[1])
	}
	var uiDs []int
	result, _ := comment.GetCommentListByPId(piDs)
	for _, v := range result {
		uiDs = append(uiDs, v.Uid)
	}
	users, _ := GetUserInfoUidKey(uiDs)

	pidCommentKey := make(map[int][]*models.CommentFormated)
	for _, v := range result {
		info := v.Format()
		if info.Uid == userId {
			info.IsDel = 1
		} else {
			info.IsDel = 0
		}
		if _, ok := users[v.Uid]; ok {
			info.FirstName = users[v.Uid].(models.User).FirstName
			info.LastName = users[v.Uid].(models.User).LastName
			info.Avatar = users[v.Uid].(models.User).Avatar
		}
		var arr []*models.CommentFormated
		arr = append(arr, info)
		pidCommentKey[v.PId] = arr
	}
	return pidCommentKey

}

// GetFriendCommentFirst 批量获取评论列表好友第一条
func GetFriendCommentFirst(pid []int, userId int, friendIDs []int) map[int][]*models.CommentFormated {
	fmt.Println("friendIDs:", friendIDs)
	comment := models.Comment{}
	data, err := comment.GetCommentFriendFirst(pid, friendIDs)
	if err != nil {
		fmt.Println("GetCommentFriendFirst error:", err)
	}
	var piDs []int
	for _, v := range data {
		//字符串分割数组
		pidArr := strings.Split(v.Id, ",")
		//字符串转int
		id := make([]int, len(pidArr))
		for i, value := range pidArr {
			id[i], _ = strconv.Atoi(value)
		}
		piDs = append(piDs, id[1])
	}
	var uiDs []int
	result, _ := comment.GetCommentListByPId(piDs)
	for _, v := range result {
		uiDs = append(uiDs, v.Uid)
	}
	users, _ := GetUserInfoUidKey(uiDs)

	pidCommentKey := make(map[int][]*models.CommentFormated)
	for _, v := range result {
		info := v.Format()
		if info.Uid == userId {
			info.IsDel = 1
		} else {
			info.IsDel = 0
		}
		if _, ok := users[v.Uid]; ok {
			info.FirstName = users[v.Uid].(models.User).FirstName
			info.LastName = users[v.Uid].(models.User).LastName
			info.Avatar = users[v.Uid].(models.User).Avatar
		}
		var arr []*models.CommentFormated
		arr = append(arr, info)
		pidCommentKey[v.PId] = arr
	}
	return pidCommentKey

}

// CheckCommentIsSelf 检查评论是否是自己的
func CheckCommentIsSelf(comment []*models.CommentFormated, uid int) []*models.CommentFormated {
	for i, v := range comment {
		if v.Uid == uid {
			comment[i].IsDel = 1
		}
	}
	return comment
}
