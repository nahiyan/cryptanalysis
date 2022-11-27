package cnc

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/encodings"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/alitto/pond"
)

func ProcessMarchLog(log string) (uint, uint) {
	lines := strings.Split(log, "\n")
	var cubeCount, refutedLeaves uint
	for _, line := range lines {
		if strings.HasPrefix(line, "c number of cubes") {
			fmt.Sscanf(line, "c number of cubes %d, including %d refuted leaves", &cubeCount, &refutedLeaves)
		}
	}

	return cubeCount, refutedLeaves
}

func FindThreshold(context types.CommandContext) (uint, time.Duration) {
	// Parse the encoding
	cnf, err := encodings.Process(context.FindCncThreshold.InstanceName)
	if err != nil {
		log.Fatal("Failed to find threshold: ", err)
	}

	numWorkers := 16
	type Cubeset struct {
		threshold uint
		cubeCount uint
	}
	cubesets := []Cubeset{}

	// * 1. Generate the cubes for the selected thresholds
	{
		decrementSize := 10
		instanceFilePath := path.Join(constants.EncodingsDirPath, context.FindCncThreshold.InstanceName+".cnf")
		lookaheadSolverMaxTime := time.Second * 5000
		// maxCubeCount := 100000
		maxCubeCount := 100
		minRefutedLeaves := 500
		pool := pond.New(numWorkers, 1000)

		threshold := int(cnf.FreeVariables) - decrementSize
		commands := make(map[int]*exec.Cmd)
		lock := sync.Mutex{}

		for threshold > 0 {
			if threshold%10 != 0 {
				threshold--
				continue
			}

			// Feed the instance to the lookahead solver
			pool.Submit(func(threshold int, pool *pond.WorkerPool) func() {
				return func() {
					// TODO: Add resume capability
					// if utils.FileExists(instanceFilePath) {
					// 	return
					// }

					outputFilePath := fmt.Sprintf("%sn%d_%s.icnf", constants.EncodingsDirPath, threshold, context.FindCncThreshold.InstanceName)
					cmd, cancel := core.MarchCmd(instanceFilePath, outputFilePath, uint(threshold), lookaheadSolverMaxTime)
					defer cancel()

					lock.Lock()
					commands[threshold] = cmd
					lock.Unlock()

					output_, err := cmd.Output()
					if err != nil {
						if err.Error() == "signal: killed" {
							return
						}
						fmt.Printf("Failed to run March for n = %d\n", threshold)
					}
					output := string(output_)

					cubeCount, refutedLeaves := ProcessMarchLog(output)
					fmt.Println(threshold, cubeCount, refutedLeaves)

					// Add to the cubeset if it follows a criteria
					if cubeCount <= uint(maxCubeCount) && refutedLeaves >= uint(minRefutedLeaves) {
						cubesets = append(cubesets, Cubeset{threshold: uint(threshold), cubeCount: cubeCount})
					} else {
						if err := os.Remove(outputFilePath); err != nil {
							fmt.Println("Failed to remove the rejected cubeset: ", err)
						}
					}

					// Stop the pool if we reach the max cube count, assuming that it'd make longer time to generate more cubes
					if cubeCount > uint(maxCubeCount) {
						fmt.Println("Max cube exceeded; stopping the pool")
						pool.Stop()
					}
				}
			}(threshold, pool))
			threshold -= decrementSize
		}

		for {
			if !pool.Stopped() {
				time.Sleep(time.Second * 1)
				continue
			}

			for _, cmd := range commands {
				if cmd != nil && cmd.Process != nil {
					cmd.Process.Kill()
				}
			}

			break
		}
		fmt.Println("All done")
	}

	// * 2. Go through the cubesets, starting from the one with possibly the least difficult cubes (most refuted leaves), and solve a subset with random cubes
	// TODO: Calculate best threshold
	// bestCubeset := cubesets[0]
	// TODO: Calculate best estimate
	// bestEstimate := time.Duration(0)
	{
		sampleSize := 100
		// solver := constants.Kissat
		for i := len(cubesets) - 1; i >= 0; i-- {
			cubeset := cubesets[i]

			// Take random sample from the cubeset
			cubeRandomSample := utils.RandomCubes(int(cubeset.cubeCount), sampleSize)

			// TODO: Get the runtimes from the solvers
			pool := pond.New(numWorkers, 1000)
			for _, cube := range cubeRandomSample {
				pool.Submit(
					func(cube int) func() {
						return func() {
							// Run the CDCL solver

							// TODO: Discard the pool if the solver times out
						}
					}(cube))
			}
			pool.StopAndWait()
		}
	}

	return uint(bestThreshold), time.Duration(0)
}
