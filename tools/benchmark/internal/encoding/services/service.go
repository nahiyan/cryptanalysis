package services

import (
	"benchmark/internal/encoding"
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

func (encodingSvc *EncodingService) GetInfo(encodingPath string) (encoding.EncodingInfo, error) {
	info := encoding.EncodingInfo{}

	instanceFile, err := os.Open(encodingPath)
	if err != nil {
		return info, err
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
				info.FreeVariables += 1
			}
		}

		info.Clauses += 1
	}

	if err := scanner.Err(); err != nil {
		return info, err
	}

	return info, nil
}
