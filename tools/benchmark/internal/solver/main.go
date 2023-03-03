package solver

import (
	"time"
)

type Solver string
type Platform int
type Result string

type Solution struct {
	Runtime      time.Duration
	Result       Result
	Solver       Solver
	ExitCode     int
	InstanceName string
	Verified     bool
	Checksum     string
}

const (
	CryptoMiniSat = "cryptominisat"
	Cadical       = "cadical"
	Kissat        = "kissat"
	MapleSat      = "maplesat"
	Glucose       = "glucose"
	XnfSat        = "xnfsat"
	YalSat        = "yalsat"
	PalSat        = "palsat"
)

const (
	Sat   = "SAT"
	Unsat = "UNSAT"
	Fail  = "FAIL"
)
