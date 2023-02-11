package services

import (
	"benchmark/internal/simplifier"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/samber/lo"
)

func (solverSvc *SimplifierService) ParseOutput(logFilePath string, simplifier_ simplifier.Simplifier) (simplifier.Result, error) {
	numVarsBefore := 0
	numVarsAfter := 0
	numClausesBefore := 0
	numClausesAfter := 0
	seconds := 0.0
	result := simplifier.Result{}

	pipe := script.File(logFilePath)
	switch simplifier_ {
	case simplifier.Satelite:
		output, err := pipe.MatchRegexp(regexp.MustCompile("(Result)|(CPU time)|(UNSATISFIABLE)|(SATISFIABLE)")).Slice()
		// Catch if it's solved by the simplifier
		{
			_, isSolved := lo.Find(output, func(item string) bool {
				return item == "UNSATISFIABLE" || item == "SATISFIABLE"
			})
			if isSolved {
				return result, errors.New("instance found to be solved by SatELite")
			}
		}

		// Initialize the parse
		if err != nil {
			return result, err
		}
		if len(output) != 2 {
			return result, errors.New("invalid SatELite output")
		}

		// Parse the number of new variables and clauses
		{
			fields := strings.Fields(output[0])
			if len(fields) < 8 {
				return result, errors.New("invalid SatELite output")
			}

			numVarsAfter, err = strconv.Atoi(fields[3])
			if err != nil {
				return result, err
			}

			numClausesAfter, err = strconv.Atoi(fields[5])
			if err != nil {
				return result, err
			}
		}

		// Parse the process time
		fields := strings.Fields(output[1])
		if len(fields) != 4 {
			return result, errors.New("invalid SatELite output")
		}

		seconds, err = strconv.ParseFloat(fields[len(fields)-2], 64)
		if err != nil {
			return result, err
		}
	case simplifier.Cadical:
		// See if it's solved by the simplifier
		output, err := pipe.MatchRegexp(regexp.MustCompile("(c writing 'p cnf)|(c total process time since initialization:)|(c exit 0)")).Slice()
		{
			_, hasSimplified := lo.Find(output, func(item string) bool {
				return item == "c exit 0"
			})

			if !hasSimplified {
				return result, errors.New("instance not found to be simplified, but probably solved")
			}
		}

		// Initialize the parse
		if err != nil {
			return result, err
		}
		if len(output) != 3 {
			return result, errors.New("invalid CaDiCaL output")
		}

		// Parse the number of new variables and clauses
		{
			_, err := fmt.Sscanf(output[0], "c writing 'p cnf %d %d' header", &numVarsAfter, &numClausesAfter)
			if err != nil {
				return result, errors.New("invalid CaDiCaL output")
			}
		}

		// Parse the process time
		fields := strings.Fields(output[1])
		seconds, err = strconv.ParseFloat(fields[len(fields)-2], 64)
		if err != nil {
			return result, err
		}
	}

	// Get the numVariables and numClauses of the original instance
	{
		segments := strings.Split(logFilePath, ".")
		segments = segments[:len(segments)-2]
		originalInstanceName := path.Base(strings.Join(segments, "."))
		originalInstancePath := path.Join(solverSvc.configSvc.Config.Paths.Encodings, originalInstanceName)
		header, err := script.File(originalInstancePath).Match("p cnf").String()
		if err != nil {
			return result, err
		}

		_, err = fmt.Sscanf(header, "p cnf %d %d", &numVarsBefore, &numClausesBefore)
		if err != nil {
			return result, err
		}
	}

	result.ProcessTime = time.Duration(seconds*1000) * time.Millisecond
	// TODO: Calculate the free variables after simplification
	result.NumVars = numVarsAfter
	result.NumClauses = numClausesAfter
	result.NumEliminatedVars = numVarsBefore - numVarsAfter
	result.NumEliminatedClauses = numClausesBefore - numClausesAfter

	return result, nil
}
