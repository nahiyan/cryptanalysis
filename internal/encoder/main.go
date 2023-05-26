package encoder

import (
	"cryptanalysis/internal/cuber"
	"cryptanalysis/internal/solver"
	"errors"
	"fmt"
	"path"

	"github.com/samber/mo"
)

type AdderType string
type Encoder string
type Function string

// Function
const (
	Md4    = "md4"
	Md5    = "md5"
	Sha256 = "sha256"
)

// Encoders
const (
	NejatiEncoder = "nejati_encoder" // Short for Saeed's Encoder
	Transalg      = "transalg"
)

// Adders
const (
	TwoOperand   = "two_operand"
	DotMatrix    = "dot_matrix"
	CounterChain = "counter_chain"
	Espresso     = "espresso"
)

type EncoderService interface {
	TestRun() []string
}

type SimplificationInfo struct {
	Simplifier string
	Conflicts  int
}

type CubingInfo struct {
	Threshold int
}

type DobbertinInfo struct {
	Bits int
}

type InstanceInfo struct {
	Encoder        Encoder
	Function       Function
	Steps          int
	AdderType      AdderType
	IsXorEnabled   bool
	TargetHash     string
	Dobbertin      mo.Option[DobbertinInfo]
	Simplification mo.Option[SimplificationInfo]
	CubeIndex      mo.Option[int]
	Cubing         mo.Option[CubingInfo]
}

type Cube struct {
	Index         int
	ThresholdType cuber.ThresholdType
	Threshold     int
}

type Encoding struct {
	BasePath string
	Cube     mo.Option[Cube]
}

func (encoding Encoding) GetLogPath(logsDir string, solver_ mo.Option[solver.Solver]) string {
	basePathInLogsDir := path.Join(logsDir, path.Base(encoding.BasePath))
	// basePathWithoutExt := basePathInLogsDir[:len(basePathInLogsDir)-4]

	solver__ := ""
	if solver___, exists := solver_.Get(); exists {
		solver__ = "." + string(solver___)
	}

	if cube, exists := encoding.Cube.Get(); exists {
		thresholdArg := "n"
		if cube.ThresholdType == cuber.CutoffDepth {
			thresholdArg = "d"
		}

		logFilePath := basePathInLogsDir + fmt.Sprintf(".march_%s%d.cubes.cube%d%s.log", thresholdArg, cube.Threshold, cube.Index, solver__)
		return logFilePath
	}

	return basePathInLogsDir + solver__ + ".log"
}

func (encoding Encoding) GetName() string {
	if cube, exists := encoding.Cube.Get(); exists {
		thresholdType := "n"
		if cube.ThresholdType == cuber.CutoffDepth {
			thresholdType = "d"
		}

		return path.Join(encoding.BasePath + fmt.Sprintf(".march_%s%d.cubes.cube%d", thresholdType, cube.Threshold, cube.Index))
	}

	return encoding.BasePath[:len(encoding.BasePath)-4]
}

func (encoding Encoding) GetCubesetPath(cubesetDir string) (string, error) {
	cube, exists := encoding.Cube.Get()
	if !exists {
		return "", errors.New("encoding isn't cubed")
	}

	thresholdType := "n"
	if cube.ThresholdType == cuber.CutoffDepth {
		thresholdType = "d"
	}

	cubesetPath := path.Join(cubesetDir, path.Base(encoding.BasePath)+fmt.Sprintf(".march_%s%d.cubes", thresholdType, cube.Threshold))

	return cubesetPath, nil
}
