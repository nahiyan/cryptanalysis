package services

import (
	services2 "benchmark/internal/cuber/services"
	services "benchmark/internal/encoder/services"
	services1 "benchmark/internal/solver/services"
	do "github.com/samber/do"
)

type PipelineService struct {
	encoderSvc *services.EncoderService
	solverSvc  *services1.SolverService
	cuberSvc   *services2.CuberService
	Properties
}

func NewPipelineService(injector *do.Injector) (*PipelineService, error) {
	encoderSvc := do.MustInvoke[*services.EncoderService](injector)
	solverSvc := do.MustInvoke[*services1.SolverService](injector)
	cuberSvc := do.MustInvoke[*services2.CuberService](injector)
	svc := &PipelineService{encoderSvc: encoderSvc, solverSvc: solverSvc, cuberSvc: cuberSvc}
	return svc, nil
}
