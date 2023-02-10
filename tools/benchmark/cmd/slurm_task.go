package cmd

import (
	"github.com/spf13/cobra"
)

func initSlurmTaskCmd() *cobra.Command {
	const (
		Solve = "solve"
		Cube  = "cube"
	)

	var type_ string

	cmd := &cobra.Command{
		Use:   "slurm-task",
		Short: "Run Slurm task",
		Run: func(cmd *cobra.Command, args []string) {
			// injector := injector.New()

			switch type_ {
			// case Solve:
			// solveSlurmTaskSvc := do.MustInvoke[*services1.SolveSlurmTaskService](injector)
			// solverSvc := do.MustInvoke[*solverServices.SolverService](injector)
			// errorSvc := do.MustInvoke[*services3.ErrorService](injector)

			// for {
			// 	startTime := time.Now()
			// 	maybeTask, taskId, err := solveSlurmTaskSvc.Book()
			// 	if err != nil {
			// 		errorSvc.Fatal(err, "Slurm task: failed to book")
			// 	}
			// 	task, exists := maybeTask.Get()
			// 	if !exists {
			// 		logrus.Println("Slurm task: none to be booked")
			// 		break
			// 	}
			// 	logrus.Println("Slurm task: booked task", taskId)
			// 	logrus.Info("Slurm task: book took", time.Since(startTime))

			// 	startTime = time.Now()
			// 	encoding := task.Encoding.BasePath
			// 	logrus.Info("Slurm task: encoding promise get took", time.Since(startTime))
			// 	timeout := int(task.Timeout.Seconds())
			// 	// * Note: The tasks are assumed to have went through a skipping phase, so we aren't doing them here
			// 	solverSvc.TrackedInvoke(encoding, solver.Solver(task.Solver), timeout)
			// 	err = solveSlurmTaskSvc.Remove(taskId)
			// 	errorSvc.Fatal(err, "Slurm task: failed to remove "+taskId+" after completion")
			// }

			// case Cube:
			// 	cubeSlurmTaskSvc := do.MustInvoke[*services2.CubeSlurmTaskService](injector)
			// 	cuberSvc := do.MustInvoke[*cuberServices.CuberService](injector)

			// 	task, err := cubeSlurmTaskSvc.GetTask(id)
			// 	if err != nil && err == errorModule.ErrKeyNotFound {
			// 		log.Fatal("Task ID not found")
			// 	}

			// 	if cuberSvc.ShouldSkip(task.Encoding, task.Threshold) {
			// 		fmt.Println("Slurk task: skipped", task.Threshold, task.Encoding)
			// 		return
			// 	}

			// 	cuberSvc.TrackedInvoke(cuberServices.InvokeParameters{
			// 		Encoding:         task.Encoding,
			// 		Threshold:        task.Threshold,
			// 		Timeout:          task.Timeout,
			// 		MaxCubes:         task.MaxCubes,
			// 		MinRefutedLeaves: task.MinRefutedLeaves,
			// 	}, cuberServices.InvokeControl{})
			}
		},
	}

	cmd.Flags().StringVarP(&type_, "type", "t", "solve", "Type of the task")

	return cmd
}
