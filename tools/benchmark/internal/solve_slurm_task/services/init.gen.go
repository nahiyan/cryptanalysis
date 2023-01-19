package services

import (
	services1 "benchmark/internal/database/services"
	services "benchmark/internal/error/services"
	services3 "benchmark/internal/filesystem/services"
	services2 "benchmark/internal/marshalling/services"
	do "github.com/samber/do"
)

type SolveSlurmTaskService struct {
	errorSvc       *services.ErrorService
	databaseSvc    *services1.DatabaseService
	marshallingSvc *services2.MarshallingService
	filesystemSvc  *services3.FilesystemService
	Properties
}

func NewSolveSlurmTaskService(injector *do.Injector) (*SolveSlurmTaskService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	databaseSvc := do.MustInvoke[*services1.DatabaseService](injector)
	marshallingSvc := do.MustInvoke[*services2.MarshallingService](injector)
	filesystemSvc := do.MustInvoke[*services3.FilesystemService](injector)
	svc := &SolveSlurmTaskService{errorSvc: errorSvc, databaseSvc: databaseSvc, marshallingSvc: marshallingSvc, filesystemSvc: filesystemSvc}
	svc.Init()
	return svc, nil
}
