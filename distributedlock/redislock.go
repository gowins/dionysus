package distributedlock

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	defaultExpiration = 0 * time.Second
	defaultRetryTTL   = time.Second
	luaRefresh        = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("expire", KEYS[1], ARGV[2]) else return 0 end`)
	luaRelease        = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("del", KEYS[1]) else return 0 end`)
	luaPTTL           = redis.NewScript(`if redis.call("get", KEYS[1]) == ARGV[1] then return redis.call("pttl", KEYS[1]) else return -3 end`)
)

type RedisLock struct {
	client         *redis.Client
	lockKey        string
	expiration     time.Duration
	retryTTL       time.Duration
	watchDogEnable bool
}

type Option func(redisLock *RedisLock)

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

func WithWatchDog(enable bool) Option {
	return func(redisLock *RedisLock) {
		redisLock.watchDogEnable = enable
	}
}

func New(rclient *redis.Client, lockKey string, opts ...Option) *RedisLock {
	redisLock := &RedisLock{
		client:         rclient,
		lockKey:        lockKey,
		expiration:     defaultExpiration,
		retryTTL:       defaultRetryTTL,
		watchDogEnable: false,
	}
	for _, opt := range opts {
		opt(redisLock)
	}
	return redisLock
}

func (rl *RedisLock) Lock(ctx context.Context) (context.Context, error) {
	lockValue, err := getLockValue()
	if err != nil {
		return ctx, fmt.Errorf("get lock value error %v", err)
	}

	for {
		ok, err := rl.client.SetNX(ctx, rl.lockKey, lockValue, rl.expiration).Result()
		if err != nil {
			log.Errorf("set lock error %v", err)
			return ctx, fmt.Errorf("set lock error %v", err)
		}
		if ok {
			log.Infof("get lock success %v time %v", lockValue, time.Now().String())
			if rl.watchDogEnable && rl.expiration > 3*time.Second {
				nctx, cancelFunc := context.WithCancel(ctx)
				go rl.watchDog(nctx, cancelFunc, lockValue)
				return nctx, nil
			} else {
				return ctx, nil
			}
		}
		if rl.retryTTL <= 0 {
			return ctx, fmt.Errorf("get lock failed")
		}
		time.Sleep(rl.retryTTL)
	}
}

func (rl *RedisLock) Unlock(ctx context.Context) error {
	lockValue, err := getLockValue()
	if err != nil {
		return fmt.Errorf("get lock value error %v", err)
	}
	res, err := luaRelease.Run(ctx, rl.client, []string{rl.lockKey}, lockValue).Result()
	if err == redis.Nil {
		return fmt.Errorf("release error redislock: lock not held")
	} else if err != nil {
		return fmt.Errorf("release error %v", err)
	}

	if i, ok := res.(int64); !ok || i != 1 {
		return fmt.Errorf("release error res value type, ok is %v, i is %v", ok, i)
	}
	return nil
}

// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
func (rl *RedisLock) TTL(ctx context.Context) (time.Duration, error) {
	lockValue, err := getLockValue()
	if err != nil {
		return 0, fmt.Errorf("get lock value error %v", err)
	}
	res, err := luaPTTL.Run(ctx, rl.client, []string{rl.lockKey}, lockValue).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if num, ok := res.(int64); ok && num > 0 {
		return time.Duration(num) * time.Millisecond, nil
	}
	return 0, nil
}

func (rl *RedisLock) watchDog(ctx context.Context, cancelFunc context.CancelFunc, lockValue string) {
	expTicker := time.NewTicker(rl.expiration - time.Second*3)
	for {
		select {
		case <-expTicker.C:
			resp := luaRefresh.Run(ctx, rl.client, []string{rl.lockKey}, lockValue, rl.expiration.Seconds())
			if result, err := resp.Result(); err != nil || result == int64(0) {
				log.Infof("expire lock failed error %v, result %v", err, result)
				cancelFunc()
				return
			}
		case <-ctx.Done():
			log.Infof(" lock cancel")
		}
	}
}

func (rl *RedisLock) ClearForce(ctx context.Context) (int64, error) {
	return rl.client.Del(ctx, rl.lockKey).Result()
}

func (rl *RedisLock) GetLockIdAndTTL(ctx context.Context) (string, time.Duration, error) {
	lockid, err := rl.client.Get(ctx, rl.lockKey).Result()
	if err != nil {
		return "", 0, err
	}
	ttl, err := rl.client.TTL(ctx, rl.lockKey).Result()
	if err != nil {
		return "", 0, err
	}
	return lockid, ttl, nil
}
