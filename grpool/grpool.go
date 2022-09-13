package grpool

import (
	"time"

	"github.com/panjf2000/ants/v2"
)

var defaultGrPool, _ = NewPool(defaultAntsPoolSize)

const (
	// DefaultAntsPoolSize is the default capacity for a default goroutine pool.
	defaultAntsPoolSize = 20000

	// DefaultCleanIntervalTime is the interval time to clean up goroutines.
	defaultCleanIntervalTime = 2 * time.Second
)

type GrPool struct {
	p *ants.Pool
}

func (p *GrPool) Submit(task func()) error {
	return p.p.Submit(task)
}

func (p *GrPool) Running() int {
	return p.p.Running()
}

func (p *GrPool) Free() int {
	return p.p.Free()
}

func (p *GrPool) Cap() int {
	return p.p.Cap()
}

func (p *GrPool) Tune(size int) {
	p.p.Tune(size)
}

func (p *GrPool) Release() error {
	p.p.Release()
	return nil
}

func panicHandler(i interface{}) {
	log.Errorf("Task error: %v", i)
}

func NewPool(size int) (*GrPool, error) {
	return NewPoolWithExpire(size, defaultCleanIntervalTime)
}

func NewPoolWithExpire(size int, expiry time.Duration) (*GrPool, error) {
	p, err := ants.NewPool(size, ants.WithExpiryDuration(expiry), ants.WithPanicHandler(panicHandler))
	if err != nil {
		return nil, err
	}

	return &GrPool{
		p: p,
	}, nil
}

func Submit(task func()) error {
	return defaultGrPool.Submit(task)
}

func Running() int {
	return defaultGrPool.Running()
}

func Free() int {
	return defaultGrPool.Free()
}

func Cap() int {
	return defaultGrPool.Cap()
}

func Tune(size int) {
	defaultGrPool.Tune(size)
}

func Release() error {
	return defaultGrPool.Release()
}

// =============================================================================
// pool with function
type GrFuncPool struct {
	*ants.PoolWithFunc
}

// Pool with a argument
func NewPoolWithFunc(size int, pf func(interface{})) (*GrFuncPool, error) {
	return NewTimingPoolWithFunc(size, defaultCleanIntervalTime, pf)
}

func NewTimingPoolWithFunc(size int, expiry time.Duration, pf func(interface{})) (*GrFuncPool, error) {
	p, err := ants.NewPoolWithFunc(size, pf, ants.WithExpiryDuration(expiry), ants.WithPanicHandler(panicHandler))
	if err != nil {
		return nil, err
	}

	return &GrFuncPool{p}, nil
}
