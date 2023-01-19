package services

import (
	"benchmark/internal/consts"
	"benchmark/internal/cubeset"
	"benchmark/internal/solver"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/samber/lo"
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

			return solutions[i].ExitCode < solutions[j].ExitCode
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

func (logSvc *LogService) WriteSimplificationLog(basePath string) {
	filePath := basePath + ".simplifications.csv"
	logSvc.WriteLog(filePath, func(writer *csv.Writer) {
		simplifications, err := logSvc.simplificationSvc.All()
		logSvc.errorSvc.Fatal(err, "Logger: failed to retrieve simplifications")
		sort.Slice(simplifications, func(i, j int) bool {
			if simplifications[i].InstanceName != simplifications[j].InstanceName {
				return simplifications[i].InstanceName > simplifications[j].InstanceName
			}

			return simplifications[i].FreeVariables < simplifications[j].FreeVariables
		})

		writer.Write([]string{"Conflicts", "Free Variables", "Elimination", "Simplifier", "Clauses", "Process Time", "Instance Name"})
		for _, simplification := range simplifications {
			conflicts := fmt.Sprintf("%d", simplification.Conflicts)
			freeVariables := fmt.Sprintf("%d", simplification.FreeVariables)
			elimination := fmt.Sprintf("%d", simplification.Eliminaton)
			simplifier := simplification.Simplifier
			clauses := fmt.Sprintf("%d", simplification.Clauses)
			processTime := fmt.Sprintf("%.3f", simplification.ProcessTime.Seconds())
			instanceName := simplification.InstanceName

			writer.Write([]string{conflicts, freeVariables, elimination, simplifier, clauses, processTime, instanceName})
		}
	})
}

func (logSvc *LogService) WriteSummaryLog(basePath string) {
	filePath := basePath + ".summary.md"
	summary := ""
	solutions, err := logSvc.solutionSvc.All()
	logSvc.errorSvc.Fatal(err, "Logger: failed to fetch solutions")
	cubesets, err := logSvc.cubesetSvc.All()
	logSvc.errorSvc.Fatal(err, "Logger: failed to fetch cubesets")

	summary += "# Solutions\n\n"
	groupedSolutions := lo.GroupBy(solutions, func(s solver.Solution) string {
		instanceName := s.InstanceName
		lastCnfIndex := strings.LastIndex(instanceName, ".cnf")
		lastCubesIndex := strings.LastIndex(instanceName, ".cubes")

		if lastCubesIndex != -1 {
			return instanceName[:lastCubesIndex]
		}

		if lastCnfIndex != -1 {
			return instanceName[:lastCnfIndex]
		}

		return instanceName
	})

	encodings := lo.Keys(groupedSolutions)
	sort.Strings(encodings)
	for _, encoding := range encodings {
		sat := 0
		unsat := 0
		others := 0
		totalTime := time.Duration(0)
		quantity := 0
		cubesCount := 0

		for _, cubeset := range cubesets {
			if fmt.Sprintf("%s.cnf.n%d", cubeset.InstanceName, cubeset.Threshold) == encoding {
				cubesCount = cubeset.Cubes
			}
		}

		solutions := groupedSolutions[encoding]
		for _, solution := range solutions {
			switch solution.ExitCode {
			case 10:
				sat += 1
			case 20:
				unsat += 1
			default:
				others += 1
			}

			if solution.ExitCode == 20 || solution.ExitCode == 10 {
				totalTime += solution.Runtime
				quantity += 1
			}
		}

		sat_ := humanize.Comma(int64(sat))
		unsat_ := humanize.Comma(int64(unsat))
		others_ := humanize.Comma(int64(others))

		summary += fmt.Sprintf("## %s\n\n%s SAT, %s UNSAT, %s Others", encoding, sat_, unsat_, others_)
		percentageComplete := float64(quantity) / float64(cubesCount) * 100
		{
			split := strings.Split(encoding, ".")
			len := len(split)
			if strings.HasPrefix(split[len-1], "n") {
				summary += fmt.Sprintf(", %.2f%% complete", percentageComplete)
			}
		}
		summary += "\n"

		if quantity > 1 {
			estimate := time.Duration(totalTime.Seconds()/float64(quantity)*float64(cubesCount)) * time.Second
			summary += fmt.Sprintf("Estimate: %s\n", estimate.Round(time.Millisecond))
		}
		summary += "\n"
	}

	summary += "# Cubesets\n"
	groupedCubesets := lo.GroupBy(cubesets, func(cubeset cubeset.CubeSet) string {
		instanceName := cubeset.InstanceName
		return instanceName
	})

	encodings = lo.Keys(groupedCubesets)
	sort.Strings(encodings)
	for _, encoding := range encodings {
		cubesets := groupedCubesets[encoding]
		summary += "\n## " + encoding + "\n\n"

		sort.Slice(cubesets, func(i, j int) bool {
			return cubesets[i].Threshold < cubesets[j].Threshold
		})

		for i, cubeset := range cubesets {
			threshold := cubeset.Threshold
			cubes := cubeset.Cubes
			refutedLeaves := cubeset.RefutedLeaves
			processTime := cubeset.Runtime

			cubes_ := humanize.Comma(int64(cubes))
			refutedLeaves_ := humanize.Comma(int64(refutedLeaves))

			summary += fmt.Sprintf("%d. n%d: %s cubes, %s refuted leaves, %s process time\n", i+1, threshold, cubes_, refutedLeaves_, processTime.Round(time.Millisecond))
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	logSvc.errorSvc.Fatal(err, "Logger: failed to write summary")
	file.WriteString(summary)
}

func (logSvc *LogService) Run(basePath string) {
	logSvc.WriteCuberLog(basePath)
	logSvc.WriteSolverLog(basePath)
	logSvc.WriteSimplificationLog(basePath)
	logSvc.WriteSummaryLog(basePath)
}
