package distributedlock

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
)

/*
func TestNew(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{})

	err := redisCli.Ping(context.Background()).Err()
	if err != nil {
		t.Errorf("ping redis error %v", err)
		return
	}
	rl := New(redisCli, WithExpiration(0))
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			lockValue, _ := getLockValue()
			fmt.Printf("this is %v, start get lock at %v\n", lockValue, time.Now().String())
			ctx, err := rl.Lock(context.Background())
			if err != nil {
				t.Errorf("lock error %v", err)
				return
			}
			fmt.Printf("this is %v, get lock at %v\n", lockValue, time.Now().String())
			defer wg.Done()
			tick := time.NewTicker(time.Second)
			count := 0
			for {
				select {
				case <-tick.C:
					le, err := rl.TTL(context.Background())
					fmt.Printf("ttt1 this is %v, get lock at %v\n", lockValue, time.Now().String())
					if count > 15 {
						err := rl.Unlock(ctx)
						fmt.Printf("unlock lockValue %v lock lost at %v, error is %v", lockValue, time.Now().String(), err)
						return
					}
					count++
					fmt.Printf("hold by %v, time %v, err %v\n", lockValue, le.String(), err)
				case <-ctx.Done():
					fmt.Printf("%v lock lost at time %v\n", lockValue, time.Now().String())
					return
				}
			}

		}()
	}
	wg.Wait()
}
*/

func TestNewMock(t *testing.T) {
	db, mock := redismock.NewClientMock()
	mock.MatchExpectationsInOrder(false)
	mock.ExpectGet("1234").SetVal("ddd")
	mock.ExpectSetNX("dsvsd", "wev", 0).SetVal(true)
	mock.ExpectEvalSha("cf0e94b2e9ffc7e04395cf88f7583fc309985910", []string{defaultLockKey}, "lockValue").SetVal(1)
	mock.ExpectEvalSha("cf0e94b2e9ffc7e04395cf88f7583fc309985910", []string{defaultLockKey}, "lockValue").SetErr(nil)

	res, err := luaRelease.Run(context.Background(), db, []string{defaultLockKey}, "lockValue").Result()

	mock.ClearExpect()
	mock.ExpectEvalSha("cf0e94b2e9ffc7e04395cf88f7583fc309985910", []string{defaultLockKey}, "lockValue").SetVal(1)
	mock.ExpectEvalSha("cf0e94b2e9ffc7e04395cf88f7583fc309985910", []string{defaultLockKey}, "lockValue").SetErr(nil)
	//db.Set(context.Background(), "1234", "ddw", 0)
	fmt.Printf("get redis value %v, err %v\n", res, err)
	res, err = luaRelease.Run(context.Background(), db, []string{defaultLockKey}, "lockValue").Result()
	//db.Set(context.Background(), "1234", "ddw", 0)
	fmt.Printf("11get redis value %v, err %v\n", res, err)
}

func TestNew(t *testing.T) {
	lockKey := "testLockKet"
	testExpiration := 11 * time.Second
	testRetryTTL := 2 * time.Second
	db, _ := redismock.NewClientMock()
	rlock := New(db, lockKey, WithExpiration(testExpiration), WithRetryTTL(testRetryTTL), WithWatchDog(false))
	if rlock.watchDogEnable != false {
		t.Errorf("want watchDogEnable false get true")
		return
	}
	if rlock.expiration != testExpiration {
		t.Errorf("want expiration %v, get expiration %v", testExpiration, rlock.expiration)
		return
	}
	if rlock.retryTTL != testRetryTTL {
		t.Errorf("want retryTTL %v, get retryTTL %v", testRetryTTL, rlock.retryTTL)
		return
	}
}
