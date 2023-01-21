package services

type EncodingPromise struct {
	Encoding string
}

func (encodingPromise EncodingPromise) Get(dependencies map[string]interface{}) string {
	return encodingPromise.Encoding
}

func (encodingPromise EncodingPromise) GetPath() string {
	return encodingPromise.Encoding
}
