package services

import (
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/pipeline"
	"cryptanalysis/internal/simplifier"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/samber/mo"
)

func getInstanceName(info encoder.InstanceInfo) string {
	encoder_ := info.Encoder
	function := info.Function
	steps := info.Steps
	adderType := info.AdderType
	targetHash := info.TargetHash

	instanceName := fmt.Sprintf("%s_%s_%d_%s",
		encoder_,
		function,
		steps,
		targetHash)
	if dobbertinInfo, enabled := info.Dobbertin.Get(); enabled {
		instanceName += fmt.Sprintf("_dobbertin%d", dobbertinInfo.Bits)
	}
	if encoder_ == encoder.NejatiEncoder {
		instanceName += "_" + string(adderType)

		if info.IsXorEnabled {
			instanceName += "_xor"
		}
	}
	instanceName += ".cnf"

	if simplificationInfo, exists := info.Simplification.Get(); exists {
		instanceName += "." + simplificationInfo.Simplifier
		if simplificationInfo.Simplifier == simplifier.Cadical {
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

func processInstanceName(instanceName string) (encoder.InstanceInfo, error) {
	instanceName = path.Base(instanceName)
	info := encoder.InstanceInfo{}
	errInvalidFormat := errors.New("instance name processor: invalid format")

	segments := strings.Split(instanceName, ".")
	if len(segments) == 0 {
		return info, errInvalidFormat
	}

	// Processing the main segment
	{
		mainSegment := segments[0]

		// Encoder
		if strings.HasPrefix(mainSegment, encoder.Transalg) {
			info.Encoder = encoder.Transalg
		} else if strings.HasPrefix(mainSegment, encoder.NejatiEncoder) {
			info.Encoder = encoder.NejatiEncoder
		}
		mainSegment = strings.TrimPrefix(mainSegment, string(info.Encoder)+"_")

		// Function
		if strings.HasPrefix(mainSegment, encoder.Md4) {
			info.Function = encoder.Md4
		} else if strings.HasPrefix(mainSegment, encoder.Sha256) {
			info.Function = encoder.Sha256
		}
		mainSegment = strings.TrimPrefix(mainSegment, string(info.Function)+"_")

		// Steps
		steps, err := strconv.Atoi(strings.Split(mainSegment, "_")[0])
		if err != nil {
			return info, errInvalidFormat
		}
		mainSegment = strings.TrimPrefix(mainSegment, fmt.Sprintf("%d_", steps))
		info.Steps = steps

		// Adder type
		if strings.HasPrefix(mainSegment, encoder.Espresso) {
			info.AdderType = encoder.Espresso
		} else if strings.HasPrefix(mainSegment, encoder.DotMatrix) {
			info.AdderType = encoder.DotMatrix
		} else if strings.HasPrefix(mainSegment, encoder.CounterChain) {
			info.AdderType = encoder.CounterChain
		} else if strings.HasPrefix(mainSegment, encoder.TwoOperand) {
			info.AdderType = encoder.TwoOperand
		}
		mainSegment = strings.TrimPrefix(mainSegment, string(info.AdderType)+"_")

		// Target hash
		info.TargetHash = strings.Split(mainSegment, "_")[0]
		mainSegment = strings.TrimPrefix(mainSegment, info.TargetHash+"_")

		// Dobbertin
		if index := strings.Index(mainSegment, "dobbertin"); index != -1 {
			bits := strings.Split(mainSegment[index+len("dobbertin"):], "_")[0]
			bits_, err := strconv.Atoi(bits)
			if err != nil {
				return info, errInvalidFormat
			}
			info.Dobbertin = mo.Some(encoder.DobbertinInfo{
				Bits: bits_,
			})
		}

		// Xor
		if index := strings.Index(mainSegment, "xor"); index != -1 {
			info.IsXorEnabled = true
		}
	}

	// Cube info
	if index := strings.Index(instanceName, ".march_n"); index != -1 {
		threshold := strings.Split(instanceName[index+len(".march_n"):], ".")[0]
		threshold_, err := strconv.Atoi(threshold)
		if err != nil {
			return info, err
		}

		info.Cubing = mo.Some(encoder.CubingInfo{
			Threshold: threshold_,
		})
	}

	// Cube index
	if index := strings.Index(instanceName, ".cubes.cube"); index != -1 {
		cubeIndex := strings.Split(instanceName[index+len(".cubes.cube"):], ".")[0]
		cubeIndex_, err := strconv.Atoi(cubeIndex)
		if err != nil {
			return info, err
		}

		info.CubeIndex = mo.Some(cubeIndex_)
	}

	// CaDiCaL Simplification
	if index := strings.Index(instanceName, ".cadical_c"); index != -1 {
		conflicts := strings.Split(instanceName[index+len(".cadical_c"):], ".")[0]
		conflicts_, err := strconv.Atoi(conflicts)
		if err != nil {
			return info, err
		}

		info.Simplification = mo.Some(encoder.SimplificationInfo{
			Simplifier: simplifier.Cadical,
			Conflicts:  conflicts_,
		})
	}

	return info, nil
}

func loopThroughVariation(params pipeline.EncodeParams, cb func(instanceInfo encoder.InstanceInfo)) {
	instances := make([]encoder.InstanceInfo, 0)
	for _, stepsCount := range params.Steps {
		instances = append(instances, encoder.InstanceInfo{
			Encoder:    params.Encoder,
			Function:   params.Function,
			AttackType: params.AttackType,
			Steps:      stepsCount,
		})
	}

	// Add the hashes
	if params.AttackType == encoder.Preimage {
		if len(params.Hashes) == 0 {
			log.Fatal("Encoder: expected target hashes in the schema")
		}

		for i, instance := range instances {
			for j, targetHash := range params.Hashes {
				if j == 0 {
					instances[i].TargetHash = targetHash
				} else {
					newInstance := instance
					newInstance.TargetHash = targetHash
					instances = append(instances, newInstance)
				}
			}
		}
	}

	// Add the XOR versions
	if params.Encoder == encoder.NejatiEncoder {
		for i, instance := range instances {
			for j, xorEnabled := range params.Xor {
				if j == 0 {
					instances[i].IsXorEnabled = xorEnabled == 1
				} else {
					newInstance := instance
					newInstance.IsXorEnabled = xorEnabled == 1
					instances = append(instances, newInstance)
				}
			}
		}
	}

	// Add the adder types
	if params.Encoder == encoder.NejatiEncoder {
		for i, instance := range instances {
			for j, adderType := range params.Adders {
				if j == 0 {
					instances[i].AdderType = adderType
				} else {
					instance.AdderType = adderType
					instances = append(instances, instance)
				}
			}
		}
	}

	// Handle the dobbertin attack
	if (params.Function == encoder.Md4 || params.Function == encoder.Md5) && params.AttackType == encoder.Preimage {
		// Add the dobbertin options
		for i, instance := range instances {
			for j, isDobbertinEnabled := range params.Dobbertin {
				dobbertinParams := mo.None[encoder.DobbertinInfo]()
				if isDobbertinEnabled == 1 {
					dobbertinParams = mo.Some(encoder.DobbertinInfo{})
				}

				if j == 0 {
					instances[i].Dobbertin = dobbertinParams
				} else {
					newInstance := instance
					newInstance.Dobbertin = dobbertinParams
					instances = append(instances, newInstance)
				}
			}
		}

		// Add the dobbertin bits
		for i, instance := range instances {
			var dobbertinParams encoder.DobbertinInfo
			var exists bool
			if dobbertinParams, exists = instance.Dobbertin.Get(); !exists {
				continue
			}

			for j, bits := range params.DobbertinBits {
				if j == 0 {
					dobbertinParams.Bits = bits
					instances[i].Dobbertin = mo.Some(dobbertinParams)
				} else {
					newInstance := instance
					dobbertinParams.Bits = bits
					newInstance.Dobbertin = mo.Some(dobbertinParams)
					instances = append(instances, newInstance)
				}
			}
		}
	}

	for _, instance := range instances {
		cb(instance)
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

func (encoderSvc *EncoderService) Run(parameters pipeline.EncodeParams) []encoder.Encoding {
	err := encoderSvc.filesystemSvc.PrepareDir(encoderSvc.configSvc.Config.Paths.Encodings)
	encoderSvc.errorSvc.Fatal(err, "Encoder: failed to prepare directory for storing the encodings")

	if parameters.Function != encoder.Md4 && parameters.Function != encoder.Md5 && parameters.Function != encoder.Sha256 {
		log.Fatal("Encoder: function not supported")
	}

	switch parameters.Encoder {
	case encoder.NejatiEncoder:
		encodings := encoderSvc.InvokeNejatiEncoder(parameters)
		log.Println("Encoder: nejati_encoder", encodings)
		return encodings
	case encoder.Transalg:
		encodings := encoderSvc.InvokeTransalg(parameters)
		log.Println("Encoder: transalg", encodings)
		return encodings
	}

	panic("Encoder not found: " + parameters.Encoder)
}
