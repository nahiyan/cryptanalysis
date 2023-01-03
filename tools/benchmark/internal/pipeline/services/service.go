package services

import (
	"benchmark/internal/consts"
	encoderServices "benchmark/internal/encoder/services"
	"benchmark/internal/pipeline"
	"benchmark/internal/solver"
)

const (
	ListOfEncodings = "list[encoding]"
	ListOfSolutions = "list[solution]"
	ListOfCubesets  = "list[cubeset]"
	None            = "none"
)

type InputOutputType string
type LoopHandler func(*pipeline.Pipe, *pipeline.Pipe)

func getInputType(pipe *pipeline.Pipe) InputOutputType {
	switch pipe.Type {
	case pipeline.Encode:
		return None
	case pipeline.Solve:
		return ListOfEncodings
	case pipeline.Cube:
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
	case pipeline.Cube:
		return ListOfCubesets
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
			panic("Incompatible pipeline: " + outputType + " can't fit into the expected input type " + nextPipelineInputType)
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

func (pipelineSvc *PipelineService) TestRun(pipes []pipeline.Pipe) {
	encodeParameters := pipeline.Encoding{
		Xor:           []int{0},
		Dobbertin:     []int{0},
		DobbertinBits: []int{32},
		Adders:        []pipeline.AdderType{encoderServices.Espresso},
		Hashes:        []string{"ffffffffffffffffffffffffffffffff", "00000000000000000000000000000000"},
		Steps:         []int{16},
	}

	solveParameters := pipeline.Solving{
		Solvers:  []solver.Solver{consts.Kissat, consts.MapleSat},
		Timeout:  5,
		Platform: consts.Slurm,
		Workers:  16,
	}

	cubeParameters := pipeline.Cubing{
		Platform:   consts.General,
		Timeout:    5,
		Thresholds: []int{2170, 2160},
	}

	newPipeline := make([]pipeline.Pipe, 0)
	for _, pipe := range pipes {
		newPipe := pipeline.Pipe{Type: pipe.Type}

		switch pipe.Type {
		case pipeline.Encode:
			newPipe.Encoding = encodeParameters
		case pipeline.Solve:
			newPipe.Solving = solveParameters
		case pipeline.Cube:
			newPipe.Cubing = cubeParameters
		}

		newPipeline = append(newPipeline, newPipe)
	}

	pipelineSvc.RealRun(newPipeline)
}

func (pipelineSvc *PipelineService) RealRun(pipes []pipeline.Pipe) {
	var lastValue interface{}
	pipelineSvc.Loop(pipes, func(pipe, nextPipe *pipeline.Pipe) {
		switch pipe.Type {
		case pipeline.Encode:
			lastValue = pipelineSvc.encoderSvc.Run(encoderServices.SaeedE, pipe.Encoding)

			if nextPipe == nil {
				return
			}

		case pipeline.Solve:
			pipelineSvc.solverSvc.Run(lastValue.([]string), pipe.Solving)

		case pipeline.Cube:
			lastValue = pipelineSvc.cuberSvc.Run(lastValue.([]string), pipe.Cubing)
		}
	})
}

func (pipelineSvc *PipelineService) Run(pipes []pipeline.Pipe) {
	pipelineSvc.Validate(pipes)
	pipelineSvc.TestRun(pipes)
	// pipelineSvc.RealRun(pipes)
}
