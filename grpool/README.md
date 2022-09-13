# pkg说明

### grpool

在gapi框架中不建议使用go goroutine开go协程，因为如果go协程如果没有recover()的话，程序panic会无法捕获。
要求使用此pkg下的```grpool.Submit(task func())```来运行go协程，此方法自带了recover()，可以捕获panic。


### TaskGroup
串行或并发地执行一组任务
```go
type TaskE = func(ctx context.Context) error

// RunTask serial run all tasks, one after the other.
// if any task has return error, other tasks will not run
// return the error.
func RunTask(ctx context.Context, tasks ...TaskE) error 

// CRunTask concurrence run all tasks.
// if any task has return error, all tasks will be killed (Don't take it too seriously)
// return the first error.
func CRunTask(ctx context.Context, tasks ...TaskE) error 

// CRunTaskE concurrence run all tasks.
// if any task has return error, all tasks will not be killed (Don't take it too seriously)
// return all errors.
func CRunTaskE(ctx context.Context, tasks ...TaskE) (timeout bool, errs error) 
```

```
`CRunTask`
所有任务将被并发执行（无序），如果任何一个任务`发生错误`或`ctx超时` ，`CRunTaskE` 将返回第一个错误。

`CRunTaskE`
所有任务将被并发执行（无序），任何一个任务`发生错误`都不会影响其他任务，并收集所有任务的error。
但如果是ctx超时，`CRunTaskE`将不再等待未完成的任务直接返回。

`RunTask`
所有任务按照`FIFO`被串行执行，如果任何一个任务在执行过程中`发生错误`或`ctx超时`，则剩余任务将不再执行，返回当前错误。
```
