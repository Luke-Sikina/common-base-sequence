package main

type HashCountPair [2]uint32

func (pair *HashCountPair) hash() uint32 {
	return pair[0]
}

func (pair *HashCountPair) count() uint32 {
	return pair[1]
}

type HashCountPairHeap []HashCountPair

func (h HashCountPairHeap) Len() int           { return len(h) }
func (h HashCountPairHeap) Less(i, j int) bool { return h[i][1] < h[j][1] }
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
