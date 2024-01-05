package entity

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type JSONStringWrapper[T any] struct {
	Content T
}

func (d *JSONStringWrapper[T]) UnmarshalJSON(data []byte) error {
	var dataStr string
	err := json.Unmarshal(data, &dataStr)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal data string")
	}
	var content T
	err = json.Unmarshal([]byte(dataStr), &content)
	d.Content = content
	return errors.Wrap(err, "failed to unmarshal content")
}

func (d JSONStringWrapper[T]) MarshalJSON() ([]byte, error) {
	content, err := json.Marshal(d.Content)
	if err != nil {
		return content, errors.Wrap(err, "failed to marshal content")
	}
	data, err := json.Marshal(string(content))
	if err != nil {
		return data, errors.Wrap(err, "failed to marshal data string")
	}
	return data, nil
}
