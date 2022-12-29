package services

import (
	"benchmark/internal/solver"
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
