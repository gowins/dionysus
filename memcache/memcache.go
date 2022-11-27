package memcache

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

type cacheStore struct {
	cache map[string]*bigcache.BigCache
	sync.RWMutex
}

var store = cacheStore{
	cache: map[string]*bigcache.BigCache{},
}

// NewBigCache return Cacher implementation
func NewBigCache(ctx context.Context, name string, opts ...ConfigOpt) error {
	cfg := bigcache.DefaultConfig(defaultLifeWindow)
	cfg.HardMaxCacheSize = defaultHardMaxCacheSize
	cfg.CleanWindow = defaultCleanWindow
	for _, opt := range opts {
		opt(&cfg)
	}
	c, err := bigcache.New(ctx, cfg)
	if err != nil {
		return err
	}
	store.Lock()
	store.cache[name] = c
	store.Unlock()
	return nil
}

// Close close bigCache
func Close(name string) error {
	store.RLock()
	bCache, ok := store.cache[name]
	store.RUnlock()
	if !ok {
		return fmt.Errorf("cache %v is not found", name)
	}
	return bCache.Close()
}

// Delete delete cache for the specific key
func Delete(name string, key string) error {
	store.RLock()
	bCache, ok := store.cache[name]
	store.RUnlock()
	if !ok {
		return fmt.Errorf("cache %v is not found", name)
	}
	return bCache.Delete(key)
}

// Get get cache for the specific key, without time to live
func Get(name string, key string) ([]byte, error) {
	store.RLock()
	bCache, ok := store.cache[name]
	store.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cache %v is not found", name)
	}
	return bCache.Get(key)
}

// Set set key without time to live
func Set(name string, key string, value []byte) error {
	store.RLock()
	bCache, ok := store.cache[name]
	store.RUnlock()
	if !ok {
		return fmt.Errorf("cache %v is not found", name)
	}
	return bCache.Set(key, value)
}

func GetCache(name string) (*bigcache.BigCache, bool) {
	store.RLock()
	bCache, ok := store.cache[name]
	store.RUnlock()
	return bCache, ok
}
