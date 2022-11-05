package core

import (
	"benchmark/constants"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"log"
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
		exiterr, _ := err.(*exec.ExitError)
		exitCode = exiterr.ExitCode()
	}

	duration := time.Since(startTime)
	instanceName := strings.TrimSuffix(path.Base(filepath), ".cnf")
	benchmarkLogFilePath := "benchmark_" + instanceName + "_" + satSolver + ".csv"
	validResultsLogFilePath := "valid_results_" + instanceName + "_" + satSolver + ".csv"
	verificationLogFilePath := "verification_" + instanceName + "_" + satSolver + ".csv"

	// Log down to a file
	logRecord := []string{satSolver, instanceName, fmt.Sprintf("%.2f", duration.Seconds()), strconv.Itoa(exitCode)}
	utils.AppendLog(benchmarkLogFilePath, logRecord)

	// Normalize the solution
	{
		command := fmt.Sprintf("%s %s%s_%s.sol normalize > /tmp/%s_%s.sol && cat /tmp/%s_%s.sol > %s%s_%s.sol", constants.SolutionAnalyzerBinPath, constants.SolutionsDirPath, satSolver, instanceName, satSolver, instanceName, satSolver, instanceName, constants.SolutionsDirPath, satSolver, instanceName)
		cmd := exec.Command("bash", "-c", command)
		if err := cmd.Run(); err != nil {
			utils.AppendLog(verificationLogFilePath, []string{satSolver, instanceName, fmt.Sprintf("Normalization failed: %s %s", err.Error(), cmd.String())})
		}
	}

	// Verify the solution
	{
		steps, err := strconv.Atoi(strings.Split(instanceName, "_")[2])
		if err != nil {
			utils.AppendLog(verificationLogFilePath, []string{satSolver, instanceName, "Verification failed"})
		}

		command := fmt.Sprintf("%s %d < %s%s_%s.sol", constants.VerifierBinPath, steps, constants.SolutionsDirPath, satSolver, instanceName)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			utils.AppendLog(verificationLogFilePath, []string{satSolver, instanceName, "Verification failed"})
		}

		if strings.Contains(string(output), "Solution's hash matches the target!") {
			utils.AppendLog(verificationLogFilePath, []string{satSolver, instanceName, "Valid"})
			utils.AppendLog(validResultsLogFilePath, logRecord)
		} else if strings.Contains(string(output), "Solution's hash DOES NOT match the target:") || strings.Contains(string(output), "Result is UNSAT!") {
			utils.AppendLog(verificationLogFilePath, []string{satSolver, instanceName, "Verification failed"})
		} else {
			utils.AppendLog(verificationLogFilePath, []string{satSolver, instanceName, fmt.Sprintf("Unknown error: %s", output)})
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
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s --verb=0 %s > %scryptominisat_%ssol", constants.CryptoMiniSatBinPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func Kissat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := KissatCmd(filepath)

	invokeSatSolver(command, constants.Kissat, context, filepath, startTime, instanceIndex, maxTime)
}

func KissatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %skissat_%ssol", constants.KissatBinPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func Cadical(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := CadicalCmd(filepath)

	invokeSatSolver(command, constants.Cadical, context, filepath, startTime, instanceIndex, maxTime)
}

func CadicalCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %scadical_%ssol", constants.CadicalBinPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func MapleSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := MapleSatCmd(filepath)

	invokeSatSolver(command, constants.MapleSat, context, filepath, startTime, instanceIndex, maxTime)
}

func MapleSatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %smaplesat_%ssol", constants.MapleSatBinPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func XnfSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := XnfSatCmd(filepath)

	invokeSatSolver(command, constants.XnfSat, context, filepath, startTime, instanceIndex, maxTime)
}

func XnfSatCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s --witness --verbose=0 %s > %sxnfsat_%ssol", constants.XnfSatBinPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func Glucose(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := GlucoseCmd(filepath)

	invokeSatSolver(command, constants.Glucose, context, filepath, startTime, instanceIndex, maxTime)
}

func GlucoseCmd(filepath string) string {
	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %sglucose_%ssol", constants.GlucoseBinPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func March(filepath string, maxDepth uint) {
	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s %s -d %d -o %s%sicnf", constants.MarchBinPath, filepath, maxDepth, constants.EncodingsDirPath, instanceName)
	if err := exec.Command("bash", "-c", command).Run(); err != nil {
		log.Fatal("Failed to generate cubes with March", err)
	}
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
