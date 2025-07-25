package infrastructure

import (
	"justus/internal/container"
	"justus/internal/global"

	"github.com/sirupsen/logrus"
)

// LoggerImpl Logger接口的实现
type LoggerImpl struct {
	logger *logrus.Logger
}

// NewLogger 创建Logger实例
func NewLogger() container.Logger {
	return &LoggerImpl{
		logger: global.Logger,
	}
}

// Debug 调试日志
func (l *LoggerImpl) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Info 信息日志
func (l *LoggerImpl) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Warn 警告日志
func (l *LoggerImpl) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Error 错误日志
func (l *LoggerImpl) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Debugf 格式化调试日志
func (l *LoggerImpl) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof 格式化信息日志
func (l *LoggerImpl) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warnf 格式化警告日志
func (l *LoggerImpl) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Errorf 格式化错误日志
func (l *LoggerImpl) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// WithFields 结构化日志
func (l *LoggerImpl) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}
