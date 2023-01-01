package services

import (
	services1 "benchmark/internal/database/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	services3 "benchmark/internal/marshalling/services"
	do "github.com/samber/do"
)

type CubesetService struct {
	errorSvc       *services.ErrorService
	databaseSvc    *services1.DatabaseService
	filesystemSvc  *services2.FilesystemService
	marshallingSvc *services3.MarshallingService
	Properties
}

func NewCubesetService(injector *do.Injector) (*CubesetService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	databaseSvc := do.MustInvoke[*services1.DatabaseService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	marshallingSvc := do.MustInvoke[*services3.MarshallingService](injector)
	svc := &CubesetService{errorSvc: errorSvc, databaseSvc: databaseSvc, filesystemSvc: filesystemSvc, marshallingSvc: marshallingSvc}
	svc.Init()
	return svc, nil
}
