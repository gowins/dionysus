package distributedlock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	defaultLockKey    = "dioRedisLockKey"
	defaultExpiration = 10 * time.Second
	defaultRetryTTL   = time.Second
	//luaRefresh = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("expire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL    = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("redislock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("redislock: lock not held")
)

type RedisLock struct {
	client     *redis.Client
	lockKey    string
	expiration time.Duration
	retryTTL   time.Duration
	unlockCh   chan struct{}
}

type Option func(redisLock *RedisLock)

func WithLockKey(lockKey string) Option {
	return func(redisLock *RedisLock) {
		redisLock.lockKey = lockKey
	}
}

func WithExpiration(expiration time.Duration) Option {
	return func(redisLock *RedisLock) {
		redisLock.expiration = expiration
	}
}

func WithRetryTTL(retryTTL time.Duration) Option {
	return func(redisLock *RedisLock) {
		redisLock.retryTTL = retryTTL
	}
}

func New(rclient *redis.Client, opts ...Option) *RedisLock {
	redisLock := &RedisLock{
		client:     rclient,
		lockKey:    defaultLockKey,
		expiration: defaultExpiration,
		retryTTL:   defaultRetryTTL,
	}
	for _, opt := range opts {
		opt(redisLock)
	}
	return redisLock
}

func (rl *RedisLock) Lock() {
	lockValue, err := getLockValue()
	if err != nil {
		fmt.Printf("get lock value error %v\n", err)
		return
	}

	for {
		ok, err := rl.client.SetNX(context.Background(), defaultLockKey, lockValue, defaultExpiration).Result()
		if err != nil {
			fmt.Printf("set lock error %v\n", err)
			return
		}
		if ok {
			fmt.Printf("get lock success %v time %v\n", lockValue, time.Now().String())
			return
		}
		fmt.Printf("--------------get lock failed %v time %v---------------\n", lockValue, time.Now().String())
		time.Sleep(defaultRetryTTL)
	}
}

func (rl *RedisLock) Unlock() {
	lockValue, err := getLockValue()
	if err != nil {
		fmt.Printf("get lock value error %v\n", err)
	}
	res, err := luaRelease.Run(context.Background(), rl.client, []string{defaultLockKey}, lockValue).Result()
	if err == redis.Nil {
		fmt.Printf("release error %v\n", ErrLockNotHeld)
		return
	} else if err != nil {
		fmt.Printf("release error %v\n", err)
		return
	}

	if i, ok := res.(int64); !ok || i != 1 {
		fmt.Printf("release error res value type\n")
		return
	}
	fmt.Printf("release lock %v at time %v\n", lockValue, time.Now().String())
	return
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (rl *RedisLock) TTL() (time.Duration, error) {
	lockValue, err := getLockValue()
	if err != nil {
		fmt.Printf("get lock value error %v\n", err)
	}
	res, err := luaPTTL.Run(context.Background(), rl.client, []string{defaultLockKey}, lockValue).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if num := res.(int64); num > 0 {
		return time.Duration(num) * time.Millisecond, nil
	}
	return 0, nil
}
