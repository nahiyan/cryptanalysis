package services

import (
	"bufio"
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/pipeline"
	"cryptanalysis/internal/simplifier"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/alitto/pond"
)

func (simplifierSvc *SimplifierService) getLogPath(instancePath string) string {
	logFileName := path.Base(instancePath[:len(instancePath)-3] + "log")
	logFilePath := path.Join(simplifierSvc.configSvc.Config.Paths.Logs, logFileName)

	return logFilePath
}

func (simplifierSvc *SimplifierService) TrackedInvoke(simplifier_ simplifier.Simplifier, encoding encoder.Encoding, outputFilePath string, conflicts int, parameters pipeline.SimplifyParams) error {
	if simplifier_ != simplifier.Cadical {
		log.Fatal("Simplifier not supported")
	}
	config := simplifierSvc.configSvc.Config

	simplifierBinPath := config.Paths.Bin.Cadical
	rsFilePath := outputFilePath + ".rs.txt"
	args := []string{"-o", outputFilePath, "-e", rsFilePath, "-c", strconv.Itoa(conflicts)}

	cmd := exec.Command(simplifierBinPath, args...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	// Provide the encoding to the solver (through stdin)
	var encodingString string
	{
		var encodingReader io.Reader
		if cube, exists := encoding.Cube.Get(); exists {
			// Handle cubes
			cubesetPath, err := encoding.GetCubesetPath(config.Paths.Cubesets)
			if err != nil {
				return err
			}

			encodingReader, err = simplifierSvc.cubeSelectorSvc.EncodingFromCube(encoding.BasePath, cubesetPath, cube.Index)
			if err != nil {
				return err
			}
		} else {
			// Handle regular files
			encodingReader, err = os.OpenFile(encoding.BasePath, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}
		}
		builder := strings.Builder{}
		io.Copy(&builder, encodingReader)
		encodingString = builder.String()
	}

	_, err = io.Copy(stdinPipe, strings.NewReader(encodingString))
	if err != nil {
		return err
	}
	stdinPipe.Close()

	logFilePath := simplifierSvc.getLogPath(outputFilePath)
	simplifierSvc.filesystemSvc.WriteFromPipe(stdoutPipe, logFilePath)
	cmd.Wait()

	// Reconstruct the instance if required
	if parameters.Reconstruct {
		err := simplifierSvc.simplificationSvc.Reconstruct(outputFilePath, rsFilePath)
		simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to reconstruct instance")
	}

	// Add back the removed comments
	if parameters.PreserveComments {
		simplifiedEncoding, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_APPEND, 0644)
		simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to open encoding for appending")
		defer simplifiedEncoding.Close()

		scanner := bufio.NewScanner(strings.NewReader(encodingString))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "c ") {
				simplifiedEncoding.WriteString(line + "\n")
			}
		}
	}

	// TODO: Add more details
	simplifierSvc.logSvc.Info("Simplifier: " + encoding.GetName())

	return nil
}

func (simplifierSvc *SimplifierService) ShouldSkip(instancePath string, simplifier_ simplifier.Simplifier) bool {
	if !simplifierSvc.filesystemSvc.FileExists(instancePath) {
		return false
	}

	logFilePath := simplifierSvc.getLogPath(instancePath)

	if simplifierSvc.combinedLogsSvc.IsLoaded() {
		_, err := simplifierSvc.ParseOutputFromCombinedLog(path.Base(logFilePath), simplifier_)
		return err == nil
	}

	_, err := simplifierSvc.ParseOutputFromFile(logFilePath, simplifier_)
	return err == nil
}

func (simplifierSvc *SimplifierService) RunWith(simplifier_ simplifier.Simplifier, encodings []encoder.Encoding, parameters pipeline.SimplifyParams) []encoder.Encoding {
	simplifierSvc.logSvc.Info("Simplifier: started with " + string(simplifier_))
	simplifiedEncodings := []encoder.Encoding{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	for _, encoding := range encodings {
		// Get the conflicts list
		info, err := simplifierSvc.encoderSvc.ProcessInstanceName(encoding.GetName())
		simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to process instance name")
		var conflictsList []int
		conflictsList, exists := parameters.ConflictsMap[info.Steps]
		if !exists {
			conflictsList = parameters.Conflicts
		}

		for _, conflicts := range conflictsList {
			var outputFilePath string
			switch simplifier_ {
			case simplifier.Cadical:
				outputFilePath = fmt.Sprintf("%s.cadical_c%d.cnf", encoding.GetEncodingPath(), conflicts)
			}
			simplifiedEncodings = append(simplifiedEncodings, encoder.Encoding{BasePath: outputFilePath})

			// Check if it should be skipped
			if simplifierSvc.ShouldSkip(outputFilePath, simplifier_) {
				// TODO: Add more details
				simplifierSvc.logSvc.Info("Simplifier: skipped " + encoding.BasePath)
				continue
			}
			os.Remove(outputFilePath)

			// encodingPath := encoding.GetEncodingPath()
			pool.Submit(func(simplifier_ simplifier.Simplifier, encoding encoder.Encoding, outputFilePath string, conflicts int, parameters pipeline.SimplifyParams) func() {
				return func() {
					err := simplifierSvc.TrackedInvoke(simplifier_, encoding, outputFilePath, conflicts, parameters)
					simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to simplify "+encoding.GetName())
				}
			}(simplifier_, encoding, outputFilePath, conflicts, parameters))
		}
	}

	pool.StopAndWait()
	simplifierSvc.logSvc.Info("Simplifier: stopped with " + string(simplifier_))

	return simplifiedEncodings
}

func (simplifierSvc *SimplifierService) Run(encodings []encoder.Encoding, parameters pipeline.SimplifyParams) []encoder.Encoding {
	err := simplifierSvc.filesystemSvc.PrepareDirs([]string{"encodings", "logs"})
	simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to prepare directories")

	switch parameters.Name {
	case simplifier.Cadical:
		return simplifierSvc.RunWith(simplifier.Cadical, encodings, parameters)
	}

	return []encoder.Encoding{}
}
