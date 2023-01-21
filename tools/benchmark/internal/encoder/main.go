package encoder

import (
	"benchmark/internal/pipeline"

	"github.com/samber/mo"
)

type Name string

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
	Encoder        pipeline.Encoder
	Function       string
	Steps          int
	AdderType      pipeline.AdderType
	IsXorEnabled   bool
	TargetHash     string
	Dobbertin      mo.Option[DobbertinInfo]
	Simplification mo.Option[SimplificationInfo]
	CubeIndex      mo.Option[int]
	Cubing         mo.Option[CubingInfo]
}
