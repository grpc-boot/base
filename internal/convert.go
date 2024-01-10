package internal

import (
	"strconv"
	"unsafe"
)

// Bytes2String 字节切片转换为字符串
func Bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

// String2Bytes 字符串转字节切片，注意：转换后不能对字节切片进行修改
func String2Bytes(data string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&data))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2Int64 字节切片转换为int64
func Bytes2Int64(data []byte) int64 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return val
}

// Bytes2Uint64 字节切片转换为uint64
func Bytes2Uint64(data []byte) uint64 {
	val, _ := strconv.ParseUint(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return val
}

// Bytes2Int32 字节切片转换为int32
func Bytes2Int32(data []byte) int32 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 32)
	return int32(val)
}

// Bytes2Uint32 字节切片转换为uint32
func Bytes2Uint32(data []byte) uint32 {
	val, _ := strconv.ParseUint(*(*string)(unsafe.Pointer(&data)), 10, 32)
	return uint32(val)
}

// Bytes2Int16 字节切片转换为int16
func Bytes2Int16(data []byte) int16 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 16)
	return int16(val)
}

// Bytes2Uint16 字节切片转换为uint16
func Bytes2Uint16(data []byte) uint16 {
	val, _ := strconv.ParseUint(*(*string)(unsafe.Pointer(&data)), 10, 16)
	return uint16(val)
}

// Bytes2Int8 字节切片转换为int8
func Bytes2Int8(data []byte) int8 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 8)
	return int8(val)
}

// Bytes2Uint8 字节切片转换为uint8
func Bytes2Uint8(data []byte) uint8 {
	val, _ := strconv.ParseUint(*(*string)(unsafe.Pointer(&data)), 10, 8)
	return uint8(val)
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
