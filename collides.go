package hdb

import "github.com/intdxdt/mbr"

func (tree *Hdb) Collides(query mbr.MBR) bool {
    var bbox = &query
    if !intersects(bbox, &tree.Data.bbox) {
        return false
    }
    var child *dbNode
    var bln  = false
    var searchList []*dbNode
    var nd = &tree.Data

    for !bln && nd != nil {
        for i, length := 0, len(nd.children); !bln && i < length; i++ {
            child = &nd.children[i]
            if intersects(bbox, &child.bbox) {
                bln =  nd.leaf || contains(bbox, &child.bbox)
                searchList = append(searchList, child)
            }
        }
        nd, searchList = popNode(searchList)
    }
    return bln
}
