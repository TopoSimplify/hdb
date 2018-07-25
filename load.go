package hdb

import (
	"github.com/intdxdt/mbr"
	"github.com/TopoSimplify/node"
)

//loadBoxes loads bounding boxes
func (tree *hdb) loadBoxes(data []mbr.MBR) *hdb {
	var items = make([]*node.Node, 0, len(data))
	for i := range data {
		items = append(items, &node.Node{Id: i, MBR: data[i]})
	}
	return tree.Load(items)
}

//Load implements bulk loading
func (tree *hdb) Load(data []*node.Node) *hdb {
	var n  = len(data)
	if n < tree.minEntries {
		for i := range data {
			tree.Insert(data[i])
		}
		return tree
	}

	//var data = make([]*node.Node, 0, n)
	//for i := range items {
	//	data = append(data, items[i])
	//}


	// recursively build the tree with the given data from stratch using OMT algorithm
	var nd = tree.buildTree(data, 0, n-1, 0)

	if len(tree.Data.children) == 0 {
		// save as is if tree is empty
		tree.Data = nd
	} else if tree.Data.height == nd.height {
		// split root if trees have the same height
		tree.splitRoot(tree.Data, nd)
	} else {
		if tree.Data.height < nd.height {
			// swap trees if inserted one is bigger
			tree.Data, nd = nd, tree.Data
		}

		// insert the small tree into the large tree at appropriate level
		tree.insertNode(nd, tree.Data.height-nd.height-1)
	}

	return tree
}
