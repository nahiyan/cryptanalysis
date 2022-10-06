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
	CRYPTOMINISAT              = "cryptominisat"
	KISSAT                     = "kissat"
	CADICAL                    = "cadical"
	GLUCOSE                    = "glucose"
	MAPLESAT                   = "maplesat"
	CRYPTOMINISAT_BIN_PATH     = "../../../sat-solvers/cryptominisat"
	KISSAT_BIN_PATH            = "../../../sat-solvers/kissat"
	CADICAL_BIN_PATH           = "../../../sat-solvers/cadical"
	GLUCOSE_BIN_PATH           = "../../../sat-solvers/glucose"
	MAPLESAT_BIN_PATH          = "../../../sat-solvers/maplesat"
	VERIFIER_BIN_PATH          = "../../encoders/saeed/crypto/verify-md4"
	SOLUTION_ANALYZER_BIN_PATH = "../solution_analyzer/target/release/solution_analyzer"
	MAX_TIME                   = 5000
	BENCHMARK_LOG_FILE_NAME    = "benchmark.log"
	VERIFICATION_LOG_FILE_NAME = "verification.log"
	BASE_PATH                  = "../../"
	SOLUTIONS_DIR_PATH         = BASE_PATH + "solutions/saeed/"
	ENCODINGS_DIR_PATH         = BASE_PATH + "encodings/saeed/"
	MAX_INSTANCES_COUNT        = 50
)

// Variations
var (
	xorOptions = []uint{0}
	hashes     = []string{"ffffffffffffffffffffffffffffffff",
		"00000000000000000000000000000000"}
	adderTypes     = []string{"counter_chain", "dot_matrix"}
	stepVariations = makeRange(16, 28)

	satSolvers = []string{CRYPTOMINISAT, KISSAT, CADICAL, GLUCOSE, MAPLESAT}
)

type Context struct {
	progress         map[string][]bool
	runningInstances uint
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

	// Log down to a file
	instanceName := strings.TrimSuffix(path.Base(filepath), ".cnf")
	logMessage := fmt.Sprintf("Time: %.2fs, instance index: %d, instance name: %s, SAT solver: %s, exit code: %d", duration.Seconds(), instanceIndex, instanceName, satSolver, exitCode)
	appendBenchmarkLog(logMessage)

	// Normalize the solution
	{
		command := fmt.Sprintf("%s %s%s/%s.sol normalize > /tmp/%s-%s.sol && cat /tmp/%s-%s.sol > %s%s/%s.sol", SOLUTION_ANALYZER_BIN_PATH, SOLUTIONS_DIR_PATH, satSolver, instanceName, satSolver, instanceName, satSolver, instanceName, SOLUTIONS_DIR_PATH, satSolver, instanceName)
		cmd := exec.Command("bash", "-c", command)
		if err := cmd.Run(); err != nil {
			appendVerificationLog("Failed to normalize " + instanceName + " " + err.Error() + " " + cmd.String())
		}
	}

	// Verify the solution
	{
		steps, err := strconv.Atoi(strings.Split(instanceName, "_")[1])
		if err != nil {
			appendVerificationLog("Failed to verify " + instanceName)
		}

		command := fmt.Sprintf("%s %d < %s%s/%s.sol", VERIFIER_BIN_PATH, steps, SOLUTIONS_DIR_PATH, satSolver, instanceName)
		// fmt.Println(command)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			appendVerificationLog("Failed to verify " + instanceName)
		}

		if strings.Contains(string(output), "Solution's hash matches the target!") {
			appendVerificationLog(fmt.Sprintf("Valid: %s %s", satSolver, instanceName))
		} else if strings.Contains(string(output), "Solution's hash DOES NOT match the target:") || strings.Contains(string(output), "Result is UNSAT!") {
			appendVerificationLog(fmt.Sprintf("Invalid: %s %s", satSolver, instanceName))
		} else {
			appendVerificationLog(fmt.Sprintf("Unknown error: %s %s %s", satSolver, instanceName, output))
		}
	}

	// Report the instance's completion
	var (
		completedInstancesCount uint = 0
		totalInstancesCount     int  = 0
	)
	for satSolver_ := range context_.progress {
		completedInstancesCount += lo.SumBy(context_.progress[satSolver_], func(b bool) uint {
			if b {
				return 1
			} else {
				return 0
			}
		})

		totalInstancesCount += len(context_.progress[satSolver_])
	}
	completedInstancesCount += 1

	fmt.Printf("[%d/%d] %s \t %s \t %.2fs \t exit code: %d\n", completedInstancesCount, totalInstancesCount, satSolver, instanceName, duration.Seconds(), exitCode)

	context_.runningInstances -= 1
	context_.progress[satSolver][instanceIndex] = true
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

func appendLog(filename, message string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to write logs")
	}
	_, err = f.WriteString(message + "\n")
	if err != nil {
		panic("Failed to write logs")
	}
	f.Close()
}

func appendBenchmarkLog(message string) {
	appendLog(BENCHMARK_LOG_FILE_NAME, message)
}

func appendVerificationLog(message string) {
	appendLog(VERIFICATION_LOG_FILE_NAME, message)
}

func main() {
	// Should be 264 for all the possible variations
	instancesCount := len(xorOptions) * len(hashes) * len(adderTypes) * len(stepVariations)

	// Define the context
	context := &Context{
		progress: make(map[string][]bool),
	}
	for _, satSolver := range satSolvers {
		context.progress[satSolver] = make([]bool, instancesCount)
	}

	// Remove the files from previous execution
	os.Remove(BENCHMARK_LOG_FILE_NAME)
	os.Remove(VERIFICATION_LOG_FILE_NAME)
	for _, satSolver := range satSolvers {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("rm %s%s/*.sol", SOLUTIONS_DIR_PATH, satSolver))
		if err := cmd.Run(); err != nil {
			fmt.Println(cmd.String())
			fmt.Println("Failed to delete the solution files: " + err.Error())
		}
	}

	// Solve the encodings for each SAT solver
	for _, satSolver := range satSolvers {
		var i uint = 0

		for _, steps := range stepVariations {
			for _, hash := range hashes {
				for _, xorOption := range xorOptions {
					for _, adderType := range adderTypes {
						for context.runningInstances > MAX_INSTANCES_COUNT {
							time.Sleep(time.Second * 1)
						}

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

						context.runningInstances += 1
						i++
					}
				}
			}
		}
	}

	for !areAllInstancesCompleted(context) {
		time.Sleep(time.Second * 1)
	}
}
