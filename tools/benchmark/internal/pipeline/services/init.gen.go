package services

import (
	services3 "benchmark/internal/cube_selector/services"
	services2 "benchmark/internal/cuber/services"
	services "benchmark/internal/encoder/services"
	services1 "benchmark/internal/solver/services"
	do "github.com/samber/do"
)

type PipelineService struct {
	encoderSvc      *services.EncoderService
	solverSvc       *services1.SolverService
	cuberSvc        *services2.CuberService
	cubeSelectorSvc *services3.CubeSelectorService
}

func NewPipelineService(injector *do.Injector) (*PipelineService, error) {
	encoderSvc := do.MustInvoke[*services.EncoderService](injector)
	solverSvc := do.MustInvoke[*services1.SolverService](injector)
	cuberSvc := do.MustInvoke[*services2.CuberService](injector)
	cubeSelectorSvc := do.MustInvoke[*services3.CubeSelectorService](injector)
	svc := &PipelineService{encoderSvc: encoderSvc, solverSvc: solverSvc, cuberSvc: cuberSvc, cubeSelectorSvc: cubeSelectorSvc}
	return svc, nil
}
