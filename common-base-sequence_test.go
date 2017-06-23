package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	sequence := Sequence{'A', 'C', 'G', 'T', 'A', 'C', 'G', 'T', 'A', 'C', 'G', 'T', 'A', 'C', 'G', 'T'}
	actual := sequence.hash()
	assert.Equal(t, uint32(656877351), actual)
}
