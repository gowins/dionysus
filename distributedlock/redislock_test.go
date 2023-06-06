package distributedlock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"testing"
	"time"
)

func Test_redisLock_Lock(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "120.27.236.102:6379",
		DB:       0,
		Password: "123456",
	})

	err := redisCli.Ping(context.Background()).Err()
	if err != nil {
		t.Errorf("ping redis error %v", err)
		return
	}
	rl := RedisLock{client: redisCli}
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			rl.Lock()
			defer wg.Done()
			defer rl.Unlock()
			for j := 0; j < 16; j++ {
				ttl, err := rl.TTL()
				va, _ := getLockValue()
				fmt.Printf("err is %v, va is %v ttl is %v, time is %v\n", err, va, ttl, time.Now().String())
				time.Sleep(time.Second)
			}
		}()
	}
	wg.Wait()
}
