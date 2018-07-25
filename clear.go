package hdb

import (
	"math"
	"github.com/intdxdt/mbr"
)

func emptyMBR() mbr.MBR {
	return mbr.MBR{
		math.Inf(1), math.Inf(1),
		math.Inf(-1), math.Inf(-1),
	}
}

func (tree *Hdb) Clear() *Hdb {
	tree.data = createDBNode(nil, 1, true, []dbNode{}, )
	return tree
}

//IsEmpty checks for empty tree
func (tree *Hdb) IsEmpty() bool {
	return len(tree.data.children) == 0
}
