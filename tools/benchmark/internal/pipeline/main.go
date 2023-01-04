package pipeline

import "benchmark/internal/solver"

const (
	Encode     = "encode"
	Solve      = "solve"
	SlurmSolve = "slurm_solve"
	Cube       = "cube"
	SlurmCube  = "slurm_cube"
)

type Type string
type Solver string
type AdderType string

type Encoding struct {
	Xor           []int
	Dobbertin     []int
	DobbertinBits []int
	Adders        []AdderType
	Hashes        []string
	Steps         []int
	Solvers       []Solver
	OutputDir     string
}
type Solving struct {
	Solvers []solver.Solver
	Timeout int
	Workers int
}

type Cubing struct {
	MaxCubes         int
	MinRefutedLeaves int
	MinThreshold     int
	Thresholds       []int
	Workers          int
	Timeout          int
}

type Pipe struct {
	Type Type

	// Type: encode
	Encoding

	// Type: solve
	Solving

	// Type: cube
	Cubing
}
