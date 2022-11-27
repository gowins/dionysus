### 配置

```go
type ConfigOpt func(*bigcache.Config)
```

* Shards int

分片数量，必须是2的幂次方， 例如2<sup>1</sup>,2<sup>2</sup>....；

```go
// 即shards切片长度
shards []*cacheShard
```

* LifeWindow time.Duration
  * 整体缓存过期时间；
  * 若CleanWindow大于0，每隔CleanWindow时间间隔，过期的缓存会被自动清理；
  * 否则，bigcache在每次设置缓存时，会判断最早的缓存是否过期，若过期，则清理；


* CleanWindow time.Duration

若大于0，则会定期自动删除缓存时间超过LifeWindow

```go
	if config.CleanWindow > 0 {
		go func() {
			ticker := time.NewTicker(config.CleanWindow)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					fmt.Println("ctx done, shutting down bigcache cleanup routine")
					return
				case t := <-ticker.C:
					cache.cleanUp(uint64(t.Unix()))
				case <-cache.close:
					return
				}
			}
		}()
	}

```

* MaxEntriesInWindow int

必须大于0，用于计算每个shard的初始大小

```go
// initialShardSize computes initial shard size
max(MaxEntriesInWindow/Shards, 10)
```

* MaxEntrySize int

每个实体最大的内存大小，单位byte，必须大于0

```go
// queue.ByteQueue.capacity大小，若此值大于HardMaxCacheSize，则取HardMaxCacheSize
max(MaxEntriesInWindow/Shards, 10) * MaxEntrySize
```

* StatsEnabled bool

若开启，各个cacheShard分片会统计每个key的命中次数

* Verbose bool

若开启，会打印新内存分配的信息

* Hasher bigchache.Hasher

为key生成无符号的64位整数哈希值

```go
// bigchache
type Hasher interface {
  Sum64(key string) uint64
}
```

* HardMaxCacheSize int

每个分片，缓存（queue.ByteQueue）最大内存大小，单位MB，**若为0，则表示没有限制**

### 使用

```go
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
```

>由于bigcache有一个整体的过期时间和清理策略，通过CleanWindow来控制清理时间间隔，以及LifeWindow来控制整体缓存的有效时间；
>
>在使用时，注意LifeWindow与这里设置的过期时间，建议把LifeWindow设置比TTL大一些，避免被bigcache自动清理，例如1分钟
>
>需要调用c.Close()或者通过取消上下文，让bigcache定时器goroutine退出

## Benchmarks
1:  50秒有效期，1分钟清理窗口配置下，1000万数据下，100并发写，100并发读，平均延时1微秒延时，最高延时30微妙。  
bench test
```
goos: darwin
goarch: arm64
pkg: github.com/gowins/dionysus/memcache
BenchmarkBigCacheSet-8           	17409649	       382.0 ns/op	     152 B/op	       2 allocs/op
BenchmarkBigCacheGet-8           	16806682	       278.9 ns/op	      85 B/op	       2 allocs/op
BenchmarkBigCacheSetParallel-8   	41394356	       134.9 ns/op	     117 B/op	       3 allocs/op
BenchmarkBigCacheGetParallel-8   	44905400	        92.14 ns/op	      23 B/op	       1 allocs/op
PASS
ok  	github.com/gowins/dionysus/memcache	49.881s
```

2:  GC pause time
```shell
go version
go version go1.19 linux/amd64

go run caches_gc_overhead_comparison.go

Number of entries:  10000000
GC pause for bigcache:  0.306077ms
Number of entries:  20000000
GC pause for bigcache:  0.506077ms
```