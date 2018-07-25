package hdb

import (
	"github.com/intdxdt/heap"
	"github.com/intdxdt/mbr"
	"github.com/TopoSimplify/node"
	"fmt"
)

func predicate(_ *KObj) (bool, bool) {
	return true, false
}

func (tree *Hdb) Knn(
	query mbr.MBR, limit int, score func(*mbr.MBR, *KObj) float64,
	predicates ...func(*KObj) (bool, bool)) []*node.Node {

	var predFn = predicate
	if len(predicates) > 0 {
		predFn = predicates[0]
	}

	var nd = &tree.Data
	var result []*node.Node
	var child *dbNode
	var stop, pred bool
	var queue = heap.NewHeap(kObjCmp, heap.NewHeapType().AsMin())

	for !stop && (nd != nil) {
		for i := range nd.children {
			child = &nd.children[i]
			fmt.Println(child.item.Geometry.WKT())
			var o = &KObj{
				dbNode: child,
				MBR:    &child.bbox,
				IsItem: len(child.children) == 0,
				Dist:   -1,
			}
			o.Dist = score(&query, o)
			queue.Push(o)
		}

		for !queue.IsEmpty() && queue.Peek().(*KObj).IsItem {
			var candidate = queue.Pop().(*KObj)
			pred, stop = predFn(candidate)
			if pred {
				result = append(result, candidate.GetNode())
			}

			if stop {
				break
			}

			if limit != 0 && len(result) == limit {
				return result
			}
		}

		if !stop {
			var q = queue.Pop()
			if q == nil {
				nd = nil
			} else {
				nd = q.(*KObj).dbNode
			}
		}
	}
	return result
}
