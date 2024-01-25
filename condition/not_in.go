package condition

import (
	"strings"

	"github.com/grpc-boot/base/v2/kind"
)

type NotIn[T any] struct {
	Field string        `json:"field"`
	Value kind.Slice[T] `json:"value"`
}

func (ni NotIn[T]) Build() (sql string, args []any) {
	if len(ni.Value) == 0 {
		return
	}

	if len(ni.Value) == 1 {
		return NotEqual{Field: ni.Field, Value: ni.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(ni.Field) + 10 + (len(ni.Value)-1)*2)
	args = make([]any, len(ni.Value), len(ni.Value))

	buffer.WriteString(ni.Field)
	buffer.WriteString(" NOT IN(?")
	args[0] = ni.Value[0]
	for index := 1; index < len(ni.Value); index++ {
		buffer.WriteString(",?")
		args[index] = ni.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}
