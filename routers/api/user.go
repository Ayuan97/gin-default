package api

//func UserCenter(c *gin.Context) {
//	appG := app.Gin{C: c}
//	userId := c.MustGet("userId").(int)
//	var param service.UserCenter
//	_ = c.ShouldBind(&param)
//	//if err != nil {
//	//	global.Logger.Error("UserCenter BindJSON err: ", err)
//	//	appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
//	//	return
//	//}
//	data, err := service.GetUserCenter(userId, param.FollowUid)
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.SUCCESS, err.Error())
//		return
//	}
//	if data == nil {
//		appG.Response(http.StatusOK, e.SUCCESS, make(map[string]interface{}))
//		return
//	}
//	appG.Response(http.StatusOK, e.SUCCESS, data)
//	return
//}
//
//func UserCenterImgList(c *gin.Context) {
//	appG := app.Gin{C: c}
//	userId := c.MustGet("userId").(int)
//	param := service.UserCenter{}
//	_ = c.ShouldBindJSON(&param)
//	data, err := service.GetOneSelfImgList(param.FollowUid, param.Page, userId)
//	if err != nil {
//		appG.Response(http.StatusInternalServerError, e.SUCCESS, err.Error())
//		return
//	}
//	appG.Response(http.StatusOK, e.SUCCESS, data)
//}
//
////用户关注
//func UserFollow(c *gin.Context) {
//	appG := app.Gin{C: c}
//	type Request struct {
//		FollowUid int `json:"follow_uid"`
//		Status    int `json:"status"`
//	}
//	request := Request{}
//	err := c.ShouldBindJSON(&request)
//	if err != nil {
//		appG.Response(http.StatusBadRequest, e.ERROR, nil)
//		return
//	}
//	uid := c.GetInt("userId")
//	res := user_service.FollowUser(uid, request.FollowUid, request.Status)
//	if res {
//		appG.Response(http.StatusOK, e.SUCCESS, nil)
//	} else {
//		appG.Response(e.ERROR, e.ERROR, nil)
//	}
//
//}
