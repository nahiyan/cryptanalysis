package services

import (
	"benchmark/internal/solution"
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/bitfield/script"
)

func scanRSLine(line string) ([]int, int, error) {
	segments := strings.Fields(line)
	clause := []int{}
	witness := 0

	isClause := true
	for _, segment := range segments {
		if segment == "0" {
			isClause = false
			continue
		}

		var err error
		if isClause {
			literal, err := strconv.Atoi(segment)
			if err != nil {
				return clause, witness, err
			}

			clause = append(clause, literal)
			continue
		}

		witness, err = strconv.Atoi(segment)
		if err != nil {
			return clause, witness, err
		}
	}

	return clause, witness, nil
}

func literalToVariable(literal int) int {
	variable := int(math.Abs(float64(literal)))
	return variable
}

func isVariableInRanges(literal int, ranges []solution.Range) bool {
	variable := literalToVariable(literal)
	for _, range_ := range ranges {
		if variable >= range_.Start && variable <= range_.End {
			return true
		}
	}

	return false
}

func countVariables(clauses [][]int) int {
	count := 0
	for _, clause := range clauses {
		for _, literal := range clause {
			variable := literalToVariable(literal)
			if variable > count {
				count = variable
			}
		}
	}
	return count
}

func (solutionSvc *SolutionService) ReconstructSolution(solutionPath string, reconstructionFilePath string, ranges []solution.Range) error {
	reader := script.File(reconstructionFilePath).Reader
	scanner := bufio.NewScanner(reader)
	clauses := [][]int{}
	for scanner.Scan() {
		line := scanner.Text()
		clause, witness, err := scanRSLine(line)
		if err != nil {
			return err
		}

		// Ignore variables that we're not interested in
		if !isVariableInRanges(witness, ranges) {
			continue
		}

		clauses = append(clauses, clause)
	}

	{
		reader := script.File(solutionPath).Reader
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "SAT" {
				continue
			}

			literals := strings.Fields(line)
			for _, literal := range literals {
				literal_, err := strconv.Atoi(literal)
				if err != nil {
					return err
				}
				if literal_ == 0 {
					continue
				}

				if isVariableInRanges(literal_, ranges) {
					continue
				}

				clauses = append(clauses, []int{literal_})
			}
		}
	}

	numClauses := len(clauses)
	numVariables := countVariables(clauses)
	instance := fmt.Sprintf("p cnf %d %d\n", numVariables, numClauses)
	for _, clause := range clauses {
		for _, literal := range clause {
			instance += fmt.Sprintf("%d ", literal)
		}
		instance += "0 \n"
	}

	_, err := script.Echo(instance).WriteFile("/tmp/test.cnf")
	if err != nil {
		return err
	}

	_, err = script.File("/tmp/test.cnf").Exec("kissat -q").WriteFile(solutionPath)
	if err != nil {
		return err
	}

	if err := solutionSvc.Normalize(solutionPath); err != nil {
		return err
	}

	return nil
}
