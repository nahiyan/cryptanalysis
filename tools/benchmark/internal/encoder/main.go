package encoder

type Name string

type EncoderService interface {
	TestRun() []string
}
