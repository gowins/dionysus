package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gowins/dionysus/memcache"
)

func main() {
	cacheName := "cacheDemo"
	err := memcache.NewBigCache(context.Background(), cacheName, memcache.WithCleanWindow(time.Minute), memcache.WithLifeWindow(50*time.Second))
	if err != nil {
		fmt.Printf("new memory cache error %v\n", err)
	}
	memcache.Set(cacheName, "key9999999", []byte("key9999999"))
	startTime := time.Now()
	data, err := memcache.Get(cacheName, "key9999999")
	fmt.Printf("spend time %v\n", time.Now().UnixMicro()-startTime.UnixMicro())
	if err != nil {
		fmt.Printf("memory cache Get error %v\n", err)
		return
	}
	fmt.Printf("data is %v\n", string(data))
}
