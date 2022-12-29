package services

import (
	"benchmark/internal/consts"
	errorModule "benchmark/internal/error"
	"benchmark/internal/slurm"
	"benchmark/internal/solver"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/alitto/pond"
)

type Properties struct {
	Encodings []string
	Settings  solver.Settings
}

func (solverSvc *SolverService) GetCmdInfo(encoding string, solver solver.Solver) (string, []string) {
	config := solverSvc.configSvc.Config

	var binPath string
	var args__ string
	switch solver {
	case consts.Kissat:
		binPath = config.Paths.Bin.Kissat
		args__ = "-q"
	case consts.Cadical:
		binPath = config.Paths.Bin.Cadical
		args__ = "-q"
	case consts.MapleSat:
		binPath = config.Paths.Bin.MapleSat
		args__ = "-verb=0"
	case consts.CryptoMiniSat:
		binPath = config.Paths.Bin.CryptoMiniSat
		args__ = "--verb=0"
	case consts.Glucose:
		binPath = config.Paths.Bin.Glucose
		args__ = "-verb=0"
	}

	args_ := args__ + " " + encoding
	args := strings.Fields(args_)

	return binPath, args
}

func (solverSvc *SolverService) Invoke(encoding string, solver_ solver.Solver) (time.Duration, solver.Result) {
	filesystemSvc := solverSvc.filesystemSvc
	errorSvc := solverSvc.errorSvc
	binPath, solverArgs := solverSvc.GetCmdInfo(encoding, solver_)
	timeout := solverSvc.Settings.Timeout

	if !filesystemSvc.FileExists(binPath) {
		log.Fatalf("%s doesn't exist. Did you forget to compile it?", binPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Run and handle result
	cmd := exec.CommandContext(ctx, binPath, solverArgs...)
	startTime := time.Now()
	err := cmd.Run()
	var result solver.Result = consts.Fail
	errorSvc.Handle(err, func(err error) {
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return
		}

		exitCode := exiterr.ExitCode()
		if exitCode == 10 {
			result = consts.Sat
		} else if exitCode == 20 {
			result = consts.Unsat
		}
	})
	runtime := time.Since(startTime)

	return runtime, result
}

func (solverSvc *SolverService) Loop(handler func(encoding string, solver solver.Solver)) {
	for _, encoding := range solverSvc.Encodings {
		for _, solver := range solverSvc.Settings.Solvers {
			handler(encoding, solver)
		}
	}
}

func (solverSvc *SolverService) ShouldSkip(encoding string, solver_ solver.Solver) bool {
	solutionSvc := solverSvc.solutionSvc
	errorSvc := solverSvc.errorSvc

	solution, err := solutionSvc.Find(encoding, solver_)
	if err != nil && err != errorModule.ErrKeyNotFound {
		errorSvc.Fatal(err, "Solver: failed to search the solution")
	}

	if err == nil && (solution.Result == consts.Sat || solution.Result == consts.Unsat) {
		return true
	}

	// 10 seconds is the threshold
	if err == nil && solution.Result == consts.Fail && (solverSvc.Settings.Timeout.Seconds()-solution.Runtime.Seconds()) < 10 {
		return true
	}

	return false
}

func (solverSvc *SolverService) RunSlurm() {
	slurmSvc := solverSvc.slurmSvc
	errorSvc := solverSvc.errorSvc
	config := solverSvc.configSvc.Config

	solverSvc.Loop(func(encoding string, solver_ solver.Solver) {
		if solverSvc.ShouldSkip(encoding, solver_) {
			fmt.Println("Solver: skipped", encoding, "with "+string(solver_))
			return
		}

		solverSvc.slurmSvc.AddTask(slurm.Task{
			Encoding: encoding,
			Solver:   solver_,
			Timeout:  solverSvc.Settings.Timeout,
		})
	})

	timeout := int(solverSvc.Settings.Timeout.Seconds())
	jobFilePath, err := slurmSvc.GenerateJob(
		1,
		1,
		300,
		timeout,
		fmt.Sprintf("%s slurm-task -j ${SLURM_ARRAY_TASK_ID} -t %d", config.Paths.Bin.Benchmark, timeout))
	errorSvc.Fatal(err, "Solver: failed to create slurm job file")

	fmt.Println(jobFilePath)
}

func (solverSvc *SolverService) RunRegular() {
	solutionSvc := solverSvc.solutionSvc

	fmt.Println("Solver: started")

	pool := pond.New(solverSvc.Settings.Workers, 1000, pond.IdleTimeout(100*time.Millisecond))
	solverSvc.Loop(func(encoding string, solver_ solver.Solver) {
		if solverSvc.ShouldSkip(encoding, solver_) {
			fmt.Println("Solver: skipped", encoding, "with "+string(solver_))
			return
		}

		pool.Submit(func() {
			// Invoke
			runtime, result := solverSvc.Invoke(encoding, solver_)

			resultString := "Fail"
			if result == consts.Sat {
				resultString = "SAT"
			} else if result == consts.Unsat {
				resultString = "UNSAT"
			}

			fmt.Println("Solver:", solver_, resultString, runtime, encoding)

			// Store in the database
			solutionSvc.Register(encoding, solver_, solver.Solution{
				Runtime: runtime,
				Result:  result,
				Solver:  solver_,
			})
		})
	})
	pool.StopAndWait()

	fmt.Println("Solver: stopped")
}

func (solverSvc *SolverService) Run(encodings []string, settings solver.Settings) {
	solverSvc.Encodings = encodings
	solverSvc.Settings = settings

	switch settings.Platform {
	case consts.Regular:
		solverSvc.RunRegular()
	case consts.Slurm:
		solverSvc.RunSlurm()
	}
}
