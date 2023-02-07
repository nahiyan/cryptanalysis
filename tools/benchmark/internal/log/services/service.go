package services

import (
	"benchmark/internal/solver"
	"time"

	"go.uber.org/zap"
)

type Properties struct {
	logger *zap.Logger
}

func (logSvc *LogService) Init() {
	logSvc.logger, _ = zap.NewDevelopment()
}

func (logSvc *LogService) Info(message string) {
	logSvc.logger.Info(message)
	logSvc.logger.Sync()
}

func (logSvc *LogService) SolverResult(encoding string, solver_ solver.Solver, exitCode int, result solver.Result, runtime time.Duration) {
	logSvc.logger.Info(
		"Solver",
		zap.String("result", string(result)),
		zap.String("runtime", runtime.String()),
		zap.String("solver", string(solver_)),
		zap.String("encoding", encoding),
		zap.Int("exit code", exitCode))
	logSvc.logger.Sync()
}
