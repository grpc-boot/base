package condition

import "strings"

type Lt struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

func (l Lt) Build() (sql string, args []interface{}) {
	var (
		buffer strings.Builder
	)

	buffer.Grow(len(l.Field) + 2)

	buffer.WriteString(l.Field)
	buffer.WriteString("<?")

	return buffer.String(), []interface{}{l.Value}
}
