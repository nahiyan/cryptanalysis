package cmd

import (
	"benchmark/internal/injector"
	summarizerServices "benchmark/internal/summarizer/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initSummarizeCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "summarize",
		Short: "Summarize the benchmark results",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			summarizerSvc := do.MustInvoke[*summarizerServices.SummarizerService](injector)
			summarizerSvc.Run(output)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "summary", "Base path of the output file. E.g. './results/summary'")

	return cmd
}
