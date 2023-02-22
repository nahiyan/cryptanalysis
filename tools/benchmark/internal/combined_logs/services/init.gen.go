package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	do "github.com/samber/do"
)

type CombinedLogsService struct {
	errorSvc      *services.ErrorService
	configSvc     *services1.ConfigService
	filesystemSvc *services2.FilesystemService
	Properties
}

func NewCombinedLogsService(injector *do.Injector) (*CombinedLogsService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	svc := &CombinedLogsService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc}
	return svc, nil
}
