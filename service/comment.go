package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"justus/dao"
	"justus/models"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"justus/pkg/setting"
	"justus/service/message_service"
	"justus/service/push_service"
	"justus/service/statistics_service"
	"strconv"
	"time"
)

type PostComment struct {
	Pid     int    `json:"pid"`
	Content string `json:"content"`
	Page    int    `json:"page"`
}
type PostDelComment struct {
	ID int `json:"id"`
}
type List struct {
	List interface{} `json:"list"`
}

// CreateComment  创建评论
func CreateComment(comment PostComment, uid int) (*PostDelComment, error) {
	var commentInfo *models.Comment
	commentInfo, _ = dao.CrateComment(comment.Pid, uid, comment.Content)
	if commentInfo != nil && commentInfo.ID > 0 {
		key := rediskey.GetPictureCommentMapKey(comment.Pid)
		var z redis.Z
		z.Score = float64(time.Now().Unix())
		z.Member = strconv.Itoa(commentInfo.ID)
		gredis.Zadd(key, z)
		//协程 添加消息
		go func() {
			message_service.AddPictureCommentMessage(uid, commentInfo.PId, comment.Content)
			statistics_service.UploadStatisticsById(comment.Pid, statistics_service.FromPicture, statistics_service.Comment)
			push_service.SendPushInList(push_service.PushSendTypePictureComment, uid, push_service.PushMessageParam{PictureId: commentInfo.PId})
		}()
	}
	return &PostDelComment{ID: commentInfo.ID}, nil
}

// DelComment 删除评论
func DelComment(postDelComment PostDelComment, uid int) (bool, error) {
	//查询评论信息 验证是否是本人评论
	commentInfo, err := dao.GetComment(postDelComment.ID)
	//查询图片信息 验证是否是自己的图片
	ImgInfo, err := models.GetPictureLibrary(commentInfo.PId)
	if err != nil {
		return false, err
	}
	if ImgInfo.Uid == uid || commentInfo.Uid == uid {
		_, error := dao.DelComment(postDelComment.ID, commentInfo.PId, commentInfo.Uid)
		if error != nil {
			return false, error
		}
		key := rediskey.GetPictureCommentMapKey(commentInfo.PId)
		gredis.Zrem(key, strconv.Itoa(postDelComment.ID))

		//协程 添加消息
		go func() {
			statistics_service.UploadStatisticsDecById(commentInfo.PId, statistics_service.FromPicture, statistics_service.Comment)
		}()
		return true, nil
	} else {
		fmt.Println("不能删除该评论,postDelComment-id", postDelComment.ID)
		return false, err
	}

}

// GetComment 获取评论
func GetComment(pid int, page int, userId int) (*List, error) {
	var List List
	var offset int
	if page > 0 {
		offset = (page - 1) * setting.AppSetting.PageSize
	} else {
		offset = 0
	}
	ImgInfo, err := models.GetPictureLibrary(pid)
	if err != nil {
		return &List, err
	}
	var isDel int
	if ImgInfo.Uid == userId {
		isDel = 1
	} else {
		isDel = 0
	}
	comments, err := dao.GetCommentList(pid, offset, setting.AppSetting.PageSize)
	if err != nil {
		return &List, err
	}
	var userIDs []int
	for _, comment := range comments {
		userIDs = append(userIDs, comment.Uid)
	}
	users, err := dao.GetUserInfoUidKey(userIDs)
	if err != nil {
		return &List, err
	}
	commentFormateds := []*models.CommentFormated{}
	for _, comment := range comments {
		commentFormated := &models.CommentFormated{}
		commentFormated.FirstName = users[comment.Uid].(models.User).FirstName
		commentFormated.LastName = users[comment.Uid].(models.User).LastName
		commentFormated.Avatar = users[comment.Uid].(models.User).Avatar
		commentFormated.Content = comment.Content
		commentFormated.Uid = comment.Uid
		commentFormated.Pid = comment.PId
		commentFormated.IsTop = comment.IsTop
		if isDel == 0 {
			if userId == comment.Uid {
				commentFormated.IsDel = 1
			} else {
				commentFormated.IsDel = 0
			}
		} else {
			commentFormated.IsDel = isDel
		}

		commentFormated.ID = comment.ID
		commentFormated.CreatedAt = comment.Model.CreatedAt
		commentFormateds = append(commentFormateds, commentFormated)
	}
	List.List = commentFormateds
	return &List, nil
}
