package cmd

import (
	databaseServices "benchmark/internal/database/services"
	"benchmark/internal/injector"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initClearCmd() *cobra.Command {
	var bucket string

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear all the data in the specified bucket",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			databaseSvc := do.MustInvoke[*databaseServices.DatabaseService](injector)
			databaseSvc.RemoveAll(bucket)
		},
	}

	cmd.Flags().StringVarP(&bucket, "bucket", "b", "solutions", "The name of the bucket you want to remove. E.g. solutions")

	return cmd
}
