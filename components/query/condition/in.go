package condition

import (
	"strings"

	"github.com/grpc-boot/base/v2/kind"
)

type In[T any] struct {
	Field string        `json:"field"`
	Value kind.Slice[T] `json:"value"`
}

func (i In[T]) Build() (sql string, args []any) {
	if len(i.Value) == 0 {
		return
	}

	if len(i.Value) == 1 {
		return Equal{Field: i.Field, Value: i.Value[0]}.Build()
	}

	var buffer strings.Builder
	buffer.Grow(len(i.Field) + 6 + (len(i.Value)-1)*2)
	args = make([]any, len(i.Value), len(i.Value))

	buffer.WriteString(i.Field)
	buffer.WriteString(" IN(?")
	args[0] = i.Value[0]
	for index := 1; index < len(i.Value); index++ {
		buffer.WriteString(",?")
		args[index] = i.Value[index]
	}

	buffer.WriteByte(')')
	return buffer.String(), args
}
