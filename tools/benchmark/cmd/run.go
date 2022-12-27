package cmd

import (
	"benchmark/internal/injector"
	schemaServices "benchmark/internal/schema/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initRunCmd() *cobra.Command {
	schemaFilePath := ""

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the benchmark based on the defined pipeline",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			schemaSvc := do.MustInvoke[*schemaServices.SchemaService](injector)
			schemaSvc.Process(schemaFilePath)
		},
	}

	runCmd.Flags().StringVarP(&schemaFilePath, "schema", "s", "schema.toml", "A schema is a TOML file that holds the pipelines for the benchmark")

	return runCmd
}
