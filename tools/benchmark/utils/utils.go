package utils

import (
	"benchmark/constants"
	"benchmark/types"
	"fmt"
	"io/ioutil"
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
	AppendLog(constants.BENCHMARK_LOG_FILE_NAME, message)
}

func AppendVerificationLog(message string) {
	AppendLog(constants.VERIFICATION_LOG_FILE_NAME, message)
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

							cb(i, satSolver, steps, hash, xorOption, adderType, dobbertin)
							i++
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
		return constants.CRYPTOMINISAT
	case constants.ArgKissat:
		return constants.KISSAT
	case constants.ArgCadical:
		return constants.CADICAL
	case constants.ArgGlucoseSyrup:
		return constants.GLUCOSE
	case constants.ArgMapleSat:
		return constants.MAPLESAT
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
	items, _ := ioutil.ReadDir(constants.RESULTS_DIR_PATH)
	for _, item := range items {
		if item.IsDir() {
			continue
		}

		if path.Ext(item.Name()) != ".log" {
			continue
		}

		data, err := os.ReadFile(path.Join(constants.RESULTS_DIR_PATH, item.Name()))
		if err != nil {
			fmt.Println("Failed to aggregate logs", err.Error())
		}

		data_ := strings.TrimSpace(string(data))

		fmt.Println(item.Name())
		if strings.HasPrefix(item.Name(), "verification") {
			AppendVerificationLog(data_)
		} else if strings.HasPrefix(item.Name(), "benchmark") {
			AppendBenchmarkLog(data_)
		}
	}
}
