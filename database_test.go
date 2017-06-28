package main

import (
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestOpenDb(t *testing.T) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	_, err := os.Stat(dir + "/" + databaseName)
	assert.True(t, os.IsNotExist(err))
	db, err := openDb()
	defer db.Close()
	defer os.Remove(dir + "/" + databaseName)
	assert.True(t, err == nil)
	assert.True(t, db != nil)
}

func TestUpdateCounts(t *testing.T) {
	db, err := openDb()
	defer db.Close()
	defer os.Remove(db.Path())
	assert.True(t, err == nil, "Error opening the database. This is an error in test setup.")
	storedCounts := map[uint32]uint32{0: 10, 1: 20, 2345: 2345, 4294967295: 4294967295}
	updateDb(storedCounts, db)
	db.Close()
	db, err = openDb()
	assert.True(t, err == nil, "Error opening the database. This is an error in test setup.")

	newCounts := map[uint32]uint32{0: 5, 1: 20, 2: 22, 2345: 1111}
	updateCounts(newCounts, db)
	assert.Equal(t, uint32(15), newCounts[0])
	assert.Equal(t, uint32(40), newCounts[1])
	assert.Equal(t, uint32(22), newCounts[2])
	assert.Equal(t, uint32(3456), newCounts[2345])
}

func TestUpdateDb(t *testing.T) {
	db, err := openDb()
	defer db.Close()
	defer os.Remove(db.Path())
	assert.True(t, err == nil, "Error opening the database. This is an error in test setup.")
	providedCounts := map[uint32]uint32{0: 10, 1: 20, 2345: 2345, 4294967295: 4294967295}

	updateDb(providedCounts, db)
	returnedCounts := retriveCountsFromDb(db)

	for key, value := range providedCounts {
		assert.Equal(t, value, returnedCounts[key])
	}
	assert.Equal(t, len(providedCounts), len(returnedCounts))
}

func retriveCountsFromDb(db *bolt.DB) (counts map[uint32]uint32) {
	counts = make(map[uint32]uint32)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(countsBucket)
		bucket.ForEach(func(k, v []byte) error {
			key := dataTransformer.Uint32(k)
			count := dataTransformer.Uint32(v)
			counts[key] = count
			return nil
		})
		return nil
	})
	return
}

func TestStoreCounts(t *testing.T) {
	initialCounts := map[uint32]uint32{0: 10, 1: 20, 2345: 2345, 4294967295: 4294967295}
	StoreCounts(initialCounts)
	nextCounts := map[uint32]uint32{0: 10, 1: 20, 2: 30, 2345: 1111}
	StoreCounts(nextCounts)

	db, err := openDb()
	defer db.Close()
	defer os.Remove(db.Path())
	assert.True(t, err == nil, "Error opening the database. This is an error in test setup.")

	totalCounts := retriveCountsFromDb(db)
	assert.Equal(t, uint32(20), totalCounts[0])
	assert.Equal(t, uint32(40), totalCounts[1])
	assert.Equal(t, uint32(30), totalCounts[2])
	assert.Equal(t, uint32(3456), totalCounts[2345])
	assert.Equal(t, uint32(4294967295), totalCounts[4294967295])
}

func TestGetTopCounts(t *testing.T) {
	counts := map[uint32]uint32{0: 10, 1: 5, 3: 1, 5: 9, 7: 2}
	err := StoreCounts(counts)
	assert.True(t, err == nil, "Error opening the database. This is an error in test setup.")

	actualTop, err := GetTopCounts(3)
	assert.True(t, err == nil)
	expectedTop := []HashCountPair{{0, 10}, {5, 9}, {1, 5}}

	assert.Equal(t, len(expectedTop), len(actualTop))
	for index, value := range expectedTop {
		assert.Equal(t, value, actualTop[index])
	}
}
