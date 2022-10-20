package cmd

import (
	"benchmark/constants"
	"benchmark/regular"
	"benchmark/slurm"
	"benchmark/types"
	"benchmark/utils"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var variationsXor_ string
var variationsHashes_ string
var variationsAdders_ string
var variationsSatSolvers_ string
var variationsDobbertin_ string
var variationsSteps_ string
var instanceMaxTime uint
var maxConcurrentInstancesCount uint
var cleanResults bool

var rootCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "benchmark - MD4 Inversion tool",
	Long:  `MD4 Inversion tool`,
	Run: func(cmd *cobra.Command, args []string) {

	},
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

var aggregateLogsCmd = &cobra.Command{
	Use:   "aggregate_logs",
	Short: "Aggregate the logs into a single file for each category",
	Run: func(cmd *cobra.Command, args []string) {
		utils.AggregateLogs()
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
				return s == constants.ArgCounterChain || s == constants.ArgDotMatrix
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

	// Max. Instance Time
	context.InstanceMaxTime = instanceMaxTime

	// Max. concurrent instances count
	context.MaxConcurrentInstancesCount = maxConcurrentInstancesCount

	// Reset data
	context.CleanResults = cleanResults

	// TODO: Improve the way this is handled
	// Remove leftover results
	if context.CleanResults {
		exec.Command("bash", "-c", "rm "+constants.ResultsDirPat+"*.log").Run()
	}

	return context
}

func init() {
	// Flags
	rootCmd.PersistentFlags().StringVar(&variationsXor_, "var-xor", "0,1", "Comma-separated variations of XOR. Possible values: 0, 1")
	rootCmd.PersistentFlags().StringVar(&variationsAdders_, "var-adders", "cc,dm", "Comma-separated variations of the adders. Possible values: cm, dm")
	rootCmd.PersistentFlags().StringVar(&variationsSatSolvers_, "var-sat-solvers", "cms,ks,cdc,gs,ms", "Comma-separated variations of the SAT solvers. Possible values: cms, ks, cdc, gs, ms, xnf")
	rootCmd.PersistentFlags().StringVar(&variationsHashes_, "var-hashes", "ffffffffffffffffffffffffffffffff,00000000000000000000000000000000", "Comma-separated variations of the hashes. Possible values: ffffffffffffffffffffffffffffffff, 00000000000000000000000000000000")
	rootCmd.PersistentFlags().StringVar(&variationsDobbertin_, "var-dobbertin", "0,1", "Comma-separated variations of the Dobbertin's attack. Possible values: 0, 1")
	rootCmd.PersistentFlags().StringVar(&variationsSteps_, "var-steps", "31-39", "Comma-separated variations of the values and/or ranges of steps")
	rootCmd.PersistentFlags().UintVar(&instanceMaxTime, "max-time", 5000, "Maximum time in seconds for each instance to run")
	regularCmd.Flags().UintVar(&maxConcurrentInstancesCount, "max-instances", 50, "Maximum number of instances to run concurrently")
	rootCmd.PersistentFlags().BoolVar(&cleanResults, "clean-results", false, "Remove leftover results from previous sessions")

	// Commands
	rootCmd.AddCommand(regularCmd)
	rootCmd.AddCommand(slurmCmd)
	rootCmd.AddCommand(aggregateLogsCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
