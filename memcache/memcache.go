package memcache

import (
	"context"
	"errors"
	"time"

	"github.com/allegro/bigcache/v3"
)

var (
	ErrEntryIsDead = errors.New("entry has expired")
	ErrTTL         = errors.New("ttl must be greater than 0")
)

var (
	// defaultHardMaxCacheSize 每个分片内存最大限制，单位MB，缓存大小不能超过此值
	defaultHardMaxCacheSize = 1024
	// defaultLifeWindow 整体缓存有效时间
	defaultLifeWindow = 1 * time.Minute
	// defaultCleanWindow 若此值大于0，则每隔CleanWindow时间间隔，清理一次过期缓存
	// 若不大于0，bigcache在每次设置缓存时，会判断最早的key是否过期，过期则清理
	defaultCleanWindow = 10 * time.Second
)

const (
	timeBinaryLenV1 = 15
	timeBinaryLenV2 = 16
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

// SetTTL set key with time to live
func (c *bigCache) SetTTL(key string, value []byte, ttl time.Duration) error {
	if ttl <= 0 {
		return ErrTTL
	}
	bestBefore := time.Now().Add(ttl)
	timeBinary, err := bestBefore.MarshalBinary()
	if err != nil {
		return err
	}
	v := append(timeBinary, value...)
	return c.cache.Set(key, v)
}

// Delete delete cache for the specific key
func (c *bigCache) Delete(key string) error {
	return c.cache.Delete(key)
}

// Get get cache for the specific key, with time to live
func (c *bigCache) GetTTL(key string) ([]byte, error) {
	value, err := c.cache.Get(key)
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return value, nil
	}
	var t time.Time
	idx := timeBinaryLenV2
	if value[0] == byte(1) {
		idx = timeBinaryLenV1
	}
	if len(value) < idx {
		return value, nil
	}
	if err = t.UnmarshalBinary(value[:idx]); err != nil {
		return nil, err
	}
	if y := t.Year(); y < 0 || y > 10000 {
		return value, nil
	}
	if time.Now().After(t) {
		return nil, ErrEntryIsDead
	}
	return value[idx:], nil
}

// Get get cache for the specific key, without time to live
func (c *bigCache) Get(key string) ([]byte, error) {
	return c.cache.Get(key)
}

// Set set key without time to live
func (c *bigCache) Set(key string, value []byte) error {
	return c.cache.Set(key, value)
}
