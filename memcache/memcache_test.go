package memcache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/smartystreets/goconvey/convey"
)

type testLogger struct{}

func (t *testLogger) Printf(msg string, args ...any) {}

func TestBigCache(t *testing.T) {
	convey.Convey("big cache", t, func() {
		convey.Convey("new bigcache error", func() {
			opts := []ConfigOpt{WithShards(10)}
			err := NewBigCache(context.Background(), "testCache1", opts...)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("new bigcache", func() {
			err := newBigCache()
			convey.So(err, convey.ShouldBeNil)
			err = newBigCache()
			fmt.Printf("newBigCache error is %v\n", err)
			convey.So(err, convey.ShouldNotBeNil)
			err = Set("testCache2", "h3", []byte("2"))
			convey.So(err, convey.ShouldBeNil)
			b, err := Get("testCache2", "h3")
			convey.So(err, convey.ShouldBeNil)
			convey.So(string(b), convey.ShouldEqual, "2")

			err = Delete("testCache2", "h3")
			convey.So(err, convey.ShouldBeNil)
			_, err = Get("testCache2", "h3")
			convey.So(err, convey.ShouldNotBeNil)
			_, ok := GetCache("testCache2")
			convey.So(ok, convey.ShouldBeTrue)
			err = Close("testCache2")
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("init bigcache", func() {
			cacheInit := InitBigCache("cachetest1")
			convey.So(cacheInit, convey.ShouldNotBeNil)
			err := cacheInit.Set("h5", []byte("3"))
			convey.So(err, convey.ShouldBeNil)
			b, err := cacheInit.Get("h5")
			convey.So(err, convey.ShouldBeNil)
			convey.So(string(b), convey.ShouldEqual, "3")
		})
	})
}

func newBigCache() error {
	opts := []ConfigOpt{
		WithShards(1024),
		WithLifeWindow(10 * time.Second),
		WithCleanWindow(1 * time.Second),
		WithMaxEntriesInWindow(1000 * 10 * 60),
		WithMaxEntrySize(500),
		WithStatsEnabled(false),
		WithVerbose(true),
		WithHardMaxCacheSize(2),
		WithOnRemove(func(s string, b []byte) {}),
		WithOnRemoveWithMetadata(func(s string, b []byte, m bigcache.Metadata) {}),
		WithOnRemoveWithReason(func(s string, b []byte, rr bigcache.RemoveReason) {}),
		WithLogger(&testLogger{}),
		WithHasher(&Hasher{}),
	}
	return NewBigCache(context.Background(), "testCache2", opts...)
}
