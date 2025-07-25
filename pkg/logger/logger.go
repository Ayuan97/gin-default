package logger

import (
	"io"
	"justus/internal/global"
	"justus/pkg/setting"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(s *setting.LoggerSettingS) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		PrettyPrint: false,
	}

	switch s.LogType {
	case setting.LogFileType:
		fileWriter := &lumberjack.Logger{
			Filename:  s.LogFileSavePath + "/" + s.LogFileName + s.LogFileExt,
			MaxSize:   600,
			MaxAge:    10,
			LocalTime: true,
		}

		logger.Out = io.MultiWriter(os.Stdout, fileWriter)

	case setting.LogFileSLS:
		// 阿里云SLS日志服务
		slsHook, err := NewSLSHook(&s.SLS)
		if err != nil {
			return nil, err
		}

		logger.Hooks.Add(slsHook)

		// 同时输出到控制台
		logger.Out = os.Stdout

	default:
		// 默认输出到控制台
		logger.Out = os.Stdout
	}

	return logger, nil
}

func Setup() {
	var err error
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}
}

func setupLogger() error {
	logger, err := New(setting.LoggerSetting)
	if err != nil {
		return err
	}
	global.Logger = logger

	return nil
}
