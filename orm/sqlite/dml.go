package sqlite

import (
	"strings"

	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/condition"
)

func repeatAndJoin(word, sep string, count int) string {
	if count < 1 {
		return ""
	}

	if count == 1 {
		return word
	}

	var buffer strings.Builder
	buffer.Grow(len(word) + (count-1)*(len(word)+len(sep)))

	buffer.WriteString(word)
	for index := 1; index < count; index++ {
		buffer.WriteString(sep)
		buffer.WriteString(word)
	}

	return buffer.String()
}

func Replace(table string, columns basis.Columns, rows []basis.Row) (sql string, args []any) {
	if len(columns) == 0 || len(rows) == 0 {
		return
	}

	var (
		buffer      strings.Builder
		fields      = strings.Join(columns, ",")
		placeHolder = repeatAndJoin("?", ",", len(columns))
	)
	length := 7 + 6 + len(table) + len(fields) + 2 + 7 + 2 + len(placeHolder) + (len(rows)-1)*(len(placeHolder)+3)

	buffer.Grow(length)

	args = make([]any, 0, len(columns)*len(rows))

	buffer.WriteString("REPLACE INTO ")
	buffer.WriteString(table)
	buffer.WriteByte('(')
	buffer.WriteString(fields)
	buffer.WriteByte(')')
	buffer.WriteString(" VALUES")
	buffer.WriteByte('(')
	buffer.WriteString(placeHolder)
	buffer.WriteByte(')')

	args = append(args, rows[0]...)
	for index := 1; index < len(rows); index++ {
		buffer.WriteString(",(")
		buffer.WriteString(placeHolder)
		buffer.WriteByte(')')
		args = append(args, rows[index]...)
	}

	return buffer.String(), args
}

func Insert(table string, columns basis.Columns, rows []basis.Row, ignore bool) (sql string, args []any) {
	if len(columns) == 0 || len(rows) == 0 {
		return
	}

	var (
		buffer      strings.Builder
		fields      = strings.Join(columns, ",")
		placeHolder = repeatAndJoin("?", ",", len(columns))
	)
	length := 6 + 6 + len(table) + len(fields) + 2 + 7 + 2 + len(placeHolder) + (len(rows)-1)*(len(placeHolder)+3)
	if ignore {
		length += 7
	}

	buffer.Grow(length)

	args = make([]any, 0, len(columns)*len(rows))

	buffer.WriteString("INSERT")
	if ignore {
		buffer.WriteString(" IGNORE")
	}

	buffer.WriteString(" INTO ")
	buffer.WriteString(table)
	buffer.WriteByte('(')
	buffer.WriteString(fields)
	buffer.WriteByte(')')
	buffer.WriteString(" VALUES")
	buffer.WriteByte('(')
	buffer.WriteString(placeHolder)
	buffer.WriteByte(')')

	args = append(args, rows[0]...)
	for index := 1; index < len(rows); index++ {
		buffer.WriteString(",(")
		buffer.WriteString(placeHolder)
		buffer.WriteByte(')')
		args = append(args, rows[index]...)
	}

	return buffer.String(), args
}

func Update(table, setters string, condition condition.Condition) (sql string, args []any) {
	var (
		where  string
		buffer strings.Builder
	)

	where, args = condition.Build()

	length := 7 + len(table) + 5 + len(setters) + len(where)
	if where != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString("UPDATE ")
	buffer.WriteString(table)
	buffer.WriteString(" SET ")
	buffer.WriteString(setters)
	if where != "" {
		buffer.WriteString(" WHERE ")
		buffer.WriteString(where)
	}

	return buffer.String(), args
}

func Delete(table string, condition condition.Condition) (sql string, args []any) {
	var (
		where  string
		buffer strings.Builder
	)

	where, args = condition.Build()

	length := 12 + len(table) + len(where)
	if where != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString("DELETE FROM ")
	buffer.WriteString(table)
	if where != "" {
		buffer.WriteString(" WHERE ")
		buffer.WriteString(where)
	}

	return buffer.String(), args
}
