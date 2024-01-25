package condition

import "strings"

type Contains struct {
	Field string `json:"field"`
	Words string `json:"words"`
}

func (c Contains) Build() (sql string, args []any) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(c.Field) + 7)

	buffer.WriteString(c.Field)
	buffer.WriteString(" LIKE ?")

	return buffer.String(), []any{"%" + c.Words + "%"}
}
