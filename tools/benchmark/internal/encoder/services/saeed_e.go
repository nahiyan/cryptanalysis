package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"fmt"
	"log"
	"path"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func (encoderSvc *EncoderService) InvokeSaeedE(parameters pipeline.Encoding) []pipeline.EncodingPromise {
	config := &encoderSvc.configSvc.Config
	filesystemSvc := encoderSvc.filesystemSvc

	// Check if the encoder exists
	if !filesystemSvc.FileExists(config.Paths.Bin.SaeedE) {
		log.Fatal("Failed to find the encoder in the path '" + config.Paths.Bin.SaeedE + "'. Can you ensure that you compiled it?")
	}

	encodings := []string{}

	// * Loop through the variations
	encoderSvc.LoopThroughVariation(parameters, func(instanceInfo encoder.InstanceInfo) {
		instanceName := encoderSvc.GetInstanceName(instanceInfo)
		encodingFilePath := path.Join(config.Paths.Encodings, instanceName)
		encodings = append(encodings, encodingFilePath)

		// Skip if encoding already exists
		if encoderSvc.filesystemSvc.FileExists(encodingFilePath) {
			logrus.Println("Encoder: skipped", encodingFilePath)
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
			config.Paths.Bin.SaeedE,
			instanceInfo.Steps,
			instanceInfo.AdderType,
			instanceInfo.TargetHash,
			dobbertinFlag,
			dobbertinBits,
			xorFlag)
		cmd := encoderSvc.commandSvc.Create(command)
		encoderSvc.OutputToFile(cmd, encodingFilePath)

		logrus.Println("Encoder:", encodingFilePath)
	})

	encodingPromises := lo.Map(encodings, func(encoding string, _ int) pipeline.EncodingPromise {
		return EncodingPromise{Encoding: encoding}
	})
	return encodingPromises
}
