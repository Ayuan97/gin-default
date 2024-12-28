package rediskey

import "fmt"

const PictureLikeNum = "picture:like:like_num_%d"                   //图片点赞数量
const PictureCommentNum = "picture:comment:comment_num_%d"          //评论点赞数量
const PictureFirstComment = "picture:last_comment:first_comment_%d" //最新评论

const FollowListMapUpdateTime = "follow:list:map:update_time_%d"   //关注列表更新时间
const FollowListMap = "user:follow:list:%d"                        //关注列表
const RecommentListReadRecordMap = "recommend:list:read_record_%d" //推荐列表已读记录
const UserCollectTopicMap = "user:collect:topic:map:%d"            //收藏话题的集合
const UserFollowUserMap = "user:follow:user:map:%d"                //关注用户的集合

const TopicUserPictureUpdateTimeMap = "topic:user:public:updateTime:map" //话题及图片集合的更新时间
const TopicPictureMap = "topic:picture:public:map:%d"                    //话题图片集合
const UserPictureMap = "user:picture:public:map:%d"                      //用户图片集合

const PictureLikeMap = "picture:like:map:%d"       //图片点赞信息集合
const PictureCommentMap = "picture:comment:map:%d" //图片评论集合

const StatisticsTopicList = "statistics:topic:list"     //话题统计队列
const StatisticsPictureList = "statistics:picture:list" //图片统计队列

const StatisticsTopicCurData = "statistics:topic:%d"     //话题实时计算数据
const StatisticsPictureCurData = "statistics:picture:%d" //图片实时计算数据

const PushMessageFriendSwitch = "push:friend:switch:%d"    //好友消息开关
const PushMessageHotBeatSwitch = "push:hot_beat:switch:%d" //好友消息开关

const WebShareLink = "web:share:link_%s"    //邀拍 web链接 秘钥key
const PushListQueue = "push_list"           //推送队列
const MessageReadStatus = "message:read:%d" //消息是否阅读

const UserFriendMap = "user:friend:map:%d" //好友集合

const WhiteList = "yq_white_list" //白名单

// GetFriendListMap 好友列表集合
func GetFriendListMap(uid int) string {
	return fmt.Sprintf(UserFriendMap, uid)
}

// GetMessageReadStatusKey 消息是否阅读
func GetMessageReadStatusKey(uid int) string {
	return fmt.Sprintf(MessageReadStatus, uid)
}

// GetWebShareLinkKey 获取分享link key
func GetWebShareLinkKey(key string) string {
	return fmt.Sprintf(WebShareLink, key)
}

// GetPushMessageFriendSwitchKey 获取好友推送开关key
func GetPushMessageFriendSwitchKey(uid int) string {
	return fmt.Sprintf(PushMessageFriendSwitch, uid)
}

// GetPushMessageHotBeatSwitchKey 获取热拍相关推送 key
func GetPushMessageHotBeatSwitchKey(uid int) string {
	return fmt.Sprintf(PushMessageHotBeatSwitch, uid)
}

// GetRecommendListReadRecordKey 获取推荐列表已读记录key
func GetRecommendListReadRecordKey(uid int) string {
	return fmt.Sprintf(RecommentListReadRecordMap, uid)
}

// GetFollowListMapKey 获取关注列表key
func GetFollowListMapKey(uid int) string {
	return fmt.Sprintf(FollowListMap, uid)
}

// FollowListMapUpdateTime 获取关注列表更新时间
func GetFollowListMapUpdateTime(uid int) string {
	return fmt.Sprintf(FollowListMapUpdateTime, uid)
}

// GetTopicUserPictureUpdateTimeMapKey 获取话题及图片集合的更新时间
func GetTopicUserPictureUpdateTimeMapKey() string {
	return TopicUserPictureUpdateTimeMap
}

// GetTopicPictureMapKey 获取话题图片合集key
func GetTopicPictureMapKey(id int) string {
	return fmt.Sprintf(TopicPictureMap, id)
}

// GetUserPictureMapKey 获取用户图片合集key
func GetUserPictureMapKey(id int) string {
	return fmt.Sprintf(UserPictureMap, id)
}

// GetPictureLikeNumkey 获取图片点赞数量
func GetPictureLikeNumkey(pid int) string {
	return fmt.Sprintf(PictureLikeNum, pid)
}

// GetPictureCommentNumkey 获取图片评论数量
func GetPictureCommentNumkey(pid int) string {
	return fmt.Sprintf(PictureCommentNum, pid)
}

// GetPictureFirstCommentkey 获取图片首条评论
func GetPictureFirstCommentkey(pid int) string {
	return fmt.Sprintf(PictureFirstComment, pid)
}

// GetUserCollectTopicKey 获取用户收藏话题的集合key
func GetUserCollectTopicKey(uid int) string {
	return fmt.Sprintf(UserCollectTopicMap, uid)
}

// GetUserFollowUserKey 关注用户的集合
func GetUserFollowUserKey(uid int) string {
	return fmt.Sprintf(UserFollowUserMap, uid)
}

// GetPictureLikeMapKey 图片点赞合集key
func GetPictureLikeMapKey(pid int) string {
	return fmt.Sprintf(PictureLikeMap, pid)
}

// GetPictureCommentMapKey 图片评论合集key
func GetPictureCommentMapKey(pid int) string {
	return fmt.Sprintf(PictureCommentMap, pid)
}

// GetStatisticsTopicCurDataKey 话题实时曝光总数
func GetStatisticsTopicCurDataKey(topicId int) string {
	return fmt.Sprintf(StatisticsTopicCurData, topicId)
}

// GetStatisticsPictureCurDataKey 图片实时曝光总数
func GetStatisticsPictureCurDataKey(pId int) string {
	return fmt.Sprintf(StatisticsPictureCurData, pId)
}

// GetWhiteListKey 获取白名单key
func GetWhiteListKey() string {
	return WhiteList
}
