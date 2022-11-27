package memcache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func value(i int) []byte {
	value := fmt.Sprintf("value-%010d", i)
	return []byte(value)
}

func BenchmarkBigCacheSet(b *testing.B) {
	cacheName := "CacheSet"
	err := initBigCache(b.N, cacheName)
	if err != nil {
		b.Errorf("want new big cache error nil, get error %v", err)
		return
	}
	for i := 0; i < b.N; i++ {
		err = Set(cacheName, key(i), value(i))
		if err != nil {
			b.Errorf("want new big cache error nil, get error %v", err)
			return
		}
	}
}

func BenchmarkBigCacheGet(b *testing.B) {
	b.StopTimer()
	cacheName := "CacheSet"
	err := initBigCache(b.N, cacheName)
	if err != nil {
		b.Errorf("want new big cache error nil, get error %v", err)
		return
	}
	for i := 0; i < b.N; i++ {
		err = Set(cacheName, key(i), value(i))
		if err != nil {
			b.Errorf("want new big cache error nil, get error %v", err)
			return
		}
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Get(cacheName, key(i))
	}
}

func BenchmarkBigCacheGetParallel(b *testing.B) {
	cacheName := "GetParallel"
	b.StopTimer()
	err := initBigCache(b.N, cacheName)
	if err != nil {
		b.Errorf("want new big cache error nil, get error %v", err)
		return
	}
	for i := 0; i < b.N; i++ {
		err = Set(cacheName, key(i), value(i))
		if err != nil {
			b.Errorf("want new big cache error nil, get error %v", err)
			return
		}
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			Get(cacheName, key(counter))
			counter = counter + 1
		}
	})
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}

func BenchmarkBigCacheSetParallel(b *testing.B) {
	cacheName := "SetParallel"
	b.StopTimer()
	err := initBigCache(b.N, cacheName)
	if err != nil {
		b.Errorf("want new big cache error nil, get error %v", err)
		return
	}

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			Set(cacheName, parallelKey(id, counter), value(id))
			counter = counter + 1
		}
	})
}

func initBigCache(entriesInWindow int, cacheName string) error {
	return NewBigCache(context.Background(), cacheName, WithCleanWindow(0), WithMaxEntriesInWindow(entriesInWindow),
		WithLifeWindow(10*time.Minute), WithVerbose(true), WithShards(256), WithMaxEntrySize(256))
}
