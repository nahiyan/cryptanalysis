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

	benchmarkLogFilePath := constants.ResultsDirPat + "benchmark_" + instanceName + "_" + satSolver + ".log"
	validResultsLogFilePath := constants.ResultsDirPat + "valid_results_" + instanceName + "_" + satSolver + ".log"
	verificationLogFilePath := constants.ResultsDirPat + "verification_" + instanceName + "_" + satSolver + ".log"

	utils.AppendLog(benchmarkLogFilePath, logMessage)

	// Normalize the solution
	{
		command := fmt.Sprintf("%s %s%s/%s.sol normalize > /tmp/%s-%s.sol && cat /tmp/%s-%s.sol > %s%s/%s.sol", constants.SolutionAnalyzerBinPath, constants.SolutionsDirPath, satSolver, instanceName, satSolver, instanceName, satSolver, instanceName, constants.SolutionsDirPath, satSolver, instanceName)
		cmd := exec.Command("bash", "-c", command)
		if err := cmd.Run(); err != nil {
			utils.AppendLog(verificationLogFilePath, "Failed to normalize "+instanceName+" "+err.Error()+" "+cmd.String())
		}
	}

	// Verify the solution
	{
		steps, err := strconv.Atoi(strings.Split(instanceName, "_")[1])
		if err != nil {
			utils.AppendLog(verificationLogFilePath, "Failed to verify "+instanceName)
		}

		command := fmt.Sprintf("%s %d < %s%s/%s.sol", constants.VerifierBinPath, steps, constants.SolutionsDirPath, satSolver, instanceName)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			utils.AppendLog(verificationLogFilePath, "Failed to verify "+instanceName)
		}

		if strings.Contains(string(output), "Solution's hash matches the target!") {
			utils.AppendLog(verificationLogFilePath, fmt.Sprintf("Valid: %s %s", satSolver, instanceName))
			utils.AppendLog(validResultsLogFilePath, logMessage)
		} else if strings.Contains(string(output), "Solution's hash DOES NOT match the target:") || strings.Contains(string(output), "Result is UNSAT!") {
			utils.AppendLog(verificationLogFilePath, fmt.Sprintf("Invalid: %s %s", satSolver, instanceName))
		} else {
			utils.AppendLog(verificationLogFilePath, fmt.Sprintf("Unknown error: %s %s %s", satSolver, instanceName, output))
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

	invokeSatSolver(command, constants.CryptoMiniSat, context, filepath, startTime, instanceIndex, maxTime)
}

func CryptoMiniSatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFileName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s --verb=0 %s > %scryptominisat_%ssol", constants.CryptoMiniSatBinPath, filepath, constants.SolutionsDirPath, solutionFileName)

	return command
}

func Kissat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := KissatCmd(filepath)

	invokeSatSolver(command, constants.Kissat, context, filepath, startTime, instanceIndex, maxTime)
}

func KissatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFileName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %skissat_%ssol", constants.KissatBinPath, filepath, constants.SolutionsDirPath, solutionFileName)

	return command
}

func Cadical(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := CadicalCmd(filepath)

	invokeSatSolver(command, constants.Cadical, context, filepath, startTime, instanceIndex, maxTime)
}

func CadicalCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFileName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %scadical_%ssol", constants.CadicalBinPath, filepath, constants.SolutionsDirPath, solutionFileName)

	return command
}

func MapleSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := MapleSatCmd(filepath)

	invokeSatSolver(command, constants.MapleSat, context, filepath, startTime, instanceIndex, maxTime)
}

func MapleSatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFileName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %smaplesat_%ssol", constants.MapleSatBinPath, filepath, constants.SolutionsDirPath, solutionFileName)

	return command
}

func XnfSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := XnfSatCmd(filepath)

	invokeSatSolver(command, constants.XnfSat, context, filepath, startTime, instanceIndex, maxTime)
}

func XnfSatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFileName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s --witness --verbose=0 %s > %sxnfsat_%ssol", constants.XnfSatBinPath, filepath, constants.SolutionsDirPath, solutionFileName)

	return command
}

func Glucose(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := GlucoseCmd(filepath)

	invokeSatSolver(command, constants.Glucose, context, filepath, startTime, instanceIndex, maxTime)
}

func GlucoseCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	solutionFileName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %sglucose_%ssol", constants.GlucoseBinPath, filepath, constants.SolutionsDirPath, solutionFileName)

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
