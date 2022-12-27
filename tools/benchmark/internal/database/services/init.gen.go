package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	do "github.com/samber/do"
)

type DatabaseService struct {
	errorSvc  *services.ErrorService
	configSvc *services1.ConfigService
	Properties
}

func NewDatabaseService(injector *do.Injector) (*DatabaseService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	svc := &DatabaseService{errorSvc: errorSvc, configSvc: configSvc}
	svc.Init()
	return svc, nil
}
