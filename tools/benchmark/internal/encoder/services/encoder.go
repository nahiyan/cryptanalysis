package services

import (
	configSvc "benchmark/internal/config/services"
	"benchmark/internal/filesystem"
	"benchmark/internal/pipeline"
	"fmt"
	"log"
	"os/exec"

	"github.com/samber/do"
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

type Name string

type EncoderService struct {
	configSvc     *configSvc.ConfigService
	filesystemSvc filesystem.FilesystemService

	Name Name
}

func NewEncoderService(i *do.Injector) *EncoderService {
	configSvc := do.MustInvoke[*configSvc.ConfigService](i)
	return &EncoderService{configSvc: configSvc}
}

func (encoderSvc *EncoderService) GetInstanceName(steps int, adderType pipeline.AdderType, xor int, hash string, dobbertin, dobbertinBits int, cubeIndex *int) string {
	return fmt.Sprintf("%smd4_%d_%s_xor%d_%s_dobbertin%d_b%d", func(cubeIndex *int) string {
		if cubeIndex != nil {
			return fmt.Sprintf("cube%d_", *cubeIndex)
		}

		return ""
	}(cubeIndex), steps, adderType, xor, hash, dobbertin, dobbertinBits)
}

func (encoderSvc *EncoderService) LoopThroughVariations(variations pipeline.Variation, cb func(int, pipeline.Solver, int, string, int, pipeline.AdderType, int, int)) {
	for _, satSolver := range variations.Solvers {
		i := 0
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

								// No XOR for SAT Solvers other than CryptoMiniSAT and XNFSAT
								if xorOption == 1 && satSolver != CryptoMiniSat && satSolver != XnfSat {
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

func (encoderSvc *EncoderService) InvokeSaeedE(variations pipeline.Variation) {
	config := &encoderSvc.configSvc.Config

	// * Loop through the variations
	encoderSvc.LoopThroughVariations(variations, func(i int, solver pipeline.Solver, steps int, hash string, xorOption int, adderType pipeline.AdderType, dobbertin, dobbertinBits int) {
		instanceName := encoderSvc.GetInstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, nil)

		encodingFilePath := fmt.Sprintf("%s%s.cnf", EncodingsDirPath, instanceName)

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
		cmd := exec.Command("bash", "-c",
			fmt.Sprintf(
				"%s%s -A %s -r %d -f md4 -a preimage -t %s%s --bits %d > %s",
				config.Paths.Bin.SaeedE,
				xorFlag,
				encoderSvc.ResolveSaeedEAdderType(adderType),
				steps,
				hash,
				dobbertinFlag,
				dobbertinBits,
				encodingFilePath))
		if err := cmd.Run(); err != nil {
			log.Fatal("Encoding generation failed:", instanceName)
		}
	})
}

func (encoderSvc *EncoderService) Run(name Name, pipe pipeline.Pipe) {
	filesystemSvc := encoderSvc.filesystemSvc
	config := &encoderSvc.configSvc.Config

	switch name {
	case SaeedE:
		// * Check if the encoder exists
		if !filesystemSvc.FileExists(config.Paths.Bin.SaeedE) {
			log.Fatal("Failed to find the encoder in the 'encoders/saeed/crypto' directory. Can you ensure that you compiled it?")
		}

		encoderSvc.InvokeSaeedE(pipe.Variation)
	}
}
