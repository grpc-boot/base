package condition

import "strings"

type And []Condition

func (ac And) Build() (sql string, args []any) {
	if len(ac) == 0 {
		return
	}

	var buffer strings.Builder

	args = make([]any, 0)

	buffer.WriteByte('(')
	sql, buildArgs := ac[0].Build()
	buffer.WriteString(sql)
	args = append(args, buildArgs...)
	for index := 1; index < len(ac); index++ {
		sql, buildArgs = ac[index].Build()
		buffer.WriteString(" AND ")
		buffer.WriteString(sql)
		args = append(args, buildArgs...)

	}
	buffer.WriteByte(')')
	return buffer.String(), args
}
