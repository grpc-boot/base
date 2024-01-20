package elasticsearch

import (
	"strings"
	"sync"

	"github.com/grpc-boot/base/v2/orm/condition"
)

var (
	queryPool = &sync.Pool{
		New: func() any {
			return &Query{}
		},
	}
)

// AcquireQuery 获取mysqlQuery对象
func AcquireQuery() *Query {
	return queryPool.Get().(*Query)
}

type Query struct {
	table   string
	columns string
	where   condition.Condition
	group   string
	having  string
	order   string
	limit   int64
}

func (q *Query) reset() *Query {
	q.table = ""
	q.columns = ""
	q.limit = 0
	q.group = ""
	q.having = ""
	q.order = ""
	q.where = nil

	return q
}

func (q *Query) Select(columns ...string) *Query {
	q.columns = strings.Join(columns, ",")
	return q
}

func (q *Query) From(table string) *Query {
	q.table = table
	return q
}

func (q *Query) HasFrom() bool {
	return q.table != ""
}

func (q *Query) Where(condition condition.Condition) *Query {
	q.where = condition
	return q
}

func (q *Query) Group(fields ...string) *Query {
	q.group = " GROUP BY " + strings.Join(fields, ",")
	return q
}

func (q *Query) Having(having string) *Query {
	q.having = " HAVING " + having
	return q
}

func (q *Query) Order(orders ...string) *Query {
	q.order = " ORDER BY " + strings.Join(orders, ",")
	return q
}

func (q *Query) Limit(limit int64) *Query {
	q.limit = limit
	return q
}

func (q *Query) Close() {
	queryPool.Put(q.reset())
}

func (q *Query) Sql() (sql string, args []any) {
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

func (q *Query) Count(field string) (sql string, args []any) {
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

func (q *Query) Sum(field string) (sql string, args []any) {
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

func (q *Query) Max(field string) (sql string, args []any) {
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

func (q *Query) Min(field string) (sql string, args []any) {
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

func (q *Query) Avg(field string) (sql string, args []any) {
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
