package basis

import (
	"strconv"
	"strings"

	"github.com/grpc-boot/base/v2/utils"
)

const (
	emptySetter = ``
)

type Setter interface {
	Build() (string, error)
}

type StringSetter string

func (ss StringSetter) Build() (string, error) {
	if len(ss) == 0 {
		return emptySetter, ErrNoSetter
	}

	return string(ss), nil
}

type MapSetter map[string]any

func (ms MapSetter) Build() (string, error) {
	if len(ms) == 0 {
		return emptySetter, ErrNoSetter
	}

	var buf strings.Builder
	for field, value := range ms {
		if buf.Len() > 0 {
			buf.WriteString(",")
		}

		buf.WriteString(field)
		buf.WriteByte('=')
		switch v := value.(type) {
		case string:
			buf.WriteByte('\'')
			buf.WriteString(v)
			buf.WriteByte('\'')
		case []byte:
			buf.WriteByte('\'')
			buf.WriteString(utils.Bytes2String(v))
			buf.WriteByte('\'')
		case int:
			buf.WriteString(strconv.Itoa(v))
		case int8:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int16:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int32:
			buf.WriteString(strconv.FormatInt(int64(v), 10))
		case int64:
			buf.WriteString(strconv.FormatInt(v, 10))
		case uint:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint8:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint16:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint32:
			buf.WriteString(strconv.FormatUint(uint64(v), 10))
		case uint64:
			buf.WriteString(strconv.FormatUint(v, 10))
		default:
			return emptySetter, ErrDataFormat
		}
	}

	return buf.String(), nil
}
