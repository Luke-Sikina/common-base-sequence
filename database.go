package main

import (
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const databaseName = "counts.db"

var countsBucket = []byte("counts")
var dataTransformer = binary.BigEndian

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
