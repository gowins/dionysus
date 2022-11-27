package main

import (
	"context"
	"fmt"
	"github.com/gowins/dionysus/memcache"
	"github.com/gowins/dionysus/step"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
)

var (
	green  = string([]byte{27, 91, 51, 50, 109})
	white  = string([]byte{27, 91, 51, 55, 109})
	yellow = string([]byte{27, 91, 51, 51, 109})
	red    = string([]byte{27, 91, 51, 49, 109})
	blue   = string([]byte{27, 91, 51, 52, 109})
	// magenta = string([]byte{27, 91, 51, 53, 109})
	// cyan    = string([]byte{27, 91, 51, 54, 109})
	// color   = []string{green, white, yellow, red, blue, magenta, cyan}
)

func main() {
	d := dionysus.NewDio()
	postSteps := []step.InstanceStep{
		{
			StepName: "PostPrint1", Func: func() error {
				fmt.Println(green, "=========== post 1 =========== ", white)
				return nil
			},
		},
		{
			StepName: "PostPrint2", Func: func() error {
				fmt.Println(green, "=========== post 2 =========== ", white)
				return nil
			},
		},
	}
	preSteps := []step.InstanceStep{
		{
			StepName: "PrePrint1", Func: func() error {
				fmt.Println(green, "=========== pre 1 =========== ", white)
				return nil
			},
		},
		{
			StepName: "PrePrint2", Func: func() error {
				fmt.Println(green, "=========== pre 2 =========== ", white)
				return nil
			},
		},
	}
	//NewBigCache()
	// PreRun exec before server start
	_ = d.PreRunStepsAppend(preSteps...)

	ctlCmd := cmd.NewCtlCommand()
	_ = ctlCmd.RegRunFunc(func() error {
		timer1 := time.NewTicker(time.Millisecond * 10)
		s := debug.GCStats{}
		m := runtime.MemStats{}
		var lastGC int64 = 0
		//go NewBigCache()
		for {
			select {
			case <-timer1.C:
				debug.ReadGCStats(&s)
				if s.NumGC != lastGC {
					lastGC = s.NumGC
					fmt.Printf("gc %d last@%v, PauseTotal %v\n", s.NumGC, s.LastGC, s.PauseTotal)
					runtime.ReadMemStats(&m)
					fmt.Printf("gc %d last@%v, next_heap_size@%vMB\n", m.NumGC, time.Unix(int64(time.Duration(m.LastGC).Seconds()), 0), m.NextGC/(1<<20))
				}
			case <-ctlCmd.Ctx.Done():
				fmt.Printf("this is stopChan %v\n", time.Now().String())
				return nil
			}
		}
	})

	ctx, cancel := context.WithCancel(ctlCmd.Ctx)
	ctlCmd.Ctx = ctx
	stopSteps := []cmd.StopStep{
		{
			StepName: "before stop",
			StopFn: func() {
				fmt.Printf("this is before stop\n")
			},
		},
		{
			StepName: "stop",
			StopFn: func() {
				fmt.Printf("this is stop\n")
				cancel()
			},
		},
	}
	ctlCmd.RegShutdownFunc(stopSteps...)

	// PostRun exec after server stop
	_ = d.PostRunStepsAppend(postSteps...)

	if err := d.DioStart("ctldemo", ctlCmd); err != nil {
		fmt.Printf("dio start error %v\n", err)
	}
}

var cacheName = "ctlMem"

func NewBigCache() {
	err := memcache.NewBigCache(context.Background(), cacheName, memcache.WithCleanWindow(time.Minute), memcache.WithLifeWindow(50*time.Second))
	if err != nil {
		fmt.Printf("new memory cache error %v\n", err)
	}
	runMemSet(20000*1000, 100)
	/*
		startTime := time.Now()
		data, err := memcache.Get(cacheName, "key9999999")
		fmt.Printf("spend time %v\n", time.Now().UnixMicro()-startTime.UnixMicro())
		if err != nil {
			fmt.Printf("memory cache Get error %v\n", err)
			return
		}
		fmt.Printf("data is %v\n", string(data))
	*/

	go runMemSetBig(10000*1000, 100)
	go runMemGet(10000*1000, 100)
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
				memcache.Set(cacheName, fmt.Sprintf("key%v", j), []byte(fmt.Sprintf("value%v", j)))
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
				//startTime := time.Now()
				memcache.Get(cacheName, fmt.Sprintf("key%v", j))
				//fmt.Printf("read spend time %v error %v\n", time.Now().UnixMicro()-startTime.UnixMicro(), err)
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
		start := page*i + 20000*1000
		end := page*i + page + 20000*1000
		go func() {
			for j := start; j < end; j++ {
				//startTime := time.Now()
				memcache.Set(cacheName, fmt.Sprintf("key%v", j), []byte(fmt.Sprintf("value%v", j)))
				//fmt.Printf("write spend time %v, error %v\n", time.Now().UnixMicro()-startTime.UnixMicro(), err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
