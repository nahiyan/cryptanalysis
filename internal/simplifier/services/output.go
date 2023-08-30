package services

import (
	"cryptanalysis/internal/simplifier"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitfield/script"
	"github.com/samber/lo"
)

func (solverSvc *SimplifierService) ParseOutputFromFile(logFilePath string, simplifier_ simplifier.Simplifier) (simplifier.Result, error) {
	file, err := os.OpenFile(logFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return simplifier.Result{}, err
	}
	return solverSvc.ParseOutput(file, logFilePath, simplifier_)
}

func (solverSvc *SimplifierService) ParseOutputFromCombinedLog(logFilePath string, simplifier_ simplifier.Simplifier) (simplifier.Result, error) {
	maybeContent, err := solverSvc.combinedLogsSvc.Get(logFilePath)
	if err != nil {
		return simplifier.Result{}, err
	}

	content, exists := maybeContent.Get()
	log.Println(exists)
	if !exists {
		return simplifier.Result{}, os.ErrNotExist
	}

	reader := strings.NewReader(content)
	return solverSvc.ParseOutput(reader, logFilePath, simplifier_)
}

func (solverSvc *SimplifierService) ParseOutput(outputReader io.Reader, logFilePath string, simplifier_ simplifier.Simplifier) (simplifier.Result, error) {
	numVarsBefore := 0
	numVarsAfter := 0
	numClausesBefore := 0
	numClausesAfter := 0
	seconds := 0.0
	result := simplifier.Result{}

	output := ""
	{
		buf := new(strings.Builder)
		_, err := io.Copy(buf, outputReader)
		if err != nil {
			return simplifier.Result{}, err
		}
		output = buf.String()
	}

	switch simplifier_ {
	case simplifier.Cadical:
		// See if it's solved by the simplifier
		matches := regexp.MustCompile(`(c writing 'p cnf [0-9]+ [0-9]+' header)|(c total process time since initialization:\s+[0-9.]+)|(c exit 0)`).FindAllString(output, len(output))
		{
			_, hasSimplified := lo.Find(matches, func(item string) bool {
				return item == "c exit 0"
			})

			if !hasSimplified {
				return result, errors.New("instance not found to be simplified, but probably solved")
			}
		}

		// Initialize the parse
		if len(matches) != 3 {
			return result, errors.New("invalid CaDiCaL output")
		}

		// Parse the number of new variables and clauses
		{
			_, err := fmt.Sscanf(matches[0], "c writing 'p cnf %d %d' header", &numVarsAfter, &numClausesAfter)
			if err != nil {
				return result, errors.New("invalid CaDiCaL output")
			}
		}

		// Parse the process time
		var err error
		fields := strings.Fields(matches[1])
		seconds, err = strconv.ParseFloat(fields[len(fields)-1], 64)
		if err != nil {
			return result, err
		}
	}

	// TODO: Support original instances existing only in memory
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
