package slurm

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os"
)

func generateJobs(context *types.CommandContext) {
	utils.LoopThroughVariations(context, func(i uint, satSolver_ string, steps uint, hash string, xorOption uint, adderType string, dobbertin uint) {
		instanceName := fmt.Sprintf("md4_%d_%s_xor%d_%s_dobbertin%d",
			steps, adderType, xorOption, hash, dobbertin)

		filepath := fmt.Sprintf("%s%s.cnf",
			constants.ENCODINGS_DIR_PATH, instanceName)

		satSolver := utils.ResolveSatSolverName(satSolver_)

		command := ""
		switch satSolver {
		case constants.CRYPTOMINISAT:
			command = core.CryptoMiniSatCmd(filepath)
		case constants.KISSAT:
			command = core.KissatCmd(filepath)
		case constants.CADICAL:
			command = core.CadicalCmd(filepath)
		case constants.MAPLESAT:
			command = core.MapleSatCmd(filepath)
		case constants.GLUCOSE:
			command = core.GlucoseCmd(filepath)
		}

		// Write the file for the job
		d := []byte("#!/bin/bash\n" + command)
		if err := os.WriteFile("./jobs/"+instanceName+".sh", d, 0644); err != nil {
			fmt.Println("Failed to create job:", instanceName)
		}
	})
}

func Run(context *types.CommandContext) {
	// TODO:Clean up the results and jobs directory
	// TODO: Generate jobs
	generateJobs(context)

	// TODO: Schedule the jobs
}
