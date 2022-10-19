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
			constants.ENCODINGS_DIR_PATH, steps, adderType, xorOption, hash, dobbertin)

		satSolver := utils.ResolveSatSolverName(satSolver_)
		startTime := time.Now()
		switch satSolver {
		case constants.CRYPTOMINISAT:
			go core.CryptoMiniSat(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.KISSAT:
			go core.Kissat(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.CADICAL:
			go core.Cadical(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.MAPLESAT:
			go core.MapleSat(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		case constants.GLUCOSE:
			go core.Glucose(filepath, benchmarkContext, i, startTime, commandContext.InstanceMaxTime)
		}

		benchmarkContext.RunningInstances += 1
	})

	for !core.AreAllInstancesCompleted(benchmarkContext) {
		time.Sleep(time.Second * 1)
	}
}
