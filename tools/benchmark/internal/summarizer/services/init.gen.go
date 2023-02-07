package services

import (
	services2 "benchmark/internal/cubeset/services"
	services4 "benchmark/internal/encoder/services"
	services "benchmark/internal/error/services"
	services3 "benchmark/internal/simplification/services"
	services1 "benchmark/internal/solution/services"
	do "github.com/samber/do"
)

type SummarizerService struct {
	errorSvc          *services.ErrorService
	solutionSvc       *services1.SolutionService
	cubesetSvc        *services2.CubesetService
	simplificationSvc *services3.SimplificationService
	encoderSvc        *services4.EncoderService
}

func NewSummarizerService(injector *do.Injector) (*SummarizerService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	solutionSvc := do.MustInvoke[*services1.SolutionService](injector)
	cubesetSvc := do.MustInvoke[*services2.CubesetService](injector)
	simplificationSvc := do.MustInvoke[*services3.SimplificationService](injector)
	encoderSvc := do.MustInvoke[*services4.EncoderService](injector)
	svc := &SummarizerService{errorSvc: errorSvc, solutionSvc: solutionSvc, cubesetSvc: cubesetSvc, simplificationSvc: simplificationSvc, encoderSvc: encoderSvc}
	return svc, nil
}
