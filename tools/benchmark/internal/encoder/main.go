package encoder

import (
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

type EncodingPromise struct {
	Encoding string
}

func (encodingPromise EncodingPromise) Get() string {
	return encodingPromise.Encoding
}

func (encodingPromise EncodingPromise) GetPath() string {
	return encodingPromise.Encoding
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
