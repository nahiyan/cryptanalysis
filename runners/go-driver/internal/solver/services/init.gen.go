package services

import (
	services9 "cryptanalysis/internal/combined_logs/services"
	services "cryptanalysis/internal/config/services"
	services5 "cryptanalysis/internal/cube_selector/services"
	services6 "cryptanalysis/internal/encoder/services"
	services2 "cryptanalysis/internal/error/services"
	services1 "cryptanalysis/internal/filesystem/services"
	services7 "cryptanalysis/internal/log/services"
	services8 "cryptanalysis/internal/random/services"
	services4 "cryptanalysis/internal/slurm/services"
	services3 "cryptanalysis/internal/solution/services"
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
	randomSvc       *services8.RandomService
	combinedLogsSvc *services9.CombinedLogsService
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
	randomSvc := do.MustInvoke[*services8.RandomService](injector)
	combinedLogsSvc := do.MustInvoke[*services9.CombinedLogsService](injector)
	svc := &SolverService{configSvc: configSvc, filesystemSvc: filesystemSvc, errorSvc: errorSvc, solutionSvc: solutionSvc, slurmSvc: slurmSvc, cubeSelectorSvc: cubeSelectorSvc, encoderSvc: encoderSvc, logSvc: logSvc, randomSvc: randomSvc, combinedLogsSvc: combinedLogsSvc}
	return svc, nil
}
