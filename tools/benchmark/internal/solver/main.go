package solver

import (
	"time"
)

type Solver string
type Platform int
type Result int

type Settings struct {
	Solvers  []Solver
	Timeout  time.Duration
	Platform Platform
	Workers  int
}

type Solution struct {
	Runtime time.Duration
	Result  Result
	Solver  Solver
}
