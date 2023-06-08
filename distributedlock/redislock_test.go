package distributedlock

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
