package memcache

import (
	"context"
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
			_, err := NewBigCache(context.Background(), opts...)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("new bigcache", func() {
			c, err := newBigCache()
			convey.So(err, convey.ShouldBeNil)
			_, err = c.GetTTL("h1")
			convey.So(err, convey.ShouldNotBeNil)
			err = c.SetTTL("h1", []byte("t"), time.Millisecond*0)
			convey.So(err, convey.ShouldNotBeNil)
			err = c.SetTTL("h1", []byte("t"), time.Millisecond*50)
			convey.So(err, convey.ShouldBeNil)
			b, _ := c.GetTTL("h1")
			convey.So(string(b), convey.ShouldEqual, "t")
			time.Sleep(time.Millisecond * 60)
			_, err = c.GetTTL("h1")
			convey.So(err, convey.ShouldEqual, ErrEntryIsDead)

			err = c.Delete("h2")
			convey.So(err, convey.ShouldNotBeNil)
			err = c.SetTTL("h2", []byte("1"), 20*time.Second)
			convey.So(err, convey.ShouldBeNil)
			err = c.Delete("h2")
			convey.So(err, convey.ShouldBeNil)
			_, err = c.GetTTL("h2")
			convey.So(err, convey.ShouldNotBeNil)

			err = c.Set("h3", []byte("2"))
			convey.So(err, convey.ShouldBeNil)
			b, err = c.Get("h3")
			convey.So(err, convey.ShouldBeNil)
			convey.So(string(b), convey.ShouldEqual, "2")

			err = c.Close()
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

func newBigCache() (*bigCache, error) {
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
		WithHasher(&hasher{}),
	}
	return NewBigCache(context.Background(), opts...)
}
