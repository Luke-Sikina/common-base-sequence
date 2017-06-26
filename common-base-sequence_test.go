package main

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGatherCommonSequences(t *testing.T) {
	sequence := "CCTAACCCTAACCCTAACCTAACCCTAACCCTA\n"
	expected := make(map[string]uint)
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
	for hash, count := range counts {
		subSequence := ReverseHash(hash)
		assert.Equal(t, expected[subSequence], count)
	}
}

func TestReverseHash(t *testing.T) {
	expected := "AGCTAGCTAGCTAGCT"
	actual := ReverseHash(454761243)
	assert.Equal(t, expected, actual)
}
