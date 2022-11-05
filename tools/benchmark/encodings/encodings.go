package encodings

import (
	"benchmark/constants"
	"benchmark/types"
	"benchmark/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Generate(context types.EncodingsGenContext) {
	// * 1. Check if the executable exists
	if _, err := os.Stat(constants.EncoderPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Failed to find the encoder in the 'encoders/saeed/crypto' directory. Can you ensure that you compiled it?")
	}

	// * 2. Loop through the variations
	for _, hash := range context.VariationsHashes {
		for _, xorOption := range context.VariationsXor {
			xorFlag := func(xorOption uint) string {
				if xorOption == 1 {
					return " --xor"
				}

				return ""
			}(xorOption)

			for _, adderType := range context.VariationsAdders {
				for _, steps := range context.VariationsSteps {
					for _, isDobbertinEnabled_ := range context.VariationsDobbertin {
						for _, dobbertinRelaxationBits := range context.VariationsDobbertinBits {
							isDobbertinEnabled := isDobbertinEnabled_ == 1 && steps >= 27

							if !isDobbertinEnabled && dobbertinRelaxationBits != 32 {
								continue
							}

							dobbertinFlag := func(enabled bool) string {
								if enabled {
									return " --dobbertin"
								}

								return ""
							}(isDobbertinEnabled)

							instanceName := utils.InstanceName(steps, adderType, xorOption, hash, func(x bool) uint {
								if x {
									return 1
								}

								return 0
							}(isDobbertinEnabled), dobbertinRelaxationBits)

							// * 3. Drive the encoder
							cmd := exec.Command("bash", "-c", fmt.Sprintf("%s%s -A %s -r %d -f md4 -a preimage -t %s%s --bits %d > %sencodings/%s.cnf", constants.EncoderPath, xorFlag, utils.ResolveAdderType(adderType), steps, hash, dobbertinFlag, dobbertinRelaxationBits, constants.ResultsDirPath, instanceName))
							if err := cmd.Run(); err != nil {
								log.Fatal("Failed to generate encodings for ", instanceName)
							}

							// * 4. Conditional: Generate cubes if flagged
							if context.IsCubeEnabled {
								if err := generateCubes(instanceName); err != nil {
									log.Fatal("Failed to generate cubes")
								}
							}
						}
					}
				}
			}
		}
	}
}
