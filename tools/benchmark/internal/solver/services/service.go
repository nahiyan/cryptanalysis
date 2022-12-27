package services

import (
	"benchmark/internal/consts"
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

func (solverSvc *SolverService) RunSlurm() {

}

func (solverSvc *SolverService) RunRegular() {
	solutionSvc := solverSvc.solutionSvc

	fmt.Println("Solver: Started")

	pool := pond.New(solverSvc.Settings.Workers, 1000)
	solverSvc.Loop(func(encoding string, solver_ solver.Solver) {
		pool.Submit(func() {
			runtime, result := solverSvc.Invoke(encoding, solver_)

			// TODO: Store in database
			solutionSvc.Register(encoding, solver.Solution{
				Runtime: runtime,
				Result:  result,
			})
		})
	})
	pool.StopAndWait()

	fmt.Println("Solver: Stopped")
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
