package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/grpc-boot/base/v3/internal"
	"github.com/grpc-boot/base/v3/kind"
)

// ToString 转为字符串类型
func ToString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return internal.Bytes2String(v)
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
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	default:
		return strconv.FormatUint(val.Uint(), 10)
	}
}

func LcFirst(str string) string {
	return internal.LcFirst(str)
}

func UcFirst(str string) string {
	return internal.UcFirst(str)
}

func SmallCamelByChar(str string, char byte) string {
	if str == "" {
		return str
	}

	name := []byte(str)
	j := 0
	u := false
	for index, b := range name {
		if b == char {
			u = true
			continue
		}

		if u {
			name[j] = bytes.ToUpper(name[index : index+1])[0]
		} else {
			name[j] = b
		}
		j++
		u = false
	}

	return Bytes2String(name[:j])
}

func BigCamelByChar(str string, char byte) string {
	if str == "" {
		return str
	}

	name := []byte(str)
	name[0] = bytes.ToUpper(name[:1])[0]
	j := 0
	u := false
	for index, b := range name {
		if b == char {
			u = true
			continue
		}

		if u {
			name[j] = bytes.ToUpper(name[index : index+1])[0]
		} else {
			name[j] = b
		}
		j++
		u = false
	}

	return Bytes2String(name[:j])
}
