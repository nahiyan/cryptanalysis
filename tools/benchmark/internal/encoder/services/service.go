package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
)

const (
	SaeedE           = "saeed_e" // Short for Saeed's Encoder
	ResultsDirPath   = "./results/"
	EncodingsDirPath = ResultsDirPath + "encodings/"
)

const (
	TwoOperand   = "two_operand"
	DotMatrix    = "dot_matrix"
	CounterChain = "counter_chain"
	Espresso     = "espresso"
)

const (
	CryptoMiniSat = "cryptominisat"
	Cadical       = "cadical"
	Kissat        = "kissat"
	MapleSat      = "maplesat"
	Glucose       = "glucose"
	XnfSat        = "xnfsat"
)

func (encoderSvc *EncoderService) GetInstanceName(steps int, adderType pipeline.AdderType, xor int, hash string, dobbertin, dobbertinBits int, cubeIndex *int) string {
	return fmt.Sprintf("%smd4_%d_%s_xor%d_%s_dobbertin%d_b%d", func(cubeIndex *int) string {
		if cubeIndex != nil {
			return fmt.Sprintf("cube%d_", *cubeIndex)
		}

		return ""
	}(cubeIndex), steps, adderType, xor, hash, dobbertin, dobbertinBits)
}

func (encoderSvc *EncoderService) LoopThroughVariation(variations pipeline.Variation, cb func(int, string, int, pipeline.AdderType, int, int)) {
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

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	err = cmd.Start()
	errorSvc.Fatal(err, "Encoding generation failed: "+instanceName)

	err = filesystemSvc.WriteFromPipe(pipe, filePath)
	errorSvc.Fatal(err, "Encoding generation failed: "+instanceName)

	err = cmd.Wait()
	errorSvc.Fatal(err, "Encoding generation failed: "+instanceName)

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

func (encoderSvc *EncoderService) InvokeSaeedE(variations pipeline.Variation) []string {
	config := &encoderSvc.configSvc.Config
	filesystemSvc := encoderSvc.filesystemSvc

	// Check if the encoder exists
	if !filesystemSvc.FileExists(config.Paths.Bin.SaeedE) {
		log.Fatal("Failed to find the encoder in the '" + config.Paths.Bin.SaeedE + "' directory. Can you ensure that you compiled it?")
	}

	encodings := []string{}

	// * Loop through the variations
	encoderSvc.LoopThroughVariation(variations, func(steps int, hash string, xorOption int, adderType pipeline.AdderType, dobbertin, dobbertinBits int) {
		instanceName := encoderSvc.GetInstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, nil)

		encodingFilePath := path.Join(EncodingsDirPath, instanceName+".cnf")
		encodings = append(encodings, encodingFilePath)

		// Skip if encoding already exists
		if encoderSvc.filesystemSvc.FileExists(encodingFilePath) {
			fmt.Println("Encoding already exists:", encodingFilePath)
			return
		}

		dobbertinFlag := func(enabled int) string {
			if enabled == 1 {
				return " --dobbertin"
			}

			return ""
		}(dobbertin)

		xorFlag := func(xorOption int) string {
			if xorOption == 1 {
				return " --xor"
			}

			return ""
		}(xorOption)

		// * Drive the encoder
		cmd := exec.Command(
			config.Paths.Bin.SaeedE,
			xorFlag,
			"-A",
			string(encoderSvc.ResolveSaeedEAdderType(adderType)),
			"-r",
			strconv.Itoa(steps),
			"-f",
			"md4",
			"-a",
			"preimage",
			"-t",
			hash,
			dobbertinFlag,
			"--bits",
			strconv.Itoa(dobbertinBits))
		encoderSvc.OutputToFile(cmd, encodingFilePath)
	})

	return encodings
}

func (encoderSvc *EncoderService) TestRun() []string {
	pipe := pipeline.Pipe{
		Variation: pipeline.Variation{
			Xor:           []int{0},
			Dobbertin:     []int{0},
			DobbertinBits: []int{32},
			Adders:        []pipeline.AdderType{Espresso},
			Hashes:        []string{"ffffffffffffffffffffffffffffffff", "00000000000000000000000000000000"},
			Steps:         []int{16},
			Solvers:       []pipeline.Solver{Kissat},
		},
	}

	return encoderSvc.Run(SaeedE, pipe)
}

func (encoderSvc *EncoderService) Run(name encoder.Name, pipe pipeline.Pipe) []string {
	switch name {
	case SaeedE:
		return encoderSvc.InvokeSaeedE(pipe.Variation)
	}

	panic("Encoder not found: " + name)
}