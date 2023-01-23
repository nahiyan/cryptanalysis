package solveslurmtask

import (
	"benchmark/internal/pipeline"
	"benchmark/internal/solver"
	"time"
)

type Task struct {
	EncodingPromise pipeline.EncodingPromise
	Solver          solver.Solver
	Timeout         time.Duration
}
