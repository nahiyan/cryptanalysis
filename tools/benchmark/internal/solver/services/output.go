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

func (solverSvc *SolverService) ParseLog(logPath string, solver_ solver.Solver, solutionLiterals *[]int) (solver.Result, time.Duration, time.Duration, error) {
	switch solver_ {
	case solver.Kissat:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c process-time:", 1, solutionLiterals)
	case solver.Cadical:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c total process time since initialization:", 1, solutionLiterals)
	case solver.Glucose:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c CPU time", 1, solutionLiterals)
	case solver.MapleSat:
		return parseOutputWith(logPath, "SATISFIABLE", "UNSATISFIABLE", "CPU time", 1, solutionLiterals)
	case solver.YalSat:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c total process time of", 1, solutionLiterals)
	case solver.PalSat:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c total wall clock time", 1, solutionLiterals)
	case solver.CryptoMiniSat:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c Total time", 0, solutionLiterals)
	}

	return solver.Fail, time.Duration(0), time.Duration(0), nil
}

func parseOutputWith(logPath, satText, unsatText, processTimeText string, processTimeFieldOffset int, solutionLiterals *[]int) (solver.Result, time.Duration, time.Duration, error) {
	processTime := time.Duration(0)
	runTime := time.Duration(0)
	result := solver.Result(solver.Fail)

	regexp_ := regexp.MustCompile(fmt.Sprintf("(%s)|(%s)|(%s)|(Info: Ended after)", satText, unsatText, processTimeText))
	matches, err := script.File(logPath).MatchRegexp(regexp_).Slice()
	if err != nil {
		return result, processTime, runTime, err
	}

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
		runTime_, err := strconv.ParseFloat(pieces[len(pieces)-2], 64)
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
	if solutionLiterals != nil {
		// TODO: Requires regexp improvement
		regexp_ := regexp.MustCompile("(v )|(-1 )|(1 )")
		lines, err := script.File(logPath).MatchRegexp(regexp_).Slice()
		if err != nil {
			return result, processTime, runTime, err
		}

		for _, line := range lines {
			if !strings.HasPrefix(line, "v ") && !strings.HasPrefix(line, "-1 ") && !strings.HasPrefix(line, "1 ") {
				continue
			}
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
