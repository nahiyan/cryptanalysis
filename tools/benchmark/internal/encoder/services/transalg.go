package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"bytes"
	"fmt"
	"math"
	"os"
	"path"
	"text/template"

	_ "embed"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

//go:embed transalg.txt
var layout string

// TODO: Add arg dobbertinConstant
func (encoderSvc *EncoderService) GenerateTransalgCode(instanceInfo encoder.InstanceInfo, dobbertinConstant uint32) (string, error) {
	tmpl := template.New("transalg.txt").Funcs(map[string]interface{}{
		"inc": func(i int) int {
			return i + 1
		},
		"quo": func(i int, q int) int {
			return i % q
		},
		"add": func(i int, b int) int {
			return i + b
		},
		"function": func(i int) string {
			if i < 16 {
				return "FF"
			}
			if i < 32 {
				return "GG"
			}

			return "HH"
		},
		"dobbertinsConstraint": func(i int, register string) string {
			dobbertinIndices := []int{
				12, 13, 14,
				16, 17, 18,
				20, 21, 22,
				24, 25, 26,
			}

			if _, exists := lo.Find(dobbertinIndices, func(index int) bool {
				return i == index
			}); exists && i != 12 {
				return fmt.Sprintf("\n\tassert(!(%s ^ K));", register)
			}

			dobbertinInfo, exists := instanceInfo.Dobbertin.Get()
			if !exists {
				return ""
			}

			if i != 12 {
				return ""
			}

			bits := dobbertinInfo.Bits
			if bits == 0 {
				return ""
			}

			if bits == 32 {
				return fmt.Sprintf("\n\tassert(!(%s ^ K));", register)
			}

			code := "\n"
			for j := 0; j < bits; j += 1 {
				code += fmt.Sprintf("\tassert(%s[%d]);\n", register, j)
			}
			return code
		},
	})
	tmpl, err := tmpl.Parse(layout)
	if err != nil {
		return "", err
	}

	m := []int{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
		0, 4, 8, 12,
		1, 5, 9, 13,
		2, 6, 10, 14,
		3, 7, 11, 15,
		0, 8, 4, 12,
		2, 10, 6, 14,
		1, 9, 5, 13,
		3, 11, 7, 15}
	n := []int{
		3, 7, 11, 19,
		3, 7, 11, 19,
		3, 7, 11, 19,
		3, 7, 11, 19,
		3, 5, 9, 13,
		3, 5, 9, 13,
		3, 5, 9, 13,
		3, 5, 9, 13,
		3, 9, 11, 15,
		3, 9, 11, 15,
		3, 9, 11, 15,
		3, 9, 11, 15}

	registers := []string{
		"a", "d", "c", "b",
		"b", "a", "d", "c",
		"c", "b", "a", "d",
		"d", "c", "b", "a",
	}

	var buffer bytes.Buffer
	steps := instanceInfo.Steps
	tmpl.Execute(&buffer, map[string]interface{}{
		"DobbertinConstant": fmt.Sprintf("0x%x", dobbertinConstant),
		"m":                 m[:steps],
		"n":                 n[:steps],
		"Registers":         registers,
		"OneTargetHash":     instanceInfo.TargetHash == "ffffffffffffffffffffffffffffffff",
	})
	code := buffer.String()

	return code, nil
}

func (encoderSvc *EncoderService) InvokeTransalg(parameters pipeline.Encoding) []pipeline.EncodingPromise {
	err := encoderSvc.filesystemSvc.PrepareTempDir()
	encoderSvc.errorSvc.Fatal(err, "Encoder: failed to create tmp dir")

	encodings := []string{}

	// * Loop through the variations
	encoderSvc.LoopThroughVariation(parameters, func(instanceInfo encoder.InstanceInfo) {
		var dobbertinConstant uint32 = math.MaxUint32
		transalgCode, err := encoderSvc.GenerateTransalgCode(instanceInfo, dobbertinConstant)
		encoderSvc.errorSvc.Fatal(err, "Encoder: failed to generate Transalg code")
		os.WriteFile("./tmp/transalg.alg", []byte(transalgCode), 0644)

		instanceName := encoderSvc.GetInstanceName(instanceInfo)
		encodingFilePath := path.Join(parameters.OutputDir, instanceName)
		encodings = append(encodings, encodingFilePath)

		// Skip if encoding already exists
		if encoderSvc.filesystemSvc.FileExists(encodingFilePath) {
			logrus.Println("Encoder: skipped", encodingFilePath)
			return
		}

		// * Drive the encoder
		command := fmt.Sprintf("%s -i ./tmp/transalg.alg -o %s", encoderSvc.configSvc.Config.Paths.Bin.Transalg, encodingFilePath)
		err = encoderSvc.commandSvc.Create(command).Run()
		encoderSvc.errorSvc.Fatal(err, "Encoder: failed to run Transalg for "+instanceName)

		logrus.Println("Encoder:", encodingFilePath)
	})

	encodingPromises := lo.Map(encodings, func(encoding string, _ int) pipeline.EncodingPromise {
		return EncodingPromise{Encoding: encoding}
	})
	return encodingPromises
}
