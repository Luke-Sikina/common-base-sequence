package main

import (
	"bufio"
	"io"
	"log"
	"strconv"
)

func GatherCommonSequences(stream *bufio.Reader) (frequencies map[uint32]uint) {
	var validTokenSequence uint16 = 0
	var hash uint32 = 0
	frequencies = make(map[uint32]uint)
	for {
		if c, _, err := stream.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			log.Print(strconv.QuoteRune(c))
			skip := false
			switch c {
			case 'A':
				hash = hash << 2
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'G':
				hash = (hash << 2) + 1
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'C':
				hash = (hash << 2) + 2
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'T':
				hash = (hash << 2) + 3
				validTokenSequence = (validTokenSequence << 1) + 1
			case '\n':
				skip = true
			default:
				validTokenSequence = validTokenSequence << 1
			}
			if validTokenSequence == 65535 && !skip {
				frequencies[hash]++
			}
		}
	}
	return
}

func ReverseHash(hash uint32) (bases string) {
	for i := 0; i < 16; i++ {
		switch hash % 4 {
		case 0:
			bases = "A" + bases
		case 1:
			bases = "G" + bases
		case 2:
			bases = "C" + bases
		case 3:
			bases = "T" + bases
		}
		hash = hash >> 2
	}
	return
}
