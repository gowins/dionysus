package algs

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGraphSearch_Search(t *testing.T) {
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
		{5, 3}}
	g := NewGraph[int]()
	for _, d := range ds {
		g.AddEdge(d[0], d[1])
	}
	gs := NewGraphSearch(g, 0)

	Convey("Test Search", t, func() {
		So(gs.Count(), ShouldNotEqual, g.V())
		So(gs.Count(), ShouldEqual, 7)
		So(gs.Search(0), ShouldEqual, true)
		So(gs.Search(1), ShouldEqual, true)
		So(gs.Search(2), ShouldEqual, true)
		So(gs.Search(3), ShouldEqual, true)
		So(gs.Search(4), ShouldEqual, true)
		So(gs.Search(5), ShouldEqual, true)
		So(gs.Search(6), ShouldEqual, true)
		So(gs.Search(7), ShouldEqual, false)
	})

	ds = [][]int{
		{0, 5},
		{2, 4},
		{2, 3},
		{1, 2},
		{0, 1},
		{3, 4},
		{3, 5},
		{0, 2}}
	g = NewGraph[int]()
	for _, d := range ds {
		g.AddEdge(d[0], d[1])
	}
	gs = NewGraphSearch(g, 0)
	Convey("Test connected", t, func() {
		So(gs.Count(), ShouldEqual, g.V())
	})
}
