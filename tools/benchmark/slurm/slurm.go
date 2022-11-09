package slurm

import (
	"benchmark/config"
	"benchmark/constants"
	"benchmark/db"
	"benchmark/encodings"
	"benchmark/types"
	"benchmark/utils"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/samber/lo"
	"gorm.io/gorm"
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
		body := fmt.Sprintf("%s regular --var-steps %d --var-xor %d --var-dobbertin %d --var-dobbertin-bits %d --var-adders %s --var-hashes %s --var-sat-solvers %s --generate-encodings 0", config.Get().Paths.Bin.Benchmark, steps, xorOption, dobbertin, dobbertinBits, adderType, hash, satSolver_)

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

func getScheduledJobCount() uint {
	output, err := exec.Command("sq").Output()
	if err != nil {
		return 0
	}

	reader := bytes.NewReader(output)
	lines, err := utils.CountLines(reader)
	if err != nil {
		return 0
	}

	return uint(lines - 1)
}

func schedule(tx *gorm.DB, jobFilePaths []string) error {
	areAllJobsScheduled := true

	// Use DB transaction to lock the database
	if err := tx.Transaction(func(_ *gorm.DB) error {
		for _, jobFilePath := range jobFilePaths {
			// Schedule
			cmd := exec.Command("sbatch", jobFilePath)
			output_, err := cmd.Output()
			output := string(output_)
			fmt.Print(output)
			if err != nil {
				areAllJobsScheduled = false
				// Commit the transaction
				return nil
			}

			// Grab the job ID
			jobId := utils.GetJobId(output)

			// Schedule trigger upon completion
			if err := exec.Command(fmt.Sprintf("strigger --set --jobid=%d --fini --program=\"%s slurm --digest %d\"", jobId, config.Get().Paths.Bin.Benchmark, jobId)).Run(); err != nil {
				return errors.New("failed to set strigger")
			}

			// Register the Slurm ID in the database
			if err := tx.Model(&types.Job{}).Where("file_name", path.Base(jobFilePath)).Update("slurm_id", jobId).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if !areAllJobsScheduled {
		return errors.New(constants.ErrOneJobScheduleFailed)
	}

	return nil
}

func scheduleToLimit() error {
	if err := db.Get().Transaction(func(tx *gorm.DB) error {
		// TODO: Use the max concurrent instance parameter for slurm too
		if emptySpots := config.Get().Slurm.MaxJobs - getScheduledJobCount(); emptySpots > 0 {
			pendingJobs := make([]types.Job, 0)

			if err := tx.Limit(int(emptySpots)).Find(pendingJobs).Error; err != nil {
				return err
			}

			pendingJobFilePaths := lo.Map(pendingJobs, func(job types.Job, i int) string {
				return fmt.Sprintf("%s%s", constants.JobsDirPath, job.FileName)
			})
			if err := schedule(tx, pendingJobFilePaths); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		// Commit transaction even if one of the jobs failed
		if err.Error() == constants.ErrOneJobScheduleFailed {
			return nil
		}

		return err
	}

	return nil
}

func Run(context *types.CommandContext) {
	// Handle any job that needs to be digested
	if context.Digest != 0 {
		if err := db.Get().Where("slurm_id", context.Digest).Delete(&types.Job{}).Error; err != nil {
			log.Fatal(err)
		}

		// Schedule more jobs
		if err := scheduleToLimit(); err != nil {
			log.Fatal(err)
		}

		return
	}

	// Generate encodings
	if context.GenerateEncodings == 1 {
		fmt.Println("Generating encodings")
		encodings.Generate(types.EncodingsGenContext{
			Variations:    context.Variations,
			IsCubeEnabled: context.IsCubeEnabled,
			CubeVars:      context.CubeVars,
		})
		fmt.Println("Done")
	}

	// Generate jobs
	jobFilePaths := generateJobs(context)

	// Add to the database
	db.Get().Create(lo.Map(jobFilePaths, func(filePath string, i int) *types.Job {
		return &types.Job{
			FileName: path.Base(filePath),
		}
	}))

	// Schedule the jobs up to the limit
	if err := scheduleToLimit(); err != nil {
		log.Fatal(err)
	}
}
