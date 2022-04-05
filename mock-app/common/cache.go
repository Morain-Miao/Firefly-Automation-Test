package common

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var cacheAdapter *cache.Cache

func init() {
	// 创建一个默认过期时间为5分钟的缓存适配器
	// 每60清除一次过期的项目
	cacheAdapter = cache.New(-1*time.Minute, 60*time.Second)
}

// SetCache 添加缓存 如果存在会覆盖
func SetCache(key string, value interface{}, d time.Duration) {
	cacheAdapter.Set(key, value, d)
}

func GetCache(key string) (interface{}, bool) {
	return cacheAdapter.Get(key)
}

// SetDefaultCache 设置cache 无时间参数
func SetDefaultCache(key string, value interface{}) {
	cacheAdapter.SetDefault(key, value)
}

// DeleteCache 删除 cache
func DeleteCache(key string) {
	cacheAdapter.Delete(key)
}

// AddCache Add() 加入缓存 如果存在则返回错误
func AddCache(key string, value interface{}, d time.Duration) {
	cacheAdapter.Add(key, value, d)
}

// IncrementIntCache IncrementInt() 对已存在的key 值自增n
func IncrementIntCache(key string, n int) (num int, err error) {
	return cacheAdapter.IncrementInt(key, n)
}
