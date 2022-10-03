package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"
)

const (
	CRYPTOMINISAT           = "cryptominisat"
	KISSAT                  = "kissat"
	CADICAL                 = "cadical"
	GLUCOSE                 = "glucose"
	MAPLESAT                = "maplesat"
	MAX_TIME                = 5000
	BENCHMARK_LOG_FILE_NAME = "benchmark.log"
	BASE_PATH               = "../../"
	SOLUTIONS_DIR_PATH      = BASE_PATH + "solutions/saeed/"
	ENCODINGS_DIR_PATH      = BASE_PATH + "encodings/saeed/"
)

type Context struct {
	completionTimes map[string][]*time.Duration
}

func invokeSatSolver(command string, satSolver string, context_ *Context, filepath string, startTime time.Time, instanceIndex uint) {
	cmd := exec.Command("timeout", strconv.Itoa(MAX_TIME), "bash", "-c", command)
	if err := cmd.Run(); err != nil {
		// TODO: Aggregate the logs
		exiterr, _ := err.(*exec.ExitError)
		exitcode := exiterr.ExitCode()
		if exitcode != 10 && exitcode != 20 && exitcode != 124 {
			// TODO: Take action
			fmt.Println("Error with solving the instance:", exitcode)
		}
	}

	duration := time.Since(startTime)
	if duration <= time.Second*MAX_TIME {
		// TODO: Validate the results
		context_.completionTimes[satSolver][instanceIndex] = &duration

		// Log down to a file
		logMessage := fmt.Sprintf("Time: %.2fs, Instance index: %d", duration.Seconds(), instanceIndex)
		appendLog(logMessage)
	}
}

func cryptoMiniSat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("cryptominisat5 --verb 0 %s > %scryptominisat/%ssol", filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, CRYPTOMINISAT, context, filepath, startTime, instanceIndex)
}

func kissat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("kissat -q %s > %skissat/%ssol", filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, KISSAT, context, filepath, startTime, instanceIndex)
}

func cadical(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("cadical -q %s > %scadical/%ssol", filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, CADICAL, context, filepath, startTime, instanceIndex)
}

func mapleSat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("maplesat -verb=0 %s %smaplesat/%ssol", filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, MAPLESAT, context, filepath, startTime, instanceIndex)
}

func glucose(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("glucose -verb=0 %s %sglucose/%ssol", filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, GLUCOSE, context, filepath, startTime, instanceIndex)
}

func areInstancesCompleted(context *Context, satSolver string) bool {
	return completedInstancesCount(context, satSolver) == uint(len(context.completionTimes[satSolver]))
}

func completedInstancesCount(context *Context, satSolver string) uint {
	var count uint = 0

	for _, duration := range context.completionTimes[satSolver] {
		if duration != nil {
			count++
		}
	}

	return count
}

func appendLog(message string) {
	f, err := os.OpenFile(BENCHMARK_LOG_FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to write logs")
	}
	_, err = f.WriteString(message + "\n")
	if err != nil {
		panic("Failed to write logs")
	}
	f.Close()
}

func main() {
	// Variations
	xorOptions := []uint{0, 1}
	hashes := []string{"ffffffffffffffffffffffffffffffff",
		"00000000000000000000000000000000"}
	adderTypes := []string{"counter_chain", "dot_matrix"}
	stepVariations := makeRange(16, 32)

	// sat_solvers = ["cryptominisat", "kissat", "cadical", "glucose", "maplesat"]
	satSolvers := []string{CRYPTOMINISAT, KISSAT, CADICAL, GLUCOSE, MAPLESAT}

	// Should be 264 for all the possible variations
	instancesCount := len(xorOptions) * len(hashes) * len(adderTypes) * len(stepVariations)

	// Define the context
	context := &Context{
		completionTimes: make(map[string][]*time.Duration),
	}
	for _, satSolver := range satSolvers {
		context.completionTimes[satSolver] = make([]*time.Duration, instancesCount)
	}

	os.Remove(BENCHMARK_LOG_FILE_NAME)

	// Solve the encodings for each SAT solver
	for _, satSolver := range satSolvers {
		var i uint = 0

		startTime := time.Now()

		appendLog("SAT Solver: " + satSolver)
		for _, steps := range stepVariations {
			for _, hash := range hashes {
				for _, xorOption := range xorOptions {
					for _, adderType := range adderTypes {
						filepath := fmt.Sprintf("%smd4_%d_%s_xor%d_%s.cnf",
							ENCODINGS_DIR_PATH, steps, adderType, xorOption, hash)

						switch satSolver {
						case CRYPTOMINISAT:
							go cryptoMiniSat(filepath, context, i, time.Now())
						case KISSAT:
							go kissat(filepath, context, i, time.Now())
						case CADICAL:
							go cadical(filepath, context, i, time.Now())
						case MAPLESAT:
							go mapleSat(filepath, context, i, time.Now())
						case GLUCOSE:
							go glucose(filepath, context, i, time.Now())
						}

						i++
					}
				}
			}
		}

		fmt.Printf("Spawned %d instances of %s.\n", instancesCount, satSolver)

		{
			interval := time.Second * 5
			lastOutputTime := time.Now().Add(-interval)
			// Wait as long as the operation didn't timeout and the instances aren't completed
			for time.Since(startTime).Seconds() <= MAX_TIME && !areInstancesCompleted(context, satSolver) {
				if time.Since(lastOutputTime) > interval {
					totalItems := len(context.completionTimes[satSolver])
					completions := 0
					for i, item := range context.completionTimes[satSolver] {
						if item != nil {
							fmt.Print("x")
							completions++
						} else {
							fmt.Print("-")
						}
						if (i+1)%32 == 0 {
							fmt.Println()
						}
					}
					fmt.Println()
					fmt.Printf("Completion: %.2f%%\n", float32(completions)/float32(totalItems)*100)
					// logMessage := fmt.Sprintf("Time: %.2fs, Instances Solved: %d, Completion: %.2f%%", time.Now().Sub(startTime).Seconds(), completedInstancesCount(context, satSolver), float32(completions)/float32(totalItems)*100)
					// fmt.Println(logMessage)
					// fmt.Println()
					fmt.Println()

					lastOutputTime = time.Now()

					// Log down to a file
					// appendLog(logMessage)
				}
			}
		}

		fmt.Printf("Results for %s:\n", satSolver)
		for _, item := range context.completionTimes[satSolver] {
			if item != nil {
				fmt.Printf("%.2f ", item.Seconds())
			} else {
				fmt.Print("- ")
			}
		}

		fmt.Println()
	}
}
