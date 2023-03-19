package config

// Important: Register new SAT Solver here
type Config struct {
	Paths struct {
		Bin struct {
			CryptoMiniSat    string
			Kissat           string
			Cadical          string
			Glucose          string
			MapleSat         string
			XnfSat           string
			YalSat           string
			PalSat           string
			LSTechMaple      string
			KissatCF         string
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
		Tmp       string
	}

	Slurm struct {
		MaxJobs       int
		WorkerTimeMul float64
		WorkerMemory  int
	}

	Solver struct {
		Slurm struct {
			NumTaskSelectWorkers int
		}

		Kissat struct {
			LocalSearch       bool
			LocalSearchEffort int
		}

		Cadical struct {
			LocalSearchRounds int
		}

		CryptoMiniSat struct {
			LocalSearch     bool
			LocalSearchType string
		}
	}
}
