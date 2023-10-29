package services

import (
	"bufio"
	"os"
	"strings"
)

func (solutionSvc *SolutionService) Normalize(solutionPath string) error {
	instanceFile, err := os.Open(solutionPath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(instanceFile)
	newBody := ""
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimPrefix(line, "v ")

		if strings.HasPrefix(line, "s SATISFIABLE") || strings.HasPrefix(line, "SAT") {
			continue
		}

		newBody += line + " "
	}
	instanceFile.Close()

	outputPath, err := os.OpenFile(solutionPath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = outputPath.WriteString("SAT" + "\n" + newBody + "\n")
	if err != nil {
		return err
	}

	return nil
}
