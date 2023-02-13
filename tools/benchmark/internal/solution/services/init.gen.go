package services

import (
	services4 "benchmark/internal/command/services"
	services1 "benchmark/internal/config/services"
	services5 "benchmark/internal/encoder/services"
	services "benchmark/internal/error/services"
	services2 "benchmark/internal/filesystem/services"
	services3 "benchmark/internal/marshalling/services"
	do "github.com/samber/do"
)

type SolutionService struct {
	errorSvc       *services.ErrorService
	configSvc      *services1.ConfigService
	filesystemSvc  *services2.FilesystemService
	marshallingSvc *services3.MarshallingService
	commandSvc     *services4.CommandService
	encoderSvc     *services5.EncoderService
}

func NewSolutionService(injector *do.Injector) (*SolutionService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	marshallingSvc := do.MustInvoke[*services3.MarshallingService](injector)
	commandSvc := do.MustInvoke[*services4.CommandService](injector)
	encoderSvc := do.MustInvoke[*services5.EncoderService](injector)
	svc := &SolutionService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc, marshallingSvc: marshallingSvc, commandSvc: commandSvc, encoderSvc: encoderSvc}
	return svc, nil
}
