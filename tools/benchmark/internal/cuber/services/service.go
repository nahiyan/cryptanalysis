package services

import (
	"benchmark/internal/command"
	cubeslurmtask "benchmark/internal/cube_slurm_task"
	"benchmark/internal/cuber"
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"benchmark/internal/slurm"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/alitto/pond"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus"
)

type InvokeParameters struct {
	Encoding         string
	Threshold        int
	Timeout          time.Duration
	MaxCubes         int
	MinRefutedLeaves int
}

type InvokeControl struct {
	CommandGroup *command.Group
	CubesetPaths *[]string
	ShouldStop   map[string]bool
}

func (cuberSvc *CuberService) getLogFilePath(cubeFilePath string) string {
	logFileName := path.Base(cubeFilePath[:len(cubeFilePath)-5] + "log")
	logFilePath := path.Join(cuberSvc.configSvc.Config.Paths.Logs, logFileName)
	return logFilePath
}

func (cuberSvc *CuberService) CubesFilePath(encoding string, threshold int) string {
	info, err := cuberSvc.encoderSvc.ProcessInstanceName(path.Base(encoding))
	cuberSvc.errorSvc.Fatal(err, "Cuber: failed to process encoding")
	info.Cubing = mo.Some(encoder.CubingInfo{
		Threshold: threshold,
	})
	newInstanceName := cuberSvc.encoderSvc.GetInstanceName(info)
	cubesFilePath := path.Join(cuberSvc.configSvc.Config.Paths.Cubesets, newInstanceName)
	return cubesFilePath
}

func (cuberSvc *CuberService) ShouldSkip(encoding string, threshold int) bool {
	cubesFilePath := cuberSvc.CubesFilePath(encoding, threshold)
	logFilePath := cuberSvc.getLogFilePath(cubesFilePath)

	if exists := cuberSvc.filesystemSvc.FileExistsNonEmpty(cubesFilePath); !exists {
		return false
	}

	_, _, _, err := cuberSvc.ParseOutput(logFilePath)
	return err == nil
}

// func (cuberSvc *CuberService) ReadMarchOutput(output string) (int, int, error) {
// 	lines := strings.Split(output, "\n")
// 	for _, line := range lines {
// 		if !strings.HasPrefix(line, "c number of cubes") {
// 			continue
// 		}

// 		var cubes, refutedLeaves int
// 		_, err := fmt.Sscanf(line, "c number of cubes %d, including %d refuted leaves", &cubes, &refutedLeaves)
// 		if err != nil {
// 			return 0, 0, err
// 		}

// 		return cubes, refutedLeaves, nil
// 	}

// 	return 0, 0, nil
// }

func (cuberSvc *CuberService) TrackedInvoke(parameters InvokeParameters, control InvokeControl) error {
	if shouldStop_, exists := control.ShouldStop[parameters.Encoding]; exists && shouldStop_ {
		return nil
	}

	cubesFilePath, logFilePath, err := cuberSvc.Invoke(parameters, control)
	if err != nil {
		return err
	}

	processTime, numCubes, numRefutedLeaves, err := cuberSvc.ParseOutput(logFilePath)
	if err != nil {
		return err
	}

	// instanceName := strings.TrimSuffix(parameters.Encoding, ".cnf")
	// cubesFilePath := cuberSvc.CubesFilePath(parameters.Encoding, parameters.Threshold)
	// err = cubesetSvc.Register(cubesFilePath, cubeset.CubeSet{
	// 	Threshold:     parameters.Threshold,
	// 	InstanceName:  instanceName,
	// 	Cubes:         cubes,
	// 	RefutedLeaves: refutedLeaves,
	// 	Runtime:       runtime,
	// })
	// if err != nil {
	// 	return err
	// }

	maxCubesExceeded := parameters.MaxCubes > 0 && numCubes > parameters.MaxCubes
	minRefutedLeavesViolated := parameters.MinRefutedLeaves > 0 && numRefutedLeaves < parameters.MinRefutedLeaves
	hasViolated := maxCubesExceeded || minRefutedLeavesViolated

	cuberSvc.logSvc.CubeResult(cubesFilePath, processTime, numCubes, numRefutedLeaves, hasViolated)

	if maxCubesExceeded {
		control.ShouldStop[parameters.Encoding] = true
		cuberSvc.commandSvc.StopGroup(control.CommandGroup)
		// logrus.Println("Cuber: Written stop signal", parameters.Threshold, parameters.Encoding)
	}

	if hasViolated {
		if err := os.Remove(cubesFilePath); err != nil {
			return err
		}
		// logrus.Println("Cuber: removed", cubesFilePath)
		return cuber.ErrCubesetViolatedConstraints
	}

	*control.CubesetPaths = append(*control.CubesetPaths, cuberSvc.CubesFilePath(parameters.Encoding, parameters.Threshold))

	return nil
}

func (cuberSvc *CuberService) Invoke(parameters InvokeParameters, control InvokeControl) (string, string, error) {
	config := cuberSvc.configSvc.Config

	ctx, cancel := context.WithTimeout(context.Background(), parameters.Timeout)
	defer cancel()

	cubesFilePath := cuberSvc.CubesFilePath(parameters.Encoding, parameters.Threshold)
	cmd := exec.CommandContext(ctx, config.Paths.Bin.March, parameters.Encoding, "-o", cubesFilePath, "-n", strconv.Itoa(parameters.Threshold))
	cuberSvc.commandSvc.AddToGroup(control.CommandGroup, cmd)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	logFilePath := cuberSvc.getLogFilePath(cubesFilePath)
	cuberSvc.filesystemSvc.WriteFromPipe(stdoutPipe, logFilePath)
	if err := cmd.Wait(); err != nil {
		return "", "", err
	}

	return cubesFilePath, logFilePath, err
}

func (cuberSvc *CuberService) Loop(encodingPromises []pipeline.EncodingPromise, parameters pipeline.Cubing, handler func(encoding string, threshold int, timeout int)) {
	for _, promise := range encodingPromises {
		encoding := promise.Get(map[string]interface{}{})
		thresholds := parameters.Thresholds
		if len(thresholds) == 0 {
			encodingInfo, err := cuberSvc.encodingSvc.GetInfo(encoding)
			freeVariables := encodingInfo.FreeVariables
			cuberSvc.errorSvc.Fatal(err, "Cuber: failed to process the encoding")

			stepSize := 10
			// Initial threshold is the nearest multiple of step size less than the num. of free variables
			threshold := int(math.Floor(float64(freeVariables)/float64(stepSize)) * float64(stepSize))

			for {
				thresholds = append(thresholds, int(threshold))
				threshold -= stepSize

				if threshold <= 0 {
					break
				}
			}
		}

		for _, threshold := range thresholds {
			handler(encoding, threshold, parameters.Timeout)
		}
	}
}

func (cuberSvc *CuberService) RunRegular(encodingPromises []pipeline.EncodingPromise, parameters pipeline.Cubing) []string {
	err := cuberSvc.filesystemSvc.PrepareDirs([]string{"cubesets", "encodings", "logs"})
	cuberSvc.errorSvc.Fatal(err, "Cuber: failed to prepare the required dirs")

	cubesFilePaths := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	shouldStop := map[string]bool{}
	commandGrps := map[string]*command.Group{}
	logrus.Println("Cuber: started")

	cuberSvc.Loop(encodingPromises, parameters, func(encoding string, threshold, timeout int) {
		if cuberSvc.ShouldSkip(encoding, threshold) {
			logrus.Println("Cuber: skipped", threshold, encoding)
			cubesFilePaths = append(cubesFilePaths, cuberSvc.CubesFilePath(encoding, threshold))
			return
		}

		if _, exists := commandGrps[encoding]; !exists {
			commandGrps[encoding] = cuberSvc.commandSvc.CreateGroup()
		}

		pool.Submit(func(encoding string, threshold int) func() {
			return func() {
				err := cuberSvc.TrackedInvoke(InvokeParameters{
					Encoding:         encoding,
					Threshold:        threshold,
					Timeout:          time.Duration(timeout) * time.Second,
					MaxCubes:         parameters.MaxCubes,
					MinRefutedLeaves: parameters.MinRefutedLeaves,
				}, InvokeControl{
					CommandGroup: commandGrps[encoding],
					CubesetPaths: &cubesFilePaths,
					ShouldStop:   shouldStop,
				})

				if err != nil {
					if err == cuber.ErrCubesetViolatedConstraints || shouldStop[encoding] {
						return
					}
				}
				cuberSvc.errorSvc.Fatal(err, "Cuber: failed to cube "+strconv.Itoa(threshold)+" "+encoding)
			}
		}(encoding, threshold))
	})

	pool.StopAndWait()
	logrus.Println("Cuber: stopped")
	return cubesFilePaths
}

// TODO: Remove this
func (cuberSvc *CuberService) RunSlurm(previousPipeOutput pipeline.SlurmPipeOutput, parameters pipeline.Cubing) pipeline.SlurmPipeOutput {
	errorSvc := cuberSvc.errorSvc
	slurmSvc := cuberSvc.slurmSvc
	config := cuberSvc.configSvc.Config
	encodingPromises, ok := previousPipeOutput.Values.([]pipeline.EncodingPromise)
	if !ok {
		log.Fatal("Cuber: invalid input")
	}
	dependencies := previousPipeOutput.Jobs

	err := cuberSvc.cubeSlurmTaskSvc.RemoveAll()
	errorSvc.Fatal(err, "Cuber: failed to clear slurm tasks")

	counter := 1
	cuberSvc.Loop(encodingPromises, parameters, func(encoding string, threshold int, timeout int) {
		if cuberSvc.ShouldSkip(encoding, threshold) {
			logrus.Println("Cuber: skipped", threshold, encoding)
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

	logrus.Println("Cuber: added", counter-1, "slurm tasks")

	numTasks := counter - 1
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
	logrus.Println("Cuber: scheduled job with ID", jobId)

	return pipeline.SlurmPipeOutput{
		Jobs:   []slurm.Job{{Id: jobId}},
		Values: []string{},
	}
}
