package cmd

import (
	"benchmark/internal/injector"
	schemaServices "benchmark/internal/schema/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initRunCmd() *cobra.Command {
	schemaFilePath := ""

	cmd := &cobra.Command{
		Use:   "run [flags] [schema_file_path]",
		Short: "Run the benchmark based on the defined pipeline",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			schemaFilePath = args[0]
			injector := injector.New()
			schemaSvc := do.MustInvoke[*schemaServices.SchemaService](injector)
			schemaSvc.Process(schemaFilePath)
		},
	}

	// cmd.Flags().StringVarP(&schemaFilePath, "schema", "s", "schema.toml", "A schema is a TOML file that holds the pipelines for the benchmark")

	return cmd
}
