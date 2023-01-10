package services

import (
	services1 "benchmark/internal/config/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	services3 "benchmark/internal/simplification/services"
	do "github.com/samber/do"
)

type SimplifierService struct {
	errorSvc          *services.ErrorService
	configSvc         *services1.ConfigService
	filesystemSvc     *services2.FilesystemService
	simplificationSvc *services3.SimplificationService
}

func NewSimplifierService(injector *do.Injector) (*SimplifierService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	simplificationSvc := do.MustInvoke[*services3.SimplificationService](injector)
	svc := &SimplifierService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc, simplificationSvc: simplificationSvc}
	return svc, nil
}
