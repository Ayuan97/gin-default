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
const LogFileSLS LogType = "sls"
const LogFileZinc LogType = "zinc"

// MiddlewareLogLevel 中间件日志级别
type MiddlewareLogLevel string

const (
	LogLevelBasic    MiddlewareLogLevel = "basic"
	LogLevelDetailed MiddlewareLogLevel = "detailed"
	LogLevelFull     MiddlewareLogLevel = "full"
)

// MiddlewareLogConfig 中间件日志配置
type MiddlewareLogConfig struct {
	Enabled            bool               `yaml:"Enabled"`
	Level              MiddlewareLogLevel `yaml:"Level"`
	EnableRequestBody  bool               `yaml:"EnableRequestBody"`
	EnableResponseBody bool               `yaml:"EnableResponseBody"`
	MaxBodySize        int                `yaml:"MaxBodySize"`
}

// SLSConfig 阿里云SLS配置
type SLSConfig struct {
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	Endpoint        string `yaml:"Endpoint"`
	Project         string `yaml:"Project"`
	Logstore        string `yaml:"Logstore"`
}

// ZincSearchConfig ZincSearch配置
type ZincSearchConfig struct {
	Host         string `yaml:"Host"`
	Username     string `yaml:"Username"`
	Password     string `yaml:"Password"`
	Timeout      int    `yaml:"Timeout"`
	DefaultIndex string `yaml:"DefaultIndex"`
}

type LoggerSettingS struct {
	LogType         LogType
	LogFileSavePath string
	LogFileName     string
	LogFileExt      string
	LogZincHost     string
	LogZincIndex    string
	LogZincUser     string
	LogZincPassword string
	// 阿里云SLS配置
	SLS SLSConfig `yaml:"SLS"`
	// 中间件日志配置
	MiddlewareLog MiddlewareLogConfig `yaml:"MiddlewareLog"`
}

var LoggerSetting = &LoggerSettingS{}

var ZincSearchSetting = &ZincSearchConfig{}

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

// GetMiddlewareLogConfig 获取中间件日志配置
func GetMiddlewareLogConfig() *MiddlewareLogConfig {
	if LoggerSetting == nil {
		return &MiddlewareLogConfig{
			Enabled:            false,
			Level:              LogLevelBasic,
			EnableRequestBody:  false,
			EnableResponseBody: false,
			MaxBodySize:        1024,
		}
	}
	return &LoggerSetting.MiddlewareLog
}

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
	// SLS 环境变量绑定
	v.BindEnv("log.SLS.AccessKeyID", "SLS_ACCESS_KEY_ID")
	v.BindEnv("log.SLS.AccessKeySecret", "SLS_ACCESS_KEY_SECRET")
	v.BindEnv("log.SLS.Endpoint", "SLS_ENDPOINT")
	v.BindEnv("log.SLS.Project", "SLS_PROJECT")
	v.BindEnv("log.SLS.Logstore", "SLS_LOGSTORE")

	// ZincSearch 环境变量绑定
	v.BindEnv("zincsearch.Host", "ZINC_HOST")
	v.BindEnv("zincsearch.Username", "ZINC_USERNAME")
	v.BindEnv("zincsearch.Password", "ZINC_PASSWORD")
	v.BindEnv("zincsearch.Timeout", "ZINC_TIMEOUT")
	v.BindEnv("zincsearch.DefaultIndex", "ZINC_DEFAULT_INDEX")

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
	if err := v.UnmarshalKey("zincsearch", ZincSearchSetting); err != nil {
		panic(fmt.Errorf("Unmarshal zincsearch config error: %w", err))
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
