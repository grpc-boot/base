package internal

import "unsafe"

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
