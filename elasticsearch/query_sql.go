package elasticsearch

import (
	"strings"
	"sync"

	"github.com/grpc-boot/base/v2/orm/condition"
)

var (
	queryPool = &sync.Pool{
		New: func() any {
			return &QuerySql{}
		},
	}
)

// AcquireQuery 获取mysqlQuery对象
func AcquireQuery() *QuerySql {
	return queryPool.Get().(*QuerySql)
}

type QuerySql struct {
	table   string
	columns string
	where   condition.Condition
	group   string
	having  string
	order   string
	limit   int64
}

func (q *QuerySql) reset() *QuerySql {
	q.table = ""
	q.columns = ""
	q.limit = 0
	q.group = ""
	q.having = ""
	q.order = ""
	q.where = nil

	return q
}

func (q *QuerySql) Select(columns ...string) *QuerySql {
	q.columns = strings.Join(columns, ",")
	return q
}

func (q *QuerySql) From(table string) *QuerySql {
	q.table = table
	return q
}

func (q *QuerySql) HasFrom() bool {
	return q.table != ""
}

func (q *QuerySql) Where(condition condition.Condition) *QuerySql {
	q.where = condition
	return q
}

func (q *QuerySql) Group(fields ...string) *QuerySql {
	q.group = " GROUP BY " + strings.Join(fields, ",")
	return q
}

func (q *QuerySql) Having(having string) *QuerySql {
	q.having = " HAVING " + having
	return q
}

func (q *QuerySql) Order(orders ...string) *QuerySql {
	q.order = " ORDER BY " + strings.Join(orders, ",")
	return q
}

func (q *QuerySql) Limit(limit int64) *QuerySql {
	q.limit = limit
	return q
}

func (q *QuerySql) Close() {
	queryPool.Put(q.reset())
}

func (q *QuerySql) Sql() (sql string, args []any) {
	var (
		whereStr  string
		sqlBuffer strings.Builder
	)

	sqlBuffer.WriteString(`SELECT `)

	if q.columns == "" {
		sqlBuffer.WriteString("*")
	} else {
		sqlBuffer.WriteString(q.columns)
	}

	sqlBuffer.WriteString(` FROM `)
	sqlBuffer.WriteString(q.table)

	if q.where != nil {
		whereStr, args = q.where.Build()
		if whereStr != "" {
			sqlBuffer.WriteString(" WHERE ")
			sqlBuffer.WriteString(whereStr)
		}
	}

	sqlBuffer.WriteString(q.group)
	sqlBuffer.WriteString(q.having)
	sqlBuffer.WriteString(q.order)

	return sqlBuffer.String(), args
}

func (q *QuerySql) Count(field string) (sql string, args []any) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if q.where != nil {
		whereStr, args = q.where.Build()
	}

	length := 13 + len(field) + 7 + len(q.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT COUNT(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(q.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (q *QuerySql) Sum(field string) (sql string, args []any) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if q.where != nil {
		whereStr, args = q.where.Build()
	}

	length := 11 + len(field) + 7 + len(q.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT SUM(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(q.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (q *QuerySql) Max(field string) (sql string, args []any) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if q.where != nil {
		whereStr, args = q.where.Build()
	}

	length := 11 + len(field) + 7 + len(q.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT Max(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(q.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (q *QuerySql) Min(field string) (sql string, args []any) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if q.where != nil {
		whereStr, args = q.where.Build()
	}

	length := 11 + len(field) + 7 + len(q.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT Min(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(q.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (q *QuerySql) Avg(field string) (sql string, args []any) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if q.where != nil {
		whereStr, args = q.where.Build()
	}

	length := 11 + len(field) + 7 + len(q.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT Avg(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(q.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}
