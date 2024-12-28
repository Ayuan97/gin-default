package message_service

import (
	"justus/dao"
	"justus/models"
	"justus/pkg/glange"
	"justus/pkg/gredis"
	"justus/pkg/rediskey"
	"justus/pkg/setting"
	"time"
)

type List struct {
	List  []Item `json:"list"`
}

type Item struct {
	Id int `json:"id"`
	Uid int `json:"uid"`
	OperateUid int `json:"operate_uid"`
	MessageKey string `json:"message_key"`
	ContentId int `json:"content_id"`
	Picture string `json:"picture"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	User user `json:"user"`
}

type user struct {
	Uid int `json:"uid"`
	Avatar          string `json:"avatar"`
	FirstName       string `json:"first_name"`
}

//获取消息列表
func GetMessageList(uid int,lange string,pageNum int) (*List,error) {
	var result []Item
	var resultList List
	resultList.List = make([]Item,0)
	message := models.Message{}
	list, err := message.GetList(uid,pageNum, setting.AppSetting.PageSize)
	if err != nil {
		return &resultList, err
	}
	if len(list) > 0{
		for _,v := range list{
			createAt := time.Unix(int64(v.CreatedAt),0).Format("01 02,2006")
			if v.MessageKey == MessageKeyFollowUser{
				content := glange.GetlangeMessage(lange,v.Content)
				result = append(result,Item{v.ID,v.Uid,v.OperateUid,v.MessageKey,v.ContentId,v.Picture,content,createAt,user{}})
			}else{
				result = append(result,Item{v.ID,v.Uid,v.OperateUid,v.MessageKey,v.ContentId,v.Picture,v.Content,createAt,user{}})
			}

		}
		var operateUidList []int
		for _,v := range result{
			operateUidList = append(operateUidList,v.OperateUid)
		}
		userKeyList, err := dao.GetUserInfoUidKey(operateUidList)
		if err != nil {
			return nil, err
		}
		for k,v := range result{
			if _,ok := userKeyList[v.OperateUid];ok{
				result[k].User.Uid = userKeyList[v.OperateUid].(models.User).Uid
				result[k].User.FirstName = userKeyList[v.OperateUid].(models.User).FirstName
				result[k].User.Avatar = userKeyList[v.OperateUid].(models.User).Avatar
			}
		}
		resultList.List = result
	}
	_, _ = gredis.Del(rediskey.GetMessageReadStatusKey(uid))
	return &resultList, nil

}

//清除消息
func ClearList(uid int)bool{
	message := models.Message{}
	message.Uid = uid
	err := message.Clear()
	if err != nil {
		return false
	}
	return true
}

//删除消息
func Delete(uid int,id int)bool{
	message := models.Message{}
	message.Uid = uid
	message.ID = id
	err := message.Delete()
	if err != nil {
		return false
	}
	return true
}