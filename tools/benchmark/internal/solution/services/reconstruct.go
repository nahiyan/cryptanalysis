package services

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitfield/script"
	"github.com/samber/lo"
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

func (solutionSvc *SolutionService) Reconstruct(solutionLiterals []int, reconstructionFilePath string, originalInstancePath string) ([]int, error) {
	// Get the RS literals
	rsLiterals := make(map[int]interface{}, 0)
	{
		reader := script.File(reconstructionFilePath).Reader
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			_, witness, err := scanRSLine(line)
			if err != nil {
				return []int{}, err
			}
			rsLiterals[witness] = nil
		}
	}

	// Pick the literals that are to be added to the original instance
	literalsToAdd := []int{}
	for _, solutionLiteral := range solutionLiterals {
		if _, exists := rsLiterals[solutionLiteral]; exists {
			continue
		}

		if _, exists := rsLiterals[-solutionLiteral]; exists {
			continue
		}

		literalsToAdd = append(literalsToAdd, solutionLiteral)
	}

	// Add the picked literals
	reader := script.File(originalInstancePath).Reader
	scanner := bufio.NewScanner(reader)
	header := ""
	var newInstance strings.Builder
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "c") {
			continue
		}

		if strings.HasPrefix(line, "p cnf") {
			header = line
			var numVars, numClauses int
			_, err := fmt.Sscanf(header, "p cnf %d %d", &numVars, &numClauses)
			if err != nil {
				return []int{}, err
			}
			newHeader := fmt.Sprintf("p cnf %d %d", numVars, numClauses+len(literalsToAdd))
			newInstance.WriteString(newHeader + "\n")
			continue
		}

		newInstance.WriteString(line + "\n")
	}

	for _, literal := range literalsToAdd {
		newInstance.WriteString(fmt.Sprintf("%d 0\n", literal))
	}

	script.Echo(newInstance.String()).WriteFile("/tmp/test.cnf")

	// Use Kissat for solving the new instance
	kissatCmd := fmt.Sprintf("%s -q", solutionSvc.configSvc.Config.Paths.Bin.Kissat)
	solution_, err := script.Echo(newInstance.String()).Exec(kissatCmd).ReplaceRegexp(regexp.MustCompile("(s SATISFIABLE)|(v)"), "").String()
	if err != nil && err.Error() != "exit status 10" {
		return []int{}, err
	}

	reconstructedSolutionLiterals := lo.Map(strings.Fields(solution_), func(item string, _ int) int {
		literal, _ := strconv.Atoi(item)
		return literal
	})

	return reconstructedSolutionLiterals, nil
}
