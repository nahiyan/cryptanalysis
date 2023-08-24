package services

import (
	"bufio"
	cubeselector "cryptanalysis/internal/cube_selector"
	"cryptanalysis/internal/cuber"
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/pipeline"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

func (cubeSelectorSvc *CubeSelectorService) RandomCubes(cubesCount, selectionSize, offset int, seed int64) []int {
	rnd := rand.New(rand.NewSource(seed))
	indexPermutation := rnd.Perm(cubesCount)[offset:]
	cubes := lo.Map(indexPermutation, func(index, _ int) int {
		return index + 1
	})

	// Return all the cubes if the selecton size is 0
	if selectionSize == 0 {
		return cubes
	}

	randomCubeSelectionCount := int(math.Min(float64(len(indexPermutation)), float64(selectionSize)))

	return cubes[:randomCubeSelectionCount]
}

// TODO: See if it should be in the cubesets or cuber module
func (cubeSelectorSvc *CubeSelectorService) EncodingFromCube(encodingFilePath string, cubesetPath string, cubeIndex int, output io.Writer) error {
	// * 1. Read the instance
	instanceReader, err := os.OpenFile(encodingFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer instanceReader.Close()
	instanceScanner := bufio.NewScanner(instanceReader)

	// * 2. Get the cube from the binary
	cubeLiterals, err := cubeSelectorSvc.cubesetSvc.GetCube(cubesetPath, cubeIndex)
	if err != nil {
		return err
	}

	// * 3. Get the num. of variables and clauses along with the body
	var numVars, numClauses int
	for instanceScanner.Scan() {
		line := instanceScanner.Text()
		if strings.HasPrefix(line, "p cnf") {
			fmt.Sscanf(line, "p cnf %d %d\n", &numVars, &numClauses)

			// * 6. Generate a new header with an increased number of clauses
			newHeader := fmt.Sprintf("p cnf %d %d", numVars, numClauses+len(cubeLiterals))
			output.Write([]byte(newHeader + "\n"))

			continue
		}

		output.Write([]byte(line + "\n"))
	}

	// * 4. Assemble the new encoding by adding the cube literals as unit clauses
	for _, cubeLiteral := range cubeLiterals {
		clause := fmt.Sprintf("%d 0\n", cubeLiteral)
		output.Write([]byte(clause))
	}

	return nil
}

func (cubeSelectorSvc *CubeSelectorService) Select(cubesets []string, parameters pipeline.CubeSelectParams) []encoder.Encoding {
	encodings := []encoder.Encoding{}
	for _, cubeset := range cubesets {
		// Generate the binary cubeset file
		if !cubeSelectorSvc.filesystemSvc.FileExistsNonEmpty(cubeset+".bcubes") || !cubeSelectorSvc.filesystemSvc.FileExistsNonEmpty(cubeset+".bcubes.map") {
			startTime := time.Now()
			err := cubeSelectorSvc.cubesetSvc.BinEncode(cubeset)
			cubeSelectorSvc.errorSvc.Fatal(err, "Cube selector: failed to binary encode "+cubeset)
			log.Printf("Cube selector: generated binary cubeset for %s in %s\n", cubeset, time.Since(startTime))
		}

		var encodingPath string
		{
			segments := strings.Split(cubeset, ".")
			encodingPath = path.Join(cubeSelectorSvc.configSvc.Config.Paths.Encodings, path.Base(strings.Join(segments[:len(segments)-2], ".")))
		}

		var cubeIndices []int
		if parameters.Type == cubeselector.Random {
			cubesCount, err := cubeSelectorSvc.filesystemSvc.CountLines(cubeset)
			cubeSelectorSvc.errorSvc.Fatal(err, "Cube selector: failed to count lines "+cubeset)

			cubeIndices = cubeSelectorSvc.RandomCubes(cubesCount, parameters.Quantity, parameters.Offset, parameters.Seed)
		} else {
			cubeIndices = parameters.Indices
		}

		threshold := 0
		thresholdType := cuber.ThresholdType(cuber.CutoffVars)
		{
			matches := regexp.MustCompile(`march_([nd])(\d+)`).FindAllStringSubmatch(cubeset, 2)

			thresholdType_ := matches[0][1]
			threshold_ := matches[0][2]

			if thresholdType_ == "d" {
				thresholdType = cuber.CutoffDepth
			}
			var err error
			threshold, err = strconv.Atoi(threshold_)
			cubeSelectorSvc.errorSvc.Fatal(err, "Cube selector: failed to get the threshold of "+cubeset)
		}

		for _, cubeIndex := range cubeIndices {
			encoding := encoder.Encoding{
				BasePath: encodingPath,
				Cube: mo.Some(
					encoder.Cube{
						Index:         cubeIndex,
						Threshold:     threshold,
						ThresholdType: thresholdType,
					},
				),
			}
			encodings = append(encodings, encoding)
		}
	}

	if parameters.Type == cubeselector.Random {
		log.Println("Cube selector: randomly selected", len(encodings), "cubes")
	} else {
		log.Println("Cube selector: selected", len(encodings), "cubes")
	}
	return encodings
}

func (cubeSelectorSvc *CubeSelectorService) Run(cubesets []string, parameters pipeline.CubeSelectParams) []encoder.Encoding {
	encodings := []encoder.Encoding{}
	switch parameters.Type {
	case cubeselector.Random:
		encodings = cubeSelectorSvc.Select(cubesets, parameters)
	case cubeselector.Specific:
		encodings = cubeSelectorSvc.Select(cubesets, parameters)
	}

	return encodings
}
