package services

import "benchmark/internal/cubeset"

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

	err = cubesetSvc.databaseSvc.Set(cubesetSvc.Bucket, []byte(checksum), data)
	return err
}
