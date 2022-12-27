package services

import "fmt"

func (solverSvc *SolverService) Run(encodings []string) {
	fmt.Println("Solve", encodings)
}
