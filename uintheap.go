package main

type HashCountPair struct {
	hash  Hash
	count uint32
}

type HashCountPairHeap []HashCountPair

func (h HashCountPairHeap) Len() int           { return len(h) }
func (h HashCountPairHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h HashCountPairHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *HashCountPairHeap) Push(x interface{}) {
	*h = append(*h, x.(HashCountPair))
}

func (h *HashCountPairHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
