package encodings

import (
	"benchmark/config"
	"benchmark/constants"
	"benchmark/types"
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"gonum.org/v1/gonum/stat/combin"
)

func ReconstructEncoding(instanceFilePath, reconstStackFilePath string, ranges []types.Range) error {
	// Clauses to preserve
	clauses := []types.Clause{}

	// * 1. Read the reconstruction stack file
	reconstStackFile, err := os.Open(reconstStackFilePath)
	if err != nil {
		return err
	}
	defer reconstStackFile.Close()

	// * 2. Go through the reconstruction stack file line by line
	scanner := bufio.NewScanner(reconstStackFile)
	for scanner.Scan() {
		line := scanner.Text()
		segments := strings.Split(line, " ")

		// * 3. Skip line if the number of segments is not at least 4
		segmentsCount := len(segments)
		if segmentsCount < 4 {
			continue
		}

		// * 4. See if the witness literal should be preserved
		{
			literal, _ := strconv.Atoi(segments[segmentsCount-2])
			variable_ := math.Abs(float64(literal))
			variable := int(variable_)

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

		// * 5. Preserve the clause for the literal
		clause := types.Clause{}
		for _, segment := range segments {
			segment = strings.TrimSpace(segment)

			// Determine if the clause ends here
			if segment == "0" {
				break
			}

			literal, _ := strconv.Atoi(segment)
			clause = append(clause, literal)
		}

		clauses = append(clauses, clause)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// * 6. Reconstruct the instance file

	// Skip if there is no clause to preserve
	if len(clauses) == 0 {
		return nil
	}

	{
		instanceFile, err := os.OpenFile(instanceFilePath, os.O_RDONLY, 0600)
		if err != nil {
			return err
		}
		defer instanceFile.Close()

		scanner := bufio.NewScanner(instanceFile)
		variablesCount, clausesCount := 0, 0
		lines := []string{}
		headerLineIndex := 0
		i := 0
		for scanner.Scan() {
			line := scanner.Text()

			// Process the header
			if strings.HasPrefix(line, "p cnf") {
				headerLineIndex = i
				fmt.Sscanf(line, "p cnf %d %d", &variablesCount, &clausesCount)
			}

			lines = append(lines, line)
			i++
		}

		// * 7. Add the removed clauses back
		for _, clause := range clauses {
			line := strings.Join(lo.Map(clause, func(i1, i2 int) string {
				return strconv.Itoa(i1)
			}), " ") + " 0"
			lines = append(lines, line)
		}

		// * 8. Write the reconstructed instance
		lines[headerLineIndex] = fmt.Sprintf("p cnf %d %d", variablesCount, clausesCount+len(clauses))

		// Overwrite the instance file
		instanceFile, err = os.Create(instanceFilePath)
		if err != nil {
			return err
		}

		newInstance := strings.Join(lines, "\n") + "\n"
		if _, err := instanceFile.WriteString(newInstance); err != nil {
			return err
		}
	}

	return nil
}

func ReconstructSolution(instanceName, reconstStackFilePath string, ranges []types.Range) error {
	variableVariations := make(map[int][]int)

	// * 1. Read the reconstruction stack file and determine the literals that need to be preserved
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
	solutionFilePath := path.Join(constants.SolutionsDirPath, instanceName+".sol")
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
		output_, err := exec.Command("bash", "-c", fmt.Sprintf("printf \"%s\" | %s %d", newSolution, config.Get().Paths.Bin.Verifier, 43)).Output()
		if err != nil {
			return err
		}
		output := string(output_)

		// * 8. Check if it's the valid solution
		if strings.Contains(output, "Solution's hash matches the target!") {
			if err := os.WriteFile(solutionFilePath, []byte(newSolution), 0644); err != nil {
				return err
			}
			fmt.Println("Successfully reconstruted the solution")
			return nil
		}
	}

	return nil
}
