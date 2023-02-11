package services

import (
	"benchmark/internal/encoder"
	errorModule "benchmark/internal/error"
	"benchmark/internal/pipeline"
	"benchmark/internal/slurm"
	solveslurmtask "benchmark/internal/solve_slurm_task"
	"benchmark/internal/solver"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/alitto/pond"
	"github.com/samber/mo"
	"github.com/sirupsen/logrus"
)

// TODO: Refactor
func (solverSvc *SolverService) GetCmdInfo(solver_ solver.Solver, solutionPath string) (string, []string) {
	config := solverSvc.configSvc.Config

	var binPath string
	args := ""
	switch solver_ {
	case solver.Kissat:
		binPath = config.Paths.Bin.Kissat
		// args = "-q"
	case solver.Cadical:
		binPath = config.Paths.Bin.Cadical
		// args = "-q"
	case solver.CryptoMiniSat:
		binPath = config.Paths.Bin.CryptoMiniSat
		// args = "--verb=0"
	case solver.MapleSat:
		binPath = config.Paths.Bin.MapleSat
		// args = "-verb=0"
	case solver.Glucose:
		binPath = config.Paths.Bin.Glucose
		// args = "-verb=0"
	}

	// args += " " + encoding
	if solver_ == solver.MapleSat || solver_ == solver.Glucose {
		args += " " + solutionPath
	}
	args_ := strings.Fields(args)

	return binPath, args_
}

func (solverSvc *SolverService) Invoke(encoding encoder.Encoding, solver_ solver.Solver, timeout int) (solver.Result, int) {
	errorSvc := solverSvc.errorSvc
	solutionsDir := solverSvc.configSvc.Config.Paths.Solutions
	solutionPath := path.Join(solutionsDir, path.Base(encoding.GetName())+"."+string(solver_)+".sol")
	binPath, solverArgs := solverSvc.GetCmdInfo(solver_, solutionPath)
	duration := time.Duration(timeout) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Run and handle the result
	cmd := exec.CommandContext(ctx, binPath, solverArgs...)

	// Output pipe
	stdoutPipe, err := cmd.StdoutPipe()
	solverSvc.errorSvc.Fatal(err, "Solver: failed to open stdout pipe")

	// Input pipe
	stdinPipe, err := cmd.StdinPipe()
	solverSvc.errorSvc.Fatal(err, "Solver: failed to open stdin pipe")

	cmd.Start()

	if cube, exists := encoding.Cube.Get(); exists {
		cubesetPath, err := encoding.GetCubesetPath(solverSvc.configSvc.Config.Paths.Cubesets)
		solverSvc.errorSvc.Fatal(err, "Solver: can't get cubeset path of an encoding that isn't cubed")

		err = solverSvc.cubeSelectorSvc.EncodingFromCube(encoding.BasePath, cubesetPath, cube.Index, stdinPipe)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to construct instance from cube")
	} else {
		reader, err := os.OpenFile(encoding.BasePath, os.O_RDONLY, 0644)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to read the instance file")
		_, err = io.Copy(stdinPipe, reader)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to provide the instance to the solver")
	}
	stdinPipe.Close()

	// Write from stdout pipe to log file
	logFilePath := encoding.GetLogPath(solverSvc.configSvc.Config.Paths.Logs, mo.Some(solver_))
	err = solverSvc.filesystemSvc.WriteFromPipe(stdoutPipe, logFilePath)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to write from pipe")

	err = cmd.Wait()
	var (
		result   solver.Result = solver.Fail
		exitCode int
	)
	errorSvc.Handle(err, func(err error) {
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return
		}

		exitCode = exiterr.ExitCode()
		if exitCode == 10 {
			result = solver.Sat
		} else if exitCode == 20 {
			result = solver.Unsat
		} else {
			logrus.Error(err)
		}
	})

	return result, exitCode
}

func (solverSvc *SolverService) TrackedInvoke(encoding encoder.Encoding, solver_ solver.Solver, timeout int) {
	result, exitCode := solverSvc.Invoke(encoding, solver_, timeout)
	solverSvc.logSvc.SolveResult(encoding, solver_, exitCode, result)
}

func (solverSvc *SolverService) Loop(encodings []encoder.Encoding, parameters pipeline.SolveParams, handler func(encoding encoder.Encoding, solver solver.Solver)) {
	for _, promise := range encodings {
		for _, solver := range parameters.Solvers {
			handler(promise, solver)
		}
	}
}

// TODO: Read the log file and see if it actually finished writing
func (solverSvc *SolverService) ShouldSkip(encoding encoder.Encoding, solver_ solver.Solver, timeout int) bool {
	logFilePath := encoding.GetLogPath(solverSvc.configSvc.Config.Paths.Logs, mo.Some(solver_))
	result, _, err := solverSvc.ParseLog(logFilePath, solver_, nil)
	if err != nil {
		return false
	}

	isSolved := result == solver.Unsat || result == solver.Sat
	return isSolved
}

func (solverSvc *SolverService) RunSlurm(previousPipeOutput pipeline.SlurmPipeOutput, parameters pipeline.SolveParams) pipeline.SlurmPipeOutput {
	dirs := []string{solverSvc.configSvc.Config.Paths.Solutions, solverSvc.configSvc.Config.Paths.Logs}
	err := solverSvc.filesystemSvc.PrepareDirs(dirs)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to prepare directory for storing the solutions and logs")

	slurmSvc := solverSvc.slurmSvc
	errorSvc := solverSvc.errorSvc
	config := solverSvc.configSvc.Config
	encodings, ok := previousPipeOutput.Values.([]encoder.Encoding)
	if !ok {
		log.Fatal("Solver: invalid input")
	}
	dependencies := previousPipeOutput.Jobs

	tasks := []solveslurmtask.Task{}
	i := 0
	numOfPromises := len(encodings)
	solverSvc.Loop(encodings, parameters, func(encoding encoder.Encoding, solver_ solver.Solver) {
		// Note: We aren't checking if this task is already solved, since we'd have to retrieve the promised encoding, triggering the generation of cube encodings that are expensive on the FS to produce
		timeout := time.Duration(parameters.Timeout) * time.Second
		task := solveslurmtask.Task{
			Encoding: encoding,
			Solver:   solver_,
			Timeout:  timeout,
		}

		i += 1
		logrus.Printf("Solver: [%d/%d] tasks processed", i, numOfPromises)

		// Check if a task should be skipped
		if !parameters.Redundant && solverSvc.ShouldSkip(encoding, solver_, parameters.Timeout) {
			return
		}

		// Prevent overwriting of any existing task
		taskId := solverSvc.solveSlurmTaskSvc.GenerateId(task)
		if _, err := solverSvc.solveSlurmTaskSvc.Get(taskId); err == nil || (err != nil && err != errorModule.ErrKeyNotFound) {
			return
		}

		tasks = append(tasks, task)
	})

	err = solverSvc.solveSlurmTaskSvc.AddMultiple(tasks)
	errorSvc.Fatal(err, "Solver: failed to add slurm task")
	logrus.Println("Solver: added", len(tasks), "slurm tasks")

	slurmMaxJobs := config.Slurm.MaxJobs
	numTasks := int(math.Min(float64(parameters.Workers), float64(slurmMaxJobs)))
	timeout := parameters.Timeout
	jobFilePath, err := slurmSvc.GenerateJob(
		numTasks,
		1,
		1,
		300,
		timeout,
		fmt.Sprintf(
			"%s slurm-task -t solve",
			config.Paths.Bin.Benchmark))
	errorSvc.Fatal(err, "Solver: failed to create slurm job file")

	jobId, err := slurmSvc.ScheduleJob(jobFilePath, dependencies)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to schedule the job")
	logrus.Println("Solver: scheduled job with ID", jobId)

	return pipeline.SlurmPipeOutput{
		Jobs:   []slurm.Job{{Id: jobId}},
		Values: []string{},
	}
}

func (solverSvc *SolverService) RunRegular(encodings []encoder.Encoding, parameters pipeline.SolveParams) {
	dirs := []string{solverSvc.configSvc.Config.Paths.Solutions, solverSvc.configSvc.Config.Paths.Logs}
	err := solverSvc.filesystemSvc.PrepareDirs(dirs)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to prepare directory for storing the solutions and logs")

	logrus.Println("Solver: started")
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	solverSvc.Loop(encodings, parameters, func(encoding encoder.Encoding, solver_ solver.Solver) {
		if !parameters.Redundant && solverSvc.ShouldSkip(encoding, solver_, parameters.Timeout) {
			logrus.Println("Solver: skipped", encoding, "with "+string(solver_))
			return
		}

		pool.Submit(func() {
			solverSvc.TrackedInvoke(encoding, solver_, parameters.Timeout)
		})
	})

	pool.StopAndWait()
	logrus.Println("Solver: stopped")
}
