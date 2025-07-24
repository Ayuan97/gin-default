package bodyLog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"justus/internal/global"
	"justus/pkg/setting"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel 日志记录级别
type LogLevel string

const (
	LogLevelBasic    LogLevel = "basic"    // 只记录基本信息
	LogLevelDetailed LogLevel = "detailed" // 记录详细信息（过滤敏感数据）
	LogLevelFull     LogLevel = "full"     // 记录完整信息（开发环境）
)

// LogConfig 日志中间件配置
type LogConfig struct {
	Level              LogLevel `yaml:"level" json:"level"`                               // 日志级别
	EnableRequestBody  bool     `yaml:"enable_request_body" json:"enable_request_body"`   // 是否记录请求体
	EnableResponseBody bool     `yaml:"enable_response_body" json:"enable_response_body"` // 是否记录响应体
	MaxBodySize        int      `yaml:"max_body_size" json:"max_body_size"`               // 最大记录的请求体大小(字节)
	SensitiveFields    []string `yaml:"sensitive_fields" json:"sensitive_fields"`         // 敏感字段列表
	SkipPaths          []string `yaml:"skip_paths" json:"skip_paths"`                     // 跳过记录的路径
}

// 日志配置常量
const (
	// 默认最大请求体大小（字节）
	DefaultMaxBodySize = 4096 // 4KB
)

// 敏感字段列表常量
var sensitiveFieldsList = []string{
	"password", "passwd", "secret", "token", "key", "authorization",
	"auth", "credential", "sign", "signature", "private", "jwt",
	"session", "cookie", "csrf", "api_key", "apikey", "access_token",
	"refresh_token", "client_secret", "private_key", "passphrase",
}

// 跳过记录的路径列表常量
var skipPathsList = []string{
	"/health", "/healthz", "/ready", "/live", "/metrics",
	"/favicon.ico", "/robots.txt", "/ping",
}

// middlewareLogger 中间件专用日志器
var middlewareLogger *logrus.Logger

// bodyLogWriter 响应体记录器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestLogData 请求日志数据结构
type RequestLogData struct {
	// 基本信息
	Timestamp  time.Time `json:"timestamp"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	Duration   int64     `json:"duration_ms"`
	ClientIP   string    `json:"client_ip"`
	UserAgent  string    `json:"user_agent"`
	RequestURI string    `json:"request_uri"`
	Protocol   string    `json:"protocol"`

	// 扩展信息（详细级别及以上）
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty"`
	RequestBody interface{}       `json:"request_body,omitempty"`

	// 完整信息（完整级别）
	ResponseBody interface{} `json:"response_body,omitempty"`
	RequestSize  int64       `json:"request_size,omitempty"`
	ResponseSize int64       `json:"response_size,omitempty"`

	// 用户信息（如果有）
	UserID   int    `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	UserRole string `json:"user_role,omitempty"`
}

// initMiddlewareLogger 初始化中间件日志器
func initMiddlewareLogger() {
	config := setting.GetMiddlewareLogConfig()

	// 如果已经初始化过或者未启用，直接返回
	if middlewareLogger != nil || !config.Enabled {
		return
	}

	middlewareLogger = logrus.New()

	// 固定使用文件输出
	setupFileOutput()
}

// setupFileOutput 设置文件输出（固定配置）
func setupFileOutput() {
	// 固定的文件输出配置
	logDir := "storage/logs"
	fileName := "middleware"
	fileExt := ".log"

	// 确保目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Errorf("Failed to create log directory: %v", err)
		return
	}

	middlewareLogger.SetOutput(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s%s", logDir, fileName, fileExt),
		MaxSize:    100, // 100MB
		MaxAge:     30,  // 30天
		MaxBackups: 10,  // 10个备份
		LocalTime:  true,
		Compress:   true,
	})

	middlewareLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
	})
}

// GinBodyLogMiddleware 创建请求日志中间件
func GinBodyLogMiddleware() gin.HandlerFunc {
	return GinBodyLogMiddlewareFromConfig()
}

// GinBodyLogMiddlewareFromConfig 从配置文件创建请求日志中间件
func GinBodyLogMiddlewareFromConfig() gin.HandlerFunc {
	config := setting.GetMiddlewareLogConfig()

	// 如果未启用中间件日志，返回空中间件
	if !config.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 初始化中间件日志器
	initMiddlewareLogger()

	// 转换配置格式
	logConfig := LogConfig{
		Level:              LogLevel(config.Level),
		EnableRequestBody:  config.EnableRequestBody,
		EnableResponseBody: config.EnableResponseBody,
		MaxBodySize:        config.MaxBodySize,
		SensitiveFields:    sensitiveFieldsList,
		SkipPaths:          skipPathsList,
	}

	return GinBodyLogMiddlewareWithConfig(logConfig)
}

// GinBodyLogMiddlewareWithConfig 使用自定义配置创建请求日志中间件
func GinBodyLogMiddlewareWithConfig(config LogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过此路径
		if shouldSkipPath(c.Request.URL.Path, config.SkipPaths) {
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

		// 读取并备份请求体
		var requestBody []byte
		var requestBodyData interface{}

		if config.EnableRequestBody && config.Level != LogLevelBasic {
			if c.Request.Body != nil {
				requestBody, _ = io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))

				// 解析请求体
				if len(requestBody) > 0 && len(requestBody) <= config.MaxBodySize {
					if strings.Contains(c.GetHeader("Content-Type"), "application/json") {
						var jsonBody interface{}
						if err := json.Unmarshal(requestBody, &jsonBody); err == nil {
							requestBodyData = filterSensitiveData(jsonBody, config.SensitiveFields)
						}
					} else {
						requestBodyData = string(requestBody)
					}
				} else if len(requestBody) > config.MaxBodySize {
					requestBodyData = "[Body too large, truncated]"
				}
			}
		}

		// 执行请求
		c.Next()

		// 计算请求耗时
		duration := time.Since(startTime)

		// 构建日志数据
		logData := buildLogData(c, blw, requestBodyData, duration, config)

		// 输出日志
		outputLogToMiddlewareLogger(logData, config.Level)
	}
}

// shouldSkipPath 检查是否应该跳过此路径的日志记录
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// buildLogData 构建日志数据
func buildLogData(c *gin.Context, blw *bodyLogWriter, requestBody interface{}, duration time.Duration, config LogConfig) *RequestLogData {
	logData := &RequestLogData{
		Timestamp:  time.Now(),
		Method:     c.Request.Method,
		Path:       c.Request.URL.Path,
		StatusCode: c.Writer.Status(),
		Duration:   duration.Milliseconds(),
		ClientIP:   c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
		RequestURI: c.Request.RequestURI,
		Protocol:   c.Request.Proto,
	}

	// 获取用户信息（如果存在）
	if userID, exists := c.Get("userId"); exists {
		if uid, ok := userID.(int); ok {
			logData.UserID = uid
		}
	}
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			logData.Username = name
		}
	}
	if userRole, exists := c.Get("userRole"); exists {
		if role, ok := userRole.(string); ok {
			logData.UserRole = role
		}
	}

	// 详细级别及以上记录更多信息
	if config.Level == LogLevelDetailed || config.Level == LogLevelFull {
		// 过滤请求头
		logData.Headers = filterHeaders(c.Request.Header, config.SensitiveFields)

		// 查询参数
		logData.QueryParams = make(map[string]string)
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 {
				if isSensitiveField(k, config.SensitiveFields) {
					logData.QueryParams[k] = "[FILTERED]"
				} else {
					logData.QueryParams[k] = v[0]
				}
			}
		}

		// 请求体
		if requestBody != nil {
			logData.RequestBody = requestBody
		}
	}

	// 完整级别记录所有信息
	if config.Level == LogLevelFull {
		// 响应体
		if config.EnableResponseBody && blw.body.Len() > 0 {
			if blw.body.Len() <= config.MaxBodySize {
				var responseBody interface{}
				if json.Unmarshal(blw.body.Bytes(), &responseBody) == nil {
					logData.ResponseBody = responseBody
				} else {
					logData.ResponseBody = blw.body.String()
				}
			} else {
				logData.ResponseBody = "[Response too large, truncated]"
			}
		}

		// 请求和响应大小
		logData.RequestSize = c.Request.ContentLength
		logData.ResponseSize = int64(blw.body.Len())
	}

	return logData
}

// filterHeaders 过滤敏感请求头
func filterHeaders(headers map[string][]string, sensitiveFields []string) map[string]string {
	filtered := make(map[string]string)
	for k, v := range headers {
		if len(v) > 0 {
			if isSensitiveField(k, sensitiveFields) {
				filtered[k] = "[FILTERED]"
			} else {
				filtered[k] = v[0]
			}
		}
	}
	return filtered
}

// filterSensitiveData 递归过滤敏感数据
func filterSensitiveData(data interface{}, sensitiveFields []string) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		filtered := make(map[string]interface{})
		for key, value := range v {
			if isSensitiveField(key, sensitiveFields) {
				filtered[key] = "[FILTERED]"
			} else {
				filtered[key] = filterSensitiveData(value, sensitiveFields)
			}
		}
		return filtered
	case []interface{}:
		filtered := make([]interface{}, len(v))
		for i, item := range v {
			filtered[i] = filterSensitiveData(item, sensitiveFields)
		}
		return filtered
	default:
		return v
	}
}

// isSensitiveField 检查字段是否为敏感字段
func isSensitiveField(field string, sensitiveFields []string) bool {
	field = strings.ToLower(field)
	for _, sensitive := range sensitiveFields {
		sensitivePattern := strings.ToLower(sensitive)
		// 精确匹配或包含匹配
		if field == sensitivePattern || strings.Contains(field, sensitivePattern) {
			return true
		}
	}
	return false
}

// outputLogToMiddlewareLogger 输出日志到中间件日志器
func outputLogToMiddlewareLogger(logData *RequestLogData, level LogLevel) {
	if middlewareLogger == nil {
		return
	}

	// 构建基础日志字段
	fields := logrus.Fields{
		"method":      logData.Method,
		"path":        logData.Path,
		"status_code": logData.StatusCode,
		"duration_ms": logData.Duration,
		"client_ip":   logData.ClientIP,
		"user_agent":  logData.UserAgent,
	}

	// 添加用户信息
	if logData.UserID > 0 {
		fields["user_id"] = logData.UserID
	}
	if logData.Username != "" {
		fields["username"] = logData.Username
	}
	if logData.UserRole != "" {
		fields["user_role"] = logData.UserRole
	}

	// 根据级别添加更多字段
	if level == LogLevelDetailed || level == LogLevelFull {
		if len(logData.Headers) > 0 {
			fields["headers"] = logData.Headers
		}
		if len(logData.QueryParams) > 0 {
			fields["query_params"] = logData.QueryParams
		}
		if logData.RequestBody != nil {
			fields["request_body"] = logData.RequestBody
		}
	}

	if level == LogLevelFull {
		if logData.ResponseBody != nil {
			fields["response_body"] = logData.ResponseBody
		}
		if logData.RequestSize > 0 {
			fields["request_size"] = logData.RequestSize
		}
		if logData.ResponseSize > 0 {
			fields["response_size"] = logData.ResponseSize
		}
		fields["request_uri"] = logData.RequestURI
		fields["protocol"] = logData.Protocol
	}

	// 根据状态码决定日志级别
	var logLevel logrus.Level
	switch {
	case logData.StatusCode >= 500:
		logLevel = logrus.ErrorLevel
	case logData.StatusCode >= 400:
		logLevel = logrus.WarnLevel
	case logData.StatusCode >= 300:
		logLevel = logrus.InfoLevel
	default:
		logLevel = logrus.InfoLevel
	}

	// 输出结构化日志
	middlewareLogger.WithFields(fields).Log(logLevel, "HTTP Request")
}

// outputLog 输出日志 (保持向后兼容)
func outputLog(logData *RequestLogData, level LogLevel) {
	if global.Logger == nil {
		return
	}

	// 构建基础日志字段
	fields := logrus.Fields{
		"method":      logData.Method,
		"path":        logData.Path,
		"status_code": logData.StatusCode,
		"duration_ms": logData.Duration,
		"client_ip":   logData.ClientIP,
		"user_agent":  logData.UserAgent,
	}

	// 添加用户信息
	if logData.UserID > 0 {
		fields["user_id"] = logData.UserID
	}
	if logData.Username != "" {
		fields["username"] = logData.Username
	}
	if logData.UserRole != "" {
		fields["user_role"] = logData.UserRole
	}

	// 根据级别添加更多字段
	if level == LogLevelDetailed || level == LogLevelFull {
		if len(logData.Headers) > 0 {
			fields["headers"] = logData.Headers
		}
		if len(logData.QueryParams) > 0 {
			fields["query_params"] = logData.QueryParams
		}
		if logData.RequestBody != nil {
			fields["request_body"] = logData.RequestBody
		}
	}

	if level == LogLevelFull {
		if logData.ResponseBody != nil {
			fields["response_body"] = logData.ResponseBody
		}
		if logData.RequestSize > 0 {
			fields["request_size"] = logData.RequestSize
		}
		if logData.ResponseSize > 0 {
			fields["response_size"] = logData.ResponseSize
		}
		fields["request_uri"] = logData.RequestURI
		fields["protocol"] = logData.Protocol
	}

	// 根据状态码决定日志级别
	var logLevel logrus.Level
	switch {
	case logData.StatusCode >= 500:
		logLevel = logrus.ErrorLevel
	case logData.StatusCode >= 400:
		logLevel = logrus.WarnLevel
	case logData.StatusCode >= 300:
		logLevel = logrus.InfoLevel
	default:
		logLevel = logrus.InfoLevel
	}

	// 输出结构化日志
	global.Logger.WithFields(fields).Log(logLevel, "HTTP Request")
}
