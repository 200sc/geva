package eda

// ComplexDiff is a pair of indices and some value representing how valuable
// the pair is
type ComplexDiff struct {
	i, j int
	gain float64
}

// ComplexHeap is a heap of ComplexDiffs, heaped on their gain value
type ComplexHeap []ComplexDiff

func (ch ComplexHeap) Len() int {
	return len(ch)
}
func (ch ComplexHeap) Less(i, j int) bool {
	return ch[i].gain < ch[j].gain
}
func (ch ComplexHeap) Swap(i, j int) {
	ch[i], ch[j] = ch[j], ch[i]
}

// Push is boilerplate for pushing to a heap
func (ch *ComplexHeap) Push(x interface{}) {
	*ch = append(*ch, x.(ComplexDiff))
}

// Pop is boilerplate for popping a heap
func (ch *ComplexHeap) Pop() interface{} {
	old := *ch
	n := len(old)
	x := old[n-1]
	*ch = old[0 : n-1]
	return x
}
