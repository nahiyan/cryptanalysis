package utils

import (
	"benchmark/constants"
	"benchmark/types"
	"bufio"
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

func MakeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func LoopThroughVariations(context *types.CommandContext, cb func(uint, string, uint, string, uint, string, uint, uint, *uint)) {
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

								// TODO: Replace the .icnf format with .cubes
								// No XOR for SAT Solvers other than CryptoMiniSAT and XNFSAT
								if xorOption == 1 && satSolver != constants.ArgCryptoMiniSat && satSolver != constants.ArgXnfSat {
									xorOption = 0
								}

								if context.CubeParams != nil {
									cubesFile, err := os.Open(fmt.Sprintf("%s%s.icnf", constants.EncodingsDirPath, InstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, nil)))
									if err != nil {
										log.Fatal("Failed to read .icnf file")
									}
									cubesCount, err := CountLines(cubesFile)
									if err != nil {
										log.Fatal("Failed to count the cubes")
									}

									if context.CubeParams.CubeIndex == 0 {
										// Shuffled list of cubes
										cubes := RandomCubes(cubesCount, func(selectionCountArg uint) int {
											if selectionCountArg == 0 {
												return cubesCount
											}

											return int(selectionCountArg)
										}(context.CubeParams.SelectionSize))

										for _, cubeIndex := range cubes {
											cb(i, satSolver, steps, hash, xorOption, adderType, dobbertin, dobbertinBits, lo.ToPtr(uint(cubeIndex)))
											i += 1
										}
									} else {
										if context.CubeParams.CubeIndex > uint(cubesCount) {
											log.Fatal("Cube doesn't exist")
										}

										cb(i, satSolver, steps, hash, xorOption, adderType, dobbertin, dobbertinBits, lo.ToPtr(context.CubeParams.CubeIndex))
										i += 1
									}
								} else {
									cb(i, satSolver, steps, hash, xorOption, adderType, dobbertin, dobbertinBits, nil)
									i += 1
								}
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
	case constants.ArgCounterChain:
		return "counter_chain"
	case constants.ArgDotMatrix:
		return "dot_matrix"
	case constants.ArgEspresso:
		return "espresso"
	case constants.ArgTwoOperand:
		return "two_operand"
	default:
		return ""
	}
}

func InstancesCount(commandContext *types.CommandContext) uint {
	var count uint = 0
	LoopThroughVariations(commandContext, func(_ uint, _ string, _ uint, _ string, _ uint, _ string, _, _ uint, _ *uint) {
		count++
	})

	return count
}

func EncodingsFileName(steps uint, adderType string, xorOption uint, hash string, dobbertin, dobbertinBits uint, cubeIndex *uint) string {
	return fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, InstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, cubeIndex))
}

func InstanceName(steps uint, adderType string, xorOption uint, hash string, dobbertin, dobbertinBits uint, cubeIndex *uint) string {
	return fmt.Sprintf("%smd4_%d_%s_xor%d_%s_dobbertin%d_b%d", func(cubeIndex *uint) string {
		if cubeIndex != nil {
			return fmt.Sprintf("cube%d_", *cubeIndex)
		}

		return ""
	}(cubeIndex), steps, adderType, xorOption, hash, dobbertin, dobbertinBits)
}

func CountLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func RandomHash() (string, error) {
	hash := sha1.New()
	if _, err := hash.Write([]byte(string(fmt.Sprintf("%d", time.Now().UnixNano())))); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash), nil
}

func GetJobId(output string) int {
	jobId, err := strconv.Atoi(strings.Split(output, " ")[3])
	if err != nil {
		return 0
	}
	return jobId
}

func ReadLine(r io.Reader, lineNum int) (line string, lastLine int, err error) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lastLine++
		if lastLine == lineNum {
			// you can return sc.Bytes() if you need output in []bytes
			return sc.Text(), lastLine, sc.Err()
		}
	}
	return line, lastLine, io.EOF
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
