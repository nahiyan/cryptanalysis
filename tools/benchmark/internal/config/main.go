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
			Satelite         string
			Verifier         string
			SolutionAnalyzer string
			SaeedE           string
			Benchmark        string
		}

		Database string
	}

	Slurm struct {
		MaxJobs int
	}
}
