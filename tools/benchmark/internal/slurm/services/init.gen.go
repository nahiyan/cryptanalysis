package services

import (
	services2 "benchmark/internal/config/services"
	services1 "benchmark/internal/database/services"
	services "benchmark/internal/error/services"
	services4 "benchmark/internal/marshalling/services"
	services3 "benchmark/internal/random/services"
	do "github.com/samber/do"
)

type SlurmService struct {
	errorSvc       *services.ErrorService
	databaseSvc    *services1.DatabaseService
	configSvc      *services2.ConfigService
	randomSvc      *services3.RandomService
	marshallingSvc *services4.MarshallingService
	Properties
}

func NewSlurmService(injector *do.Injector) (*SlurmService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	databaseSvc := do.MustInvoke[*services1.DatabaseService](injector)
	configSvc := do.MustInvoke[*services2.ConfigService](injector)
	randomSvc := do.MustInvoke[*services3.RandomService](injector)
	marshallingSvc := do.MustInvoke[*services4.MarshallingService](injector)
	svc := &SlurmService{errorSvc: errorSvc, databaseSvc: databaseSvc, configSvc: configSvc, randomSvc: randomSvc, marshallingSvc: marshallingSvc}
	svc.Init()
	return svc, nil
}
