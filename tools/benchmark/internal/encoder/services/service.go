package services

import (
	"benchmark/internal/consts"
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"fmt"
	"log"
	"os/exec"
	"path"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus"
)

// Encoders
const (
	SaeedE   = "saeed_e" // Short for Saeed's Encoder
	Transalg = "transalg"
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

func (encodingPromise EncodingPromise) Get(dependencies map[string]interface{}) string {
	return encodingPromise.Encoding
}

func (encodingPromise EncodingPromise) GetPath() string {
	return encodingPromise.Encoding
}

func (encoderSvc *EncoderService) GetInstanceName(info encoder.InstanceInfo) string {
	encoder := info.Encoder
	function := info.Function
	steps := info.Steps
	adderType := info.AdderType
	targetHash := info.TargetHash

	instanceName := fmt.Sprintf("%s_%s_%d_%s_%s",
		encoder,
		function,
		steps,
		adderType,
		targetHash)
	if dobbertinInfo, enabled := info.Dobbertin.Get(); enabled {
		instanceName += fmt.Sprintf("_dobbertin%d", dobbertinInfo.Bits)
	}
	if info.IsXorEnabled {
		instanceName += "_xor"
	}
	instanceName += ".cnf"

	if simplificationInfo, exists := info.Simplification.Get(); exists {
		instanceName += "." + simplificationInfo.Simplifier
		if simplificationInfo.Simplifier == consts.Cadical {
			instanceName += fmt.Sprintf("_c%d.cnf", simplificationInfo.Conflicts)
		}
	}

	if cubingInfo, exists := info.Cubing.Get(); exists {
		instanceName += fmt.Sprintf(".march_n%d.cubes", cubingInfo.Threshold)
	}

	if cubeIndex, exists := info.CubeIndex.Get(); exists {
		instanceName += fmt.Sprintf(".cube%d.cnf", cubeIndex)
	}

	return instanceName
}

func (encoderSvc *EncoderService) LoopThroughVariation(params pipeline.Encoding, cb func(instanceInfo encoder.InstanceInfo)) {
	for _, steps := range params.Steps {
		for _, hash := range params.Hashes {
			for _, xorOption := range params.Xor {
				for _, adderType := range params.Adders {
					for _, dobbertin := range params.Dobbertin {
						for _, dobbertinBits := range params.DobbertinBits {
							// Skip any dobbertin bit variation when dobbertin's attack isn't on
							if dobbertin == 0 && dobbertinBits != 32 {
								continue
							}

							// Skip dobbertin's attacks when steps count < 27
							if steps < 27 && dobbertin == 1 {
								continue
							}

							dobbertin_ := mo.None[encoder.DobbertinInfo]()
							if dobbertin == 1 {
								dobbertin_ = mo.Some(encoder.DobbertinInfo{
									Bits: dobbertinBits,
								})
							}

							cb(encoder.InstanceInfo{
								Encoder:      params.Encoder,
								Function:     "md4",
								Steps:        steps,
								TargetHash:   hash,
								AdderType:    adderType,
								IsXorEnabled: xorOption == 1,
								Dobbertin:    dobbertin_,
							})
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
		encodingFilePath := path.Join(parameters.OutputDir, instanceName)
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

func (encoderSvc *EncoderService) Run(parameters pipeline.Encoding) []pipeline.EncodingPromise {
	switch parameters.Encoder {
	case SaeedE:
		promises := encoderSvc.InvokeSaeedE(parameters)
		logrus.Println("Encoder: saeed_e", promises)
		return promises
	case Transalg:
		promises := encoderSvc.InvokeTransalg(parameters)
		logrus.Println("Encoder: transalg", promises)
		return promises
	}

	panic("Encoder not found: " + parameters.Encoder)
}
