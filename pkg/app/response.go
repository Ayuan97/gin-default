package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"justus/pkg/e"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
	return
}

// Success 成功响应
func (g *Gin) Success(data interface{}) {
	g.Response(http.StatusOK, e.SUCCESS, data)
}

// Error 错误响应
func (g *Gin) Error(errCode int) {
	g.Response(http.StatusOK, errCode, nil)
}

// ErrorWithData 带数据的错误响应
func (g *Gin) ErrorWithData(errCode int, data interface{}) {
	g.Response(http.StatusOK, errCode, data)
}

// InvalidParams 参数错误响应
func (g *Gin) InvalidParams() {
	g.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
}

// Unauthorized 未授权响应
func (g *Gin) Unauthorized(errCode int) {
	g.Response(http.StatusUnauthorized, errCode, nil)
}
