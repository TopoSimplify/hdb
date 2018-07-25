package hdb

import (
	"github.com/intdxdt/mbr"
	"github.com/TopoSimplify/node"
)

//insert - private
func (tree *hdb) Insert(item *node.Node) *hdb {
	if item == nil {
		return tree
	}
	var nd *dbNode
	var level = tree.Data.height - 1
	var insertPath = make([]*dbNode, 0, tree.maxEntries)

	// find the best dbNode for accommodating the item, saving all nodes along the path too
	nd, insertPath = chooseSubtree(
		&item.MBR, &tree.Data, level, insertPath,
	)

	// put the item into the dbNode item_bbox
	nd.addChild(newLeafNode(item))
	extend(&nd.bbox, &item.MBR)

	// split on dbNode overflow propagate upwards if necessary
	level, insertPath = tree.splitOnOverflow(level, insertPath)

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(&item.MBR, insertPath, level)
	return tree
}

//insert - private
func (tree *hdb) insertNode(item dbNode, level int) {
	var nd *dbNode
	var insertPath []*dbNode

	// find the best dbNode for accommodating the item, saving all nodes along the path too
	nd, insertPath = chooseSubtree(&item.bbox, &tree.Data, level, insertPath)

	nd.children = append(nd.children, item)
	extend(&nd.bbox, &item.bbox)

	// split on dbNode overflow propagate upwards if necessary
	level, insertPath = tree.splitOnOverflow(level, insertPath)

	// adjust bboxes along the insertion path
	tree.adjustParentBBoxes(&item.bbox, insertPath, level)
}

// split on dbNode overflow propagate upwards if necessary
func (tree *hdb) splitOnOverflow(level int, insertPath []*dbNode) (int, []*dbNode) {
	for (level >= 0) && (len(insertPath[level].children) > tree.maxEntries) {
		tree.split(insertPath, level)
		level--
	}
	return level, insertPath
}

//_chooseSubtree select child of dbNode and updates path to selected dbNode.
func chooseSubtree(bbox *mbr.MBR, nd *dbNode, level int, path []*dbNode) (*dbNode, []*dbNode) {
	var child *dbNode
	var targetNode *dbNode
	var minArea, minEnlargement, area, enlargement float64

	for {
		path = append(path, nd)
		if nd.leaf || (len(path)-1 == level) {
			break
		}
		minArea, minEnlargement = inf, inf

		for i, length := 0, len(nd.children); i < length; i++ {
			child = &nd.children[i]
			area = bboxArea(&child.bbox)
			enlargement = enlargedArea(bbox, &child.bbox) - area

			// choose entry with the least area enlargement
			if enlargement < minEnlargement {
				minEnlargement = enlargement
				if area < minArea {
					minArea = area
				}
				targetNode = child
			} else if feq(enlargement, minEnlargement) {
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

//extend bounding box
func extend(a, b *mbr.MBR) *mbr.MBR {
	a[x1] = min(a[x1], b[x1])
	a[y1] = min(a[y1], b[y1])
	a[x2] = max(a[x2], b[x2])
	a[y2] = max(a[y2], b[y2])
	return a
}

//computes area of bounding box
func bboxArea(a *mbr.MBR) float64 {
	return (a[x2] - a[x1]) * (a[y2] - a[y1])
}

//computes box margin
func bboxMargin(a *mbr.MBR) float64 {
	return (a[x2] - a[x1]) + (a[y2] - a[y1])
}

//computes enlarged area given two mbrs
func enlargedArea(a, b *mbr.MBR) float64 {
	return (max(a[x2], b[x2]) - min(a[x1], b[x1])) * (max(a[y2], b[y2]) - min(a[y1], b[y1]))
}

//computes the intersection area of two mbrs
func intersectionArea(a, b *mbr.MBR) float64 {
	var minx, miny, maxx, maxy = a[x1], a[y1], a[x2], a[y2]

	if !intersects(a, b) {
		return 0.0
	}

	if b[x1] > minx {
		minx = b[x1]
	}

	if b[y1] > miny {
		miny = b[y1]
	}

	if b[x2] < maxx {
		maxx = b[x2]
	}

	if b[y2] < maxy {
		maxy = b[y2]
	}

	return (maxx - minx) * (maxy - miny)
}

//contains tests whether a contains b
func contains(a, b *mbr.MBR) bool {
	return b[x1] >= a[x1] && b[x2] <= a[x2] && b[y1] >= a[y1] && b[y2] <= a[y2]
}

//intersects tests a intersect b (MBR)
func intersects(a, b *mbr.MBR) bool {
	return !(b[x1] > a[x2] || b[x2] < a[x1] || b[y1] > a[y2] || b[y2] < a[y1])
}
