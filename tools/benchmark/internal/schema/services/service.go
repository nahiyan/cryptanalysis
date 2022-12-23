package services

import (
	"benchmark/internal/pipeline/services"
	"benchmark/internal/schema"
	"fmt"
	"log"

	"github.com/samber/do"
	"github.com/spf13/viper"
)

type SchemaService struct {
	Schema      schema.Schema
	PipelineSvc *services.PipelineService
}

func NewSchemaService(i *do.Injector) (*SchemaService, error) {
	pipelineSvc := do.MustInvoke[*services.PipelineService](i)
	return &SchemaService{PipelineSvc: pipelineSvc}, nil
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
	schemaSvc.PipelineSvc.Run()
}
