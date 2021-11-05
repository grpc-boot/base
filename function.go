package base

import (
	"fmt"
	"hash/crc32"
	"math"
)

func Index(key interface{}) uint8 {
	switch key.(type) {
	//优先使用自定义hash
	case CanHash:
		return uint8(key.(CanHash).HashCode() & math.MaxUint8)
	case string:
		return uint8(crc32.ChecksumIEEE([]byte(key.(string))) & math.MaxUint8)
	case []byte:
		return uint8(crc32.ChecksumIEEE(key.([]byte)) & math.MaxUint8)
	case uint8:
		return key.(uint8) & math.MaxUint8
	case uint16:
		return uint8(key.(uint16) & math.MaxUint8)
	case uint32:
		return uint8(key.(uint32) & math.MaxUint8)
	case uint64:
		return uint8(key.(uint64) & math.MaxUint8)
	case uint:
		return uint8(key.(uint) & math.MaxUint8)
	case int8:
		return uint8(key.(int8)) & math.MaxUint8
	case int16:
		return uint8(uint16(key.(int16)) & math.MaxUint8)
	case int32:
		return uint8(uint32(key.(int32)) & math.MaxUint8)
	case int64:
		return uint8(uint64(key.(int64)) & math.MaxUint8)
	case int:
		return uint8(uint(key.(int)) & math.MaxUint8)
	case float64:
		return uint8(int64(key.(float64)) & math.MaxUint8)
	case float32:
		return uint8(int64(key.(float32)) & math.MaxUint8)
	}

	return uint8(crc32.ChecksumIEEE([]byte(fmt.Sprintln(key))) & math.MaxUint8)
}
