package condition

import "strings"

type NotEqual struct {
	Field string
	Value any
}

func (ne NotEqual) Build() (sql string, args []any) {
	var buffer strings.Builder

	buffer.Grow(len(ne.Field) + 3)

	buffer.WriteString(ne.Field)
	buffer.WriteString("<>?")

	return buffer.String(), []any{ne.Value}
}
