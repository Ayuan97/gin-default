package logger

import (
	"github.com/kenshaw/sdhook"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"justus/internal/global"
	"justus/pkg/setting"
	"log"
)

func New(s *setting.LoggerSettingS) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{
		PrettyPrint: true,
	}
	switch s.LogType {
	case setting.LogFileType:
		logger.Out = &lumberjack.Logger{
			Filename:  s.LogFileSavePath + "/" + s.LogFileName + s.LogFileExt,
			MaxSize:   600,
			MaxAge:    10,
			LocalTime: true,
		}
	case setting.LogFileLogging:
		hook, err := sdhook.New(
			sdhook.GoogleServiceAccountCredentialsFile("google-tactile-alloy-283608.json"),
			sdhook.LogName("justus.api.production"),
		)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Hooks.Add(hook)
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
