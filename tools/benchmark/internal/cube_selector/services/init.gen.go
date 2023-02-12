package services

import (
	services2 "benchmark/internal/config/services"
	services3 "benchmark/internal/cubeset/services"
	services "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	do "github.com/samber/do"
)

type CubeSelectorService struct {
	errorSvc      *services.ErrorService
	filesystemSvc *services1.FilesystemService
	configSvc     *services2.ConfigService
	cubesetSvc    *services3.CubesetService
}

func NewCubeSelectorService(injector *do.Injector) (*CubeSelectorService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	configSvc := do.MustInvoke[*services2.ConfigService](injector)
	cubesetSvc := do.MustInvoke[*services3.CubesetService](injector)
	svc := &CubeSelectorService{errorSvc: errorSvc, filesystemSvc: filesystemSvc, configSvc: configSvc, cubesetSvc: cubesetSvc}
	return svc, nil
}
