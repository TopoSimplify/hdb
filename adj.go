package hdb

import "github.com/intdxdt/mbr"

// adjust bboxes along the given tree path
func (tree *hdb) adjustParentBBoxes(bbox *mbr.MBR, path []*dbNode, level int) {
	for i := level; i >= 0; i-- {
		extend(&path[i].bbox, bbox)
	}
}
