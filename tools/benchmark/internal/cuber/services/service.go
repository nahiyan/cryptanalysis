package services

import (
	"benchmark/internal/command"
	"benchmark/internal/cuber"
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/alitto/pond"
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
	cubesetFileName := path.Base(encoding) + fmt.Sprintf(".march_n%d.cubes", threshold)
	cubesFilePath := path.Join(cuberSvc.configSvc.Config.Paths.Cubesets, cubesetFileName)
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

func (cuberSvc *CuberService) Loop(encodings []encoder.Encoding, parameters pipeline.CubeParams, handler func(encoding string, threshold int, timeout int)) {
	for _, encoding := range encodings {
		thresholds := parameters.Thresholds
		if len(thresholds) == 0 {
			encodingInfo, err := cuberSvc.encodingSvc.GetInfo(encoding.BasePath)
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
			handler(encoding.BasePath, threshold, parameters.Timeout)
		}
	}
}

func (cuberSvc *CuberService) RunRegular(encodings []encoder.Encoding, parameters pipeline.CubeParams) []string {
	err := cuberSvc.filesystemSvc.PrepareDirs([]string{"cubesets", "encodings", "logs"})
	cuberSvc.errorSvc.Fatal(err, "Cuber: failed to prepare the required dirs")

	cubesFilePaths := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	shouldStop := map[string]bool{}
	commandGrps := map[string]*command.Group{}
	logrus.Println("Cuber: started")

	cuberSvc.Loop(encodings, parameters, func(encoding string, threshold, timeout int) {
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
