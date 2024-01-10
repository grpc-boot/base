package base

var ModelTemplate = `package {package}

import (
	"context"
	"database/sql"

	"github.com/grpc-boot/base/v2/orm"
	"github.com/grpc-boot/base/v2/orm/base"
	"github.com/grpc-boot/base/v2/orm/condition"
	"github.com/grpc-boot/base/v2/orm/mysql"
)

type {struct} struct {{fields}
}

func ({this} *{struct}) TableName() string {
	return "{table}"
}

func ({this} *{struct}) PrimaryField() string {
	return "{primary}"
}

func ({this} *{struct}) PrimaryValue() any {
	return {this}.{Primary}
}

func ({this} *{struct}) GetDefault(name string) string {
	return orm.ParseMapping({this}).Default(name)
}

func ({this} *{struct}) GetLabel(name string) string {
	return orm.ParseMapping({this}).Label(name)
}

func ({this} *{struct}) GetEnums(name string) []string {
	return orm.ParseMapping({this}).Enums(name)
}

func ({this} *{struct}) GetSize(name string) int {
	return orm.ParseMapping({this}).Size(name)
}

func ({this} *{struct}) Query() base.Query {
	return mysql.AcquireQuery()
}

func ({this} *{struct}) Add(ctx context.Context, executor orm.Executor, {smallCamel} {struct}) (res sql.Result, err error) {
	fm := orm.ParseMapping(&{smallCamel})
	sqlStr, args := orm.Insert({smallCamel}.TableName(), orm.Columns{{columns}
	}, []orm.Row{{{rows}
	}}, false)

	return orm.Exec(ctx, executor, sqlStr, args...)
}

func ({this} *{struct}) FindOne(ctx context.Context, executor orm.Executor, {smallCamel} {struct}) ({struct}, error) {
	qObj := {this}.Query().
		From({smallCamel}.TableName()).
		Where(condition.Equal{
			Field: "{primary}",
			Value: {smallCamel}.{Primary},
		})
	defer qObj.Close()

	list, err := {this}.FindByQuery(ctx, executor, qObj)
	if err != nil || len(list) == 0 {
		return {struct}{}, err
	}
	
	return list[0], nil
}

func ({this} *{struct}) FindSome(ctx context.Context, executor orm.Executor, idList ...any) ([]{struct}, error) {
	qObj := {this}.Query().
		From({this}.TableName()).
		Where(condition.In[any]{
			Field: "id",
			Value: idList,
		})
	defer qObj.Close()

	return {this}.FindByQuery(ctx, executor, qObj)
}

func ({this} *{struct}) FindByQuery(ctx context.Context, executor orm.Executor, query base.Query) ([]{struct}, error) {
	rows, err := orm.QueryWithQuery(ctx, executor, query)
	if err != nil {
		return nil, err
	}

	list, err := {this}.scanRows(rows)
	if err != nil || len(list) == 0 {
		return nil, err
	}

	return list, err
}

func ({this} *{struct}) scanRows(rows *sql.Rows) ([]{struct}, error) {
	records, err := orm.ScanRecords(rows)
	if err != nil {
		return nil, err
	}

	fm := orm.ParseMapping({this})
	list := make([]{struct}, len(records))
	for index, row := range records {
		list[index] = {struct}{{assignment}
		}
	}

	return list, nil
}
`
