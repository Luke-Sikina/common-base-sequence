package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		iterateOverFastaFiles(args[0])
	} else {
		log.Fatal("Incorrect number of args. Need a file to read from.")
	}
}

func iterateOverFastaFiles(dir string) {
	counts := make(map[uint32]uint32)
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		name := dir + "/" + f.Name()
		if strings.HasSuffix(name, ".fa") {
			log.Print("Finding sequences in file: " + name)
			toRead, err := os.Open(name)
			if err == nil {
				CombineMaps(counts, GatherCommonSequences(bufio.NewReader(toRead)))
			} else {
				log.Fatal("Error opening file: " + name + ", ignoring.")
				log.Fatal(err)
			}
		}
	}

	var max uint32
	var sequence string
	for key, value := range counts {
		if max < value {
			max = value
			sequence = reverseHash(key)
		}
	}

	println("Most frequenct sequence, " + sequence + ", occurs " + strconv.Itoa(int(max)) + " times.")
}

func GatherCommonSequences(stream *bufio.Reader) (frequencies map[uint32]uint32) {
	var validTokenSequence uint16 = 0
	var hash uint32 = 0
	frequencies = make(map[uint32]uint32)
	for {
		if c, _, err := stream.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
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

func reverseHash(hash uint32) (bases string) {
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

func CombineMaps(target, toAdd map[uint32]uint32) {
	for key, valueToAdd := range toAdd {
		target[key] += valueToAdd
	}
}
