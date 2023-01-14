package services

import (
	errorModule "benchmark/internal/error"
	"strconv"

	"github.com/boltdb/bolt"
)

type Properties struct {
	db *bolt.DB
}

func (databaseSvc *DatabaseService) CreateBuckets() error {
	buckets := []string{"solutions", "tasks", "cubesets", "simplifications"}

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

func (databaseSvc *DatabaseService) Open(readOnly bool) error {
	db, err := bolt.Open("database.db", 0600, &bolt.Options{Timeout: 0, ReadOnly: readOnly})
	if err != nil {
		return err
	}
	databaseSvc.db = db

	return nil
}

func (databaseSvc *DatabaseService) Close() error {
	if err := databaseSvc.db.Close(); err != nil {
		return err
	}

	return nil
}

func (databaseSvc *DatabaseService) Use(isReadOnly bool, handler func(db *bolt.DB) error) error {
	if err := databaseSvc.Open(isReadOnly); err != nil {
		return err
	}
	defer databaseSvc.Close()

	if err := handler(databaseSvc.db); err != nil {
		return err
	}

	return nil
}

func (databaseSvc *DatabaseService) UseReadOnly(handler func(db *bolt.DB) error) error {
	return databaseSvc.Use(true, handler)
}

func (databaseSvc *DatabaseService) UseReadWrite(handler func(db *bolt.DB) error) error {
	return databaseSvc.Use(false, handler)
}

func (databaseSvc *DatabaseService) Init() {
	errorSvc := databaseSvc.errorSvc

	// Open the database
	err := databaseSvc.Open(false)
	errorSvc.Fatal(err, "Database: failed to open")

	// Buckets
	err = databaseSvc.CreateBuckets()
	errorSvc.Fatal(err, "Database: failed to create buckets")

	err = databaseSvc.Close()
	errorSvc.Fatal(err, "Database: failed to close")
}

func (databaseSvc *DatabaseService) Get(bucket string, key []byte) ([]byte, error) {
	var value []byte
	err := databaseSvc.UseReadOnly(func(db *bolt.DB) error {
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

		return err
	})

	return value, err
}

func (databaseSvc *DatabaseService) Set(bucket string, key []byte, value []byte) error {
	err := databaseSvc.UseReadWrite(func(db *bolt.DB) error {
		err := databaseSvc.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))

			// Set key to auto-incrementing ID if empty
			if key == nil {
				id, err := b.NextSequence()
				if err != nil {
					return err
				}

				key = []byte(strconv.Itoa(int(id)))
			}

			err := b.Put([]byte(key), []byte(value))
			return err
		})

		return err
	})

	return err
}

func (databaseSvc *DatabaseService) All(bucket string, handler func(key, value []byte)) error {
	err := databaseSvc.UseReadOnly(func(db *bolt.DB) error {
		err := databaseSvc.db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			bucket := tx.Bucket([]byte(bucket))
			err := bucket.ForEach(func(key, value []byte) error {
				handler(key, value)
				return nil
			})
			return err
		})

		return err
	})

	return err
}

func (databaseSvc *DatabaseService) RemoveAll(bucket string) error {
	err := databaseSvc.UseReadWrite(func(db *bolt.DB) error {
		err := databaseSvc.db.Update(func(tx *bolt.Tx) error {
			if err := tx.DeleteBucket([]byte(bucket)); err != nil {
				return err
			}

			_, err := tx.CreateBucket([]byte(bucket))
			return err
		})

		return err
	})

	return err
}
