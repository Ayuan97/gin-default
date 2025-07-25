package logger

import (
	"encoding/json"
	"fmt"
	"justus/pkg/setting"
	"reflect"
	"strconv"
	"strings"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/sirupsen/logrus"
)

// SLSHook 阿里云SLS Hook
type SLSHook struct {
	client   sls.ClientInterface
	project  string
	logstore string
}

// NewSLSHook 创建新的SLS Hook
func NewSLSHook(config *setting.SLSConfig) (*SLSHook, error) {
	client := sls.CreateNormalInterface(
		config.Endpoint,
		config.AccessKeyID,
		config.AccessKeySecret,
		"",
	)

	return &SLSHook{
		client:   client,
		project:  config.Project,
		logstore: config.Logstore,
	}, nil
}

// Fire 实现logrus.Hook接口
func (hook *SLSHook) Fire(entry *logrus.Entry) error {
	timestamp := uint32(entry.Time.Unix())
	topic := "justus-api"
	source := "justus-go"

	// 创建基础的日志内容
	var contents []*sls.LogContent

	// 添加基础字段
	timestampStr := strconv.FormatInt(entry.Time.Unix(), 10)
	level := entry.Level.String()
	message := entry.Message

	contents = append(contents, &sls.LogContent{
		Key:   strPtr("timestamp"),
		Value: strPtr(timestampStr),
	})
	contents = append(contents, &sls.LogContent{
		Key:   strPtr("level"),
		Value: strPtr(level),
	})
	contents = append(contents, &sls.LogContent{
		Key:   strPtr("message"),
		Value: strPtr(message),
	})

	// 如果有调用者信息，添加到日志中
	if entry.HasCaller() {
		caller := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		contents = append(contents, &sls.LogContent{
			Key:   strPtr("caller"),
			Value: strPtr(caller),
		})
	}

	// 添加自定义字段
	for k, v := range entry.Data {
		key := k
		value := formatValue(v)
		contents = append(contents, &sls.LogContent{
			Key:   strPtr(key),
			Value: strPtr(value),
		})
	}

	// 创建日志条目
	log := &sls.Log{
		Time:     &timestamp,
		Contents: contents,
	}

	// 创建日志组
	logGroup := &sls.LogGroup{
		Topic:  &topic,
		Source: &source,
		Logs:   []*sls.Log{log},
	}

	// 发送日志到SLS
	err := hook.client.PostLogStoreLogs(hook.project, hook.logstore, logGroup, nil)
	if err != nil {
		// 如果发送失败，不要返回错误，避免影响正常业务
		fmt.Printf("Failed to send log to SLS: %v\n", err)
	}

	return nil
}

// Levels 实现logrus.Hook接口，返回hook处理的日志级别
func (hook *SLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// strPtr 辅助函数，返回字符串指针
func strPtr(s string) *string {
	return &s
}

// formatValue 格式化值，将复杂类型转换为JSON格式
func formatValue(v interface{}) string {
	if v == nil {
		return "<nil>"
	}

	// 首先尝试直接JSON序列化
	if jsonBytes, err := json.Marshal(v); err == nil {
		// 检查是否是字符串类型
		if str, ok := v.(string); ok {
			// 如果是字符串，检查是否是Go的map格式
			if isGoMapFormat(str) {
				// 尝试清理格式并返回原字符串或更好的格式
				return cleanGoMapFormat(str)
			}
			return str
		}
		return string(jsonBytes)
	}

	// 获取值的反射类型
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
		// 强制使用fmt.Sprintf然后清理格式
		goFormat := fmt.Sprintf("%v", v)
		return cleanGoMapFormat(goFormat)
	case reflect.String:
		str := rv.String()
		if isGoMapFormat(str) {
			return cleanGoMapFormat(str)
		}
		return str
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	default:
		// 最后回退到默认格式
		return fmt.Sprintf("%v", v)
	}
}

// isGoMapFormat 检查字符串是否是Go的map格式
func isGoMapFormat(s string) bool {
	return strings.HasPrefix(s, "map[") && strings.HasSuffix(s, "]")
}

// cleanGoMapFormat 清理Go map格式，转换为更友好的JSON格式
func cleanGoMapFormat(s string) string {
	// 如果不是map格式，直接返回
	if !isGoMapFormat(s) {
		return s
	}

	// 移除 "map[" 前缀和 "]" 后缀
	content := s[4 : len(s)-1]
	if content == "" {
		return "{}"
	}

	// 构建JSON格式
	var jsonPairs []string

	// 简单解析键值对（这里做基础处理）
	pairs := parseMapPairs(content)
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 为键添加引号
			jsonKey := fmt.Sprintf(`"%s"`, key)

			// 处理值
			var jsonValue string
			if value == "<nil>" {
				jsonValue = "null"
			} else if isNumeric(value) {
				jsonValue = value
			} else if value == "true" || value == "false" {
				jsonValue = value
			} else {
				jsonValue = fmt.Sprintf(`"%s"`, value)
			}

			jsonPairs = append(jsonPairs, fmt.Sprintf("%s:%s", jsonKey, jsonValue))
		}
	}

	return fmt.Sprintf("{%s}", strings.Join(jsonPairs, ","))
}

// parseMapPairs 解析map中的键值对
func parseMapPairs(content string) []string {
	var pairs []string
	var current strings.Builder
	depth := 0

	for _, r := range content {
		switch r {
		case '[', '(':
			depth++
			current.WriteRune(r)
		case ']', ')':
			depth--
			current.WriteRune(r)
		case ' ':
			if depth == 0 && current.Len() > 0 {
				pairs = append(pairs, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		pairs = append(pairs, current.String())
	}

	return pairs
}

// isNumeric 检查字符串是否是数字
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
