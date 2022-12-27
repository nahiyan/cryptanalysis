package services

import (
	"benchmark/internal/schema"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Properties struct {
	Schema schema.Schema
}

func (schemaSvc *SchemaService) Process(filePath string) {
	// Set config file
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read the config file")
	}

	// Unwrap the structure
	if err := viper.Unmarshal(&schemaSvc.Schema); err != nil {
		log.Fatal("Failed to unmarshal viper config")
	}

	fmt.Println(schemaSvc.Schema)

	// Run the pipeline
	schemaSvc.pipelineSvc.Run(schemaSvc.Schema.Pipeline)
}
