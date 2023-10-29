package cmd

import (
	services2 "cryptanalysis/internal/combined_logs/services"
	"cryptanalysis/internal/injector"
	services1 "cryptanalysis/internal/schema/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initRunCmd() *cobra.Command {
	useCombinedLogs := false

	cmd := &cobra.Command{
		Use:   "run [flags] [schema_file_path]",
		Short: "Run the cryptanalysis based on the defined pipeline",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			schemaFilePath := args[0]
			injector := injector.New()
			schemaSvc := do.MustInvoke[*services1.SchemaService](injector)
			if useCombinedLogs {
				combinedLogsSvc := do.MustInvoke[*services2.CombinedLogsService](injector)
				combinedLogsSvc.Load()
			}
			schemaSvc.Process(schemaFilePath)
		},
	}

	cmd.Flags().BoolVarP(&useCombinedLogs, "combined-logs", "c", false, "Load all the logs from the all.clog file")

	return cmd
}
