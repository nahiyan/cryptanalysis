package types

type CommandContext struct {
	VariationsXor               []uint
	VariationsAdders            []string
	VariationsSatSolvers        []string
	VariationsDobbertin         []uint
	VariationsSteps             []uint
	VariationsHashes            []string
	InstanceMaxTime             uint
	MaxConcurrentInstancesCount uint
}

type BenchmarkContext struct {
	Progress         map[string][]bool
	RunningInstances uint
}
