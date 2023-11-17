package condition

import "strings"

type Gte struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (g Gte) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(g.Field) + 3)

	buffer.WriteString(g.Field)
	buffer.WriteString(">=?")

	return buffer.String(), []interface{}{g.Value}
}
