package setting

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
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

var v *viper.Viper

// Setup initialize the configuration instance
func Setup() {
	// 优先加载根目录下的 .env 文件
	_ = godotenv.Load(".env")
	// 读取环境变量 APP_ENV，默认为 dev
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	configFile := fmt.Sprintf("conf/app.%s.yaml", env)

	v = viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 设置环境变量映射
	v.BindEnv("app.JwtSecret", "JWT_SECRET")
	v.BindEnv("app.PrefixUrl", "APP_PREFIX_URL")
	v.BindEnv("app.ImageUrl", "APP_IMAGE_URL")
	v.BindEnv("app.AesKey", "AES_KEY")
	v.BindEnv("server.RunMode", "APP_RUN_MODE")
	v.BindEnv("server.HttpPort", "APP_PORT")
	v.BindEnv("database.User", "DB_USER")
	v.BindEnv("database.Password", "DB_PASSWORD")
	v.BindEnv("database.Host", "DB_HOST")
	v.BindEnv("database.Name", "DB_NAME")
	v.BindEnv("database.TablePrefix", "DB_TABLE_PREFIX")
	v.BindEnv("redis.DB", "REDIS_DB")
	v.BindEnv("redis.Host", "REDIS_HOST")
	v.BindEnv("redis.Password", "REDIS_PASSWORD")
	v.BindEnv("redis.Prefix", "REDIS_PREFIX")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w", err))
	}

	if err := v.UnmarshalKey("app", AppSetting); err != nil {
		panic(fmt.Errorf("Unmarshal app config error: %w", err))
	}
	if err := v.UnmarshalKey("server", ServerSetting); err != nil {
		panic(fmt.Errorf("Unmarshal server config error: %w", err))
	}
	if err := v.UnmarshalKey("database", DatabaseSetting); err != nil {
		panic(fmt.Errorf("Unmarshal database config error: %w", err))
	}
	if err := v.UnmarshalKey("redis", RedisSetting); err != nil {
		panic(fmt.Errorf("Unmarshal redis config error: %w", err))
	}
	if err := v.UnmarshalKey("log", LoggerSetting); err != nil {
		panic(fmt.Errorf("Unmarshal log config error: %w", err))
	}

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second

	// 配置热加载
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		_ = v.UnmarshalKey("app", AppSetting)
		_ = v.UnmarshalKey("server", ServerSetting)
		_ = v.UnmarshalKey("database", DatabaseSetting)
		_ = v.UnmarshalKey("redis", RedisSetting)
		_ = v.UnmarshalKey("log", LoggerSetting)
	})
}
