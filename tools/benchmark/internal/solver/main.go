package solver

import (
	"time"
)

type Solver string
type Platform int
type Result int

type Solution struct {
	Runtime      time.Duration
	Result       Result
	Solver       Solver
	ExitCode     int
	InstanceName string
	Verified     bool
	Checksum     string
}
