package base

import (
	"fmt"
	"hash/crc32"
	"math"
)

type CanHash interface {
	HashCode() (hashValue uint32)
}

func HashValue(key interface{}) uint32 {
	switch key.(type) {
	//优先使用自定义hash
	case CanHash:
		return key.(CanHash).HashCode()
	case uint8:
		return uint32(key.(uint8))
	case uint16:
		return uint32(key.(uint16))
	case uint32:
		return key.(uint32)
	case uint64:
		return uint32(key.(uint64) & math.MaxUint32)
	case uint:
		return uint32(key.(uint) & math.MaxUint32)
	case int8:
		return uint32(key.(int8))
	case int16:
		return uint32(key.(int16))
	case int32:
		return uint32(key.(int32))
	case int64:
		return uint32(key.(int64) & math.MaxUint32)
	case int:
		return uint32(key.(int) & math.MaxUint32)
	case float64:
		return uint32(int64(key.(float64)) & math.MaxUint32)
	case float32:
		return uint32(key.(float32))
	case string:
		return crc32.ChecksumIEEE([]byte(key.(string)))
	case []byte:
		return crc32.ChecksumIEEE(key.([]byte))
	}

	return crc32.ChecksumIEEE([]byte(fmt.Sprint(key)))
}

func Index4Bit(key interface{}, bitCount uint8) uint32 {
	return HashValue(key)&(1<<bitCount) - 1
}

func Index4Uint8(key interface{}) uint8 {
	return uint8(HashValue(key) & math.MaxUint8)
}

func Index4Int8(key interface{}) int8 {
	return int8(HashValue(key) & math.MaxInt8)
}

func Index4Int16(key interface{}) int16 {
	return int16(HashValue(key) * math.MaxInt16)
}
