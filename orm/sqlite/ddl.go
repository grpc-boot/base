package sqlite

import (
	"fmt"
	"strings"
)

func Pragma(key string, values ...any) (sql string, args []any) {
	if len(values) == 0 {
		return fmt.Sprintf("PRAGMA %s;", key), nil
	}

	return fmt.Sprintf("PRAGMA %s = %v;", key, values[0]), nil
}

func ShowIndexs(tableName string) (sql string, args []any) {
	return fmt.Sprintf("SELECT name, type FROM sqlite_master WHERE type='index' AND tbl_name='%s';", tableName), nil
}

func DropIndex(indexName string) (sql string, args []any) {
	return fmt.Sprintf("DROP INDEX %s;", indexName), nil
}

func AddIndex(indexName, tableName string, unique bool, cols ...string) (sql string, args []any) {
	if len(cols) < 1 {
		return
	}

	var buf strings.Builder

	buf.WriteString("CREATE ")

	if unique {
		buf.WriteString("UNIQUE ")
	}

	buf.WriteString("INDEX ")
	buf.WriteString(indexName)
	buf.WriteString(" ON ")
	buf.WriteString(tableName)
	buf.WriteByte('(')
	buf.WriteString(cols[0])

	for i := 1; i < len(cols); i++ {
		buf.WriteByte(',')
		buf.WriteString(cols[i])
	}
	buf.
		WriteByte(')')

	return buf.String(), nil
}
