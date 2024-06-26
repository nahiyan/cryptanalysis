package services

import (
	"cryptanalysis/internal/config/services"
	"cryptanalysis/internal/cuber"
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/solver"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/samber/mo"
)

// Important: Register new SAT Solver here
func TestOverall(t *testing.T) {
	svc := SolverService{
		configSvc: &services.ConfigService{},
	}

	tasksSetPath, err := svc.AddTasks([]Task{
		{
			Encoding: encoder.Encoding{
				BasePath: "lorem_ipsum.cnf",
			},
			Solver:     solver.Kissat,
			MaxRuntime: time.Duration(5000) * time.Second,
		},
		{
			Encoding: encoder.Encoding{
				BasePath: "transalg_md4_41_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf",
				Cube: mo.Some(
					encoder.Cube{
						Threshold: 1234,
						Index:     11000000,
					},
				),
			},
			Solver:     solver.YalSat,
			MaxRuntime: time.Duration(30) * time.Second,
		},
		{
			Encoding: encoder.Encoding{
				BasePath: "lorem_ipsum.cnf",
			},
			Solver:     solver.LSTechMaple,
			MaxRuntime: time.Duration(5000) * time.Second,
		},
		{
			Encoding: encoder.Encoding{
				BasePath: "lorem_ipsum.cnf",
			},
			Solver:     solver.KissatCF,
			MaxRuntime: time.Duration(5000) * time.Second,
		},
		{
			Encoding: encoder.Encoding{
				BasePath: "transalg_md4_41_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf",
				Cube: mo.Some(
					encoder.Cube{
						Threshold:     1234,
						Index:         11000000,
						ThresholdType: cuber.CutoffDepth,
					},
				),
			},
			Solver:     solver.YalSat,
			MaxRuntime: time.Duration(30) * time.Second,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tasksSetPath)
	defer os.Remove(tasksSetPath + ".map")

	// Task 1
	{
		task, err := svc.GetTask(tasksSetPath, 1)
		if err != nil {
			t.Fatal(err)
		}

		// Cube info.
		if _, exists := task.Encoding.Cube.Get(); exists {
			t.Fatalf("Cube info. shouldn't exist")
		}
		// Solver
		if task.Solver != solver.Kissat {
			t.Fatalf("Expected Kissat, got %s", task.Solver)
		}
		// Max. Runtime
		if task.MaxRuntime.Seconds() != 5000 {
			t.Fatalf("Expected max. runtime = 5000, got %s", task.MaxRuntime)
		}
		// Base path
		if task.Encoding.BasePath != "lorem_ipsum.cnf" {
			t.Fatalf("Expected base path = lorem_ipsum.cnf, got '%s'", task.Encoding.BasePath)
		}
	}

	// Task 2
	{
		task, err := svc.GetTask(tasksSetPath, 2)
		if err != nil {
			t.Fatal(err)
		}

		// Cube info.
		if cube, exists := task.Encoding.Cube.Get(); exists {
			if cube.Threshold != 1234 {
				t.Fatalf("Expected threshold = 1234, got %d", cube.Threshold)
			}
			if cube.Index != 11000000 {
				t.Fatalf("Expected index = 11000000, got %d", cube.Index)
			}
		} else {
			t.Fatal("Cube info. should exist")
		}
		// Solver
		if task.Solver != solver.YalSat {
			t.Fatalf("Expected YalSAT, got %s", task.Solver)
		}
		// Max. Runtime
		if task.MaxRuntime.Seconds() != 30 {
			t.Fatalf("Expected max. runtime = 5000, got %s", task.MaxRuntime)
		}
		// Base path
		if task.Encoding.BasePath != "transalg_md4_41_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf" {
			t.Fatalf("Expected base path = transalg_md4_41_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf, got '%s'", task.Encoding.BasePath)
		}
	}

	// Task 3
	{
		task, err := svc.GetTask(tasksSetPath, 3)
		if err != nil {
			t.Fatal(err)
		}

		if task.Solver != solver.LSTechMaple {
			t.Fatalf("Expected LSTech_Maple but got %s", task.Solver)
		}
	}

	// Task 4
	{
		task, err := svc.GetTask(tasksSetPath, 4)
		if err != nil {
			t.Fatal(err)
		}

		if task.Solver != solver.KissatCF {
			t.Fatalf("Expected Kissat_CF but got %s", task.Solver)
		}
	}

	// Task 5
	{
		task, err := svc.GetTask(tasksSetPath, 5)
		if err != nil {
			t.Fatal("The call to fetch task 5 should fail")
		}

		if cube, exists := task.Encoding.Cube.Get(); exists {
			if cube.Threshold != 1234 {
				t.Fatalf("Expected threshold = 1234, got %d", cube.Threshold)
			}
			if cube.Index != 11000000 {
				t.Fatalf("Expected index = 11000000, got %d", cube.Index)
			}
			if cube.ThresholdType != cuber.CutoffDepth {
				t.Fatalf("Expected threshold type of cutoff_depth, got %s", cube.ThresholdType)
			}
		} else {
			t.Fatal("Cube info. should exist")
		}
	}

	// Task 6
	{
		_, err := svc.GetTask(tasksSetPath, 6)
		if err == nil {
			t.Fatal("The call to fetch task 6 should fail")
		}

		if !errors.Is(err, io.EOF) {
			t.Fatalf("expected EOF but got %s", err)
		}
	}
}
