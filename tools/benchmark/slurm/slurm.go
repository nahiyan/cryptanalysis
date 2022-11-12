package slurm

import (
	"benchmark/config"
	"benchmark/constants"
	"benchmark/encodings"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"html/template"
	"log"
	"math/rand"
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
		body := fmt.Sprintf("%s regular --max-time %d  --var-steps %d --var-xor %d --var-dobbertin %d --var-dobbertin-bits %d --var-adders %s --var-hashes %s --var-sat-solvers %s --generate-encodings 0 %s --seed %d", config.Get().Paths.Bin.Benchmark, context.InstanceMaxTime, steps, xorOption, dobbertin, dobbertinBits, adderType, hash, satSolver_, func(cubeParams *types.CubeParams) string {
			if cubeParams == nil {
				return ""
			}

			return fmt.Sprintf("--cube --cube-index %d", *cubeIndex)
		}(context.CubeParams), context.Seed)

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

func scheduleJobs(paths []string) {
	for _, jobFilePath := range paths {
		// Schedule
		cmd := exec.Command("sbatch", jobFilePath)
		output_, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		output := string(output_)
		fmt.Print(output)

		// Grab the job ID
		// jobId := utils.GetJobId(output)

		// Register the Slurm ID in the database
		// if err := tx.Model(&types.Job{}).Where("file_name", path.Base(jobFilePath)).Update("slurm_id", jobId).Error; err != nil {
		// 	return err
		// }
	}
}

func Run(context *types.CommandContext) {
	// Set the seed
	rand.Seed(context.Seed)

	// TODO: Reduce redudant code
	// Generate encodings
	if context.GenerateEncodings == 1 {
		fmt.Println("Generating encodings")
		encodings.Generate(types.EncodingsGenContext{
			Variations: context.Variations,
			CubeParams: context.CubeParams,
		})
		fmt.Println("Done")
	}

	// Generate jobs
	jobFilePaths := generateJobs(context)

	// Add to the database
	// db.Get().Create(lo.Map(jobFilePaths, func(filePath string, i int) *types.Job {
	// 	return &types.Job{
	// 		FileName: path.Base(filePath),
	// 	}
	// }))

	// Schedule the jobs
	scheduleJobs(jobFilePaths)
}
