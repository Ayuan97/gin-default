package zincsearch

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ZincHook ZincSearch 日志钩子
type ZincHook struct {
	client      *Client
	index       string
	levels      []logrus.Level
	mutex       sync.RWMutex
	buffer      []map[string]interface{}
	batchSize   int
	flushTicker *time.Ticker
	stopChan    chan struct{}
}

// NewZincHook 创建新的 ZincSearch 日志钩子
func NewZincHook(client *Client, index string) *ZincHook {
	hook := &ZincHook{
		client:    client,
		index:     index,
		levels:    logrus.AllLevels,
		buffer:    make([]map[string]interface{}, 0),
		batchSize: 10, // 默认批量大小
		stopChan:  make(chan struct{}),
	}

	// 启动定时刷新
	hook.flushTicker = time.NewTicker(5 * time.Second)
	go hook.flushLoop()

	return hook
}

// Levels 返回支持的日志级别
func (hook *ZincHook) Levels() []logrus.Level {
	return hook.levels
}

// Fire 处理日志条目
func (hook *ZincHook) Fire(entry *logrus.Entry) error {
	hook.mutex.Lock()
	defer hook.mutex.Unlock()

	// 转换日志条目为 map
	logEntry := make(map[string]interface{})
	logEntry["@timestamp"] = entry.Time.Format(time.RFC3339)
	logEntry["level"] = entry.Level.String()
	logEntry["message"] = entry.Message
	logEntry["logger"] = "justus-go"

	// 添加字段
	for key, value := range entry.Data {
		logEntry[key] = value
	}

	// 添加到缓冲区
	hook.buffer = append(hook.buffer, logEntry)

	// 如果达到批量大小，立即刷新
	if len(hook.buffer) >= hook.batchSize {
		return hook.flush()
	}

	return nil
}

// SetLevels 设置支持的日志级别
func (hook *ZincHook) SetLevels(levels []logrus.Level) {
	hook.levels = levels
}

// SetBatchSize 设置批量大小
func (hook *ZincHook) SetBatchSize(size int) {
	hook.mutex.Lock()
	defer hook.mutex.Unlock()
	hook.batchSize = size
}

// flush 刷新缓冲区到 ZincSearch
func (hook *ZincHook) flush() error {
	if len(hook.buffer) == 0 {
		return nil
	}

	// 确保索引存在
	exists, err := hook.client.IndexExists(hook.index)
	if err != nil {
		return fmt.Errorf("check index exists: %w", err)
	}

	if !exists {
		if err := hook.client.CreateIndex(hook.index, nil); err != nil {
			return fmt.Errorf("create index: %w", err)
		}
	}

	// 批量索引日志
	_, err = hook.client.IndexDocuments(hook.index, hook.buffer)
	if err != nil {
		return fmt.Errorf("index log documents: %w", err)
	}

	// 清空缓冲区
	hook.buffer = hook.buffer[:0]
	return nil
}

// flushLoop 定时刷新循环
func (hook *ZincHook) flushLoop() {
	for {
		select {
		case <-hook.flushTicker.C:
			hook.mutex.Lock()
			if err := hook.flush(); err != nil {
				// 记录刷新错误，但不阻塞日志系统
				fmt.Printf("ZincHook flush error: %v\n", err)
			}
			hook.mutex.Unlock()
		case <-hook.stopChan:
			return
		}
	}
}

// Close 关闭钩子，刷新剩余日志
func (hook *ZincHook) Close() error {
	close(hook.stopChan)
	hook.flushTicker.Stop()

	hook.mutex.Lock()
	defer hook.mutex.Unlock()

	return hook.flush()
}

// LogEntry 日志条目结构
type LogEntry struct {
	Timestamp time.Time              `json:"@timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Logger    string                 `json:"logger"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// LogToZinc 直接记录日志到 ZincSearch
func LogToZinc(client *Client, index, level, message string, fields map[string]interface{}) error {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Logger:    "justus-go",
		Fields:    fields,
	}

	doc := map[string]interface{}{
		"@timestamp": entry.Timestamp.Format(time.RFC3339),
		"level":      entry.Level,
		"message":    entry.Message,
		"logger":     entry.Logger,
	}

	// 添加额外字段
	for key, value := range entry.Fields {
		doc[key] = value
	}

	_, err := client.IndexDocument(index, doc)
	return err
}
