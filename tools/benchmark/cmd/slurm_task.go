package cmd

import (
	services2 "benchmark/internal/cube_slurm_task/services"
	cuberServices "benchmark/internal/cuber/services"
	errorModule "benchmark/internal/error"
	"benchmark/internal/injector"
	services1 "benchmark/internal/solve_slurm_task/services"
	"benchmark/internal/solver"
	solverServices "benchmark/internal/solver/services"
	"fmt"
	"log"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initSlurmTaskCmd() *cobra.Command {
	const (
		Solve = "solve"
		Cube  = "cube"
	)

	var id int
	var type_ string

	cmd := &cobra.Command{
		Use:   "slurm-task",
		Short: "Run Slurm task",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()

			switch type_ {
			case Solve:
				solveSlurmTaskSvc := do.MustInvoke[*services1.SolveSlurmTaskService](injector)
				solverSvc := do.MustInvoke[*solverServices.SolverService](injector)

				task, err := solveSlurmTaskSvc.GetTask(id)
				if err != nil && err == errorModule.ErrKeyNotFound {
					log.Fatal("Task ID not found")
				}

				timeout := int(task.Timeout.Seconds())
				if solverSvc.ShouldSkip(task.Encoding, task.Solver, timeout) {
					fmt.Println("Slurk task: skipped", task.Solver, task.Encoding)
					return
				}

				solverSvc.TrackedInvoke(task.Encoding, solver.Solver(task.Solver), timeout)
			case Cube:
				cubeSlurmTaskSvc := do.MustInvoke[*services2.CubeSlurmTaskService](injector)
				cuberSvc := do.MustInvoke[*cuberServices.CuberService](injector)

				task, err := cubeSlurmTaskSvc.GetTask(id)
				if err != nil && err == errorModule.ErrKeyNotFound {
					log.Fatal("Task ID not found")
				}

				if cuberSvc.ShouldSkip(task.Encoding, task.Threshold) {
					fmt.Println("Slurk task: skipped", task.Threshold, task.Encoding)
					return
				}

				cuberSvc.TrackedInvoke(task.Encoding, task.Threshold, task.Timeout)
			}
		},
	}

	cmd.Flags().IntVarP(&id, "id", "i", 1, "ID of the task")
	cmd.Flags().StringVarP(&type_, "type", "t", "solve", "Type of the task")

	return cmd
}
