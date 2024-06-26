package services

import (
	"bufio"
	"cryptanalysis/internal/cuber"
	"cryptanalysis/internal/encoder"
	"cryptanalysis/internal/solver"
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

// Important: Register new SAT Solver here
func solverToUint8(solver_ solver.Solver) uint8 {
	switch solver_ {
	case solver.Kissat:
		return uint8(0)
	case solver.Cadical:
		return uint8(1)
	case solver.MapleSat:
		return uint8(2)
	case solver.MapleSatCrypto:
		return uint8(3)
	case solver.Glucose:
		return uint8(4)
	case solver.CryptoMiniSat:
		return uint8(5)
	case solver.YalSat:
		return uint8(6)
	case solver.PalSat:
		return uint8(7)
	case solver.LSTechMaple:
		return uint8(8)
	case solver.KissatCF:
		return uint8(9)
	case solver.Lingeling:
		return uint8(10)
	}

	log.Fatal("Solver: couldn't identify the SAT solver for the task")
	return uint8(0)
}

// Important: Register new SAT Solver here
func uint8ToSolver(solver_ uint8) solver.Solver {
	switch solver_ {
	case 0:
		return solver.Kissat
	case 1:
		return solver.Cadical
	case 2:
		return solver.MapleSat
	case 3:
		return solver.MapleSatCrypto
	case 4:
		return solver.Glucose
	case 5:
		return solver.CryptoMiniSat
	case 6:
		return solver.YalSat
	case 7:
		return solver.PalSat
	case 8:
		return solver.LSTechMaple
	case 9:
		return solver.KissatCF
	case 10:
		return solver.Lingeling
	}

	log.Fatal("Solver: couldn't identify the SAT solver for the task")
	return solver.Kissat
}

// Important: Register new threshold type here
func thresholdTypeToUint8(thresholdType cuber.ThresholdType) uint8 {
	switch thresholdType {
	case cuber.CutoffVars:
		return uint8(0)
	case cuber.CutoffDepth:
		return uint8(1)
	}

	return uint8(0)
}

// Important: Register new threshold type here
func uint8ToThresholdType(thresholdType_ uint8) cuber.ThresholdType {
	switch thresholdType_ {
	case 0:
		return cuber.CutoffVars
	case 1:
		return cuber.CutoffDepth
	}

	return cuber.CutoffVars
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
		log.Println("Solver: Add task ", task)
		cubeThreshold := 0
		cubeThresholdType := cuber.ThresholdType(cuber.CutoffVars)
		cubeIndex := 0
		if cube, exists := task.Encoding.Cube.Get(); exists {
			cubeThresholdType = cube.ThresholdType
			cubeThreshold = cube.Threshold
			cubeIndex = cube.Index
		}

		// Threshold type
		thresholdTypeBytes := make([]byte, 1)
		thresholdTypeBytes[0] = thresholdTypeToUint8(cubeThresholdType)
		tasksWriter.Write(thresholdTypeBytes)

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
		solverBytes[0] = solverToUint8(task.Solver)
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
		bytesCount := len(thresholdTypeBytes) + len(thresholdBytes) + len(indexBytes) + len(solverBytes) + len(timeoutBytes) + len(basePathBytes)
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
	cubeThresholdType := uint8ToThresholdType(taskBytes[0])
	cubeThreshold := int(binary.BigEndian.Uint16(taskBytes[1:3]))
	cubeIndex := int(binary.BigEndian.Uint32(taskBytes[3:7]))
	solver_ := uint8(taskBytes[7])
	task.MaxRuntime = time.Duration(int(binary.BigEndian.Uint32(taskBytes[8:12]))) * time.Second
	basePath := string(taskBytes[12:])

	task.Encoding.BasePath = basePath
	if cubeIndex != 0 && cubeThreshold != 0 {
		task.Encoding.Cube = mo.Some(encoder.Cube{
			Threshold:     cubeThreshold,
			Index:         cubeIndex,
			ThresholdType: cubeThresholdType,
		})
	}

	// Get the solver
	task.Solver = uint8ToSolver(solver_)

	return task, nil
}
