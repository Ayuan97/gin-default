package infrastructure

import (
	"justus/internal/container"
	"justus/pkg/gredis"
	"time"
)

// CacheImpl Cache接口的实现
type CacheImpl struct{}

// NewCache 创建Cache实例
func NewCache() container.Cache {
	return &CacheImpl{}
}

// Set 设置缓存
func (c *CacheImpl) Set(key string, data interface{}, expiration time.Duration) error {
	return gredis.Set(key, data, expiration)
}

// Get 获取缓存
func (c *CacheImpl) Get(key string) string {
	return gredis.Get(key)
}

// Del 删除缓存
func (c *CacheImpl) Del(key string) (int64, error) {
	return gredis.Del(key)
}
