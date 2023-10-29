package services

import (
	"bufio"
	"context"
	"cryptanalysis/internal/command"
	"cryptanalysis/internal/cuber"
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/pipeline"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alitto/pond"
)

// TODO: Go through the entire code of this service and rewrite if necessary (since there are lots of race conditions)
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
	// TODO: Add it back based on a parameter (removeCubesets)
	minCubesExceeded := numCubes < parameters.MinCubes
	minRefutedLeavesViolated := parameters.MinRefutedLeaves > 0 && numRefutedLeaves < parameters.MinRefutedLeaves
	hasViolated := maxCubesExceeded || minRefutedLeavesViolated || minCubesExceeded || numCubes <= 1

	cuberSvc.logSvc.CubeResult(cubesFilePath, processTime, numCubes, numRefutedLeaves, hasViolated)

	if maxCubesExceeded {
		control.ShouldStop[parameters.Encoding] = true
		cuberSvc.commandSvc.StopGroup(control.CommandGroup)
	}

	// Stop if there's no upper bound and min. cubes has reached
	if numCubes >= parameters.MinCubes && parameters.MaxCubes == 0 {
		control.ShouldStop[parameters.Encoding] = true
		cuberSvc.commandSvc.StopGroup(control.CommandGroup)
	}

	if hasViolated {
		// TODO: Take the decision based on a parameter (removeCubesets)
		if err := os.Remove(cubesFilePath); err != nil {
			return err
		}
		return cuber.ErrCubesetViolatedConstraints
	}

	*control.CubesetPaths = append(*control.CubesetPaths, cubesFilePath)

	return nil
}

// TODO: March doesn't like leading comments, strip 'em!
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
	cmd := exec.CommandContext(ctx, config.Paths.Bin.March, parameters.Encoding, thresholdArg, strconv.Itoa(parameters.Threshold), "-o", cubesFilePath)
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
	config := cuberSvc.configSvc.Config
	for _, encoding := range encodings {
		// Get the thresholds list
		info, err := cuberSvc.encoderSvc.ProcessInstanceName(encoding.GetName())
		cuberSvc.errorSvc.Fatal(err, "Cuber: failed to process instance name")
		var thresholds []int
		thresholds, exists := parameters.ThresholdsMap[info.Steps]
		if !exists {
			thresholds = parameters.Thresholds
		}

		if len(thresholds) == 0 {
			// Command to call March
			cmd := exec.Command(config.Paths.Bin.March, encoding.GetName()+".cnf")
			stdoutPipe, err := cmd.StdoutPipe()
			cuberSvc.errorSvc.Fatal(err, "failed to get the stdout pipe")

			// Execute the command
			cuberSvc.errorSvc.Fatal(cmd.Start(), "failed to start the program")

			// Read March's output for free variables
			freeVariables := -1
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				line := scanner.Text()

				if strings.Contains(line, "c number of free variables") {
					fields := strings.Fields(line)
					freeVariables, err = strconv.Atoi(fields[len(fields)-1])
					cuberSvc.errorSvc.Fatal(err, "failed to get the number of free variables from March")
					cuberSvc.errorSvc.Fatal(cmd.Process.Kill(), "failed to kill the process")
					break
				}
			}

			if freeVariables == -1 {
				log.Fatal("failed to get the free variables using March")
			}

			stepChange := parameters.StepChange
			if stepChange == 0 {
				stepChange = 10
			}
			// Initial threshold is the nearest multiple of step size less than the num. of free variables
			threshold := int(math.Floor(float64(freeVariables)/float64(stepChange)) * float64(stepChange))
			if parameters.InitialThreshold != 0 {
				threshold = parameters.InitialThreshold
			}

			// Try 100 different thresholds
			for i := 0; i < 100; i++ {
				thresholds = append(thresholds, int(threshold))
				threshold -= stepChange

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
