package condition

import "strings"

type Null struct {
	Field string `json:"field"`
}

func (n Null) Build() (sql string, args []interface{}) {
	var buffer strings.Builder
	buffer.Grow(len(n.Field) + 8)

	buffer.WriteString(n.Field)
	buffer.WriteString(" IS NULL")

	return buffer.String(), args
}
