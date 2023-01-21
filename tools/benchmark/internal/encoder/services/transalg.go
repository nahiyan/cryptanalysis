package services

import (
	"benchmark/internal/pipeline"
	"log"

	"github.com/samber/lo"
)

func (encoderSvc *EncoderService) InvokeTransalg(parameters pipeline.Encoding) []pipeline.EncodingPromise {
	config := &encoderSvc.configSvc.Config
	filesystemSvc := encoderSvc.filesystemSvc

	// Check if the encoder exists
	if !filesystemSvc.FileExists(config.Paths.Bin.Transalg) {
		log.Fatal("Failed to find the encoder in the path '" + config.Paths.Bin.Transalg + "'. Can you ensure that you compiled it?")
	}

	encodings := []string{}

	// * Loop through the variations
	// encoderSvc.LoopThroughVariation(parameters, func(steps int, hash string, xorOption int, adderType pipeline.AdderType, dobbertin, dobbertinBits int) {
	// 	instanceName := encoderSvc.GetInstanceName(steps, adderType, xorOption, hash, dobbertin, dobbertinBits, nil)

	// 	encodingFilePath := path.Join(parameters.OutputDir, instanceName+".cnf")
	// 	encodings = append(encodings, encodingFilePath)

	// 	// Skip if encoding already exists
	// 	if encoderSvc.filesystemSvc.FileExists(encodingFilePath) {
	// 		logrus.Println("Encoder: skipped", encodingFilePath)
	// 		return
	// 	}

	// 	dobbertinFlag := func(enabled int) string {
	// 		if enabled == 1 {
	// 			return "--dobbertin"
	// 		}

	// 		return ""
	// 	}(dobbertin)

	// 	xorFlag := func(xorOption int) string {
	// 		if xorOption == 1 {
	// 			return "--xor"
	// 		}

	// 		return ""
	// 	}(xorOption)

	// 	// * Drive the encoder
	// 	cmd := exec.Command(
	// 		config.Paths.Bin.SaeedE,
	// 		strings.Fields(
	// 			fmt.Sprintf(
	// 				"-A %s -r %d -f md4 -a preimage -t %s --bits %d %s %s",
	// 				string(encoderSvc.ResolveSaeedEAdderType(adderType)),
	// 				steps,
	// 				hash,
	// 				dobbertinBits,
	// 				xorFlag,
	// 				dobbertinFlag))...)
	// 	encoderSvc.OutputToFile(cmd, encodingFilePath)

	// 	logrus.Println("Encoder:", encodingFilePath)
	// })

	encodingPromises := lo.Map(encodings, func(encoding string, _ int) pipeline.EncodingPromise {
		return EncodingPromise{Encoding: encoding}
	})
	return encodingPromises
}
