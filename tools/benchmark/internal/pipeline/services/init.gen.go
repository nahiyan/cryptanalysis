package services

import (
	services "benchmark/internal/encoder/services"
	services1 "benchmark/internal/solver/services"
	do "github.com/samber/do"
)

type PipelineService struct {
	encoderSvc *services.EncoderService
	solverSvc  *services1.SolverService
	Properties
}

func NewPipelineService(injector *do.Injector) (*PipelineService, error) {
	encoderSvc := do.MustInvoke[*services.EncoderService](injector)
	solverSvc := do.MustInvoke[*services1.SolverService](injector)
	svc := &PipelineService{encoderSvc: encoderSvc, solverSvc: solverSvc}
	return svc, nil
}
