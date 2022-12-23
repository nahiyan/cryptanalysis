package services

import (
	"benchmark/internal/pipelines"

	"github.com/samber/do"
)

type PipelinesService struct {
	Pipelines []pipelines.Pipeline
}

const (
	ListOfEncodings = iota
	ListOfSolutions
	None
)

type InputOutputType int

func NewPipelinesService(i *do.Injector) (*PipelinesService, error) {
	return &PipelinesService{}, nil
}

func getInputType(pipeline pipelines.Pipeline) InputOutputType {
	switch pipeline.Type {
	case pipelines.Encode:
		return None
	case pipelines.Solve:
		return ListOfEncodings
	}
}

func getOutputType(pipeline pipelines.Pipeline) InputOutputType {
	switch pipeline.Type {
	case pipelines.Encode:
		return ListOfEncodings
	case pipelines.Solve:
		return ListOfSolutions
	}
}

func (pipelinesSvc *PipelinesService) Run() {
	for _, pipeline := range pipelinesSvc.Pipelines {
		switch pipeline.Type {
		case pipelines.Encode:

		}
	}
}
