package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cryptanalysis",
	Short: "cryptanalysis",
	Long:  `Cryptanalysis tool for MD4, MD5, SHA-256, etc.`,
}

func init() {
	// Commands
	rootCmd.AddCommand(initRunCmd())
	rootCmd.AddCommand(initTaskCmd())
	rootCmd.AddCommand(initSummarizeCmd())
	rootCmd.AddCommand(initCombineLogsCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
