package cmd

import (
	"benchmark/internal/injector"
	summarizerServices "benchmark/internal/summarizer/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initSummarizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summarize",
		Short: "Summarize the benchmark results",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			summarizerSvc := do.MustInvoke[*summarizerServices.SummarizerService](injector)
			summarizerSvc.Run()
		},
	}

	return cmd
}
