package services

import (
	services "benchmark/internal/encoder/services"
	do "github.com/samber/do"
)

type PipelineService struct {
	encoderSvc *services.EncoderService
	Properties
}

func NewPipelineService(injector *do.Injector) (*PipelineService, error) {
	encoderSvc := do.MustInvoke[*services.EncoderService](injector)
	svc := &PipelineService{encoderSvc: encoderSvc}
	return svc, nil
}
