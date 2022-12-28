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

func (solutionSvc *SolutionService) Find(encoding string) (solver.Solution, error) {
	filesystemSvc := solutionSvc.filesystemSvc
	databaseSvc := solutionSvc.databaseSvc

	checksum, err := filesystemSvc.Checksum(encoding)
	if err != nil {
		return solver.Solution{}, err
	}

	data, err := databaseSvc.Get(solutionSvc.Bucket, checksum)
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

func (solutionSvc *SolutionService) Register(encoding string, solution solver.Solution) error {
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

	if err := databaseSvc.Set(solutionSvc.Bucket, checksum, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}
