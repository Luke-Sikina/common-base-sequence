package main

import (
	"bufio"
	"io"
	"log"
	"unicode"
)

type SequenceGatherer struct {
	stream *bufio.Reader
	batch  uint32
}

type Hash uint32

const (
	A = iota
	G
	C
	T
	max16BitU    uint16 = 0xFFFF
	batchBitMask uint32 = 0xFF000000
)

//Gathers counts for sequences with hash within range batch <= hash < batch + 0x01000000
func (gatherer *SequenceGatherer) gatherSequences() (counts map[Hash]uint32) {
	counts = make(map[Hash]uint32)
	var validTokenSequence uint16 = 0
	var hash Hash = 0
	var c rune
	var err error
	for err != io.EOF {
		if c, _, err = gatherer.stream.ReadRune(); err != nil && err != io.EOF {
			log.Fatal(err)
		} else if !unicode.IsSpace(c) {
			switch c {
			case 'A':
				hash = (hash << 2) + A
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'G':
				hash = (hash << 2) + G
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'C':
				hash = (hash << 2) + C
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'T':
				hash = (hash << 2) + T
				validTokenSequence = (validTokenSequence << 1) + 1
			default:
				validTokenSequence = validTokenSequence << 1
			}
			if validTokenSequence == max16BitU && hash.inBatch(gatherer.batch) {
				counts[hash]++
			}
		}
	}
	return
}

func (hash Hash) inBatch(batch uint32) bool {
	return uint32(hash)&batchBitMask == batch
}

func (hash Hash) reverse() (bases string) {
	tmp := hash
	for i := 0; i < 16; i++ {
		switch tmp % 4 {
		case A:
			bases = "A" + bases
		case G:
			bases = "G" + bases
		case C:
			bases = "C" + bases
		case T:
			bases = "T" + bases
		}
		tmp = tmp >> 2
	}
	return
}
