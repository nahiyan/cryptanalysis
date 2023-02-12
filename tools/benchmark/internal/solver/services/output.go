package services

import (
	"benchmark/internal/solver"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitfield/script"
)

func (solverSvc *SolverService) ParseLog(logPath string, solver_ solver.Solver, solutionLiterals *[]int) (solver.Result, time.Duration, error) {
	switch solver_ {
	case solver.Kissat:
		return parseOutputWith(logPath, "c exit 10", "c exit 20", "c process-time:", 1, solutionLiterals)
	case solver.Cadical:
		return parseOutputWith(logPath, "c exit 10", "c exit 20", "c total process time since initialization:", 1, solutionLiterals)
	case solver.Glucose:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c CPU time", 1, solutionLiterals)
	case solver.MapleSat:
		return parseOutputWith(logPath, "SATISFIABLE", "UNSATISFIABLE", "CPU time", 1, solutionLiterals)
	case solver.CryptoMiniSat:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c time", 0, solutionLiterals)
	}

	return solver.Fail, time.Duration(0), nil
}

func parseOutputWith(logPath, satText, unsatText, processTimeText string, processTimeFieldOffset int, solutionLiterals *[]int) (solver.Result, time.Duration, error) {
	processTime := time.Duration(0)
	result := solver.Result(solver.Fail)

	regexp := regexp.MustCompile(fmt.Sprintf("(%s)|(%s)|(%s)", processTimeText, satText, unsatText))

	output, err := script.File(logPath).MatchRegexp(regexp).Slice()
	if err != nil {
		return result, processTime, err
	}

	if len(output) != 2 {
		return result, processTime, errors.New("invalid solver output format")
	}

	// Process the time it took to process the instance
	processTimeFields := strings.Fields(output[0])
	seconds, _ := strconv.ParseFloat(processTimeFields[len(processTimeFields)-(processTimeFieldOffset+1)], 64)
	processTime = time.Duration(seconds*1000) * time.Millisecond

	// Process the result
	switch output[1] {
	case satText:
		result = solver.Sat
	case unsatText:
		result = solver.Unsat
	default:
		result = solver.Fail
	}

	// Extract the solution literals
	if solutionLiterals != nil {
		lines, err := script.File(logPath).Match("v ").Slice()
		if err != nil {
			return result, processTime, err
		}

		for _, line := range lines {
			if !strings.HasPrefix(line, "v") {
				continue
			}
			segments := strings.Fields(line)
			for _, segment := range segments {
				if segment == "0" || segment == "v" {
					continue
				}

				literal, err := strconv.Atoi(segment)
				if err != nil {
					return result, processTime, err
				}
				*solutionLiterals = append(*solutionLiterals, literal)
			}
		}
	}

	return result, processTime, nil
}