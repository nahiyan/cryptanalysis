package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	do "github.com/samber/do"
)

type DatabaseService struct {
	errorSvc      *services.ErrorService
	configSvc     *services1.ConfigService
	filesystemSvc *services2.FilesystemService
	Properties
}

func NewDatabaseService(injector *do.Injector) (*DatabaseService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	svc := &DatabaseService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc}
	svc.Init()
	return svc, nil
}
