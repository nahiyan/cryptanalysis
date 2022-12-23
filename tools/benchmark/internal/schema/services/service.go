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
