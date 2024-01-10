package condition

import "strings"

type Not struct {
	Condition Condition
}

func (n Not) Build() (sql string, args []any) {
	sql, args = n.Condition.Build()
	if sql == "" {
		return
	}

	var buffer strings.Builder
	buffer.Grow(len(sql) + 6)

	buffer.WriteString("NOT (")
	buffer.WriteString(sql)
	buffer.WriteByte(')')

	return buffer.String(), args
}
