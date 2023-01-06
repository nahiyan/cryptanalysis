package services

import (
	services "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	do "github.com/samber/do"
)

type CubeSelectorService struct {
	errorSvc      *services.ErrorService
	filesystemSvc *services1.FilesystemService
}

func NewCubeSelectorService(injector *do.Injector) (*CubeSelectorService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	svc := &CubeSelectorService{errorSvc: errorSvc, filesystemSvc: filesystemSvc}
	return svc, nil
}
