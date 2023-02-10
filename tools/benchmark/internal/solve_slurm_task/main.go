package solveslurmtask

import (
	"benchmark/internal/encoder"
	"benchmark/internal/solver"
	"time"
)

type Task struct {
	Encoding encoder.Encoding
	Solver   solver.Solver
	Timeout  time.Duration
}
