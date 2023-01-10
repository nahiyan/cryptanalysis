package services

import (
	"benchmark/internal/pipeline"
	"benchmark/internal/simplification"
	"benchmark/internal/simplifier"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/alitto/pond"
)

type CadicalOutput struct {
	FreeVariables int
	Clauses       int
	Eliminations  int
	ProcessTime   time.Duration
}

func (simplifierSvc *SimplifierService) ReadCadicalOutput(output string) CadicalOutput {
	summaryIndex := strings.Index(output, "c ?")
	if summaryIndex == -1 {
		return CadicalOutput{}
	}
	summary := output[summaryIndex:]

	eliminations := 0
	eliminatedIndex := strings.Index(summary, "c eliminated")
	if eliminatedIndex == -1 {
		return CadicalOutput{}
	}
	fmt.Sscanf(summary[eliminatedIndex:], "c eliminated: %d", &eliminations)

	freeVariables := 0
	if index := strings.Index(summary, "c ?"); index != -1 {
		segments := strings.Fields(summary[index:])
		freeVariables, _ = strconv.Atoi(segments[12])
	}

	var processTime time.Duration
	if index := strings.Index(summary, "c total process time since initialization:"); index != -1 {
		seconds := 0.0
		fmt.Sscanf(summary[index:], "c total process time since initialization: %f", &seconds)
		processTime = time.Duration(seconds) * time.Second
	}

	clauses := 0
	variables := 0
	if index := strings.Index(summary, "c writing 'p cnf"); index != -1 {
		fmt.Sscanf(summary[index:], "c writing 'p cnf %d %d' header", &variables, &clauses)
	}

	return CadicalOutput{
		FreeVariables: freeVariables,
		Eliminations:  eliminations,
		ProcessTime:   processTime,
		Clauses:       clauses,
	}
}

func (simplifierSvc *SimplifierService) TrackedInvoke(encoding, outputFilePath string, conflicts int, parameters pipeline.Simplifying) error {
	config := simplifierSvc.configSvc.Config
	args := []string{encoding, "-o", outputFilePath, "-e", outputFilePath + ".rs.txt"}
	args = append(args, "-c", fmt.Sprintf("%d", conflicts))
	if parameters.Timeout > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", parameters.Timeout))
	}

	cmd := exec.Command(config.Paths.Bin.Cadical, args...)
	output_, err := cmd.Output()
	if err != nil {
		return err
	}

	cadicalOutput := simplifierSvc.ReadCadicalOutput(string(output_))
	eliminations := cadicalOutput.Eliminations
	freeVariables := cadicalOutput.FreeVariables
	processTime := cadicalOutput.ProcessTime
	clauses := cadicalOutput.Clauses
	time := fmt.Sprintf("%.3f", processTime.Seconds())

	fmt.Println("Simplifier:", conflicts, "conflicts", eliminations, "eliminated", freeVariables, "remaining", clauses, "clauses", time, encoding)

	err = simplifierSvc.simplificationSvc.Register(outputFilePath, simplification.Simplification{
		FreeVariables: freeVariables,
		Simplifier:    "cadical",
		ProcessTime:   processTime,
		Eliminaton:    eliminations,
		Conflicts:     conflicts,
		Clauses:       clauses,
	})
	if err != nil {
		return err
	}

	return nil
}

func (simplifierSvc *SimplifierService) RunCadical(encodings []string, parameters pipeline.Simplifying) []string {
	fmt.Println("Simplifier: started")
	simplifiedEncodings := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	for _, encoding := range encodings {
		for _, conflicts := range parameters.Conflicts {
			outputFilePath := fmt.Sprintf("%s.cadical_c%d.cnf", encoding, conflicts)

			if simplifierSvc.filesystemSvc.FileExists(outputFilePath) {
				fmt.Println("Simplifier: skipped", encoding)
				simplifiedEncodings = append(simplifiedEncodings, outputFilePath)

				continue
			}

			pool.Submit(func(encoding, outputFilePath string, conflicts int, parameters pipeline.Simplifying) func() {
				return func() {
					err := simplifierSvc.TrackedInvoke(encoding, outputFilePath, conflicts, parameters)
					simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to simplify "+encoding)
				}
			}(encoding, outputFilePath, conflicts, parameters))
			simplifiedEncodings = append(simplifiedEncodings, outputFilePath)
		}
	}

	pool.StopAndWait()
	fmt.Println("Simplifier: stopped")
	return simplifiedEncodings
}

func (simplifierSvc *SimplifierService) Run(encodings []string, parameters pipeline.Simplifying) []string {
	switch parameters.Name {
	case simplifier.Cadical:
		return simplifierSvc.RunCadical(encodings, parameters)
	}

	return []string{}
}
