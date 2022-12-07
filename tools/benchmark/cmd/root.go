package cmd

import (
	"benchmark/cnc"
	"benchmark/config"
	"benchmark/constants"
	"benchmark/encodings"
	"benchmark/regular"
	"benchmark/slurm"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var variationsXor_, variationsHashes_, variationsAdders_, variationsSatSolvers_, variationsDobbertin_, variationsDobbertinBits_, variationsSteps_, simplifier, simplificationInstanceName, reconstructInstanceName, reconstructReconstructionStackPath string
var instanceMaxTime, maxConcurrentInstancesCount, digest, generateEncodings, sessionId, cubeCutoffVars, cubeSelectionCount, cubeIndex, simplifierPasses, simplifierPassDuration uint
var cleanResults, isCubeEnabled, simplifierReconstruct bool
var seed int64
var genSubProblem types.GenSubProblem
var findCncThreshold types.FindCncThreshold

var rootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "benchmark - MD4 Inversion tool",
	Long:  `MD4 Inversion tool`,
}

var slurmCmd = &cobra.Command{
	Use:   "slurm",
	Short: "Run the benchmark instances as small jobs managed by Slurm",
	Run: func(cmd *cobra.Command, args []string) {
		context := processFlags()
		slurm.Run(&context)
	},
}

var regularCmd = &cobra.Command{
	Use:   "regular",
	Short: "Run the benchmark instances as concurrent operations",
	Run: func(cmd *cobra.Command, args []string) {
		context := processFlags()
		regular.Run(&context)
	},
}

var simplifyCmd = &cobra.Command{
	Use:   "simplify",
	Short: "Simplify the encodings through preprocessing or in-processing",
	Run: func(cmd *cobra.Command, args []string) {
		context := processFlags()

		if context.Simplification.Simplifier == constants.ArgCadical {
			encodings.CadicalSimplify(
				fmt.Sprintf("%s%s.cnf", constants.EncodingsDirPath, context.Simplification.InstanceName),
				context.Simplification.Passes,
				time.Duration(context.Simplification.PassDuration)*time.Second,
				simplifierReconstruct)
		}
	},
}

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct",
	Short: "Reconstruct a solution from a reconstruction stack",
	Run: func(cmd *cobra.Command, args []string) {
		context := processFlags()

		encodings.ReconstructSolution(context.Reconstruction.InstanceName, context.Reconstruction.StackFilePath, []types.Range{{Start: 1, End: 512}, {Start: 641, End: 768}})
	},
}

var findCncThresholdCmd = &cobra.Command{
	Use:   "find-cnc-threshold",
	Short: "Find cutoff threshold for Cube & Conquer",
	Run: func(cmd *cobra.Command, args []string) {
		context := processFlags()

		bestThreshold, bestEstimate := cnc.FindThreshold(context)
		fmt.Printf("Best threshold, n = %d, with an estimated time of %s\n", bestThreshold, bestEstimate.String())
	},
}

var genSubProblemCmd = &cobra.Command{
	Use:   "gen-subproblem",
	Short: "Generate a subproblem from a specific cube",
	Run: func(cmd *cobra.Command, args []string) {
		var threshold *uint
		if genSubProblem.Threshold != 0 {
			threshold = lo.ToPtr(genSubProblem.Threshold)
		}

		subproblem, err := encodings.GenerateSubProblemAsStringWithThreshold(genSubProblem.InstanceName, int(genSubProblem.CubeIndex), threshold)
		if err != nil {
			panic(err)
		}

		fmt.Println(subproblem)
	},
}

func processFlags() types.CommandContext {
	context := types.CommandContext{}

	// XOR Variations
	{
		var variationsXor []uint
		pieces := strings.Split(variationsXor_, ",")
		if len(pieces) == 1 || len(pieces) == 2 {
			variationsXor = lo.Map(pieces, func(s string, i int) uint {
				value, _ := strconv.Atoi(s)
				return uint(value)
			})
			context.VariationsXor = variationsXor
		}
	}

	// Adder Variations
	{
		var variationsAdders []string
		pieces := strings.Split(variationsAdders_, ",")
		if len(pieces) == 1 || len(pieces) == 2 {
			variationsAdders = lo.Filter(pieces, func(s string, i int) bool {
				return s == constants.ArgCounterChain || s == constants.ArgDotMatrix || s == constants.ArgEspresso || s == constants.ArgTwoOperand
			})
			context.VariationsAdders = variationsAdders
		}
	}

	// Dobbertin Variations
	{
		var variationsDobbertin []uint
		pieces := strings.Split(variationsDobbertin_, ",")
		if len(pieces) == 1 || len(pieces) == 2 {
			variationsDobbertin = lo.Map(pieces, func(s string, i int) uint {
				value, _ := strconv.Atoi(s)
				return uint(value)
			})

			context.VariationsDobbertin = variationsDobbertin
		}
	}

	// SAT Solver Variations
	{
		satSolvers := []string{constants.ArgCryptoMiniSat, constants.ArgKissat, constants.ArgCadical, constants.ArgGlucoseSyrup, constants.ArgMapleSat, constants.ArgXnfSat}

		var variationSatSolvers []string
		pieces := strings.Split(variationsSatSolvers_, ",")
		if len(pieces) > 0 {
			variationSatSolvers = lo.Filter(pieces, func(s string, i int) bool {
				return lo.Contains(satSolvers, s)
			})
			context.VariationsSatSolvers = variationSatSolvers
		}
	}

	// Hash variations
	{
		var variationsHashes []string
		pieces := strings.Split(variationsHashes_, ",")
		if len(pieces) > 0 {
			variationsHashes = lo.Filter(pieces, func(s string, i int) bool {
				return len(s) == 32
			})
			context.VariationsHashes = variationsHashes
		}
	}

	// Step Variations
	{
		var variationsSteps []uint
		pieces := strings.Split(variationsSteps_, ",")
		if len(pieces) > 0 {
			for _, piece := range pieces {
				isRange := len(strings.Split(piece, "-")) == 2
				if isRange {
					tuple := make([]int, 2)
					rangePieces := strings.Split(piece, "-")
					{
						minValue, err := strconv.Atoi(rangePieces[0])
						if err != nil {
							continue
						}
						tuple[0] = (minValue)
					}
					{
						maxValue, err := strconv.Atoi(rangePieces[1])
						if err != nil {
							continue
						}
						tuple[1] = (maxValue)
					}

					values := utils.MakeRange(tuple[0], tuple[1])

					for _, value := range values {
						variationsSteps = append(variationsSteps, uint(value))
					}
				} else {
					value, err := strconv.Atoi(piece)
					if err != nil {
						continue
					}
					variationsSteps = append(variationsSteps, uint(value))
				}
			}

			context.VariationsSteps = variationsSteps
		}
	}

	// Dobbertin Bits Variations
	{
		var variationsDobbertinBits []uint
		pieces := strings.Split(variationsDobbertinBits_, ",")
		if len(pieces) > 0 {
			for _, piece := range pieces {
				isRange := len(strings.Split(piece, "-")) == 2
				if isRange {
					tuple := make([]int, 2)
					rangePieces := strings.Split(piece, "-")
					{
						minValue, err := strconv.Atoi(rangePieces[0])
						if err != nil {
							continue
						}
						tuple[0] = (minValue)
					}
					{
						maxValue, err := strconv.Atoi(rangePieces[1])
						if err != nil {
							continue
						}
						tuple[1] = (maxValue)
					}

					values := utils.MakeRange(tuple[0], tuple[1])

					for _, value := range values {
						variationsDobbertinBits = append(variationsDobbertinBits, uint(value))
					}
				} else {
					value, err := strconv.Atoi(piece)
					if err != nil {
						continue
					}
					variationsDobbertinBits = append(variationsDobbertinBits, uint(value))
				}
			}

			context.VariationsDobbertinBits = lo.Reverse(variationsDobbertinBits)
		}
	}

	// Max. Instance Time
	context.InstanceMaxTime = instanceMaxTime

	// Max. concurrent instances count
	context.MaxConcurrentInstancesCount = maxConcurrentInstancesCount

	// Reset data
	context.CleanResults = cleanResults

	// TODO: Prevent executing commands here
	// Remove leftover results
	if context.CleanResults {
		exec.Command("bash", "-c", fmt.Sprintf("rm %s*.csv", constants.LogsDirPath)).Run()
		exec.Command("bash", "-c", fmt.Sprintf("rm %s*.cnf", constants.EncodingsDirPath)).Run()
		exec.Command("bash", "-c", fmt.Sprintf("rm %s*.icnf", constants.EncodingsDirPath)).Run()
		exec.Command("bash", "-c", fmt.Sprintf("rm %s*.sol", constants.SolutionsDirPath)).Run()
	}

	// Cubing
	if isCubeEnabled {
		context.CubeParams = new(types.CubeParams)
		context.CubeParams.CutoffVars = cubeCutoffVars
		context.CubeParams.SelectionSize = cubeSelectionCount
		context.CubeParams.CubeIndex = cubeIndex
	}

	// Seed
	context.Seed = seed

	// Digest
	context.Digest = digest

	// Generate encodings
	context.GenerateEncodings = generateEncodings

	// Session ID
	context.SessionId = sessionId

	// Simplification
	context.Simplification.Simplifier = simplifier
	context.Simplification.Passes = simplifierPasses
	context.Simplification.InstanceName = simplificationInstanceName
	context.Simplification.PassDuration = simplifierPassDuration
	context.Simplification.Reconstruct = simplifierReconstruct

	// Reconstruction
	context.Reconstruction.InstanceName = reconstructInstanceName
	context.Reconstruction.StackFilePath = reconstructReconstructionStackPath

	context.FindCncThreshold = findCncThreshold

	return context
}

func init() {
	// TODO: Feed directly to the command context
	// Flags and arguments
	rootCmd.PersistentFlags().StringVar(&variationsXor_, "var-xor", "0", "Comma-separated variations of XOR. Possible values: 0, 1")
	rootCmd.PersistentFlags().StringVar(&variationsAdders_, "var-adders", "cc,dm", "Comma-separated variations of the adders. Possible values: cm, dm, esp, 2op")
	rootCmd.PersistentFlags().StringVar(&variationsSatSolvers_, "var-sat-solvers", "cms,ks,cdc,gs,ms", "Comma-separated variations of the SAT solvers. Possible values: cms, ks, cdc, gs, ms, xnf")
	rootCmd.PersistentFlags().StringVar(&variationsHashes_, "var-hashes", "ffffffffffffffffffffffffffffffff,00000000000000000000000000000000", "Comma-separated variations of the hashes. Possible values: ffffffffffffffffffffffffffffffff, 00000000000000000000000000000000")
	rootCmd.PersistentFlags().StringVar(&variationsDobbertin_, "var-dobbertin", "0,1", "Comma-separated variations of the Dobbertin's attack. Possible values: 0, 1")
	rootCmd.PersistentFlags().StringVar(&variationsDobbertinBits_, "var-dobbertin-bits", "32", "Comma-separated variations of the values and/or ranges of the number of significant bits to constrain in Dobbertin's attack (The order of the values evaluated is reversed)")
	rootCmd.PersistentFlags().StringVar(&variationsSteps_, "var-steps", "31-39", "Comma-separated variations of the values and/or ranges of steps")
	rootCmd.PersistentFlags().UintVar(&instanceMaxTime, "max-time", 5000, "Maximum time in seconds for each instance to run")
	rootCmd.PersistentFlags().BoolVar(&cleanResults, "clean-results", false, "Remove leftover results from previous sessions")

	rootCmd.PersistentFlags().BoolVar(&isCubeEnabled, "cube", false, "Produce cubes from the instances and solve them")
	rootCmd.PersistentFlags().UintVar(&cubeCutoffVars, "cube-cutoff-vars", 3000, "Number of variables as a threshold for cube generation")
	rootCmd.PersistentFlags().UintVar(&cubeSelectionCount, "cube-selection-count", 0, "Number of cubes to select randomly for solving")
	rootCmd.PersistentFlags().UintVar(&cubeIndex, "cube-index", 0, "Index of a specific cube to solve")

	rootCmd.PersistentFlags().Int64Var(&seed, "seed", 1, "Seed for the randomization")
	rootCmd.PersistentFlags().UintVar(&generateEncodings, "generate-encodings", 1, "Flag whether to generate encodings prior to solving")
	rootCmd.PersistentFlags().UintVar(&sessionId, "session-id", 0, "ID of a pre-existing session")

	regularCmd.Flags().UintVar(&maxConcurrentInstancesCount, "max-instances", 50, "Maximum number of instances to run concurrently")
	slurmCmd.Flags().UintVar(&digest, "digest", 0, "The ID of the finished slurm job that needs to be digested")

	simplifyCmd.Flags().StringVar(&simplifier, "simplifier", "cdc", "Name of the simplifier. Possible values: cdc")
	simplifyCmd.Flags().UintVar(&simplifierPasses, "passes", 0, "Number of passes (100s) for simplification; 0 for auto.")
	simplifyCmd.Flags().UintVar(&simplifierPassDuration, "pass-duration", 100, "Duration of simplifier passes in seconds")
	simplifyCmd.Flags().StringVar(&simplificationInstanceName, "instance-name", "", "Name of the instance to simplify")
	simplifyCmd.Flags().BoolVar(&simplifierReconstruct, "reconstruct", false, "Reconstruct the CNF after every simplification pass")

	reconstructCmd.Flags().StringVar(&reconstructInstanceName, "instance-name", "", "Instance name of the solution that needs reconstruction")
	reconstructCmd.Flags().StringVarP(&reconstructReconstructionStackPath, "reconstruction-stack-path", "r", "reconstruction-stack.txt", "Path to the reconstruction stack")

	findCncThresholdCmd.Flags().StringVar(&findCncThreshold.InstanceName, "instance-name", "", "Name of the instance to find the CnC threshold for")
	findCncThresholdCmd.Flags().UintVar(&findCncThreshold.NumWorkersLookahead, "num-workers-lookahead", 16, "Number of workers in the Lookahead pool")
	findCncThresholdCmd.Flags().UintVar(&findCncThreshold.NumWorkersCdcl, "num-workers-cdcl", 16, "Number of workers in the CDCL pool")
	findCncThresholdCmd.Flags().UintVar(&findCncThreshold.SampleSize, "sample-size", 1000, "Size of the random sample from each cubeset for estimating the runtime on the full cubeset")
	findCncThresholdCmd.Flags().UintVar(&findCncThreshold.MaxCubes, "max-cubes", 1000000, "Max number of cubes to consider for finding the estimate of a cubeset")
	findCncThresholdCmd.Flags().UintVar(&findCncThreshold.MinRefutedLeaves, "min-refuted-leaves", 500, "Min number of refuted to consider for finding the estimate of a cubeset")
	findCncThresholdCmd.Flags().UintVar(&findCncThreshold.CdclTimeout, "cdcl-timeout", 5000, "Max. number of seconds given to a CDCL instance")

	genSubProblemCmd.Flags().StringVar(&genSubProblem.InstanceName, "instance-name", "", "Name of the instance to generate the sub-problem for")
	genSubProblemCmd.Flags().UintVar(&genSubProblem.CubeIndex, "cube-index", 0, "Index of the cube to generate the sub-problem for")
	genSubProblemCmd.Flags().UintVar(&genSubProblem.Threshold, "threshold", 0, "Threshold of the cubeset using which to generate the sub-problem")

	// Commands
	rootCmd.AddCommand(regularCmd)
	rootCmd.AddCommand(slurmCmd)
	rootCmd.AddCommand(simplifyCmd)
	rootCmd.AddCommand(reconstructCmd)
	rootCmd.AddCommand(findCncThresholdCmd)
	rootCmd.AddCommand(genSubProblemCmd)
}

func Execute() {
	// Process the configuration from config.toml
	config.ProcessConfig()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
