package services

import (
	services1 "cryptanalysis/internal/config/services"
	services "cryptanalysis/internal/error/services"
	services3 "cryptanalysis/internal/marshalling/services"
	services2 "cryptanalysis/internal/random/services"
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
