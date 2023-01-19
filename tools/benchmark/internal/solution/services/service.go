package services

import (
	"benchmark/internal/solver"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO: Should use a repository for DB operations

type Properties struct {
	Bucket string
}

func (solutionSvc *SolutionService) Init() {
	solutionSvc.Bucket = "solutions"
}

func (solutionSvc *SolutionService) Find(encoding string, solver_ solver.Solver) (solver.Solution, error) {
	filesystemSvc := solutionSvc.filesystemSvc
	databaseSvc := solutionSvc.databaseSvc

	checksum, err := filesystemSvc.Checksum(encoding)
	if err != nil {
		return solver.Solution{}, err
	}

	key := []byte(checksum + "_" + string(solver_))
	data, err := databaseSvc.Get(solutionSvc.Bucket, key)
	if err != nil {
		return solver.Solution{}, err
	}

	solution := solver.Solution{}
	if err := solutionSvc.marshallingSvc.BinDecode(data, &solution); err != nil {
		return solution, err
	}

	return solution, nil
}

func (solutionSvc *SolutionService) Register(encoding string, solver_ solver.Solver, solution solver.Solution) error {
	startTime := time.Now()
	defer solutionSvc.filesystemSvc.LogInfo("Solution: register took", time.Since(startTime).String())

	databaseSvc := solutionSvc.databaseSvc
	filesystemSvc := solutionSvc.filesystemSvc

	checksum, err := filesystemSvc.Checksum(encoding)
	if err != nil {
		return err
	}

	value, err := solutionSvc.marshallingSvc.BinEncode(solution)
	if err != nil {
		return err
	}

	key := []byte(checksum + "_" + string(solver_))
	if err := databaseSvc.Set(solutionSvc.Bucket, key, value); err != nil {
		return err
	}

	return nil
}

func (solutionSvc *SolutionService) All() ([]solver.Solution, error) {
	solutions := []solver.Solution{}
	solutionSvc.databaseSvc.All(solutionSvc.Bucket, func(key, value []byte) {
		var solution solver.Solution
		if err := solutionSvc.marshallingSvc.BinDecode(value, &solution); err != nil {
			return
		}

		solutions = append(solutions, solution)
	})

	return solutions, nil
}

func (solutionSvc *SolutionService) Normalize(encodingPath string) error {
	instanceFile, err := os.Open(encodingPath)
	if err != nil {
		return err
	}

	literals := []int{}
	scanner := bufio.NewScanner(instanceFile)
	header := ""
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimPrefix(line, "v ")

		if strings.HasPrefix(line, "p cnf") {
			header = line
		}

		if strings.HasPrefix(line, "c ") {
			continue
		}

		segments := strings.Fields(line)
		for _, segment := range segments {
			literal, err := strconv.Atoi(segment)
			if err != nil {
				return err
			}

			if literal == 0 {
				continue
			}

			literals = append(literals, literal)
		}
	}
	instanceFile.Close()

	outputPath, err := os.OpenFile(encodingPath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	line := ""
	for _, literal := range literals {
		line = fmt.Sprintf("%d ", literal)
	}
	_, err = outputPath.WriteString(header + "\n" + line + " 0\n")
	if err != nil {
		return err
	}

	return nil
}
