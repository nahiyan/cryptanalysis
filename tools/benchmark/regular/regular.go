package regular

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/encodings"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"math/rand"
	"time"
)

func Run(context *types.CommandContext) {
	// Set the seed
	rand.Seed(context.Seed)

	// Generate encodings
	if context.GenerateEncodings == 1 {
		fmt.Println("Generating encodings")
		encodings.Generate(types.EncodingsGenContext{
			Variations: context.Variations,
			CubeParams: context.CubeParams,
		})
		fmt.Println("Done")
	}

	// TODO: Register session in the database

	// Count the number of instances for determining the progress
	instancesCount := utils.InstancesCount(context)

	// Define the context
	benchmarkContext := &types.BenchmarkContext{
		Progress: make(map[string][]bool),
	}
	for _, satSolver_ := range context.VariationsSatSolvers {
		satSolver := utils.ResolveSatSolverName(satSolver_)
		benchmarkContext.Progress[satSolver] = make([]bool, instancesCount)
	}

	// Loop through the instances
	utils.LoopThroughVariations(context, func(i uint, satSolver_ string, steps uint, hash string, xorOption uint, adderType string, dobbertin, dobbertinBits uint, cubeIndex *uint) {
		for uint(benchmarkContext.RunningInstances) >= context.MaxConcurrentInstancesCount {
			time.Sleep(time.Second * 1)
		}

		satSolver := utils.ResolveSatSolverName(satSolver_)
		startTime := time.Now()
		filepath := utils.EncodingsFileName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, cubeIndex)

		// Generate the subproblems from the cubes on the fly if cubes are enabled
		if cubeIndex != nil {
			instanceName := utils.InstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, nil)
			encodings.GenerateSubProblem(instanceName, int(*cubeIndex))
		}

		switch satSolver {
		case constants.CryptoMiniSat:
			go core.CryptoMiniSat(filepath, benchmarkContext, i, startTime, context.InstanceMaxTime)
		case constants.Kissat:
			go core.Kissat(filepath, benchmarkContext, i, startTime, context.InstanceMaxTime)
		case constants.Cadical:
			go core.Cadical(filepath, benchmarkContext, i, startTime, context.InstanceMaxTime)
		case constants.MapleSat:
			go core.MapleSat(filepath, benchmarkContext, i, startTime, context.InstanceMaxTime)
		case constants.Glucose:
			go core.Glucose(filepath, benchmarkContext, i, startTime, context.InstanceMaxTime)
		case constants.XnfSat:
			go core.XnfSat(filepath, benchmarkContext, i, startTime, context.InstanceMaxTime)
		}

		benchmarkContext.RunningInstances += 1
	})

	for !core.AreAllInstancesCompleted(benchmarkContext) {
		time.Sleep(time.Second * 1)
	}
}
