package services

import (
	"benchmark/internal/command"
	"benchmark/internal/cuber"
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alitto/pond"
)

func (cuberSvc *CuberService) EncodingsFromCubes(encodingPath, cubesetPath string, depth int) []string {
	var encodings []string
	for i := 1; i <= 2; i++ {
		newBaseEncoding := path.Base(cubesetPath)
		newEncodingsPath := path.Join(cuberSvc.configSvc.Config.Paths.Encodings, newBaseEncoding+fmt.Sprintf(".cube%d.cnf", i-1))

		encodingWriter, err := os.OpenFile(newEncodingsPath, os.O_CREATE|os.O_WRONLY, 0644)
		cuberSvc.errorSvc.Fatal(err, "Inc. cuber: failed to write encoding from cube")
		err = cuberSvc.cubeSelectorSvc.EncodingFromCube(encodingPath, cubesetPath, i, encodingWriter)
		cuberSvc.errorSvc.Fatal(err, "Inc. cuber: failed to write encoding from cube")
		encodings = append(encodings, newEncodingsPath)
	}

	return encodings
}

func (cuberSvc *CuberService) Depth1Cube(encodingPath string, cubeParams pipeline.CubeParams, suffix string) string {
	cubesetPaths := make([]string, 0)
	cmdGroup := command.Group{}

	// Cube up to a depth of 1
	err := cuberSvc.TrackedInvoke(InvokeParameters{
		Encoding:         encodingPath,
		ThresholdType:    cuber.CutoffDepth,
		Threshold:        1,
		Timeout:          time.Duration(cubeParams.Timeout) * time.Second,
		MinCubes:         0,
		MaxCubes:         math.MaxInt,
		MinRefutedLeaves: 0,
		Suffix:           suffix,
		MaxVariable:      512,
		SkipLogs:         true,
	}, InvokeControl{
		CommandGroup: &cmdGroup,
		ShouldStop:   map[string]bool{},
		CubesetPaths: &cubesetPaths,
	})
	cuberSvc.errorSvc.Fatal(err, "Inc. cuber: failed to cube "+encodingPath)

	return cubesetPaths[0]
}

func adjoinCubes(cubesTree [][]int) string {
	bt := BinaryTree{}
	bt.Insert(0)
	for i := range cubesTree {
		for j := range cubesTree[i] {
			bt.Insert(cubesTree[i][j])
		}
	}

	branches := [][]int{}
	GetBranches(bt.root, []int{}, &branches)

	cubeset := strings.Builder{}
	for _, branch := range branches {
		branch = branch[1:]
		branch = append(branch, 0)
		assumption := "a "
		for _, item := range branch {
			assumption += fmt.Sprintf("%d ", item)
		}
		assumption = assumption[:len(assumption)-1]
		assumption += "\n"
		cubeset.WriteString(assumption)
	}

	return cubeset.String()
}

func addToCubesTree(cubesTree [][]int, cubeset string, depth, index int) error {
	file, err := os.Open(cubeset)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		literal_ := line[2 : len(line)-2]
		literal, err := strconv.Atoi(literal_)
		if err != nil {
			return err
		}

		cubesTree[depth][index*2+i] = literal
		i++
	}

	return nil
}

// TODO: Make the process resumable
func (cuberSvc *CuberService) RunIncremental(encodings []encoder.Encoding, cubeParams pipeline.CubeParams, simplifyParams pipeline.SimplifyParams, solveParams pipeline.SolveParams) []string {
	encodings_ := cuberSvc.simplifierSvc.Run(encodings, simplifyParams)
	cubesets := []string{}
	cuberSvc.filesystemSvc.PrepareDir(cuberSvc.configSvc.Config.Paths.Cubesets)
	startTime := time.Now()
	// debugFile, _ := os.OpenFile("/tmp/debug.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	for _, encoding := range encodings_ {
		// Threshold is assumed to be cutoff depth for incremental cubing
		for _, threshold := range cubeParams.Thresholds {
			encodingPath := encoding.GetName() + ".cnf"
			{
				newBase := path.Join(cuberSvc.configSvc.Config.Paths.Cubesets, path.Base(encodingPath))
				cubesetPath := fmt.Sprintf("%s.march_d%d.cubes", newBase, threshold)
				if cuberSvc.filesystemSvc.FileExists(cubesetPath) {
					cubesets = append(cubesets, cubesetPath)
					continue
				}
			}

			endDepth := threshold - 1
			cubesTree := [][]int{}
			// tree := new(node)

			// sample := []encoder.Encoding{}

			// There should be 1 cubeset
			cubesetPath := cuberSvc.Depth1Cube(encodingPath, cubeParams, "inc_d0i0")
			err := cuberSvc.cubesetSvc.BinEncode(cubesetPath)
			cuberSvc.errorSvc.Fatal(err, "Inc. cuber: failed to binary encode the cubesets file")
			previousEncodings := cuberSvc.EncodingsFromCubes(encodingPath, cubesetPath, 0)
			previousCubesets := []string{cubesetPath}

			// Insert into the tree
			cubesTree = append(cubesTree, make([]int, 2))
			addToCubesTree(cubesTree, cubesetPath, 0, 0)

			// Each iteration will result in a 1 level deeper cubing
			for depth := 1; depth < endDepth+1; depth++ {
				pool := pond.New(cubeParams.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))
				lock := sync.Mutex{}

				// encodingsCount := int(math.Pow(2, float64(depth+1)))
				newEncodings := []string{}
				newCubesets := []string{}
				cubesTree = append(cubesTree, make([]int, int(math.Pow(2, float64(depth+1)))))

				// Cube each new encoding by 1 depth
				for _, previousEncoding := range previousEncodings {
					pool.Submit(func(previousEncoding string) func() {
						return func() {
							matches := regexp.MustCompile(`march_inc_d\d+i(\d+).cubes.cube(\d+).cnf`).FindAllStringSubmatch(previousEncoding, 1)
							cubesetIndex, _ := strconv.Atoi(matches[0][1])
							cubeIndex, _ := strconv.Atoi(matches[0][2])
							index := cubesetIndex*2 + cubeIndex
							suffix := fmt.Sprintf("inc_d%di%d", depth, index)

							cubesetPath := cuberSvc.Depth1Cube(previousEncoding, cubeParams, suffix)
							err := cuberSvc.cubesetSvc.BinEncode(cubesetPath)
							cuberSvc.errorSvc.Fatal(err, "Inc. cuber: failed to binary encode the cubesets file")

							// Insert into the tree
							addToCubesTree(cubesTree, cubesetPath, depth, index)

							newEncodings_ := cuberSvc.EncodingsFromCubes(previousEncoding, cubesetPath, depth)

							lock.Lock()
							newEncodings = append(newEncodings, newEncodings_...)
							newCubesets = append(newCubesets, cubesetPath)
							lock.Unlock()
						}
					}(previousEncoding))
				}
				pool.StopAndWait()

				// Remove the previous encodings and cubesets
				for _, encoding := range previousEncodings {
					os.Remove(encoding)
				}
				for _, cubeset := range previousCubesets {
					os.Remove(cubeset)
					os.Remove(cubeset + ".bcubes")
					os.Remove(cubeset + ".bcubes.map")
				}

				previousEncodings = newEncodings
				previousCubesets = newCubesets

				cubes := adjoinCubes(cubesTree)
				newCubesFilePath := path.Join(cuberSvc.configSvc.Config.Paths.Cubesets, path.Base(fmt.Sprintf("%s.march_d%d.cubes", encodingPath, depth+1)))
				err := os.WriteFile(newCubesFilePath, []byte(cubes), 0644)
				cuberSvc.errorSvc.Fatal(err, "Inc. cuber: failed to write the cubes file")
				cubesets = append(cubesets, newCubesFilePath)
			}

			// Remove the previous encodings and cubesets
			for _, encoding := range previousEncodings {
				os.Remove(encoding)
			}
			for _, cubeset := range previousCubesets {
				os.Remove(cubeset)
				os.Remove(cubeset + ".bcubes")
				os.Remove(cubeset + ".bcubes.map")
			}
		}
	}

	log.Println("Process took: ", time.Since(startTime))

	return cubesets[len(cubesets)-1:]
}
