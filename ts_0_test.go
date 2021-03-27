package hdb

import (
	"github.com/TopoSimplify/node"
	"github.com/franela/goblin"
	"github.com/intdxdt/iter"
	"github.com/intdxdt/math"
	"github.com/intdxdt/mbr"
	"sort"
	"testing"
	"time"
)

type Boxes []mbr.MBR

//Len for sort interface
func (o Boxes) Len() int {
	return len(o)
}

//Swap for sort interface
func (o Boxes) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

//Less sorts boxes lexicographically
func (o Boxes) Less(i, j int) bool {
	var d = o[i].MinX - o[j].MinX
	if d == 0 || math.Abs(d) < math.EPSILON {
		d = o[i].MinY - o[j].MinY
	}
	return d < 0
}

func someData(n int) []mbr.MBR {
	var data = make([]mbr.MBR, n)
	for i := 0; i < n; i++ {
		data[i] = mbr.CreateMBR(float64(i), float64(i), float64(i), float64(i))
	}
	return data
}

func testResults(g *goblin.G, objects []*node.Node, boxes Boxes) {
	var results = make([]mbr.MBR, 0, len(objects))
	for i := range objects {
		results = append(results, objects[i].MBR)
	}

	sort.Sort(Boxes(results))
	sort.Sort(boxes)
	g.Assert(len(results)).Equal(len(boxes))
	for i, n := range results {
		g.Assert(n.Equals(&boxes[i])).IsTrue()
	}
}

func getObjs(nodes []dbNode) []*node.Node {
	var objs = make([]*node.Node, 0, len(nodes))
	for _, o := range nodes {
		objs = append(objs, o.item)
	}
	return objs
}

//data from rbush 1.4.2
var data = []mbr.MBR{{0, 0, 0, 0}, {10, 10, 10, 10}, {20, 20, 20, 20}, {25, 0, 25, 0}, {35, 10, 35, 10}, {45, 20, 45, 20}, {0, 25, 0, 25}, {10, 35, 10, 35},
	{20, 45, 20, 45}, {25, 25, 25, 25}, {35, 35, 35, 35}, {45, 45, 45, 45}, {50, 0, 50, 0}, {60, 10, 60, 10}, {70, 20, 70, 20}, {75, 0, 75, 0},
	{85, 10, 85, 10}, {95, 20, 95, 20}, {50, 25, 50, 25}, {60, 35, 60, 35}, {70, 45, 70, 45}, {75, 25, 75, 25}, {85, 35, 85, 35}, {95, 45, 95, 45},
	{0, 50, 0, 50}, {10, 60, 10, 60}, {20, 70, 20, 70}, {25, 50, 25, 50}, {35, 60, 35, 60}, {45, 70, 45, 70}, {0, 75, 0, 75}, {10, 85, 10, 85},
	{20, 95, 20, 95}, {25, 75, 25, 75}, {35, 85, 35, 85}, {45, 95, 45, 95}, {50, 50, 50, 50}, {60, 60, 60, 60}, {70, 70, 70, 70}, {75, 50, 75, 50},
	{85, 60, 85, 60}, {95, 70, 95, 70}, {50, 75, 50, 75}, {60, 85, 60, 85}, {70, 95, 70, 95}, {75, 75, 75, 75}, {85, 85, 85, 85}, {95, 95, 95, 95}}

func TestRtreeRbush(t *testing.T) {
	var g = goblin.Goblin(t)
	var id = iter.NewIgen(0)

	g.Describe("Hdb Tests - From Rbush", func() {
		g.It("should test load 9 & 10", func() {
			var tree0 = NewHdb().loadBoxes(id, someData(0))
			g.Assert(tree0.data.height).Equal(1)

			var tree1 = NewHdb(8).loadBoxes(id, someData(8))
			g.Assert(tree1.data.height).Equal(1)

			var tree2 = NewHdb(8).loadBoxes(id, someData(10))
			g.Assert(tree2.data.height).Equal(2)
		})

		g.It("tests search with some other", func() {
			var data = []mbr.MBR{
				{-115, 45, -105, 55}, {105, 45, 115, 55}, {105, -55, 115, -45}, {-115, -55, -105, -45},
			}
			var tree = NewHdb(4)
			tree.loadBoxes(id, data)
			testResults(g, tree.Search(mbr.CreateMBR(-180, -90, 180, 90)), []mbr.MBR{
				{-115, 45, -105, 55}, {105, 45, 115, 55}, {105, -55, 115, -45}, {-115, -55, -105, -45},
			})

			testResults(g, tree.Search(mbr.CreateMBR(-180, -90, 0, 90)), []mbr.MBR{
				{-115, 45, -105, 55}, {-115, -55, -105, -45},
			})
			testResults(g, tree.Search(mbr.CreateMBR(0, -90, 180, 90)), []mbr.MBR{
				{105, 45, 115, 55}, {105, -55, 115, -45},
			})
			testResults(g, tree.Search(mbr.CreateMBR(-180, 0, 180, 90)), []mbr.MBR{
				{-115, 45, -105, 55}, {105, 45, 115, 55},
			})
			testResults(g, tree.Search(mbr.CreateMBR(-180, -90, 180, 0)), []mbr.MBR{
				{105, -55, 115, -45}, {-115, -55, -105, -45},
			})
		})

		g.It("#load uses standard insertion when given a low number of items", func() {
			var id = iter.NewIgen(0)
			var id2 = iter.NewIgen(0)

			var tree = NewHdb(8).loadBoxes(id, data).loadBoxes(id, data[0:3])

			var tree2 = NewHdb(8).loadBoxes(id2, data).insert(
				&node.Node{Id: id2.Next(), MBR: data[0]},
			).insert(
				&node.Node{Id: id2.Next(), MBR: data[1]},
			).insert(
				&node.Node{Id: id2.Next(), MBR: data[2]},
			)

			g.Assert(tree.data).Eql(tree2.data)
		})

		g.It("#load does nothing if loading empty data", func() {
			var tree = NewHdb(0).Load(make([]node.Node, 0))
			g.Assert(tree.IsEmpty()).IsTrue()
		})

		g.It("#load properly splits tree root when merging trees of the same height", func() {
			var cloneData = make([]mbr.MBR, len(data))
			for i := 0; i < len(data); i++ {
				cloneData[i] = data[i].Clone()
			}
			for i := range data {
				cloneData = append(cloneData, data[i].Clone())
			}
			var tree = NewHdb(4).loadBoxes(id, data).loadBoxes(id, data)
			testResults(g, tree.All(), cloneData)
		})

		g.It("#load properly merges data of smaller or bigger tree heights", func() {
			var smaller = someData(10)
			var cloneData = make([]mbr.MBR, len(data))
			for i := 0; i < len(data); i++ {
				cloneData[i] = data[i].Clone()
			}
			for i := 0; i < len(smaller); i++ {
				cloneData = append(cloneData, smaller[i].Clone())
			}

			var tree1 = NewHdb(4).loadBoxes(id, data).loadBoxes(id, smaller)
			var tree2 = NewHdb(4).loadBoxes(id, smaller).loadBoxes(id, data)
			g.Assert(tree1.data.height).Equal(tree2.data.height)
			testResults(g, tree1.All(), cloneData)
			testResults(g, tree2.All(), cloneData)
		})

		g.It("#load properly merges data of smaller or bigger tree heights 2", func() {
			N = 8020
			var smaller = GenDataItems(N, 1)
			var larger = GenDataItems(2*N, 1)
			var cloneData = make([]mbr.MBR, len(larger))

			for i := 0; i < len(larger); i++ {
				cloneData[i] = larger[i].Clone()
			}
			for i := 0; i < len(smaller); i++ {
				cloneData = append(cloneData, smaller[i].Clone())
			}

			var tree1 = NewHdb(64).loadBoxes(id, larger).loadBoxes(id, smaller)
			var tree2 = NewHdb(64).loadBoxes(id, smaller).loadBoxes(id, larger)
			g.Assert(tree1.data.height).Equal(tree2.data.height)
			testResults(g, tree1.All(), cloneData)
			testResults(g, tree2.All(), cloneData)
		})

		g.It("#search finds matching points in the tree given a bbox", func() {
			var tree = NewHdb(4).loadBoxes(id, data)
			var result = tree.Search(mbr.CreateMBR(40, 20, 80, 70))
			testResults(g, result, []mbr.MBR{
				{70, 20, 70, 20}, {75, 25, 75, 25}, {45, 45, 45, 45}, {50, 50, 50, 50}, {60, 60, 60, 60}, {70, 70, 70, 70},
				{45, 20, 45, 20}, {45, 70, 45, 70}, {75, 50, 75, 50}, {50, 25, 50, 25}, {60, 35, 60, 35}, {70, 45, 70, 45},
			})
		})

		g.It("#collides returns true when search finds matching points", func() {
			var tree = NewHdb(4).loadBoxes(id, data)
			g.Assert(tree.Collides(mbr.CreateMBR(40, 20, 80, 70))).IsTrue()
			g.Assert(tree.Collides(mbr.CreateMBR(200, 200, 210, 210))).IsFalse()
		})

		g.It("#search returns an empty array if nothing found", func() {
			var result = NewHdb(4).loadBoxes(id, data).Search(
				mbr.CreateMBR(200, 200, 210, 210),
			)
			g.Assert(len(result)).Equal(0)
		})

		g.It("#all <==>.data returns all points in the tree", func() {
			var cloneData = make([]mbr.MBR, len(data))
			for i := 0; i < len(data); i++ {
				cloneData[i] = data[i]
			}

			var tree = NewHdb(4).loadBoxes(id, data)
			var result = tree.Search(mbr.CreateMBR(0, 0, 100, 100))
			testResults(g, result, cloneData)
		})

		g.It("#insert adds an item to an existing tree correctly", func() {
			var data = []mbr.MBR{{0, 0, 0, 0}, {2, 2, 2, 2}, {1, 1, 1, 1}}
			var tree = NewHdb(4)
			tree.loadBoxes(id, data)
			tree.insert(&node.Node{Id: id.Next(), MBR: mbr.CreateMBR(3, 3, 3, 3)})
			g.Assert(tree.data.leaf).IsTrue()
			g.Assert(tree.data.height).Equal(1)

			var box = mbr.CreateMBR(0, 0, 3, 3)
			g.Assert(tree.data.bbox.Equals(&box)).IsTrue()
			testResults(g, getObjs(tree.data.children), []mbr.MBR{
				{0, 0, 0, 0}, {1, 1, 1, 1}, {2, 2, 2, 2}, {3, 3, 3, 3},
			})
		})

		g.It("#insert does nothing if given nil", func() {
			var o *node.Node
			var id = iter.NewIgen(0)
			var id2 = iter.NewIgen(0)
			var tree = NewHdb(4).loadBoxes(id, data)
			g.Assert(tree.data).Eql(NewHdb(4).loadBoxes(id2, data).insert(o).data)
		})

		g.It("#insert forms a valid tree if items are inserted one by one", func() {
			var tree = NewHdb(4)
			for i := 0; i < len(data); i++ {
				tree.insert(&node.Node{Id: id.Next(), MBR: data[i]})
			}

			var tree2 = NewHdb(4).loadBoxes(id, data)
			g.Assert(tree.data.height-tree2.data.height <= 1).IsTrue()

			var boxes2 = make([]mbr.MBR, 0)
			var all2 = tree2.All()
			for i := 0; i < len(all2); i++ {
				boxes2 = append(boxes2, all2[i].MBR)
			}
			testResults(g, tree.All(), boxes2)
		})

		g.It("#remove removes items correctly", func() {
			var tree = NewHdb(4).loadBoxes(id, data)
			var N = len(data)
			tree.removeMBR(&data[0])
			tree.removeMBR(&data[1])
			tree.removeMBR(&data[2])

			tree.removeMBR(&data[N-1])
			tree.removeMBR(&data[N-2])
			tree.removeMBR(&data[N-3])
			var cloneData []mbr.MBR
			for i := 3; i < len(data)-3; i++ {
				cloneData = append(cloneData, data[i].Clone())
			}

			testResults(g, tree.All(), cloneData)
		})

		g.It("#remove does nothing if nothing found", func() {
			var item *node.Node
			var id = iter.NewIgen(0)
			var id2 = iter.NewIgen(0)

			var tree = NewHdb(0).loadBoxes(id, data)
			var tree2 = NewHdb(0).loadBoxes(id2, data)

			var query = mbr.CreateMBR(13, 13, 13, 13)
			var results13 = tree2.Search(query)
			var node13 *node.Node
			for i := range results13 {
				if results13[i].MBR.Equals(&query) {
					node13 = results13[i]
				}
			}
			var querybox = node13
			g.Assert(tree.data).Eql(tree2.removeMBR(&query).data)
			g.Assert(tree.data).Eql(tree2.Remove(querybox).data)
			g.Assert(tree.data).Eql(tree2.Remove(item).data)
		})

		g.It("#remove brings the tree to a clear state when removing everything one by one", func() {
			var tree = NewHdb(4).loadBoxes(id, data)
			var result = tree.Search(mbr.CreateMBR(0, 0, 100, 100))
			for i := range result {
				tree.Remove(result[i])
			}
			g.Assert(tree.Remove(&node.Node{}).IsEmpty()).IsTrue()
		})

		g.It("#clear should clear all the data in the tree", func() {
			var tree = NewHdb(4).loadBoxes(id, data).Clear()
			g.Assert(tree.IsEmpty()).IsTrue()
		})

		g.It("should have chainable API", func() {
			g.Assert(NewHdb(4).loadBoxes(id, data).insert(
				&node.Node{Id: id.Next(), MBR: data[0]},
			).removeMBR(&data[0]).Clear().IsEmpty()).IsTrue()
		})
	})

}

func TestRtreeUtil(t *testing.T) {
	var g = goblin.Goblin(t)
	var id = iter.NewIgen(0)
	g.Describe("Hdb Util", func() {
		g.It("tests pop nodes", func() {
			g.Timeout(1 * time.Hour)

			var a = createDBNode(&node.Node{Id: id.Next(), MBR: emptyMBR()}, 0, true, nil)
			var b = createDBNode(&node.Node{Id: id.Next(), MBR: emptyMBR()}, 1, true, nil)
			var c = createDBNode(&node.Node{Id: id.Next(), MBR: emptyMBR()}, 1, true, nil)
			var nodes = make([]*dbNode, 0)
			var n *dbNode

			n, nodes = popNode(nodes)
			g.Assert(n == nil).IsTrue()

			nodes = []*dbNode{&a, &b, &c}
			g.Assert(len(nodes)).Equal(3)

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(2)
			g.Assert(n == &c).IsTrue()

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(1)
			g.Assert(n == &b).IsTrue()

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(0)
			g.Assert(n == &a).IsTrue()

			n, nodes = popNode(nodes)
			g.Assert(len(nodes)).Equal(0)
			g.Assert(n == nil).IsTrue()

			var nodesABC = []dbNode{a, b, c}
			g.Assert(len(nodesABC)).Equal(3)
			nodesABC = removeNode(nodesABC, 1)
			g.Assert(len(nodesABC)).Equal(2)
			nodesABC = removeNode(nodesABC, 4)
			g.Assert(len(nodesABC)).Equal(2)

		})

		g.It("tests pop index", func() {
			a := 0
			b := 1
			c := 2
			var indexes = make([]int, 0)
			var n int

			n = popIndex(&indexes)
			g.Assert(n == 0).IsTrue()

			indexes = []int{a, b, c}
			g.Assert(len(indexes)).Equal(3)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(2)
			g.Assert(n).Eql(c)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(1)
			g.Assert(n).Eql(b)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(0)
			g.Assert(n).Eql(a)

			n = popIndex(&indexes)
			g.Assert(len(indexes)).Equal(0)
			g.Assert(n == 0).IsTrue()
		})

		g.It("util max min", func() {
			g.Assert(max(2.3, 4.5)).Equal(4.5)
			g.Assert(max(4.5, 2.3)).Equal(4.5)
			g.Assert(min(2.3, 4.5)).Equal(2.3)
			g.Assert(min(4.5, 2.3)).Equal(2.3)
			g.Assert(min(4.50233, 4.50233)).Equal(4.50233)
			g.Assert(max(4.50233, 4.50233)).Equal(4.50233)
		})
	})
}
