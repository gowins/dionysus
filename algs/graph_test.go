package algs

import (
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGraph(t *testing.T) {
	ds := [][]int{
		{0, 5},
		{4, 3},
		{1, 0},
		{9, 12},
		{6, 4},
		{5, 4},
		{0, 2},
		{11, 12},
		{9, 10},
		{0, 6},
		{7, 8},
		{9, 11},
		{9, 11},
		{5, 3}}
	g := NewGraph[int]()
	for _, d := range ds {
		g.AddEdge(d[0], d[1])
	}
	Convey("Test graph", t, func() {
		So(g.V(), ShouldEqual, 13)
		So(g.E(), ShouldEqual, 13)
		s0 := g.Adj(0)
		sort.SliceStable(s0, func(i, j int) bool { return s0[i] < s0[j] })
		So(s0, ShouldResemble, []int{1, 2, 5, 6})
	})
}
