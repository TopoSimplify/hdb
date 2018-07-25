package hdb

import (
	"math"
	"github.com/intdxdt/mbr"
	"github.com/TopoSimplify/node"
)

func emptyMBR() mbr.MBR {
	return mbr.MBR{
		math.Inf(1), math.Inf(1),
		math.Inf(-1), math.Inf(-1),
	}
}

func emptyObject() *node.Node {
	return &node.Node{
		MBR: emptyMBR(),
	}
}



func (tree *hdb) Clear() *hdb {
	tree.Data = newNode(
		emptyObject(), 1, true, []dbNode{},
	)
	return tree
}

//IsEmpty checks for empty tree
func (tree *hdb) IsEmpty() bool {
	return len(tree.Data.children) == 0
}
