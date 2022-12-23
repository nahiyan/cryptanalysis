package services

import (
	configSvc "benchmark/internal/config/services"

	"github.com/samber/do"
)

const (
	SaeedE = "saeed_e" // Short for Saeed's Encoder
)

type Name string

type EncoderService struct {
	configSvc *configSvc.ConfigService

	Name Name
}

func NewEncoderService(i *do.Injector) *EncoderService {
	configSvc := do.MustInvoke[*configSvc.ConfigService](i)
	return &EncoderService{configSvc: configSvc}
}

func (encoderSvc *EncoderService) InvokeSaeedE() {
	// // * 1. Check if the executable exists
	// if !utils.FileExists(config.Get().Paths.Bin.Encoder) {
	// 	log.Fatal("Failed to find the encoder in the 'encoders/saeed/crypto' directory. Can you ensure that you compiled it?")
	// }

	// // * 2. Loop through the variations
	// for _, hash := range context.VariationsHashes {
	// 	for _, xorOption := range context.VariationsXor {
	// 		xorFlag := func(xorOption uint) string {
	// 			if xorOption == 1 {
	// 				return " --xor"
	// 			}

	// 			return ""
	// 		}(xorOption)

	// 		for _, adderType := range context.VariationsAdders {
	// 			for _, steps := range context.VariationsSteps {
	// 				for _, isDobbertinEnabled_ := range context.VariationsDobbertin {
	// 					for _, dobbertinRelaxationBits := range context.VariationsDobbertinBits {
	// 						isDobbertinEnabled := isDobbertinEnabled_ == 1 && steps >= 27

	// 						if !isDobbertinEnabled && dobbertinRelaxationBits != 32 {
	// 							continue
	// 						}

	// 						dobbertinFlag := func(enabled bool) string {
	// 							if enabled {
	// 								return " --dobbertin"
	// 							}

	// 							return ""
	// 						}(isDobbertinEnabled)

	// 						instanceName := utils.InstanceName(steps, adderType, xorOption, hash, func(x bool) uint {
	// 							if x {
	// 								return 1
	// 							}

	// 							return 0
	// 						}(isDobbertinEnabled), dobbertinRelaxationBits, nil)

	// 						// * 3. Drive the encoder if the encoding doesn't exist
	// 						encodingFilePath := fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, instanceName)
	// 						if !utils.FileExists(encodingFilePath) {
	// 							cmd := exec.Command("bash", "-c", fmt.Sprintf("%s%s -A %s -r %d -f md4 -a preimage -t %s%s --bits %d > %s", config.Get().Paths.Bin.Encoder, xorFlag, utils.ResolveAdderType(adderType), steps, hash, dobbertinFlag, dobbertinRelaxationBits, encodingFilePath))
	// 							if err := cmd.Run(); err != nil {
	// 								log.Fatal("Failed to generate encodings for ", instanceName)
	// 							}
	// 						} else {
	// 							fmt.Println("Encoding already exists")
	// 						}

	// 						// * 4. Conditional: Generate cubes if flagged
	// 						if context.CubeParams != nil {
	// 							// Generate cubes
	// 							generateCubes(instanceName, context.CubeParams.CutoffVars)
	// 						}
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

func (encoderSvc *EncoderService) Run(name Name) {
	switch name {
	case SaeedE:
		encoderSvc.InvokeSaeedE()
	}
}
