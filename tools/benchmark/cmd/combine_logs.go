package cmd

import (
	"benchmark/internal/combined_logs/services"
	"benchmark/internal/injector"

	"github.com/spf13/cobra"
)

func initCombineLogsCmd() *cobra.Command {
	var workers int

	cmd := &cobra.Command{
		Use:   "combine-logs",
		Short: "Combine log files into one .clog file",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			combinedLogsSvc, err := services.NewCombinedLogsService(injector)
			if err != nil {
				panic(err)
			}
			combinedLogsSvc.Generate(workers)
		},
	}

	cmd.Flags().IntVarP(&workers, "workers", "w", 100, "Number of workers to read the log files in parallel")

	return cmd
}
