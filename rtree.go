package hdb
/*
 (c) 2015, Titus Tienaah
 A library for 2D spatial indexing of points and rectangles.
 https://github.com/mourner/rbush
 @after  (c) 2015, Vladimir Agafonkin
*/

//Hdb type
type Hdb struct {
	data       dbNode
	maxEntries int
	minEntries int
}

func NewHdb(nodeCap ...int) *Hdb {
	var bucketSize = 8
	var tree = Hdb{}
	tree.Clear()
	if len(nodeCap) > 0 {
		bucketSize = nodeCap[0]
	}
	// bucket size(dbNode) == 8 by default
	tree.maxEntries = maxEntries(bucketSize)
	tree.minEntries = minEntries(tree.maxEntries)
	return &tree
}
