package services

import (
	services9 "cryptanalysis/internal/combined_logs/services"
	services6 "cryptanalysis/internal/command/services"
	services2 "cryptanalysis/internal/config/services"
	services12 "cryptanalysis/internal/cube_selector/services"
	services3 "cryptanalysis/internal/cubeset/services"
	services7 "cryptanalysis/internal/encoder/services"
	services4 "cryptanalysis/internal/encoding/services"
	services "cryptanalysis/internal/error/services"
	services1 "cryptanalysis/internal/filesystem/services"
	services8 "cryptanalysis/internal/log/services"
	services13 "cryptanalysis/internal/random/services"
	services10 "cryptanalysis/internal/simplifier/services"
	services5 "cryptanalysis/internal/slurm/services"
	services11 "cryptanalysis/internal/solver/services"
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
