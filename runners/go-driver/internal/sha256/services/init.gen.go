package services

import (
	services1 "cryptanalysis/internal/config/services"
	services "cryptanalysis/internal/error/services"
	do "github.com/samber/do"
)

type Sha256Service struct {
	errorSvc  *services.ErrorService
	configSvc *services1.ConfigService
}

func NewSha256Service(injector *do.Injector) (*Sha256Service, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	svc := &Sha256Service{errorSvc: errorSvc, configSvc: configSvc}
	return svc, nil
}
