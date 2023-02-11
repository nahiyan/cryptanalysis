package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/simplifier"
	"benchmark/internal/solver"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type solution struct {
	name        string
	processTime time.Duration
	result      solver.Result
	solver      solver.Solver
	verified    bool
}

type cubeset struct {
	name               string
	processTime        time.Duration
	cubesCount         int
	refutedLeavesCount int
	threshold          int
}

type combination struct {
	solutions *[]solution
	cubesets  *[]cubeset
}

type simplification struct {
	name                 string
	processTime          time.Duration
	numVars              int
	numClauses           int
	numEliminatedVars    int
	numEliminatedClauses int
	conflicts            int
	simplifier           simplifier.Simplifier
}

func (summarizerSvc *SummarizerService) GetCubesets(logFiles []string) []cubeset {
	config := summarizerSvc.configSvc.Config
	cubesets := []cubeset{}
	for _, logFile := range logFiles {
		processTime, cubesCount, refutedLeavesCount, err := summarizerSvc.cuberSvc.ParseOutput(path.Join(config.Paths.Logs, logFile))
		if err != nil {
			continue
		}

		name := path.Base(logFile)[:len(logFile)-4]
		threshold_ := regexp.MustCompile("(march_n)[0-9]+").Find([]byte(name))[7:]
		threshold, err := strconv.Atoi(string(threshold_))
		summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to get the threshold of the cubeset")

		cubesets = append(cubesets, cubeset{
			name:               name,
			processTime:        processTime,
			cubesCount:         cubesCount,
			refutedLeavesCount: refutedLeavesCount,
			threshold:          threshold,
		})
	}
	return cubesets
}

func (summarizerSvc *SummarizerService) GetSolutions(logFiles []string) []solution {
	config := summarizerSvc.configSvc.Config
	solutions := []solution{}
	for _, logFile := range logFiles {
		name := path.Base(logFile)[:len(logFile)-4]
		var solver_ solver.Solver
		{
			segments := strings.Split(name, ".")
			solver_ = solver.Solver(segments[len(segments)-1])
		}
		solutionLiterals := make([]int, 0)
		result, processTime, err := summarizerSvc.solverSvc.ParseLog(path.Join(config.Paths.Logs, logFile), solver_, &solutionLiterals)
		if err != nil {
			continue
		}

		solution := solution{
			name:        name,
			processTime: processTime,
			result:      result,
			solver:      solver_,
		}

		if result != solver.Sat {
			solutions = append(solutions, solution)
			continue
		}

		fileName := path.Base(logFile)

		// TODO: Remap SatELite simplifications
		// Reconstruct CaDiCaL simplifications
		if strings.Contains(name, simplifier.Cadical+"_c") {
			segments := strings.Split(fileName, ".")
			instanceFilePath := path.Join(config.Paths.Encodings, strings.Join(segments[:3], ".")+".cnf")
			originalFilePath := path.Join(config.Paths.Encodings, strings.Join(segments[:2], "."))
			rsFilePath := instanceFilePath + ".rs.txt"

			solutionLiterals, err = summarizerSvc.solutionSvc.Reconstruct(solutionLiterals, rsFilePath, originalFilePath)
			summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to reconstruct solution")
		}

		// Extract the message from the solution
		message, err := summarizerSvc.solutionSvc.ExtractMessage(solutionLiterals)
		summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed extract message from the solution literal")

		// Take the steps derived from the instance name
		var (
			encoder_   encoder.Encoder
			step       int
			targetHash string
		)
		{
			match := regexp.MustCompile("[a-z_]+_md4_[0-9]+_[a-z]+_").FindString(fileName)
			encoder_ = encoder.Encoder(strings.Split(match, "_md4")[0])
			segments := strings.Split(match[len(encoder_):], "_")
			step, _ = strconv.Atoi(segments[2])
			targetHash = segments[3]
		}

		// Verify the solution
		addChainingVars := encoder_ == encoder.SaeedE
		hash, err := summarizerSvc.md4Svc.Run(message, step, addChainingVars)
		summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to generate the hash")
		solution.verified = hash == targetHash

		solutions = append(solutions, solution)
	}
	return solutions
}

func (summarizerSvc *SummarizerService) GetSimplifications(logFiles []string) []simplification {
	config := summarizerSvc.configSvc.Config
	simplifications := []simplification{}
	for _, logFile := range logFiles {
		name := path.Base(logFile)[:len(logFile)-4]
		var (
			simplifier_ simplifier.Simplifier
			conflicts   int
			err         error
		)
		{
			segments := strings.Split(name, ".")
			segment := segments[len(segments)-1]
			segments_ := strings.Split(segment, "_")
			segment_ := segments_[0]
			simplifier_ = simplifier.Simplifier(segment_)
			if simplifier_ == simplifier.Cadical {
				conflicts, err = strconv.Atoi(segments_[1][1:])
				summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to get the conflicts")
			}
		}
		result, err := summarizerSvc.simplifierSvc.ParseOutput(path.Join(config.Paths.Logs, logFile), simplifier_)
		if err != nil {
			continue
		}

		simplifications = append(simplifications, simplification{
			name:                 name,
			processTime:          result.ProcessTime,
			numVars:              result.NumVars,
			numClauses:           result.NumClauses,
			numEliminatedVars:    result.NumEliminatedVars,
			numEliminatedClauses: result.NumEliminatedClauses,
			conflicts:            conflicts,
			simplifier:           simplifier_,
		})
	}
	return simplifications
}

func (summarizerSvc *SummarizerService) GetCombinations(solutions []solution, cubesets []cubeset) map[string]combination {
	combinations := make(map[string]combination)

	for _, solution_ := range solutions {
		baseName := strings.Split(solution_.name, ".")[0]

		if _, exists := combinations[baseName]; !exists {
			combinations[baseName] = combination{
				solutions: &[]solution{},
				cubesets:  &[]cubeset{},
			}
		}

		combination_ := combinations[baseName]
		*combination_.solutions = append(*combination_.solutions, solution_)
	}

	for _, cubeset_ := range cubesets {
		baseName := strings.Split(cubeset_.name, ".")[0]

		if _, exists := combinations[baseName]; !exists {
			combinations[baseName] = combination{
				solutions: &[]solution{},
				cubesets:  &[]cubeset{},
			}
		}

		combination_ := combinations[baseName]
		*combination_.cubesets = append(*combination_.cubesets, cubeset_)
	}

	return combinations
}

func (summarizerSvc *SummarizerService) WriteLog(filePath string, handler func(writer *csv.Writer)) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	summarizerSvc.errorSvc.Fatal(err, "Logger: failed to open "+filePath)
	defer file.Close()
	writer := csv.NewWriter(file)

	handler(writer)

	writer.Flush()
	err = writer.Error()
	summarizerSvc.errorSvc.Fatal(err, "Logger: failed to write to the log file")
}

func (summarizerSvc *SummarizerService) WriteCubesetsLog(cubesets []cubeset) {
	filePath := "summary.cubesets.csv"

	summarizerSvc.WriteLog(filePath, func(writer *csv.Writer) {
		sort.Slice(cubesets, func(i, j int) bool {
			if cubesets[i].name != cubesets[j].name {
				return cubesets[i].name > cubesets[j].name
			}

			return cubesets[i].threshold < cubesets[j].threshold
		})

		writer.Write([]string{"Threshold", "Cubes", "Refuted Leaves", "Runtime", "Name"})
		for _, cubeset := range cubesets {
			name := cubeset.name
			threshold := strconv.Itoa(cubeset.threshold)
			cubesCount := strconv.Itoa(cubeset.cubesCount)
			refutedLeavesCount := strconv.Itoa(cubeset.refutedLeavesCount)
			processTime := fmt.Sprintf("%.3f", cubeset.processTime.Seconds())
			writer.Write([]string{threshold, cubesCount, refutedLeavesCount, processTime, name})
		}
	})
}

func (summarizerSvc *SummarizerService) WriteSolutionsLog(solutions []solution) {
	filePath := "summary.solutions.csv"

	summarizerSvc.WriteLog(filePath, func(writer *csv.Writer) {
		sort.Slice(solutions, func(i, j int) bool {
			if solutions[i].result != solutions[j].result {
				return solutions[i].result < solutions[j].result
			}

			return solutions[i].name < solutions[j].name
		})

		writer.Write([]string{"Result", "Process Time", "Solver", "Instance Name"})
		for _, solution := range solutions {
			name := solution.name
			result := string(solution.result)
			if solution.result == solver.Sat {
				if solution.verified {
					result += "✔"
				} else {
					result += "✖"
				}
			}
			processTime := fmt.Sprintf("%.3f", solution.processTime.Seconds())
			solver_ := solution.solver
			writer.Write([]string{result, processTime, string(solver_), name})
		}
	})
}

func (summarizerSvc *SummarizerService) WriteSimplificationsLog(simplifications []simplification) {
	filePath := "summary.simplifications.csv"

	summarizerSvc.WriteLog(filePath, func(writer *csv.Writer) {
		sort.Slice(simplifications, func(i, j int) bool {
			if simplifications[i].name != simplifications[j].name {
				return simplifications[i].name < simplifications[j].name
			}

			return simplifications[i].numVars < simplifications[j].numVars
		})

		writer.Write([]string{"Conflicts", "Variables", "Eliminated Variables", "Simplifier", "Clauses", "Eliminated Clauses", "Process Time", "Name"})
		for _, simplification := range simplifications {
			name := simplification.name
			processTime := fmt.Sprintf("%.3f", simplification.processTime.Seconds())
			simplifier_ := simplification.simplifier
			conflicts := strconv.Itoa(simplification.conflicts)
			variables := strconv.Itoa(simplification.numVars)
			clauses := strconv.Itoa(simplification.numClauses)
			elimination := strconv.Itoa(simplification.numEliminatedVars)
			eliminationClauses := strconv.Itoa(simplification.numEliminatedClauses)
			writer.Write([]string{conflicts, variables, elimination, string(simplifier_), clauses, eliminationClauses, processTime, name})
		}
	})
}

func (summarizerSvc *SummarizerService) WriteCombinationsLog(combinations map[string]combination) {
	filePath := "summary.combined.md"
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to write summary.combined.md file")

	baseEncodingNames := make([]string, 0, len(combinations))
	for k := range combinations {
		baseEncodingNames = append(baseEncodingNames, k)
	}
	sort.Slice(baseEncodingNames, func(i, j int) bool {
		return baseEncodingNames[i] > baseEncodingNames[j]
	})

	for _, baseEncodingName := range baseEncodingNames {
		file.WriteString("# " + baseEncodingName + "\n")
		combination := combinations[baseEncodingName]

		solutions := *combination.solutions
		file.WriteString("\n## Solutions\n\n")
		satCount := 0
		satVerifiedCount := 0
		unsatCount := 0
		failsCount := 0
		totalTime := time.Duration(0)

		sort.Slice(solutions, func(i, j int) bool {
			return solutions[i].name > solutions[j].name
		})

		for _, solution := range solutions {
			switch solution.result {
			case solver.Sat:
				satCount++
			case solver.Unsat:
				unsatCount++
			case solver.Fail:
				failsCount++
			}
			if solution.verified {
				satVerifiedCount++
			}
			totalTime += solution.processTime
		}
		satCount_ := humanize.Comma(int64(satCount))
		satVerifiedComment := ""
		if satCount > 0 {
			if satCount == satVerifiedCount {
				satVerifiedComment = " (✔)"
			} else {
				satVerifiedComment = fmt.Sprintf(" (✖%d)", satCount-satVerifiedCount)
			}
		}
		unsatCount_ := humanize.Comma(int64(unsatCount))
		failCount_ := humanize.Comma(int64(failsCount))

		file.WriteString(fmt.Sprintf("%s SAT%s, %s UNSAT, %s Fails\n", satCount_, satVerifiedComment, unsatCount_, failCount_))
		file.WriteString(fmt.Sprintf("Process time: %s\n\n", totalTime))

		cubesets := *combination.cubesets
		sort.Slice(cubesets, func(i, j int) bool {
			return cubesets[i].name > cubesets[j].name
		})
		file.WriteString("## Cubesets\n\n")
		for _, cubeset := range cubesets {
			file.WriteString(cubeset.name + "\n")
		}
		if len(cubesets) > 0 {
			file.WriteString("\n")
		}
	}
}

// func (summarizerSvc *SummarizerService) WriteSummaryLog(basePath string) {
// 	filePath := basePath + ".summary.md"
// 	summary := ""
// 	solutions, err := summarizerSvc.solutionSvc.All()
// 	summarizerSvc.errorSvc.Fatal(err, "Logger: failed to fetch solutions")
// 	cubesets, err := summarizerSvc.cubesetSvc.All()
// 	summarizerSvc.errorSvc.Fatal(err, "Logger: failed to fetch cubesets")

// 	summary += "# Solutions\n\n"
// 	groupedSolutions := lo.GroupBy(solutions, func(s solver.Solution) string {
// 		instanceName := s.InstanceName
// 		lastCnfIndex := strings.LastIndex(instanceName, ".cnf")
// 		lastCubesIndex := strings.LastIndex(instanceName, ".cubes")

// 		if lastCubesIndex != -1 {
// 			return instanceName[:lastCubesIndex]
// 		}

// 		if lastCnfIndex != -1 {
// 			return instanceName[:lastCnfIndex]
// 		}

// 		return instanceName
// 	})

// 	encodings := lo.Keys(groupedSolutions)
// 	sort.Strings(encodings)
// 	for _, encoding := range encodings {
// 		sat := 0
// 		unsat := 0
// 		others := 0
// 		totalTime := time.Duration(0)
// 		quantity := 0
// 		cubesCount := 0
// 		encodingInfo, err := summarizerSvc.encoderSvc.ProcessInstanceName(encoding)
// 		summarizerSvc.errorSvc.Fatal(err, "Logger: failed to process instance name")

// 		for _, cubeset := range cubesets {
// 			info, err := summarizerSvc.encoderSvc.ProcessInstanceName(cubeset.InstanceName)
// 			summarizerSvc.errorSvc.Fatal(err, "Logger: failed to process instance name")
// 			info.Cubing = mo.Some(encoder.CubingInfo{
// 				Threshold: cubeset.Threshold,
// 			})
// 			if encoding == strings.TrimSuffix(summarizerSvc.encoderSvc.GetInstanceName(info), ".cubes") {
// 				cubesCount = cubeset.Cubes
// 			}
// 		}

// 		solutions := groupedSolutions[encoding]
// 		for _, solution := range solutions {
// 			switch solution.ExitCode {
// 			case 10:
// 				sat += 1
// 			case 20:
// 				unsat += 1
// 			default:
// 				others += 1
// 			}

// 			if solution.ExitCode == 20 || solution.ExitCode == 10 {
// 				totalTime += solution.Runtime
// 				quantity += 1
// 			}
// 		}

// 		sat_ := humanize.Comma(int64(sat))
// 		unsat_ := humanize.Comma(int64(unsat))
// 		others_ := humanize.Comma(int64(others))
// 		quantity_ := humanize.Comma(int64(quantity))
// 		cubesCount_ := humanize.Comma(int64(cubesCount))

// 		summary += fmt.Sprintf("## %s\n\n%s SAT, %s UNSAT, %s Others", encoding, sat_, unsat_, others_)
// 		if _, exists := encodingInfo.Cubing.Get(); exists {
// 			percentageComplete := float64(quantity) / float64(cubesCount) * 100
// 			summary += fmt.Sprintf(", %.2f%% complete\n", percentageComplete)

// 			estimate := time.Duration(totalTime.Seconds()/float64(quantity)*float64(cubesCount)) * time.Second
// 			summary += fmt.Sprintf("Estimate (1 CPU): %s for %s cubes\n", estimate.Round(time.Millisecond), cubesCount_)

// 			estimate12Cpu := time.Duration(totalTime.Seconds()/float64(quantity)*float64(cubesCount)) * time.Second / 12
// 			summary += fmt.Sprintf("Estimate (12 CPU): %s for %s cubes\n", estimate12Cpu.Round(time.Millisecond), cubesCount_)

// 			estimateForNRemaining12Cpu := time.Duration(totalTime.Seconds()/float64(quantity)*float64(cubesCount-1000)) * time.Second / 12
// 			summary += fmt.Sprintf("Estimate (12 CPU): %s for %s cubes\n", estimateForNRemaining12Cpu.Round(time.Millisecond), humanize.Comma(int64(cubesCount)-1000))

// 			summary += fmt.Sprintf("Real time (1 CPU): %s for %s cubes\n", totalTime.Round(time.Millisecond), quantity_)
// 			summary += fmt.Sprintf("Real time (12 CPU): %s for %s cubes\n", (totalTime / 12).Round(time.Millisecond), quantity_)
// 		} else {
// 			summary += "\n"
// 		}

// 		summary += "\n"
// 	}

// 	summary += "# Cubesets\n"
// 	groupedCubesets := lo.GroupBy(cubesets, func(cubeset cubeset.CubeSet) string {
// 		instanceName := cubeset.InstanceName
// 		return instanceName
// 	})

// 	encodings = lo.Keys(groupedCubesets)
// 	sort.Strings(encodings)
// 	for _, encoding := range encodings {
// 		cubesets := groupedCubesets[encoding]
// 		summary += "\n## " + encoding + "\n\n"

// 		sort.Slice(cubesets, func(i, j int) bool {
// 			return cubesets[i].Threshold < cubesets[j].Threshold
// 		})

// 		for i, cubeset := range cubesets {
// 			threshold := cubeset.Threshold
// 			cubes := cubeset.Cubes
// 			refutedLeaves := cubeset.RefutedLeaves
// 			processTime := cubeset.Runtime

// 			cubes_ := humanize.Comma(int64(cubes))
// 			refutedLeaves_ := humanize.Comma(int64(refutedLeaves))

// 			summary += fmt.Sprintf("%d. n%d: %s cubes, %s refuted leaves, %s process time\n", i+1, threshold, cubes_, refutedLeaves_, processTime.Round(time.Millisecond))
// 		}
// 	}

// 	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
// 	summarizerSvc.errorSvc.Fatal(err, "Logger: failed to write summary")
// 	file.WriteString(summary)
// }

func (summarizerSvc *SummarizerService) Run() {
	startTime := time.Now()
	files, err := os.ReadDir(summarizerSvc.configSvc.Config.Paths.Logs)
	summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to find log files")

	solutionLogFiles := []string{}
	cubesetLogFiles := []string{}
	simplificationLogFiles := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if strings.Contains(fileName, solver.Kissat) || strings.Contains(fileName, solver.Cadical+".log") || strings.Contains(fileName, solver.CryptoMiniSat) || strings.Contains(fileName, solver.Glucose) || strings.Contains(fileName, solver.MapleSat) {
			solutionLogFiles = append(solutionLogFiles, fileName)
		} else if strings.Contains(fileName, "march") {
			cubesetLogFiles = append(cubesetLogFiles, fileName)
		} else if strings.Contains(fileName, "satelite") || strings.Contains(fileName, "cadical_") {
			simplificationLogFiles = append(simplificationLogFiles, fileName)
		}
	}

	cubesets := summarizerSvc.GetCubesets(cubesetLogFiles)
	summarizerSvc.WriteCubesetsLog(cubesets)

	solutions := summarizerSvc.GetSolutions(solutionLogFiles)
	summarizerSvc.WriteSolutionsLog(solutions)

	simplifications := summarizerSvc.GetSimplifications(simplificationLogFiles)
	summarizerSvc.WriteSimplificationsLog(simplifications)

	combinations := summarizerSvc.GetCombinations(solutions, cubesets)
	summarizerSvc.WriteCombinationsLog(combinations)

	log.Println("Time taken:", time.Since(startTime))
}
