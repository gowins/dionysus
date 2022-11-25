package memcache

/*  benchmark test 1000万数据下，100并发写，100并发读，平均延时1微秒延时，最高延时30微妙
func TestNewBigCache(t *testing.T) {
	memCache, err := NewBigCache(context.Background())
	if err != nil {
		fmt.Printf("new memory cache error %v\n", err)
	}
	fmt.Printf("before cap %v, len %v\n", memCache.Capacity(), memCache.Len())
	runMemSet(10000*1000, 100, memCache)
	fmt.Printf("after cap %v, len %v\n", memCache.Capacity(), memCache.Len())
	startTime := time.Now()
	data, err := memCache.Get("key9999999")
	fmt.Printf("spend time %v\n", time.Now().UnixMicro()-startTime.UnixMicro())
	if err != nil {
		fmt.Printf("memory cache Get error %v\n", err)
		return
	}
	fmt.Printf("data is %v\n", string(data))

	go runMemSetBig(10000*1000, 100, memCache)
	go runMemGet(10000*1000, 100, memCache)
	select {}
}

func runMemSet(dataTotal int, job int, memCache *bigCache) {
	var wg sync.WaitGroup
	wg.Add(job)
	page := dataTotal / job
	for i := 0; i < job; i++ {
		start := page * i
		end := page*i + page
		go func() {
			for j := start; j < end; j++ {
				err := memCache.Set(fmt.Sprintf("key%v", j), []byte(fmt.Sprintf("value%v", j)))
				if err != nil {
					fmt.Printf("memCache set error")
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func runMemGet(dataTotal int, job int, memCache *bigCache) {
	var wg sync.WaitGroup
	wg.Add(job)
	page := dataTotal / job
	for i := 0; i < job; i++ {
		start := page * i
		end := page*i + page
		go func() {
			for j := start; j < end; j++ {
				startTime := time.Now()
				_, err := memCache.Get(fmt.Sprintf("key%v", j))
				if err != nil {
					fmt.Printf("memCache get error")
				}
				fmt.Printf("read spend time %v\n", time.Now().UnixMicro()-startTime.UnixMicro())
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func runMemSetBig(dataTotal int, job int, memCache *bigCache) {
	var wg sync.WaitGroup
	wg.Add(job)
	page := dataTotal / job
	for i := 0; i < job; i++ {
		start := page*i + 10000*1000
		end := page*i + page + 10000*1000
		go func() {
			for j := start; j < end; j++ {
				startTime := time.Now()
				err := memCache.Set(fmt.Sprintf("key%v", j), []byte(fmt.Sprintf("value%v", j)))
				if err != nil {
					fmt.Printf("memCache set error")
				}
				fmt.Printf("write spend time %v\n", time.Now().UnixMicro()-startTime.UnixMicro())
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
*/
