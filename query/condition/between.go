package condition

import "strings"

type Between struct {
	Field string `json:"field"`
	Start any    `json:"start"`
	End   any    `json:"end"`
}

func (b Between) Build() (sql string, args []any) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(b.Field) + 16)

	buffer.WriteString(b.Field)
	buffer.WriteString(" BETWEEN ? AND ?")

	return buffer.String(), []any{b.Start, b.End}
}
