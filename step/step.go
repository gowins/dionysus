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

type InstanceStep struct {
	StepName string
	Func     FuncStep
}

func New() *Steps {
	return &Steps{
		q:           algs.NewPQ(),
		m:           sync.Map{},
		appendIndex: UserAppendPrioritySteps,
	}
}

func (s *Steps) RegSysFirstSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 1, instanceStep.Func)
}

func (s *Steps) RegSysSecondSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 2, instanceStep.Func)
}

func (s *Steps) RegSysThirdSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 3, instanceStep.Func)
}

func (s *Steps) RegSysFourthSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 4, instanceStep.Func)
}

func (s *Steps) RegSysFifthSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 5, instanceStep.Func)
}

func (s *Steps) RegSysSixthSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 6, instanceStep.Func)
}

func (s *Steps) RegSysSeventhSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 7, instanceStep.Func)
}

func (s *Steps) RegSysEighthSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 8, instanceStep.Func)
}

func (s *Steps) RegSysNinethSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 9, instanceStep.Func)
}

func (s *Steps) RegSysTenthSteps(instanceStep InstanceStep) {
	s.RegActionSteps(instanceStep.StepName, 10, instanceStep.Func)
}

func (s *Steps) RegActionSteps(value string, priority int, fn FuncStep) error {
	if fn == nil {
		return fmt.Errorf("func stop should not be nil")
	}
	item := algs.NewItem(value, priority)
	s.m.Store(item, fn)
	s.q.Push(item)
	return nil
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
