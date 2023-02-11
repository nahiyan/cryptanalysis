package services

import (
	services2 "benchmark/internal/config/services"
	services3 "benchmark/internal/cuber/services"
	services1 "benchmark/internal/encoder/services"
	services "benchmark/internal/error/services"
	services7 "benchmark/internal/md4/services"
	services5 "benchmark/internal/simplifier/services"
	services6 "benchmark/internal/solution/services"
	services4 "benchmark/internal/solver/services"
	do "github.com/samber/do"
)

type SummarizerService struct {
	errorSvc      *services.ErrorService
	encoderSvc    *services1.EncoderService
	configSvc     *services2.ConfigService
	cuberSvc      *services3.CuberService
	solverSvc     *services4.SolverService
	simplifierSvc *services5.SimplifierService
	solutionSvc   *services6.SolutionService
	md4Svc        *services7.Md4Service
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
	svc := &SummarizerService{errorSvc: errorSvc, encoderSvc: encoderSvc, configSvc: configSvc, cuberSvc: cuberSvc, solverSvc: solverSvc, simplifierSvc: simplifierSvc, solutionSvc: solutionSvc, md4Svc: md4Svc}
	return svc, nil
}
