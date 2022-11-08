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

type CommandContext struct {
	Variations
	IsCubeEnabled               bool
	CubeVars                    uint
	InstanceMaxTime             uint
	MaxConcurrentInstancesCount uint
	CleanResults                bool
	Digest                      uint
	GenerateEncodings           uint
}

type BenchmarkContext struct {
	Progress         map[string][]bool
	RunningInstances uint
}

type EncodingsGenContext struct {
	Variations
	IsCubeEnabled bool
	CubeVars      uint
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
