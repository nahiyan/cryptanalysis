package encodings

import (
	"benchmark/types"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func Reconstruct(instanceFilePath, reconstStackFilePath string, ranges []types.Range) error {
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

		// * 3. Skip line if the number of segments is at least 4
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

		newInstance := strings.Join(lines, "\n")
		if _, err := instanceFile.WriteString(newInstance); err != nil {
			return err
		}
	}

	return nil
}
