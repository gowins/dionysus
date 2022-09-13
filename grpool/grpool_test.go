package grpool

import (
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	panicStr    = "This is panic"
	gotAnyPanic = false
)

func testPanicHandler(i interface{}) {
	gotAnyPanic = true
	fmt.Printf("%v", i)
}

func panicTask() {
	panic(panicStr)
}

func noPanicTask() {
	time.Sleep(time.Microsecond * 10)
}

func TestPool(t *testing.T) {
	task := func() {
		fmt.Println("task s")
	}

	Convey("test pool", t, func() {
		nNum := 10
		pool, err := NewPool(nNum)
		So(err, ShouldBeNil)

		_ = pool.Submit(task)

		rNum := pool.Running()
		So(rNum, ShouldEqual, 1)

		fNum := pool.Free()
		So(fNum, ShouldEqual, nNum-rNum)

		err = pool.Release()
		So(err, ShouldBeNil)
	})
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		task      func()
		wantPanic bool
		handler   func(interface{})
		size      int
		wantError bool
	}{
		{
			name:      "Default Panic Handler",
			task:      panicTask,
			wantPanic: true,
			size:      10,
		},
		{
			name:      "No Panic",
			task:      noPanicTask,
			wantPanic: false,
			handler:   testPanicHandler,
			size:      10,
		},
		{
			name:      "Default Panic Handler",
			task:      noPanicTask,
			wantPanic: false,
			size:      10,
		},
		{
			name:      "Invalid size",
			task:      noPanicTask,
			wantError: true,
			size:      0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPool(tt.size)
			if err != nil {
				if tt.wantError {
					return
				}
				t.Error(err)
			}

			if got == nil {
				t.Errorf("New() = %v, err = %v", got, err)
			}

			var wg sync.WaitGroup
			f := func() {
				defer func() {
					wg.Done()
				}()
				tt.task()
			}

			wg.Add(1)
			gotAnyPanic = false
			if err := got.Submit(f); err != nil {
				t.Errorf("Submit error: %v", err)
			}
			wg.Wait()

			// If use the test handler. We need to check the panic has been recovered or not.
			if tt.handler != nil && tt.wantPanic != gotAnyPanic {
				t.Error("panic not be recover.")
			}

			_ = got.Release()
		})
	}
}

func TestSubmit(t *testing.T) {
	type args struct {
		task func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				task: panicTask,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Submit(tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("Submit() error = %v, wantErr %v", err, tt.wantErr)
			}

			if r := Running(); r > defaultAntsPoolSize || r < 0 {
				t.Error("Invalid size.")
			}

			if r := Free(); r > defaultAntsPoolSize || r < 0 {
				t.Error("Invalid size.")
			}

			if r := Cap(); r > defaultAntsPoolSize || r < 0 {
				t.Error("Invalid size.")
			}

			Tune(10)
			if r := Cap(); r > 10 || r < 0 {
				t.Error("Invalid size.")
			}

			_ = Release()
		})
	}
}

func TestForMap(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("cover 11")
		}
	}()

	ids := make([]int64, 0)
	for i := 1; i <= 1000; i++ {
		ids = append(ids, int64(i))
	}
	var bindIds = make(map[int64]int64)
	var wg sync.WaitGroup
	var mtx sync.Mutex
	for _, id := range ids {
		wg.Add(1)
		tmpID := id
		task := func() {
			defer func() {
				wg.Done()
			}()

			mtx.Lock() // need lock
			bindIds[tmpID] = tmpID
			mtx.Unlock()
		}

		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Println("cover 22")
				}
			}()

			task()
		}()
	}

	wg.Wait()
}

func TestForSyncMap(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("cover 11")
		}
	}()

	ids := make([]int64, 0)
	for i := 1; i <= 1000; i++ {
		ids = append(ids, int64(i))
	}
	var bindIds sync.Map
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		tmpID := id
		task := func() {
			defer func() {
				wg.Done()
			}()
			bindIds.Store(tmpID, tmpID)
		}

		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Println("cover 22")
				}
			}()

			task()
		}()
	}

	wg.Wait()

	var max int64
	bindIds.Range(func(key, value interface{}) bool {
		uid, ok1 := key.(int64)
		_, ok2 := value.(int64)
		if !ok1 || !ok2 {
			fmt.Println("error")
		}

		if max < uid {
			max = uid
		}

		return true
	})

	fmt.Println(max)
}
