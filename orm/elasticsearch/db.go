package elasticsearch

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/grpc-boot/base/v2/orm/basis"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	client *http.Client
	opts   Options
}

func (db *Db) Options() Options {
	return db.opts
}

func (db *Db) LoadTableSchema(table string) (t *basis.Table, err error) {
	var columns []basis.Column
	columns, err = db.FetchColumns(table)
	if err != nil {
		return
	}

	t = basis.NewTable(table, columns)
	return
}

func (db *Db) FetchColumns(table string) (columns []basis.Column, err error) {
	return
}

func (db *Db) ShowTables(pattern string) (tables []string, err error) {
	return
}

func NewDb(opt Options) (mysql *Db, err error) {
	db, err := sql.Open("mysql", opt.Dsn())
	if err != nil {
		return
	}

	if opt.MaxIdleConns > 0 {
		db.SetMaxIdleConns(opt.MaxIdleConns)
	}

	if opt.MaxOpenConns > 0 {
		db.SetMaxOpenConns(opt.MaxOpenConns)
	}

	if opt.ConnMaxIdleTimeSecond > 0 {
		db.SetConnMaxIdleTime(time.Duration(opt.ConnMaxIdleTimeSecond) * time.Second)
	}

	if opt.ConnMaxLifetimeSecond > 0 {
		db.SetConnMaxLifetime(time.Duration(opt.ConnMaxLifetimeSecond) * time.Second)
	}

	return &Db{client: &http.Client{}, opts: opt}, nil
}

func (db *Db) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, basis.ErrUnsupportedDriver
}

func (db *Db) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, basis.ErrUnsupportedDriver
}

func (db *Db) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return nil
}

func (db *Db) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, basis.ErrUnsupportedDriver
}

func (db *Db) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, basis.ErrUnsupportedDriver
}
