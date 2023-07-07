package condition

import "strings"

type NotEqual struct {
	Field string
	Value interface{}
}

func (ne NotEqual) Build() (sql string, args []interface{}) {
	var buffer strings.Builder

	buffer.Grow(len(ne.Field) + 3)

	buffer.WriteString(ne.Field)
	buffer.WriteString("<>?")

	return buffer.String(), []interface{}{ne.Value}
}
