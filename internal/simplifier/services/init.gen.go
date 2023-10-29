package services

import (
	services6 "cryptanalysis/internal/combined_logs/services"
	services1 "cryptanalysis/internal/config/services"
	services4 "cryptanalysis/internal/cube_selector/services"
	services7 "cryptanalysis/internal/encoder/services"
	services "cryptanalysis/internal/error/services"
	services2 "cryptanalysis/internal/filesystem/services"
	services5 "cryptanalysis/internal/log/services"
	services3 "cryptanalysis/internal/simplification/services"
	do "github.com/samber/do"
)

type SimplifierService struct {
	errorSvc          *services.ErrorService
	configSvc         *services1.ConfigService
	filesystemSvc     *services2.FilesystemService
	simplificationSvc *services3.SimplificationService
	cubeSelectorSvc   *services4.CubeSelectorService
	logSvc            *services5.LogService
	combinedLogsSvc   *services6.CombinedLogsService
	encoderSvc        *services7.EncoderService
}

func NewSimplifierService(injector *do.Injector) (*SimplifierService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	configSvc := do.MustInvoke[*services1.ConfigService](injector)
	filesystemSvc := do.MustInvoke[*services2.FilesystemService](injector)
	simplificationSvc := do.MustInvoke[*services3.SimplificationService](injector)
	cubeSelectorSvc := do.MustInvoke[*services4.CubeSelectorService](injector)
	logSvc := do.MustInvoke[*services5.LogService](injector)
	combinedLogsSvc := do.MustInvoke[*services6.CombinedLogsService](injector)
	encoderSvc := do.MustInvoke[*services7.EncoderService](injector)
	svc := &SimplifierService{errorSvc: errorSvc, configSvc: configSvc, filesystemSvc: filesystemSvc, simplificationSvc: simplificationSvc, cubeSelectorSvc: cubeSelectorSvc, logSvc: logSvc, combinedLogsSvc: combinedLogsSvc, encoderSvc: encoderSvc}
	return svc, nil
}
