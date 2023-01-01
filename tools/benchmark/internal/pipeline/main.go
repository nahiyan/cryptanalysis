package pipeline

import "benchmark/internal/solver"

const (
	Encode = "encode"
	Solve  = "solve"
	Cube   = "cube"
)

type Type string
type Platform string
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
}
type Solving struct {
	Platform
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
	Platform
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
