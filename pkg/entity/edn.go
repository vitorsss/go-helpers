package entity

import (
	"bytes"

	"github.com/pkg/errors"
	"olympos.io/encoding/edn"
)

type Set[T any] []T

func (e Set[T]) MarshalEDN() ([]byte, error) {
	bb := bytes.NewBuffer([]byte("#{"))
	for _, value := range e {
		data, err := edn.Marshal(value)
		if err != nil {
			return nil, err
		}
		_, err = bb.Write(append(data, ' '))
		if err != nil {
			return nil, err
		}
	}
	data := bb.Bytes()
	data[len(data)-1] = '}'
	return data, nil
}

func (e *Set[T]) UnmarshalEDN(data []byte) error {
	if data[0] != '#' || data[1] != '{' || data[len(data)-1] != '}' {
		return errors.New("missing boundary chars")
	}
	data[0] = ' '
	data[1] = '['
	data[len(data)-1] = ']'
	result := []T{}
	err := edn.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	*e = result
	return nil
}
