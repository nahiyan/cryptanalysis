package services

import (
	"benchmark/internal/pipeline"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"

	"github.com/samber/lo"
)

const (
	Random = "random"
)

type EncodingPromise struct {
	BaseEncodingPath string
	CubeIndex        int
	CubesetPath      string
	SubProblemPath   string
	CubeSelectorSvc  *CubeSelectorService
}

func (encodingPromise EncodingPromise) Get() string {
	cubeSelectorSvc := encodingPromise.CubeSelectorSvc
	subProblemFilePath := encodingPromise.SubProblemPath
	encoding := encodingPromise.BaseEncodingPath
	cubeset := encodingPromise.CubesetPath
	cubeIndex := encodingPromise.CubeIndex

	if cubeSelectorSvc.filesystemSvc.FileExists(subProblemFilePath) {
		fmt.Println("Cube selector: skipped", cubeIndex, subProblemFilePath)
		return subProblemFilePath
	}

	err := cubeSelectorSvc.EncodingFromCube(subProblemFilePath, encoding, cubeset, cubeIndex)
	cubeSelectorSvc.errorSvc.Fatal(err, "Cube selector: failed to generate the encoding from a cube")

	// fmt.Println("Cube selector:", cubeIndex, subProblemFilePath)
	return subProblemFilePath
}

func (cubeSelectorSvc *CubeSelectorService) RandomCubes(cubesCount, selectionSize, offset int, seed int64) []int {
	rand.Seed(seed)
	indexPermutation := rand.Perm(cubesCount)[offset:]
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

func (cubeSelectorSvc *CubeSelectorService) EncodingFromCube(subProblemFilePath, encodingFilePath, cubesetFilePath string, cubeIndex int) error {
	// * 1. Read the instance
	var instance string
	{
		instance_, err := os.ReadFile(encodingFilePath)
		if err != nil {
			return err
		}
		instance = string(instance_)
	}

	// * 2. Open the cubes file
	cubesFile, err := os.Open(cubesetFilePath)
	if err != nil {
		return err
	}
	defer cubesFile.Close()

	// * 3. Get the cube
	cube, _, err := cubeSelectorSvc.filesystemSvc.ReadLine(cubesFile, cubeIndex)
	if err != nil {
		return err
	}

	// * 4. Generate unit clauses from the literals in the cube
	cubeClauses_ := strings.Split(strings.TrimPrefix(cube, "a "), " ")
	cubeClauses := lo.Map(cubeClauses_[:len(cubeClauses_)-1], func(s string, _ int) string {
		return s + " 0"
	})

	// * 5. Get the num. of variables and clauses along with the body
	var numVars, numClauses int
	fmt.Sscanf(instance, "p cnf %d %d\n", &numVars, &numClauses)
	n := len(fmt.Sprintf("p cnf %d %d\n", numVars, numClauses))
	body := instance[n:]

	// * 6. Generate a new header with an increased number of clauses
	newHeader := fmt.Sprintf("p cnf %d %d", numVars, numClauses+len(cubeClauses))

	// * 7. Assemble the new encoding
	newEncoding := fmt.Sprintf("%s\n%s%s\n",
		newHeader,
		body,
		strings.Join(cubeClauses, "\n"))

	// * 8. Write the file
	if err := os.WriteFile(subProblemFilePath, []byte(newEncoding), 0644); err != nil {
		return err
	}

	return nil
}

func (cubeSelectorSvc *CubeSelectorService) GetInfo(cubeset string) (string, error) {
	index := strings.LastIndex(cubeset, ".cnf")
	if index == -1 {
		return "", errors.New("invalid cubeset filename")
	}
	encoding := cubeset[:index+4]
	return encoding, nil
}

func (cubeSelectorSvc *CubeSelectorService) RunRandom(cubesets []string, parameters pipeline.CubeSelecting) []pipeline.PromiseString {
	encodings := []pipeline.PromiseString{}

	// temporary directory for holding the cubes
	if !cubeSelectorSvc.filesystemSvc.FileExists("tmp") {
		err := os.Mkdir("tmp", os.ModePerm)
		cubeSelectorSvc.errorSvc.Fatal(err, "Cube selector: failed to create the ./tmp dir")
	}

	for _, cubeset := range cubesets {
		encoding, err := cubeSelectorSvc.GetInfo(cubeset)
		cubeSelectorSvc.errorSvc.Fatal(err, "Cube selector: failed to get cubeset details "+cubeset)

		cubesCount, err := cubeSelectorSvc.filesystemSvc.CountLines(cubeset)
		cubeSelectorSvc.errorSvc.Fatal(err, "Cube selector: failed to count lines "+cubeset)

		cubeIndices := cubeSelectorSvc.RandomCubes(cubesCount, parameters.Quantity, parameters.Offset, parameters.Seed)

		for _, cubeIndex := range cubeIndices {
			subProblemFilePath := path.Join("./tmp", fmt.Sprintf("%s.cube%d.cnf", cubeset, cubeIndex))

			encodingPromise := EncodingPromise{
				BaseEncodingPath: encoding,
				CubeIndex:        cubeIndex,
				CubesetPath:      cubeset,
				SubProblemPath:   subProblemFilePath,
				CubeSelectorSvc:  cubeSelectorSvc,
			}
			encodings = append(encodings, encodingPromise)
		}
	}

	fmt.Println("Cube selector: randomly selected", len(encodings), "cubes")
	return encodings
}

func (cubeSelectorSvc *CubeSelectorService) Run(cubesets []string, parameters pipeline.CubeSelecting) []pipeline.PromiseString {
	switch parameters.Type {
	case Random:
		encodings := cubeSelectorSvc.RunRandom(cubesets, parameters)
		return encodings
	}

	return []pipeline.PromiseString{}
}
