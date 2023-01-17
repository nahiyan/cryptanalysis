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
	"github.com/samber/lo"
)

type CadicalOutput struct {
	FreeVariables int
	Clauses       int
	Eliminations  int
	ProcessTime   time.Duration
}

type SateliteOutput struct {
	FreeVariables int
	Clauses       int
	ProcessTime   time.Duration
}

type Simplifier string

type EncodingPromise struct {
	Encoding string
}

func (promise EncodingPromise) GetPath() string {
	return promise.Encoding
}

func (promise EncodingPromise) Get() string {
	return promise.Encoding
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
		processTime = time.Duration(seconds*1000) * time.Millisecond
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

func (simplifierSvc *SimplifierService) ReadSateliteOutput(output string) SateliteOutput {
	result := SateliteOutput{}

	cpuTimeIndex := strings.Index(output, "CPU time:")
	if cpuTimeIndex != -1 {
		processTime := 0.0
		x := ""
		fmt.Sscanf(output[cpuTimeIndex:], "CPU time%s%f s", &x, &processTime)
		processTime *= 1000
		result.ProcessTime = time.Duration(processTime) * time.Millisecond
	}

	resultIndex := strings.Index(output, "Result")
	if resultIndex != -1 {
		x := ""
		fmt.Sscanf(output[resultIndex:], "Result%s#vars: %d #clauses: %d", &x, &result.FreeVariables, &result.Clauses)
	}

	return result
}

func (simplifierSvc *SimplifierService) TrackedCadicalInvoke(encoding, outputFilePath string, conflicts int, parameters pipeline.Simplifying) error {
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
		Simplifier:    simplifier.Cadical,
		ProcessTime:   processTime,
		Eliminaton:    eliminations,
		Conflicts:     conflicts,
		Clauses:       clauses,
		InstanceName:  encoding,
	})
	if err != nil {
		return err
	}

	return nil
}

func (simplifierSvc *SimplifierService) TrackedSateliteInvoke(encoding, outputFilePath string, parameters pipeline.Simplifying) error {
	config := simplifierSvc.configSvc.Config
	args := []string{encoding, outputFilePath, outputFilePath + ".var_map.txt", "+pre"}

	cmd := exec.Command(config.Paths.Bin.Satelite, args...)
	output_, err := cmd.Output()
	if err != nil {
		return err
	}

	output := simplifierSvc.ReadSateliteOutput(string(output_))

	freeVariables := output.FreeVariables
	processTime := output.ProcessTime
	clauses := output.Clauses
	time := fmt.Sprintf("%.3f", processTime.Seconds())

	fmt.Println("Simplifier:", freeVariables, "remaining", clauses, "clauses", time, encoding)

	err = simplifierSvc.simplificationSvc.Register(outputFilePath, simplification.Simplification{
		FreeVariables: freeVariables,
		Simplifier:    simplifier.Satelite,
		ProcessTime:   processTime,
		Clauses:       clauses,
		InstanceName:  encoding,
	})
	if err != nil {
		return err
	}

	return nil
}

func (simplifierSvc *SimplifierService) RunCadical(encodingPromises []pipeline.EncodingPromise, parameters pipeline.Simplifying) []pipeline.EncodingPromise {
	fmt.Println("Simplifier: started with CaDiCaL")
	simplifiedEncodings := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	for _, encodingPromise := range encodingPromises {
		encoding := encodingPromise.Get()

		for _, conflicts := range parameters.Conflicts {
			outputFilePath := fmt.Sprintf("%s.cadical_c%d.cnf", encoding, conflicts)

			if simplifierSvc.filesystemSvc.FileExists(outputFilePath) {
				fmt.Println("Simplifier: skipped", encoding)
				simplifiedEncodings = append(simplifiedEncodings, outputFilePath)

				continue
			}

			pool.Submit(func(encoding, outputFilePath string, conflicts int, parameters pipeline.Simplifying) func() {
				return func() {
					err := simplifierSvc.TrackedCadicalInvoke(encoding, outputFilePath, conflicts, parameters)
					simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to simplify "+encoding)
				}
			}(encoding, outputFilePath, conflicts, parameters))
			simplifiedEncodings = append(simplifiedEncodings, outputFilePath)
		}
	}

	pool.StopAndWait()
	fmt.Println("Simplifier: stopped with CaDiCaL")

	simplifiedEncodingPromises := lo.Map(simplifiedEncodings, func(simplifiedEncoding string, _ int) pipeline.EncodingPromise {
		return EncodingPromise{Encoding: simplifiedEncoding}
	})

	return simplifiedEncodingPromises
}

func (simplifierSvc *SimplifierService) RunSatelite(encodingPromises []pipeline.EncodingPromise, parameters pipeline.Simplifying) []pipeline.EncodingPromise {
	fmt.Println("Simplifier: started with SatELite")
	simplifiedEncodings := []string{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	for _, encodingPromise := range encodingPromises {
		encoding := encodingPromise.Get()

		outputFilePath := fmt.Sprintf("%s.satelite.cnf", encoding)

		if simplifierSvc.filesystemSvc.FileExists(outputFilePath) {
			fmt.Println("Simplifier: skipped", encoding)
			simplifiedEncodings = append(simplifiedEncodings, outputFilePath)

			continue
		}

		pool.Submit(func(encoding, outputFilePath string, parameters pipeline.Simplifying) func() {
			return func() {
				err := simplifierSvc.TrackedSateliteInvoke(encoding, outputFilePath, parameters)
				simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to simplify "+encoding)
			}
		}(encoding, outputFilePath, parameters))
		simplifiedEncodings = append(simplifiedEncodings, outputFilePath)
	}

	pool.StopAndWait()
	fmt.Println("Simplifier: stopped with SatELite")

	simplifiedEncodingPromises := lo.Map(simplifiedEncodings, func(simplifiedEncoding string, _ int) pipeline.EncodingPromise {
		return EncodingPromise{Encoding: simplifiedEncoding}
	})

	return simplifiedEncodingPromises
}

func (simplifierSvc *SimplifierService) Run(encodingPromises []pipeline.EncodingPromise, parameters pipeline.Simplifying) []pipeline.EncodingPromise {
	switch parameters.Name {
	case simplifier.Cadical:
		return simplifierSvc.RunCadical(encodingPromises, parameters)
	case simplifier.Satelite:
		return simplifierSvc.RunSatelite(encodingPromises, parameters)
	}

	return []pipeline.EncodingPromise{}
}
