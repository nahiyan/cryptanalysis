package services

import (
	services2 "benchmark/internal/cubeset/services"
	services "benchmark/internal/error/services"
	services3 "benchmark/internal/simplification/services"
	services1 "benchmark/internal/solution/services"
	do "github.com/samber/do"
)

type LogService struct {
	errorSvc          *services.ErrorService
	solutionSvc       *services1.SolutionService
	cubesetSvc        *services2.CubesetService
	simplificationSvc *services3.SimplificationService
}

func NewLogService(injector *do.Injector) (*LogService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	solutionSvc := do.MustInvoke[*services1.SolutionService](injector)
	cubesetSvc := do.MustInvoke[*services2.CubesetService](injector)
	simplificationSvc := do.MustInvoke[*services3.SimplificationService](injector)
	svc := &LogService{errorSvc: errorSvc, solutionSvc: solutionSvc, cubesetSvc: cubesetSvc, simplificationSvc: simplificationSvc}
	return svc, nil
}
