package sqlite

import (
	"strconv"
	"strings"
	"sync"

	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/condition"
)

var (
	queryPool = &sync.Pool{
		New: func() any {
			return &query{}
		},
	}
)

// AcquireQuery 获取mysqlQuery对象
func AcquireQuery() basis.Query {
	return queryPool.Get().(basis.Query)
}

type query struct {
	table   string
	columns string
	where   condition.Condition
	group   string
	having  string
	order   string
	offset  int64
	limit   int64
}

func (q *query) reset() basis.Query {
	q.table = ""
	q.columns = ""
	q.offset = 0
	q.limit = 0
	q.group = ""
	q.having = ""
	q.order = ""
	q.where = nil

	return q
}

func (q *query) Select(columns ...string) basis.Query {
	q.columns = strings.Join(columns, ",")
	return q
}

func (q *query) From(table string) basis.Query {
	q.table = table
	return q
}

func (q *query) Where(condition condition.Condition) basis.Query {
	q.where = condition
	return q
}

func (q *query) Group(fields ...string) basis.Query {
	q.group = " GROUP BY " + strings.Join(fields, ",")
	return q
}

func (q *query) Having(having string) basis.Query {
	q.having = " HAVING " + having
	return q
}

func (q *query) Order(orders ...string) basis.Query {
	q.order = " ORDER BY " + strings.Join(orders, ",")
	return q
}

func (q *query) Offset(offset int64) basis.Query {
	q.offset = offset
	return q
}

func (q *query) Limit(limit int64) basis.Query {
	q.limit = limit
	return q
}

func (q *query) Close() {
	queryPool.Put(q.reset())
}

func (q *query) Sql() (sql string, args []any) {
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

	if q.limit < 1 {
		return sqlBuffer.String(), args
	}

	sqlBuffer.WriteString(" LIMIT ")
	sqlBuffer.WriteString(strconv.FormatInt(q.offset, 10))
	sqlBuffer.WriteByte(',')
	sqlBuffer.WriteString(strconv.FormatInt(q.limit, 10))
	return sqlBuffer.String(), args
}

func (q *query) Count(field string) (sql string, args []any) {
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

func (q *query) Sum(field string) (sql string, args []any) {
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

func (q *query) Max(field string) (sql string, args []any) {
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

func (q *query) Min(field string) (sql string, args []any) {
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

func (q *query) Avg(field string) (sql string, args []any) {
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
