package query

import (
	"strconv"
	"strings"
	"sync"

	"github.com/grpc-boot/base/core/query/condition"
)

var (
	mysqlQueryPool = &sync.Pool{
		New: func() interface{} {
			return &mysqlQuery{}
		},
	}
)

// Query Query对象
type Query interface {
	// Select Select表达式
	Select(columns ...string) Query
	// From From表达式
	From(table string) Query
	// Where where表达式
	Where(condition condition.Condition) Query
	// Group Group表达式
	Group(fields ...string) Query
	// Having Having表达式
	Having(having string) Query
	// Order Order表达式
	Order(orders ...string) Query
	// Offset Offset表达式
	Offset(offset int64) Query
	// Limit Limit表达式
	Limit(limit int64) Query
	// Sql 生成sql和参数
	Sql() (sql string, args []interface{})
	// Count 计数
	Count(field string) (sql string, args []interface{})
	// Sum 求和
	Sum(field string) (sql string, args []interface{})
	// Max 最大值
	Max(field string) (sql string, args []interface{})
	// Min 最小值
	Min(field string) (sql string, args []interface{})
	// Avg 平均值
	Avg(field string) (sql string, args []interface{})
	// Close 释放Query
	Close()
}

// Acquire4Mysql 获取mysqlQuery对象
func Acquire4Mysql() Query {
	return mysqlQueryPool.Get().(Query)
}

type mysqlQuery struct {
	table   string
	columns string
	where   condition.Condition
	group   string
	having  string
	order   string
	offset  int64
	limit   int64
}

func (mq *mysqlQuery) reset() Query {
	mq.table = ""
	mq.columns = ""
	mq.offset = 0
	mq.limit = 0
	mq.group = ""
	mq.having = ""
	mq.order = ""
	mq.where = nil

	return mq
}

func (mq *mysqlQuery) Select(columns ...string) Query {
	mq.columns = strings.Join(columns, ",")
	return mq
}

func (mq *mysqlQuery) From(table string) Query {
	mq.table = table
	return mq
}

func (mq *mysqlQuery) Where(condition condition.Condition) Query {
	mq.where = condition
	return mq
}

func (mq *mysqlQuery) Group(fields ...string) Query {
	mq.group = " GROUP BY " + strings.Join(fields, ",")
	return mq
}

func (mq *mysqlQuery) Having(having string) Query {
	mq.having = " HAVING " + having
	return mq
}

func (mq *mysqlQuery) Order(orders ...string) Query {
	mq.order = " ORDER BY " + strings.Join(orders, ",")
	return mq
}

func (mq *mysqlQuery) Offset(offset int64) Query {
	mq.offset = offset
	return mq
}

func (mq *mysqlQuery) Limit(limit int64) Query {
	mq.limit = limit
	return mq
}

func (mq *mysqlQuery) Close() {
	mq.reset()
	mysqlQueryPool.Put(mq)
}

func (mq *mysqlQuery) Sql() (sql string, args []interface{}) {
	var (
		whereStr  string
		sqlBuffer strings.Builder
	)

	sqlBuffer.WriteString(`SELECT `)

	if mq.columns == "" {
		sqlBuffer.WriteString("*")
	} else {
		sqlBuffer.WriteString(mq.columns)
	}

	sqlBuffer.WriteString(` FROM `)
	sqlBuffer.WriteString(mq.table)

	if mq.where != nil {
		whereStr, args = mq.where.Build()
		if whereStr != "" {
			sqlBuffer.WriteString(" WHERE ")
			sqlBuffer.WriteString(whereStr)
		}
	}

	sqlBuffer.WriteString(mq.group)
	sqlBuffer.WriteString(mq.having)
	sqlBuffer.WriteString(mq.order)

	if mq.limit < 1 {
		return sqlBuffer.String(), args
	}

	sqlBuffer.WriteString(" LIMIT ")
	sqlBuffer.WriteString(strconv.FormatInt(mq.offset, 10))
	sqlBuffer.WriteByte(',')
	sqlBuffer.WriteString(strconv.FormatInt(mq.limit, 10))
	return sqlBuffer.String(), args
}

func (mq *mysqlQuery) Count(field string) (sql string, args []interface{}) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if mq.where != nil {
		whereStr, args = mq.where.Build()
	}

	length := 13 + len(field) + 7 + len(mq.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT COUNT(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(mq.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (mq *mysqlQuery) Sum(field string) (sql string, args []interface{}) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if mq.where != nil {
		whereStr, args = mq.where.Build()
	}

	length := 11 + len(field) + 7 + len(mq.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT SUM(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(mq.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (mq *mysqlQuery) Max(field string) (sql string, args []interface{}) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if mq.where != nil {
		whereStr, args = mq.where.Build()
	}

	length := 11 + len(field) + 7 + len(mq.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT Max(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(mq.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (mq *mysqlQuery) Min(field string) (sql string, args []interface{}) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if mq.where != nil {
		whereStr, args = mq.where.Build()
	}

	length := 11 + len(field) + 7 + len(mq.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT Min(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(mq.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}

func (mq *mysqlQuery) Avg(field string) (sql string, args []interface{}) {
	var (
		buffer   strings.Builder
		whereStr string
	)

	if mq.where != nil {
		whereStr, args = mq.where.Build()
	}

	length := 11 + len(field) + 7 + len(mq.table) + len(whereStr)
	if whereStr != "" {
		length += 7
	}

	buffer.Grow(length)

	buffer.WriteString(`SELECT Avg(`)
	buffer.WriteString(field)
	buffer.WriteString(`) FROM `)
	buffer.WriteString(mq.table)

	if whereStr != "" {
		buffer.WriteString(` WHERE `)
		buffer.WriteString(whereStr)
	}

	return buffer.String(), args
}
