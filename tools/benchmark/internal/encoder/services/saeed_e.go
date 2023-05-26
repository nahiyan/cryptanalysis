package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"fmt"
	"log"
	"path"
)

func (encoderSvc *EncoderService) InvokeNejatiEncoder(parameters pipeline.EncodeParams) []encoder.Encoding {
	config := &encoderSvc.configSvc.Config

	// * Loop through the variations
	encodings := []encoder.Encoding{}
	encoderSvc.LoopThroughVariation(parameters, func(instanceInfo encoder.InstanceInfo) {
		instanceName := encoderSvc.GetInstanceName(instanceInfo)
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
		command := fmt.Sprintf(
			"%s -f md4 -a preimage -r %d -A %s -t %s %s --bits %d %s",
			config.Paths.Bin.NejatiEncoder,
			instanceInfo.Steps,
			instanceInfo.AdderType,
			instanceInfo.TargetHash,
			dobbertinFlag,
			dobbertinBits,
			xorFlag)
		cmd := encoderSvc.commandSvc.Create(command)
		encoderSvc.OutputToFile(cmd, encodingPath)

		log.Println("Encoder:", encodingPath)
	})
	return encodings
}
