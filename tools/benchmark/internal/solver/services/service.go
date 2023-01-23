package services

import (
	"benchmark/internal/consts"
	"benchmark/internal/encoder"
	errorModule "benchmark/internal/error"
	"benchmark/internal/pipeline"
	"benchmark/internal/simplifier"
	"benchmark/internal/slurm"
	"benchmark/internal/solution"
	solveslurmtask "benchmark/internal/solve_slurm_task"
	"benchmark/internal/solver"
	"context"
	"fmt"
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

func (solverSvc *SolverService) GetCmdInfo(encoding string, solver solver.Solver, solutionPath string) (string, []string) {
	config := solverSvc.configSvc.Config

	var binPath string
	var args string
	switch solver {
	case consts.Kissat:
		binPath = config.Paths.Bin.Kissat
		args = "-q"
	case consts.Cadical:
		binPath = config.Paths.Bin.Cadical
		args = "-q"
	case consts.CryptoMiniSat:
		binPath = config.Paths.Bin.CryptoMiniSat
		args = "--verb=0"
	case consts.MapleSat:
		binPath = config.Paths.Bin.MapleSat
		args = "-verb=0"
	case consts.Glucose:
		binPath = config.Paths.Bin.Glucose
		args = "-verb=0"
	}

	args += " " + encoding
	if solver == consts.MapleSat || solver == consts.Glucose {
		args += " " + solutionPath
	}
	args_ := strings.Fields(args)

	return binPath, args_
}

func (solverSvc *SolverService) Invoke(encoding string, solver_ solver.Solver, timeout int) (string, time.Duration, solver.Result, int) {
	errorSvc := solverSvc.errorSvc
	solutionPath := path.Join("./tmp", path.Base(encoding)+"."+string(solver_)+".sol")
	binPath, solverArgs := solverSvc.GetCmdInfo(encoding, solver_, solutionPath)
	duration := time.Duration(timeout) * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Run and handle result
	cmd := exec.CommandContext(ctx, binPath, solverArgs...)
	pipe, err := cmd.StdoutPipe()
	solverSvc.errorSvc.Fatal(err, "Solver: failed to open pipe")
	startTime := time.Now()
	cmd.Start()

	if solver_ != consts.MapleSat && solver_ != consts.Glucose {
		err = solverSvc.filesystemSvc.WriteFromPipe(pipe, solutionPath)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to write from pipe")
	}

	var (
		result   solver.Result = consts.Fail
		exitCode int
	)
	err = cmd.Wait()
	runtime := time.Since(startTime)
	errorSvc.Handle(err, func(err error) {
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return
		}

		exitCode = exiterr.ExitCode()
		if exitCode == 10 {
			result = consts.Sat
		} else if exitCode == 20 {
			result = consts.Unsat
		} else {
			logrus.Error(err)
		}
	})

	return solutionPath, runtime, result, exitCode
}

func (solverSvc *SolverService) TrackedInvoke(encoding string, solver_ solver.Solver, timeout int) {
	solutionSvc := solverSvc.solutionSvc

	// Invoke
	solutionPath, runtime, result, exitCode := solverSvc.Invoke(encoding, solver_, timeout)
	runtime = runtime.Round(time.Millisecond)

	resultString := "Fail"
	if result == consts.Sat {
		resultString = "SAT"
	} else if result == consts.Unsat {
		resultString = "UNSAT"
	}

	verified := false
	checksum := ""
	if result == consts.Sat {
		err := solutionSvc.Normalize(solutionPath)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to normalize solution")
		info, err := solverSvc.encoderSvc.ProcessInstanceName(encoding)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to process instance name")

		needsReconstruction := false
		needsRemapping := false
		{
			if simplificationInfo, exists := info.Simplification.Get(); exists {
				needsReconstruction = simplificationInfo.Simplifier == simplifier.Cadical
				needsRemapping = simplificationInfo.Simplifier == simplifier.Satelite
			}
		}

		if needsReconstruction {
			info.Cubing = mo.None[encoder.CubingInfo]()
			info.CubeIndex = mo.None[int]()
			instance := path.Join(path.Dir(encoding), "..", solverSvc.encoderSvc.GetInstanceName(info))
			reconstructionPath := instance + ".rs.txt"
			err := solutionSvc.ReconstructAndVerify(solutionPath, reconstructionPath, []solution.Range{{Start: 1, End: 512}, {Start: 641, End: 768}})
			solverSvc.errorSvc.Fatal(err, "Solver: failed to reconstruct solution")
		} else if needsRemapping {
			info.Cubing = mo.None[encoder.CubingInfo]()
			info.CubeIndex = mo.None[int]()
			instance := path.Join(path.Dir(encoding), "..", solverSvc.encoderSvc.GetInstanceName(info))
			varMapPath := instance + ".var_map.txt"
			err := solutionSvc.RemapAndVerify(solutionPath, varMapPath)
			solverSvc.errorSvc.Fatal(err, "Solver: failed to remap variables in the solution")
		} else {
			info, err := solverSvc.encoderSvc.ProcessInstanceName(encoding)
			solverSvc.errorSvc.Fatal(err, "Solver: failed to process instance name")
			solutionFile, err := os.Open(solutionPath)
			solverSvc.errorSvc.Fatal(err, "Solver: failed to read the solution")

			verified, err = solutionSvc.Verify(solutionFile, info.Steps)
			solverSvc.errorSvc.Check(err, "Solver: verification failed")
		}

		checksum, err = solverSvc.filesystemSvc.Checksum(solutionPath)
		solverSvc.errorSvc.Fatal(err, "Solver: failed to calculate checksum of the solution "+solutionPath)
	}

	message := []any{"Solver:", solver_, resultString, exitCode, runtime, encoding}
	if verified {
		message = append(message, "verified")
	}
	logrus.Println(message)

	// Store in the database
	instanceName := strings.TrimSuffix(path.Base(encoding), ".cnf")
	solutionSvc.Register(encoding, solver_, solver.Solution{
		Runtime:      runtime,
		Result:       result,
		Solver:       solver_,
		ExitCode:     exitCode,
		InstanceName: instanceName,
		Verified:     verified,
		Checksum:     checksum,
	})
}

func (solverSvc *SolverService) Loop(encodingPromises []pipeline.EncodingPromise, parameters pipeline.Solving, handler func(encodingPromise pipeline.EncodingPromise, solver solver.Solver)) {
	for _, promise := range encodingPromises {
		for _, solver := range parameters.Solvers {
			handler(promise, solver)
		}
	}
}

func (solverSvc *SolverService) ShouldSkip(encoding string, solver_ solver.Solver, timeout int) bool {
	solutionSvc := solverSvc.solutionSvc
	errorSvc := solverSvc.errorSvc

	solution, err := solutionSvc.Find(encoding, solver_)
	// Don't skip if there is no solution
	if err == errorModule.ErrKeyNotFound || err == os.ErrNotExist {
		return false
	}

	// Handle errors
	errorSvc.Fatal(err, "Solver: failed to search the solution")

	// Skip solved solutions
	if err == nil && (solution.Result == consts.Sat || solution.Result == consts.Unsat) {
		return true
	}

	// Skip failed solutions: 10 seconds is the threshold
	if err == nil && solution.Result == consts.Fail && (timeout-int(solution.Runtime.Seconds())) < 10 {
		return true
	}

	return false
}

func (solverSvc *SolverService) RunSlurm(previousPipeOutput pipeline.SlurmPipeOutput, parameters pipeline.Solving) pipeline.SlurmPipeOutput {
	slurmSvc := solverSvc.slurmSvc
	errorSvc := solverSvc.errorSvc
	config := solverSvc.configSvc.Config
	encodingPromises, ok := previousPipeOutput.Values.([]pipeline.EncodingPromise)
	if !ok {
		log.Fatal("Solver: invalid input")
	}
	dependencies := previousPipeOutput.Jobs

	tasks := []solveslurmtask.Task{}
	i := 0
	numOfPromises := len(encodingPromises)
	solverSvc.Loop(encodingPromises, parameters, func(encodingPromise pipeline.EncodingPromise, solver_ solver.Solver) {
		// Note: We aren't checking if this task is already solved, since we'd have to retrieve the promised encoding, triggering the generation of cube encodings that are expensive on the FS to produce
		timeout := time.Duration(parameters.Timeout) * time.Second
		task := solveslurmtask.Task{
			EncodingPromise: encodingPromise,
			Solver:          solver_,
			Timeout:         timeout,
		}

		i += 1
		logrus.Printf("Solver: [%d/%d] tasks processed", i, numOfPromises)

		// Check if a task should be skipped
		if !parameters.Redundant && solverSvc.ShouldSkip(encodingPromise.GetPath(), solver_, parameters.Timeout) {
			return
		}

		// Prevent overwriting of any existing task
		taskId := solverSvc.solveSlurmTaskSvc.GenerateId(task)
		if _, err := solverSvc.solveSlurmTaskSvc.Get(taskId); err == nil || (err != nil && err != errorModule.ErrKeyNotFound) {
			return
		}

		tasks = append(tasks, task)
	})

	err := solverSvc.solveSlurmTaskSvc.AddMultiple(tasks)
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

func (solverSvc *SolverService) RunRegular(encodingPromises []pipeline.EncodingPromise, parameters pipeline.Solving) {
	logrus.Println("Solver: started")
	pool := pond.New(parameters.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	dependencies := map[string]interface{}{
		"CubeSelectorService": solverSvc.cubeSelectorSvc,
	}

	err := solverSvc.filesystemSvc.PrepareTempDir()
	solverSvc.errorSvc.Fatal(err, "Solver: failed to prepare tmp dir")

	solverSvc.Loop(encodingPromises, parameters, func(encodingPromise pipeline.EncodingPromise, solver_ solver.Solver) {
		encoding := encodingPromise.Get(dependencies)
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
