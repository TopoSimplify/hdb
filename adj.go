package hdb

import "github.com/intdxdt/mbr"

// adjust bboxes along the given tree path
func (tree *Hdb) adjustParentBBoxes(bbox *mbr.MBR, path []*dbNode, level int) {
	for i := level; i >= 0; i-- {
		path[i].bbox.ExpandIncludeMBR(bbox)
	}
}
