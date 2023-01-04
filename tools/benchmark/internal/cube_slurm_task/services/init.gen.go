package services

import (
	services1 "benchmark/internal/database/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/marshalling/services"
	do "github.com/samber/do"
)

type CubeSlurmTaskService struct {
	errorSvc       *services.ErrorService
	databaseSvc    *services1.DatabaseService
	marshallingSvc *services2.MarshallingService
	Properties
}

func NewCubeSlurmTaskService(injector *do.Injector) (*CubeSlurmTaskService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	databaseSvc := do.MustInvoke[*services1.DatabaseService](injector)
	marshallingSvc := do.MustInvoke[*services2.MarshallingService](injector)
	svc := &CubeSlurmTaskService{errorSvc: errorSvc, databaseSvc: databaseSvc, marshallingSvc: marshallingSvc}
	svc.Init()
	return svc, nil
}
