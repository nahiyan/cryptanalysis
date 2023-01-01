package services

import (
	"benchmark/internal/consts"
	"benchmark/internal/cubeset"
	"benchmark/internal/pipeline"
	"context"
	"fmt"
	"math"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/alitto/pond"
)

func (cuberSvc *CuberService) CubesFilePath(encoding string) string {
	encodingDir, fileName := path.Split(encoding)
	cubesFilePath := path.Join(encodingDir, fileName+".cubes")

	return cubesFilePath
}

func (cuberSvc *CuberService) ShouldCube(encoding string) bool {
	filesystemSvc := cuberSvc.filesystemSvc
	cubesFilePath := cuberSvc.CubesFilePath(encoding)

	return !filesystemSvc.FileExists(cubesFilePath)
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

	cubesFilePath := cuberSvc.CubesFilePath(encoding)
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

	cubesFilePath := cuberSvc.CubesFilePath(encoding)
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

func (cuberSvc *CuberService) RunRegular(encodings []string, parameters pipeline.Cubing) []string {
	cubesFilePaths := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	fmt.Println("Cuber: started")
	for _, encoding := range encodings {
		if !cuberSvc.ShouldCube(encoding) {
			fmt.Println("Cuber: skipped", encoding)
			continue
		}

		thresholds := parameters.Thresholds
		if len(thresholds) == 0 {
			freeVariables, _, err := cuberSvc.encodingSvc.Process(encoding)
			cuberSvc.errorSvc.Fatal(err, "Cuber: failed to process the encoding")

			var stepSize float64 = 10
			// starting threshold is the nearest multiple of step size less than the num. of free variables
			threshold := math.Floor(float64(freeVariables)/stepSize) * stepSize

			for {
				thresholds = append(thresholds, int(threshold))
				threshold -= stepSize

				if threshold < 0 {
					break
				}
			}
		}

		for _, threshold := range thresholds {
			pool.Submit(func() {
				err := cuberSvc.TrackedInvoke(encoding, threshold, time.Duration(parameters.Timeout)*time.Second)
				cuberSvc.errorSvc.Fatal(err, "Cuber: failed to cube")
			})
		}
	}
	pool.StopAndWait()
	fmt.Println("Cuber: stopped")

	return cubesFilePaths
}

func (cuberSvc *CuberService) Run(encodings []string, parameters pipeline.Cubing) []string {
	var cubesets []string

	switch parameters.Platform {
	case consts.General:
		cubesets = cuberSvc.RunRegular(encodings, parameters)
	}

	return cubesets
}
