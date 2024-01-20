package basis

import (
	"github.com/grpc-boot/base/v2/orm/condition"
)

// Query Query对象
type Query interface {
	// Select Select表达式
	Select(columns ...string) Query
	// From From表达式
	From(table string) Query
	// 是否设置了table
	HasFrom() bool
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
	Sql() (sql string, args []any)
	// Count 计数
	Count(field string) (sql string, args []any)
	// Sum 求和
	Sum(field string) (sql string, args []any)
	// Max 最大值
	Max(field string) (sql string, args []any)
	// Min 最小值
	Min(field string) (sql string, args []any)
	// Avg 平均值
	Avg(field string) (sql string, args []any)
	// Close 释放Query
	Close()
}
