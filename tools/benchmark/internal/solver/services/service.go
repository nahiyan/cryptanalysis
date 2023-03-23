package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/pipeline"
	"benchmark/internal/solver"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/bitfield/script"
	"github.com/samber/mo"
)

// Important: Register new SAT Solver here
func (solverSvc *SolverService) GetCmdInfo(solver_ solver.Solver, solutionPath string, timeout int) (string, []string) {
	config := solverSvc.configSvc.Config

	var binPath string
	args := []string{}
	switch solver_ {
	case solver.Kissat:
		binPath = config.Paths.Bin.Kissat
		args = append(args, fmt.Sprintf("--time=%d", timeout))
	case solver.Cadical:
		binPath = config.Paths.Bin.Cadical
		args = append(args, "-t", fmt.Sprintf("%d", timeout))
	case solver.CryptoMiniSat:
		binPath = config.Paths.Bin.CryptoMiniSat
		args = append(args, "--maxtime", fmt.Sprintf("%d", timeout))
	case solver.MapleSat:
		binPath = config.Paths.Bin.MapleSat
		args = append(args, "-model", fmt.Sprintf("-cpu-lim=%d", timeout))
	case solver.Glucose:
		binPath = config.Paths.Bin.Glucose
		args = append(args, "-model", fmt.Sprintf("-cpu-lim=%d", timeout))
	case solver.YalSat:
		binPath = config.Paths.Bin.YalSat
		args = append(args, "--witness=1")
	case solver.PalSat:
		binPath = config.Paths.Bin.PalSat
		args = append(args, "--witness=1")
	case solver.LSTechMaple:
		binPath = config.Paths.Bin.LSTechMaple
		args = append(args, fmt.Sprintf("-cpu-lim=%d", timeout))
	case solver.KissatCF:
		binPath = config.Paths.Bin.KissatCF
		args = append(args, "-v")
	}

	return binPath, args
}

func (solverSvc *SolverService) Invoke(encoding encoder.Encoding, solver_ solver.Solver, timeout int) (solver.Result, int) {
	config := solverSvc.configSvc.Config
	solutionsDir := solverSvc.configSvc.Config.Paths.Solutions
	solutionPath := path.Join(solutionsDir, path.Base(encoding.GetName())+"."+string(solver_)+".sol")
	binPath, solverArgs := solverSvc.GetCmdInfo(solver_, solutionPath, timeout)
	duration := time.Duration(timeout) * time.Second

	// Local search
	if solver_ == solver.Kissat && config.Solver.Kissat.LocalSearch {
		solverArgs = append(solverArgs, "--walkinitially=true", fmt.Sprintf("--walkeffort=%d", config.Solver.Kissat.LocalSearchEffort))
	} else if solver_ == solver.Cadical && config.Solver.Cadical.LocalSearchRounds > 0 {
		solverArgs = append(solverArgs, fmt.Sprintf("-L%d", config.Solver.Cadical.LocalSearchRounds))
	} else if solver_ == solver.CryptoMiniSat && config.Solver.CryptoMiniSat.LocalSearch {
		solverArgs = append(solverArgs, "--sls=1", "--slstype="+config.Solver.CryptoMiniSat.LocalSearchType)
	}

	// Command context
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

	// Start the command
	startTime := time.Now()
	cmd.Start()

	if cube, exists := encoding.Cube.Get(); exists {
		// Handle cubes
		cubesetPath, err := encoding.GetCubesetPath(solverSvc.configSvc.Config.Paths.Cubesets)
		solverSvc.errorSvc.Fatal(err, "Solver: can't get cubeset path of an encoding that isn't cubed")

		err = solverSvc.cubeSelectorSvc.EncodingFromCube(encoding.BasePath, cubesetPath, cube.Index, stdinPipe)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to construct instance from cube")
	} else {
		// Handle regular files
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

	cmd.Wait()
	exitCode := cmd.ProcessState.ExitCode()
	result := solver.Result(solver.Fail)
	if exitCode == 10 {
		result = solver.Sat
	} else if exitCode == 20 {
		result = solver.Unsat
	}

	runtimeSeconds := time.Since(startTime).Round(time.Millisecond).Seconds()
	script.Echo(fmt.Sprintf("\nInfo: Ended after %.2f seconds with exit code %d", runtimeSeconds, exitCode)).AppendFile(logFilePath)

	return result, exitCode
}

func (solverSvc *SolverService) TrackedInvoke(encoding encoder.Encoding, solver_ solver.Solver, timeout int) {
	result, exitCode := solverSvc.Invoke(encoding, solver_, timeout)
	solverSvc.logSvc.SolveResult(encoding, solver_, exitCode, result)
}

func (solverSvc *SolverService) Loop(encodings []encoder.Encoding, parameters pipeline.SolveParams, handler func(encoding encoder.Encoding, solver solver.Solver)) {
	for _, encoding := range encodings {
		for _, solver := range parameters.Solvers {
			handler(encoding, solver)
		}
	}
}

func (solverSvc *SolverService) ShouldSkip(encoding encoder.Encoding, solver_ solver.Solver, maxRunTime time.Duration) bool {
	logFilePath := encoding.GetLogPath(solverSvc.configSvc.Config.Paths.Logs, mo.Some(solver_))
	var (
		result  solver.Result
		runTime time.Duration
		err     error
	)

	if solverSvc.combinedLogsSvc.IsLoaded() {
		result, _, runTime, err = solverSvc.ParseLogFromCombinedLog(path.Base(logFilePath), solver_, nil)
		if err != nil {
			return false
		}
	} else {
		result, _, runTime, err = solverSvc.ParseLogFromFile(logFilePath, solver_, nil)
		if err != nil {
			return false
		}
	}

	isSolved := result == solver.Unsat || result == solver.Sat
	if isSolved {
		return true
	}

	// See if it was attempted but timed out, being strict up to a second
	if diff := maxRunTime - runTime; diff < 1 {
		return true
	}

	return false
}

func (solverSvc *SolverService) RunSlurm(encodings []encoder.Encoding, parameters pipeline.SolveParams) {
	config := solverSvc.configSvc.Config
	dirs := []string{config.Paths.Solutions, solverSvc.configSvc.Config.Paths.Logs, solverSvc.configSvc.Config.Paths.Tmp}
	err := solverSvc.filesystemSvc.PrepareDirs(dirs)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to prepare directory for storing the solutions, logs, and tasks")

	// Select the unfinished tasks and skip the rest
	tasks := []Task{}
	pool := pond.New(config.Solver.Slurm.NumTaskSelectWorkers, 1000, pond.IdleTimeout(100*time.Millisecond))
	lock := sync.Mutex{}
	counter := 1
	numEncodings := len(encodings)
	solverSvc.Loop(encodings, parameters, func(encoding encoder.Encoding, solver_ solver.Solver) {
		pool.Submit(func(encoding encoder.Encoding, solver_ solver.Solver, index int, numEncodings int) func() {
			return func() {
				log.Printf("Solver: [%d/%d] processing task", index, numEncodings)
				if !parameters.Redundant && solverSvc.ShouldSkip(encoding, solver_, time.Duration(parameters.Timeout)*time.Second) {
					return
				}

				lock.Lock()
				tasks = append(tasks, Task{
					Encoding:   encoding,
					Solver:     solver_,
					MaxRuntime: time.Duration(parameters.Timeout) * time.Second,
				})
				lock.Unlock()
			}
		}(encoding, solver_, counter, numEncodings))
		counter++
	})
	pool.StopAndWait()
	log.Printf("Solver: Skipped %d tasks", (counter-1)-len(tasks))

	tasksSetPath, err := solverSvc.AddTasks(tasks)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to generate the taskset file")
	slurmMaxJobs := config.Slurm.MaxJobs
	numConcurrentTasks := int(math.Min(float64(parameters.Workers), float64(slurmMaxJobs)))
	timeout := parameters.Timeout
	// TODO: Prove that there will always be at least 1 task for the last worker
	tasksPerWorker := int(math.Ceil(float64(len(tasks)) / float64(parameters.Workers)))
	command := fmt.Sprintf(
		"%s task -t solve -i %s -n %d -g ${SLURM_ARRAY_TASK_ID}",
		config.Paths.Bin.Benchmark,
		tasksSetPath,
		tasksPerWorker)
	jobFilePath, err := solverSvc.slurmSvc.GenerateJob(
		command,
		numConcurrentTasks,
		1,
		1,
		config.Slurm.WorkerMemory,
		timeout)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to create slurm job file")

	jobId, err := solverSvc.slurmSvc.ScheduleJob(jobFilePath, nil)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to schedule the job")
	log.Printf("Solver: scheduled job with ID %d with %d tasks per worker\n", jobId, tasksPerWorker)
}

func (solverSvc *SolverService) RunRegular(encodings []encoder.Encoding, parameters pipeline.SolveParams) {
	dirs := []string{solverSvc.configSvc.Config.Paths.Solutions, solverSvc.configSvc.Config.Paths.Logs}
	err := solverSvc.filesystemSvc.PrepareDirs(dirs)
	solverSvc.errorSvc.Fatal(err, "Solver: failed to prepare directory for storing the solutions and logs")

	log.Println("Solver: started")
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))

	solverSvc.Loop(encodings, parameters, func(encoding encoder.Encoding, solver_ solver.Solver) {
		if !parameters.Redundant && solverSvc.ShouldSkip(encoding, solver_, time.Duration(parameters.Timeout)*time.Second) {
			log.Println("Solver: skipped", encoding, "with "+string(solver_))
			return
		}

		pool.Submit(func() {
			solverSvc.TrackedInvoke(encoding, solver_, parameters.Timeout)
		})
	})

	pool.StopAndWait()
	log.Println("Solver: stopped")
}

func (solverSvc *SolverService) Run(encodings []encoder.Encoding, asSlurmJobs bool, parameters pipeline.SolveParams) {
	if asSlurmJobs {
		solverSvc.RunSlurm(encodings, parameters)
		return
	}

	solverSvc.RunRegular(encodings, parameters)
}
