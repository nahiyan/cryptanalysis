package services

import (
	"benchmark/internal/cubeset"
	"fmt"
	"time"
)

type Properties struct {
	Bucket string
}

func (cubesetSvc *CubesetService) Init() {
	cubesetSvc.Bucket = "cubesets"
}

func (cubesetSvc *CubesetService) Register(cubesetFilePath string, cubeSet cubeset.CubeSet) error {
	checksum, err := cubesetSvc.filesystemSvc.Checksum(cubesetFilePath)
	if err != nil {
		return err
	}

	data, err := cubesetSvc.marshallingSvc.BinEncode(cubeSet)
	if err != nil {
		return err
	}

	startTime := time.Now()
	err = cubesetSvc.databaseSvc.Set(cubesetSvc.Bucket, []byte(checksum), data)
	fmt.Println("Stored", cubesetFilePath, checksum, time.Since(startTime))
	return err
}

func (cubesetSvc *CubesetService) All() ([]cubeset.CubeSet, error) {
	cubesets := []cubeset.CubeSet{}
	cubesetSvc.databaseSvc.All(cubesetSvc.Bucket, func(key, value []byte) {
		var cubeset cubeset.CubeSet
		if err := cubesetSvc.marshallingSvc.BinDecode(value, &cubeset); err != nil {
			return
		}

		cubesets = append(cubesets, cubeset)
	})

	return cubesets, nil
}
