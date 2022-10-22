package utils

import (
	"benchmark/constants"
	"benchmark/types"
	"fmt"
	"os"
	"path"
	"strings"
)

func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func AppendLog(filename, message string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to write logs")
	}
	_, err = f.WriteString(message + "\n")
	if err != nil {
		panic("Failed to write logs")
	}
	f.Close()
}

func AppendBenchmarkLog(message string) {
	AppendLog(constants.BenchmarkLogFileName, message)
}
func AppendValidResultsLog(message string) {
	AppendLog(constants.BenchmarkLogFileName, message)
}

func AppendVerificationLog(message string) {
	AppendLog(constants.VerificationLogFileName, message)
}

func LoopThroughVariations(context *types.CommandContext, cb func(uint, string, uint, string, uint, string, uint)) {
	for _, satSolver := range context.VariationsSatSolvers {
		var i uint = 0
		for _, steps := range context.VariationsSteps {
			for _, hash := range context.VariationsHashes {
				for _, xorOption := range context.VariationsXor {
					for _, adderType := range context.VariationsAdders {
						for _, dobbertin := range context.VariationsDobbertin {
							// Skip dobbertin's attacks when steps count < 27
							if steps < 27 && dobbertin == 1 {
								continue
							}

							// No XOR for SAT Solvers other than CryptoMiniSAT and XNFSAT
							if xorOption == 1 && satSolver != constants.ArgCryptoMiniSat && satSolver != constants.ArgXnfSat {
								xorOption = 0
							}

							cb(i, satSolver, steps, hash, xorOption, adderType, dobbertin)
							i += 1
						}
					}
				}
			}
		}
	}
}

func ResolveSatSolverName(shortcut string) string {
	switch shortcut {
	case constants.ArgCryptoMiniSat:
		return constants.CryptoMiniSat
	case constants.ArgKissat:
		return constants.Kissat
	case constants.ArgCadical:
		return constants.Cadical
	case constants.ArgGlucoseSyrup:
		return constants.Glucose
	case constants.ArgMapleSat:
		return constants.MapleSat
	case constants.ArgXnfSat:
		return constants.XnfSat
	}

	return ""
}

func ResolveAdderType(shortcut string) string {
	switch shortcut {
	case "cc":
		return "counter_chain"
	case "dm":
		return "dot_matrix"
	}

	return ""
}

func InstancesCount(commandContext *types.CommandContext) uint {
	var count uint = 0
	for _, steps := range commandContext.VariationsSteps {
		for range commandContext.VariationsHashes {
			for range commandContext.VariationsXor {
				for range commandContext.VariationsAdders {
					for _, dobbertin := range commandContext.VariationsDobbertin {
						// Skip dobbertin's attacks when steps count < 27
						if steps < 27 && dobbertin == 1 {
							continue
						}

						count++
					}
				}
			}
		}
	}

	return count
}

func AggregateLogs() {
	items, _ := os.ReadDir(constants.ResultsDirPat)
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if path.Ext(item.Name()) != ".log" {
			continue
		}

		data, err := os.ReadFile(path.Join(constants.ResultsDirPat, item.Name()))
		if err != nil {
			fmt.Println("Failed to aggregate logs", err.Error())
		}

		data_ := strings.TrimSpace(string(data))

		fmt.Println(item.Name())
		if strings.HasPrefix(item.Name(), "verification") {
			AppendVerificationLog(data_)
		} else if strings.HasPrefix(item.Name(), "benchmark") {
			AppendBenchmarkLog(data_)
		} else if strings.HasPrefix(item.Name(), "valid_results") {
			AppendValidResultsLog(data_)
		}
	}
}
