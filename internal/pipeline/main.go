package pipeline

import (
	cubeselector "cryptanalysis/internal/cube_selector"
	"cryptanalysis/internal/cuber"
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/solver"
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
	Encoder    encoder.Encoder
	Function   encoder.Function
	Steps      []int
	AttackType encoder.AttackType

	// Nejati Encoder
	Xor    []int
	Adders []encoder.AdderType

	// Preimage attack
	Dobbertin     []int // TODO: Bring Dobbertin into a single entity
	DobbertinBits []int
	Hashes        []string

	Redundant bool
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
	Name             string
	Conflicts        []int
	ConflictsMap     map[int][]int
	Timeout          int
	Workers          int
	Reconstruct      bool
	PreserveComments bool
	// TODO: Add support for redundant simplification
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
