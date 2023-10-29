package services

import (
	services "cryptanalysis/internal/error/services"
	services1 "cryptanalysis/internal/filesystem/services"
	services2 "cryptanalysis/internal/marshalling/services"
	do "github.com/samber/do"
)

type CubesetService struct {
	errorSvc       *services.ErrorService
	filesystemSvc  *services1.FilesystemService
	marshallingSvc *services2.MarshallingService
}

func NewCubesetService(injector *do.Injector) (*CubesetService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	marshallingSvc := do.MustInvoke[*services2.MarshallingService](injector)
	svc := &CubesetService{errorSvc: errorSvc, filesystemSvc: filesystemSvc, marshallingSvc: marshallingSvc}
	return svc, nil
}
