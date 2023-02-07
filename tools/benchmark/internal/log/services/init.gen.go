package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	do "github.com/samber/do"
)

type LogService struct {
	errorSvc      *services.ErrorService
	configSvc     *services1.ConfigService
	filesystemSvc *services2.FilesystemService
	Properties
}

func NewLogService(injector *do.Injector) (*LogService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	svc := &LogService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc}
	svc.Init()
	return svc, nil
}
