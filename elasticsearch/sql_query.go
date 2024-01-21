package elasticsearch

import (
	"strings"
	"sync"

	"github.com/grpc-boot/base/v2/orm/condition"
)

var (
	queryPool = &sync.Pool{
		New: func() any {
			return &SqlQuery{}
		},
	}
)

// AcquireQuery 获取mysqlQuery对象
func AcquireQuery() *SqlQuery {
	return queryPool.Get().(*SqlQuery)
}

type SqlQuery struct {
	table   string
	columns string
	where   condition.Condition
	group   string
	having  string
	order   string
	limit   int64
}

func (q *SqlQuery) reset() *SqlQuery {
	q.table = ""
	q.columns = ""
	q.limit = 0
	q.group = ""
	q.having = ""
	q.order = ""
	q.where = nil

	return q
}

func (q *SqlQuery) Select(columns ...string) *SqlQuery {
	q.columns = strings.Join(columns, ",")
	return q
}

func (q *SqlQuery) From(table string) *SqlQuery {
	q.table = table
	return q
}

func (q *SqlQuery) HasFrom() bool {
	return q.table != ""
}

func (q *SqlQuery) Where(condition condition.Condition) *SqlQuery {
	q.where = condition
	return q
}

func (q *SqlQuery) Group(fields ...string) *SqlQuery {
	q.group = " GROUP BY " + strings.Join(fields, ",")
	return q
}

func (q *SqlQuery) Having(having string) *SqlQuery {
	q.having = " HAVING " + having
	return q
}

func (q *SqlQuery) Order(orders ...string) *SqlQuery {
	q.order = " ORDER BY " + strings.Join(orders, ",")
	return q
}

func (q *SqlQuery) Limit(limit int64) *SqlQuery {
	q.limit = limit
	return q
}

func (q *SqlQuery) Close() {
	queryPool.Put(q.reset())
}

func (q *SqlQuery) Sql() (sql string, args []any) {
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

func (q *SqlQuery) Count(field string) (sql string, args []any) {
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

func (q *SqlQuery) Sum(field string) (sql string, args []any) {
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

func (q *SqlQuery) Max(field string) (sql string, args []any) {
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

func (q *SqlQuery) Min(field string) (sql string, args []any) {
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

func (q *SqlQuery) Avg(field string) (sql string, args []any) {
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
