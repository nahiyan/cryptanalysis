package config

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
			SaeedE           string
			Benchmark        string
		}
	}
	Slurm struct {
		MaxJobs uint
	}
}
