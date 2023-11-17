package condition

import "strings"

type Or []Condition

func (oc Or) Build() (sql string, args []interface{}) {
	if len(oc) == 0 {
		return
	}

	var buffer strings.Builder

	args = make([]interface{}, 0)

	buffer.WriteByte('(')
	sql, buildArgs := oc[0].Build()
	buffer.WriteString(sql)
	args = append(args, buildArgs...)
	for index := 1; index < len(oc); index++ {
		sql, buildArgs = oc[index].Build()
		buffer.WriteString(" OR ")
		buffer.WriteString(sql)
		args = append(args, buildArgs...)

	}
	buffer.WriteByte(')')
	return buffer.String(), args
}
