package algs

import (
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Item struct {
	// The value of the item; arbitrary.
	value string

	// The priority of the item in the queue.
	priority int

	// The index is needed by update and is maintained by the heap.Interface methods.
	// The index of the item in the heap.
	index int
}

func (i Item) Value() string {
	return i.value
}

// PriorityQueue wrapper
type PriorityQueue struct {
	pq *priorityQueue
}

// Pop update modifies the priority and value of an Item in the queue.
func (p *PriorityQueue) Pop() (*Item, bool) {
	if p.Len() <= 0 {
		return nil, false
	}
	i := heap.Pop(p.pq)
	item, ok := i.(*Item)
	return item, ok
}

// Update update modifies the priority and value of an Item in the queue.
func (p *PriorityQueue) Update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(p.pq, item.index)
}

func (p *PriorityQueue) Push(item *Item) {
	heap.Push(p.pq, item)
}
func (p *PriorityQueue) Len() int { return p.pq.Len() }

// A PriorityQueue implements heap.Interface and holds Items.
type priorityQueue []*Item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func NewPQ() *PriorityQueue {
	pq := PriorityQueue{pq: &priorityQueue{}}
	heap.Init(pq.pq)
	return &pq
}

func NewItem(value string, priority int) *Item {
	return &Item{value: value, priority: priority}
}
