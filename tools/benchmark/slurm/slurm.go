package slurm

import (
	"benchmark/constants"
	"benchmark/encodings"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
)

func writeSlurmJob(job types.SlurmJob, satSolver, instanceName string) (string, error) {
	template, err := template.New("job.tmpl").ParseFiles("slurm/job.tmpl")
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s%s_%s.sh", constants.JobsDirPath, utils.ResolveSatSolverName(satSolver), instanceName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}

	if err := template.Execute(file, job); err != nil {
		return "", err
	}

	return filePath, nil
}

func generateJobs(context *types.CommandContext) []string {
	filePaths := make([]string, 0)

	utils.LoopThroughVariations(context, func(i uint, satSolver_ string, steps uint, hash string, xorOption uint, adderType string, dobbertin, dobbertinBits uint, cubeIndex *uint) {
		instanceName := utils.InstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, cubeIndex)

		// Write the file for the job
		body := fmt.Sprintf("%s regular --var-steps %d --var-xor %d --var-dobbertin %d --var-dobbertin-bits %d --var-adders %s --var-hashes %s --var-sat-solvers %s", config.Get().Paths.Bin.Benchmark, steps, xorOption, dobbertin, dobbertinBits, adderType, hash, satSolver_)

		job := types.SlurmJob{
			Body: body,
		}
		job.Head.Nodes = 1
		job.Head.CpuCores = 1
		job.Head.Memory = 300
		job.Head.Time = context.InstanceMaxTime

		filePath, err := writeSlurmJob(job, satSolver_, instanceName)
		if err != nil {
			log.Fatal("Failed to create job", err)
		}

		filePaths = append(filePaths, filePath)
	})

	return filePaths
}

func Run(context *types.CommandContext) {
	// Generate encodings
	encodings.Generate(types.EncodingsGenContext{
		Variations:    context.Variations,
		IsCubeEnabled: context.IsCubeEnabled,
		CubeDepth:     context.CubeDepth,
	})

	// Generate jobs
	jobFilePaths := generateJobs(context)

	// Schedule the jobs
	for _, jobFilePath := range jobFilePaths {
		cmd := exec.Command("sbatch", jobFilePath)
		output, err := cmd.Output()
		fmt.Print(string(output))
		if err != nil {
			fmt.Println("Job schedule failed:", jobFilePath)
		}
	}
}
