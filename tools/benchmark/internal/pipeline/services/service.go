package services

import (
	"benchmark/internal/consts"
	encoderServices "benchmark/internal/encoder/services"
	"benchmark/internal/pipeline"
	"benchmark/internal/solver"
	"log"
)

const (
	ListOfEncodings         = "list[encoding]"
	ListOfSlurmJobEncodings = "list[slurm_job[encoding]]"
	ListOfSolutions         = "list[solution]"
	ListOfSlurmJobSolutions = "list[slurm_job[solution]]"
	ListOfCubesets          = "list[cubeset]"
	ListOfSlurmJobCubesets  = "list[slurm_job[cubeset]]"
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
		return ListOfSlurmJobEncodings
	case pipeline.Cube:
		return ListOfEncodings
	case pipeline.SlurmCube:
		return ListOfSlurmJobEncodings
	case pipeline.CubeSelect:
		return ListOfCubesets
	case pipeline.SlurmCubeSelect:
		return ListOfSlurmJobCubesets
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
	case pipeline.SlurmCube:
		return ListOfSlurmJobCubesets
	case pipeline.CubeSelect:
		return ListOfEncodings
	case pipeline.SlurmCubeSelect:
		return ListOfSlurmJobEncodings
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
		Solvers: []solver.Solver{consts.Kissat, consts.MapleSat},
		Timeout: 5,
		Workers: 16,
	}

	cubeParameters := pipeline.Cubing{
		Timeout:    5,
		Thresholds: []int{2170, 2160},
	}

	cubeSelectParameters := pipeline.CubeSelecting{
		Type:     "random",
		Quantity: 3,
		Seed:     1,
	}

	newPipeline := make([]pipeline.Pipe, 0)
	for _, pipe := range pipes {
		newPipe := pipeline.Pipe{Type: pipe.Type}

		switch pipe.Type {
		case pipeline.Encode:
			newPipe.Encoding = encodeParameters
		case pipeline.Solve:
			newPipe.Solving = solveParameters
		case pipeline.SlurmSolve:
			newPipe.Solving = solveParameters
		case pipeline.Cube:
			newPipe.Cubing = cubeParameters
		case pipeline.SlurmCube:
			newPipe.Cubing = cubeParameters
		case pipeline.CubeSelect:
			newPipe.CubeSelecting = cubeSelectParameters
		case pipeline.SlurmCubeSelect:
			newPipe.CubeSelecting = cubeSelectParameters
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
			pipelineSvc.solverSvc.RunRegular(lastValue.([]pipeline.PromiseString), pipe.Solving)

		case pipeline.SlurmSolve:
			input, ok := lastValue.(pipeline.SlurmPipeOutput)
			if !ok {
				log.Fatal("Slurm-based solver expects a slurm-based input")
			}

			lastValue = pipelineSvc.solverSvc.RunSlurm(input, pipe.Solving)

		case pipeline.Cube:
			lastValue = pipelineSvc.cuberSvc.RunRegular(lastValue.([]string), pipe.Cubing)
		case pipeline.SlurmCube:
			lastValue = pipelineSvc.cuberSvc.RunSlurm(lastValue.(pipeline.SlurmPipeOutput), pipe.Cubing)
		case pipeline.CubeSelect:
			input, ok := lastValue.([]string)
			if !ok {
				log.Fatal("Cube selector expects a list of cubesets")
			}
			lastValue = pipelineSvc.cubeSelectorSvc.Run(input, pipe.CubeSelecting)
		case pipeline.Simplify:
			input, ok := lastValue.([]string)
			if !ok {
				log.Fatal("Simplifier expects a list of encodings")
			}
			lastValue = pipelineSvc.simplifierSvc.Run(input, pipe.Simplifying)
		}
	})
}

func (pipelineSvc *PipelineService) Run(pipes []pipeline.Pipe) {
	pipelineSvc.Validate(pipes)
	// pipelineSvc.TestRun(pipes)
	pipelineSvc.RealRun(pipes)
}
