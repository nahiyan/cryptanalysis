package services

import (
	services "cryptanalysis/internal/pipeline/services"
	do "github.com/samber/do"
)

type SchemaService struct {
	pipelineSvc *services.PipelineService
	Properties
}

func NewSchemaService(injector *do.Injector) (*SchemaService, error) {
	pipelineSvc := do.MustInvoke[*services.PipelineService](injector)
	svc := &SchemaService{pipelineSvc: pipelineSvc}
	return svc, nil
}
