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

func (logSvc *LogService) Debug(message string) {
	logSvc.logger.Debug(message)
	logSvc.logger.Sync()
}

func (logSvc *LogService) SolveResult(encoding string, solver_ solver.Solver, exitCode int, result solver.Result, runtime time.Duration) {
	logSvc.logger.Info(
		"Solve",
		zap.String("result", string(result)),
		zap.String("runtime", runtime.String()),
		zap.String("solver", string(solver_)),
		zap.String("encoding", encoding),
		zap.Int("exit code", exitCode))
	logSvc.logger.Sync()
}

func (logSvc *LogService) SimplifyResult(encoding string, runtime time.Duration, numVariables, numClauses, eliminatedVars int) {
	logSvc.logger.Info(
		"Simplify",
		zap.String("runtime", runtime.String()),
		zap.Int("variables", numVariables),
		zap.Int("clauses", numClauses),
		zap.Int("eliminated", eliminatedVars),
		zap.String("encoding", encoding))
	logSvc.logger.Sync()
}

func (logSvc *LogService) CubeResult(cubesFilePath string, processTime time.Duration, numCubes, numRefutedLeaves int, isRemoved bool) {
	logSvc.logger.Info(
		"Cube",
		zap.String("cubes file", cubesFilePath),
		zap.String("process time", processTime.String()),
		zap.Int("cubes", numCubes),
		zap.Int("refuted leaves", numRefutedLeaves),
		zap.Bool("removed", isRemoved))
	logSvc.logger.Sync()
}
