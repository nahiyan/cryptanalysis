package config

import (
	"benchmark/types"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var config *types.Config

func ProcessConfig() {
	viper.SetDefault("Paths.Bin.CryptoMiniSat", "../../../sat-solvers/cryptominisat")
	viper.SetDefault("Paths.Bin.Kissat", "../../../sat-solvers/kissat")
	viper.SetDefault("Paths.Bin.Cadical", "../../../sat-solvers/cadical")
	viper.SetDefault("Paths.Bin.Glucose", "../../../sat-solvers/glucose")
	viper.SetDefault("Paths.Bin.MapleSat", "../../../sat-solvers/maplesat")
	viper.SetDefault("Paths.Bin.XnfSat", "../../../sat-solvers/xnfsat")
	viper.SetDefault("Paths.Bin.March", "../../../sat-solvers/march_cu")
	viper.SetDefault("Paths.Bin.SolutionAnalyzer", "../solution_analyzer/target/release/solution_analyzer")
	viper.SetDefault("Paths.Bin.Encoder", "../../encoders/saeed/crypto/main")
	viper.SetDefault("Paths.Bin.Validator", "../../encoders/saeed/crypto/verify-md4")
	viper.SetDefault("Slurm.MaxJobs", 1000)

	benchmarkAbsPath, err := filepath.Abs("./benchmark")
	if err != nil {
		benchmarkAbsPath = "./benchmark"
	}
	viper.SetDefault("Paths.Bin.Benchmark", benchmarkAbsPath)

	viper.SetConfigFile("./config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read the config file")
	}

	config = new(types.Config)
	if err := viper.Unmarshal(config); err != nil {
		log.Fatal("Failed to unmarshal viper config")
	}
}

func Get() *types.Config {
	return config
}
