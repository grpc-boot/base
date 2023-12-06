package msg

import "time"

type Handler[V Value] func() (V, error)

func Bool(value interface{}) bool {
	val, _ := value.(bool)
	return val
}

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

func String(value interface{}) string {
	val, _ := value.(string)
	return val
}

func Bytes(value interface{}) []byte {
	val, _ := value.([]byte)
	return val
}

func Slice(value interface{}) []interface{} {
	val, _ := value.([]interface{})
	return val
}

func MsgMap(value interface{}) Map {
	switch val := value.(type) {
	case map[string]interface{}:
		return Map(val)
	case Map:
		return val
	}

	return Map{}
}

func StringMap(value interface{}) map[string]string {
	val, _ := value.(map[string]string)
	return val
}

func Time(value interface{}) time.Time {
	t, _ := value.(time.Time)
	return t
}
