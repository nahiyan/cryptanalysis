package cnc

import (
	"benchmark/config"
	"benchmark/constants"
	"benchmark/core"
	"benchmark/encodings"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/samber/lo"
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

// Watches workers and the pool and kills the commands if necessary when the pool is stopped
func ManageWorkers(pool *pond.WorkerPool, commands *[]*exec.Cmd) {
	for {
		if !pool.Stopped() && pool.RunningWorkers() > 0 {
			time.Sleep(time.Second * 1)
			continue
		}

		for {
			time.Sleep(1 * time.Second)

			for _, cmd := range *commands {
				if cmd != nil && cmd.Process != nil {
					cmd.Process.Kill()
				}
			}

			if pool.RunningWorkers() == 0 {
				break
			}
		}

		break
	}
}

func FindThreshold(context types.CommandContext) (uint, time.Duration) {
	numWorkers := int(context.FindCncThreshold.NumWorkers)
	type Cubeset struct {
		threshold uint
		cubeCount uint
	}
	cubesets := []Cubeset{}
	lock := sync.Mutex{}

	fmt.Println("Generating cubesets for various thresholds")
	// * 1. Generate the cubes for the selected thresholds
	{
		// Parse the encoding
		cnf, err := encodings.Process(context.FindCncThreshold.InstanceName)
		if err != nil {
			log.Fatal("Failed to find threshold: ", err)
		}

		decrementSize := 10
		instanceFilePath := path.Join(constants.EncodingsDirPath, context.FindCncThreshold.InstanceName+".cnf")
		lookaheadSolverMaxTime := time.Second * 5000
		maxCubeCount := int(context.FindCncThreshold.MaxCubes)
		minRefutedLeaves := int(context.FindCncThreshold.MinRefutedLeaves)
		pool := pond.New(numWorkers, 1000)

		threshold := int(cnf.FreeVariables) - decrementSize
		commands := make([]*exec.Cmd, 0)

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
					commands = append(commands, cmd)
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
					fmt.Printf("Generated cubeset with n = %d, cube count = %d, and refuted leaves = %d\n", threshold, cubeCount, refutedLeaves)

					// Add to the cubeset if it satisfies the constraints
					if cubeCount <= uint(maxCubeCount) && refutedLeaves >= uint(minRefutedLeaves) {
						cubesets = append(cubesets, Cubeset{threshold: uint(threshold), cubeCount: cubeCount})
					} else {
						if err := os.Remove(outputFilePath); err != nil {
							fmt.Println("Failed to remove the rejected cubeset: ", err)
						}
					}

					// Stop the pool if we reach the max cube count, assuming that it'd take longer time to generate more cubes
					if cubeCount > uint(maxCubeCount) {
						fmt.Println("Max cube exceeded; stopping the pool")
						pool.Stop()
					}
				}
			}(threshold, pool))
			threshold -= decrementSize
		}

		ManageWorkers(pool, &commands)
		fmt.Println("Finished generating the cubesets")
	}

	// * 2. Go through the cubesets, starting from the one with possibly the least difficult cubes (most refuted leaves), and solve a subset with random cubes

	// cubesets = []Cubeset{
	// 	{
	// 		threshold: 2260,
	// 		cubeCount: 344,
	// 	},
	// 	{
	// 		threshold: 2250,
	// 		cubeCount: 528,
	// 	},
	// 	{
	// 		threshold: 2240,
	// 		cubeCount: 820,
	// 	},
	// 	{
	// 		threshold: 2230,
	// 		cubeCount: 1263,
	// 	},
	// 	{
	// 		threshold: 2220,
	// 		cubeCount: 1914,
	// 	},
	// 	{
	// 		threshold: 2210,
	// 		cubeCount: 2941,
	// 	},
	// 	{
	// 		threshold: 2200,
	// 		cubeCount: 4497,
	// 	},
	// 	{
	// 		threshold: 2190,
	// 		cubeCount: 6870,
	// 	},
	// }

	type BestResult struct {
		threshold uint
		estimate  time.Duration
	}
	bestResult := BestResult{
		threshold: 0,
		estimate:  time.Duration(math.MaxInt64),
	}

	{
		cdclSolverMaxDuration := time.Duration(5000 * time.Second)
		sampleSize := int(context.FindCncThreshold.SampleSize)
		// solver := constants.Kissat
		timedOut := false
		for i := len(cubesets) - 1; i >= 0; i-- {
			cubeset := cubesets[i]

			// Take a random sample from the cubeset
			cubeRandomSample := utils.RandomCubes(int(cubeset.cubeCount), int(sampleSize))

			// Benchmark the subset of the cubeset
			runtimes := make([]time.Duration, 0)
			pool := pond.New(numWorkers, 1000)
			commands := make([]*exec.Cmd, 0)

			fmt.Printf("Benchmarking sample from cubeset with n = %d\n", cubeset.threshold)

			// TODO: Manage workers for each cubeset in its own task group

			// Add the CDCL commands to the pool
			for _, cube := range cubeRandomSample {
				pool.Submit(
					func(cube int, threshold uint) func() {
						return func() {
							// Generate the sub-problems command
							subproblemCmd := fmt.Sprintf("%s gen-subproblem --instance-name %s --cube-index %d --threshold %d", config.Get().Paths.Bin.Benchmark, context.FindCncThreshold.InstanceName, cube, threshold)

							command := utils.NewCommand().AddCommand(subproblemCmd).AddPipe(utils.PipeVl).AddPlaceholder()

							// Run the CDCL solver
							exitCode, duration := core.KissatWithStream(command, cdclSolverMaxDuration, &commands)
							if exitCode == 10 || exitCode == 20 {
								fmt.Printf("Tried CDCL on cube index = %d, n = %d with exit code = %d\n", cube, threshold, exitCode)
								// Add the runtime
								lock.Lock()
								runtimes = append(runtimes, duration)
								lock.Unlock()
							} else {
								fmt.Println("Unexpected exit code: ", exitCode)
							}

							// Discard the pool if the solver times out
							if duration.Seconds() > cdclSolverMaxDuration.Seconds() {
								fmt.Printf("Timed out for n = %d, stopping pool\n", threshold)
								pool.Stop()
								timedOut = true
							}
						}
					}(cube, cubeset.threshold))
			}

			ManageWorkers(pool, &commands)

			if timedOut {
				break
			}

			// Note: Oleg finds estimate for the rest of the cubeset, but here we're finding for all
			// Calculate the estimate
			estimateInSeconds := (lo.SumBy(runtimes, func(duration time.Duration) float64 {
				return duration.Seconds()
			}) / (float64(numWorkers * sampleSize))) * float64(cubeset.cubeCount)

			// See if we can register it as the best estimate
			if estimateInSeconds < bestResult.estimate.Seconds() {
				bestResult.estimate = time.Duration(time.Second * time.Duration(estimateInSeconds))
				bestResult.threshold = cubeset.threshold
			}

			fmt.Printf("Found estimate of %s for n = %d and cube count of %d\n\n", bestResult.estimate.String(), bestResult.threshold, cubeset.cubeCount)
		}
	}

	return bestResult.threshold, bestResult.estimate
}
