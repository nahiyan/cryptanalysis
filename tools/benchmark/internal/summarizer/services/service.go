package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/simplifier"
	"benchmark/internal/solver"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/dustin/go-humanize"
	"github.com/samber/lo"
	"gonum.org/v1/gonum/stat"
)

type solution struct {
	name        string
	processTime time.Duration
	runTime     time.Duration
	result      solver.Result
	solver      solver.Solver
	verified    bool
	message     string
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
		var (
			processTime        time.Duration
			cubesCount         int
			refutedLeavesCount int
			err                error
		)
		if summarizerSvc.combinedLogsSvc.IsLoaded() {
			processTime, cubesCount, refutedLeavesCount, err = summarizerSvc.cuberSvc.ParseOutputFromCombinedLog(logFile)
			if err != nil {
				continue
			}
		} else {
			processTime, cubesCount, refutedLeavesCount, err = summarizerSvc.cuberSvc.ParseOutputFromFile(path.Join(config.Paths.Logs, logFile))
			if err != nil {
				continue
			}
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

func parseSolutionLogName(name string) (encoder.Encoder, encoder.Function, int, string, error) {
	matches := regexp.MustCompile("([a-z_]+)_(md4|md5|sha256)_([0-9]+)_([a-z0-9]+)[_.]").FindAllStringSubmatch(name, len(name))
	groups := matches[0][1:]
	encoder_ := encoder.Encoder(groups[0])
	function := encoder.Function(groups[1])
	step, err := strconv.Atoi(groups[2])
	if err != nil {
		return encoder.Encoder(""), encoder.Md4, 0, "", err
	}
	targetHash := groups[3]

	return encoder_, function, step, targetHash, nil
}

func (summarizerSvc *SummarizerService) GetSolutions(logFiles []string, workers int) []solution {
	config := summarizerSvc.configSvc.Config
	solutions := []solution{}
	logFilesCount := len(logFiles)
	lock := sync.Mutex{}
	pool := pond.New(workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	for i, logFile := range logFiles {
		pool.Submit(func(i int, logFile string) func() {
			return func() {
				defer log.Printf("Solution: Read [%d/%d] file\n", i+1, logFilesCount)

				name := path.Base(logFile)[:len(logFile)-4]
				var solver_ solver.Solver
				{
					segments := strings.Split(name, ".")
					solver_ = solver.Solver(segments[len(segments)-1])
				}
				solutionLiterals := make([]int, 0)
				var (
					processTime time.Duration
					runTime     time.Duration
					result      solver.Result
					err         error
				)
				if summarizerSvc.combinedLogsSvc.IsLoaded() {
					result, processTime, runTime, err = summarizerSvc.solverSvc.ParseLogFromCombinedLog(logFile, solver_, &solutionLiterals)
				} else {
					result, processTime, runTime, err = summarizerSvc.solverSvc.ParseLogFromFile(path.Join(config.Paths.Logs, logFile), solver_, &solutionLiterals)
				}
				if err != nil {
					return
				}

				solution := solution{
					name:        name,
					processTime: processTime,
					runTime:     runTime,
					result:      result,
					solver:      solver_,
				}

				if result != solver.Sat {
					lock.Lock()
					solutions = append(solutions, solution)
					lock.Unlock()
					return
				}

				// The follow parts are executed only if it's a SAT solution
				if len(solutionLiterals) == 0 {
					log.Println("Summarizer: WARNING; failed to retrieve solution literals from", logFile)
				}

				fileName := path.Base(logFile)

				// TODO: Drop support for SatELite
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
				message, err := summarizerSvc.solutionSvc.ExtractFromLiterals(solutionLiterals[:512])
				summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed extract the message from the solution literal")

				// Take the steps derived from the instance name
				_, function, step, targetHash, err := parseSolutionLogName(fileName)
				summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to extract information from the log file name")

				// Verify the solution
				var hash string
				if function == encoder.Md4 {
					hash, err = summarizerSvc.md4Svc.Run(message, step, false)
					summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to generate the md4 hash")
				} else if function == encoder.Md5 {
					hash, err = summarizerSvc.md5Svc.Run(message, step, false)
					summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to generate the md5 hash")
				} else if function == encoder.Sha256 {
					hash, err = summarizerSvc.sha256Svc.Run(message, step, false)
					summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to generate the md5 hash")
				}
				solution.verified = hash == targetHash
				if solution.verified {
					solution.message = hex.EncodeToString(message)
				}

				// Add the solution to the list
				lock.Lock()
				solutions = append(solutions, solution)
				lock.Unlock()
			}
		}(i, logFile))
	}

	pool.StopAndWait()

	return solutions
}

func (summarizerSvc *SummarizerService) GetSimplifications(logFiles []string) []simplification {
	config := summarizerSvc.configSvc.Config
	simplifications := []simplification{}
	combinedLogsLoaded := summarizerSvc.combinedLogsSvc.IsLoaded()
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

		var result simplifier.Result
		if combinedLogsLoaded {
			result, err = summarizerSvc.simplifierSvc.ParseOutputFromCombinedLog(logFile, simplifier_)
			if err != nil {
				continue
			}
		} else {
			result, err = summarizerSvc.simplifierSvc.ParseOutputFromFile(path.Join(config.Paths.Logs, logFile), simplifier_)
			if err != nil {
				continue
			}
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

		writer.Write([]string{"Result", "Process Time", "Run Time", "Solver", "Instance Name"})
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
			runTime := fmt.Sprintf("%.3f", solution.runTime.Seconds())
			solver_ := solution.solver
			writer.Write([]string{result, processTime, runTime, string(solver_), name})
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

func printSolutionsStat(name string, solutions []solution, cubesets []cubeset, file *os.File) {
	satCount := 0
	satVerifiedCount := 0
	unsatCount := 0
	failsCount := 0
	solvedCount := 0
	totalProcessTime := time.Duration(0)
	totalRunTime := time.Duration(0)
	messages := []string{}
	processTimes := []float64{}

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
			messages = append(messages, solution.message)
		}

		if solution.result == solver.Sat || solution.result == solver.Unsat {
			totalProcessTime += solution.processTime
			processTimes = append(processTimes, solution.processTime.Seconds())
			solvedCount++
		}

		totalRunTime += solution.runTime
	}
	totalProcessTime = totalProcessTime.Round(time.Millisecond)
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
	solvedCount_ := humanize.Comma(int64(solvedCount))
	allCount_ := humanize.Comma(int64(solvedCount + failsCount))

	remaining := 0
	if len(cubesets) > 0 {
		remaining = cubesets[0].cubesCount - solvedCount
	} else {
		remaining = 1 - solvedCount
	}
	remaining_ := humanize.Comma(int64(remaining))

	file.WriteString(fmt.Sprintf("%s SAT%s, %s UNSAT, %s Fails, %s Remaining\n", satCount_, satVerifiedComment, unsatCount_, failCount_, remaining_))

	// Stats
	if len(processTimes) > 1 {
		sort.Float64s(processTimes)
		mean, stdDev := stat.MeanStdDev(processTimes, nil)
		median := stat.Quantile(0.5, stat.Empirical, processTimes, nil)
		lowest, highest := processTimes[0], processTimes[len(processTimes)-1]
		file.WriteString(fmt.Sprintf("Mean: %0.2fs, Median: %.2fs, Std. Deviation: %0.2f, Range: %.2fs to %.2fs\n\n", mean, median, stdDev, lowest, highest))
	}

	if solvedCount > 0 {
		file.WriteString(fmt.Sprintf("Process time (1 CPU, %s instances): %s\n", solvedCount_, totalProcessTime))
		file.WriteString(fmt.Sprintf("Process time (12 CPU, %s instances): %s\n", solvedCount_, time.Duration(totalProcessTime/12).Round(time.Millisecond)))
	}

	// TODO: Include runtimes
	file.WriteString(fmt.Sprintf("Run time (1 CPU, %s instances): %s\n", allCount_, totalRunTime))
	file.WriteString(fmt.Sprintf("Run time (12 CPU, %s instances): %s\n", allCount_, time.Duration(totalRunTime/12).Round(time.Millisecond)))

	if len(cubesets) > 0 && solvedCount > 0 {
		cubeset := cubesets[0]
		cubesCount := humanize.Comma(int64(cubeset.cubesCount))
		estimatedTime := time.Duration((int(totalProcessTime) / solvedCount) * cubeset.cubesCount).Round(time.Millisecond)
		estimatedTime12Cpu := time.Duration((int(totalProcessTime) / (solvedCount * 12)) * cubeset.cubesCount).Round(time.Millisecond)
		file.WriteString(fmt.Sprintf("\nEstimated time (1 CPU, %s instances): %s\n", cubesCount, estimatedTime))
		file.WriteString(fmt.Sprintf("Estimated time (12 CPU, %s instances): %s\n", cubesCount, estimatedTime12Cpu))
	}

	if len(messages) > 0 {
		file.WriteString("\nMessages:\n")
		for i, message := range messages {
			if len(message) != 128 {
				log.Fatalln("Message expected to be 128 hex. chars. long")
			}

			formattedMessage := strings.Builder{}
			for j := 0; j < 128; j += 8 {
				formattedMessage.WriteString(message[j:j+8] + " ")
			}

			file.WriteString(fmt.Sprintf("%d. %s\n", i+1, formattedMessage.String()))
		}
	}

	file.WriteString("\n")
}

func printCubesetsStat(cubesets []cubeset, file *os.File) {
	file.WriteString("## Cubesets\n\n")

	sort.Slice(cubesets, func(i, j int) bool {
		return cubesets[i].threshold < cubesets[j].threshold
	})
	for i, cubeset := range cubesets {
		cubesCount := humanize.Comma(int64(cubeset.cubesCount))
		refutedLeavesCount := humanize.Comma(int64(cubeset.refutedLeavesCount))
		file.WriteString(fmt.Sprintf("%d. n%d: %s cubes, %s refuted leaves, %s process time, %s\n", i+1, cubeset.threshold, cubesCount, refutedLeavesCount, cubeset.processTime, cubeset.name))
	}
	if len(cubesets) > 0 {
		file.WriteString("\n")
	}
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

	for i, baseEncodingName := range baseEncodingNames {
		file.WriteString(fmt.Sprintf("# %d. %s\n\n", i+1, baseEncodingName))
		combination := combinations[baseEncodingName]

		solutions := *combination.solutions
		solutionGroups := lo.GroupBy(solutions, func(s solution) string {
			name := regexp.MustCompile("(.cube[0-9]+)|(.log)").ReplaceAllString(s.name, "")
			return name
		})

		cubesets := *combination.cubesets
		cubesetGroups := lo.GroupBy(cubesets, func(c cubeset) string {
			return c.name
		})

		solutionNames := make([]string, 0, len(solutionGroups))
		for k := range solutionGroups {
			solutionNames = append(solutionNames, k)
		}
		sort.Slice(solutionNames, func(i, j int) bool {
			return solutionNames[i] > solutionNames[j]
		})
		for j, name := range solutionNames {
			solutions_ := solutionGroups[name]
			romanNumeral := integerToRoman(j + 1)
			file.WriteString(fmt.Sprintf("## %s. %s\n\n", romanNumeral, name))
			searchName_ := strings.Split(name, ".")
			searchName := strings.Join(searchName_[:len(searchName_)-2], ".")
			printSolutionsStat(name, solutions_, cubesetGroups[searchName], file)
		}

		if len(cubesets) > 0 {
			printCubesetsStat(cubesets, file)
		}
	}
}

func (summarizerSvc *SummarizerService) Run(workers int) {
	startTime := time.Now()

	solutionLogFiles := []string{}
	cubesetLogFiles := []string{}
	simplificationLogFiles := []string{}

	// Get the file entries
	fileEntries := []string{}
	if summarizerSvc.combinedLogsSvc.IsLoaded() {
		for entry := range summarizerSvc.combinedLogsSvc.LogFiles {
			fileEntries = append(fileEntries, entry)
		}
	} else {
		files, err := os.ReadDir(summarizerSvc.configSvc.Config.Paths.Logs)
		summarizerSvc.errorSvc.Fatal(err, "Summarizer: failed to find the log files")
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			fileName := file.Name()
			if path.Ext(fileName) != ".log" {
				continue
			}

			fileEntries = append(fileEntries, fileName)
		}
	}

	// Important: Register new SAT Solver here
	for _, fileEntry := range fileEntries {
		regexp_ := regexp.MustCompile(fmt.Sprintf("(%s.log)|(%s.log)|(%s.log)|(%s.log)|(%s.log)|(%s.log)|(%s.log)|(%s.log)|(%s.log)", solver.Kissat, simplifier.Cadical, solver.CryptoMiniSat, solver.Glucose, solver.MapleSat, solver.YalSat, solver.PalSat, solver.LSTechMaple, solver.KissatCF))
		if regexp_.Match([]byte(fileEntry)) {
			solutionLogFiles = append(solutionLogFiles, fileEntry)
		} else if strings.Contains(fileEntry, "march") {
			cubesetLogFiles = append(cubesetLogFiles, fileEntry)
		} else if strings.Contains(fileEntry, "satelite") || strings.Contains(fileEntry, "cadical_") {
			simplificationLogFiles = append(simplificationLogFiles, fileEntry)
		}
	}
	log.Printf("Processed %d items\n", len(fileEntries))

	cubesets := summarizerSvc.GetCubesets(cubesetLogFiles)
	summarizerSvc.WriteCubesetsLog(cubesets)
	log.Printf("Written summary for %d cubesets\n", len(cubesets))

	simplifications := summarizerSvc.GetSimplifications(simplificationLogFiles)
	summarizerSvc.WriteSimplificationsLog(simplifications)
	log.Printf("Written summary for %d simplifications\n", len(simplifications))

	solutions := summarizerSvc.GetSolutions(solutionLogFiles, workers)
	summarizerSvc.WriteSolutionsLog(solutions)
	log.Printf("Written summary for %d solutions\n", len(solutions))

	combinations := summarizerSvc.GetCombinations(solutions, cubesets)
	summarizerSvc.WriteCombinationsLog(combinations)
	log.Printf("Written summary for %d combinations\n", len(combinations))

	log.Println("Time taken:", time.Since(startTime))
}
