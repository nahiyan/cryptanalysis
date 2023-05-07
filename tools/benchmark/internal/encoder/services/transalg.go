package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"text/template"

	_ "embed"

	"github.com/samber/lo"
)

//go:embed transalg_md4.txt
var layoutMd4 string

//go:embed transalg_md5.txt
var layoutMd5 string

//go:embed transalg_sha256.txt
var layoutSha256 string

func (encoderSvc *EncoderService) GenerateTransalgMd4Code(instanceInfo encoder.InstanceInfo) (string, error) {
	dobbertinSteps := []int{
		13, 14, 15,
		17, 18, 19,
		21, 22, 23,
		25, 26, 27,
	}
	dobbertinInfo, dobbertinAttackEnabled := instanceInfo.Dobbertin.Get()

	tmpl := template.New("transalg_md4.txt").Funcs(map[string]interface{}{
		"step": func(step int, body string) string {
			if step <= instanceInfo.Steps {
				dobbertinConstraint := ""
				if dobbertinAttackEnabled && lo.Contains(dobbertinSteps, step) {
					registers := []byte{'a', 'd', 'c', 'b'}
					register := registers[(step-1)%4]

					if step == 13 && dobbertinInfo.Bits < 32 {
						for i := 0; i < dobbertinInfo.Bits; i++ {
							dobbertinConstraint += fmt.Sprintf("\n\tassert(!(%c[%d] ^ K[%d]));", register, i, i)
						}
					} else {
						dobbertinConstraint = fmt.Sprintf("\n\tassert(!(%c ^ K));", register)
					}
				}

				return body + fmt.Sprintf(" // Step %d", step) + dobbertinConstraint
			}
			return ""
		},
	})
	tmpl, err := tmpl.Parse(layoutMd4)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	tmpl.Execute(&buffer, map[string]interface{}{
		"Steps":             instanceInfo.Steps,
		"DobbertinConstant": math.MaxUint32,
		"OneTargetHash":     instanceInfo.TargetHash == "ffffffffffffffffffffffffffffffff",
	})
	code := buffer.String()

	return code, nil
}

func (encoderSvc *EncoderService) GenerateTransalgMd5Code(instanceInfo encoder.InstanceInfo) (string, error) {
	tmpl := template.New("transalg_md5.txt").Funcs(map[string]interface{}{
		"step": func(step int, body string) string {
			if step <= instanceInfo.Steps {
				return body + fmt.Sprintf(" // Step %d", step)
			}
			return ""
		},
	})
	tmpl, err := tmpl.Parse(layoutMd5)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	tmpl.Execute(&buffer, map[string]interface{}{
		"Steps":         instanceInfo.Steps,
		"OneTargetHash": instanceInfo.TargetHash == "ffffffffffffffffffffffffffffffff",
	})
	code := buffer.String()

	return code, nil
}

func (encoderSvc *EncoderService) GenerateTransalgSha256Code(instanceInfo encoder.InstanceInfo) (string, error) {
	tmpl := template.New("transalg_sha256.txt")
	tmpl, err := tmpl.Parse(layoutSha256)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	tmpl.Execute(&buffer, map[string]interface{}{
		"Steps":         instanceInfo.Steps,
		"OneTargetHash": instanceInfo.TargetHash == "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	})
	code := buffer.String()

	return code, nil
}

func (encoderSvc *EncoderService) ShouldSkip(encodingPath string) bool {
	return encoderSvc.filesystemSvc.FileExists(encodingPath)
}

// TODO: Reduce shared redundant code with SaeedE invokation
func (encoderSvc *EncoderService) InvokeTransalg(parameters pipeline.EncodeParams) []encoder.Encoding {
	err := encoderSvc.filesystemSvc.PrepareDir("tmp")
	encoderSvc.errorSvc.Fatal(err, "Encoder: failed to create tmp dir")

	encodings := []encoder.Encoding{}

	// * Loop through the variations
	encoderSvc.LoopThroughVariation(parameters, func(instanceInfo encoder.InstanceInfo) {
		instanceName := encoderSvc.GetInstanceName(instanceInfo)
		encodingPath := path.Join(encoderSvc.configSvc.Config.Paths.Encodings, instanceName)
		encodings = append(encodings, encoder.Encoding{BasePath: encodingPath})

		// Skip if encoding already exists
		if !parameters.Redundant && encoderSvc.ShouldSkip(encodingPath) {
			log.Println("Encoder: skipped", encodingPath)
			return
		}

		var (
			transalgCode string
			err          error
		)
		if parameters.Function == encoder.Md4 {
			transalgCode, err = encoderSvc.GenerateTransalgMd4Code(instanceInfo)
		} else if parameters.Function == encoder.Md5 {
			transalgCode, err = encoderSvc.GenerateTransalgMd5Code(instanceInfo)
		} else if parameters.Function == encoder.Sha256 {
			transalgCode, err = encoderSvc.GenerateTransalgSha256Code(instanceInfo)
		}
		encoderSvc.errorSvc.Fatal(err, "Encoder: failed to generate Transalg code")
		transalgFileName := fmt.Sprintf("%s.alg", encoderSvc.randomSvc.RandString(16))
		transalgFilePath := path.Join(encoderSvc.configSvc.Config.Paths.Tmp, transalgFileName)
		os.WriteFile(transalgFilePath, []byte(transalgCode), 0644)

		// * Drive the encoder
		command := fmt.Sprintf("%s -i %s -o %s", encoderSvc.configSvc.Config.Paths.Bin.Transalg, transalgFilePath, encodingPath)
		err = encoderSvc.commandSvc.Create(command).Run()
		defer os.Remove(transalgFilePath)
		encoderSvc.errorSvc.Fatal(err, "Encoder: failed to run Transalg for "+instanceName)

		log.Println("Encoder:", encodingPath)
	})

	return encodings
}
