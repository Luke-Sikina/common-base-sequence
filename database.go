package main

import (
	"container/heap"
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"path/filepath"
)

const databaseName = "counts.db"

var countsBucket = []byte("counts")
var dataTransformer = binary.BigEndian

func Clear() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	_, err := os.Stat(dir + "/" + databaseName)
	if os.IsExist(err) {
		os.Remove(dir + "/" + databaseName)
	}
}

func openDb() (*bolt.DB, error) {
	db, err := bolt.Open(databaseName, 0644, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
func StoreCounts(counts map[uint32]uint32) error {
	db, err := openDb()
	defer db.Close()
	err = updateCounts(counts, db)
	err = updateDb(counts, db)

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func updateCounts(counts map[uint32]uint32, db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(countsBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", countsBucket)
		}

		bucket.ForEach(func(k, v []byte) error {
			sequence := dataTransformer.Uint32(k)
			additionalCount := dataTransformer.Uint32(v)
			oldCount, exists := counts[sequence]
			if exists {
				counts[sequence] = oldCount + additionalCount
			}
			return nil
		})
		return nil
	})
}

func updateDb(counts map[uint32]uint32, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(countsBucket)
		if err != nil {
			return err
		}
		for key, value := range counts {
			keySlice, valueSlice := make([]byte, 4), make([]byte, 4)
			dataTransformer.PutUint32(keySlice, key)
			dataTransformer.PutUint32(valueSlice, value)
			err = bucket.Put(keySlice, valueSlice)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetTopCounts(mapSize int) (top []HashCountPair, err error) {
	db, err := openDb()
	defer db.Close()
	if err != nil {
		return
	}

	top = make([]HashCountPair, mapSize)

	h := &HashCountPairHeap{}
	heap.Init(h)

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(countsBucket)
		bucket.ForEach(func(k, v []byte) error {
			hash := dataTransformer.Uint32(k)
			count := dataTransformer.Uint32(v)
			heap.Push(h, HashCountPair{hash, count})
			if mapSize < h.Len() {
				heap.Pop(h)
			}
			return nil
		})
		return nil
	})

	for i := 0; i < mapSize; i++ {
		top[i] = h.Pop().(HashCountPair)
	}
	return
}
