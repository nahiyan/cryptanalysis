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

func isLiteralInRanges(literal int, ranges []solution.Range) bool {
	variable := literalToVariable(literal)
	for _, range_ := range ranges {
		if variable >= range_.Start && variable <= range_.End {
			return true
		}
	}

	return false
}

func (encodingSvc *SimplificationService) Reconstruct(instancePath, reconstructionPath string, ranges []solution.Range) error {
	removedClauses := [][]int{}
	{
		reader := script.File(reconstructionPath).Reader
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			clause, witness, err := scanRSLine(line)
			if err != nil {
				return err
			}

			if !isLiteralInRanges(witness, ranges) {
				continue
			}

			removedClauses = append(removedClauses, clause)
		}
	}

	{
		newInstance := ""
		reader := script.File(instancePath).Reader
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "c") {
				continue
			}

			// Header
			if strings.HasPrefix(line, "p cnf") {
				varCount := 0
				clausesCount := 0
				fmt.Sscanf(line, "p cnf %d %d", &varCount, &clausesCount)

				newClausesCount := clausesCount + len(removedClauses)
				newInstance += fmt.Sprintf("p cnf %d %d\n", varCount, newClausesCount)
				continue
			}

			newInstance += line + "\n"
		}

		for _, clause := range removedClauses {
			clause_ := ""
			for _, literal := range clause {
				clause_ += strconv.Itoa(literal) + " "
			}
			newInstance += clause_ + "0\n"
		}

		_, err := script.Echo(newInstance).WriteFile(instancePath)
		if err != nil {
			return err
		}
	}

	return nil
}
