package components

import (
	"fmt"
	"hash/crc32"
	"math"
)

// CanHash hash接口
type CanHash interface {
	// HashCode 计算hash值
	HashCode() (hashValue uint32)
}

// HashValue 计算任意类型hash值
func HashValue(key interface{}) uint32 {
	switch v := key.(type) {
	case CanHash:
		return v.HashCode()
	case bool:
		if v {
			return 1
		}
		return 0
	case uint8:
		return uint32(v)
	case uint16:
		return uint32(v)
	case uint32:
		return key.(uint32)
	case uint64:
		return uint32(v & math.MaxUint32)
	case uint:
		return uint32(v & math.MaxUint32)
	case int8:
		return uint32(v)
	case int16:
		return uint32(v)
	case int32:
		return uint32(v)
	case int64:
		return uint32(v & math.MaxUint32)
	case int:
		return uint32(key.(int) & math.MaxUint32)
	case float64:
		return uint32(int64(v) & math.MaxUint32)
	case float32:
		return uint32(v)
	case string:
		return crc32.ChecksumIEEE([]byte(v))
	case []byte:
		return crc32.ChecksumIEEE(v)
	}

	return crc32.ChecksumIEEE([]byte(fmt.Sprint(key)))
}

// Index4Bit 索引路由方法，值范围为uint32
func Index4Bit(key interface{}, bitCount uint8) uint32 {
	return HashValue(key) & ((1 << bitCount) - 1)
}

// Index4Uint8 索引路由方法，值范围为uint8
func Index4Uint8(key interface{}) uint8 {
	return uint8(HashValue(key) & math.MaxUint8)
}

// Index4Int8 索引路由方法，值范围为int8
func Index4Int8(key interface{}) int8 {
	return int8(HashValue(key) & math.MaxInt8)
}

// Index4Int16 索引路由方法，值范围为int16
func Index4Int16(key interface{}) int16 {
	return int16(HashValue(key) * math.MaxInt16)
}
