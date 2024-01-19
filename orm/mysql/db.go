package mysql

import (
	"database/sql"
	"fmt"

	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/utils"

	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	*sql.DB
	opts Options
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
	var rows *sql.Rows
	rows, err = db.Query(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", table))
	if err != nil {
		return
	}

	defer rows.Close()

	columns = []basis.Column{}

	var (
		f  []byte
		t  []byte
		_c []byte
		n  []byte
		k  []byte
		d  []byte
		e  []byte
		_p []byte
		c  []byte
	)

	for rows.Next() {
		if err = rows.Scan(&f, &t, &_c, &n, &k, &d, &e, &_p, &c); err != nil {
			return
		}
		col := &column{
			f: utils.Bytes2String(f),
			t: utils.Bytes2String(t),
			n: utils.Bytes2String(n),
			k: utils.Bytes2String(k),
			d: utils.Bytes2String(d),
			e: utils.Bytes2String(e),
			c: utils.Bytes2String(c),
		}
		col.format()

		columns = append(columns, col)
	}

	return
}

func (db *Db) ShowTables(pattern string) (tables []string, err error) {
	var rows *sql.Rows
	if pattern == "" {
		rows, err = db.Query("SHOW TABLES")
	} else {
		rows, err = db.Query(fmt.Sprintf("SHOW TABLES LIKE '%s'", pattern))
	}

	if err != nil {
		return
	}

	defer rows.Close()

	tables = []string{}

	for rows.Next() {
		var tbl string
		if err = rows.Scan(&tbl); err != nil {
			return
		}

		tables = append(tables, tbl)
	}

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
		db.SetConnMaxIdleTime(opt.ConnMaxIdleTime())
	}

	if opt.ConnMaxLifetimeSecond > 0 {
		db.SetConnMaxLifetime(opt.ConnMaxLifetime())
	}

	return &Db{DB: db, opts: opt}, nil
}
