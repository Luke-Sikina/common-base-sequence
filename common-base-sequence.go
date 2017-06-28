package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
				fmt.Printf("Sequence: %s count: %d", reverseHash(pair.hash()), int(pair.count()))
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
			log.Printf("Finding sequences in file: %s", name)
			toRead, err := os.Open(name)
			if err == nil {
				log.Print("Storing sequences found in file: " + name)
				StoreCounts(gatherCommonSequences(bufio.NewReader(toRead)))
			} else {
				log.Fatalf("Error opening file: %s, ignoring.", name)
				log.Fatal(err)
			}
		}
	}
}

const (
	max16BitU uint16 = 65535
	a                = iota
	g
	c
	t
)

func gatherCommonSequences(stream *bufio.Reader) (counts map[uint32]uint32) {
	counts = make(map[uint32]uint32)
	var validTokenSequence uint16
	var hash uint32
	var char rune
	var err error
	for err != io.EOF {
		if char, _, err = stream.ReadRune(); err != nil && err != io.EOF {
			log.Fatal(err)
		} else if !unicode.IsSpace(char) {
			switch char {
			case 'A':
				hash = (hash << 2) + a
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'G':
				hash = (hash << 2) + g
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'C':
				hash = (hash << 2) + c
				validTokenSequence = (validTokenSequence << 1) + 1
			case 'T':
				hash = (hash << 2) + t
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
		case a:
			bases = "A" + bases
		case g:
			bases = "G" + bases
		case c:
			bases = "C" + bases
		case t:
			bases = "T" + bases
		}
		hash = hash >> 2
	}
	return
}
