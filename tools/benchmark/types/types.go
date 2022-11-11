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
	Seed          int  // Seed for random selection of cubes
	SelectionSize uint // Number of cubes in the selection
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
			Verifier         string
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
