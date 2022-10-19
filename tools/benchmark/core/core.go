package core

// var (
// 	maxTime           = uint(5000)
// 	maxInstancesCount = uint(50)

// 	// Variations
// 	xorOptions = []uint{0}
// 	hashes     = []string{"ffffffffffffffffffffffffffffffff",
// 		"00000000000000000000000000000000"}
// 	adderTypes           = []string{"counter_chain", "dot_matrix"}
// 	stepVariations       = utils.MakeRange(16, 48)
// 	dobbertionVariations = []uint{0, 1}
// 	satSolvers           = []string{constants.CRYPTOMINISAT, constants.KISSAT, constants.CADICAL, constants.GLUCOSE, constants.MAPLESAT}
// )

// func main() {
// 	// Get the steps from the CLI arguments
// 	if len(os.Args) >= 3 {
// 		stepsStart, err1 := strconv.Atoi(os.Args[1])
// 		stepsEnd, err2 := strconv.Atoi(os.Args[2])

// 		if err1 == nil && err2 == nil {
// 			stepVariations = utils.MakeRange(stepsStart, stepsEnd)
// 		}
// 	}

// 	if len(os.Args) >= 4 {
// 		if maxTime_, err := strconv.Atoi(os.Args[3]); err == nil {
// 			maxTime = uint(maxTime_)
// 		}
// 	}

// 	if len(os.Args) >= 5 {
// 		if maxInstancesCount_, err := strconv.Atoi(os.Args[4]); err == nil {
// 			maxInstancesCount = uint(maxInstancesCount_)
// 		}
// 	}

// 	// Count the number of instances for determining the progress
// 	instancesCount := len(xorOptions) * len(hashes) * len(adderTypes) * len(stepVariations) * (len(dobbertionVariations) * len(lo.Filter(stepVariations, func(steps, _ int) bool {
// 		return steps >= 27
// 	})))

// 	// Define the context
// 	context := &types.BenchmarkContext{
// 		Progress: make(map[string][]bool),
// 	}
// 	for _, satSolver := range satSolvers {
// 		context.Progress[satSolver] = make([]bool, instancesCount)
// 	}

// 	// Remove the files from previous execution
// 	os.Remove(constants.BENCHMARK_LOG_FILE_NAME)
// 	os.Remove(constants.VERIFICATION_LOG_FILE_NAME)
// 	for _, satSolver := range satSolvers {
// 		cmd := exec.Command("bash", "-c", fmt.Sprintf("rm %s%s/*.sol", constants.SOLUTIONS_DIR_PATH, satSolver))
// 		if err := cmd.Run(); err != nil {
// 			fmt.Println(cmd.String())
// 			fmt.Println("Failed to delete the solution files: " + err.Error())
// 		}
// 	}

// 	// Solve the encodings for each SAT solver
// 	for _, satSolver := range satSolvers {
// 		var i uint = 0

// 		for _, steps := range stepVariations {
// 			for _, hash := range hashes {
// 				for _, xorOption := range xorOptions {
// 					for _, adderType := range adderTypes {
// 						for _, dobbertin := range dobbertionVariations {
// 							// Skip dobbertin's attacks when steps count < 28
// 							if steps < 28 && dobbertin == 1 {
// 								dobbertin = 0
// 							}

// 							for context.RunningInstances > maxInstancesCount {
// 								time.Sleep(time.Second * 1)
// 							}

// 							filepath := fmt.Sprintf("%smd4_%d_%s_xor%d_%s_dobbertin%d.cnf",
// 								constants.ENCODINGS_DIR_PATH, steps, adderType, xorOption, hash, dobbertin)

// 							startTime := time.Now()
// 							switch satSolver {
// 							case constants.CRYPTOMINISAT:
// 								go cryptoMiniSat(filepath, context, i, startTime)
// 							case constants.KISSAT:
// 								go kissat(filepath, context, i, startTime)
// 							case constants.CADICAL:
// 								go cadical(filepath, context, i, startTime)
// 							case constants.MAPLESAT:
// 								go mapleSat(filepath, context, i, startTime)
// 							case constants.GLUCOSE:
// 								go glucose(filepath, context, i, startTime)
// 							}

// 							context.runningInstances += 1
// 							i++
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}

// 	for !areAllInstancesCompleted(context) {
// 		time.Sleep(time.Second * 1)
// 	}
// }
