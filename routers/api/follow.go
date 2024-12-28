package api

import (
	"github.com/gin-gonic/gin"
	"justus/global"
	"justus/pkg/app"
	"justus/pkg/e"
	"justus/service/follow_service"
	"net/http"
)

// FollowList  获取关注列表
func FollowList(c *gin.Context) {
	appG := app.Gin{C: c}
	userId := c.MustGet("userId").(int)
	var param follow_service.PostFollowList
	err := c.ShouldBind(&param)
	if err != nil {
		global.Logger.Error("GetComment BindJSON err: ", err)
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}
	data := follow_service.GetFollowList(userId, param.Page, userId)
	appG.Response(http.StatusOK, e.SUCCESS, data)
	return
}

func Test(c *gin.Context) {
	//global.Logger.WithFields(logrus.Fields{
	//	"userId": c.MustGet("userId").(int),
	//}).Error("错误")
	//global.Logger.Error("error test")
	//global.Logger.Info("info test")
	//global.Logger.Println("print test")

	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, e.SUCCESS, "test")
	return
	//origData := []byte("Hello World") // 待加密的数据
	//encrypted := aes.AesEncryptCBC(origData)
	//str := base64.StdEncoding.EncodeToString(encrypted)
	//fmt.Println("str:", str)
	//
	//decrypted, _ := base64.StdEncoding.DecodeString(str)
	//decrypted = aes.AesDecryptCBC(decrypted)
	//log.Println("解密结果：", string(decrypted))

}
