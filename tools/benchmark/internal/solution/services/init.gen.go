package services

import (
	services2 "benchmark/internal/config/services"
	services1 "benchmark/internal/database/services"
	services "benchmark/internal/error/services"
	services3 "benchmark/internal/filesystem/services"
	services4 "benchmark/internal/marshalling/services"
	do "github.com/samber/do"
)

type SolutionService struct {
	errorSvc       *services.ErrorService
	databaseSvc    *services1.DatabaseService
	configSvc      *services2.ConfigService
	filesystemSvc  *services3.FilesystemService
	marshallingSvc *services4.MarshallingService
	Properties
}

func NewSolutionService(injector *do.Injector) (*SolutionService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	databaseSvc := do.MustInvoke[*services1.DatabaseService](injector)
	configSvc := do.MustInvoke[*services2.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services3.FilesystemService](injector)
	marshallingSvc := do.MustInvoke[*services4.MarshallingService](injector)
	svc := &SolutionService{errorSvc: errorSvc, databaseSvc: databaseSvc, configSvc: configSvc, filesystemSvc: filesystemSvc, marshallingSvc: marshallingSvc}
	svc.Init()
	return svc, nil
}
