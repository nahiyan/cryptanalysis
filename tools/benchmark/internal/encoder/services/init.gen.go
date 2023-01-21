package services

import (
	services3 "benchmark/internal/command/services"
	services "benchmark/internal/config/services"
	services2 "benchmark/internal/error/services"
	services1 "benchmark/internal/filesystem/services"
	do "github.com/samber/do"
)

type EncoderService struct {
	configSvc     *services.ConfigService
	filesystemSvc *services1.FilesystemService
	errorSvc      *services2.ErrorService
	commandSvc    *services3.CommandService
}

func NewEncoderService(injector *do.Injector) (*EncoderService, error) {
	configSvc := do.MustInvoke[*services.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services1.FilesystemService](injector)
	errorSvc := do.MustInvoke[*services2.ErrorService](injector)
	commandSvc := do.MustInvoke[*services3.CommandService](injector)
	svc := &EncoderService{configSvc: configSvc, filesystemSvc: filesystemSvc, errorSvc: errorSvc, commandSvc: commandSvc}
	return svc, nil
}
