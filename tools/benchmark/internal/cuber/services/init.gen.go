package services

import (
	services6 "benchmark/internal/command/services"
	services2 "benchmark/internal/config/services"
	services3 "benchmark/internal/cubeset/services"
	services7 "benchmark/internal/encoder/services"
	services4 "benchmark/internal/encoding/services"
	services "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	services8 "benchmark/internal/log/services"
	services5 "benchmark/internal/slurm/services"
	do "github.com/samber/do"
)

type CuberService struct {
	errorSvc      *services.ErrorService
	filesystemSvc *services1.FilesystemService
	configSvc     *services2.ConfigService
	cubesetSvc    *services3.CubesetService
	encodingSvc   *services4.EncodingService
	slurmSvc      *services5.SlurmService
	commandSvc    *services6.CommandService
	encoderSvc    *services7.EncoderService
	logSvc        *services8.LogService
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
	svc := &CuberService{errorSvc: errorSvc, filesystemSvc: filesystemSvc, configSvc: configSvc, cubesetSvc: cubesetSvc, encodingSvc: encodingSvc, slurmSvc: slurmSvc, commandSvc: commandSvc, encoderSvc: encoderSvc, logSvc: logSvc}
	return svc, nil
}
