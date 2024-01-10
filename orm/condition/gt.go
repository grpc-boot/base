package condition

import "strings"

type Gt struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

func (g Gt) Build() (sql string, args []any) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(g.Field) + 2)

	buffer.WriteString(g.Field)
	buffer.WriteString(">?")

	return buffer.String(), []any{g.Value}
}
