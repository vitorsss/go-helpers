package entity

import (
	"fmt"
)

type Map map[string]any

func (e Map) Get(keys ...string) (any, bool) {
	switch len(keys) {
	case 0:
		return nil, false
	case 1:
		value, ok := e[keys[0]]
		return value, ok
	default:
		value, ok := e[keys[0]]
		if !ok {
			return nil, false
		}
		var mValue Map
		switch v := value.(type) {
		case map[string]any:
			mValue = v
		case Map:
			mValue = v
		default:
			return nil, false
		}
		return mValue.Get(keys[1:]...)
	}
}

func (e Map) GetString(keys ...string) (string, bool) {
	valueIf, ok := e.Get(keys...)
	if !ok {
		return "", ok
	}
	value, ok := valueIf.(string)
	return value, ok
}

func (e Map) GetStringer(keys ...string) (string, bool) {
	valueIf, ok := e.Get(keys...)
	if !ok {
		return "", ok
	}
	switch v := valueIf.(type) {
	case string:
		return v, ok
	case fmt.Stringer:
		return v.String(), ok
	case float64, float32:
		return fmt.Sprintf("%f", v), ok
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v), ok
	case bool:
		return fmt.Sprintf("%t", v), ok
	}
	return "", false
}

func (e Map) GetInt(keys ...string) (int, bool) {
	valueIf, ok := e.Get(keys...)
	if !ok {
		return 0, ok
	}
	switch value := valueIf.(type) {
	case float64:
		return int(value), true
	case int:
		return value, true
	}
	return 0, false
}

func (e Map) Contains(keys ...string) bool {
	_, ok := e.Get(keys...)
	return ok
}

func (e Map) Create(value any, keys ...string) bool {
	if e.Contains(keys...) {
		return false
	}
	return e.Set(value, keys...)
}

func (e Map) Set(value any, keys ...string) bool {
	switch len(keys) {
	case 0:
		return false
	case 1:
		e[keys[0]] = value
		return true
	default:
		var iValue any
		var ok bool
		if iValue, ok = e[keys[0]]; !ok {
			iValue = Map{}
		}
		var mValue Map
		switch v := iValue.(type) {
		case map[string]any:
			mValue = v
		case Map:
			mValue = v
		default:
			return false
		}
		if ok = mValue.Set(value, keys[1:]...); !ok {
			return false
		}
		e[keys[0]] = mValue
		return true
	}
}
