package cmd

import (
	errorServices "benchmark/internal/error/services"
	"benchmark/internal/injector"
	solverServices "benchmark/internal/solver/services"
	"errors"
	"io"
	"log"

	"github.com/samber/do"
	"github.com/spf13/cobra"
)

func initTaskCmd() *cobra.Command {
	const (
		Solve = "solve"
	)

	var (
		type_         string
		inputFilePath string
		tasksPerGroup int
		groupId       int
	)

	cmd := &cobra.Command{
		Use:   "task",
		Short: "Run task",
		Run: func(cmd *cobra.Command, args []string) {
			injector := injector.New()

			switch type_ {
			case Solve:
				solverSvc := do.MustInvoke[*solverServices.SolverService](injector)
				errorSvc := do.MustInvoke[*errorServices.ErrorService](injector)

				taskIdsStart := (groupId-1)*tasksPerGroup + 1
				taskIdsEnd := groupId * tasksPerGroup
				tasksCounter := 0
				log.Printf("Started with task IDs %d to %d\n", taskIdsStart, taskIdsEnd)
				for i := taskIdsStart; i <= taskIdsEnd; i++ {
					// TODO: Do a batch search
					task, err := solverSvc.GetTask(inputFilePath, i)
					if err != nil && errors.Is(err, io.EOF) {
						continue
					}
					errorSvc.Fatal(err, "Task: failed to get task from the taskset")

					solverSvc.TrackedInvoke(task.Encoding, task.Solver, int(task.MaxRuntime.Seconds()))

					tasksCounter++
				}
				log.Printf("Ran %d tasks\n", tasksCounter)
			}
		},
	}

	cmd.Flags().StringVarP(&type_, "type", "t", "solve", "Type of the task")
	cmd.Flags().StringVarP(&inputFilePath, "input-file", "i", "x.tasks", "Input file path that holds the tasks")
	cmd.Flags().IntVarP(&tasksPerGroup, "tasks-per-group", "n", 1, "Number of tasks per group")
	cmd.Flags().IntVarP(&groupId, "group-id", "g", 1, "ID of the group")

	return cmd
}
