package utils

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/grpc-boot/base/v2/kind"
)

// ToString 转为字符串类型
func ToString(val any) string {
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
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func StringInteger[T kind.Integer](value T) string {
	return fmt.Sprintf("%d", value)
}

func LcFirst(str string) string {
	strBytes := []byte(str)
	strBytes[0] = bytes.ToLower(strBytes[:1])[0]
	return Bytes2String(strBytes)
}

func UcFirst(str string) string {
	strBytes := []byte(str)
	strBytes[0] = bytes.ToUpper(strBytes[:1])[0]
	return Bytes2String(strBytes)
}
