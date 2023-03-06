package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (cuberSvc *CuberService) ParseOutputFromFile(logFilePath string) (time.Duration, int, int, error) {
	file, err := os.OpenFile(logFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return time.Duration(0), 0, 0, err
	}
	return cuberSvc.ParseOutput(file)
}

func (cuberSvc *CuberService) ParseOutputFromCombinedLog(logFilePath string) (time.Duration, int, int, error) {
	maybeContent, err := cuberSvc.combinedLogsSvc.Get(logFilePath)
	if err != nil {
		return time.Duration(0), 0, 0, err
	}

	content, exists := maybeContent.Get()
	if !exists {
		return time.Duration(0), 0, 0, os.ErrNotExist
	}

	reader := strings.NewReader(content)
	return cuberSvc.ParseOutput(reader)
}

func (cuberSvc *CuberService) ParseOutput(outputReader io.Reader) (time.Duration, int, int, error) {
	cubes := 0
	refutedLeaves := 0
	processTime := time.Duration(0)

	output := ""
	{
		buf := new(strings.Builder)
		_, err := io.Copy(buf, outputReader)
		if err != nil {
			return time.Duration(0), cubes, refutedLeaves, err
		}
		output = buf.String()
	}

	matches := regexp.MustCompile("(c time.*)|(c number of cubes.*)").FindAllString(output, len(output))
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
