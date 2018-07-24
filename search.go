package hdb

import (
	"github.com/intdxdt/mbr"
)

//Search item
func (tree *hdb) Search(query mbr.MBR) []*Obj {
	var bbox = &query
	var result []*Obj
	var nd = &tree.Data

	if !intersects(bbox, &nd.bbox) {
		return []*Obj{}
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

	//var objs = make([]*Obj, 0, len(result))
	//for i := range result {
	//	objs = append(objs, result[i].item)
	//}
	return result
}

//All items from  root dbNode
func (tree *hdb) All() []*Obj {
	return all(&tree.Data, []*Obj{})
}

//all - fetch all items from dbNode
func all(nd *dbNode, result []*Obj) []*Obj {
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
