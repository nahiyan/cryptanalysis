package services

import (
	"benchmark/internal/encoder"
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/bitfield/script"
	"github.com/samber/mo"
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

func (solutionSvc *SolutionService) Reconstruct(solutionPath string, reconstructionFilePath string, instancePath string) error {
	info, err := solutionSvc.encoderSvc.ProcessInstanceName(instancePath)
	if err != nil {
		return err
	}
	info.Simplification = mo.None[encoder.SimplificationInfo]()
	originalInstance := solutionSvc.encoderSvc.GetInstanceName(info)

	// Get the RS literals
	rsLiterals := make(map[int]interface{}, 0)
	{
		reader := script.File(reconstructionFilePath).Reader
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			_, witness, err := scanRSLine(line)
			if err != nil {
				return err
			}
			rsLiterals[witness] = nil
		}
	}

	// Get the solution literals
	solutionLiterals := []int{}
	{
		reader := script.File(solutionPath).Reader
		scanner := bufio.NewScanner(reader)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			word := scanner.Text()
			if word == "SAT" || word == "0" {
				continue
			}

			literal, err := strconv.Atoi(word)
			if err != nil {
				return err
			}
			solutionLiterals = append(solutionLiterals, literal)
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
	reader := script.File(originalInstance).Reader
	scanner := bufio.NewScanner(reader)
	header := ""
	newInstance := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "c") {
			continue
		}

		if strings.HasPrefix(line, "p cnf") {
			header = line
			continue
		}

		newInstance += line + "\n"
	}

	var numVars, numClauses int
	_, err = fmt.Sscanf(header, "p cnf %d %d", &numVars, &numClauses)
	if err != nil {
		return err
	}
	newHeader := fmt.Sprintf("p cnf %d %d", numVars, numClauses+len(literalsToAdd))
	newInstance = newHeader + "\n" + newInstance
	for _, literal := range literalsToAdd {
		newInstance += fmt.Sprintf("%d 0\n", literal)
	}

	// Use Kissat for solving the new instance
	kissatCmd := fmt.Sprintf("%s -q", solutionSvc.configSvc.Config.Paths.Bin.Kissat)
	_, err = script.Echo(newInstance).Exec(kissatCmd).Replace("v ", "").Replace("s SATISFIABLE", "SAT").WriteFile(solutionPath)
	return err
}
