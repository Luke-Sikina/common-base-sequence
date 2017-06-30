package main

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type ReverseHashCase struct {
	hash     Hash
	expected string
}

var reverseHashCases = []*ReverseHashCase{
	{0, "AAAAAAAAAAAAAAAA"},
	{0xFFFFFFFF, "TTTTTTTTTTTTTTTT"},
	{454761243, "AGCTAGCTAGCTAGCT"},
}

func TestReverseHash(t *testing.T) {
	for _, testCase := range reverseHashCases {
		actual := testCase.hash.reverse()
		assert.Equal(t, testCase.expected, actual)
	}
}

type InBatchCase struct {
	hash     Hash
	batch    uint32
	expected bool
}

var InBatchCases = []*InBatchCase{
	{0x00000000, 0x00000000, true},
	{0x00FFFFFF, 0x00000000, true},
	{0x01000000, 0x00000000, false},
	{0x3A000000, 0x3A000000, true},
	{0x3AFFFFFF, 0x3A000000, true},
	{0x1A000000, 0x3A000000, false},
	{0x2AFFFFFF, 0x3A000000, false},
	{0x4A000000, 0x3A000000, false},
	{0xFF000000, 0xFF000000, true},
	{0xFFFFFFFF, 0xFF000000, true},
	{0xFEFFFFFF, 0xFF000000, false},
}

func TestInBatch(t *testing.T) {
	for _, testCase := range InBatchCases {
		actual := testCase.hash.inBatch(testCase.batch)
		assert.Equal(t, testCase.expected, actual)
	}
}

type GatherSequencesCase struct {
	sequence string
	batch    uint32
	expected map[Hash]uint32
}

var gatherSequencesCases = []*GatherSequencesCase{
	{"AAAAAAAAAAAAAAAAAAAAA", 0, map[Hash]uint32{0: 6}},
	{"TTTTTTTTTTTTTTTTTTTTT", 0, map[Hash]uint32{}},
	{"AAAATTTTTTTTTTTTTTTTT", 0, map[Hash]uint32{0x00FFFFFF: 1}},
	{"AAAAATTTTTTTTTTTTTTTT", 0, map[Hash]uint32{0x003FFFFF: 1, 0x00FFFFFF: 1}},
}

func TestGatherSequences(t *testing.T) {
	for _, testCase := range gatherSequencesCases {
		reader := bufio.NewReader(strings.NewReader(testCase.sequence))
		gatherer := SequenceGatherer{reader, testCase.batch}
		actual := gatherer.gatherSequences()
		assert.Equal(t, len(testCase.expected), len(actual))
		for key, value := range testCase.expected {
			assert.Equal(t, value, actual[key])
		}
	}
}
