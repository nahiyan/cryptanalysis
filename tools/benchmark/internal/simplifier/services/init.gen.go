package services

import (
	services1 "benchmark/internal/config/services"
	services4 "benchmark/internal/cube_selector/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	services5 "benchmark/internal/log/services"
	services3 "benchmark/internal/simplification/services"
	do "github.com/samber/do"
)

type SimplifierService struct {
	errorSvc          *services.ErrorService
	configSvc         *services1.ConfigService
	filesystemSvc     *services2.FilesystemService
	simplificationSvc *services3.SimplificationService
	cubeSelectorSvc   *services4.CubeSelectorService
	logSvc            *services5.LogService
}

func NewSimplifierService(injector *do.Injector) (*SimplifierService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	simplificationSvc := do.MustInvoke[*services3.SimplificationService](injector)
	cubeSelectorSvc := do.MustInvoke[*services4.CubeSelectorService](injector)
	logSvc := do.MustInvoke[*services5.LogService](injector)
	svc := &SimplifierService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc, simplificationSvc: simplificationSvc, cubeSelectorSvc: cubeSelectorSvc, logSvc: logSvc}
	return svc, nil
}
