package services

import (
	"benchmark/internal/config"
	"log"

	"github.com/spf13/viper"
)

type Properties struct {
	Config config.Config
}

func (configSvc *ConfigService) Init() {
	configSvc.Process()
}

// Important: Register new SAT Solver here
func (configSvc *ConfigService) Process() {
	configFilePath := "./config.toml"

	// Binaries
	viper.SetDefault("Paths.Bin.CryptoMiniSat", "cryptominisat")
	viper.SetDefault("Paths.Bin.Kissat", "kissat")
	viper.SetDefault("Paths.Bin.Cadical", "cadical")
	viper.SetDefault("Paths.Bin.Glucose", "glucose")
	viper.SetDefault("Paths.Bin.MapleSat", "maplesat")
	viper.SetDefault("Paths.Bin.XnfSat", "xnfsat")
	viper.SetDefault("Paths.Bin.YalSat", "yalsat")
	viper.SetDefault("Paths.Bin.PalSat", "palsat")
	viper.SetDefault("Paths.Bin.LSTechMaple", "lstech_maple")
	viper.SetDefault("Paths.Bin.KissatCF", "kissat_cf")
	viper.SetDefault("Paths.Bin.March", "march_cu")
	viper.SetDefault("Paths.Bin.Satelite", "satelite")
	viper.SetDefault("Paths.Bin.SolutionAnalyzer", "solution_analyzer")
	viper.SetDefault("Paths.Bin.SaeedE", "saeed_e")
	viper.SetDefault("Paths.Bin.Transalg", "transalg")
	viper.SetDefault("Paths.Bin.Verifier", "saeed_e_verify")
	viper.SetDefault("Paths.Bin.Benchmark", "benchmark")

	// Database
	viper.SetDefault("Paths.Database", "database.db")

	// Output
	viper.SetDefault("Paths.Logs", "logs")
	viper.SetDefault("Paths.Solutions", "solutions")
	viper.SetDefault("Paths.Encodings", "encodings")
	viper.SetDefault("Paths.Cubesets", "cubesets")
	viper.SetDefault("Paths.Tmp", "tmp")

	// Slurm
	viper.SetDefault("Slurm.MaxJobs", 1000)
	viper.SetDefault("Slurm.WorkerTimeMul", 1)
	viper.SetDefault("Slurm.WorkerMemory", 300)

	// Solver
	viper.SetDefault("Solver.Slurm.NumTaskSelectWorkers", 1000)
	viper.SetDefault("Solver.Kissat.LocalSearch", false)
	viper.SetDefault("Solver.Kissat.LocalSearchEffort", 50)
	viper.SetDefault("Solver.Cadical.LocalSearchRounds", 0)
	viper.SetDefault("Solver.CryptoMiniSat.LocalSearch", false)
	viper.SetDefault("Solver.CryptoMiniSat.LocalSearchType", "ccnr")

	// Set config file
	viper.SetConfigFile(configFilePath)
	viper.ReadInConfig()

	// Unwrap the structure
	if err := viper.Unmarshal(&configSvc.Config); err != nil {
		log.Fatal("Failed to unmarshal viper config")
	}
}
