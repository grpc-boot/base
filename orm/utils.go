package orm

import (
	"context"
	"database/sql"
	"time"

	"github.com/grpc-boot/base/v2/condition"
	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/orm/mysql"

	"go.uber.org/zap/zapcore"
)

func PrepareTimeout(timeout time.Duration, executor basis.Executor, query string) (*sql.Stmt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Prepare(ctx, executor, query)
}

func Prepare(ctx context.Context, executor basis.Executor, query string) (*sql.Stmt, error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("prepare query start"),
			logger.Sql(query),
		)
	}

	stmt, err := executor.PrepareContext(ctx, query)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("prepare query error"),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("prepare query done"),
			logger.Sql(query),
			logger.Duration(time.Since(start)),
		)
	}

	return stmt, err
}

func ExecStmtTimeout(timeout time.Duration, stmt basis.Stmt, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return ExecStmt(ctx, stmt, args...)
}

func ExecStmt(ctx context.Context, stmt basis.Stmt, args ...any) (sql.Result, error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("exec stmt start"),
			logger.Args(args...),
		)
	}

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("exec stmt error"),
			logger.Args(args...),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("exec stmt done"),
			logger.Args(args...),
			logger.Duration(time.Since(start)),
		)
	}

	return res, err
}

func QueryStmtTimeout(timeout time.Duration, stmt basis.Stmt, args ...any) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return QueryStmt(ctx, stmt, args...)
}

func QueryStmt(ctx context.Context, stmt basis.Stmt, args ...any) (*sql.Rows, error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("query stmt start"),
			logger.Args(args...),
		)
	}

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("query stmt error"),
			logger.Args(args...),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("query stmt done"),
			logger.Args(args...),
			logger.Duration(time.Since(start)),
		)
	}

	return rows, err
}

func QueryRowStmtTimeout(timeout time.Duration, stmt basis.Stmt, args ...any) (*sql.Row, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return QueryRowStmt(ctx, stmt, args...)
}

func QueryRowStmt(ctx context.Context, stmt basis.Stmt, args ...any) (*sql.Row, error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("query row stmt start"),
			logger.Args(args...),
		)
	}

	row := stmt.QueryRowContext(ctx, args...)
	err := row.Err()
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("query row stmt error"),
			logger.Args(args...),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("query row stmt done"),
			logger.Args(args...),
			logger.Duration(time.Since(start)),
		)
	}

	return row, err
}

func BeginTimeout(timeout time.Duration, executor basis.Executor, opts *sql.TxOptions) (*sql.Tx, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Begin(ctx, executor, opts)
}

func Begin(ctx context.Context, executor basis.Executor, opts *sql.TxOptions) (*sql.Tx, error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("begin transaction start"),
		)
	}

	tx, err := executor.BeginTx(ctx, opts)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("begin transaction error"),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("begin transaction done"),
			logger.Duration(time.Since(start)),
		)
	}

	return tx, err
}

func Commit(tx basis.Tx) error {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("commit transaction start"),
		)
	}

	err := tx.Commit()
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("commit transaction error"),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("commit transaction done"),
			logger.Duration(time.Since(start)),
		)
	}

	return err
}

func Rollback(tx basis.Tx) error {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("rollback transaction start"),
		)
	}

	err := tx.Rollback()
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("rollback transaction error"),
			logger.Duration(time.Since(start)),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("rollback transaction done"),
			logger.Duration(time.Since(start)),
		)
	}

	return err
}

func QueryWithQueryTimeout(timeout time.Duration, executor basis.Executor, q basis.Query) (rows *sql.Rows, err error) {
	sqlStr, args := q.Sql()

	return QueryTimeout(timeout, executor, sqlStr, args...)
}

func QueryWithQuery(ctx context.Context, executor basis.Executor, q basis.Query) (rows *sql.Rows, err error) {
	sqlStr, args := q.Sql()

	return Query(ctx, executor, sqlStr, args...)
}

func QueryRowTimeout(timeout time.Duration, executor basis.Executor, query string, args ...any) (row *sql.Row, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return QueryRow(ctx, executor, query, args...)
}

func QueryRow(ctx context.Context, executor basis.Executor, query string, args ...any) (row *sql.Row, err error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapError("exec sql",
			logger.Event("query row start"),
			logger.Sql(query),
			logger.Args(args...),
		)
	}

	row = executor.QueryRowContext(ctx, query, args...)
	err = row.Err()

	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("query row error"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapError("exec sql",
			logger.Event("query row done"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Duration(time.Since(start)),
		)
	}

	return row, err
}

func QueryTimeout(timeout time.Duration, executor basis.Executor, query string, args ...any) (*sql.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Query(ctx, executor, query, args...)
}

func Query(ctx context.Context, executor basis.Executor, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("query start"),
			logger.Sql(query),
			logger.Args(args...),
		)
	}

	rows, err := executor.QueryContext(ctx, query, args...)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("query error"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("query done"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Duration(time.Since(start)),
		)
	}

	return rows, err
}

func ExecTimeout(timeout time.Duration, executor basis.Executor, query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Exec(ctx, executor, query, args...)
}

func Exec(ctx context.Context, executor basis.Executor, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("exec start"),
			logger.Sql(query),
			logger.Args(args...),
		)
	}

	res, err := executor.ExecContext(ctx, query, args...)
	if err != nil {
		logger.ZapError("exec sql failed",
			logger.Event("exec error"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Error(err),
			logger.Duration(time.Since(start)),
		)
	} else if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("exec sql",
			logger.Event("exec done"),
			logger.Sql(query),
			logger.Args(args...),
			logger.Duration(time.Since(start)),
		)
	}

	return res, err
}

func InsertRowTimeoutWithMysql(timeout time.Duration, executor basis.Executor, table string, columns basis.Columns, row basis.Row, ignore bool) (lastInsertId int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return InsertRow(ctx, executor, table, columns, row, ignore, mysql.Insert)
}

func InsertRow(ctx context.Context, executor basis.Executor, table string, columns basis.Columns, row basis.Row, ignore bool, inserter basis.Insert) (lastInsertId int64, err error) {
	sqlStr, args := inserter(table, columns, []basis.Row{row}, ignore)
	res, err := Exec(ctx, executor, sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func InsertRowsTimeoutWithMysql(timeout time.Duration, executor basis.Executor, table string, columns basis.Columns, rows []basis.Row, ignore bool) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return InsertRows(ctx, executor, table, columns, rows, ignore, mysql.Insert)
}

func InsertRows(ctx context.Context, executor basis.Executor, table string, columns basis.Columns, rows []basis.Row, ignore bool, inserter basis.Insert) (rowsAffected int64, err error) {
	sqlStr, args := inserter(table, columns, rows, ignore)
	res, err := Exec(ctx, executor, sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func DeleteTimeoutWithMysql(timeout time.Duration, executor basis.Executor, table string, condition condition.Condition) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return Delete(ctx, executor, table, condition, mysql.Delete)
}

func Delete(ctx context.Context, executor basis.Executor, table string, condition condition.Condition, deleter basis.Delete) (rowsAffected int64, err error) {
	sqlStr, args := deleter(table, condition)
	res, err := Exec(ctx, executor, sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func UpdateTimeoutWithMysql(timeout time.Duration, executor basis.Executor, table string, setter basis.Setter, condition condition.Condition) (rowsAffected int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return Update(ctx, executor, table, setter, condition, mysql.Update)
}

func Update(ctx context.Context, executor basis.Executor, table string, setter basis.Setter, condition condition.Condition, updater basis.Update) (rowsAffected int64, err error) {
	setterStr, err := setter.Build()
	if err != nil {
		return 0, err
	}

	sqlStr, args := updater(table, setterStr, condition)
	res, err := Exec(ctx, executor, sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func UpdateModel(ctx context.Context, executor basis.Executor, setter basis.Setter, model basis.Model) (rowsAffected int64, err error) {
	return Update(ctx, executor, model.TableName(), setter, condition.Equal{
		Field: model.PrimaryField(),
		Value: model.PrimaryValue(),
	}, model.Updater())
}

func DeleteModel(ctx context.Context, executor basis.Executor, model basis.Model) (rowsAffected int64, err error) {
	return Delete(ctx, executor, model.TableName(), condition.Equal{
		Field: model.PrimaryField(),
		Value: model.PrimaryValue(),
	}, model.Deleter())
}
