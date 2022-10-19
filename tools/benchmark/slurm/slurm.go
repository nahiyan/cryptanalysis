package slurm

import (
	"benchmark/constants"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os"
	"os/exec"
)

func generateJobs(context *types.CommandContext) []string {
	filepaths := make([]string, 0)

	utils.LoopThroughVariations(context, func(i uint, satSolver_ string, steps uint, hash string, xorOption uint, adderType string, dobbertin uint) {
		instanceName := fmt.Sprintf("md4_%d_%s_xor%d_%s_dobbertin%d",
			steps, adderType, xorOption, hash, dobbertin)

		slurmArgs := fmt.Sprintf("#SBATCH --nodes=1\n#SBATCH --cpus-per-task=1\n#SBATCH --mem=300M\n#SBATCH --time=00:%d\n", context.InstanceMaxTime)

		command := fmt.Sprintf("%s\n./benchmark --var-steps %d --var-xor %d --var-dobbertin %d --var-adders %s --var-hashes %s --var-sat-solvers %s --reset-data 0 regular", slurmArgs, steps, xorOption, dobbertin, adderType, hash, satSolver_)

		// Write the file for the job
		d := []byte("#!/bin/bash\n\n" + command)
		filepath := "./jobs/" + instanceName + ".sh"
		if err := os.WriteFile(filepath, d, 0644); err != nil {
			fmt.Println("Failed to create job:", instanceName)
		}

		filepaths = append(filepaths, filepath)
	})

	return filepaths
}

func Run(context *types.CommandContext) {
	// TODO:Clean up the results and jobs directory
	os.Remove(constants.BENCHMARK_LOG_FILE_NAME)
	os.Remove(constants.VERIFICATION_LOG_FILE_NAME)

	// Generate jobs
	jobFilePaths := generateJobs(context)

	// TODO: Schedule the jobs
	for _, jobFilePath := range jobFilePaths {
		cmd := exec.Command("sbatch", jobFilePath)
		output, err := cmd.Output()
		fmt.Println(output)
		if err != nil {
			fmt.Println("Job schedule failed:", jobFilePath)
		}
	}
}
