package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gowins/dionysus/memcache"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	// exit bigcache time.Ticker goroutine
	defer cancel()
	c, err := memcache.NewBigCache(ctx, memcache.WithCleanWindow(1*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	err = c.SetTTL("h1", []byte("t"), time.Millisecond*500)
	if err != nil {
		log.Fatal(err)
	}
	b, err := c.GetTTL("h1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	time.Sleep(time.Millisecond * 600)
	_, err = c.GetTTL("h1")
	fmt.Println(err == memcache.ErrEntryIsDead)
}
