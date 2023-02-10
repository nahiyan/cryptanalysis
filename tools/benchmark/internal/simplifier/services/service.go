package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"benchmark/internal/simplifier"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/alitto/pond"
)

type EncodingPromise struct {
	Encoding string
}

func (promise EncodingPromise) GetPath() string {
	return promise.Encoding
}

func (promise EncodingPromise) Get(dependencies map[string]interface{}) string {
	return promise.Encoding
}

func (simplifierSvc *SimplifierService) getLogPath(instancePath string) string {
	logFileName := path.Base(instancePath[:len(instancePath)-3] + "log")
	logFilePath := path.Join(simplifierSvc.configSvc.Config.Paths.Logs, logFileName)

	return logFilePath
}

// TODO: Finalize the code
func (simplifierSvc *SimplifierService) TrackedInvoke(simplifier_ simplifier.Simplifier, encoding, outputFilePath string, conflicts int, parameters pipeline.SimplifyParams) error {
	config := simplifierSvc.configSvc.Config

	var simplifierBinPath string
	args := []string{encoding}
	if simplifier_ == simplifier.Cadical {
		simplifierBinPath = config.Paths.Bin.Cadical
		args = append(args, "-o", outputFilePath, "-e", outputFilePath+".rs.txt", "-c", strconv.Itoa(conflicts))
	} else {
		simplifierBinPath = config.Paths.Bin.Satelite
		args = append(args, outputFilePath, outputFilePath+".var_map.txt", "+pre")
	}

	cmd := exec.Command(simplifierBinPath, args...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	logFilePath := simplifierSvc.getLogPath(outputFilePath)
	simplifierSvc.filesystemSvc.WriteFromPipe(stdoutPipe, logFilePath)
	cmd.Wait()

	// TODO: Add more details
	simplifierSvc.logSvc.Info("Simplifier: " + encoding)

	return nil
}

func (simplifierSvc *SimplifierService) ShouldSkip(instancePath string, simplifier_ simplifier.Simplifier) bool {
	logFilePath := simplifierSvc.getLogPath(instancePath)
	_, err := simplifierSvc.ParseOutput(logFilePath, simplifier_)
	return err == nil
}

func (simplifierSvc *SimplifierService) RunWith(simplifier_ simplifier.Simplifier, encodings []encoder.Encoding, parameters pipeline.SimplifyParams) []encoder.Encoding {
	simplifierSvc.logSvc.Info("Simplifier: started with " + string(simplifier_))
	simplifiedEncodings := []encoder.Encoding{}
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	for _, encoding := range encodings {
		conflictsList := parameters.Conflicts
		if simplifier_ == simplifier.Satelite {
			conflictsList = []int{0}
		}

		for _, conflicts := range conflictsList {
			var outputFilePath string
			switch simplifier_ {
			case simplifier.Cadical:
				outputFilePath = fmt.Sprintf("%s.cadical_c%d.cnf", encoding.BasePath, conflicts)
			case simplifier.Satelite:
				outputFilePath = fmt.Sprintf("%s.satelite.cnf", encoding.BasePath)
			}
			simplifiedEncodings = append(simplifiedEncodings, encoder.Encoding{BasePath: outputFilePath})

			// Check if it should be skipped
			if simplifierSvc.ShouldSkip(outputFilePath, simplifier_) {
				// TODO: Add more details
				simplifierSvc.logSvc.Info("Simplifier: skipped " + encoding.BasePath)
				continue
			}
			os.Remove(outputFilePath)

			pool.Submit(func(simplifier_ simplifier.Simplifier, encoding, outputFilePath string, conflicts int, parameters pipeline.SimplifyParams) func() {
				return func() {
					err := simplifierSvc.TrackedInvoke(simplifier_, encoding, outputFilePath, conflicts, parameters)
					simplifierSvc.errorSvc.Fatal(err, "Simplifier: failed to simplify "+encoding)
				}
			}(simplifier_, encoding.BasePath, outputFilePath, conflicts, parameters))
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
	case simplifier.Satelite:
		return simplifierSvc.RunWith(simplifier.Satelite, encodings, parameters)
	}

	return []encoder.Encoding{}
}
