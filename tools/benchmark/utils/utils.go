package utils

import (
	"benchmark/constants"
	"benchmark/types"
	"os"
)

func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func appendLog(filename, message string) {
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
	appendLog(constants.BENCHMARK_LOG_FILE_NAME, message)
}

func AppendVerificationLog(message string) {
	appendLog(constants.VERIFICATION_LOG_FILE_NAME, message)
}

func LoopThroughVariations(context *types.CommandContext, cb func(uint, string, uint, string, uint, string, uint)) {
	var i uint = 0
	for _, satSolver := range context.VariationsSatSolvers {
		for _, steps := range context.VariationsSteps {
			for _, hash := range context.VariationsHashes {
				for _, xorOption := range context.VariationsXor {
					for _, adderType := range context.VariationsAdders {
						for _, dobbertin := range context.VariationsDobbertin {
							// Skip dobbertin's attacks when steps count < 28
							if steps < 28 && dobbertin == 1 {
								dobbertin = 0
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
