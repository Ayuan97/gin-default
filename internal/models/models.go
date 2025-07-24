package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"justus/internal/global"
	"justus/pkg/setting"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// GormTime 自定义时间类型，与MySQL的timestamp兼容
type GormTime struct {
	time.Time
}

// GormDate 自定义日期类型，与MySQL的date兼容
type GormDate struct {
	time.Time
}

// Scan 实现sql.Scanner接口
func (gt *GormTime) Scan(value interface{}) error {
	if value == nil {
		*gt = GormTime{Time: time.Time{}}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*gt = GormTime{Time: v}
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(v))
		if err != nil {
			return err
		}
		*gt = GormTime{Time: t}
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*gt = GormTime{Time: t}
	default:
		return fmt.Errorf("cannot scan %T into GormTime", value)
	}
	return nil
}

// Value 实现driver.Valuer接口
func (gt GormTime) Value() (driver.Value, error) {
	if gt.Time.IsZero() {
		return nil, nil
	}
	return gt.Time, nil
}

// MarshalJSON 实现JSON序列化
func (gt GormTime) MarshalJSON() ([]byte, error) {
	if gt.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, gt.Time.Format("2006-01-02 15:04:05"))), nil
}

// UnmarshalJSON 实现JSON反序列化
func (gt *GormTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	t, err := time.Parse(`"2006-01-02 15:04:05"`, string(data))
	if err != nil {
		return err
	}
	*gt = GormTime{Time: t}
	return nil
}

// Scan 实现sql.Scanner接口
func (gd *GormDate) Scan(value interface{}) error {
	if value == nil {
		*gd = GormDate{Time: time.Time{}}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		*gd = GormDate{Time: v}
	case []byte:
		t, err := time.Parse("2006-01-02", string(v))
		if err != nil {
			return err
		}
		*gd = GormDate{Time: t}
	case string:
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return err
		}
		*gd = GormDate{Time: t}
	default:
		return fmt.Errorf("cannot scan %T into GormDate", value)
	}
	return nil
}

// Value 实现driver.Valuer接口
func (gd GormDate) Value() (driver.Value, error) {
	if gd.Time.IsZero() {
		return nil, nil
	}
	return gd.Time.Format("2006-01-02"), nil
}

// MarshalJSON 实现JSON序列化
func (gd GormDate) MarshalJSON() ([]byte, error) {
	if gd.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, gd.Time.Format("2006-01-02"))), nil
}

// UnmarshalJSON 实现JSON反序列化
func (gd *GormDate) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	t, err := time.Parse(`"2006-01-02"`, string(data))
	if err != nil {
		return err
	}
	*gd = GormDate{Time: t}
	return nil
}

// Setup 初始化数据库连接
func Setup() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name)

	// 配置GORM日志级别
	var logLevel logger.LogLevel
	if setting.ServerSetting.RunMode == "debug" {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		global.Logger.Fatalf("models.Setup err: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		global.Logger.Fatalf("db.DB() err: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// GetDb 获取数据库连接
func GetDb() *gorm.DB {
	return db
}
