package cache

func Int(value interface{}) int64 {
	switch val := value.(type) {
	case int64:
		return val
	case uint64:
		return int64(val)
	case int:
		return int64(val)
	case uint:
		return int64(val)
	case int32:
		return int64(val)
	case uint32:
		return int64(val)
	case int16:
		return int64(val)
	case uint16:
		return int64(val)
	case int8:
		return int64(val)
	case uint8:
		return int64(val)
	}

	return 0
}

func Uint(value interface{}) uint64 {
	switch val := value.(type) {
	case int64:
		return uint64(val)
	case uint64:
		return val
	case int:
		return uint64(val)
	case uint:
		return uint64(val)
	case int32:
		return uint64(val)
	case uint32:
		return uint64(val)
	case int16:
		return uint64(val)
	case uint16:
		return uint64(val)
	case int8:
		return uint64(val)
	case uint8:
		return uint64(val)
	}

	return 0
}

func Float(value interface{}) float64 {
	switch val := value.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	}

	return 0
}
