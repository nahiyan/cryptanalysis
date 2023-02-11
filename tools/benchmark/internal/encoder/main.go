package encoder

import (
	"benchmark/internal/solver"
	"errors"
	"fmt"
	"path"

	"github.com/samber/mo"
)

type AdderType string
type Encoder string

// Encoders
const (
	SaeedE   = "saeed_e" // Short for Saeed's Encoder
	Transalg = "transalg"
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
	Function       string
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
	Index     int
	Threshold int
}

type Encoding struct {
	BasePath string
	Cube     mo.Option[Cube]
}

func (encoding Encoding) GetLogPath(logsDir string, solver_ mo.Option[solver.Solver]) string {
	basePathInLogsDir := path.Join(logsDir, path.Base(encoding.BasePath))
	basePathWithoutExt := basePathInLogsDir[:len(basePathInLogsDir)-4]

	solver__ := ""
	if solver___, exists := solver_.Get(); exists {
		solver__ = "." + string(solver___)
	}

	if cube, exists := encoding.Cube.Get(); exists {
		logFilePath := basePathWithoutExt + fmt.Sprintf(".march_n%d.cube%d%s.log", cube.Threshold, cube.Index, solver__)
		return logFilePath
	}

	return basePathWithoutExt + solver__ + ".log"
}

func (encoding Encoding) GetName() string {
	if cube, exists := encoding.Cube.Get(); exists {
		return path.Join(encoding.BasePath + fmt.Sprintf(".march_n%d.cube%d", cube.Threshold, cube.Index))
	}

	return encoding.BasePath[:len(encoding.BasePath)-4]
}

func (encoding Encoding) GetCubesetPath(cubesetDir string) (string, error) {
	cube, exists := encoding.Cube.Get()
	if !exists {
		return "", errors.New("encoding isn't cubed")
	}

	cubesetPath := path.Join(cubesetDir, path.Base(encoding.BasePath)+fmt.Sprintf(".march_n%d.cubes", cube.Threshold))

	return cubesetPath, nil
}
