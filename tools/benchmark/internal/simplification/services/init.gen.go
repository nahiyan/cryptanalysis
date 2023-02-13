package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	services3 "benchmark/internal/marshalling/services"
	do "github.com/samber/do"
)

type SimplificationService struct {
	errorSvc       *services.ErrorService
	configSvc      *services1.ConfigService
	filesystemSvc  *services2.FilesystemService
	marshallingSvc *services3.MarshallingService
}

func NewSimplificationService(injector *do.Injector) (*SimplificationService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	marshallingSvc := do.MustInvoke[*services3.MarshallingService](injector)
	svc := &SimplificationService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc, marshallingSvc: marshallingSvc}
	return svc, nil
}
