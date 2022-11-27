package core

import (
	"benchmark/config"
	"benchmark/constants"
	"benchmark/types"
	"benchmark/utils"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

func invokeSatSolver(command string, satSolver string, context_ *types.BenchmarkContext, filepath string, startTime time.Time, instanceIndex uint, maxTime uint) {
	messages := make([]string, 0)
	validity := constants.Undetermined
	// TODO: Improve the way the instance name is generated
	instanceName := strings.TrimSuffix(path.Base(filepath), ".cnf")
	solutionFilePath := fmt.Sprintf("%s%s_%s.sol", constants.SolutionsDirPath, satSolver, instanceName)

	// * 1. Check if the solution analyzer exists
	if !utils.FileExists(config.Get().Paths.Bin.SolutionAnalyzer) {
		log.Fatal("Solution analyzer doesn't exist. Did you forget to compile it?")
	}

	// * 2. Check if the solution already exists
	isPreviouslySolved := func(solutionFilePath string) bool {
		if !utils.FileExists(solutionFilePath) {
			return false
		}

		stat, _ := os.Stat(solutionFilePath)
		return stat.Size() == 0
	}(solutionFilePath)
	if isPreviouslySolved {
		messages = append(messages, "Solution already exists")
	}

	// * 3. Invoke the SAT solver
	exitCode := 0
	duration := time.Since(time.Now())
	if !isPreviouslySolved {
		cmd := exec.Command("timeout", strconv.Itoa(int(maxTime)), "bash", "-c", command)
		if err := cmd.Run(); err != nil {
			exiterr, _ := err.(*exec.ExitError)
			exitCode = exiterr.ExitCode()
		}
		duration = time.Since(startTime)
	}

	// * 4. Normalize the solution
	isNormalized := false
	{
		command := fmt.Sprintf("%s %s normalize > /tmp/%s.sol && cat /tmp/%s.sol > %s", config.Get().Paths.Bin.SolutionAnalyzer, solutionFilePath, path.Base(solutionFilePath), path.Base(solutionFilePath), solutionFilePath)
		cmd := exec.Command("bash", "-c", command)
		if err := cmd.Run(); err != nil {
			messages = append(messages, fmt.Sprintf("Normalization failed: %s %s", err.Error(), cmd.String()))
		} else {
			isNormalized = true
		}
	}

	// * 5. Validate the solution
	if isNormalized {
		// TODO: Write a better method to get the number of steps
		steps, err := strconv.Atoi(strings.Split(instanceName, "_")[2])
		if err != nil {
			messages = append(messages, "Error in the validation process: "+err.Error())
		}

		command := fmt.Sprintf("%s %d < %s", config.Get().Paths.Bin.Validator, steps, solutionFilePath)
		cmd := exec.Command("bash", "-c", command)
		output, err := cmd.Output()
		if err != nil {
			messages = append(messages, "Error in the validation process: "+err.Error())
		}

		if strings.Contains(string(output), "Solution's hash matches the target!") {
			validity = constants.Valid
		} else if strings.Contains(string(output), "Solution's hash DOES NOT match the target:") || strings.Contains(string(output), "Result is UNSAT!") {
			validity = constants.Invalid
		} else {
			messages = append(messages, "Error in the validation process: unknown")
		}
	}

	// * 6. Report the instance's completion
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

	// * 7. Cleanup failed solutions (not SAT or UNSAT)
	if validity == constants.Undetermined {
		if err := exec.Command("bash", "-c", fmt.Sprintf("rm %s", solutionFilePath)).Run(); err != nil {
			fmt.Println("Failed to cleanup an invalid solution file", err)
		}
	}

	fmt.Printf("[%d/%d] %s_%s %.2fs %d %s\n", completedInstancesCount, totalInstancesCount, satSolver, instanceName, duration.Seconds(), exitCode, validity)

	utils.AppendLog(satSolver, instanceName, duration, messages, exitCode, validity)

	context_.RunningInstances -= 1
	context_.Progress[satSolver][instanceIndex] = true
}

func CryptoMiniSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := CryptoMiniSatCmd(filepath)

	invokeSatSolver(command, constants.CryptoMiniSat, context, filepath, startTime, instanceIndex, maxTime)
}

func CryptoMiniSatCmd(filepath string) string {
	binPath := config.Get().Paths.Bin.CryptoMiniSat
	if !utils.FileExists(binPath) {
		log.Fatal("CryptoMiniSAT doesn't exist. Did you forget to compile it?")
	}

	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s --verb=0 %s > %scryptominisat_%ssol",
		binPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func Kissat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := KissatCmd(filepath)

	invokeSatSolver(command, constants.Kissat, context, filepath, startTime, instanceIndex, maxTime)
}

func KissatCmd(filepath string) string {
	binPath := config.Get().Paths.Bin.Kissat
	if !utils.FileExists(binPath) {
		log.Fatal("Kissat doesn't exist. Did you forget to compile it?")
	}

	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %skissat_%ssol", binPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func Cadical(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := CadicalCmd(filepath)

	invokeSatSolver(command, constants.Cadical, context, filepath, startTime, instanceIndex, maxTime)
}

func CadicalCmd(filepath string) string {
	binPath := config.Get().Paths.Bin.Cadical
	if !utils.FileExists(binPath) {
		log.Fatal("CaDiCaL doesn't exist. Did you forget to compile it?")
	}

	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -q %s > %scadical_%ssol", binPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func MapleSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := MapleSatCmd(filepath)

	invokeSatSolver(command, constants.MapleSat, context, filepath, startTime, instanceIndex, maxTime)
}

func MapleSatCmd(filepath string) string {
	binPath := config.Get().Paths.Bin.MapleSat
	if !utils.FileExists(binPath) {
		log.Fatal("MapleSAT doesn't exist. Did you forget to compile it?")
	}

	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %smaplesat_%ssol", binPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func XnfSat(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := XnfSatCmd(filepath)

	invokeSatSolver(command, constants.XnfSat, context, filepath, startTime, instanceIndex, maxTime)
}

func XnfSatCmd(filepath string) string {
	binPath := config.Get().Paths.Bin.XnfSat
	if !utils.FileExists(binPath) {
		log.Fatal("XNFSAT doesn't exist. Did you forget to compile it?")
	}

	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s --witness --verbose=0 %s > %sxnfsat_%ssol", binPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func Glucose(filepath string, context *types.BenchmarkContext, instanceIndex uint, startTime time.Time, maxTime uint) {
	command := GlucoseCmd(filepath)

	invokeSatSolver(command, constants.Glucose, context, filepath, startTime, instanceIndex, maxTime)
}

func GlucoseCmd(filepath string) string {
	binPath := config.Get().Paths.Bin.Glucose
	if !utils.FileExists(binPath) {
		log.Fatal("Glucose doesn't exist. Did you forget to compile it?")
	}

	baseFileName := path.Base(filepath)
	instanceName := baseFileName[:len(baseFileName)-3]

	command := fmt.Sprintf("%s -verb=0 %s %sglucose_%ssol", binPath, filepath, constants.SolutionsDirPath, instanceName)

	return command
}

func March(filePath string, outputPath string, cubeCutoffVars uint, maxTime time.Duration) string {
	cmd, cancel := MarchCmd(filePath, outputPath, cubeCutoffVars, maxTime)
	output_, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed to generate cubes with March", err)
	}
	defer cancel()

	return string(output_)
}

func MarchCmd(filePath string, outputPath string, cubeCutoffVars uint, maxTime time.Duration) (*exec.Cmd, context.CancelFunc) {
	binPath := config.Get().Paths.Bin.March
	if !utils.FileExists(binPath) {
		log.Fatal("March doesn't exist. Did you forget to compile it?")
	}

	command := fmt.Sprintf("%s %s -n %d -o %s", binPath, filePath, cubeCutoffVars, outputPath)
	context, cancel := context.WithTimeout(context.Background(), maxTime)
	return exec.CommandContext(context, "bash", "-c", command), cancel
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
