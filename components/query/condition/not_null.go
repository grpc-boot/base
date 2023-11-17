package condition

import "strings"

type NotNull struct {
	Field string `json:"field"`
}

func (nn NotNull) Build() (sql string, args []interface{}) {
	var buffer strings.Builder
	buffer.Grow(len(nn.Field) + 12)

	buffer.WriteString(nn.Field)
	buffer.WriteString(" IS NOT NULL")

	return buffer.String(), args
}
