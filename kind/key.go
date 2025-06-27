package kind

import (
	"math"
)

type Key interface {
	~string | Integer
}

func KeyHash(key interface{}) uint32 {
	switch value := key.(type) {
	case int:
		return uint32(value & math.MaxUint32)
	case int8:
		return uint32(value) & math.MaxUint32
	case int16:
		return uint32(value) & math.MaxUint32
	case int32:
		return uint32(value) & math.MaxUint32
	case int64:
		return uint32(value) & math.MaxUint32
	case uint:
		return uint32(value) & math.MaxUint32
	case uint8:
		return uint32(value) & math.MaxUint32
	case uint16:
		return uint32(value) & math.MaxUint32
	case uint32:
		return uint32(value) & math.MaxUint32
	case uint64:
		return uint32(value) & math.MaxUint32
	case string:
		return Uint32Hash(string2Bytes(value))
	}
	return 0
}
