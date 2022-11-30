package algs

import (
	"sync"
	"sync/atomic"
)

type GraphSearch[T comparable] struct {
	count  uint64
	marked sync.Map
	g      *Graph[T]
	st     T
}

func NewGraphSearch[T comparable](g *Graph[T], v T) GraphSearch[T] {
	gs := GraphSearch[T]{count: 0, g: g, st: v}
	gs.dfs(v)
	return gs
}

func (gs *GraphSearch[T]) Search(w T) bool {
	if _, ok := gs.marked.Load(w); ok {
		return true
	}
	return false
}

func (gs *GraphSearch[T]) Count() uint64 {
	return atomic.LoadUint64(&gs.count)
}

func (gs *GraphSearch[T]) dfs(v T) {
	gs.marked.Store(v, true)
	atomic.AddUint64(&gs.count, 1)
	for _, w := range gs.g.Adj(v) {
		if _, ok := gs.marked.Load(w); !ok {
			gs.dfs(w)
		}
	}
}
