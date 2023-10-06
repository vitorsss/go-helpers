package entity

type Map map[string]interface{}

func (e Map) Get(keys ...string) (interface{}, bool) {
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
		case map[string]interface{}:
			mValue = v
		case Map:
			mValue = v
		default:
			return nil, false
		}
		return mValue.Get(keys[1:]...)
	}
}

func (e Map) Set(value interface{}, keys ...string) bool {
	switch len(keys) {
	case 0:
		return false
	case 1:
		e[keys[0]] = value
		return true
	default:
		var iValue interface{}
		var ok bool
		if iValue, ok = e[keys[0]]; !ok {
			iValue = Map{}
		}
		var mValue Map
		switch v := iValue.(type) {
		case map[string]interface{}:
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
