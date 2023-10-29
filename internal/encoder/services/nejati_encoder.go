package services

import (
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/pipeline"
	"fmt"
	"log"
	"path"
)

func (encoderSvc *EncoderService) InvokeNejatiEncoder(parameters pipeline.EncodeParams) []encoder.Encoding {
	config := &encoderSvc.configSvc.Config

	// * Loop through the variations
	encodings := []encoder.Encoding{}
	loopThroughVariation(parameters, func(instanceInfo encoder.InstanceInfo) {
		instanceName := getInstanceName(instanceInfo)
		encodingPath := path.Join(config.Paths.Encodings, instanceName)
		encodings = append(encodings, encoder.Encoding{BasePath: encodingPath})

		// Skip if encoding already exists
		if !parameters.Redundant && encoderSvc.ShouldSkip(encodingPath) {
			log.Println("Encoder: skipped", encodingPath)
			return
		}

		dobbertinFlag := ""
		dobbertinBits := 0
		dobbertinInfo, isDobbertinEnabled := instanceInfo.Dobbertin.Get()
		if isDobbertinEnabled {
			dobbertinFlag = "--dobbertin"
			dobbertinBits = dobbertinInfo.Bits
		}

		xorFlag := ""
		if instanceInfo.IsXorEnabled {
			xorFlag = "--xor"
		}

		// * Drive the encoder
		var command string
		if parameters.AttackType == encoder.Preimage {
			command = fmt.Sprintf(
				"%s -f md4 -a preimage -r %d -A %s -t %s %s --bits %d %s",
				config.Paths.Bin.NejatiPreimageEncoder,
				instanceInfo.Steps,
				instanceInfo.AdderType,
				instanceInfo.TargetHash,
				dobbertinFlag,
				dobbertinBits,
				xorFlag)
		} else if parameters.AttackType == encoder.Collision {
			command = fmt.Sprintf(
				"%s -r %d -A %s %s --diff_desc --diff_impl -d HARD_CODED",
				config.Paths.Bin.NejatiCollisionEncoder,
				instanceInfo.Steps,
				instanceInfo.AdderType,
				xorFlag)
		} else {
			log.Fatal("Encoder: failed to recognize attack type")
		}
		cmd := encoderSvc.commandSvc.Create(command)
		encoderSvc.OutputToFile(cmd, encodingPath)

		log.Println("Encoder:", encodingPath)
	})
	return encodings
}
