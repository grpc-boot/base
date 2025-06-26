package utils

import (
	"github.com/grpc-boot/base/v3/internal"
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
	return internal.Bytes2Int64(data)
}

// Bytes2Uint64 字节切片转换为uint64
func Bytes2Uint64(data []byte) uint64 {
	return internal.Bytes2Uint64(data)
}

// Bytes2Int32 字节切片转换为int32
func Bytes2Int32(data []byte) int32 {
	return internal.Bytes2Int32(data)
}

// Bytes2Uint32 字节切片转换为uint32
func Bytes2Uint32(data []byte) uint32 {
	return internal.Bytes2Uint32(data)
}

// Bytes2Int16 字节切片转换为int16
func Bytes2Int16(data []byte) int16 {
	return internal.Bytes2Int16(data)
}

// Bytes2Uint16 字节切片转换为uint16
func Bytes2Uint16(data []byte) uint16 {
	return internal.Bytes2Uint16(data)
}

// Bytes2Int8 字节切片转换为int8
func Bytes2Int8(data []byte) int8 {
	return internal.Bytes2Int8(data)
}

// Bytes2Uint8 字节切片转换为uint8
func Bytes2Uint8(data []byte) uint8 {
	return internal.Bytes2Uint8(data)
}

// Bytes2Float64 字节切片转换为float64
func Bytes2Float64(data []byte) float64 {
	return internal.Bytes2Float64(data)
}

// Bytes2Bool 字节切片转换为bool
func Bytes2Bool(data []byte) bool {
	return internal.Bytes2Bool(data)
}
