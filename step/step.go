package step

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/gowins/dionysus/algs"
)

type Steps struct {
	q           *algs.PriorityQueue // 全局启动依赖项
	m           sync.Map
	appendIndex int
}

// SystemPrioritySteps system Priority 0-100
const SystemPrioritySteps = 100

// UserPrioritySteps user Priority 100--10100
const UserPrioritySteps = 10100

// UserAppendPrioritySteps user Append steps 10101-无穷大
const UserAppendPrioritySteps = 10101

type FuncStep func() error

func New() *Steps {
	return &Steps{
		q:           algs.NewPQ(),
		m:           sync.Map{},
		appendIndex: UserAppendPrioritySteps,
	}
}

func (s *Steps) RegFirstSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 1, fn)
}

func (s *Steps) RegSecondSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 2, fn)
}

func (s *Steps) RegThirdSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 3, fn)
}

func (s *Steps) RegFourthSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 4, fn)
}

func (s *Steps) RegFifthSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 5, fn)
}

func (s *Steps) RegSixthSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 6, fn)
}

func (s *Steps) RegSeventhSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 7, fn)
}

func (s *Steps) RegEighthSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 8, fn)
}

func (s *Steps) RegNinethSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 9, fn)
}

func (s *Steps) RegTenthSteps(value string, fn FuncStep) {
	s.RegActionSteps(value, 10, fn)
}

func (s *Steps) RegActionSteps(value string, priority int, fn FuncStep) {
	item := algs.NewItem(value, priority)
	s.m.Store(item, fn)
	s.q.Push(item)
}

func (s *Steps) RegActionStepsE(value string, priority int, fn FuncStep) error {
	if priority < 0 {
		return fmt.Errorf(" Priority can not be negtive: %d ", priority)
	}

	if fn == nil {
		return fmt.Errorf(" Function can not be nil: %T ", fn)
	}

	s.RegActionSteps(value, priority, fn)
	return nil
}

func (s *Steps) ActionStepsAppend(value string, fn FuncStep) error {
	if fn == nil {
		return fmt.Errorf(" Function can not be nil: %T ", fn)
	}
	s.RegActionSteps(value, s.appendIndex, fn)
	s.appendIndex++
	return nil
}

// 初始化加载router middle afterstart等等
func (s *Steps) Run() error {
	// Take the items out; they arrive in decreasing priority order.
	i := 1
	pqLen := s.q.Len()
	for s.q.Len() > 0 {
		item, _ := s.q.Pop()
		if fn, ok := s.m.Load(item); ok && !reflect.ValueOf(fn).IsNil() {
			if f, ok := fn.(FuncStep); ok {
				if err := f(); err != nil {
					ef := fmt.Errorf("[step %d/%d] %s err: %v", i, pqLen, item.Value(), err)
					log.Print(ef)
					return ef
				} else {
					log.Printf("[step %d/%d] %s success", i, pqLen, item.Value())
				}
			}
		} else {
			log.Printf("[warn] load step false %v \n", item)
		}
		i++
	}
	return nil
}
