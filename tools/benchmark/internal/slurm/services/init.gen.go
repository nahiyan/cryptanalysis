package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	services3 "benchmark/internal/marshalling/services"
	services2 "benchmark/internal/random/services"
	do "github.com/samber/do"
)

type SlurmService struct {
	errorSvc       *services.ErrorService
	configSvc      *services1.ConfigService
	randomSvc      *services2.RandomService
	marshallingSvc *services3.MarshallingService
}

func NewSlurmService(injector *do.Injector) (*SlurmService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	randomSvc := do.MustInvoke[*services2.RandomService](injector)
	marshallingSvc := do.MustInvoke[*services3.MarshallingService](injector)
	svc := &SlurmService{errorSvc: errorSvc, configSvc: configSvc, randomSvc: randomSvc, marshallingSvc: marshallingSvc}
	return svc, nil
}
