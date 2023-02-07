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

type Encoding struct {
	Encoder       encoder.Encoder
	Xor           []int
	Dobbertin     []int
	DobbertinBits []int
	Adders        []encoder.AdderType
	Hashes        []string
	Steps         []int
	Solvers       []Solver
}
type Solving struct {
	Solvers   []solver.Solver
	Timeout   int
	Workers   int
	Redundant bool
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
	Seed     int64
	Offset   int
	Indices  []int
}

type Simplifying struct {
	Name      string
	Conflicts []int
	Timeout   int
	Workers   int
}

type Pipe struct {
	Type Type

	// Type: encode
	Encoding

	// Type: simplifying
	Simplifying

	// Type: cube
	Cubing

	// Type: cube_select
	CubeSelecting

	// Type: solve
	Solving
}

type Value interface {
	string
}

type SlurmPipeOutput struct {
	Jobs   []slurm.Job
	Values interface{}
}

type EncodingPromise interface {
	Get(dependencies map[string]interface{}) string
	GetPath() string
}
