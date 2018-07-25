package hdb

import (
	"github.com/intdxdt/mbr"
	"github.com/TopoSimplify/node"
)

//dbNode type for internal dbNode
type dbNode struct {
	children []dbNode
	item     *node.Node
	height   int
	leaf     bool
	bbox     mbr.MBR
}

//createDBNode creates a new dbNode
func createDBNode(item *node.Node, height int, leaf bool, children []dbNode) dbNode {
	var box mbr.MBR
	if item == nil {
		box = emptyMBR()
	} else {
		box = item.MBR
	}

	return dbNode{
		children: children,
		item:     item,
		height:   height,
		leaf:     leaf,
		bbox:     box,
	}
}

//dbNode type for internal dbNode
func newLeafNode(item *node.Node) dbNode {
	return dbNode{
		children: []dbNode{},
		item:     item,
		height:   1,
		leaf:     true,
		bbox:     item.MBR,
	}
}

//MBR returns bbox property
func (nd *dbNode) BBox() *mbr.MBR {
	return &nd.bbox
}

//add child
func (nd *dbNode) addChild(child dbNode) {
	nd.children = append(nd.children, child)
}

//Constructs children of dbNode
func makeChildren(items []*node.Node) []dbNode {
	var chs = make([]dbNode, 0, len(items))
	for i := range items {
		chs = append(chs, newLeafNode(items[i]))
	}
	return chs
}
