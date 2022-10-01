package main

import (
	"fmt"
	"os"
	"os/exec"
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
)

type Context struct {
	completionTimes map[string][]*time.Duration
}

func cryptoMiniSat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	// cmd := exec.Command("cryptominisat5", "--verb 0", filepath)
	cmd := exec.Command("cryptominisat5", filepath)
	// fmt.Println(cmd.String())
	_, err := cmd.Output()
	if err != nil && err.Error() != "exit status 10" && err.Error() != "exit status 20" {
		fmt.Println("Failed to run the command: ", err.Error())
	}

	// fmt.Printf("Solved instance of index %d\n", instanceIndex)

	// fmt.Println(string(output))

	// TODO: Gather the output

	duration := time.Now().Sub(startTime)
	context.completionTimes[CRYPTOMINISAT][instanceIndex] = &duration
}

func kissat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	cmd := exec.Command("kissat", "-q", filepath)
	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to run the command: ", err.Error())
	}

	fmt.Println(string(output))

	// TODO: Gather the output

	duration := time.Now().Sub(startTime)
	context.completionTimes[KISSAT][instanceIndex] = &duration
}

func mapleSat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	cmd := exec.Command("maplesat", "-verb=0", filepath)
	fmt.Println(cmd.String())
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to run the command: ", err.Error())
	}

	// fmt.Println(string(output))

	// TODO: Gather the output

	duration := time.Now().Sub(startTime)
	context.completionTimes[MAPLESAT][instanceIndex] = &duration
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
	satSolvers := []string{"cryptominisat"}

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
						filepath := fmt.Sprintf("../../encodings/saeed/md4_%d_%s_xor%d_%s.cnf",
							steps, adderType, xorOption, hash)

						switch satSolver {
						case CRYPTOMINISAT:
							go cryptoMiniSat(filepath, context, i, time.Now())
						case KISSAT:
							go kissat(filepath, context, i, time.Now())
						case MAPLESAT:
							go mapleSat(filepath, context, i, time.Now())
						}

						i++
					}
				}
			}
		}

		fmt.Printf("Spawned %d the instances of %s.\n", instancesCount, satSolver)

		{
			lastOutputTime := time.Now().Add(-time.Second * 1)
			// Wait as long as the operation didn't timeout and the instances aren't completed
			for time.Now().Sub(startTime).Seconds() <= 5000 && !areInstancesCompleted(context, satSolver) {
				if time.Now().Sub(lastOutputTime).Seconds() > 1 {
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
					logMessage := fmt.Sprintf("Time: %.2fs, Instances Solved: %d, Completion: %.2f%%", time.Now().Sub(startTime).Seconds(), completedInstancesCount(context, satSolver), float32(completions)/float32(totalItems)*100)
					fmt.Println(logMessage)
					fmt.Println()
					fmt.Println()

					lastOutputTime = time.Now()

					// Log down to a file
					appendLog(logMessage)
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
