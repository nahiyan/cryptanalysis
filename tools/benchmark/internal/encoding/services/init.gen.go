package services

import (
	services "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	do "github.com/samber/do"
)

type EncodingService struct {
	errorSvc      *services.ErrorService
	filesystemSvc *services1.FilesystemService
}

func NewEncodingService(injector *do.Injector) (*EncodingService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	svc := &EncodingService{errorSvc: errorSvc, filesystemSvc: filesystemSvc}
	return svc, nil
}
