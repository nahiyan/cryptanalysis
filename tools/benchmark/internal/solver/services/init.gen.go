package services

import (
	services "benchmark/internal/config/services"
	services5 "benchmark/internal/cube_selector/services"
	services6 "benchmark/internal/encoder/services"
	services2 "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	services7 "benchmark/internal/log/services"
	services4 "benchmark/internal/slurm/services"
	services3 "benchmark/internal/solution/services"
	do "github.com/samber/do"
)

type SolverService struct {
	configSvc       *services.ConfigService
	filesystemSvc   *services1.FilesystemService
	errorSvc        *services2.ErrorService
	solutionSvc     *services3.SolutionService
	slurmSvc        *services4.SlurmService
	cubeSelectorSvc *services5.CubeSelectorService
	encoderSvc      *services6.EncoderService
	logSvc          *services7.LogService
}

func NewSolverService(injector *do.Injector) (*SolverService, error) {
	configSvc := do.MustInvoke[*services.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	errorSvc := do.MustInvoke[*services2.ErrorService](injector)
	solutionSvc := do.MustInvoke[*services3.SolutionService](injector)
	slurmSvc := do.MustInvoke[*services4.SlurmService](injector)
	cubeSelectorSvc := do.MustInvoke[*services5.CubeSelectorService](injector)
	encoderSvc := do.MustInvoke[*services6.EncoderService](injector)
	logSvc := do.MustInvoke[*services7.LogService](injector)
	svc := &SolverService{configSvc: configSvc, filesystemSvc: filesystemSvc, errorSvc: errorSvc, solutionSvc: solutionSvc, slurmSvc: slurmSvc, cubeSelectorSvc: cubeSelectorSvc, encoderSvc: encoderSvc, logSvc: logSvc}
	return svc, nil
}
