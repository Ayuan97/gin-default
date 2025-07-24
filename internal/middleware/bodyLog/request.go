package bodyLog

import (
	"bytes"
	"encoding/json"
	"io"
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

	// 构建基础日志字段
	fields := logrus.Fields{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"status_code": c.Writer.Status(),
		"duration_ms": duration.Milliseconds(),
		"client_ip":   c.ClientIP(),
	}

	// 添加用户信息（如果有）
	if userID, exists := c.Get("userId"); exists {
		if uid, ok := userID.(int); ok && uid > 0 {
			fields["user_id"] = uid
		}
	}
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok && name != "" {
			fields["username"] = name
		}
	}

	// 记录请求头
	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	fields["headers"] = headers

	// 记录查询参数
	if len(c.Request.URL.RawQuery) > 0 {
		fields["query"] = c.Request.URL.RawQuery
	}

	// 记录请求体
	if requestBody != "" {
		// 尝试解析为JSON
		var jsonBody interface{}
		if err := json.Unmarshal([]byte(requestBody), &jsonBody); err == nil {
			fields["request_body"] = jsonBody
		} else {
			fields["request_body"] = requestBody
		}
	}

	// 记录响应体
	if blw.body.Len() > 0 {
		responseBody := blw.body.String()
		// 尝试解析为JSON
		var jsonResponse interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &jsonResponse); err == nil {
			fields["response_body"] = jsonResponse
		} else {
			fields["response_body"] = responseBody
		}
	}

	// 根据状态码决定日志级别
	var logLevel logrus.Level
	switch {
	case c.Writer.Status() >= 500:
		logLevel = logrus.ErrorLevel
	case c.Writer.Status() >= 400:
		logLevel = logrus.WarnLevel
	default:
		logLevel = logrus.InfoLevel
	}

	// 输出完整的结构化日志
	global.Logger.WithFields(fields).Log(logLevel, "HTTP Request")
}
