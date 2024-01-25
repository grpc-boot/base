package condition

import "strings"

type Lte struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

func (l Lte) Build() (sql string, args []any) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(l.Field) + 3)

	buffer.WriteString(l.Field)
	buffer.WriteString("<=?")

	return buffer.String(), []any{l.Value}
}
