package pipeline

const (
	Encode = "encode"
	Solve  = "solve"
)

const (
	Slurm   = "slurm"
	Regular = "regular"
)

type Type string
type Platform string
type Solver string

type Variations struct {
	Xor           []int
	Dobbertin     []int
	DobbertinBits []int
	Adders        []string
	Hashes        []string
	Steps         []int
}

type Pipe struct {
	Type Type

	// Type: encode
	Variations

	// Type: solve
	Platform
	Solvers []Solver
	Timeout int
}
