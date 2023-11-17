package condition

import "strings"

type Lte struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (l Lte) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(l.Field) + 3)

	buffer.WriteString(l.Field)
	buffer.WriteString("<=?")

	return buffer.String(), []interface{}{l.Value}
}
