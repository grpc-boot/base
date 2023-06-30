package query

import (
	"strconv"
	"strings"
	"sync"
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
	// Where Where表达式
	Where(where Where) Query
	// And 附加And where
	And(condition Condition) Query
	// Or 附加Or where
	Or(condition Condition) Query
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
	// Sql 生成sql
	Sql(arguments *[]interface{}) (sql string)
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
	where   []interface{}
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

	if mq.where != nil {
		mq.where = mq.where[:0]
	}

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

func (mq *mysqlQuery) Where(where Where) Query {
	mq.where = where
	return mq
}

func (mq *mysqlQuery) And(condition Condition) Query {
	mq.where.And(condition)
	return mq
}

func (mq *mysqlQuery) Or(condition Condition) Query {
	mq.where.Or(condition)
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

func (mq *mysqlQuery) Sql(arguments *[]interface{}) (sql string) {
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

	if mq.where != nil && mq.where.HasWhere() {
		whereStr = mq.where.Sql(arguments)
		sqlBuffer.WriteString(whereStr)
	}

	sqlBuffer.WriteString(mq.group)
	sqlBuffer.WriteString(mq.having)
	sqlBuffer.WriteString(mq.order)

	if mq.limit < 1 {
		return sqlBuffer.String()
	}

	sqlBuffer.WriteString(" LIMIT ")
	sqlBuffer.WriteString(strconv.FormatInt(mq.offset, 10))
	sqlBuffer.WriteString(",")
	sqlBuffer.WriteString(strconv.FormatInt(mq.limit, 10))
	return sqlBuffer.String()
}
