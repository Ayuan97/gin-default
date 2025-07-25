package container

import (
	"justus/internal/models"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger 日志接口
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
}

// Cache 缓存接口
type Cache interface {
	Set(key string, data interface{}, expiration time.Duration) error
	Get(key string) string
	Del(key string) (int64, error)
}

// UserRepository 用户数据访问接口
type UserRepository interface {
	GetByID(id int) (*models.User, error)
	GetByIDs(ids []int) ([]*models.User, error)
	GetUsers(page, limit int, keyword, status string) ([]*models.User, int64, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
}

// AdminUserRepository 管理员用户数据访问接口
type AdminUserRepository interface {
	GetByID(id int) (*models.AdminUser, error)
	GetByUsername(username string) (*models.AdminUser, error)
	Create(user *models.AdminUser) error
	Update(user *models.AdminUser) error
	Delete(id int) error
}

// UserService 用户服务接口
type UserService interface {
	GetUserInfo(id int) (*models.User, error)
	GetUsersByIDs(ids []int) ([]*models.User, error)
	GetUsers(page, limit int, keyword, status string) ([]*models.User, int64, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

// AdminUserService 管理员用户服务接口
type AdminUserService interface {
	GetAdminUserInfo(id int) (*models.AdminUser, error)
	GetByUsername(username string) (*models.AdminUser, error)
	CreateAdminUser(user *models.AdminUser) error
	UpdateAdminUser(user *models.AdminUser) error
	DeleteAdminUser(id int) error
}

// Container 依赖注入容器
type Container struct {
	// Infrastructure
	Logger Logger
	Cache  Cache

	// Repositories
	UserRepo      UserRepository
	AdminUserRepo AdminUserRepository

	// Services
	UserService      UserService
	AdminUserService AdminUserService
}

// NewContainer 创建新的依赖注入容器
func NewContainer() *Container {
	return &Container{}
}

// 全局容器实例
var GlobalContainer *Container

// InitContainer 初始化全局容器
func InitContainer() {
	GlobalContainer = NewContainer()
}
