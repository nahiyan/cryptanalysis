package services

import (
	services8 "benchmark/internal/command/services"
	services3 "benchmark/internal/config/services"
	services6 "benchmark/internal/cube_slurm_task/services"
	services4 "benchmark/internal/cubeset/services"
	services1 "benchmark/internal/database/services"
	services5 "benchmark/internal/encoding/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	services7 "benchmark/internal/slurm/services"
	do "github.com/samber/do"
)

type CuberService struct {
	errorSvc         *services.ErrorService
	databaseSvc      *services1.DatabaseService
	filesystemSvc    *services2.FilesystemService
	configSvc        *services3.ConfigService
	cubesetSvc       *services4.CubesetService
	encodingSvc      *services5.EncodingService
	cubeSlurmTaskSvc *services6.CubeSlurmTaskService
	slurmSvc         *services7.SlurmService
	commandSvc       *services8.CommandService
}

func NewCuberService(injector *do.Injector) (*CuberService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	databaseSvc := do.MustInvoke[*services1.DatabaseService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	configSvc := do.MustInvoke[*services3.ConfigService](injector)
	cubesetSvc := do.MustInvoke[*services4.CubesetService](injector)
	encodingSvc := do.MustInvoke[*services5.EncodingService](injector)
	cubeSlurmTaskSvc := do.MustInvoke[*services6.CubeSlurmTaskService](injector)
	slurmSvc := do.MustInvoke[*services7.SlurmService](injector)
	commandSvc := do.MustInvoke[*services8.CommandService](injector)
	svc := &CuberService{errorSvc: errorSvc, databaseSvc: databaseSvc, filesystemSvc: filesystemSvc, configSvc: configSvc, cubesetSvc: cubesetSvc, encodingSvc: encodingSvc, cubeSlurmTaskSvc: cubeSlurmTaskSvc, slurmSvc: slurmSvc, commandSvc: commandSvc}
	return svc, nil
}
