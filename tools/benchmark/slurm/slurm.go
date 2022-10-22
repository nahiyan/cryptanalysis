package slurm

import (
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os"
	"os/exec"
)

func generateJobs(context *types.CommandContext) []string {
	filepaths := make([]string, 0)

	utils.LoopThroughVariations(context, func(i uint, satSolver_ string, steps uint, hash string, xorOption uint, adderType string, dobbertin, dobbertinBits uint) {
		instanceName := fmt.Sprintf("md4_%d_%s_xor%d_%s_dobbertin%d",
			steps, adderType, xorOption, hash, dobbertin)

		slurmArgs := fmt.Sprintf("#SBATCH --nodes=1\n#SBATCH --cpus-per-task=1\n#SBATCH --mem=300M\n#SBATCH --time=00:%d\n", context.InstanceMaxTime)

		command := fmt.Sprintf("%s\n./benchmark --var-steps %d --var-xor %d --var-dobbertin %d --var-adders %s --var-hashes %s --var-sat-solvers %s regular", slurmArgs, steps, xorOption, dobbertin, adderType, hash, satSolver_)

		satSolver := utils.ResolveSatSolverName(satSolver_)

		// Write the file for the job
		d := []byte("#!/bin/bash\n\n" + command)
		filepath := "./jobs/" + satSolver + "_" + instanceName + ".sh"
		if err := os.WriteFile(filepath, d, 0644); err != nil {
			fmt.Println("Failed to create job:", instanceName)
		}

		filepaths = append(filepaths, filepath)
	})

	return filepaths
}

func Run(context *types.CommandContext) {
	// Generate jobs
	jobFilePaths := generateJobs(context)

	// TODO: Schedule the jobs
	for _, jobFilePath := range jobFilePaths {
		cmd := exec.Command("sbatch", jobFilePath)
		output, err := cmd.Output()
		fmt.Print(string(output))
		if err != nil {
			fmt.Println("Job schedule failed:", jobFilePath)
		}
	}
}
