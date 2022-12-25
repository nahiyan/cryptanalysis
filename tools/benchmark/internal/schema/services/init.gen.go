package services

import (
    "github.com/samber/do"
    pipelineServices "benchmark/internal/pipeline/services"
    
)

type SchemaService struct {
    Properties
    pipelineSvc *pipelineServices.PipelineService
}

func NewSchemaService(i *do.Injector) (*SchemaService, error) {
    pipelineSvc := do.MustInvoke[*pipelineServices.PipelineService](i)

    svc := &SchemaService{
        pipelineSvc: pipelineSvc,
    }

    

	return svc, nil
}
