package services

import (
	"benchmark/internal/encoder"
	"benchmark/internal/solver"
	"bufio"
	"encoding/binary"
	"log"
	"math"
	"os"
	"path"
	"time"

	"github.com/samber/mo"
)

type Task struct {
	Encoding   encoder.Encoding
	Solver     solver.Solver
	MaxRuntime time.Duration
}

func (solverSvc *SolverService) AddTasks(tasks []Task) (string, error) {
	// Tasks file
	name := solverSvc.randomSvc.RandString(10)
	tasksFilePath := path.Join(solverSvc.configSvc.Config.Paths.Tmp, name+".tasks")
	tasksFile, err := os.OpenFile(tasksFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer tasksFile.Close()
	tasksWriter := bufio.NewWriter(tasksFile)

	// Tasks map file
	tasksFileMapPath := tasksFilePath + ".map"
	tasksFileMap, err := os.OpenFile(tasksFileMapPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return "", err
	}
	defer tasksFileMap.Close()
	tasksMapWriter := bufio.NewWriter(tasksFileMap)

	addressAccumulator := 0
	for _, task := range tasks {
		cubeThreshold := 0
		cubeIndex := 0
		if cube, exists := task.Encoding.Cube.Get(); exists {
			cubeThreshold = cube.Threshold
			cubeIndex = cube.Index
		}

		// Threshold
		thresholdBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(thresholdBytes, uint16(cubeThreshold))
		tasksWriter.Write(thresholdBytes)

		// Index
		indexBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(indexBytes, uint32(cubeIndex))
		tasksWriter.Write(indexBytes)

		// Solver
		solverBytes := make([]byte, 1)
		switch task.Solver {
		case solver.Kissat:
			solverBytes[0] = uint8(0)
		case solver.Cadical:
			solverBytes[0] = uint8(1)
		case solver.MapleSat:
			solverBytes[0] = uint8(2)
		case solver.Glucose:
			solverBytes[0] = uint8(3)
		case solver.CryptoMiniSat:
			solverBytes[0] = uint8(4)
		}
		tasksWriter.Write(solverBytes)

		// Timeout
		timeoutBytes := make([]byte, 4)
		timeoutSeconds := uint32(math.Round(task.MaxRuntime.Seconds()))
		binary.BigEndian.PutUint32(timeoutBytes, timeoutSeconds)
		tasksWriter.Write(timeoutBytes)

		// Base path
		basePathBytes := []byte(task.Encoding.BasePath)
		tasksWriter.Write(basePathBytes)

		// Add the ending address to the maps
		bytesCount := len(thresholdBytes) + len(indexBytes) + len(solverBytes) + len(timeoutBytes) + len(basePathBytes)
		addressAccumulator += bytesCount
		addressBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(addressBytes, uint32(addressAccumulator))
		tasksMapWriter.Write(addressBytes)
	}
	tasksWriter.Flush()
	tasksMapWriter.Flush()

	return tasksFilePath, nil
}

// TODO: Implement a multiple get version of this
func (solverSvc *SolverService) GetTask(tasksSetPath string, index int) (Task, error) {
	task := Task{}

	tasksFilePath := path.Join(tasksSetPath)
	tasksFile, err := os.OpenFile(tasksFilePath, os.O_RDONLY, 0644)
	if err != nil {
		return task, err
	}
	defer tasksFile.Close()

	// Tasks map file
	tasksFileMapPath := tasksFilePath + ".map"
	tasksFileMap, err := os.OpenFile(tasksFileMapPath, os.O_RDONLY, 0644)
	if err != nil {
		return task, err
	}
	defer tasksFileMap.Close()

	// Find the task bytes count
	addressBytes := make([]byte, 4)
	_, err = tasksFileMap.ReadAt(addressBytes, int64(index-1)*4)
	if err != nil {
		return task, err
	}
	endAddress := int64(binary.BigEndian.Uint32(addressBytes))

	// Find the offset
	var offset int64 = 0
	if index > 1 {
		addressBytes := make([]byte, 4)
		_, err = tasksFileMap.ReadAt(addressBytes, int64(index-2)*4)
		if err != nil {
			return task, err
		}
		offset = int64(binary.BigEndian.Uint32(addressBytes))
	}

	// Read the task bytes
	taskBytes := make([]byte, endAddress-offset)
	_, err = tasksFile.ReadAt(taskBytes, offset)
	if err != nil {
		log.Println(offset, endAddress)
		return task, err
	}

	// Construct the task
	cubeThreshold := int(binary.BigEndian.Uint16(taskBytes[0:2]))
	cubeIndex := int(binary.BigEndian.Uint32(taskBytes[2:6]))
	solver_ := uint8(taskBytes[6])
	task.MaxRuntime = time.Duration(int(binary.BigEndian.Uint32(taskBytes[7:11]))) * time.Second
	basePath := string(taskBytes[11:])

	task.Encoding.BasePath = basePath
	if cubeIndex != 0 && cubeThreshold != 0 {
		task.Encoding.Cube = mo.Some(encoder.Cube{
			Threshold: cubeThreshold,
			Index:     cubeIndex,
		})
	}

	// Get the solver

	switch solver_ {
	case 0:
		task.Solver = solver.Kissat
	case 1:
		task.Solver = solver.Cadical
	case 2:
		task.Solver = solver.MapleSat
	case 3:
		task.Solver = solver.Glucose
	case 4:
		task.Solver = solver.CryptoMiniSat
	}

	return task, nil
}
