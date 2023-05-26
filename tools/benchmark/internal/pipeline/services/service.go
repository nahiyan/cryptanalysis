package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"log"
)

const (
	ListOfEncodings         = "list[encoding]"
	ListOfSolutions         = "list[solution]"
	ListOfSlurmJobSolutions = "list[slurm_job[solution]]"
	ListOfCubesets          = "list[cubeset]"
	None                    = "none"
)

type InputOutputType string
type LoopHandler func(*pipeline.Pipe, *pipeline.Pipe)

func getInputType(pipe *pipeline.Pipe) InputOutputType {
	switch pipe.Type {
	case pipeline.Encode:
		return None
	case pipeline.Solve:
		return ListOfEncodings
	case pipeline.SlurmSolve:
		return ListOfEncodings
	case pipeline.Cube:
		return ListOfEncodings
	case pipeline.IncrementalCube:
		return ListOfEncodings
	case pipeline.CubeSelect:
		return ListOfCubesets
	case pipeline.Simplify:
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
	case pipeline.SlurmSolve:
		return ListOfSlurmJobSolutions
	case pipeline.Cube:
		return ListOfCubesets
	case pipeline.IncrementalCube:
		return ListOfCubesets
	case pipeline.CubeSelect:
		return ListOfEncodings
	case pipeline.Simplify:
		return ListOfEncodings
	}

	return None
}

// Check if the pipelines can be connected
func (pipelineSvc *PipelineService) Validate(pipes []pipeline.Pipe) {
	pipelineSvc.Loop(pipes, func(pipe, nextPipe *pipeline.Pipe) {
		if nextPipe == nil {
			return
		}

		outputType := getOutputType(pipe)
		nextPipelineInputType := getInputType(nextPipe)
		if nextPipelineInputType != outputType {
			log.Fatal("Incompatible pipeline: " + outputType + " can't fit into the expected input type " + nextPipelineInputType)
		}
	})
}

func (pipelineSvc *PipelineService) Loop(pipes []pipeline.Pipe, handler LoopHandler) {
	for i, pipe := range pipes {
		var nextPipe *pipeline.Pipe
		if len(pipes) > i+1 {
			nextPipe = &pipes[i+1]
		}

		handler(&pipe, nextPipe)
	}
}

func (pipelineSvc *PipelineService) RunPipes(pipes []pipeline.Pipe) {
	var lastValue interface{}
	pipelineSvc.Loop(pipes, func(pipe, nextPipe *pipeline.Pipe) {
		switch pipe.Type {
		case pipeline.Encode:
			// NejatiEncoder is the default encoder
			if pipe.EncodeParams.Encoder == "" {
				pipe.EncodeParams.Encoder = encoder.NejatiEncoder
			}
			lastValue = pipelineSvc.encoderSvc.Run(pipe.EncodeParams)

			if nextPipe == nil {
				return
			}

		case pipeline.Simplify:
			input, ok := lastValue.([]encoder.Encoding)
			if !ok {
				log.Fatal("Simplifier expects a list of encodings")
			}
			lastValue = pipelineSvc.simplifierSvc.Run(input, pipe.SimplifyParams)

		case pipeline.Cube:
			input, ok := lastValue.([]encoder.Encoding)
			if !ok {
				log.Fatal("Cuber expects a list of encodings")
			}

			lastValue = pipelineSvc.cuberSvc.Run(input, pipe.CubeParams)

		case pipeline.IncrementalCube:
			input, ok := lastValue.([]encoder.Encoding)
			if !ok {
				log.Fatal("Incremental cuber expects a list of encodings")
			}

			lastValue = pipelineSvc.cuberSvc.RunIncremental(input, pipe.CubeParams, pipe.SimplifyParams, pipe.SolveParams)

		case pipeline.CubeSelect:
			input, ok := lastValue.([]string)
			if !ok {
				log.Fatal("Cube selector expects a list of cubesets")
			}
			lastValue = pipelineSvc.cubeSelectorSvc.Run(input, pipe.CubeSelectParams)

		case pipeline.Solve:
			pipelineSvc.solverSvc.Run(lastValue.([]encoder.Encoding), false, pipe.SolveParams)

		case pipeline.SlurmSolve:
			input, ok := lastValue.([]encoder.Encoding)
			if !ok {
				log.Fatal("Slurm-based solver expects a slurm-based input")
			}

			pipelineSvc.solverSvc.Run(input, true, pipe.SolveParams)
		}
	})
}

func (pipelineSvc *PipelineService) Run(pipes []pipeline.Pipe) {
	pipelineSvc.Validate(pipes)
	pipelineSvc.RunPipes(pipes)
}
