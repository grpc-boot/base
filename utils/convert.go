package utils

import (
	"strconv"
	"unsafe"

	"github.com/grpc-boot/base/v2/internal"
)

// Bytes2String 字节切片转换为字符串
func Bytes2String(data []byte) string {
	return internal.Bytes2String(data)
}

// String2Bytes 字符串转字节切片，注意：转换后不能对字节切片进行修改
func String2Bytes(data string) []byte {
	return internal.String2Bytes(data)
}

// Bytes2Int64 字节切片转换为int64
func Bytes2Int64(data []byte) int64 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return val
}

// Bytes2Uint32 字节切片转换为uint32
func Bytes2Uint32(data []byte) uint32 {
	return uint32(Bytes2Uint64(data))
}

// Bytes2Uint64 字节切片转换为uint64
func Bytes2Uint64(data []byte) uint64 {
	val, _ := strconv.ParseUint(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return val
}

// Bytes2Float64 字节切片转换为float64
func Bytes2Float64(data []byte) float64 {
	val, _ := strconv.ParseFloat(Bytes2String(data), 64)
	return val
}

// Bytes2Bool 字节切片转换为bool
func Bytes2Bool(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	val, _ := strconv.ParseBool(Bytes2String(data))
	return val
}
