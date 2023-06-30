package query

import (
	"fmt"
	"github.com/grpc-boot/base/internal"
	"strconv"
)

type Item interface {
	Build() string
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case string:
		//buf := strings.Builder{}

		return "'" + v + "'"
	case []byte:
		return "'" + internal.Bytes2String(v) + "'"
	}

	return fmt.Sprint(value)
}

type Equal struct {
	Field string
	Value interface{}
}

func (e *Equal) Build() string {

}

type In struct {
	Name  string        `json:"name"`
	Value []interface{} `json:"value"`
}

type Between struct {
}
