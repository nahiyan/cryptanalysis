package cmd

import (
	"benchmark/internal/injector"
	"benchmark/internal/solver"
	solverSvc "benchmark/internal/solver/services"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initSlurmTaskCmd() *cobra.Command {
	var encoding, solver_ string
	var taskId, timeout int

	cmd := &cobra.Command{
		Use:   "slurm-task",
		Short: "Solve slurm task",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			solverSvc := do.MustInvoke[*solverSvc.SolverService](injector)
			solverSvc.Invoke(encoding, solver.Solver(solver_))
		},
	}

	cmd.Flags().IntVarP(&taskId, "job-id", "j", 1, "ID of the job/task")
	cmd.Flags().IntVarP(&timeout, "timeout", "t", 5000, "Timeout in seconds")

	return cmd
}
