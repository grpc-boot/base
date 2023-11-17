package condition

import "strings"

type Between struct {
	Field string      `json:"field"`
	Start interface{} `json:"start"`
	End   interface{} `json:"end"`
}

func (b Between) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(b.Field) + 16)

	buffer.WriteString(b.Field)
	buffer.WriteString(" BETWEEN ? AND ?")

	return buffer.String(), []interface{}{b.Start, b.End}
}
