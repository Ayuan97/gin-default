package api

import (
	"github.com/gin-gonic/gin"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service/message_service"
	"net/http"
)

//消息列表
func MessageList(c *gin.Context) {
	appG := app.Gin{C: c}
	type Request struct{
		PageNum int `json:"page_num"`
	}
	uid := c.GetInt("userId")
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil{
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS,nil)
		return
	}
	lange := c.GetString("lange")
	list, err := message_service.GetMessageList(uid, lange,request.PageNum)
	if err != nil {
		appG.Response(e.ERROR, e.ERROR,nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, list)
}


//清空消息
func MessageClear(c *gin.Context){
	appG := app.Gin{C: c}
	uid := c.GetInt("userId")
	res := message_service.ClearList(uid)
	if !res {
		appG.Response(e.ERROR, e.ERROR,nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

//删除消息
func MessageDelete(c *gin.Context){
	appG := app.Gin{C: c}
	uid := c.GetInt("userId")
	type Request struct{
		Id int `json:"id"`
	}
	request := Request{}
	err := c.ShouldBindJSON(&request)
	if err != nil{
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS,nil)
		return
	}

	res := message_service.Delete(uid,request.Id)
	if !res {
		appG.Response(e.ERROR, e.ERROR,nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}