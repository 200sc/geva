package pairing

// The first element in each heap node
// is the value of that node, and the
// second is the index that value
// represents in some external structure.
type ValIndexHeap [][2]int

func (h ValIndexHeap) Len() int           { return len(h) }
func (h ValIndexHeap) Less(i, j int) bool { return h[i][0] < h[j][0] }
func (h ValIndexHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *ValIndexHeap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.([2]int))
}

func (h *ValIndexHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
