package services

import (
	services "benchmark/internal/error/services"
	do "github.com/samber/do"
)

type CuberService struct {
	errorSvc *services.ErrorService
}

func NewCuberService(injector *do.Injector) (*CuberService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	svc := &CuberService{errorSvc: errorSvc}
	return svc, nil
}
