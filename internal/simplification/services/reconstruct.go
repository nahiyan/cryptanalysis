package services

import (
	"bufio"
	"fmt"
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

func (encodingSvc *SimplificationService) Reconstruct(instancePath, reconstructionPath string) error {
	clausesToPreserve := [][]int{}
	{
		reader := script.File(reconstructionPath).Reader
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			clause, _, err := scanRSLine(line)
			if err != nil {
				return err
			}

			clausesToPreserve = append(clausesToPreserve, clause)
		}
	}

	{
		var newInstanceBuilder strings.Builder
		reader := script.File(instancePath).Reader
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			// Header
			if strings.HasPrefix(line, "p cnf") {
				varCount := 0
				clausesCount := 0
				fmt.Sscanf(line, "p cnf %d %d", &varCount, &clausesCount)

				newClausesCount := clausesCount + len(clausesToPreserve)
				newInstanceBuilder.WriteString(fmt.Sprintf("p cnf %d %d\n", varCount, newClausesCount))
				continue
			}

			newInstanceBuilder.WriteString(line + "\n")
		}

		for _, clause := range clausesToPreserve {
			clause_ := ""
			for _, literal := range clause {
				clause_ += strconv.Itoa(literal) + " "
			}
			newInstanceBuilder.WriteString(clause_ + "0\n")
		}

		_, err := script.Echo(newInstanceBuilder.String()).WriteFile(instancePath)
		if err != nil {
			return err
		}
	}

	return nil
}
