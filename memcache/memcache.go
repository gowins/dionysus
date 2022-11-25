package memcache

import (
	"context"
	"errors"
	"time"

	"github.com/allegro/bigcache/v3"
)

var (
	ErrEntryIsDead = errors.New("entry has expired")
)

var (
	// defaultHardMaxCacheSize 每个分片内存最大限制，单位MB，缓存大小不能超过此值
	defaultHardMaxCacheSize = 1024
	// defaultLifeWindow 整体缓存有效时间
	defaultLifeWindow = 1 * time.Hour
	// defaultCleanWindow 若此值大于0，则每隔CleanWindow时间间隔，清理一次过期缓存
	// 若不大于0，bigcache在每次设置缓存时，会判断最早的key是否过期，过期则清理
	defaultCleanWindow = 10 * time.Minute
)

type bigCache struct {
	cache *bigcache.BigCache
}

// NewBigCache return Cacher implementation
func NewBigCache(ctx context.Context, opts ...ConfigOpt) (*bigCache, error) {
	cfg := bigcache.DefaultConfig(defaultLifeWindow)
	cfg.HardMaxCacheSize = defaultHardMaxCacheSize
	cfg.CleanWindow = defaultCleanWindow
	for _, opt := range opts {
		opt(&cfg)
	}
	c, err := bigcache.New(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &bigCache{cache: c}, nil
}

// Close close bigCache
func (c *bigCache) Close() error {
	return c.cache.Close()
}

// Delete delete cache for the specific key
func (c *bigCache) Delete(key string) error {
	return c.cache.Delete(key)
}

// Get get cache for the specific key, without time to live
func (c *bigCache) Get(key string) ([]byte, error) {
	return c.cache.Get(key)
}

// Set set key without time to live
func (c *bigCache) Set(key string, value []byte) error {
	return c.cache.Set(key, value)
}

func (c *bigCache) Len() int {
	return c.cache.Len()
}

func (c *bigCache) Capacity() int {
	return c.cache.Capacity()
}
