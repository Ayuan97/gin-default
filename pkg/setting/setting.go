package setting

import (
	"fmt"
	"github.com/go-ini/ini"
	"gorm.io/gorm/logger"
	"time"
)

type LogType string

const LogFileType LogType = "file"
const LogFileLogging LogType = "logging"

type LoggerSettingS struct {
	LogType         LogType
	LogFileSavePath string
	LogFileName     string
	LogFileExt      string
	LogZincHost     string
	LogZincIndex    string
	LogZincUser     string
	LogZincPassword string
}

var LoggerSetting = &LoggerSettingS{}

type App struct {
	AesKey    string
	JwtSecret string
	PageSize  int
	PrefixUrl string
	ImageUrl  string
	H5Url     string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
	Charset     string
	ParseTime   bool
	LogLevel    logger.LogLevel
}

var DatabaseSetting = &Database{}

type Redis struct {
	DB          int
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Prefix      string
}

var RedisSetting = &Redis{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/connor-local.ini")
	if err != nil {
		fmt.Printf("setting.Setup, fail to parse 'conf/app.ini': %v\n", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("log", LoggerSetting)
	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		fmt.Printf("Cfg.MapTo %s err: %v\n", section, err)
	}
}
