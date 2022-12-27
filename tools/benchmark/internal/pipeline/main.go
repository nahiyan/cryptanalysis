package pipeline

const (
	Encode = "encode"
	Solve  = "solve"
)

type Type string
type Platform string
type Solver string
type AdderType string

type Variation struct {
	Xor           []int
	Dobbertin     []int
	DobbertinBits []int
	Adders        []AdderType
	Hashes        []string
	Steps         []int
	Solvers       []Solver
}

type Pipe struct {
	Type Type

	// Type: encode
	Variation

	// Type: solve
	Platform
	Solvers []Solver
	Timeout int
	Workers int
}
