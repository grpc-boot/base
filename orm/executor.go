package orm

import (
	"context"
	"database/sql"
	"github.com/grpc-boot/base/v2/orm/condition"
	"time"

	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/orm/base"

	"go.uber.org/zap/zapcore"
)

type Executor interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func QueryWithQuery(ctx context.Context, executor Executor, q base.Query) (rows *sql.Rows, err error) {
	sqlStr, args := q.Sql()

	return Query(ctx, executor, sqlStr, args...)
}

func QueryWithQueryTimeout(timeout time.Duration, executor Executor, q base.Query) (rows *sql.Rows, err error) {
	sqlStr, args := q.Sql()

	return QueryTimeout(timeout, executor, sqlStr, args...)
}

func QueryRowTimeout(timeout time.Duration, executor Executor, query string, args ...any) (row *sql.Row, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return QueryRow(ctx, executor, query, args...)
}

func QueryRow(ctx context.Context, executor Executor, query string, args ...any) (row *sql.Row, err error) {
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapError("exec sql",
			logger.Event("query"),
			logger.Sql(query),
			logger.Args(args...),
		)
	}

	row = executor.QueryRowContext(ctx, query, args...)
	err = row.Err()

	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("query"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Error(err),
		)
	}

	return row, err
}

func QueryTimeout(timeout time.Duration, executor Executor, query string, args ...any) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Query(ctx, executor, query, args...)
}

func Query(ctx context.Context, executor Executor, query string, args ...any) (*sql.Rows, error) {
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapError("exec sql",
			logger.Event("query"),
			logger.Sql(query),
			logger.Args(args...),
		)
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("query"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Error(err),
		)
	}

	return rows, err
}

func ExecTimeout(timeout time.Duration, executor Executor, query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Exec(ctx, executor, query, args...)
}

func Exec(ctx context.Context, executor Executor, query string, args ...any) (sql.Result, error) {
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapError("exec sql",
			logger.Sql(query),
			logger.Args(args...),
		)
	}

	res, err := executor.ExecContext(ctx, query)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Sql(query),
			logger.Args(args...),
			logger.Error(err),
		)
	}

	return res, err
}

func UpdateModel(ctx context.Context, executor Executor, setter string, model Model) (rowsAffected int64, err error) {
	sqlStr, args := Update(model.TableName(), setter, condition.Equal{
		Field: model.PrimaryField(),
		Value: model.PrimaryValue(),
	})

	res, err := Exec(ctx, executor, sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func DeleteModel(ctx context.Context, executor Executor, model Model) (rowsAffected int64, err error) {
	sqlStr, args := Delete(model.TableName(), condition.Equal{
		Field: model.PrimaryField(),
		Value: model.PrimaryValue(),
	})

	res, err := Exec(ctx, executor, sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
