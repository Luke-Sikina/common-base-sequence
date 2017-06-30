package main

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddCountMapToEmptyHeap(t *testing.T) {
	testHeapSize := 5
	counts := map[Hash]uint32{0: 10, 1: 9, 2: 8, 3: 7, 4: 6, 5: 5, 6: 4, 7: 3, 8: 2, 9: 1}
	subject := &HashCountPairHeap{}
	heap.Init(subject)
	addCountMapToHeap(counts, subject, testHeapSize)

	expectedPairs := []HashCountPair{{4, 6}, {3, 7}, {2, 8}, {1, 9}, {0, 10}}
	for _, expected := range expectedPairs {
		actual := heap.Pop(subject).(HashCountPair)
		assert.Equal(t, expected.hash, actual.hash)
		assert.Equal(t, expected.count, actual.count)
	}
}
