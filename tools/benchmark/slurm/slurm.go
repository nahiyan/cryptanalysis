package slurm

import (
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os"
)

func generateJobs(context *types.CommandContext) {
	utils.LoopThroughVariations(context, func(i uint, satSolver_ string, steps uint, hash string, xorOption uint, adderType string, dobbertin uint) {
		instanceName := fmt.Sprintf("md4_%d_%s_xor%d_%s_dobbertin%d",
			steps, adderType, xorOption, hash, dobbertin)

		// filepath := fmt.Sprintf("%s%s.cnf", constants.ENCODINGS_DIR_PATH, instanceName)

		slurmArgs := fmt.Sprintf("#SBATCH --nodes=1\n#SBATCH --cpus-per-task=1\n#SBATCH --mem=300M\n#SBATCH --time=00:%d\n", context.InstanceMaxTime)

		command := fmt.Sprintf("%s\n./benchmark --var-steps %d --var-xor %d --var-dobbertin %d --var-adders %s --var-hashes %s --var-sat-solvers %s --reset-data 0 regular", slurmArgs, steps, xorOption, dobbertin, adderType, hash, satSolver_)

		// Write the file for the job
		d := []byte("#!/bin/bash\n\n" + command)
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
