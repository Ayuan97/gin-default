package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service/component_service"
	"justus/service/friend_service"
	"justus/service/picture_service"
	"net/http"
)


func GetLatestPicture(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Query("component_id")).MustInt()
	//fmt.Println("get_component_id:",id)
	uid := c.GetInt("userId")
	//fmt.Println("get_userid",uid)
	var friendUids []int
	var curFriendUids []int
	if id > 0{
		friendUids, _ = component_service.GetGroupFriendListId(id)
	}else{
		friendUids,_ = friend_service.GetFriendUids(uid)
	}

	for _,v := range friendUids{
		if v != uid{
			curFriendUids = append(curFriendUids,v)
		}
	}
	var library interface{}
	if len(friendUids) != 0{
		library, _ = picture_service.GetLatestLibrary(uid, curFriendUids)
	}else{
		library = nil
	}

	token := c.GetHeader("Authorization")
	uuid := c.GetHeader("uuid")
	fmt.Println("get_component_id:",id,",get_userid:",uid,",uuid:",uuid,",token:",token,",library:",library)

	appG.Response(http.StatusOK, e.SUCCESS, library)
}


func PostLatestPicture(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct{
		ComponentId string `json:"component_id"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil{
		appG.Response(http.StatusBadRequest, e.ERROR,nil)
	}
	id := com.StrTo(request.ComponentId).MustInt()
	//fmt.Println("post_component_id:",id)
	uid := c.GetInt("userId")
	//fmt.Println("post_userid",uid)
	var friendUids []int
	var curFriendUids []int
	if id > 0{
		friendUids, _ = component_service.GetGroupFriendListId(id)
	}else{
		friendUids,_ = friend_service.GetFriendUids(uid)
	}

	for _,v := range friendUids{
		if v != uid{
			curFriendUids = append(curFriendUids,v)
		}
	}

	var library interface{}
	if len(friendUids) != 0{
		library, _ = picture_service.GetLatestLibrary(uid, curFriendUids)
	}else{
		library = nil
	}
	token := c.GetHeader("Authorization")
	uuid := c.GetHeader("uuid")
	fmt.Println("post_component_id:",id,",post_userid:",uid,",uuid:",uuid,",token:",token,",library:",library)
	appG.Response(http.StatusOK, e.SUCCESS, library)
}
