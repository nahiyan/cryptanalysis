package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"
)

// Paths
const (
	SaeedE = "saeed_e" // Short for Saeed's Encoder
)

// Adders
const (
	TwoOperand   = "two_operand"
	DotMatrix    = "dot_matrix"
	CounterChain = "counter_chain"
	Espresso     = "espresso"
)

type EncodingPromise struct {
	Encoding string
}

func (encodingPromise EncodingPromise) Get() string {
	return encodingPromise.Encoding
}

func (encoderSvc *EncoderService) GetInstanceName(steps int, adderType pipeline.AdderType, xor int, hash string, dobbertin, dobbertinBits int, cubeIndex *int) string {
	return fmt.Sprintf("%smd4_%d_%s_xor%d_%s_dobbertin%d_b%d", func(cubeIndex *int) string {
		if cubeIndex != nil {
			return fmt.Sprintf("cube%d_", *cubeIndex)
		}

		return ""
	}(cubeIndex), steps, adderType, xor, hash, dobbertin, dobbertinBits)
}

func (encoderSvc *EncoderService) LoopThroughVariation(variations pipeline.Encoding, cb func(int, string, int, pipeline.AdderType, int, int)) {
	for _, steps := range variations.Steps {
		for _, hash := range variations.Hashes {
			for _, xorOption := range variations.Xor {
				for _, adderType := range variations.Adders {
					for _, dobbertin := range variations.Dobbertin {
						for _, dobbertinBits := range variations.DobbertinBits {
							// Skip any dobbertin bit variation when dobbertin's attack isn't on
							if dobbertin == 0 && dobbertinBits != 32 {
								continue
							}

							// Skip dobbertin's attacks when steps count < 27
							if steps < 27 && dobbertin == 1 {
								continue
							}

							cb(steps, hash, xorOption, adderType, dobbertin, dobbertinBits)
						}
					}
				}
			}
		}
	}
}

func (encoderSvc *EncoderService) OutputToFile(cmd *exec.Cmd, filePath string) {
	filesystemSvc := encoderSvc.filesystemSvc
	errorSvc := encoderSvc.errorSvc
	instanceName := path.Base(filePath)
	failureMsg := "Encoding generation failed: " + instanceName

	pipe, err := cmd.StdoutPipe()
	errorSvc.Fatal(err, failureMsg)

	err = cmd.Start()
	errorSvc.Fatal(err, failureMsg)

	err = filesystemSvc.WriteFromPipe(pipe, filePath)
	errorSvc.Fatal(err, failureMsg)

	err = cmd.Wait()
	errorSvc.Fatal(err, failureMsg)
}

func (encoderSvc *EncoderService) ResolveSaeedEAdderType(adderType pipeline.AdderType) pipeline.AdderType {
	switch adderType {
	case CounterChain:
		return "counter_chain"
	case DotMatrix:
		return "dot_matrix"
	case Espresso:
		return "espresso"
	case TwoOperand:
		return "two_operand"
	default:
		return ""
	}
}

func (encoderSvc *EncoderService) InvokeSaeedE(parameters pipeline.Encoding) []string {
	config := &encoderSvc.configSvc.Config
	filesystemSvc := encoderSvc.filesystemSvc

	// Check if the encoder exists
	if !filesystemSvc.FileExists(config.Paths.Bin.SaeedE) {
		log.Fatal("Failed to find the encoder in the path '" + config.Paths.Bin.SaeedE + "'. Can you ensure that you compiled it?")
	}

	encodings := []string{}

	// * Loop through the variations
	encoderSvc.LoopThroughVariation(parameters, func(steps int, hash string, xorOption int, adderType pipeline.AdderType, dobbertin, dobbertinBits int) {
		instanceName := encoderSvc.GetInstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, nil)

		encodingFilePath := path.Join(parameters.OutputDir, instanceName+".cnf")
		encodings = append(encodings, encodingFilePath)

		// Skip if encoding already exists
		if encoderSvc.filesystemSvc.FileExists(encodingFilePath) {
			fmt.Println("Encoder: skipped", encodingFilePath)
			return
		}

		dobbertinFlag := func(enabled int) string {
			if enabled == 1 {
				return "--dobbertin"
			}

			return ""
		}(dobbertin)

		xorFlag := func(xorOption int) string {
			if xorOption == 1 {
				return "--xor"
			}

			return ""
		}(xorOption)

		// * Drive the encoder
		cmd := exec.Command(
			config.Paths.Bin.SaeedE,
			strings.Fields(
				fmt.Sprintf(
					"-A %s -r %d -f md4 -a preimage -t %s --bits %d %s %s",
					string(encoderSvc.ResolveSaeedEAdderType(adderType)),
					steps,
					hash,
					dobbertinBits,
					xorFlag,
					dobbertinFlag))...)
		encoderSvc.OutputToFile(cmd, encodingFilePath)

		fmt.Println("Encoder:", encodingFilePath)
	})

	return encodings
}

func (encoderSvc *EncoderService) Run(name encoder.Name, parameters pipeline.Encoding) []string {
	switch name {
	case SaeedE:
		encodings := encoderSvc.InvokeSaeedE(parameters)
		fmt.Println("Encoder:", encodings)
		return encodings
	}

	panic("Encoder not found: " + name)
}
