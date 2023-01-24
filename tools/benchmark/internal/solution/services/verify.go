package services

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func (solutionSvc *SolutionService) Verify(solution io.Reader, steps int) (bool, error) {
	command := fmt.Sprintf("%s %d", solutionSvc.configSvc.Config.Paths.Bin.Verifier, steps)
	cmd := solutionSvc.commandSvc.Create(command)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	inPipe, err := cmd.StdinPipe()
	if err != nil {
		return false, err
	}

	if err := cmd.Start(); err != nil {
		return false, err
	}

	_, err = io.Copy(inPipe, solution)
	if err != nil {
		return false, err
	}
	inPipe.Close()

	scanner := bufio.NewScanner(outPipe)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Solution's hash matches the target!") {
			return true, nil
		} else if strings.Contains(line, "Solution's hash DOES NOT match the target:") || strings.Contains(line, "Result is UNSAT!") {
			return false, nil
		}
	}

	err = cmd.Wait()
	if err != nil {
		return false, err
	}

	return false, nil
}
