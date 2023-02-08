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

func (solverSvc *SolverService) ParseLog(logPath string, solver_ solver.Solver, extractSolution bool) (solver.Result, time.Duration, error) {
	// fmt.Println(solver_)

	switch solver_ {
	case solver.Kissat:
		return parseOutputWith(logPath, "c exit 10", "c exit 20", "c process-time:", 1)
	case solver.Cadical:
		return parseOutputWith(logPath, "c exit 10", "c exit 20", "c total process time since initialization:", 1)
	case solver.Glucose:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c CPU time", 1)
	case solver.MapleSat:
		return parseOutputWith(logPath, "SATISFIABLE", "UNSATISFIABLE", "CPU time", 1)
	case solver.CryptoMiniSat:
		return parseOutputWith(logPath, "s SATISFIABLE", "s UNSATISFIABLE", "c time", 0)
	}

	return solver.Fail, time.Duration(0), nil
}

func parseOutputWith(logPath, satText, unsatText, processTimeText string, processTimeFieldOffset int) (solver.Result, time.Duration, error) {
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

	// fmt.Println(processTime, result, seconds)

	return result, processTime, nil
}
