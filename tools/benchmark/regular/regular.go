package regular

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"time"
)

func Run(commandContext *types.CommandContext) {
	// Count the number of instances for determining the progress
	instancesCount := utils.InstancesCount(commandContext)

	// Define the context
	benchmarkContext := &types.BenchmarkContext{
		Progress: make(map[string][]bool),
	}
	for _, satSolver_ := range commandContext.VariationsSatSolvers {
		satSolver := utils.ResolveSatSolverName(satSolver_)
		benchmarkContext.Progress[satSolver] = make([]bool, instancesCount)
	}

	// Loop through the instances
	utils.LoopThroughVariations(commandContext, func(i uint, satSolver_ string, steps uint, hash string, xorOption uint, adderType_ string, dobbertin uint) {
		for uint(benchmarkContext.RunningInstances) > commandContext.MaxConcurrentInstancesCount {
			time.Sleep(time.Second * 1)
		}

		adderType := utils.ResolveAdderType(adderType_)
		filepath := fmt.Sprintf("%smd4_%d_%s_xor%d_%s_dobbertin%d.cnf",
			constants.EncodingsDirPath, steps, adderType, xorOption, hash, dobbertin)

		satSolver := utils.ResolveSatSolverName(satSolver_)
		startTime := time.Now()
		switch satSolver {
		case constants.CryptoMiniSat:
			go core.CryptoMiniSat(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.Kissat:
			go core.Kissat(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.Cadical:
			go core.Cadical(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.MapleSat:
			go core.MapleSat(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.Glucose:
			go core.Glucose(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.XnfSat:
			go core.XnfSat(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		}

		benchmarkContext.RunningInstances += 1
	})

	for !core.AreAllInstancesCompleted(benchmarkContext) {
		time.Sleep(time.Second * 1)
	}
}
