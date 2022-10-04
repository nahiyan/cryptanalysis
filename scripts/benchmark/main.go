package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

const (
	CRYPTOMINISAT           = "cryptominisat"
	KISSAT                  = "kissat"
	CADICAL                 = "cadical"
	GLUCOSE                 = "glucose"
	MAPLESAT                = "maplesat"
	CRYPTOMINISAT_BIN_PATH  = "../../../sat-solvers/cryptominisat"
	KISSAT_BIN_PATH         = "../../../sat-solvers/kissat"
	CADICAL_BIN_PATH        = "../../../sat-solvers/cadical"
	GLUCOSE_BIN_PATH        = "../../../sat-solvers/glucose"
	MAPLESAT_BIN_PATH       = "../../../sat-solvers/maplesat"
	VERIFIER_BIN_PATH       = "../../encoders/saeed/crypto/verify-md4"
	MAX_TIME                = 5000
	BENCHMARK_LOG_FILE_NAME = "benchmark.log"
	BASE_PATH               = "../../"
	SOLUTIONS_DIR_PATH      = BASE_PATH + "solutions/saeed/"
	ENCODINGS_DIR_PATH      = BASE_PATH + "encodings/saeed/"
)

type Context struct {
	progress map[string][]bool
}

func invokeSatSolver(command string, satSolver string, context_ *Context, filepath string, startTime time.Time, instanceIndex uint) {
	cmd := exec.Command("timeout", strconv.Itoa(MAX_TIME), "bash", "-c", command)
	exitCode := 0
	if err := cmd.Run(); err != nil {
		// TODO: Aggregate the logs
		exiterr, _ := err.(*exec.ExitError)
		exitCode = exiterr.ExitCode()
	}

	duration := time.Since(startTime)
	// TODO: Validate the results
	context_.progress[satSolver][instanceIndex] = true

	// Log down to a file
	instanceName := strings.TrimSuffix(path.Base(filepath), ".cnf")
	logMessage := fmt.Sprintf("Time: %.2fs, instance index: %d, instance name: %s, SAT solver: %s, exit code: %d", duration.Seconds(), instanceIndex, instanceName, satSolver, exitCode)
	appendLog(logMessage)

	fmt.Printf("%s completed %s in %.2fs with exit code: %d\n", satSolver, instanceName, duration.Seconds(), exitCode)
}

func cryptoMiniSat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("%s --verb=0 %s > %scryptominisat/%ssol", CRYPTOMINISAT_BIN_PATH, filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, CRYPTOMINISAT, context, filepath, startTime, instanceIndex)
}

func kissat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("%s -q %s > %skissat/%ssol", KISSAT_BIN_PATH, filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, KISSAT, context, filepath, startTime, instanceIndex)
}

func cadical(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("%s -q %s > %scadical/%ssol", CADICAL_BIN_PATH, filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, CADICAL, context, filepath, startTime, instanceIndex)
}

func mapleSat(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("%s -verb=0 %s %smaplesat/%ssol", MAPLESAT_BIN_PATH, filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, MAPLESAT, context, filepath, startTime, instanceIndex)
}

func glucose(filepath string, context *Context, instanceIndex uint, startTime time.Time) {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]
	command := fmt.Sprintf("%s -verb=0 %s %sglucose/%ssol", GLUCOSE_BIN_PATH, filepath, SOLUTIONS_DIR_PATH, solutionFilePath)

	invokeSatSolver(command, GLUCOSE, context, filepath, startTime, instanceIndex)
}

func areAllInstancesCompleted(context *Context) bool {
	for _, progressEntries := range context.progress {
		if lo.SomeBy(progressEntries, func(done bool) bool {
			return !done
		}) {
			return false
		}
	}

	return true
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
	xorOptions := []uint{0}
	hashes := []string{"ffffffffffffffffffffffffffffffff",
		"00000000000000000000000000000000"}
	adderTypes := []string{"counter_chain", "dot_matrix"}
	stepVariations := makeRange(16, 32)

	satSolvers := []string{CRYPTOMINISAT, KISSAT, CADICAL, GLUCOSE, MAPLESAT}

	// Should be 264 for all the possible variations
	instancesCount := len(xorOptions) * len(hashes) * len(adderTypes) * len(stepVariations)

	// Define the context
	context := &Context{
		progress: make(map[string][]bool),
	}
	for _, satSolver := range satSolvers {
		context.progress[satSolver] = make([]bool, instancesCount)
	}

	os.Remove(BENCHMARK_LOG_FILE_NAME)

	// Solve the encodings for each SAT solver
	for _, satSolver := range satSolvers {
		var i uint = 0

		for _, steps := range stepVariations {
			for _, hash := range hashes {
				for _, xorOption := range xorOptions {
					for _, adderType := range adderTypes {
						filepath := fmt.Sprintf("%smd4_%d_%s_xor%d_%s.cnf",
							ENCODINGS_DIR_PATH, steps, adderType, xorOption, hash)

						startTime := time.Now()
						switch satSolver {
						case CRYPTOMINISAT:
							go cryptoMiniSat(filepath, context, i, startTime)
						case KISSAT:
							go kissat(filepath, context, i, startTime)
						case CADICAL:
							go cadical(filepath, context, i, startTime)
						case MAPLESAT:
							go mapleSat(filepath, context, i, startTime)
						case GLUCOSE:
							go glucose(filepath, context, i, startTime)
						}

						i++
					}
				}
			}
		}

		fmt.Printf("Spawned %d instances of %s.\n", instancesCount, satSolver)
	}

	for !areAllInstancesCompleted(context) {

	}
}
