package services

import (
	"benchmark/internal/solver"
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

func (solutionSvc *SolutionService) RemapAndVerify(solutionPath string, varMapPath string) error {
	// info, err := solutionSvc.encoderSvc.ProcessInstanceName(solutionPath)
	// if err != nil {
	// 	return err
	// }
	// steps := info.Steps

	// // * 1. Read the reconstruction stack file and determine the literals that need to be preserved
	// varMapFile, err := os.OpenFile(varMapPath, os.O_RDONLY, 0600)
	// if err != nil {
	// 	return err
	// }
	// defer varMapFile.Close()

	// replacements := make(map[int]bool)
	// scanner := bufio.NewScanner(varMapFile)

	// TODO: Implement remapping and verification
	return nil
}
