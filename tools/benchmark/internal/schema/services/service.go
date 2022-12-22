package services

import (
	"benchmark/internal/schema"
	"fmt"
	"log"

	"github.com/samber/do"
	"github.com/spf13/viper"
)

type SchemaService struct {
	Schema schema.Schema
}

func NewSchemaService(i *do.Injector) (*SchemaService, error) {
	return new(SchemaService), nil
}

func (c *SchemaService) Process(filePath string) {
	// benchmarkBinAbsPath, err := filepath.Abs("./benchmark")
	// if err != nil {
	// 	benchmarkBinAbsPath = "./benchmark"
	// }

	// // Binaries
	// viper.SetDefault("Paths.Bin.CryptoMiniSat", "../../../sat-solvers/cryptominisat")
	// viper.SetDefault("Paths.Bin.Kissat", "../../../sat-solvers/kissat")
	// viper.SetDefault("Paths.Bin.Cadical", "../../../sat-solvers/cadical")
	// viper.SetDefault("Paths.Bin.Glucose", "../../../sat-solvers/glucose")
	// viper.SetDefault("Paths.Bin.MapleSat", "../../../sat-solvers/maplesat")
	// viper.SetDefault("Paths.Bin.XnfSat", "../../../sat-solvers/xnfsat")
	// viper.SetDefault("Paths.Bin.March", "../../../sat-solvers/march_cu")
	// viper.SetDefault("Paths.Bin.SolutionAnalyzer", "../solution_analyzer/target/release/solution_analyzer")
	// viper.SetDefault("Paths.Bin.Encoder", "../../encoders/saeed/crypto/main")
	// viper.SetDefault("Paths.Bin.Verifier", "../../encoders/saeed/crypto/verify-md4")
	// viper.SetDefault("Paths.Bin.Benchmark", benchmarkBinAbsPath)

	// Slurm
	// viper.SetDefault("Slurm.MaxJobs", 1000)

	// Set config file
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read the config file")
	}

	// Unwrap the structure
	if err := viper.Unmarshal(&c.Schema); err != nil {
		log.Fatal("Failed to unmarshal viper config")
	}

	fmt.Println(c.Schema)
}
