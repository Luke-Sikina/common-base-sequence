package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		Clear()
		iterateOverFastaFiles(args[0])
		top, err := GetTopCounts(10)
		if err == nil {
			log.Fatal("Error retrieving counts from database.")
		} else {
			for _, pair := range top {
				println("Sequence: " + reverseHash(pair.hash()) + " count: " + strconv.Itoa(int(pair.count())))
			}
		}
	} else {
		log.Fatal("Incorrect number of args. Need a file to read from.")
	}
}

func iterateOverFastaFiles(dir string) {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		name := dir + "/" + f.Name()
		if strings.HasSuffix(name, ".fa") {
			log.Print("Finding sequences in file: " + name)
			toRead, err := os.Open(name)
			if err == nil {
				StoreCounts(GatherCommonSequences(bufio.NewReader(toRead)))
			} else {
				log.Fatal("Error opening file: " + name + ", ignoring.")
				log.Fatal(err)
			}
		}
	}
}

const (
	max16BitU uint16 = 65535
	A                = iota
	G
	C
	T
)

func GatherCommonSequences(stream *bufio.Reader) (counts map[uint32]uint32) {
	counts = make(map[uint32]uint32)
	var validTokenSequence uint16 = 0
	var hash uint32 = 0
	var c rune
	var err error
	for err != io.EOF {
		if c, _, err = stream.ReadRune(); err != nil && err != io.EOF {
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
			if validTokenSequence == max16BitU {
				counts[hash]++
			}
		}
	}
	return
}

func reverseHash(hash uint32) (bases string) {
	for i := 0; i < 16; i++ {
		switch hash % 4 {
		case A:
			bases = "A" + bases
		case G:
			bases = "G" + bases
		case C:
			bases = "C" + bases
		case T:
			bases = "T" + bases
		}
		hash = hash >> 2
	}
	return
}
