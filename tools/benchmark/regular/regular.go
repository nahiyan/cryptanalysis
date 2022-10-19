package regular

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/samber/lo"
)

func Run(commandContext *types.CommandContext) {
	// Count the number of instances for determining the progress
	instancesCount := len(commandContext.VariationsXor) * len(commandContext.VariationsHashes) * len(commandContext.VariationsAdders) * len(commandContext.VariationsSteps) * (len(commandContext.VariationsDobbertin) * len(lo.Filter(commandContext.VariationsSteps, func(steps uint, _ int) bool {
		return steps >= 27
	})))

	// Define the context
	benchmarkContext := &types.BenchmarkContext{
		Progress: make(map[string][]bool),
	}
	for _, satSolver_ := range commandContext.VariationsSatSolvers {
		satSolver := utils.ResolveSatSolverName(satSolver_)
		benchmarkContext.Progress[satSolver] = make([]bool, instancesCount)
	}

	// Remove the files from previous execution
	os.Remove(constants.BENCHMARK_LOG_FILE_NAME)
	os.Remove(constants.VERIFICATION_LOG_FILE_NAME)
	for _, satSolver_ := range commandContext.VariationsSatSolvers {
		satSolver := utils.ResolveSatSolverName(satSolver_)

		cmd := exec.Command("bash", "-c", fmt.Sprintf("rm %s%s/*.sol", constants.SOLUTIONS_DIR_PATH, satSolver))
		if err := cmd.Run(); err != nil {
			// fmt.Println(cmd.String())
			fmt.Println("Failed to delete the solution files: " + err.Error())
		}
	}

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
}
