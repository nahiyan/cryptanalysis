package services

import (
	services1 "benchmark/internal/config/services"
	services3 "benchmark/internal/database/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	services4 "benchmark/internal/marshalling/services"
	do "github.com/samber/do"
)

type SimplificationService struct {
	errorSvc       *services.ErrorService
	configSvc      *services1.ConfigService
	filesystemSvc  *services2.FilesystemService
	databaseSvc    *services3.DatabaseService
	marshallingSvc *services4.MarshallingService
	Properties
}

func NewSimplificationService(injector *do.Injector) (*SimplificationService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	databaseSvc := do.MustInvoke[*services3.DatabaseService](injector)
	marshallingSvc := do.MustInvoke[*services4.MarshallingService](injector)
	svc := &SimplificationService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc, databaseSvc: databaseSvc, marshallingSvc: marshallingSvc}
	svc.Init()
	return svc, nil
}
