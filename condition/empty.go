package condition

import "strings"

type Empty struct {
	Field string `json:"field"`
}

func (e Empty) Build() (sql string, args []any) {
	var buffer strings.Builder
	buffer.Grow(len(e.Field) + 3)

	buffer.WriteString(e.Field)
	buffer.WriteString("=''")

	return buffer.String(), args
}
