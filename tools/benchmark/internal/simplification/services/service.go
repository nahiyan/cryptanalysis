package services

import (
	"benchmark/internal/simplification"
)

type Properties struct {
	Bucket string
}

func (simplificationSvc *SimplificationService) Init() {
	simplificationSvc.Bucket = "simplifications"
}

func (simplificationSvc *SimplificationService) Register(encoding string, simplification simplification.Simplification) error {
	databaseSvc := simplificationSvc.databaseSvc
	filesystemSvc := simplificationSvc.filesystemSvc

	checksum, err := filesystemSvc.Checksum(encoding)
	if err != nil {
		return err
	}

	value, err := simplificationSvc.marshallingSvc.BinEncode(simplification)
	if err != nil {
		return err
	}

	key := []byte(checksum)
	if err := databaseSvc.Set(simplificationSvc.Bucket, key, value); err != nil {
		return err
	}

	return nil
}

func (simplificationSvc *SimplificationService) All() ([]simplification.Simplification, error) {
	simplifications := []simplification.Simplification{}
	simplificationSvc.databaseSvc.All(simplificationSvc.Bucket, func(key, value []byte) {
		var simplification simplification.Simplification
		if err := simplificationSvc.marshallingSvc.BinDecode(value, &simplification); err != nil {
			return
		}

		simplifications = append(simplifications, simplification)
	})

	return simplifications, nil
}
