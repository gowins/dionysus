package memcache

/*
var cacheName = "benchCache"

// benchmark test 1000万数据下，100并发写，100并发读，平均延时1微秒延时，最高延时30微妙, 50秒有效期，1分钟清理窗口
func TestNewBigCache(t *testing.T) {
	err := NewBigCache(context.Background(), cacheName, WithCleanWindow(time.Minute), WithLifeWindow(50*time.Second))
	if err != nil {
		fmt.Printf("new memory cache error %v\n", err)
	}
	runMemSet(10000*1000, 100)
	startTime := time.Now()
	data, err := Get(cacheName, "key9999999")
	fmt.Printf("spend time %v\n", time.Now().UnixMicro()-startTime.UnixMicro())
	if err != nil {
		fmt.Printf("memory cache Get error %v\n", err)
		return
	}
	fmt.Printf("data is %v\n", string(data))

	go runMemSetBig(10000*1000, 100)
	go runMemGet(10000*1000, 100)
	select {}
}

func runMemSet(dataTotal int, job int) {
	var wg sync.WaitGroup
	wg.Add(job)
	page := dataTotal / job
	for i := 0; i < job; i++ {
		start := page * i
		end := page*i + page
		go func() {
			for j := start; j < end; j++ {
				err := Set(cacheName, fmt.Sprintf("key%v", j), []byte(fmt.Sprintf("value%v", j)))
				if err != nil {
					fmt.Printf("memCache set error %v\n", err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func runMemGet(dataTotal int, job int) {
	var wg sync.WaitGroup
	wg.Add(job)
	page := dataTotal / job
	for i := 0; i < job; i++ {
		start := page * i
		end := page*i + page
		go func() {
			for j := start; j < end; j++ {
				startTime := time.Now()
				_, err := Get(cacheName, fmt.Sprintf("key%v", j))
				fmt.Printf("read spend time %v error %v\n", time.Now().UnixMicro()-startTime.UnixMicro(), err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func runMemSetBig(dataTotal int, job int) {
	var wg sync.WaitGroup
	wg.Add(job)
	page := dataTotal / job
	for i := 0; i < job; i++ {
		start := page*i + 10000*1000
		end := page*i + page + 10000*1000
		go func() {
			for j := start; j < end; j++ {
				startTime := time.Now()
				err := Set(cacheName, fmt.Sprintf("key%v", j), []byte(fmt.Sprintf("value%v", j)))
				fmt.Printf("write spend time %v, error %v\n", time.Now().UnixMicro()-startTime.UnixMicro(), err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
*/
