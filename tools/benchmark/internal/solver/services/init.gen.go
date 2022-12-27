package services

import do "github.com/samber/do"

type SolverService struct{}

func NewSolverService(injector *do.Injector) (*SolverService, error) {
	svc := &SolverService{}
	return svc, nil
}
