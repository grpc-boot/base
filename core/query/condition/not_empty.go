package condition

import "strings"

type NotEmpty struct {
	Field string `json:"field"`
}

func (ne NotEmpty) Build() (sql string, args []interface{}) {
	var buffer strings.Builder
	buffer.Grow(len(ne.Field) + 4)

	buffer.WriteString(ne.Field)
	buffer.WriteString("<>''")

	return buffer.String(), args
}
