package services

import (
	"benchmark/internal/consts"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func (logSvc *LogService) WriteLog(filePath string, handler func(writer *csv.Writer)) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	logSvc.errorSvc.Fatal(err, "Logger: failed to open "+filePath)
	defer file.Close()
	writer := csv.NewWriter(file)

	handler(writer)

	writer.Flush()
	err = writer.Error()
	logSvc.errorSvc.Fatal(err, "Logger: failed to write to the log file")
}

func (logSvc *LogService) WriteCuberLog(basePath string) {
	filePath := basePath + ".cubes.csv"
	logSvc.WriteLog(filePath, func(writer *csv.Writer) {
		cubesets, err := logSvc.cubesetSvc.All()
		logSvc.errorSvc.Fatal(err, "Logger: failed to retrieve cubesets")
		sort.Slice(cubesets, func(i, j int) bool {
			if cubesets[i].InstanceName != cubesets[j].InstanceName {
				return cubesets[i].InstanceName > cubesets[j].InstanceName
			}

			return cubesets[i].Threshold < cubesets[j].Threshold
		})

		writer.Write([]string{"Threshold", "Cubes", "Refuted Leaves", "Runtime", "Encoding"})
		for _, cubeset := range cubesets {
			threshold := strconv.Itoa(cubeset.Threshold)
			cubes := strconv.Itoa(cubeset.Cubes)
			refutedLeaves := strconv.Itoa(cubeset.RefutedLeaves)
			runtime := fmt.Sprintf("%.3f", cubeset.Runtime.Seconds())
			instanceName := cubeset.InstanceName
			writer.Write([]string{threshold, cubes, refutedLeaves, runtime, instanceName})
		}
	})
}

func (logSvc *LogService) WriteSolverLog(basePath string) {
	filePath := basePath + ".solutions.csv"
	logSvc.WriteLog(filePath, func(writer *csv.Writer) {
		solutions, err := logSvc.solutionSvc.All()
		logSvc.errorSvc.Fatal(err, "Logger: failed to retrieve solutions")
		sort.Slice(solutions, func(i, j int) bool {
			if solutions[i].InstanceName != solutions[j].InstanceName {
				return solutions[i].InstanceName > solutions[j].InstanceName
			}

			return int(solutions[i].Result) < int(solutions[j].Result)
		})

		writer.Write([]string{"Result", "Exit Code", "Runtime", "Solver", "Encoding"})
		for _, solution := range solutions {
			runtime := fmt.Sprintf("%.3f", solution.Runtime.Seconds())
			exitCode := strconv.Itoa(solution.ExitCode)
			var result string
			switch solution.Result {
			case consts.Sat:
				result = "SAT"
			case consts.Unsat:
				result = "UNSAT"
			default:
				result = "Fail"
			}
			solver := solution.Solver
			writer.Write([]string{result, exitCode, runtime, string(solver), solution.InstanceName})
		}
	})
}

func (logSvc *LogService) Run(basePath string) {
	logSvc.WriteCuberLog(basePath)
	logSvc.WriteSolverLog(basePath)
}
