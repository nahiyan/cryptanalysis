package pipeline

import (
	cubeselector "benchmark/internal/cube_selector"
	"benchmark/internal/cuber"
	"benchmark/internal/encoder"
	"benchmark/internal/solver"
)

const (
	Encode          = "encode"
	Simplify        = "simplify"
	Cube            = "cube"
	IncrementalCube = "inc_cube"
	CubeSelect      = "cube_select"
	Solve           = "solve"
	SlurmSolve      = "slurm_solve"
)

type Type string
type Solver string

type EncodeParams struct {
	Encoder       encoder.Encoder
	Function      encoder.Function
	Xor           []int
	Dobbertin     []int
	DobbertinBits []int
	Adders        []encoder.AdderType
	Hashes        []string
	Steps         []int
	Solvers       []Solver
	Redundant     bool
}
type SolveParams struct {
	Solvers   []solver.Solver
	Timeout   int
	Workers   int
	Redundant bool
}

type CubeParams struct {
	MaxCubes         int
	MinCubes         int
	MinRefutedLeaves int
	MinThreshold     int
	StepChange       int
	InitialThreshold int
	// TODO: Add support for cutoff depth apart from in increamental cubing
	ThresholdType cuber.ThresholdType
	Thresholds    []int
	Workers       int
	Timeout       int
}

type CubeSelectParams struct {
	Type     cubeselector.CubeSelectionType
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
