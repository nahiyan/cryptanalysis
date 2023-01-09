package pipeline

import (
	"benchmark/internal/slurm"
	"benchmark/internal/solver"
)

const (
	Encode          = "encode"
	Solve           = "solve"
	SlurmSolve      = "slurm_solve"
	Cube            = "cube"
	SlurmCube       = "slurm_cube"
	CubeSelect      = "cube_select"
	SlurmCubeSelect = "slurm_cube_select"
	Simplify        = "simplify"
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

type CubeSelecting struct {
	Type     string
	Quantity int
	Seed     int
}

type Simplifying struct {
	Name      string
	Conflicts int
	Timeout   int
}

type Pipe struct {
	Type Type

	// Type: encode
	Encoding

	// Type: solve
	Solving

	// Type: cube
	Cubing

	// Type: cube_select
	CubeSelecting

	// Type: simplifying
	Simplifying
}

type Value interface {
	string
}

type SlurmPipeOutput struct {
	Jobs   []slurm.Job
	Values interface{}
}
