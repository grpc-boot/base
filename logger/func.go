package logger

import "unsafe"

func bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}
