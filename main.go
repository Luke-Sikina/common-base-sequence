package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

const heapSize = 100

//Batch based attempt: break into batches of hashes, since hashes have a finite numerical range
//For each batch, make the counts map, then take the top 100 and merge that list with the overall one,
//Keeping the overall top 100
func iterateOverFastaFiles(dir string) {
	hashHeep := &HashCountPairHeap{}
	heap.Init(hashHeep)
	for i := 0; i < 255; i++ {
		batch := uint32(i) << 24
		log.Printf("Finding sequences with hashes where %d <= hash < %d", batch, batch+0x01000000)
		files, _ := ioutil.ReadDir(dir)
		for _, f := range files {
			name := dir + "/" + f.Name()
			if strings.HasSuffix(name, ".fa") {
				log.Print("Finding sequences in file: " + name)
				toRead, err := os.Open(name)
				if err == nil {
					counts := (&SequenceGatherer{bufio.NewReader(toRead), batch}).gatherSequences()
					addCountMapToHeap(counts, hashHeep, heapSize)
				} else {
					log.Fatal("Error opening file: " + name + ", ignoring.")
					log.Fatal(err)
				}
			}
		}
	}
	fmt.Printf("Top %d sequences, sorted least to most common:\n", heapSize)
	for hashHeep.Len() > 0 {
		pair := hashHeep.Pop().(HashCountPair)
		fmt.Printf("%s: %d\n", pair.hash.reverse(), pair.count)
	}
}

func addCountMapToHeap(counts map[Hash]uint32, h *HashCountPairHeap, size int) {
	for key, value := range counts {
		heap.Push(h, HashCountPair{key, value})
		if h.Len() > size {
			heap.Pop(h)
		}
	}
}
