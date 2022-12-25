package services

import (
    "github.com/samber/do"
    configServices "benchmark/internal/config/services"
    filesystemServices "benchmark/internal/filesystem/services"
    errorServices "benchmark/internal/error/services"
    
)

type EncoderService struct {
    
    configSvc *configServices.ConfigService
    filesystemSvc *filesystemServices.FilesystemService
    errorSvc *errorServices.ErrorService
}

func NewEncoderService(i *do.Injector) (*EncoderService, error) {
    configSvc := do.MustInvoke[*configServices.ConfigService](i)
    filesystemSvc := do.MustInvoke[*filesystemServices.FilesystemService](i)
    errorSvc := do.MustInvoke[*errorServices.ErrorService](i)

    svc := &EncoderService{
        configSvc: configSvc,
        filesystemSvc: filesystemSvc,
        errorSvc: errorSvc,
    }

    

	return svc, nil
}
