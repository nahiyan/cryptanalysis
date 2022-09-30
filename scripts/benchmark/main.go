package main

import (
	"fmt"
	"os/exec"
	"time"
)

const (
	CRYPTOMINISAT = "cryptominisat"
	KISSAT        = "kissat"
	CADICAL       = "cadical"
	GLUCOSE       = "glucose"
	MAPLESAT      = "maplesat"
	MAX_TIME      = 5000
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
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to run the command: ", err.Error())
	}

	fmt.Println(string(output))

	// TODO: Gather the output

	duration := time.Now().Sub(startTime)
	context.completionTimes[MAPLESAT][instanceIndex] = &duration
}

func areInstancesSolved(context *Context, satSolver string) bool {
	for _, duration := range context.completionTimes[satSolver] {
		if duration == nil {
			return false
		}
	}

	return false
}

func main() {

	// Variations
	xorOptions := []uint{0, 1}
	hashes := []string{"ffffffffffffffffffffffffffffffff",
		"00000000000000000000000000000000"}
	adderTypes := []string{"counter_chain", "dot_matrix"}
	stepVariations := makeRange(16, 48)

	// sat_solvers = ["cryptominisat", "kissat", "cadical", "glucose", "maplesat"]
	satSolvers := []string{"cryptominisat"}

	// Should be 264
	instancesCount := len(xorOptions) * len(hashes) * len(adderTypes) * len(stepVariations)

	// Define the context
	context := &Context{
		completionTimes: make(map[string][]*time.Duration),
	}
	for _, satSolver := range satSolvers {
		context.completionTimes[satSolver] = make([]*time.Duration, instancesCount)
	}

	// Solve the encodings for each SAT solver
	for _, satSolver := range satSolvers {
		var i uint = 0

		startTime := time.Now()

		for _, steps := range stepVariations {
			for _, hash := range hashes {
				for _, xorOption := range xorOptions {
					for _, adderType := range adderTypes {
						filepath := fmt.Sprintf("encodings/saeed/md4_%d_%s_xor%d_%s.cnf",
							steps, adderType, xorOption, hash)

						switch satSolver {
						case CRYPTOMINISAT:
							go cryptoMiniSat(filepath, context, i, time.Now())
						case KISSAT:
							go kissat(filepath, context, i, time.Now())
						case MAPLESAT:
							go mapleSat(filepath, context, i, time.Now())
						}

						// time.Sleep(time.Second * 5)

						i++

					}
				}
			}
		}

		fmt.Printf("Spawned all the instances of %s\n", satSolver)

		{
			lastOutputTime := time.Now()
			// Wait for 5000 seconds or for all the instances to be solved, whichever happens first
			for time.Now().Sub(startTime).Seconds() <= 5000 || !areInstancesSolved(context, satSolver) {
				if time.Now().Sub(lastOutputTime).Seconds() > 5 {
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
					fmt.Printf("Completion: %.2f%%", float32(completions)/float32(totalItems)*100)
					fmt.Println()
					fmt.Println()

					lastOutputTime = time.Now()
				}

			}
		}

		fmt.Printf("Results for %s\n", satSolver)
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
