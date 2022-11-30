package encodings

import (
	"benchmark/config"
	"benchmark/types"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func CadicalSimplify(filePath string, passes uint, duration time.Duration, reconstrut bool) {
	var pass uint = 1
	inputFilePath := filePath
	for {
		fmt.Printf("Simplifier pass %d\n", pass)

		simpInstancePath := filePath[:len(filePath)-4] + fmt.Sprintf("_cadical_simp%d.cnf", pass)
		simpReconstStackFilePath := simpInstancePath[:len(simpInstancePath)-4] + "_reconst_stack.txt"

		command := fmt.Sprintf("%s %s -o %s -e %s -t %.0f",
			config.Get().Paths.Bin.Cadical,
			inputFilePath,
			simpInstancePath,
			simpReconstStackFilePath,
			duration.Seconds())
		output, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			fmt.Println(command)
			panic(fmt.Sprintf("Failed to simplify at pass %d: %s", pass, err.Error()))
		}

		// Optional: Reconstruct any removed clauses containing the message or target hash variables
		if reconstrut {
			if err := ReconstructEncoding(simpInstancePath, simpReconstStackFilePath, []types.Range{{Start: 1, End: 512}, {Start: 641, End: 768}}); err != nil {
				panic("Failed to reconstruct the simplified file: " + err.Error())
			}
		}

		{
			eliminatedVars := 0
			if index := strings.Index(string(output), "c eliminated:"); index != -1 {
				fmt.Sscanf(string(output)[index:], "c eliminated: %d", &eliminatedVars)
			}
			fmt.Printf("Eliminated %d variables\n", eliminatedVars)
		}

		// Break if this is the last pass
		if passes != 0 && pass == passes {
			break
		}

		// Break if the passes aren't eliminating any variable and the number of passes is set to auto
		if passes == 0 {
			eliminatedVars := 0
			if index := strings.Index(string(output), "c eliminated:"); index != -1 {
				fmt.Sscanf(string(output)[index:], "c eliminated: %d", &eliminatedVars)
			}

			if eliminatedVars == 0 {
				break
			}
		}

		inputFilePath = simpInstancePath

		pass++
	}
}
