package services

import (
	services "cryptanalysis/internal/error/services"
	do "github.com/samber/do"
)

type RandomService struct {
	errorSvc *services.ErrorService
}

func NewRandomService(injector *do.Injector) (*RandomService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	svc := &RandomService{errorSvc: errorSvc}
	return svc, nil
}
