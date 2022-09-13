package grpool

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	C "github.com/smartystreets/goconvey/convey"
)

func TestRunTask(t *testing.T) {
	defaultGrPool, _ = NewPool(defaultAntsPoolSize)
	var (
		Err  = errors.New("just for testing")
		Err2 = errors.New("just for testing {2}")
	)

	C.Convey("[TestRunTask]", t, func() {
		C.Convey("[TestRunTask] First Err", func() {
			var tasks []TaskE
			tasks = append(tasks, func(ctx context.Context) error {
				return Err
			})
			tasks = append(tasks, func(ctx context.Context) error {
				panic("1")
			})
			err := RunTask(context.Background(), tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(err.Error(), C.ShouldEqual, Err.Error())
		})

		C.Convey("[TestRunTask] Deadline exceeded", func() {
			var tasks []TaskE
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			tasks = make([]TaskE, 0)
			tasks = append(tasks, func(ctx context.Context) error {
				return nil
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return Err2
			})
			err := RunTask(ctx, tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(err.Error(), C.ShouldEqual, context.DeadlineExceeded.Error())
		})

		C.Convey("[TestRunTask] Recover", func() {
			var tasks []TaskE
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			tasks = make([]TaskE, 0)
			tasks = append(tasks, func(ctx context.Context) error {
				panic(1)
				return nil
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return Err2
			})
			err := RunTask(ctx, tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(err.Error(), C.ShouldStartWith, "got error:")
		})
	})
}

func TestCRunTask(t *testing.T) {
	defaultGrPool, _ = NewPool(defaultAntsPoolSize)

	var (
		Err  = errors.New("just for testing")
		Err2 = errors.New("just for testing {2}")
	)

	C.Convey("[TestCRunTask]", t, func() {
		C.Convey("[TestCRunTask] First Err", func() {
			var tasks []TaskE
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second)
				return Err
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return nil
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				panic("2")
			})
			err := CRunTask(context.Background(), tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(err.Error(), C.ShouldContainSubstring, Err.Error())
			C.So(err.Error(), C.ShouldNotContainSubstring, "got error:")
		})

		C.Convey("[TestCRunTask] Recover", func() {
			var tasks []TaskE
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return Err
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				panic("2")
			})
			err := CRunTask(context.Background(), tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(err.Error(), C.ShouldNotContainSubstring, Err.Error())
			C.So(err.Error(), C.ShouldContainSubstring, "got error:")
		})

		C.Convey("[TestCRunTask] Deadline exceeded", func() {
			var tasks []TaskE
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			tasks = make([]TaskE, 0)
			tasks = append(tasks, func(ctx context.Context) error {
				return nil
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return Err2
			})
			err := CRunTask(ctx, tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(err.Error(), C.ShouldEqual, context.DeadlineExceeded.Error())
			C.So(err.Error(), C.ShouldNotContainSubstring, Err2.Error())
		})
	})
}

func TestCRunTaskE(t *testing.T) {
	defaultGrPool, _ = NewPool(defaultAntsPoolSize)

	var (
		Err  = errors.New("just for testing")
		Err2 = errors.New("just for testing {2}")
	)

	C.Convey("[TestCRunTaskE]", t, func() {
		C.Convey("[TestCRunTaskE] First Err", func() {
			var tasks []TaskE
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second)
				return Err
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return nil
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				panic("2")
			})
			timeout, err := CRunTaskE(context.Background(), tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(timeout, C.ShouldBeFalse)
			C.So(err.Error(), C.ShouldContainSubstring, Err.Error())
			C.So(err.Error(), C.ShouldContainSubstring, "got error:")
		})

		C.Convey("[TestCRunTaskE] Recover", func() {
			var tasks []TaskE
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 3)
				return Err
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				panic("2")
			})
			timeout, err := CRunTaskE(context.Background(), tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(timeout, C.ShouldBeFalse)
			C.So(err.Error(), C.ShouldContainSubstring, Err.Error())
			C.So(err.Error(), C.ShouldContainSubstring, "got error:")
		})

		C.Convey("[TestCRunTask] Deadline exceeded", func() {
			var tasks []TaskE
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			tasks = make([]TaskE, 0)
			tasks = append(tasks, func(ctx context.Context) error {
				return Err
			})
			tasks = append(tasks, func(ctx context.Context) error {
				time.Sleep(time.Second * 2)
				return Err2
			})
			timeout, err := CRunTaskE(ctx, tasks...)
			C.So(err, C.ShouldNotBeNil)
			C.So(timeout, C.ShouldBeTrue)
			C.So(err.Error(), C.ShouldNotContainSubstring, context.DeadlineExceeded.Error())
			C.So(err.Error(), C.ShouldNotContainSubstring, Err2.Error())
			C.So(err.Error(), C.ShouldContainSubstring, Err.Error())
		})
	})
}
