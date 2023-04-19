package services

import (
	services9 "benchmark/internal/combined_logs/services"
	services6 "benchmark/internal/command/services"
	services2 "benchmark/internal/config/services"
	services12 "benchmark/internal/cube_selector/services"
	services3 "benchmark/internal/cubeset/services"
	services7 "benchmark/internal/encoder/services"
	services4 "benchmark/internal/encoding/services"
	services "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	services8 "benchmark/internal/log/services"
	services13 "benchmark/internal/random/services"
	services10 "benchmark/internal/simplifier/services"
	services5 "benchmark/internal/slurm/services"
	services11 "benchmark/internal/solver/services"
	do "github.com/samber/do"
)

type CuberService struct {
	errorSvc        *services.ErrorService
	filesystemSvc   *services1.FilesystemService
	configSvc       *services2.ConfigService
	cubesetSvc      *services3.CubesetService
	encodingSvc     *services4.EncodingService
	slurmSvc        *services5.SlurmService
	commandSvc      *services6.CommandService
	encoderSvc      *services7.EncoderService
	logSvc          *services8.LogService
	combinedLogsSvc *services9.CombinedLogsService
	simplifierSvc   *services10.SimplifierService
	solverSvc       *services11.SolverService
	cubeSelectorSvc *services12.CubeSelectorService
	randomSvc       *services13.RandomService
}

func NewCuberService(injector *do.Injector) (*CuberService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	configSvc := do.MustInvoke[*services2.ConfigService](injector)
	cubesetSvc := do.MustInvoke[*services3.CubesetService](injector)
	encodingSvc := do.MustInvoke[*services4.EncodingService](injector)
	slurmSvc := do.MustInvoke[*services5.SlurmService](injector)
	commandSvc := do.MustInvoke[*services6.CommandService](injector)
	encoderSvc := do.MustInvoke[*services7.EncoderService](injector)
	logSvc := do.MustInvoke[*services8.LogService](injector)
	combinedLogsSvc := do.MustInvoke[*services9.CombinedLogsService](injector)
	simplifierSvc := do.MustInvoke[*services10.SimplifierService](injector)
	solverSvc := do.MustInvoke[*services11.SolverService](injector)
	cubeSelectorSvc := do.MustInvoke[*services12.CubeSelectorService](injector)
	randomSvc := do.MustInvoke[*services13.RandomService](injector)
	svc := &CuberService{errorSvc: errorSvc, filesystemSvc: filesystemSvc, configSvc: configSvc, cubesetSvc: cubesetSvc, encodingSvc: encodingSvc, slurmSvc: slurmSvc, commandSvc: commandSvc, encoderSvc: encoderSvc, logSvc: logSvc, combinedLogsSvc: combinedLogsSvc, simplifierSvc: simplifierSvc, solverSvc: solverSvc, cubeSelectorSvc: cubeSelectorSvc, randomSvc: randomSvc}
	return svc, nil
}
