package app

import (
	"justus/internal/global"

	"github.com/astaxie/beego/validation"
	"github.com/sirupsen/logrus"
)

// MarkErrors logs validation errors using structured logging
func MarkErrors(errors []*validation.Error) {
	if global.Logger == nil {
		return
	}

	// 使用结构化日志记录所有验证错误
	for _, err := range errors {
		global.Logger.WithFields(logrus.Fields{
			"key":     err.Key,
			"message": err.Message,
			"field":   err.Field,
		}).Warn("参数验证失败")
	}

	return
}
