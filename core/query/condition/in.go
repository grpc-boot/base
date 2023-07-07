package condition

import (
	"strings"
)

type In struct {
	Field string `json:"field"`
	Value Values `json:"value"`
}

func (i In) Build() (sql string, args []interface{}) {
	if len(i.Value) == 0 {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(i.Field) + 6 + (len(i.Value)-1)*2)
	args = make([]interface{}, len(i.Value), len(i.Value))

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
