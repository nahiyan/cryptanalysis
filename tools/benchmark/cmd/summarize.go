package cmd

import (
	"benchmark/internal/injector"
	summarizerServices "benchmark/internal/summarizer/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initSummarizeCmd() *cobra.Command {
	var workers int

	cmd := &cobra.Command{
		Use:   "summarize",
		Short: "Summarize the benchmark results",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			summarizerSvc := do.MustInvoke[*summarizerServices.SummarizerService](injector)
			summarizerSvc.Run(workers)
		},
	}

	cmd.Flags().IntVarP(&workers, "workers", "w", 100, "Number of workers to read the log files in parallel")

	return cmd
}
