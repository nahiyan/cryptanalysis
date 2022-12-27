package services

import (
	badger "github.com/dgraph-io/badger/v3"
)

type Properties struct {
	db *badger.DB
}

func (databaseSvc *DatabaseService) Init() {
	errorSvc := databaseSvc.errorSvc
	config := databaseSvc.configSvc.Config

	db, err := badger.Open(badger.DefaultOptions(config.Paths.Database).WithValueLogFileSize(1024 * 1024 * 1))
	errorSvc.Fatal(err, "Database: failed to open")
	databaseSvc.db = db
}

func (databaseSvc *DatabaseService) Set(bucket string, key []byte, value []byte) error {
	err := databaseSvc.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, value)
		return err
	})

	return err
}

func (databaseSvc *DatabaseService) Get(bucket string, key []byte) ([]byte, error) {
	value := []byte{}

	err := databaseSvc.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		err = item.Value(func(value_ []byte) error {
			value = value_

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	return value, err
}
