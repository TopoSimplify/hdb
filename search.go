package hdb

import (
	"github.com/intdxdt/mbr"
	"github.com/TopoSimplify/node"
)

//Search item
func (tree *Hdb) Search(query mbr.MBR) []*node.Node {
	var bbox = &query
	var result []*node.Node
	var nd = &tree.data

	if !intersects(bbox, &nd.bbox) {
		return []*node.Node{}
	}

	var nodesToSearch []*dbNode
	var child *dbNode
	var childBBox *mbr.MBR

	for {
		for i, length := 0, len(nd.children); i < length; i++ {
			child = &nd.children[i]
			childBBox = &child.bbox

			if intersects(bbox, childBBox) {
				if nd.leaf {
					result = append(result, child.item)
				} else if contains(bbox, childBBox) {
					result = all(child, result)
				} else {
					nodesToSearch = append(nodesToSearch, child)
				}
			}
		}

		nd, nodesToSearch = popNode(nodesToSearch)
		if nd == nil {
			break
		}
	}
	return result
}

//All items from  root dbNode
func (tree *Hdb) All() []*node.Node {
	return all(&tree.data, []*node.Node{})
}

//all - fetch all items from dbNode
func all(nd *dbNode, result []*node.Node) []*node.Node {
	var nodesToSearch []*dbNode
	for {
		if nd.leaf {
			for i := range nd.children {
				result = append(result, nd.children[i].item)
			}
		} else {
			for i := range nd.children {
				nodesToSearch = append(nodesToSearch, &nd.children[i])
			}
		}

		nd, nodesToSearch = popNode(nodesToSearch)
		if nd == nil {
			break
		}
	}

	return result
}
