package bodyLog

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"time"

	"justus/internal/global"
	"justus/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 跳过记录的路径列表
var skipPathsList = []string{
	"/health", "/healthz", "/ready", "/live", "/metrics",
	"/favicon.ico", "/robots.txt", "/ping",
}

// bodyLogWriter 响应体记录器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// GinBodyLogMiddleware 创建请求日志中间件
func GinBodyLogMiddleware() gin.HandlerFunc {
	config := setting.GetMiddlewareLogConfig()

	// 如果未启用中间件日志，返回空中间件
	if !config.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// 检查是否跳过此路径
		if shouldSkipPath(c.Request.URL.Path, skipPathsList) {
			c.Next()
			return
		}

		startTime := time.Now()

		// 创建响应体记录器
		blw := &bodyLogWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			requestBody = string(bodyBytes)
		}

		// 执行请求
		c.Next()

		// 计算请求耗时
		duration := time.Since(startTime)

		// 输出日志
		outputRequestLog(c, blw, requestBody, duration)
	}
}

// shouldSkipPath 检查是否应该跳过此路径的日志记录
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// outputRequestLog 输出完整的请求日志
func outputRequestLog(c *gin.Context, blw *bodyLogWriter, requestBody string, duration time.Duration) {
	if global.Logger == nil {
		return
	}

	// 构建完整的请求日志结构
	logData := map[string]interface{}{
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"level":       getLogLevel(c.Writer.Status()),
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"status_code": c.Writer.Status(),
		"duration_ms": duration.Milliseconds(),
		"client_ip":   c.ClientIP(),
	}

	// 添加用户信息（如果有）
	if userID, exists := c.Get("userId"); exists {
		if uid, ok := userID.(int); ok && uid > 0 {
			logData["user_id"] = uid
		}
	}
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok && name != "" {
			logData["username"] = name
		}
	}

	// 记录请求头（只记录自定义头部）
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		// 跳过常见的标准HTTP头部和无关参数
		switch k {
		case "Accept", "Accept-Encoding", "Accept-Language", "Cache-Control",
			"Connection", "Cookie", "Host", "Pragma", "Referer",
			"User-Agent", "Content-Type", "Content-Length", "Origin",
			"Sec-Fetch-Site", "Sec-Fetch-Mode", "Sec-Fetch-Dest",
			"X-Requested-With", "X-Real-IP", "X-Forwarded-For",
			"If-None-Match", "If-Modified-Since", "DNT", "Keep-Alive",
			"Sec-Ch-Ua", "Sec-Ch-Ua-Mobile", "Sec-Ch-Ua-Platform",
			"Sec-Fetch-User", "Upgrade-Insecure-Requests":
			continue
		}
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	if len(headers) > 0 {
		logData["headers"] = headers
	}

	// 记录查询参数 - 解析为JSON对象
	if len(c.Request.URL.RawQuery) > 0 {
		if queryParams, err := url.ParseQuery(c.Request.URL.RawQuery); err == nil {
			queryMap := make(map[string]string)
			for k, v := range queryParams {
				if len(v) > 0 {
					queryMap[k] = v[0]
				}
			}
			logData["query"] = queryMap
		} else {
			logData["query"] = c.Request.URL.RawQuery
		}
	}

	// 记录请求体 - 解析为JSON对象
	if requestBody != "" {
		var parsedBody interface{}
		// 首先尝试解析为JSON
		if err := json.Unmarshal([]byte(requestBody), &parsedBody); err == nil {
			logData["request_body"] = parsedBody
		} else {
			// JSON解析失败，尝试解析为表单数据
			if formData, err := url.ParseQuery(requestBody); err == nil {
				formMap := make(map[string]string)
				for k, v := range formData {
					if len(v) > 0 {
						formMap[k] = v[0]
					}
				}
				logData["request_body"] = formMap
			} else {
				logData["request_body"] = requestBody
			}
		}
	}

	// 记录响应体 - 解析为JSON对象
	if blw.body.Len() > 0 {
		var parsedResponse interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &parsedResponse); err == nil {
			logData["response_body"] = parsedResponse
		} else {
			logData["response_body"] = blw.body.String()
		}
	}

	// 根据状态码决定日志级别并输出JSON格式日志
	logLevel := getLogLevel(c.Writer.Status())

	// 输出标准JSON格式日志
	global.Logger.WithFields(logrus.Fields(logData)).Log(logLevel, "HTTP Request")
}

// getLogLevel 根据HTTP状态码获取日志级别
func getLogLevel(statusCode int) logrus.Level {
	switch {
	case statusCode >= 500:
		return logrus.ErrorLevel
	case statusCode >= 400:
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}
