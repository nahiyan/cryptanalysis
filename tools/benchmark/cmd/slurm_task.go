package cmd

import (
	errorModule "benchmark/internal/error"
	"benchmark/internal/injector"
	slurmServices "benchmark/internal/slurm/services"
	"benchmark/internal/solver"
	solverServices "benchmark/internal/solver/services"
	"log"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initSlurmTaskCmd() *cobra.Command {
	var taskId int

	cmd := &cobra.Command{
		Use:   "slurm-task",
		Short: "Run Slurm task",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()
			slurmSvc := do.MustInvoke[*slurmServices.SlurmService](injector)
			solverSvc := do.MustInvoke[*solverServices.SolverService](injector)

			task, err := slurmSvc.GetTask(taskId)
			if err != nil && err == errorModule.ErrKeyNotFound {
				log.Fatal("Task ID not found")
			}

			solverSvc.Settings.Timeout = task.Timeout
			solverSvc.TrackedInvoke(task.Encoding, solver.Solver(task.Solver))
		},
	}

	cmd.Flags().IntVarP(&taskId, "job-id", "j", 1, "ID of the job/task")

	return cmd
}
