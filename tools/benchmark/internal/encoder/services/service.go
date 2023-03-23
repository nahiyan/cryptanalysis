package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"benchmark/internal/simplifier"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/samber/mo"
)

func (encoderSvc *EncoderService) GetInstanceName(info encoder.InstanceInfo) string {
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
	if encoder_ == encoder.SaeedE {
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

func (encoderSvc *EncoderService) ProcessInstanceName(instanceName string) (encoder.InstanceInfo, error) {
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
		} else if strings.HasPrefix(mainSegment, encoder.SaeedE) {
			info.Encoder = encoder.SaeedE
		}
		mainSegment = strings.TrimPrefix(mainSegment, string(info.Encoder)+"_")

		// Function
		if strings.HasPrefix(mainSegment, "md4") {
			info.Function = "md4"
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

	// SatELite Simplification
	if index := strings.Index(instanceName, ".satelite"); index != -1 {
		info.Simplification = mo.Some(encoder.SimplificationInfo{
			Simplifier: simplifier.Satelite,
		})
	}

	return info, nil
}

func (encoderSvc *EncoderService) LoopThroughVariation(params pipeline.EncodeParams, cb func(instanceInfo encoder.InstanceInfo)) {
	for _, steps := range params.Steps {
		for _, hash := range params.Hashes {
			for _, xorOption := range params.Xor {
				for _, adderType := range params.Adders {
					if params.Function == encoder.Md4 {
						for _, dobbertin := range params.Dobbertin {
							// TODO: Ignore dobbertin stuff for MD5
							for _, dobbertinBits := range params.DobbertinBits {
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
									Function:     params.Function,
									Steps:        steps,
									TargetHash:   hash,
									AdderType:    adderType,
									IsXorEnabled: xorOption == 1,
									Dobbertin:    dobbertin_,
								})

								// Skip any following dobbertin bit variation when dobbertin's attack isn't on
								if dobbertin == 0 {
									break
								}
							}
						}
					} else if params.Function == encoder.Md5 {
						cb(encoder.InstanceInfo{
							Encoder:      params.Encoder,
							Function:     params.Function,
							Steps:        steps,
							TargetHash:   hash,
							AdderType:    adderType,
							IsXorEnabled: xorOption == 1,
						})
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

func (encoderSvc *EncoderService) Run(parameters pipeline.EncodeParams) []encoder.Encoding {
	err := encoderSvc.filesystemSvc.PrepareDir(encoderSvc.configSvc.Config.Paths.Encodings)
	encoderSvc.errorSvc.Fatal(err, "Encoder: failed to prepare directory for storing the encodings")

	// TODO: Add MD5 to SaeedE
	if parameters.Function != encoder.Md4 && parameters.Function != encoder.Md5 {
		log.Fatal("Encoder: function not supported")
	}

	switch parameters.Encoder {
	case encoder.SaeedE:
		encodings := encoderSvc.InvokeSaeedE(parameters)
		log.Println("Encoder: saeed_e", encodings)
		return encodings
	case encoder.Transalg:
		encodings := encoderSvc.InvokeTransalg(parameters)
		log.Println("Encoder: transalg", encodings)
		return encodings
	}

	panic("Encoder not found: " + parameters.Encoder)
}
