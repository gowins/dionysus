package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	dl "github.com/gowins/dionysus/distributedlock"
	"github.com/gowins/dionysus/log"
	"time"
)

func main() {
	log.Setup()

	redisCli := redis.NewClient(&redis.Options{})

	err := redisCli.Ping(context.Background()).Err()
	if err != nil {
		log.Errorf("connect redis client failed %w", err)
		return
	}

	// lock will expiration after 10s automatic
	rlock := dl.New(redisCli, "diodemolockkey")

	re, err := rlock.ClearForce(context.Background())
	if re >= 0 {
		log.Infof("re is %v, err is %v", re, err)
	}

	_, err = rlock.Lock(context.Background())
	if err != nil {
		log.Errorf("dio demo lock error %v", err)
		return
	}
	// do something
	for i := 0; i < 20; i++ {
		id, ttl, err := rlock.GetLockIdAndTTL(context.Background())
		log.Infof("id %v time left ttl %v, error %v", id, ttl.String(), err)
		time.Sleep(time.Second)
	}
	err = rlock.Unlock(context.Background())
	if err != nil {
		log.Errorf("redis unlock error %v", err)
	}
}
