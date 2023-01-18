package encoder

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
