package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"unicode"
	"unsafe"
)

// ToString 转为字符串类型
func ToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return Bytes2String(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Bytes2String 字节切片转换为字符串
func Bytes2String(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

// Bytes2Int64 字节切片转换为int64
func Bytes2Int64(data []byte) int64 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return val
}

// Bytes2Uint32 字节切片转换为uint32
func Bytes2Uint32(data []byte) uint32 {
	val, _ := strconv.ParseInt(*(*string)(unsafe.Pointer(&data)), 10, 64)
	return uint32(val)
}

// Bytes2Float64 字节切片转换为float64
func Bytes2Float64(data []byte) float64 {
	val, _ := strconv.ParseFloat(Bytes2String(data), 64)
	return val
}

// Int64ToHex _
func Int64ToHex(value int64) string {
	return strconv.FormatInt(value, 16)
}

// Uint64ToHex _
func Uint64ToHex(value uint64) string {
	return strconv.FormatUint(value, 16)
}

// Int64ToHexWithPad _
func Int64ToHexWithPad(value int64, padLength int) string {
	hexStr := Int64ToHex(value)
	if len(hexStr) >= padLength {
		return hexStr
	}

	buffer := make([]byte, padLength, padLength)

	for i := 0; i < (padLength - len(hexStr)); i++ {
		buffer[i] = '0'
	}

	start := padLength - len(hexStr)
	for j := start; j < padLength; j++ {
		buffer[j] = hexStr[j-start]
	}

	return Bytes2String(buffer)
}

func PackUin32(value uint32) []byte {
	buffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(buffer, binary.BigEndian, value)
	return buffer.Bytes()
}

func PackIn32(value int32) []byte {
	buffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(buffer, binary.BigEndian, value)
	return buffer.Bytes()
}

func UnpackUint32(data []byte) (value uint32, err error) {
	buffer := bytes.NewBuffer(data)
	var val uint32
	err = binary.Read(buffer, binary.BigEndian, &val)
	return val, err
}

func UnpackInt32(data []byte) (value int32, err error) {
	buffer := bytes.NewBuffer(data)
	var val int32
	err = binary.Read(buffer, binary.BigEndian, &val)
	return val, err
}

// Hex2Int64 _
func Hex2Int64(value string) int64 {
	val, _ := strconv.ParseInt(value, 16, 64)
	return val
}

// Hex2Uint64 _
func Hex2Uint64(value string) uint64 {
	val, _ := strconv.ParseUint(value, 16, 64)
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
