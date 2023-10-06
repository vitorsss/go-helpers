package entity

import "encoding/json"

type JSONStringWrapper[T any] struct {
	Content T
}

func (d *JSONStringWrapper[T]) UnmarshalJSON(data []byte) error {
	var dataStr string
	err := json.Unmarshal(data, &dataStr)
	if err != nil {
		return err
	}
	var content T
	err = json.Unmarshal([]byte(dataStr), &content)
	d.Content = content
	return err
}

func (d JSONStringWrapper[T]) MarshalJSON() ([]byte, error) {
	content, err := json.Marshal(d.Content)
	if err != nil {
		return content, err
	}
	return json.Marshal(string(content))
}
