package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"justus/internal/global"
	"justus/pkg/setting"
	"time"
)

var db *gorm.DB

type Model struct {
	ID        int `gorm:"primary_key" json:"id"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
}

func Setup() {
	var err error
	db, err = NewDBEngine()
	if err != nil {
		global.Logger.Error("NewDb err :", err)
	}

}

func NewDBEngine() (*gorm.DB, error) {
	var err error
	newLogger := logger.New(
		global.Logger, // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             time.Second,                      // 慢 SQL 阈值
			LogLevel:                  setting.DatabaseSetting.LogLevel, // 日志级别
			IgnoreRecordNotFoundError: true,                             // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                            // 禁用彩色打印
		},
	)
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(mysql.Open(fmt.Sprintf(s,
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name,
		setting.DatabaseSetting.Charset,
		setting.DatabaseSetting.ParseTime,
	)), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.DatabaseSetting.TablePrefix,
			SingularTable: true,
		},
	})
	if err != nil {
		global.Logger.Error("NewDb err :", err)
		return nil, err
	}

	db.Use(dbresolver.Register(dbresolver.Config{}).
		SetConnMaxIdleTime(time.Hour).
		SetConnMaxLifetime(24 * time.Hour).
		SetMaxIdleConns(10).
		SetMaxOpenConns(100))

	return db, nil
}

func GetDb() *gorm.DB {
	return db
}
