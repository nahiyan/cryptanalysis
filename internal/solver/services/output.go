package services

import (
	"cryptanalysis/internal/solver"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (solverSvc *SolverService) ParseLogFromFile(logPath string, solver_ solver.Solver, solutionLiterals *[]int) (solver.Result, time.Duration, time.Duration, error) {
	file, err := os.OpenFile(logPath, os.O_RDONLY, 0644)
	if err != nil {
		return solver.Fail, time.Duration(0), time.Duration(0), err
	}
	return solverSvc.ParseLog(file, solver_, solutionLiterals)
}

func (solverSvc *SolverService) ParseLogFromCombinedLog(logFilePath string, solver_ solver.Solver, solutionLiterals *[]int) (solver.Result, time.Duration, time.Duration, error) {
	maybeContent, err := solverSvc.combinedLogsSvc.Get(logFilePath)
	if err != nil {
		return solver.Fail, time.Duration(0), time.Duration(0), err
	}

	content, exists := maybeContent.Get()
	if !exists {
		return solver.Fail, time.Duration(0), time.Duration(0), os.ErrNotExist
	}

	reader := strings.NewReader(content)
	return solverSvc.ParseLog(reader, solver_, solutionLiterals)
}

// Important: Register new SAT Solver here
func (solverSvc *SolverService) ParseLog(outputReader io.Reader, solver_ solver.Solver, solutionLiterals *[]int) (solver.Result, time.Duration, time.Duration, error) {
	switch solver_ {
	case solver.Kissat:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c process-time:", 1, solutionLiterals)
	case solver.Cadical:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c total process time since initialization:", 1, solutionLiterals)
	case solver.Glucose:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c CPU time", 1, solutionLiterals)
	case solver.MapleSat:
		return parseOutputWith(outputReader, "SATISFIABLE", "UNSATISFIABLE", "CPU time", 1, solutionLiterals)
	case solver.YalSat:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c total process time of", 1, solutionLiterals)
	case solver.PalSat:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c total wall clock time", 1, solutionLiterals)
	case solver.CryptoMiniSat:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c Total time", 0, solutionLiterals)
	case solver.LSTechMaple:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c CPU time", 1, solutionLiterals)
	case solver.KissatCF:
		return parseOutputWith(outputReader, "s SATISFIABLE", "s UNSATISFIABLE", "c process-time:", 1, solutionLiterals)
	}

	return solver.Fail, time.Duration(0), time.Duration(0), nil
}

func parseOutputWith(outputReader io.Reader, satText, unsatText, processTimeText string, processTimeFieldOffset int, solutionLiterals *[]int) (solver.Result, time.Duration, time.Duration, error) {
	processTime := time.Duration(0)
	runTime := time.Duration(0)
	result := solver.Result(solver.Fail)

	output := ""
	{
		buf := new(strings.Builder)
		_, err := io.Copy(buf, outputReader)
		if err != nil {
			return result, processTime, runTime, err
		}
		output = buf.String()
	}

	matches := regexp.MustCompile(fmt.Sprintf("(.*%s.*)|(.*%s.*)|(.*%s.*)|(Info: Ended after.*)", satText, unsatText, processTimeText)).FindAllString(output, len(output))

	// If none of the expected matches are met
	if len(matches) == 0 {
		return result, processTime, runTime, errors.New("invalid solver output format")
	}

	for _, match := range matches {
		// Process the result
		switch match {
		case satText:
			result = solver.Sat
		case unsatText:
			result = solver.Unsat
		}

		// Process the time it took to process the instance
		if strings.Contains(match, processTimeText) {
			processTimeFields := strings.Fields(match)
			seconds, _ := strconv.ParseFloat(processTimeFields[len(processTimeFields)-(processTimeFieldOffset+1)], 64)
			processTime = time.Duration(seconds*1000) * time.Millisecond
			continue
		}

		// Parse runtime
		if !strings.HasPrefix(match, "Info: Ended after") {
			continue
		}

		pieces := strings.Fields(match)
		runTime_, err := strconv.ParseFloat(pieces[3], 64)
		if err != nil {
			return result, processTime, runTime, errors.New("failed to parse runtime")
		}
		runTime = time.Duration(runTime_*1000) * time.Millisecond
		break
	}

	// See if it failed
	if len(matches) != 3 {
		return result, processTime, runTime, nil
	}

	// Extract the solution literals
	if solutionLiterals != nil && result == solver.Sat {
		// TODO: Improve regexp
		lines := regexp.MustCompile(`(?m)(^v.*)|(?m)(^-1\s.*)|(?m)(^1\s.*)`).FindAllString(output, len(output))
		for _, line := range lines {
			segments := strings.Fields(line)
			for _, segment := range segments {
				if segment == "0" || segment == "v" {
					continue
				}

				literal, err := strconv.Atoi(segment)
				if err != nil {
					return result, processTime, runTime, err
				}
				*solutionLiterals = append(*solutionLiterals, literal)
			}
		}
	}

	return result, processTime, runTime, nil
}
