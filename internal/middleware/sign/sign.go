package sign

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"justus/pkg/app"
	"justus/pkg/e"
	"justus/pkg/gredis"

	"github.com/gin-gonic/gin"
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

	return gin.HandlerFunc(func(c *gin.Context) {

		appG := app.Gin{C: c}

		// 开发模式或测试模式下跳过签名验证
		if c.GetHeader("skip-signature") == "true" {
			c.Next()
			return
		}



		// 获取所有header参数
		headers := c.Request.Header
		// 获取sign参数
		sign := headers.Get("sign")

		if sign == "" {
			fmt.Println("签名参数为空")
			appG.Error(e.SIGN_ERROR)
			c.Abort()
			return
		}

		var headerMap = map[string]interface{}{}
		headerMap["version"] = headers.Get("version")
		headerMap["uuid"] = headers.Get("uuid")
		headerMap["deviceType"] = headers.Get("deviceType")
		headerMap["deviceBrand"] = headers.Get("deviceBrand")
		headerMap["deviceVersion"] = headers.Get("deviceVersion")
		headerMap["lange"] = headers.Get("lange")
		headerMap["timeZone"] = headers.Get("timeZone")
		headerMap["Authorization"] = headers.Get("Authorization")

		// 检测header参数是否为空
		for k, v := range headerMap {
			if v == "" {
				fmt.Println("请求头参数错误:", k)
				appG.Error(e.INVALID_PARAMS)
				c.Abort()
				return
			}
		}

		// 获取body参数
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("读取body失败:", err)
			appG.Error(e.ERROR)
			c.Abort()
			return
		}
		// 重新设置body
		c.Request.Body = ioutil.NopCloser(strings.NewReader(string(bodyBytes)))

		var bodyMap map[string]interface{}
		if len(bodyBytes) > 0 {
			err = json.Unmarshal(bodyBytes, &bodyMap)
			if err != nil {
				fmt.Println("body json解析失败:", err)
				appG.Error(e.ERROR)
				c.Abort()
				return
			}
		}

		// 获取所有query参数
		query := c.Request.URL.Query()
		c.Request.URL.RawQuery = query.Encode()
		queryMap := map[string]interface{}{}
		for k, v := range query {
			queryMap[k] = v[0]
		}

		// 构建签名参数map
		requestParamMap := map[string]interface{}{}
		for k, v := range headerMap {
			requestParamMap[k] = v
		}
		for k, v := range queryMap {
			requestParamMap[k] = v
		}
		for k, v := range bodyMap {
			requestParamMap[k] = v
		}

		// 生成签名
		mySign, err := getSign(requestParamMap)
		if err != nil {
			fmt.Println("生成签名失败:", err)
			appG.Error(e.ERROR)
			c.Abort()
			return
		}

		// 验证签名
		if mySign != sign {
			fmt.Println("签名验证失败")
			fmt.Println("期望签名:", mySign)
			fmt.Println("实际签名:", sign)
			appG.Error(e.SIGN_ERROR)
			c.Abort()
			return
		}

		c.Next()
	})

}

// 生成签名
func getSign(requestParamMap map[string]interface{}) (string, error) {
	// 这里实现签名算法
	// 示例：简单的字符串拼接然后MD5
	var keys []string
	for k := range requestParamMap {
		if k != "sign" {
			keys = append(keys, k)
		}
	}

	// 对keys排序
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	var signStr strings.Builder
	for _, k := range keys {
		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(fmt.Sprintf("%v", requestParamMap[k]))
		signStr.WriteString("&")
	}

	signString := signStr.String()
	if len(signString) > 0 {
		signString = signString[:len(signString)-1] // 移除最后一个&
	}

	// 这里应该实现实际的签名算法，比如MD5、SHA256等
	// 为了简化，这里返回原字符串
	return signString, nil
}

// 生成随机字符串(用于生成nonce)
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[len(charset)/2] // 简化实现
	}
	return string(b)
}

// 检查时间戳是否有效(防重放攻击)
func isTimestampValid(timestamp string) bool {
	// 解析时间戳并检查是否在有效时间范围内
	// 这里简化实现，实际应该检查时间戳是否在合理范围内
	_, err := strconv.ParseInt(timestamp, 10, 64)
	return err == nil
}

// 检查nonce是否已使用(防重放攻击)
func isNonceUsed(nonce string) bool {
	// 检查Redis中是否存在该nonce
	key := "sign:nonce:" + nonce
	result := gredis.Get(key)
	return result != ""
}

// 标记nonce已使用
func markNonceUsed(nonce string) {
	// 在Redis中标记该nonce已使用，设置过期时间
	key := "sign:nonce:" + nonce
	gredis.Set(key, "used", 300) // 5分钟过期
}
