package memcache

import (
	"context"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkBigCacheSet(b *testing.B) {
	initBigCache(b.N, "CacheSet")
	for i := 0; i < b.N; i++ {
		Set("CacheSet", key(i), value())
	}
}

func BenchmarkBigCacheGet(b *testing.B) {
	b.StopTimer()
	initBigCache(b.N, "CacheGet")
	for i := 0; i < b.N; i++ {
		Set("CacheGet", key(i), value())
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Get("CacheGet", key(i))
	}
}

func BenchmarkBigCacheSetParallel(b *testing.B) {
	initBigCache(b.N, "SetParallel")
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			Set("SetParallel", parallelKey(id, counter), value())
			counter = counter + 1
		}
	})
}

func BenchmarkBigCacheGetParallel(b *testing.B) {
	b.StopTimer()
	initBigCache(b.N, "GetParallel")
	for i := 0; i < b.N; i++ {
		Set("GetParallel", key(i), value())
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			Get("GetParallel", key(counter))
			counter = counter + 1
		}
	})
}

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}

func initBigCache(entriesInWindow int, name string) *bigcache.BigCache {
	NewBigCache(context.Background(), name, WithMaxEntriesInWindow(entriesInWindow), WithVerbose(false))
	cache, _ := GetCache(name)
	return cache
}
