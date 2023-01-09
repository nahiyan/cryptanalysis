package services

import (
	"benchmark/internal/pipeline"
	"benchmark/internal/simplifier"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func (simplifierSvc *SimplifierService) ReadCadicalOutput(output string) (int, int) {
	eliminatedVars := 0
	if index := strings.Index(string(output), "c eliminated:"); index != -1 {
		fmt.Sscanf(string(output)[index:], "c eliminated: %d", &eliminatedVars)
	}

	freeVars := 0
	if index := strings.Index(string(output), "c ?"); index != -1 {
		segments := strings.Split(string(output)[index:], " ")
		if len(segments) >= 14 {
			freeVars, _ = strconv.Atoi(segments[13])
		}
	}

	return freeVars, eliminatedVars
}

func (simplifierSvc *SimplifierService) RunCadical(encodings []string, parameters pipeline.Simplifying) []string {
	config := simplifierSvc.configSvc.Config
	simplifiedEncodings := []string{}
	for _, encoding := range encodings {
		outputFilePath := fmt.Sprintf("%s.cadical_c%d.cnf", encoding, parameters.Conflicts)

		if simplifierSvc.filesystemSvc.FileExists(outputFilePath) {
			fmt.Println("Simplifier: skipped", encoding)
			simplifiedEncodings = append(simplifiedEncodings, outputFilePath)

			continue
		}

		args := []string{encoding, "-o", outputFilePath, "-e", outputFilePath + ".rs.txt"}
		if parameters.Conflicts > 0 {
			args = append(args, "-c", fmt.Sprintf("%d", parameters.Conflicts))
		}
		if parameters.Timeout > 0 {
			args = append(args, "-t", fmt.Sprintf("%d", parameters.Timeout))
		}

		cmd := exec.Command(config.Paths.Bin.Cadical, args...)
		output_, err := cmd.Output()
		simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to simplify "+encoding)

		freeVars, eliminatedVars := simplifierSvc.ReadCadicalOutput(string(output_))
		fmt.Println("Simplifier:", eliminatedVars, "eliminated", freeVars, "remaining", encoding)

		simplifiedEncodings = append(simplifiedEncodings, outputFilePath)
	}

	return simplifiedEncodings
}

func (simplifierSvc *SimplifierService) Run(encodings []string, parameters pipeline.Simplifying) []string {
	switch parameters.Name {
	case simplifier.Cadical:
		return simplifierSvc.RunCadical(encodings, parameters)
	}

	return []string{}
}
