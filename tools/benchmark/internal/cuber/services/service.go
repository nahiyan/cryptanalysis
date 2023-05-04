package services

import (
	"benchmark/internal/command"
	"benchmark/internal/cuber"
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/alitto/pond"
)

type InvokeParameters struct {
	Encoding         string
	ThresholdType    cuber.ThresholdType
	Threshold        int
	Timeout          time.Duration
	MaxCubes         int
	MinCubes         int
	MinRefutedLeaves int
	Suffix           string
	MaxVariable      int
	SkipLogs         bool
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

func (cuberSvc *CuberService) CubesFilePath(encoding string, thresholdType cuber.ThresholdType, threshold int, suffix string) string {
	thresholdChar := "n"
	if thresholdType == cuber.CutoffDepth {
		thresholdChar = "d"
	}

	var cubesetFileName string
	if len(suffix) == 0 {
		cubesetFileName = path.Base(encoding) + fmt.Sprintf(".march_%s%d.cubes", thresholdChar, threshold)
	} else {
		newBaseEncoding := regexp.MustCompile(`.march.+`).ReplaceAllString(path.Base(encoding), "")
		cubesetFileName = newBaseEncoding + fmt.Sprintf(".march_%s.cubes", suffix)
	}

	cubesFilePath := path.Join(cuberSvc.configSvc.Config.Paths.Cubesets, cubesetFileName)
	return cubesFilePath
}

func (cuberSvc *CuberService) ShouldSkip(encoding string, thresholdType cuber.ThresholdType, threshold int, hash string) bool {
	cubesFilePath := cuberSvc.CubesFilePath(encoding, thresholdType, threshold, hash)
	logFilePath := cuberSvc.getLogFilePath(cubesFilePath)

	if exists := cuberSvc.filesystemSvc.FileExistsNonEmpty(cubesFilePath); !exists {
		return false
	}

	var err error
	if cuberSvc.combinedLogsSvc.IsLoaded() {
		_, _, _, err = cuberSvc.ParseOutputFromCombinedLog(path.Base(logFilePath))
	} else {
		_, _, _, err = cuberSvc.ParseOutputFromFile(logFilePath)
	}
	return err == nil
}

func (cuberSvc *CuberService) TrackedInvoke(parameters InvokeParameters, control InvokeControl) error {
	if shouldStop_, exists := control.ShouldStop[parameters.Encoding]; exists && shouldStop_ {
		return nil
	}

	cubesFilePath, logFilePath, err := cuberSvc.Invoke(parameters, control)
	if err != nil {
		return err
	}

	if parameters.SkipLogs {
		*control.CubesetPaths = append(*control.CubesetPaths, cubesFilePath)
		return nil
	}

	var (
		processTime      time.Duration
		numCubes         int
		numRefutedLeaves int
	)

	if cuberSvc.combinedLogsSvc.IsLoaded() {
		processTime, numCubes, numRefutedLeaves, err = cuberSvc.ParseOutputFromCombinedLog(path.Base(logFilePath))
	} else {
		processTime, numCubes, numRefutedLeaves, err = cuberSvc.ParseOutputFromFile(logFilePath)
	}
	if err != nil {
		return err
	}

	maxCubesExceeded := parameters.MaxCubes > 0 && numCubes > parameters.MaxCubes
	minCubesExceeded := numCubes < parameters.MinCubes
	minRefutedLeavesViolated := parameters.MinRefutedLeaves > 0 && numRefutedLeaves < parameters.MinRefutedLeaves
	hasViolated := maxCubesExceeded || minRefutedLeavesViolated || minCubesExceeded || numCubes <= 1

	cuberSvc.logSvc.CubeResult(cubesFilePath, processTime, numCubes, numRefutedLeaves, hasViolated)

	if maxCubesExceeded {
		control.ShouldStop[parameters.Encoding] = true
		cuberSvc.commandSvc.StopGroup(control.CommandGroup)
	}

	if hasViolated {
		if err := os.Remove(cubesFilePath); err != nil {
			return err
		}
		return cuber.ErrCubesetViolatedConstraints
	}

	*control.CubesetPaths = append(*control.CubesetPaths, cubesFilePath)

	return nil
}

// TODO: March doesn't like comments, strip 'em!
func (cuberSvc *CuberService) Invoke(parameters InvokeParameters, control InvokeControl) (string, string, error) {
	config := cuberSvc.configSvc.Config

	ctx, cancel := context.WithTimeout(context.Background(), parameters.Timeout)
	defer cancel()

	var thresholdArg string
	if parameters.ThresholdType == cuber.CutoffDepth {
		thresholdArg = "-d"
	} else {
		thresholdArg = "-n"
	}

	cubesFilePath := cuberSvc.CubesFilePath(parameters.Encoding, parameters.ThresholdType, parameters.Threshold, parameters.Suffix)
	cmd := exec.CommandContext(ctx, config.Paths.Bin.March, parameters.Encoding, "-o", cubesFilePath, thresholdArg, strconv.Itoa(parameters.Threshold))
	log.Println(cmd)
	cuberSvc.commandSvc.AddToGroup(control.CommandGroup, cmd)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	logFilePath := ""
	if !parameters.SkipLogs {
		logFilePath = cuberSvc.getLogFilePath(cubesFilePath)
		cuberSvc.filesystemSvc.WriteFromPipe(stdoutPipe, logFilePath)
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Exit code", cmd.ProcessState.ExitCode(), cmd)
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

func (cuberSvc *CuberService) Run(encodings []encoder.Encoding, parameters pipeline.CubeParams) []string {
	err := cuberSvc.filesystemSvc.PrepareDirs([]string{"cubesets", "encodings", "logs"})
	cuberSvc.errorSvc.Fatal(err, "Cuber: failed to prepare the required dirs")

	cubesFilePaths := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	shouldStop := map[string]bool{}
	commandGrps := map[string]*command.Group{}
	log.Println("Cuber: started")

	cuberSvc.Loop(encodings, parameters, func(encoding string, threshold, timeout int) {
		if cuberSvc.ShouldSkip(encoding, cuber.CutoffVars, threshold, "") {
			log.Println("Cuber: skipped", threshold, encoding)
			cubesFilePaths = append(cubesFilePaths, cuberSvc.CubesFilePath(encoding, cuber.CutoffVars, threshold, ""))
			return
		}

		if _, exists := commandGrps[encoding]; !exists {
			commandGrps[encoding] = cuberSvc.commandSvc.CreateGroup()
		}

		pool.Submit(func(encoding string, threshold int) func() {
			return func() {
				// TODO: Simplify
				err := cuberSvc.TrackedInvoke(InvokeParameters{
					Encoding:         encoding,
					ThresholdType:    cuber.CutoffVars,
					Threshold:        threshold,
					Timeout:          time.Duration(timeout) * time.Second,
					MaxCubes:         parameters.MaxCubes,
					MinCubes:         parameters.MinCubes,
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
	log.Println("Cuber: stopped")
	return cubesFilePaths
}
