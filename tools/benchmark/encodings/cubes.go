package encodings

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/utils"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"
)

func GenerateSubProblemAsStringWithThreshold(instanceName string, i int, threshold *uint) (string, error) {
	subProblemFilePath := fmt.Sprintf("%scube%d_%s.cnf", constants.EncodingsDirPath, i, instanceName)
	if utils.FileExists(subProblemFilePath) {
		return "", nil
	}

	// * 1. Read the instance
	instance_, err := os.ReadFile(fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, instanceName))
	if err != nil {
		return "", err
	}
	instance := string(instance_)

	// * 2. Open the cubes file
	cubesFilePath := func(threshold *uint) string {
		if threshold == nil {
			return fmt.Sprintf("%s%s.icnf", constants.EncodingsDirPath, instanceName)
		}
		return fmt.Sprintf("%sn%d_%s.icnf", constants.EncodingsDirPath, *threshold, instanceName)
	}(threshold)
	cubesFile, err := os.Open(cubesFilePath)
	if err != nil {
		return "", err
	}
	defer cubesFile.Close()

	// * 3. Get the cube
	cube, _, _ := utils.ReadLine(cubesFile, i)

	// * 4. Generate unit clauses from the literals in the cube
	cubeClauses_ := strings.Split(strings.TrimPrefix(cube, "a "), " ")
	cubeClauses := lo.Map(cubeClauses_[:len(cubeClauses_)-1], func(s string, _ int) string {
		return s + " 0"
	})

	// * 5. Get the num. of variables and clauses along with the body
	var numVars, numClauses int
	fmt.Sscanf(instance, "p cnf %d %d\n", &numVars, &numClauses)
	n := len(fmt.Sprintf("p cnf %d %d\n", numVars, numClauses))
	body := instance[n:]

	// * 6. Generate a new header with an increased number of clauses
	newHeader := fmt.Sprintf("p cnf %d %d", numVars, numClauses+len(cubeClauses))

	// * 7. Assemble the new encoding
	newEncoding := fmt.Sprintf("%s\n%s%s\n",
		newHeader,
		body,
		strings.Join(cubeClauses, "\n"))

	return newEncoding, nil
}

func GenerateSubProblemWithThreshold(instanceName string, i int, threshold *uint) error {
	subProblemFilePath := fmt.Sprintf("%scube%d_%s.cnf", constants.EncodingsDirPath, i, instanceName)
	if utils.FileExists(subProblemFilePath) {
		return nil
	}

	// * 1. Generate the encoding as a string
	encoding, err := GenerateSubProblemAsStringWithThreshold(instanceName, i, threshold)
	if err != nil {
		return err
	}

	// * 2. Write the subproblem file
	subProblemFile, err := os.OpenFile(subProblemFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New("failed to create the subproblem file")
	}
	defer subProblemFile.Close()

	if _, err := subProblemFile.WriteString(encoding); err != nil {
		return errors.New("failed to add the clauses to the subproblem file")
	}

	return nil
}

func GenerateSubProblem(instanceName string, i int) error {
	if err := GenerateSubProblemWithThreshold(instanceName, i, nil); err != nil {
		return err
	}

	return nil
}

func generateCubes(instanceName string, cutoffVars uint) {
	// Skip if the cubes already exist
	if utils.FileExists(fmt.Sprintf("%s%s.icnf", constants.EncodingsDirPath, instanceName)) {
		return
	}

	// Invoke March for generating the cubes that is held in an .icnf file
	core.March(fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, instanceName), fmt.Sprintf("%s%s.icnf", constants.EncodingsDirPath, instanceName), cutoffVars, time.Second*5000)
}
