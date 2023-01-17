package solveslurmtask

import (
	"benchmark/internal/solver"
	"time"
)

type Task struct {
	Encoding string
	Solver   solver.Solver
	Timeout  time.Duration
	Booked   bool
	PingTime time.Time
}
