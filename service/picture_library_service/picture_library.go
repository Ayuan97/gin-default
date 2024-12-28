package picture_library_service

import (
	"justus/dao"
	"justus/models"
	"justus/pkg/util"
	"justus/service/user_service"
	"strconv"
	"strings"
)

var PictureLibraryType = map[string]int{
	"照片": 1,
	"视频": 2,
	"合拍": 3,
}
var PictureLibraryTuneType = map[int]int{
	2: 2,
	3: 2,
	4: 3,
	5: 3,
	6: 3,
	7: 4,
}

// PictureLibraryList 图片列表数据处理
func PictureLibraryList(data []*models.PictureLibrary, userId int, Source int) ([]*models.PictureLibraryList, error) {
	var userIDs []int
	var imgIDs []int
	for _, v := range data {
		imgIDs = append(imgIDs, v.ID)
		userIDs = append(userIDs, v.Uid)
	}
	result := make([]*models.PictureLibraryList, 0)
	users, err := dao.GetUserInfoUidKey(util.RemoveRepeatedElement(userIDs))
	if err != nil {
		return result, err
	}
	var comment map[int][]*models.CommentFormated
	var like map[int][]*models.LikeFormated
	if Source == 1 {
		//获取好友信息
		friendUids, _ := models.GetFriendList(userId)
		var friendIDs []int
		//加入自己的uid
		friendIDs = append(friendIDs, userId)
		for _, v := range friendUids {
			friendIDs = append(friendIDs, v.Uid)
		}
		comment = dao.GetFriendCommentFirst(imgIDs, userId, friendIDs) //好友评论
		like = dao.GetFriendLikeInfo(imgIDs, userId, friendIDs)        //好友点赞
	} else {
		comment = dao.GetCommentFirst(imgIDs, userId) //评论
		like = dao.GetBatchLikeList(imgIDs, userId)   //点赞
	}
	topic := dao.GetTopicPidKey(imgIDs)            //话题
	IsLike := dao.GetBatchLikeInfo(imgIDs, userId) //点赞
	statisticslist, _ := dao.GetPictureStatisticsKeyList(imgIDs)

	for _, v := range data {
		likeNum := 0
		commentNum := 0
		impressionNum := 0
		if _, ok := statisticslist[v.ID]; ok {
			likeNum = statisticslist[v.ID].LikeNum
			commentNum = statisticslist[v.ID].CommentNum
			impressionNum = statisticslist[v.ID].ImpressionNum
		}
		//查询用户图片点赞数量 GetLikeNum(c,v.ID)
		//查询用户图片评论数量
		info := v.Format()
		info.HotNum = impressionNum
		info.LikeNum = util.NumTransform(likeNum)
		info.IsLike = util.In(v.ID, IsLike)
		info.FirstName = users[v.Uid].(models.User).FirstName
		info.LastName = users[v.Uid].(models.User).LastName
		info.Avatar = users[v.Uid].(models.User).Avatar
		info.CommentNum = util.NumTransform(commentNum)
		if _, ok := like[v.ID]; ok {
			info.Like = like[v.ID]
		} else {
			info.Like = []*models.LikeFormated{}
		}
		if _, ok := comment[v.ID]; ok {
			info.Comment = comment[v.ID]
		} else {
			info.Comment = []*models.CommentFormated{}
		}
		info.Comment = dao.CheckCommentIsSelf(info.Comment, info.Uid)
		if _, ok := topic[v.ID]; ok {
			info.Topics = topic[v.ID]
		} else {
			info.Topics = []models.TopicPicture{}
		}
		//是否关注该用户
		info.IsFollow = user_service.GetFollowUserStatus(userId, v.Uid)
		//是否是好友
		info.IsFriend = user_service.GetFriendStatus(userId, v.Uid)
		//判断是否是合拍类型  是合拍类型判断合拍是否完成
		if info.Type == PictureLibraryType["合拍"] {
			//分割字符串
			arr := strings.Split(v.Position, ",")
			inTuneType, _ := strconv.Atoi(info.InTuneType)
			if len(arr) == PictureLibraryTuneType[inTuneType] {
				info.IsTuneSuccess = 1
			} else {
				info.IsTuneSuccess = 0
			}
		}
		result = append(result, info)
	}
	return result, nil
}
