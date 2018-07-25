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

func (tree *Hdb) Clear() *Hdb {
	tree.Data = newNode(nil, 1, true, []dbNode{}, )
	return tree
}

//IsEmpty checks for empty tree
func (tree *Hdb) IsEmpty() bool {
	return len(tree.Data.children) == 0
}
