package services

import (
	"benchmark/internal/pipeline"
	"fmt"
)

const (
	ListOfEncodings = "list[encoding]"
	ListOfSolutions = "list[solution]"
	None            = "none"
)

type Properties struct {
	Pipeline []pipeline.Pipe
}

type InputOutputType string
type LoopHandler func(*pipeline.Pipe, *pipeline.Pipe)

func getInputType(pipe *pipeline.Pipe) InputOutputType {
	switch pipe.Type {
	case pipeline.Encode:
		return None
	case pipeline.Solve:
		return ListOfEncodings
	}

	return None
}

func getOutputType(pipe *pipeline.Pipe) InputOutputType {
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
	pipelineSvc.Loop(pipelineSvc.Pipeline, func(pipe, nextPipe *pipeline.Pipe) {
		if nextPipe == nil {
			return
		}

		outputType := getOutputType(pipe)
		nextPipelineInputType := getInputType(nextPipe)
		if nextPipelineInputType != outputType {
			panic("Incompatible pipeline: " + outputType + " can't fit into the expected input type " + nextPipelineInputType)
		}
	})
}

func (pipelineSvc *PipelineService) Loop(pipes []pipeline.Pipe, handler LoopHandler) {
	for i, pipe := range pipes {
		var nextPipe *pipeline.Pipe
		if len(pipelineSvc.Pipeline) > i+1 {
			nextPipe = &pipelineSvc.Pipeline[i+1]
		}

		handler(&pipe, nextPipe)
	}
}

func (pipelineSvc *PipelineService) TestRun() {
	var lastValue interface{}

	pipelineSvc.Loop(pipelineSvc.Pipeline, func(pipe, nextPipe *pipeline.Pipe) {
		switch pipe.Type {
		case pipeline.Encode:
			lastValue = pipelineSvc.encoderSvc.TestRun()
			fmt.Println("Encode", lastValue)

			if nextPipe == nil {
				return
			}

		case pipeline.Solve:
			pipelineSvc.solverSvc.Run(lastValue.([]string))
		}
	})
}

func (pipelineSvc *PipelineService) Run(pipes []pipeline.Pipe) {
	pipelineSvc.Pipeline = pipes

	pipelineSvc.Validate()
	pipelineSvc.TestRun()
}
