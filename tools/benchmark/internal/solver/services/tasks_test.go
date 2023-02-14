package services

import (
	"benchmark/internal/config/services"
	"benchmark/internal/encoder"
	"benchmark/internal/solver"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/samber/mo"
)

func TestOverall(t *testing.T) {
	svc := SolverService{
		configSvc: &services.ConfigService{},
	}

	tasksSetPath, err := svc.AddTasks([]task{
		{
			encoding: encoder.Encoding{
				BasePath: "lorem_ipsum.cnf",
			},
			solver:     solver.Kissat,
			maxRuntime: time.Duration(5000) * time.Second,
		},
		{
			encoding: encoder.Encoding{
				BasePath: "transalg_md4_41_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf",
				Cube: mo.Some(
					encoder.Cube{
						Threshold: 1234,
						Index:     11000000,
					},
				),
			},
			solver:     solver.Glucose,
			maxRuntime: time.Duration(30) * time.Second,
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
		if _, exists := task.encoding.Cube.Get(); exists {
			t.Fatalf("Cube info. shouldn't exist")
		}
		// Solver
		if task.solver != solver.Kissat {
			t.Fatalf("Expected Kissat, got %s", task.solver)
		}
		// Max. Runtime
		if task.maxRuntime.Seconds() != 5000 {
			t.Fatalf("Expected max. runtime = 5000, got %s", task.maxRuntime)
		}
		// Base path
		if task.encoding.BasePath != "lorem_ipsum.cnf" {
			t.Fatalf("Expected base path = lorem_ipsum.cnf, got '%s'", task.encoding.BasePath)
		}
	}

	// Task 2
	{
		task, err := svc.GetTask(tasksSetPath, 2)
		if err != nil {
			t.Fatal(err)
		}

		// Cube info.
		if cube, exists := task.encoding.Cube.Get(); exists {
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
		if task.solver != solver.Glucose {
			t.Fatalf("Expected Glucose, got %s", task.solver)
		}
		// Max. Runtime
		if task.maxRuntime.Seconds() != 30 {
			t.Fatalf("Expected max. runtime = 5000, got %s", task.maxRuntime)
		}
		// Base path
		if task.encoding.BasePath != "transalg_md4_41_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf" {
			t.Fatalf("Expected base path = transalg_md4_41_00000000000000000000000000000000_dobbertin31.cnf.cadical_c1000000.cnf, got '%s'", task.encoding.BasePath)
		}
	}

	// Task 3
	{
		_, err := svc.GetTask(tasksSetPath, 3)
		if err == nil {
			t.Fatal("The call to fetch task 3 should fail")
		}

		if !errors.Is(err, io.EOF) {
			t.Fatalf("expected EOF but got %s", err)
		}
	}
}
