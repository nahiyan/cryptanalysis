package services

import (
	"benchmark/internal/encoder"
	encoderServices "benchmark/internal/encoder/services"
	"benchmark/internal/pipeline"
	"fmt"

	"github.com/samber/do"
)

type PipelineService struct {
	EncoderSvc encoder.EncoderService
	Pipeline   []pipeline.Pipe
}

const (
	ListOfEncodings = "list_of_encodings"
	ListOfSolutions = "list_of_solutions"
	None            = "none"
)

type InputOutputType string

func NewPipelineService(i *do.Injector) (*PipelineService, error) {
	return &PipelineService{
		EncoderSvc: do.MustInvoke[*encoderServices.EncoderService](i),
	}, nil
}

func getInputType(pipe pipeline.Pipe) InputOutputType {
	switch pipe.Type {
	case pipeline.Encode:
		return None
	case pipeline.Solve:
		return ListOfEncodings
	}

	return None
}

func getOutputType(pipe pipeline.Pipe) InputOutputType {
	switch pipe.Type {
	case pipeline.Encode:
		return ListOfEncodings
	case pipeline.Solve:
		return ListOfSolutions
	}

	return None
}

// Check if the pipelines can be connected
func (pipelineSvc *PipelineService) Validate() {
	for i, pipe := range pipelineSvc.Pipeline {
		var nextPipeline *pipeline.Pipe
		if len(pipelineSvc.Pipeline) < i+2 {
			nextPipeline = &pipelineSvc.Pipeline[i+1]
		}

		if nextPipeline == nil {
			break
		}

		outputType := getOutputType(pipe)
		nextPipelineInputType := getInputType(*nextPipeline)
		if nextPipelineInputType != outputType {
			panic("Incompatible pipeline: " + outputType + " can't fit into the expected input type " + nextPipelineInputType)
		}
	}
}

func (pipelineSvc *PipelineService) TestRun() {
	for _, pipe := range pipelineSvc.Pipeline {
		switch pipe.Type {
		case pipeline.Encode:
			x := pipelineSvc.EncoderSvc.TestRun()
			fmt.Println(x)
		}
	}
}

func (pipelineSvc *PipelineService) Run() {
	pipelineSvc.Validate()
	pipelineSvc.TestRun()
}
