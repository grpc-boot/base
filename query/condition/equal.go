package condition

import "strings"

type Equal struct {
	Field string
	Value any
}

func (e Equal) Build() (sql string, args []any) {
	var buffer strings.Builder

	buffer.Grow(len(e.Field) + 2)

	buffer.WriteString(e.Field)
	buffer.WriteString("=?")

	return buffer.String(), []any{e.Value}
}
