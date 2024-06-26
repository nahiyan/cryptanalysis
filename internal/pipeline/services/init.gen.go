package services

import (
	services3 "cryptanalysis/internal/cube_selector/services"
	services2 "cryptanalysis/internal/cuber/services"
	services "cryptanalysis/internal/encoder/services"
	services4 "cryptanalysis/internal/simplifier/services"
	services5 "cryptanalysis/internal/slurm/services"
	services1 "cryptanalysis/internal/solver/services"
	do "github.com/samber/do"
)

type PipelineService struct {
	encoderSvc      *services.EncoderService
	solverSvc       *services1.SolverService
	cuberSvc        *services2.CuberService
	cubeSelectorSvc *services3.CubeSelectorService
	simplifierSvc   *services4.SimplifierService
	slurmSvc        *services5.SlurmService
}

func NewPipelineService(injector *do.Injector) (*PipelineService, error) {
	encoderSvc := do.MustInvoke[*services.EncoderService](injector)
	solverSvc := do.MustInvoke[*services1.SolverService](injector)
	cuberSvc := do.MustInvoke[*services2.CuberService](injector)
	cubeSelectorSvc := do.MustInvoke[*services3.CubeSelectorService](injector)
	simplifierSvc := do.MustInvoke[*services4.SimplifierService](injector)
	slurmSvc := do.MustInvoke[*services5.SlurmService](injector)
	svc := &PipelineService{encoderSvc: encoderSvc, solverSvc: solverSvc, cuberSvc: cuberSvc, cubeSelectorSvc: cubeSelectorSvc, simplifierSvc: simplifierSvc, slurmSvc: slurmSvc}
	return svc, nil
}
