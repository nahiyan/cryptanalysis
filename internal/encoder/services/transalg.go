package services

import (
	"bytes"
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/pipeline"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"text/template"

	_ "embed"
)

func generateEqualityAssertion(variableA, variableB string, bits int) string {
	if bits == 32 {
		return fmt.Sprintf("assert(!(%s ^ %s));", variableA, variableB)
	}

	statements := ""
	for i := 0; i < bits; i++ {
		statements += fmt.Sprintf("assert(!(%s[%d] ^ %s[%d]));\n", variableA, i, variableB, i)
	}
	return statements
}

//go:embed transalg_md4.txt
var layoutMd4 string

//go:embed transalg_md5.txt
var layoutMd5 string

//go:embed transalg_sha256.txt
var layoutSha256 string

func (encoderSvc *EncoderService) GenerateTransalgMd4Code(instanceInfo encoder.InstanceInfo) (string, error) {
	dobbertinInfo, dobbertinAttackEnabled := instanceInfo.Dobbertin.Get()
	registers := []string{"a", "d", "c", "b"}

	tmpl := template.New("transalg_md4.txt").Funcs(map[string]interface{}{
		"step": func(step int, body string) string {
			if step <= instanceInfo.Steps {
				return body + fmt.Sprintf(" // Step %d", step)
			}
			return ""
		},
		"constraints": func() string {
			if !dobbertinAttackEnabled {
				return ""
			}

			constraints := ""
			dobbertinSteps := []int{
				14, 15,
				17, 18, 19,
				21, 22, 23,
				25, 26, 27,
				// 23,
				// 25, 26, 27,
				// 29, 30, 31,
				// 33, 34, 35,
				// 37, 38, 39,
			}
			constraints += generateEqualityAssertion("a_13", "K", dobbertinInfo.Bits)
			for _, dobbertinStep := range dobbertinSteps {
				register := fmt.Sprintf("%s_%d", registers[(dobbertinStep-1)%4], dobbertinStep)
				constraints += "\n\t" + generateEqualityAssertion(register, "K", 32)
			}
			return constraints
		},
	})
	tmpl, err := tmpl.Parse(layoutMd4)
	if err != nil {
		return "", err
	}

	lastRegVar1 := "a"
	lastRegVar2 := "b"
	lastRegVar3 := "c"
	lastRegVar4 := "d"
	for i := 3; i >= 0; i-- {
		variable := registers[(instanceInfo.Steps+i)%4]
		switch variable {
		case "a":
			lastRegVar1 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		case "b":
			lastRegVar2 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		case "c":
			lastRegVar3 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		case "d":
			lastRegVar4 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		}
	}

	var buffer bytes.Buffer
	tmpl.Execute(&buffer, map[string]interface{}{
		"Steps":             instanceInfo.Steps,
		"DobbertinConstant": math.MaxUint32,
		"OneTargetHash":     instanceInfo.TargetHash == "ffffffffffffffffffffffffffffffff",
		"LastRegVar1":       lastRegVar1,
		"LastRegVar2":       lastRegVar2,
		"LastRegVar3":       lastRegVar3,
		"LastRegVar4":       lastRegVar4,
	})
	code := buffer.String()

	return code, nil
}

func (encoderSvc *EncoderService) GenerateTransalgMd5Code(instanceInfo encoder.InstanceInfo) (string, error) {
	_, dobbertinAttackEnabled := instanceInfo.Dobbertin.Get()
	registers := []string{"a", "d", "c", "b"}

	tmpl := template.New("transalg_md5.txt").Funcs(map[string]interface{}{
		"step": func(step int, body string) string {
			if step > instanceInfo.Steps {
				return ""
			}

			return body + fmt.Sprintf(" // Step %d", step)
		},
		"constraints": func() string {
			if !dobbertinAttackEnabled {
				return ""
			}

			constraints := ""
			dobbertinSteps := []int{
				// 1, 2, 3,
				// 6,
				// 11,
				13, 14, 15,
				17, 18, 19,
				21, 22, 23,
			}
			// constraints += generateEqualityAssertion("a_5", "L", 32)
			for _, dobbertinStep := range dobbertinSteps {
				register := fmt.Sprintf("%s_%d", registers[(dobbertinStep-1)%4], dobbertinStep)
				constraints += "\n\t" + generateEqualityAssertion(register, "K", 32)
			}
			return constraints
		},
	})
	tmpl, err := tmpl.Parse(layoutMd5)
	if err != nil {
		return "", err
	}

	lastRegVar1 := "a"
	lastRegVar2 := "b"
	lastRegVar3 := "c"
	lastRegVar4 := "d"
	for i := 3; i >= 0; i-- {
		variable := registers[(instanceInfo.Steps+i)%4]
		switch variable {
		case "a":
			lastRegVar1 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		case "b":
			lastRegVar2 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		case "c":
			lastRegVar3 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		case "d":
			lastRegVar4 = fmt.Sprintf("%s_%d", variable, instanceInfo.Steps-(3-i))
		}
	}

	var buffer bytes.Buffer
	tmpl.Execute(&buffer, map[string]interface{}{
		"Steps":         instanceInfo.Steps,
		"OneTargetHash": instanceInfo.TargetHash == "ffffffffffffffffffffffffffffffff",
		"LastRegVar1":   lastRegVar1,
		"LastRegVar2":   lastRegVar2,
		"LastRegVar3":   lastRegVar3,
		"LastRegVar4":   lastRegVar4,
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

// TODO: Reduce shared redundant code with NejatiEncoder invokation
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
		transalgFileName := fmt.Sprintf("%s.alg", parameters.Function)
		transalgFilePath := path.Join(encoderSvc.configSvc.Config.Paths.Tmp, transalgFileName)
		os.WriteFile(transalgFilePath, []byte(transalgCode), 0644)

		// * Drive the encoder
		command := fmt.Sprintf("%s -i %s -o %s", encoderSvc.configSvc.Config.Paths.Bin.Transalg, transalgFilePath, encodingPath)
		err = encoderSvc.commandSvc.Create(command).Run()
		// defer os.Remove(transalgFilePath)
		encoderSvc.errorSvc.Fatal(err, "Encoder: failed to run Transalg for "+instanceName)

		log.Println("Encoder:", encodingPath)
	})

	return encodings
}
