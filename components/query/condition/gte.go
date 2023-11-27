package condition

import "strings"

type Gte struct {
	Field string `json:"field"`
	Value any    `json:"value"`
}

func (g Gte) Build() (sql string, args []any) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(g.Field) + 3)

	buffer.WriteString(g.Field)
	buffer.WriteString(">=?")

	return buffer.String(), []any{g.Value}
}
