package utils

import (
	"benchmark/constants"
	"benchmark/types"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/samber/lo"
)

func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func AppendLog(filename string, records []string) {
	f, err := os.OpenFile(path.Join(constants.LogsDirPath, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to write logs: " + err.Error())
	}

	csvWriter := csv.NewWriter(f)
	csvWriter.Write(records)
	csvWriter.Flush()

	f.Close()
}

func AppendBenchmarkLog(records []string) {
	AppendLog(constants.BenchmarkLogFileName, records)
}
func AppendValidResultsLog(records []string) {
	AppendLog(constants.ValidResultsLogFileName, records)
}

func AppendVerificationLog(records []string) {
	AppendLog(constants.VerificationLogFileName, records)
}

func LoopThroughVariations(context *types.CommandContext, cb func(uint, string, uint, string, uint, string, uint, uint)) {
	for _, satSolver := range context.VariationsSatSolvers {
		var i uint = 0
		for _, steps := range context.VariationsSteps {
			for _, hash := range context.VariationsHashes {
				for _, xorOption := range context.VariationsXor {
					for _, adderType := range context.VariationsAdders {
						for _, dobbertin := range context.VariationsDobbertin {
							for _, dobbertinBits := range context.VariationsDobbertinBits {
								// Skip any dobbertin bit variation when dobbertin's attack isn't on
								if dobbertin == 0 && dobbertinBits != 32 {
									continue
								}

								// Skip dobbertin's attacks when steps count < 27
								if steps < 27 && dobbertin == 1 {
									continue
								}

								// No XOR for SAT Solvers other than CryptoMiniSAT and XNFSAT
								if xorOption == 1 && satSolver != constants.ArgCryptoMiniSat && satSolver != constants.ArgXnfSat {
									xorOption = 0
								}

								cb(i, satSolver, steps, hash, xorOption, adderType, dobbertin, dobbertinBits)
								i += 1
							}
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
	default:
		return ""
	}
}

func InstancesCount(commandContext *types.CommandContext) uint {
	var count uint = 0
	LoopThroughVariations(commandContext, func(u1 uint, s1 string, u2 uint, s2 string, u3 uint, s3 string, u4, u5 uint) {
		count++
	})

	return count
}

func AggregateLogs() {
	AppendBenchmarkLog([]string{"SAT Solver", "Instance Name", "Time", "Exit Code"})
	AppendValidResultsLog([]string{"SAT Solver", "Instance Name", "Time", "Exit Code"})
	AppendVerificationLog([]string{"SAT Solver", "Instance Name", "Result"})

	items, _ := os.ReadDir(constants.LogsDirPath)
	for _, item := range items {
		if item.IsDir() || path.Ext(item.Name()) != ".log" || lo.Contains([]string{"benchmark.log", "verification.log", "valid_results.log"}, item.Name()) {
			continue
		}

		// Open the file as CSV
		filePath := path.Join(constants.LogsDirPath, item.Name())
		fileReader, err := os.Open(filePath)
		if err != nil {
			continue
		}

		csvReader := csv.NewReader(fileReader)
		record, err := csvReader.Read()
		if err != nil {
			continue
		}
		fmt.Println(record)

		fileReader.Close()

		fmt.Println(item.Name())
		if strings.HasPrefix(item.Name(), "verification") {
			AppendVerificationLog(record)
		} else if strings.HasPrefix(item.Name(), "benchmark") {
			AppendBenchmarkLog(record)
		} else if strings.HasPrefix(item.Name(), "valid_results") {
			AppendValidResultsLog(record)
		}
	}

	// Remove the individual logs
	for _, item := range items {
		if lo.Contains([]string{"benchmark.log", "verification.log", "valid_results.log"}, item.Name()) {
			continue
		}

		exec.Command("bash", "-c", fmt.Sprintf("rm %s%s", constants.LogsDirPath, item.Name())).Run()
	}
}

func EncodingsFileName(steps uint, adderType string, xorOption uint, hash string, dobbertin, dobbertinBits uint) string {
	return fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, InstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits))
}

func InstanceName(steps uint, adderType string, xorOption uint, hash string, dobbertin, dobbertinBits uint) string {
	return fmt.Sprintf("md4_%d_%s_xor%d_%s_dobbertin%d_b%d", steps, adderType, xorOption, hash, dobbertin, dobbertinBits)
}