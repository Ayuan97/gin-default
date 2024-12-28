package sign

import (
	"encoding/json"
	"fmt"
	"gin-default/pkg/e"
	"gin-default/pkg/gredis"
	"gin-default/pkg/rediskey"
	"gin-default/pkg/util"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
	"strings"
)

//	VerifySignature type Header struct {
//		Sign string `json:"sign"`
//		Version string `json:"version"`
//		Uuid string `json:"uuid"`
//		DeviceType string `json:"DeviceType"`
//		DeviceBrand string `json:"DeviceBrand"`
//		DeviceVersion string `json:"DeviceVersion"`
//		Lange string `json:"lange"`
//		TimeZone string `json:"timeZone"`
//		Authorization string `json:"Authorization"`
//	}
func VerifySignature() gin.HandlerFunc {

	return func(c *gin.Context) {
		//是否是白名单
		whiteList := GetWhiteList()
		url := c.Request.URL.Path
		method := c.Request.Method
		if _, ok := whiteList[url]; ok {
			if whiteList[url] == method {
				c.Next()
				return
			}
		}
		//设备号是否是 白名单
		uuid := c.Request.Header.Get("uuid")
		if GetUUidWhiteList(uuid) {
			c.Next()
			return
		}
		var code int
		var msg string
		data := map[string]interface{}{}
		//code = e.SUCCESS
		//获取所有header参数
		headers := c.Request.Header
		//获取sign参数
		sign := headers.Get("sign")
		if sign == "" {
			code = e.SIGN_ERROR
			msg = e.GetMsg(code)
		}
		var headerMap = map[string]interface{}{}
		//headerMap["sign"] = header.Sign
		headerMap["version"] = headers.Get("version")
		headerMap["uuid"] = headers.Get("uuid")
		headerMap["deviceType"] = headers.Get("deviceType")
		headerMap["deviceBrand"] = headers.Get("deviceBrand")
		headerMap["deviceVersion"] = headers.Get("deviceVersion")
		headerMap["lange"] = headers.Get("lange")
		headerMap["timeZone"] = headers.Get("timeZone")
		headerMap["Authorization"] = headers.Get("Authorization")
		//检测header参数是否为空 为空则返回错误
		for k, v := range headerMap {
			if v == "" {
				code = e.SIGN_ERROR
				msg = e.GetMsg(code) + ":" + k
				return
			}
		}
		//获取所有body参数
		body := c.Request.Body
		bodyBytes, err := ioutil.ReadAll(body)
		c.Request.Body = ioutil.NopCloser(strings.NewReader(string(bodyBytes)))
		if err != nil {
			fmt.Println(err)
		}
		var bodyMap map[string]interface{}
		err = json.Unmarshal(bodyBytes, &bodyMap)
		if err != nil {
			code = e.ERROR
			msg = e.GetMsg(code) + ":" + err.Error()
		}
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
		//根据键名dui map排序
		_, sortedMap := util.SortMapByKey(params)
		var signStr string
		for _, v := range sortedMap {
			if v == "" {
				continue
			}
			//fmt.Println("k:",v,"v:",params[v],"type:",reflect.TypeOf(v))
			switch params[v].(type) {
			case int: //int 转 string
				params[v] = fmt.Sprintf("%d", params[v])
			case float64: //float64 转 string 保留小数
				params[v] = strconv.FormatFloat(params[v].(float64), 'f', -1, 64)
			default:
				//其他类型转string
				params[v] = fmt.Sprintf("%s", params[v])
			}
			signStr += v + "=" + params[v].(string) + "&"
		}
		//截取最后一个&
		signStr = signStr[:len(signStr)-1]
		//md5加密
		mySign := util.Md5(util.Md5(signStr) + "5092wxwx1l05yuma")

		if mySign != sign {

			fmt.Println("签名不正确")
			fmt.Println("mySign:", mySign)
			fmt.Println("appSign:", sign)
			code = e.SIGN_ERROR
			msg = e.GetMsg(code) + " : " + "sign"
			c.JSON(200, gin.H{
				"code": code,
				"msg":  msg,
				"data": data,
			})
			c.Abort()
			return
		}

		c.Next()
	}

}

// GetWhiteList 获取白名单信息
func GetWhiteList() map[string]interface{} {
	var whiteList = map[string]interface{}{}
	whiteList["/api/v2/test"] = "POST"
	return whiteList
}

// GetUUidWhiteList getUUidWhiteList 获取uuid白名单
func GetUUidWhiteList(uuid string) bool {
	key := rediskey.GetWhiteListKey()
	value := gredis.Get(key)
	if value != "" {
		return false
	}
	var whiteList []string
	err := json.Unmarshal([]byte(value), &whiteList)
	if err != nil {
		return false
	}
	for _, v := range whiteList {
		if v == uuid {
			return true
		}
	}

	return false
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
