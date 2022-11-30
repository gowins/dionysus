package algs

import (
	"sync"
)

type Graph[T comparable] struct {
	mu  sync.RWMutex
	v   int //points num
	e   int //edges  num
	adj map[T][]T
}

func NewGraph[T comparable]() *Graph[T] {
	return &Graph[T]{v: 0, e: 0, adj: make(map[T][]T)}
}

// V points num
func (g *Graph[T]) V() int {
	return g.v
}

// E edge num
func (g *Graph[T]) E() int {
	return g.e
}

// AddPoint add point v for graph
func (g *Graph[T]) AddPoint(v T) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.adj[v]; !ok {
		g.v += 1
		g.adj[v] = []T{}
	}
}

// AddEdge add Edge with two points v, w
// todo nil is untyped, so v,w will not be nil.  unnecessary check v==nil or w==nil

func (g *Graph[T]) AddEdge(v T, w T) {
	//todo optimize lock
	g.AddPoint(v)
	g.AddPoint(w)
	g.mu.Lock()
	defer g.mu.Unlock()
	var v2w bool
	for _, x := range g.adj[v] {
		if x == w {
			v2w = true
			break
		}
	}
	if !v2w {
		g.e += 1
		g.adj[v] = append(g.adj[v], w)
		g.adj[w] = append(g.adj[w], v)
	}

}

// Adj return all points linked with v
func (g *Graph[T]) Adj(v T) (points []T) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	points = g.adj[v]
	return
}
