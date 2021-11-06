package base

import (
	"strconv"
	"unsafe"
)

func Bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

func Bytes2Int64(data []byte) int64 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return val
}
