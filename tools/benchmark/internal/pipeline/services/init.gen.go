package services

import (
    "github.com/samber/do"
    encoderServices "benchmark/internal/encoder/services"
    
)

type PipelineService struct {
    Properties
    encoderSvc *encoderServices.EncoderService
}

func NewPipelineService(i *do.Injector) (*PipelineService, error) {
    encoderSvc := do.MustInvoke[*encoderServices.EncoderService](i)

    svc := &PipelineService{
        encoderSvc: encoderSvc,
    }

    

	return svc, nil
}
