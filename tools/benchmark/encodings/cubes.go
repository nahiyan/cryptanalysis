package encodings

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/utils"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func GenerateSubProblem(instanceName string, i int) error {
	// Read the instance
	instance_, err := os.ReadFile(fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, instanceName))
	if err != nil {
		return err
	}
	instance := string(instance_)

	// Open the .icnf file
	iCnfFile, err := os.Open(fmt.Sprintf("%s%s.icnf", constants.EncodingsDirPath, instanceName))
	if err != nil {
		return err
	}
	defer iCnfFile.Close()

	// Generate the subproblem
	cube, _, _ := utils.ReadLine(iCnfFile, i)

	// Generate the clause
	clause := strings.TrimPrefix(cube, "a ")

	// Grab the CNF as a starting point for the cube
	head, body, _ := strings.Cut(instance, "\n")
	headerSegments := strings.Split(head, " ")
	numVars, _ := strconv.Atoi(headerSegments[2])
	numClauses, _ := strconv.Atoi(headerSegments[3])

	// Generate a new header with an incremented number of clauses
	clauseSegments := strings.Split(clause, " ")
	newNumVar := numVars
	for _, clauseSegment := range clauseSegments {
		var_, _ := strconv.Atoi(clauseSegment)
		newNumVar = int(math.Max(float64(var_), float64(numVars)))
	}
	newHeader := fmt.Sprintf("p cnf %d %d", newNumVar, numClauses+1)

	// Write the cube
	cubeFileName := fmt.Sprintf("%scube%d_%s.cnf", constants.EncodingsDirPath, i, instanceName)
	cubeFile, err := os.OpenFile(cubeFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New("failed to create the cubes file")
	}
	defer cubeFile.Close()

	_, err = cubeFile.WriteString(fmt.Sprintf("%s\n%s%s\n", newHeader, body, clause))
	if err != nil {
		return errors.New("failed to add the clause to the cube file")
	}

	return nil
}

func generateCubes(instanceName string, cubeVars uint) {
	// Invoke March for generating the cubes that is held in an .icnf file
	core.March(fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, instanceName), cubeVars)
}
