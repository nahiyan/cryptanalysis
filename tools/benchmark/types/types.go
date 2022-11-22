package types

type Variations struct {
	VariationsXor           []uint
	VariationsAdders        []string
	VariationsSatSolvers    []string
	VariationsDobbertin     []uint
	VariationsDobbertinBits []uint
	VariationsSteps         []uint
	VariationsHashes        []string
}

type CubeParams struct {
	CutoffVars    uint // The max. variable count for cutoff
	SelectionSize uint // Number of cubes in the selection
	CubeIndex     uint // Index of a specific cube to solve
}

type Simplification struct {
	Passes       uint
	Simplifier   string
	InstanceName string
}

type CommandContext struct {
	Variations
	InstanceMaxTime             uint
	MaxConcurrentInstancesCount uint
	CleanResults                bool
	Digest                      uint
	GenerateEncodings           uint
	SessionId                   uint
	CubeParams                  *CubeParams
	Seed                        int64 // Seed for random selection of cubes
	Simplification
}

type BenchmarkContext struct {
	Progress         map[string][]bool
	RunningInstances uint
}

type EncodingsGenContext struct {
	Variations
	CubeParams *CubeParams
}

type Config struct {
	Paths struct {
		Bin struct {
			CryptoMiniSat    string
			Kissat           string
			Cadical          string
			Glucose          string
			MapleSat         string
			XnfSat           string
			March            string
			Validator        string
			SolutionAnalyzer string
			Encoder          string
			Benchmark        string
		}
	}
	Slurm struct {
		MaxJobs uint
	}
}

type SlurmJob struct {
	Head struct {
		Nodes    uint
		CpuCores uint
		Memory   uint
		Time     uint
	}
	Body string
}
