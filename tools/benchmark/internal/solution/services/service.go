package services

import (
	"benchmark/internal/solution"
	"benchmark/internal/solver"
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"gonum.org/v1/gonum/stat/combin"
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

	scanner := bufio.NewScanner(instanceFile)
	newBody := ""
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimPrefix(line, "v ")

		if strings.HasPrefix(line, "s SATISFIABLE") {
			continue
		}

		newBody += line + " "
	}
	instanceFile.Close()

	outputPath, err := os.OpenFile(encodingPath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = outputPath.WriteString("SAT" + "\n" + newBody + "\n")
	if err != nil {
		return err
	}

	return nil
}

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

func (solutionSvc *SolutionService) ReconstructCadical(instancePath string, ranges []solution.Range) error {
	variableVariations := make(map[int][]int)

	// * 1. Read the reconstruction stack file and determine the literals that need to be preserved
	reconstStackFilePath := instancePath + ".rs.txt"
	reconstStackFile, err := os.OpenFile(reconstStackFilePath, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer reconstStackFile.Close()

	scanner := bufio.NewScanner(reconstStackFile)
	for scanner.Scan() {
		line := scanner.Text()

		segments := strings.Split(line, " ")

		// * 2. Skip line if the number of segments is not at least 4
		segmentsCount := len(segments)
		if segmentsCount < 4 {
			continue
		}

		// * 3. See if the witness literal should be preserved
		var literal, variable int
		{
			literal, _ = strconv.Atoi(segments[segmentsCount-2])
			variable_ := math.Abs(float64(literal))
			variable = int(variable_)

			inRange := false
			for _, range_ := range ranges {
				if variable >= range_.Start && variable <= range_.End {
					inRange = true
				}
			}

			if !inRange {
				continue
			}
		}

		// * 4. Mark the literal
		if _, exists := variableVariations[variable]; !exists {
			variableVariations[variable] = make([]int, 0)
		}

		if !lo.Contains(variableVariations[variable], literal) {
			variableVariations[variable] = append(variableVariations[variable], literal)
		}
	}

	// End if we have no variable to correct
	if len(variableVariations) == 0 {
		return nil
	}

	// * 5. Parse the solution of the given instance
	solutionFilePath := instancePath + ".sol"
	solutionContent_, err := os.ReadFile(solutionFilePath)
	if err != nil {
		return err
	}

	solutionContent := string(solutionContent_)

	solutionLiterals := lo.Filter(lo.Map(strings.Fields(string(solutionContent)), func(literal_ string, _ int) int {
		literal, err := strconv.Atoi(literal_)
		if err != nil {
			return 0
		}

		return literal
	}), func(value, i2 int) bool {
		return value != 0
	})

	// * 6. Generate a combination of the variables with varying values
	varyingVariables := []uint{}
	for variable, variation := range variableVariations {
		if len(variation) == 2 {
			varyingVariables = append(varyingVariables, uint(variable))
		}
	}

	factors := []int{}
	for i := 0; i < len(varyingVariables); i++ {
		factors = append(factors, 2)
	}

	list := combin.Cartesian(factors)
	for _, v := range list {
		literalCombination := []int{}
		for i, value := range v {
			var literal int
			if value == 0 {
				literal = int(varyingVariables[i])
			} else {
				literal = -int(varyingVariables[i])
			}
			literalCombination = append(literalCombination, literal)
		}

		// * 7. Generate a new solution with overriden literals
		newSolutionLiterals := lo.Map(solutionLiterals, func(value, _ int) string {
			for _, literal := range literalCombination {
				if literal == value || literal == -value {
					return fmt.Sprintf("%d", literal)
				}
			}

			absValue := value
			if absValue < 0 {
				absValue = -absValue
			}

			if literal, exists := variableVariations[absValue]; exists {
				return fmt.Sprintf("%d", literal[0])
			}

			return fmt.Sprintf("%d", value)
		})

		newSolution := "SAT\n" + strings.Join(newSolutionLiterals, " ") + " 0"
		err := os.WriteFile(solutionFilePath, []byte(newSolution), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
