package services

import (
	services10 "cryptanalysis/internal/combined_logs/services"
	services2 "cryptanalysis/internal/config/services"
	services3 "cryptanalysis/internal/cuber/services"
	services1 "cryptanalysis/internal/encoder/services"
	services "cryptanalysis/internal/error/services"
	services7 "cryptanalysis/internal/md4/services"
	services8 "cryptanalysis/internal/md5/services"
	services9 "cryptanalysis/internal/sha256/services"
	services5 "cryptanalysis/internal/simplifier/services"
	services6 "cryptanalysis/internal/solution/services"
	services4 "cryptanalysis/internal/solver/services"
	do "github.com/samber/do"
)

type SummarizerService struct {
	errorSvc        *services.ErrorService
	encoderSvc      *services1.EncoderService
	configSvc       *services2.ConfigService
	cuberSvc        *services3.CuberService
	solverSvc       *services4.SolverService
	simplifierSvc   *services5.SimplifierService
	solutionSvc     *services6.SolutionService
	md4Svc          *services7.Md4Service
	md5Svc          *services8.Md5Service
	sha256Svc       *services9.Sha256Service
	combinedLogsSvc *services10.CombinedLogsService
}

func NewSummarizerService(injector *do.Injector) (*SummarizerService, error) {
	errorSvc := do.MustInvoke[*services.ErrorService](injector)
	encoderSvc := do.MustInvoke[*services1.EncoderService](injector)
	configSvc := do.MustInvoke[*services2.ConfigService](injector)
	cuberSvc := do.MustInvoke[*services3.CuberService](injector)
	solverSvc := do.MustInvoke[*services4.SolverService](injector)
	simplifierSvc := do.MustInvoke[*services5.SimplifierService](injector)
	solutionSvc := do.MustInvoke[*services6.SolutionService](injector)
	md4Svc := do.MustInvoke[*services7.Md4Service](injector)
	md5Svc := do.MustInvoke[*services8.Md5Service](injector)
	sha256Svc := do.MustInvoke[*services9.Sha256Service](injector)
	combinedLogsSvc := do.MustInvoke[*services10.CombinedLogsService](injector)
	svc := &SummarizerService{errorSvc: errorSvc, encoderSvc: encoderSvc, configSvc: configSvc, cuberSvc: cuberSvc, solverSvc: solverSvc, simplifierSvc: simplifierSvc, solutionSvc: solutionSvc, md4Svc: md4Svc, md5Svc: md5Svc, sha256Svc: sha256Svc, combinedLogsSvc: combinedLogsSvc}
	return svc, nil
}
