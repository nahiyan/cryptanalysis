package types

type CommandContext struct {
	VariationsXor               []uint
	VariationsAdders            []string
	VariationsSatSolvers        []string
	VariationsDobbertin         []uint
	VariationsDobbertinBits     []uint
	VariationsSteps             []uint
	VariationsHashes            []string
	InstanceMaxTime             uint
	MaxConcurrentInstancesCount uint
	CleanResults                bool
}

type BenchmarkContext struct {
	Progress         map[string][]bool
	RunningInstances uint
}
