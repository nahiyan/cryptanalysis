package services

import (
	errorModule "benchmark/internal/error"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

type Properties struct {
	db *bolt.DB
}

func (databaseSvc *DatabaseService) CreateBuckets() error {
	buckets := []string{"solutions", "tasks"}

	err := databaseSvc.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (databaseSvc *DatabaseService) Init() {
	errorSvc := databaseSvc.errorSvc
	db, err := bolt.Open("database.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	errorSvc.Fatal(err, "Database: failed to open")
	databaseSvc.db = db

	{
		err := databaseSvc.CreateBuckets()
		errorSvc.Fatal(err, "Database: failed to create buckets")
	}
}

func (databaseSvc *DatabaseService) Set(bucket string, key []byte, value []byte) error {
	err := databaseSvc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(key), []byte(value))
		return err
	})

	return err
}

func (databaseSvc *DatabaseService) Get(bucket string, key []byte) ([]byte, error) {
	var value []byte
	err := databaseSvc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		value_ := b.Get([]byte(key))

		// Treat empty value as non-existent
		if len(value_) == 0 {
			value = []byte{}
			return errorModule.ErrKeyNotFound
		}

		// Copy the byte slice since BoltDB is a zero-copy database and the memory allocation may get reclaimed
		value = make([]byte, len(value_))
		copy(value, value_)

		return nil
	})

	return value, err
}
