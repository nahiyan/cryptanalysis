package core

import (
	"benchmark/constants"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

func invokeSatSolver(command string, satSolver string, context_ *types.BenchmarkContext, filepath string, startTime time.Time, instanceIndex uint, maxTime uint) {
	cmd := exec.Command("timeout", strconv.Itoa(int(maxTime)), "bash", "-c", command)
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
	utils.AppendBenchmarkLog(logMessage)

	// Normalize the solution
	{
		command := fmt.Sprintf("%s %s%s/%s.sol normalize > /tmp/%s-%s.sol && cat /tmp/%s-%s.sol > %s%s/%s.sol", constants.SOLUTION_ANALYZER_BIN_PATH, constants.SOLUTIONS_DIR_PATH, satSolver, instanceName, satSolver, instanceName, satSolver, instanceName, constants.SOLUTIONS_DIR_PATH, satSolver, instanceName)
		cmd := exec.Command("bash", "-c", command)
		if err := cmd.Run(); err != nil {
			utils.AppendVerificationLog("Failed to normalize " + instanceName + " " + err.Error() + " " + cmd.String())
		}
	}

	// Verify the solution
	{
		steps, err := strconv.Atoi(strings.Split(instanceName, "_")[1])
		if err != nil {
			utils.AppendVerificationLog("Failed to verify " + instanceName)
		}

		command := fmt.Sprintf("%s %d < %s%s/%s.sol", constants.VERIFIER_BIN_PATH, steps, constants.SOLUTIONS_DIR_PATH, satSolver, instanceName)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			utils.AppendVerificationLog("Failed to verify " + instanceName)
		}

		if strings.Contains(string(output), "Solution's hash matches the target!") {
			utils.AppendVerificationLog(fmt.Sprintf("Valid: %s %s", satSolver, instanceName))
		} else if strings.Contains(string(output), "Solution's hash DOES NOT match the target:") || strings.Contains(string(output), "Result is UNSAT!") {
			utils.AppendVerificationLog(fmt.Sprintf("Invalid: %s %s", satSolver, instanceName))
		} else {
			utils.AppendVerificationLog(fmt.Sprintf("Unknown error: %s %s %s", satSolver, instanceName, output))
		}
	}

	// Report the instance's completion
	var (
		completedInstancesCount uint = 0
		totalInstancesCount     int  = 0
	)
	for satSolver_ := range context_.Progress {
		completedInstancesCount += lo.SumBy(context_.Progress[satSolver_], func(b bool) uint {
			if b {
				return 1
			} else {
				return 0
			}
		})

		totalInstancesCount += len(context_.Progress[satSolver_])
	}
	completedInstancesCount += 1

	fmt.Printf("[%d/%d] %s \t %s \t %.2fs \t exit code: %d\n", completedInstancesCount, totalInstancesCount, satSolver, instanceName, duration.Seconds(), exitCode)

	context_.RunningInstances -= 1
	context_.Progress[satSolver][instanceIndex] = true
}

func CryptoMiniSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := CryptoMiniSatCmd(filepath)

	invokeSatSolver(command, constants.CRYPTOMINISAT, context, filepath, startTime, instanceIndex, maxTime)
}

func CryptoMiniSatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s --verb=0 %s > %scryptominisat/%ssol", constants.CRYPTOMINISAT_BIN_PATH, filepath, constants.SOLUTIONS_DIR_PATH, solutionFilePath)

	return command
}

func Kissat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := KissatCmd(filepath)

	invokeSatSolver(command, constants.KISSAT, context, filepath, startTime, instanceIndex, maxTime)
}

func KissatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %skissat/%ssol", constants.KISSAT_BIN_PATH, filepath, constants.SOLUTIONS_DIR_PATH, solutionFilePath)

	return command
}

func Cadical(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := CadicalCmd(filepath)

	invokeSatSolver(command, constants.CADICAL, context, filepath, startTime, instanceIndex, maxTime)
}

func CadicalCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %scadical/%ssol", constants.CADICAL_BIN_PATH, filepath, constants.SOLUTIONS_DIR_PATH, solutionFilePath)

	return command
}

func MapleSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := MapleSatCmd(filepath)

	invokeSatSolver(command, constants.MAPLESAT, context, filepath, startTime, instanceIndex, maxTime)
}

func MapleSatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %smaplesat/%ssol", constants.MAPLESAT_BIN_PATH, filepath, constants.SOLUTIONS_DIR_PATH, solutionFilePath)

	return command
}

func Glucose(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := GlucoseCmd(filepath)

	invokeSatSolver(command, constants.GLUCOSE, context, filepath, startTime, instanceIndex, maxTime)
}

func GlucoseCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFilePath := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %sglucose/%ssol", constants.GLUCOSE_BIN_PATH, filepath, constants.SOLUTIONS_DIR_PATH, solutionFilePath)

	return command
}

func AreAllInstancesCompleted(context *types.BenchmarkContext) bool {
	for _, progressEntries := range context.Progress {
		if lo.SomeBy(progressEntries, func(done bool) bool {
			return !done
		}) {
			return false
		}
	}

	return true
}
