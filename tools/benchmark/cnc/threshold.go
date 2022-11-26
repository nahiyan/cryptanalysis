package cnc

import (
	"benchmark/constants"
	"benchmark/core"
	"benchmark/encodings"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/alitto/pond"
)

func ProcessMarchLog(log string) (uint, uint) {
	lines := strings.Split(log, "\n")
	var cubeCount, refutedLeaves uint
	for _, line := range lines {
		if strings.HasPrefix(line, "c number of cubes") {
			fmt.Sscanf(line, "c number of cubes %d, including %d refuted leaves", &cubeCount, &refutedLeaves)
		}
	}

	return cubeCount, refutedLeaves
}

func FindThreshold(context types.CommandContext) (uint, time.Duration) {
	// Parse the encoding
	cnf, err := encodings.Process(context.FindCncThreshold.InstanceName)
	if err != nil {
		log.Fatal("Failed to find threshold: ", err)
	}

	// Initialize the variables
	decrementSize := 10
	instanceFilePath := path.Join(constants.EncodingsDirPath, context.FindCncThreshold.InstanceName+".cnf")
	numWorkers := 16
	pool := pond.New(numWorkers, 1000)
	lookaheadSolverMaxTime := time.Second * 5000
	maxCubeCount := 100000
	minRefutedLeaves := 500
	// sampleSize := 100

	type Cubeset struct {
		threshold uint
		cubeCount uint
	}

	// Generate the cubes for varying thresholds
	cubesets := []Cubeset{}
	threshold := int(cnf.FreeVariables) - decrementSize
	for threshold > 0 {
		if threshold%10 != 0 {
			threshold--
			continue
		}

		// Feed the instance to a lookahead solver
		pool.Submit(func(threshold int) func() {
			return func() {
				if utils.FileExists(instanceFilePath) {
					return
				}

				outputFilePath := fmt.Sprintf("%sn%d%s.icnf", constants.EncodingsDirPath, threshold, context.FindCncThreshold.InstanceName)
				output := core.March(instanceFilePath, outputFilePath, uint(threshold), lookaheadSolverMaxTime)

				cubeCount, refutedLeaves := ProcessMarchLog(output)
				fmt.Println(threshold, cubeCount, refutedLeaves)

				// Add to the cubeset if it follows a criteria
				if cubeCount <= uint(maxCubeCount) && refutedLeaves >= uint(minRefutedLeaves) {
					cubesets = append(cubesets, Cubeset{threshold: uint(threshold), cubeCount: cubeCount})
				} else {
					if err := os.Remove(outputFilePath); err != nil {
						fmt.Println("Failed to remove the rejected cubeset: ", err)
					}
				}
			}
		}(threshold))
		threshold -= decrementSize
	}
	pool.StopAndWait()

	// Go through the cubesets, starting from the one with possibly the least difficult cubes (most refuted leaves)
	// for i := len(cubesets) - 1; i >= 0; i-- {
	// 	cubeset := cubesets[i]

	// 	// TODO: Take random sample from the cubeset
	// 	// TODO: Get the runtimes from the solvers
	// }

	return uint(threshold), time.Duration(0)
}
