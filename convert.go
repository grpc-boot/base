package base

import (
	"strconv"
	"unicode"
	"unsafe"
)

// Bytes2String 字节切片转换为字符串
func Bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

// Bytes2Int64 字节切片转换为int64
func Bytes2Int64(data []byte) int64 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return val
}

// Bytes2Float64 字节切片转换为float64
func Bytes2Float64(data []byte) float64 {
	val, _ := strconv.ParseFloat(*(*string)(unsafe.Pointer(&data)), 64)
	return val
}

// BigCamels 转换为大驼峰
func BigCamels(sep byte, data string) string {
	var (
		fd    = []byte(data)
		upper = true
	)

	for index := 0; index < len(fd); index++ {
		if upper {
			fd[index] = byte(unicode.ToUpper(rune(fd[index])))
			upper = false
			continue
		}

		if fd[index] == sep {
			fd = append(fd[0:index], fd[index+1:]...)
			upper = true
			index--
		}
	}

	return Bytes2String(fd)
}

// SmallCamels 转换为小驼峰
func SmallCamels(sep byte, data string) string {
	var (
		fd    = []byte(data)
		upper = false
	)

	fd[0] = byte(unicode.ToLower(rune(fd[0])))

	for index := 1; index < len(fd); index++ {
		if upper {
			fd[index] = byte(unicode.ToUpper(rune(fd[index])))
			upper = false
			continue
		}

		if fd[index] == sep {
			fd = append(fd[0:index], fd[index+1:]...)
			upper = true
			index--
		}
	}

	return Bytes2String(fd)
}
