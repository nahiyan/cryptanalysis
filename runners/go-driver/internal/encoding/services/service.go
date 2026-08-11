package services

import (
	"bufio"
	"cryptanalysis/internal/encoding"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (encodingSvc *EncodingService) GetInfo(encodingPath string) (encoding.EncodingInfo, error) {
	info := encoding.EncodingInfo{}

	instanceFile, err := os.Open(encodingPath)
	if err != nil {
		return info, err
	}
	defer instanceFile.Close()

	variables := map[int]bool{}
	scanner := bufio.NewScanner(instanceFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip comments and problem definition lines
		if strings.HasPrefix(line, "c") || strings.HasPrefix(line, "p") {
			continue
		}

		info.Clauses += 1

		matches := regexp.MustCompile(`[+-]?\d+`).FindAllString(line, -1)
		for _, match := range matches {
			variable, err := strconv.Atoi(match)
			if err != nil {
				return info, err
			}
			if variable != 0 {
				variables[abs(variable)] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return info, err
	}

	info.FreeVariables = len(variables)

	return info, nil
}
