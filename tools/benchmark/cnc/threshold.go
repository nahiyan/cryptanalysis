package cnc

import (
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

// Watches for stop signal and kills the command if such a signal is detected
func WatchForStopSignal(stopSignal chan struct{}, cmd *exec.Cmd, started *bool, killed *bool) {
	<-stopSignal
	for !*started || cmd.Process == nil {
		time.Sleep(time.Second)
	}

	cmd.Process.Kill()
	*killed = true
}

func FindThreshold(context types.CommandContext) (uint, time.Duration) {
	type Cubeset struct {
		threshold uint
		cubeCount uint
	}
	cubesets := []Cubeset{}
	lock := sync.Mutex{}

	// * 1. Generate the cubes for the selected thresholds
	{
		numWorkersLookahead := int(context.FindCncThreshold.NumWorkersLookahead)
		fmt.Println("Generating cubesets for various thresholds")
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
		pool := pond.New(numWorkersLookahead, 1000, pond.IdleTimeout(100*time.Millisecond))

		threshold := int(cnf.FreeVariables) - decrementSize

		// Generate the thresholds
		thresholds := []int{}
		for threshold >= 0 {
			if threshold%10 != 0 {
				threshold--
				continue
			}

			thresholds = append(thresholds, threshold)
			threshold -= decrementSize

		}

		// Generate the stop signal channels
		channels := make([]chan struct{}, 0)
		for i := 0; i < len(thresholds); i++ {
			channels = append(channels, make(chan struct{}))
		}

		// Create and submit the tasks
		for i, threshold := range thresholds {
			// Feed the instance to the lookahead solver
			pool.Submit(func(threshold int, stopSignal chan struct{}) func() {
				return func() {
					// TODO: Add resume capability
					// if utils.FileExists(instanceFilePath) {
					// 	return
					// }

					// Start the stop signal watcher
					cmd := new(exec.Cmd)
					var started, killed = false, false
					// go WatchForStopSignal(stopSignal, cmd, &started, &killed)
					go func() {
						<-stopSignal
						for !started || cmd.Process == nil {
							time.Sleep(time.Second)
						}

						cmd.Process.Kill()
						killed = true
					}()

					outputFilePath := fmt.Sprintf("%sn%d_%s.icnf", constants.EncodingsDirPath, threshold, context.FindCncThreshold.InstanceName)
					cmd, cancel := core.MarchCmd(instanceFilePath, outputFilePath, uint(threshold), lookaheadSolverMaxTime)
					defer cancel()

					started = true
					startTime := time.Now()
					output_, err := cmd.Output()
					if err != nil && !killed {
						// Any signal of termination indicates that the process was killed
						if strings.HasPrefix(err.Error(), "signal: ") {
							return
						}
						fmt.Println(err)
					}
					if killed {
						return
					}
					output := string(output_)
					duration := time.Since(startTime)

					cubeCount, refutedLeaves := ProcessMarchLog(output)
					fmt.Printf("Generated cubeset with n = %d, cube count = %d, and refuted leaves = %d in %s\n", threshold, cubeCount, refutedLeaves, duration)

					// Add to the cubeset if it satisfies the constraints
					if cubeCount <= uint(maxCubeCount) && refutedLeaves >= uint(minRefutedLeaves) {
						lock.Lock()
						cubesets = append(cubesets, Cubeset{threshold: uint(threshold), cubeCount: cubeCount})
						lock.Unlock()
					} else {
						if err := os.Remove(outputFilePath); err != nil && !killed {
							fmt.Println("Failed to remove the rejected cubeset: ", err)
						}
					}

					// Stop acceptng new jobs in the pool if we reach the max cube count, while also sending stop signals to the lower threshold jobs, since lower threshold means more cubes
					if cubeCount > uint(maxCubeCount) {
						// Stop the pool
						fmt.Println("Max cube exceeded; stopping the pool")
						pool.Stop()

						// Send stop signal to all workers that will generate more cubes than the limit
						_, channelIndex, _ := lo.FindIndexOf(channels, func(c chan struct{}) bool {
							return c == stopSignal
						})
						for i, channel := range channels {
							if i > channelIndex {
								channel <- struct{}{}
							}
						}
					}

					// Stop the child goroutine
					stopSignal <- struct{}{}
				}
			}(threshold, channels[i]))
		}

		for pool.RunningWorkers() > 0 {
			time.Sleep(time.Second)
		}

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

	// cubesets = []Cubeset{
	// 	{
	// 		threshold: 3390,
	// 		cubeCount: 31200,
	// 	},
	// 	{
	// 		threshold: 3380,
	// 		cubeCount: 46421,
	// 	},
	// 	{
	// 		threshold: 3370,
	// 		cubeCount: 68751,
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

	if len(cubesets) == 0 {
		fmt.Println("No cubeset found eligible for the benchmark")
	}

	{
		cdclSolverMaxDuration := time.Duration(context.FindCncThreshold.CdclTimeout * uint(time.Second))
		sampleSize := int(context.FindCncThreshold.SampleSize)
		solver := constants.Kissat
		timedout := false
		stopOnTimeout := false
		numWorkersCdcl := int(context.FindCncThreshold.NumWorkersCdcl)
		for i := len(cubesets) - 1; i >= 0; i-- {
			cubeset := cubesets[i]

			// Take a random sample from the cubeset
			cubeRandomSample := utils.RandomCubes(int(cubeset.cubeCount), int(sampleSize))

			// Prepare the benchmark data structure and the pool
			type Result struct {
				Runtime  time.Duration
				ExitCode uint
			}
			results := make([]Result, 0)
			pool := pond.New(numWorkersCdcl, sampleSize, pond.IdleTimeout(10*time.Millisecond))

			// Create the stop signal channels
			channels := make([]chan struct{}, 0)
			for i := 0; i < sampleSize; i++ {
				channels = append(channels, make(chan struct{}))
			}

			fmt.Printf("Benchmarking sample from cubeset with n = %d\n", cubeset.threshold)

			// TODO: Manage workers for each cubeset in its own task group

			// Add the CDCL commands to the pool
			for i, cube := range cubeRandomSample {
				pool.Submit(
					func(cube int, threshold uint, stopSignal chan struct{}) func() {
						return func() {
							// Generate the sub-problems command
							subproblem, err := encodings.GenerateSubProblemAsStringWithThreshold(context.FindCncThreshold.InstanceName, cube, &threshold)
							if err != nil {
								fmt.Println("Failed to generate subproblem", err)
								return
							}
							subproblemReader := strings.NewReader(subproblem)

							// Spawn a watcher for the stop signal
							var cmd *exec.Cmd
							var started, killed = false, false
							go func() {
								<-stopSignal
								for !started || cmd == nil || cmd.Process == nil {
									time.Sleep(time.Second)
								}

								cmd.Process.Kill()
								killed = true
							}()

							// Run the CDCL solver
							exitCode, duration := core.RunSatSolver(subproblemReader, cdclSolverMaxDuration, solver, types.SatSolverConfig[string]{}, func(c *exec.Cmd) {
								cmd = c
							})

							if killed {
								return
							}

							// Add runtime if the CDCL resulted in SAT or UNSAT
							message := fmt.Sprintf("CDCL on cube index = %d, n = %d with exit code = %d", cube, threshold, exitCode)
							if exitCode == 10 || exitCode == 20 {
								fmt.Printf(message+" took %s\n", duration.String())

								// Add the runtime
								lock.Lock()
								results = append(results, Result{
									Runtime:  duration,
									ExitCode: uint(exitCode),
								})
								lock.Unlock()
							} else if duration.Seconds() > cdclSolverMaxDuration.Seconds() {
								fmt.Println(message + " timed out")
							} else {
								// fmt.Println(cmd.String())
								fmt.Println(message + ", which is unexpected")
							}

							// Discard the pool if the solver times out
							if stopOnTimeout && duration.Seconds() > cdclSolverMaxDuration.Seconds() {
								lock.Lock()
								fmt.Printf("Timed out for n = %d, stopping pool\n", threshold)
								pool.Stop()
								timedout = true

								// Stop all the workers
								for _, channel := range channels {
									channel <- struct{}{}
								}
								lock.Lock()

								return
							}

							stopSignal <- struct{}{}
						}
					}(cube, cubeset.threshold, channels[i]))
			}

			// TODO: The pool should be reused instead of waiting for all the workers to die
			for pool.RunningWorkers() > 0 {
				time.Sleep(time.Second)
			}

			if stopOnTimeout && timedout {
				continue
			}

			totalDuration := lo.SumBy(results, func(item Result) time.Duration {
				return item.Runtime
			})
			averageDuration := totalDuration / time.Duration(numWorkersCdcl*sampleSize)

			// Note: Oleg finds the estimate for the rest of the cubeset, but here we're finding for all
			estimatedTotalDuration := averageDuration * time.Duration(cubeset.cubeCount)

			// See if we can register it as the best estimate
			if estimatedTotalDuration < bestResult.estimate {
				bestResult.estimate = estimatedTotalDuration
				bestResult.threshold = cubeset.threshold
			}

			totalSat := lo.SumBy(results, func(item Result) uint {
				if item.ExitCode == 10 {
					return 1
				}
				return 0
			})
			totalUnsat := lo.SumBy(results, func(item Result) uint {
				if item.ExitCode == 20 {
					return 1
				}

				return 0
			})

			fmt.Printf("Found estimate of %s for n = %d and cube count of %d with %d SAT and %d UNSAT solutions\n\n", estimatedTotalDuration, cubeset.threshold, cubeset.cubeCount, totalSat, totalUnsat)
		}
	}

	return bestResult.threshold, bestResult.estimate
}
