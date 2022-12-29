package services

import (
	services "benchmark/internal/error/services"
	do "github.com/samber/do"
)

type MarshallingService struct {
	errorSvc *services.ErrorService
}

func NewMarshallingService(injector *do.Injector) (*MarshallingService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	svc := &MarshallingService{errorSvc: errorSvc}
	return svc, nil
}
