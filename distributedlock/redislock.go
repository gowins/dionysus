package distributedlock

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	defaultLockKey    = "dioRedisLockKey"
	defaultExpiration = 10 * time.Second
)

type redisLock struct {
	client *redis.Client
}

func (rl *redisLock) Lock() {
	lockValue, err := getLockValue()
	if err != nil {
		panic(err)
	}

	ok, err := rl.client.SetNX(context.Background(), defaultLockKey, lockValue, defaultExpiration).Result()
	if err != nil {
		panic(err)
	}
}
