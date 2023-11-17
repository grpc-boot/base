package condition

import "strings"

type Gt struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (g Gt) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(g.Field) + 2)

	buffer.WriteString(g.Field)
	buffer.WriteString(">?")

	return buffer.String(), []interface{}{g.Value}
}
