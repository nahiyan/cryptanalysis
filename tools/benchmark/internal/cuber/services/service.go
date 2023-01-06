package services

import (
	cubeslurmtask "benchmark/internal/cube_slurm_task"
	"benchmark/internal/cubeset"
	"benchmark/internal/pipeline"
	"benchmark/internal/slurm"
	"context"
	"fmt"
	"log"
	"math"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/alitto/pond"
)

func (cuberSvc *CuberService) CubesFilePath(encoding string, threshold int) string {
	encodingDir, fileName := path.Split(encoding)
	cubesFilePath := path.Join(encodingDir, fmt.Sprintf("%s.n%d.cubes", fileName, threshold))

	return cubesFilePath
}

func (cuberSvc *CuberService) ShouldSkip(encoding string, threshold int) bool {
	filesystemSvc := cuberSvc.filesystemSvc
	cubesFilePath := cuberSvc.CubesFilePath(encoding, threshold)

	return filesystemSvc.FileExists(cubesFilePath)
}

func (cuberSvc *CuberService) ReadMarchOutput(output string) (int, int, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "c number of cubes") {
			continue
		}

		var cubes, refutedLeaves int
		_, err := fmt.Sscanf(line, "c number of cubes %d, including %d refuted leaves", &cubes, &refutedLeaves)
		if err != nil {
			return 0, 0, err
		}

		return cubes, refutedLeaves, nil
	}

	return 0, 0, nil
}

func (cuberSvc *CuberService) TrackedInvoke(encoding string, threshold int, timeout time.Duration) error {
	cubesetSvc := cuberSvc.cubesetSvc
	output, runtime, err := cuberSvc.Invoke(encoding, threshold, timeout)
	if err != nil {
		return err
	}

	cubes, refutedLeaves, err := cuberSvc.ReadMarchOutput(output)
	if err != nil {
		return err
	}

	fmt.Println("Cuber:", threshold, cubes, "cubes", refutedLeaves, "refuted leaves", runtime, encoding)

	cubesFilePath := cuberSvc.CubesFilePath(encoding, threshold)
	err = cubesetSvc.Register(cubesFilePath, cubeset.CubeSet{
		Cubes:         cubes,
		RefutedLeaves: refutedLeaves,
		Runtime:       runtime,
	})

	return err
}

func (cuberSvc *CuberService) Invoke(encoding string, threshold int, timeout time.Duration) (string, time.Duration, error) {
	config := cuberSvc.configSvc.Config

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cubesFilePath := cuberSvc.CubesFilePath(encoding, threshold)
	cmd := exec.CommandContext(ctx, config.Paths.Bin.March, encoding, "-o", cubesFilePath, "-n", strconv.Itoa(threshold))

	startTime := time.Now()
	output_, err := cmd.Output()
	if err != nil {
		return "", time.Duration(0), err
	}
	runtime := time.Since(startTime)

	output := string(output_)
	return output, runtime, err
}

func (cuberSvc *CuberService) Loop(encodings []string, parameters pipeline.Cubing, handler func(encoding string, threshold int, timeout int)) {
	for _, encoding := range encodings {
		thresholds := parameters.Thresholds
		if len(thresholds) == 0 {
			freeVariables, _, err := cuberSvc.encodingSvc.Process(encoding)
			cuberSvc.errorSvc.Fatal(err, "Cuber: failed to process the encoding")

			stepSize := 10
			// Initial threshold is the nearest multiple of step size less than the num. of free variables
			threshold := int(math.Floor(float64(freeVariables)/float64(stepSize)) * float64(stepSize))

			for {
				if cuberSvc.ShouldSkip(encoding, threshold) {
					fmt.Println("Cuber: skipped", threshold, encoding)
					continue
				}

				thresholds = append(thresholds, int(threshold))
				threshold -= stepSize

				if threshold < 0 {
					break
				}
			}
		}

		for _, threshold := range thresholds {
			handler(encoding, threshold, parameters.Timeout)
		}
	}
}

func (cuberSvc *CuberService) RunRegular(encodings []string, parameters pipeline.Cubing) []string {
	cubesFilePaths := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	fmt.Println("Cuber: started")

	cuberSvc.Loop(encodings, parameters, func(encoding string, threshold, timeout int) {
		cubesFilePaths = append(cubesFilePaths, cuberSvc.CubesFilePath(encoding, threshold))

		pool.Submit(func(encoding string, threshold int) func() {
			return func() {
				err := cuberSvc.TrackedInvoke(encoding, threshold, time.Duration(parameters.Timeout)*time.Second)
				cuberSvc.errorSvc.Fatal(err, "Cuber: failed to cube")
			}
		}(encoding, threshold))
	})

	pool.StopAndWait()
	fmt.Println("Cuber: stopped")
	return cubesFilePaths
}

func (cuberSvc *CuberService) RunSlurm(previousPipeOutput pipeline.SlurmPipeOutput, parameters pipeline.Cubing) pipeline.SlurmPipeOutput {
	errorSvc := cuberSvc.errorSvc
	slurmSvc := cuberSvc.slurmSvc
	config := cuberSvc.configSvc.Config
	encodings, ok := previousPipeOutput.Values.([]string)
	if !ok {
		log.Fatal("Cuber: invalid input")
	}
	dependencies := previousPipeOutput.Jobs

	err := cuberSvc.cubeSlurmTaskSvc.RemoveAll()
	errorSvc.Fatal(err, "Cuber: failed to clear slurm tasks")

	counter := 1
	cuberSvc.Loop(encodings, parameters, func(encoding string, threshold int, timeout int) {
		if cuberSvc.ShouldSkip(encoding, threshold) {
			fmt.Println("Cuber: skipped", encoding, "with threshold of", threshold)
			return
		}

		err := cuberSvc.cubeSlurmTaskSvc.AddTask(counter, cubeslurmtask.Task{
			Encoding:  encoding,
			Threshold: threshold,
			Timeout:   time.Duration(parameters.Timeout) * time.Second,
		})
		errorSvc.Fatal(err, "Cuber: failed to add slurm task")

		counter++
	})

	fmt.Println("Cuber: added", counter-1, "slurm tasks")

	numTasks := counter
	timeout := parameters.Timeout
	jobFilePath, err := slurmSvc.GenerateJob(
		numTasks,
		1,
		1,
		1024,
		timeout,
		fmt.Sprintf(
			"%s slurm-task -t solve -i ${SLURM_ARRAY_TASK_ID}",
			config.Paths.Bin.Benchmark))
	errorSvc.Fatal(err, "Solver: failed to create slurm job file")

	jobId, err := slurmSvc.ScheduleJob(jobFilePath, dependencies)
	errorSvc.Fatal(err, "Solver: failed to schedule the job")
	fmt.Println("Cuber: scheduled job with ID", jobId)

	return pipeline.SlurmPipeOutput{
		Jobs:   []slurm.Job{{Id: jobId}},
		Values: []string{},
	}
}
