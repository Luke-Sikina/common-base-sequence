package main

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGatherCommonSequences(t *testing.T) {
	sequence := "\nNNNNCCTAACCCTAACCCTAACCTAACCCTAACCCTA\n"
	expected := make(map[string]uint32)
	expected["CCTAACCCTAACCCTA"] = 2
	expected["CTAACCCTAACCCTAA"] = 1
	expected["TAACCCTAACCCTAAC"] = 1
	expected["AACCCTAACCCTAACC"] = 1
	expected["ACCCTAACCCTAACCT"] = 1
	expected["CCCTAACCCTAACCTA"] = 1
	expected["CCTAACCCTAACCTAA"] = 1
	expected["CTAACCCTAACCTAAC"] = 1
	expected["TAACCCTAACCTAACC"] = 1
	expected["AACCCTAACCTAACCC"] = 1
	expected["ACCCTAACCTAACCCT"] = 1
	expected["CCCTAACCTAACCCTA"] = 1
	expected["CCTAACCTAACCCTAA"] = 1
	expected["CTAACCTAACCCTAAC"] = 1
	expected["TAACCTAACCCTAACC"] = 1
	expected["AACCTAACCCTAACCC"] = 1
	expected["ACCTAACCCTAACCCT"] = 1

	counts := GatherCommonSequences(bufio.NewReader(strings.NewReader(sequence)))
	assert.Equal(t, len(expected), len(counts))
	for hash, count := range counts {
		subSequence := reverseHash(hash)
		assert.Equal(t, expected[subSequence], count)
	}
}

func TestReverseHash(t *testing.T) {
	expected := "AGCTAGCTAGCTAGCT"
	actual := reverseHash(454761243)
	assert.Equal(t, expected, actual)
}

func TestCombineMaps(t *testing.T) {
	first := make(map[uint32]uint32)
	first[0] = 1
	first[1] = 10
	second := make(map[uint32]uint32)
	second[0] = 1
	second[2] = 5

	CombineMaps(first, second)
	assert.Equal(t, first[0], uint32(2))
	assert.Equal(t, first[1], uint32(10))
	assert.Equal(t, first[2], uint32(5))
}
