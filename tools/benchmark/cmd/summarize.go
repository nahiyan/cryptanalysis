package cmd

import (
	services2 "benchmark/internal/combined_logs/services"
	"benchmark/internal/injector"
	services1 "benchmark/internal/summarizer/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initSummarizeCmd() *cobra.Command {
	var workers int
	useCombinedLogs := false

	cmd := &cobra.Command{
		Use:   "summarize",
		Short: "Summarize the benchmark results",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			summarizerSvc := do.MustInvoke[*services1.SummarizerService](injector)
			if useCombinedLogs {
				combinedLogsSvc := do.MustInvoke[*services2.CombinedLogsService](injector)
				combinedLogsSvc.Load()
			}
			summarizerSvc.Run(workers)
		},
	}

	cmd.Flags().IntVarP(&workers, "workers", "w", 100, "Number of workers to read the log files in parallel")
	cmd.Flags().BoolVarP(&useCombinedLogs, "combined-logs", "c", false, "Load all the logs from the all.clog file")

	return cmd
}
