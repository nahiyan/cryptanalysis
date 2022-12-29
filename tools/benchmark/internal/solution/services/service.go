package services

import (
	"benchmark/internal/solver"
	"bytes"
	"encoding/gob"
)

// TODO: Should use a repository for DB operations

type Properties struct {
	Bucket string
}

func (solutionSvc *SolutionService) Init() {
	solutionSvc.Bucket = "solutions"
}

func (solutionSvc *SolutionService) Find(encoding string, solver_ string) (solver.Solution, error) {
	filesystemSvc := solutionSvc.filesystemSvc
	databaseSvc := solutionSvc.databaseSvc

	checksum, err := filesystemSvc.Checksum(encoding)
	if err != nil {
		return solver.Solution{}, err
	}

	key := checksum
	key = append(key, []byte("_"+solver_)...)
	data, err := databaseSvc.Get(solutionSvc.Bucket, key)
	if err != nil {
		return solver.Solution{}, err
	}

	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	solution := solver.Solution{}
	if err := decoder.Decode(&solution); err != nil {
		return solver.Solution{}, err
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

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(solution); err != nil {
		return err
	}

	key := checksum
	key = append(key, []byte("_"+solver_)...)

	if err := databaseSvc.Set(solutionSvc.Bucket, key, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}
