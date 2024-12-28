package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
	Redis  *redis.Client
)
