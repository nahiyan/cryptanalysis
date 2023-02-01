package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	do "github.com/samber/do"
)

type Md4Service struct {
	errorSvc  *services.ErrorService
	configSvc *services1.ConfigService
}

func NewMd4Service(injector *do.Injector) (*Md4Service, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	svc := &Md4Service{errorSvc: errorSvc, configSvc: configSvc}
	return svc, nil
}
