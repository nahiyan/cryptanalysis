package encodings

import (
	"benchmark/constants"
	"benchmark/types"
	"bufio"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type Cnf struct {
	Clauses       []types.Clause
	FreeVariables uint
	ClausesCount  uint
}

func Process(instanceName string) (Cnf, error) {
	c := Cnf{}
	instanceFile, err := os.Open(path.Join(constants.EncodingsDirPath, instanceName+".cnf"))
	if err != nil {
		return c, err
	}
	defer instanceFile.Close()

	variables := []uint{}

	scanner := bufio.NewScanner(instanceFile)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if it's the header
		if strings.HasPrefix(line, "p cnf ") {
			continue
		}

		// Check if it's a comment
		if strings.HasPrefix(line, "c ") {
			continue
		}

		// Divide the line into segments
		segments := strings.Fields(line)

		// Ignore an empty clause
		if len(segments) < 2 {
			continue
		}

		// Check if the clause is zero-terminated
		if lastSegment, err := lo.Last(segments); err != nil || strings.TrimSpace(lastSegment) != "0" {
			continue
		}

		// Parse the literals
		literals_ := lo.Map(segments, func(s string, i int) int {
			v, err := strconv.Atoi(s)
			if err != nil {
				return 0
			}

			return v
		})
		literals := lo.Filter(literals_, func(i1, i2 int) bool {
			return i1 != 0
		})

		for _, literal := range literals {
			variable := uint(int(math.Abs(float64(literal))))
			// New variable
			if !lo.Contains(variables, variable) {
				variables = append(variables, uint(variable))
				c.FreeVariables += 1
			}
		}

		c.Clauses = append(c.Clauses, literals)
		c.ClausesCount += 1
	}

	if err := scanner.Err(); err != nil {
		return c, err
	}

	return c, nil
}
