package services

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bitfield/script"
)

func (cuberSvc *CuberService) ParseOutput(logFilePath string) (time.Duration, int, int, error) {
	cubes := 0
	refutedLeaves := 0
	processTime := time.Duration(0)

	matches, err := script.File(logFilePath).MatchRegexp(regexp.MustCompile("(c time)|(c number of cubes)")).Slice()
	if err != nil {
		return processTime, cubes, refutedLeaves, err
	}

	if len(matches) != 2 {
		return processTime, cubes, refutedLeaves, errors.New("invalid Match output")
	}

	// Parse the process time
	fields := strings.Fields(matches[0])
	if len(fields) != 5 {
		return processTime, cubes, refutedLeaves, errors.New("invalid Match output")
	}

	seconds, err := strconv.ParseFloat(fields[len(fields)-2], 64)
	if err != nil {
		return processTime, cubes, refutedLeaves, err
	}
	processTime = time.Duration(seconds*1000) * time.Millisecond

	// Parse the number of cubes and refuted leaves
	_, err = fmt.Sscanf(matches[1], "c number of cubes %d, including %d refuted leaves", &cubes, &refutedLeaves)
	if err != nil {
		return processTime, cubes, refutedLeaves, err
	}

	return processTime, cubes, refutedLeaves, nil
}
