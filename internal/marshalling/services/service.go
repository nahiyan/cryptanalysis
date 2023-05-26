package services

import (
	"bytes"
	"encoding/gob"
)

func (marshallingSvc *MarshallingService) BinEncode(data interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (marshallingSvc *MarshallingService) BinDecode(source []byte, destinaton interface{}) error {
	buffer := bytes.NewBuffer(source)
	decoder := gob.NewDecoder(buffer)

	err := decoder.Decode(destinaton)
	return err
}
