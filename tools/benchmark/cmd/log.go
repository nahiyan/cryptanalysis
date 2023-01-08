package cmd

import (
	"benchmark/internal/injector"
	logServices "benchmark/internal/log/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initlogCmd() *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Log the tasks performed",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			logSvc := do.MustInvoke[*logServices.LogService](injector)
			logSvc.Run(output)
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "log", "Base path of the output file. E.g. './results/log'")

	return cmd
}
