package condition

import "strings"

type Equal struct {
	Field string
	Value interface{}
}

func (e Equal) Build() (sql string, args []interface{}) {
	var buffer strings.Builder

	buffer.Grow(len(e.Field) + 2)

	buffer.WriteString(e.Field)
	buffer.WriteString("=?")

	return buffer.String(), []interface{}{e.Value}
}
