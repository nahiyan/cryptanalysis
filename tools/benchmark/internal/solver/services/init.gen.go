package services

import (
	services "benchmark/internal/config/services"
	services6 "benchmark/internal/cube_selector/services"
	services7 "benchmark/internal/encoder/services"
	services2 "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	services4 "benchmark/internal/slurm/services"
	services3 "benchmark/internal/solution/services"
	services5 "benchmark/internal/solve_slurm_task/services"
	do "github.com/samber/do"
)

type SolverService struct {
	configSvc         *services.ConfigService
	filesystemSvc     *services1.FilesystemService
	errorSvc          *services2.ErrorService
	solutionSvc       *services3.SolutionService
	slurmSvc          *services4.SlurmService
	solveSlurmTaskSvc *services5.SolveSlurmTaskService
	cubeSelectorSvc   *services6.CubeSelectorService
	encoderSvc        *services7.EncoderService
}

func NewSolverService(injector *do.Injector) (*SolverService, error) {
	configSvc := do.MustInvoke[*services.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	errorSvc := do.MustInvoke[*services2.ErrorService](injector)
	solutionSvc := do.MustInvoke[*services3.SolutionService](injector)
	slurmSvc := do.MustInvoke[*services4.SlurmService](injector)
	solveSlurmTaskSvc := do.MustInvoke[*services5.SolveSlurmTaskService](injector)
	cubeSelectorSvc := do.MustInvoke[*services6.CubeSelectorService](injector)
	encoderSvc := do.MustInvoke[*services7.EncoderService](injector)
	svc := &SolverService{configSvc: configSvc, filesystemSvc: filesystemSvc, errorSvc: errorSvc, solutionSvc: solutionSvc, slurmSvc: slurmSvc, solveSlurmTaskSvc: solveSlurmTaskSvc, cubeSelectorSvc: cubeSelectorSvc, encoderSvc: encoderSvc}
	return svc, nil
}
