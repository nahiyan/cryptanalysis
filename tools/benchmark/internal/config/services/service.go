package services

import (
	"benchmark/internal/config"
	"log"
	"path/filepath"

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
	benchmarkBinAbsPath, err := filepath.Abs("./benchmark")
	if err != nil {
		benchmarkBinAbsPath = "./benchmark"
	}

	// Binaries
	viper.SetDefault("Paths.Bin.CryptoMiniSat", "../../../sat-solvers/cryptominisat")
	viper.SetDefault("Paths.Bin.Kissat", "../../../sat-solvers/kissat")
	viper.SetDefault("Paths.Bin.Cadical", "../../../sat-solvers/cadical")
	viper.SetDefault("Paths.Bin.Glucose", "../../../sat-solvers/glucose")
	viper.SetDefault("Paths.Bin.MapleSat", "../../../sat-solvers/maplesat")
	viper.SetDefault("Paths.Bin.XnfSat", "../../../sat-solvers/xnfsat")
	viper.SetDefault("Paths.Bin.March", "../../../sat-solvers/march_cu")
	viper.SetDefault("Paths.Bin.Satelite", "satelite")
	viper.SetDefault("Paths.Bin.SolutionAnalyzer", "../solution_analyzer/target/release/solution_analyzer")
	viper.SetDefault("Paths.Bin.SaeedE", "../../encoders/saeed/crypto/main")
	viper.SetDefault("Paths.Bin.Verifier", "../../encoders/saeed/crypto/verify-md4")
	viper.SetDefault("Paths.Bin.Benchmark", benchmarkBinAbsPath)

	// Database
	viper.SetDefault("Paths.Database", "./database")

	// Slurm
	viper.SetDefault("Slurm.MaxJobs", 1000)

	// Set config file
	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read the config file")
	}

	// Unwrap the structure
	if err := viper.Unmarshal(&configSvc.Config); err != nil {
		log.Fatal("Failed to unmarshal viper config")
	}
}
