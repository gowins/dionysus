package grpool

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type (
	TaskE        func(ctx context.Context) error
	DispatchType int
)

const (
	BreakByError = true
	NoBreak      = false

	SerialDispatch DispatchType = iota
	ConcurrentDispatch
)

var defaultTaskGroup = NewTaskGroup()

type TaskGroup struct{}

func NewTaskGroup() *TaskGroup {
	return &TaskGroup{}
}

// RunTask serial run all tasks, one after the other.
// if any task has return error, other tasks will not run
// return the first error.
func RunTask(ctx context.Context, tasks ...TaskE) error {
	return defaultTaskGroup.RunTask(defaultGrPool, ctx, tasks...)
}

func (*TaskGroup) RunTask(pool *GrPool, ctx context.Context, tasks ...TaskE) error {
	_, err := runTaskE(pool, ctx, SerialDispatch, BreakByError, tasks...)
	return err
}

// CRunTask concurrence run all tasks.
// if any task has return error, all tasks will be killed (Don't take it too seriously)
// return the first error.
func CRunTask(ctx context.Context, tasks ...TaskE) error {
	return defaultTaskGroup.CRunTask(defaultGrPool, ctx, tasks...)
}

func (*TaskGroup) CRunTask(pool *GrPool, ctx context.Context, tasks ...TaskE) error {
	_, err := runTaskE(pool, ctx, ConcurrentDispatch, BreakByError, tasks...)
	return err
}

// CRunTaskE concurrence run all tasks.
// if any task has return error, all tasks will not be killed (Don't take it too seriously)
// return all errors.
func CRunTaskE(ctx context.Context, tasks ...TaskE) (timeout bool, errs error) {
	return defaultTaskGroup.CRunTaskE(defaultGrPool, ctx, tasks...)
}

func (*TaskGroup) CRunTaskE(pool *GrPool, ctx context.Context, tasks ...TaskE) (timeout bool, errs error) {
	return runTaskE(pool, ctx, ConcurrentDispatch, NoBreak, tasks...)
}

func runTaskE(pool *GrPool, ctx context.Context, dtype DispatchType, isBreak bool, tasks ...TaskE) (timeout bool, errs error) {
	errCh := make(chan error, len(tasks))
	doneCh := make(chan struct{}, len(tasks))
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, task := range tasks {

		call := task
		t := func() {
			defer func() {
				if e := recover(); e != nil {
					errCh <- fmt.Errorf("got error: %+v", e)
				}
			}()
			if err := call(cctx); err != nil {
				errCh <- err
				return
			}
			doneCh <- struct{}{}
		}

		if err := pool.Submit(t); err != nil {
			if isBreak {
				return false, err
			}
			errs = ErrWrap(errs, err)
			continue
		}

		if dtype == ConcurrentDispatch {
			continue
		}

		// Serial Dispatch
		select {
		case <-cctx.Done():
			return true, ctx.Err()
		case err := <-errCh:
			return false, err
		case <-doneCh:
		}
	}

	if dtype == SerialDispatch {
		return
	}

	// Concurrent Dispatch
	for range tasks {
		select {
		case <-cctx.Done():
			if isBreak {
				errs = ctx.Err()
			}
			return true, errs
		case err := <-errCh:
			if isBreak {
				cancel()
				return false, err
			}
			errs = ErrWrap(errs, err)
		case <-doneCh:
		}
	}

	return timeout, errs
}

func ErrWrap(old, new error) error {
	if old == nil && new == nil {
		return nil
	}

	if old != nil && new != nil {
		return errors.Wrap(old, new.Error())
	}

	if old == nil {
		return new
	} else {
		return old
	}
}
