package pipeline

import (
	"benchmark/internal/encoder"
	"benchmark/internal/slurm"
	"benchmark/internal/solver"
)

const (
	Encode           = "encode"
	Simplify         = "simplify"
	Cube             = "cube"
	SlurmCube        = "slurm_cube"
	CubeSelect       = "cube_select"
	SlurmCubeSelect  = "slurm_cube_select"
	Solve            = "solve"
	SlurmSolve       = "slurm_solve"
	EncodingSlurmify = "encoding_slurmify"
)

type Type string
type Solver string

type EncodeParams struct {
	Encoder       encoder.Encoder
	Xor           []int
	Dobbertin     []int
	DobbertinBits []int
	Adders        []encoder.AdderType
	Hashes        []string
	Steps         []int
	Solvers       []Solver
}
type SolveParams struct {
	Solvers   []solver.Solver
	Timeout   int
	Workers   int
	Redundant bool
}

type CubeParams struct {
	MaxCubes         int
	MinRefutedLeaves int
	MinThreshold     int
	Thresholds       []int
	Workers          int
	Timeout          int
}

type CubeSelectParams struct {
	Type     string
	Quantity int
	Seed     int64
	Offset   int
	Indices  []int
}

type SimplifyParams struct {
	Name      string
	Conflicts []int
	Timeout   int
	Workers   int
}

type Pipe struct {
	Type Type

	// Type: encode
	EncodeParams

	// Type: simplifying
	SimplifyParams

	// Type: cube
	CubeParams

	// Type: cube_select
	CubeSelectParams

	// Type: solve
	SolveParams
}

type Value interface {
	string
}

type SlurmPipeOutput struct {
	Jobs   []slurm.Job
	Values interface{}
}
