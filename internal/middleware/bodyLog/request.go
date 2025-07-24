package bodyLog

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"justus/internal/global"
	"strings"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		//statusCode := c.Writer.Status()
		//请求内容
		headers := c.Request.Header
		var headerMap = map[string]interface{}{}
		headerMap["sign"] = headers.Get("sign")
		headerMap["version"] = headers.Get("version")
		headerMap["uuid"] = headers.Get("uuid")
		headerMap["deviceType"] = headers.Get("deviceType")
		headerMap["deviceBrand"] = headers.Get("deviceBrand")
		headerMap["deviceVersion"] = headers.Get("deviceVersion")
		headerMap["lange"] = headers.Get("lange")
		headerMap["timeZone"] = headers.Get("timeZone")
		headerMap["Authorization"] = headers.Get("Authorization")
		//获取所有body参数
		body := c.Request.Body
		bodyBytes, err := ioutil.ReadAll(body)
		c.Request.Body = ioutil.NopCloser(strings.NewReader(string(bodyBytes)))
		if err != nil {
			global.Logger.Error("request get body", err.Error())
		}
		bodyMap := map[string]interface{}{}

		_ = json.Unmarshal(bodyBytes, &bodyMap)
		//获取所有query参数
		query := c.Request.URL.Query()
		c.Request.URL.RawQuery = query.Encode()
		queryMap := map[string]interface{}{}
		for k, v := range query {
			queryMap[k] = v[0]
		}
		//获取所有form参数
		form := c.Request.PostForm
		c.Request.PostForm = form
		formMap := map[string]interface{}{}
		for k, v := range form {
			formMap[k] = v[0]
		}
		//合并所有参数
		params := MergeMap(map[string]interface{}{}, queryMap, formMap, bodyMap, headerMap)
		Info := logrus.Fields{}
		Info["response"] = string(blw.body.Bytes())
		//errInfo["Params"] = string(paramsJson)
		Info["status"] = c.Writer.Status()
		Info["ip"] = c.ClientIP()
		Info["method"] = c.Request.Method
		Info["path"] = c.Request.URL.Path
		Info["remoteAddr"] = c.Request.RemoteAddr
		Info["requestURI"] = c.Request.RequestURI
		Info["proto"] = c.Request.Proto
		for k, v := range params {
			Info[k] = v
		}
		//global.Logger.Info("path:", c.Request.URL.String(), "  response:", string(blw.body.Bytes()))
	}

}

func MergeMap(mObj ...map[string]interface{}) map[string]interface{} {
	newObj := map[string]interface{}{}
	for _, m := range mObj {
		for k, v := range m {
			newObj[k] = v
		}
	}
	return newObj
}
