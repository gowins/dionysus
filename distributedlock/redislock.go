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
	client        *redis.Client
	lockKey       string
	expiration    time.Duration
	retryPeriod   time.Duration
	refreshPeriod time.Duration
	detailLog     bool
}

type Option func(redisLock *RedisLock)

func WithExpiration(expiration time.Duration) Option {
	return func(redisLock *RedisLock) {
		redisLock.expiration = expiration
	}
}

func WithRetryPeriod(retryPeriod time.Duration) Option {
	return func(redisLock *RedisLock) {
		redisLock.retryPeriod = retryPeriod
	}
}

func WithWatchDog(refreshPeriod time.Duration) Option {
	return func(redisLock *RedisLock) {
		redisLock.refreshPeriod = refreshPeriod
	}
}

func WithDetailLog(enable bool) Option {
	return func(redisLock *RedisLock) {
		redisLock.detailLog = enable
	}
}

func New(rclient *redis.Client, lockKey string, opts ...Option) *RedisLock {
	redisLock := &RedisLock{
		client:      rclient,
		lockKey:     lockKey,
		expiration:  defaultExpiration,
		retryPeriod: defaultRetryTTL,
	}
	for _, opt := range opts {
		opt(redisLock)
	}
	return redisLock
}

func (rl *RedisLock) Lock(ctx context.Context) (context.Context, error) {
	if rl.refreshPeriod >= rl.expiration && rl.refreshPeriod != 0 {
		return ctx, fmt.Errorf("refreshPeriod ")
	}

	lockValue, err := GetLockValue()
	if err != nil {
		return ctx, fmt.Errorf("get lock value error %v", err)
	}

	for {
		if rl.detailLog {
			log.Infof("try to get lock %v : %v time %v", rl.lockKey, lockValue, time.Now().String())
		}
		ok, err := rl.client.SetNX(ctx, rl.lockKey, lockValue, rl.expiration).Result()
		if err != nil {
			log.Errorf("set lock error %v", err)
		}
		if err == nil && ok {
			if rl.detailLog {
				log.Infof("get lock %v success %v time %v", rl.lockKey, lockValue, time.Now().String())
			}
			if rl.refreshPeriod > 0 {
				nctx, cancelFunc := context.WithCancel(ctx)
				go rl.watchDog(nctx, cancelFunc, lockValue)
				return nctx, nil
			} else {
				return ctx, nil
			}
		}
		if rl.retryPeriod <= 0 {
			return ctx, fmt.Errorf("get lock failed")
		}
		time.Sleep(rl.retryPeriod)
	}
}

func (rl *RedisLock) Unlock(ctx context.Context) error {
	lockValue, err := GetLockValue()
	if err != nil {
		return fmt.Errorf("get lock value error %v", err)
	}
	if rl.detailLog {
		log.Infof("try to release lock %v : %v time %v", rl.lockKey, lockValue, time.Now().String())
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
	lockValue, err := GetLockValue()
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
	expTicker := time.NewTicker(rl.refreshPeriod)
	for {
		select {
		case <-expTicker.C:
			resp := luaRefresh.Run(ctx, rl.client, []string{rl.lockKey}, lockValue, rl.expiration.Seconds())
			if result, err := resp.Result(); err != nil || result == int64(0) {
				if rl.detailLog {
					log.Infof("expire lock failed error %v, result %v, lockkey %v, lockid %v, time %v", err, result, rl.lockKey, lockValue, time.Now().String())
				}
				cancelFunc()
				return
			} else {
				if rl.detailLog {
					log.Infof("set expire lock success lock key %v, lockid %v, time %v ", rl.lockKey, lockValue, time.Now().String())
				}
			}
		case <-ctx.Done():
			if rl.detailLog {
				log.Infof("watchDog cancel lock key %v, lockid %v, time %v ", rl.lockKey, lockValue, time.Now().String())
			}
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
