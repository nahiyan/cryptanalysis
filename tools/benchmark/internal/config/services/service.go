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

func (configSvc *ConfigService) Process() {
	configFilePath := "./config.toml"

	// Binaries
	viper.SetDefault("Paths.Bin.CryptoMiniSat", "cryptominisat")
	viper.SetDefault("Paths.Bin.Kissat", "kissat")
	viper.SetDefault("Paths.Bin.Cadical", "cadical")
	viper.SetDefault("Paths.Bin.Glucose", "glucose")
	viper.SetDefault("Paths.Bin.MapleSat", "maplesat")
	viper.SetDefault("Paths.Bin.XnfSat", "xnfsat")
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

	// Set config file
	viper.SetConfigFile(configFilePath)
	viper.ReadInConfig()

	// Unwrap the structure
	if err := viper.Unmarshal(&configSvc.Config); err != nil {
		log.Fatal("Failed to unmarshal viper config")
	}
}
