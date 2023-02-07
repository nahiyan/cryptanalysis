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
			Transalg         string
		}

		Database  string
		Encodings string
		Logs      string
		Solutions string
		Cubesets  string
	}

	Slurm struct {
		MaxJobs   int
		ExtraTime float64
	}
}
