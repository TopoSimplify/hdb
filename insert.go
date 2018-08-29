package hdb

import (
	"github.com/intdxdt/mbr"
	"github.com/TopoSimplify/node"
	"github.com/intdxdt/math"
)

//insert - private
func (tree *Hdb) insert(item *node.Node) *Hdb {
	if item == nil {
		return tree
	}
	var level = tree.data.height - 1

	var nd *dbNode
	var insertPath = make([]*dbNode, 0, tree.maxEntries)

	// find the best dbNode for accommodating the item, saving all nodes along the path too
	nd, insertPath = chooseSubtree(&item.MBR, &tree.data, level, insertPath)

	// put the item into the dbNode item_bbox
	nd.addChild(newLeafNode(item))
	nd.bbox.ExpandIncludeMBR(&item.MBR)

	// split on dbNode overflow propagate upwards if necessary
	level, insertPath = tree.splitOnOverflow(level, insertPath)

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(&item.MBR, insertPath, level)
	return tree
}

//insert - private
func (tree *Hdb) insertNode(item dbNode, level int) {
	var nd *dbNode
	var insertPath []*dbNode

	// find the best dbNode for accommodating the item, saving all nodes along the path too
	nd, insertPath = chooseSubtree(&item.bbox, &tree.data, level, insertPath)

	nd.children = append(nd.children, item)
	nd.bbox.ExpandIncludeMBR(&item.bbox)

	// split on dbNode overflow propagate upwards if necessary
	level, insertPath = tree.splitOnOverflow(level, insertPath)

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(&item.bbox, insertPath, level)
}

// split on dbNode overflow propagate upwards if necessary
func (tree *Hdb) splitOnOverflow(level int, insertPath []*dbNode) (int, []*dbNode) {
	for (level >= 0) && (len(insertPath[level].children) > tree.maxEntries) {
		tree.split(insertPath, level)
		level--
	}
	return level, insertPath
}

//_chooseSubtree select child of dbNode and updates path to selected dbNode.
func chooseSubtree(bbox *mbr.MBR, nd *dbNode, level int, path []*dbNode) (*dbNode, []*dbNode) {
	var child, targetNode *dbNode
	var minArea, minEnlargement float64
	var area, enlargement, d float64
	var minx, miny float64
	var maxx, maxy float64
	var ch_minx, ch_miny float64
	var ch_maxx, ch_maxy float64
	var b_minx, b_miny = bbox.MinX, bbox.MinY
	var b_maxx, b_maxy = bbox.MaxX, bbox.MaxY

	var chbox *mbr.MBR

	for {
		path = append(path, nd)
		if nd.leaf || (len(path)-1 == level) {
			break
		}
		minArea, minEnlargement = inf, inf

		for i, length := 0, len(nd.children); i < length; i++ {
			child = &nd.children[i]
			chbox = &child.bbox

			minx, miny = b_minx, b_miny
			maxx, maxy = b_maxx, b_maxy

			ch_minx, ch_miny = chbox.MinX, chbox.MinY
			ch_maxx, ch_maxy = chbox.MaxX, chbox.MaxY

			if ch_minx < minx {
				minx = ch_minx
			}
			if ch_miny < miny {
				miny = ch_miny
			}
			if ch_maxx > maxx {
				maxx = ch_maxx
			}
			if ch_maxy > maxy {
				maxy = ch_maxy
			}

			area = (ch_maxx - ch_minx) * (ch_maxy - ch_miny)
			enlargement = (maxx - minx) * (maxy - miny)

			d = enlargement - minEnlargement
			// choose entry with the least area enlargement
			if d < 0 {
				minEnlargement = enlargement
				if area < minArea {
					minArea = area
				}
				targetNode = child
			} else if d == 0 || math.Abs(d) < math.EPSILON {
				// otherwise choose one with the smallest area
				if area < minArea {
					minArea = area
					targetNode = child
				}
			}
		}

		nd = targetNode
	}

	return nd, path
}

//computes box_g margin
func bboxMargin(a *mbr.MBR) float64 {
	return (a.MaxX - a.MinX) + (a.MaxY - a.MinY)
}

//computes the intersection area of two mbrs
func intersectionArea(a, b *mbr.MBR) float64 {
	var minx, miny, maxx, maxy = a.MinX, a.MinY, a.MaxX, a.MaxY

	if !intersects(a, b) {
		return 0
	}

	if b.MinX > minx {
		minx = b.MinX
	}

	if b.MinY > miny {
		miny = b.MinY
	}

	if b.MaxX < maxx {
		maxx = b.MaxX
	}

	if b.MaxY < maxy {
		maxy = b.MaxY
	}

	return (maxx - minx) * (maxy - miny)
}

//contains tests whether a contains b
func contains(a, b *mbr.MBR) bool {
	return b.MinX >= a.MinX && b.MaxX <= a.MaxX && b.MinY >= a.MinY && b.MaxY <= a.MaxY
}

//intersects tests a intersect b (MBR)
func intersects(a, b *mbr.MBR) bool {
	return !(b.MinX > a.MaxX || b.MaxX < a.MinX || b.MinY > a.MaxY || b.MaxY < a.MinY)
}
