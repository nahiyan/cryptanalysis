package services

import (
	services "benchmark/internal/error/services"
	do "github.com/samber/do"
)

type FilesystemService struct {
	errorSvc *services.ErrorService
}

func NewFilesystemService(injector *do.Injector) (*FilesystemService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	svc := &FilesystemService{errorSvc: errorSvc}
	return svc, nil
}
