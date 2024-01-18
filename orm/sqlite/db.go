package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/grpc-boot/base/v2/orm/basis"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	*sql.DB
	opts Options
}

func (db *Db) Options() Options {
	return db.opts
}

func (db *Db) ShowCreateTable(table string) (tableSql string, err error) {
	var rows *sql.Rows
	rows, err = db.Query(fmt.Sprintf("SELECT `sql` FROM sqlite_master WHERE type='table' AND name='%s'", table))
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&tableSql); err != nil {
			return
		}
	}

	return
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
	tableSql, err := db.ShowCreateTable(table)
	if err != nil {
		return
	}

	if len(tableSql) < 1 {
		return
	}

	start := strings.Index(tableSql, "(")
	tableSql = strings.TrimSuffix(tableSql, ")")
	tableSql = strings.ReplaceAll(tableSql[start+1:], "\n", "")
	tableSql = strings.ReplaceAll(tableSql, "\r\n", "")
	tableSql = strings.ReplaceAll(tableSql, "\t", "")
	items := strings.Split(strings.TrimSpace(tableSql), ",")

	columns = make([]basis.Column, 0, len(items))

	for _, item := range items {
		info := strings.SplitN(strings.TrimSpace(item), " ", 2)
		if len(info) == 1 {
			continue
		}

		col := &column{
			f: info[0],
			t: strings.ToLower(strings.TrimSpace(info[1])),
		}
		col.format()

		columns = append(columns, col)
	}

	return
}

func (db *Db) ShowTables(pattern string) (tables []string, err error) {
	var rows *sql.Rows
	if pattern == "" {
		rows, err = db.Query("SELECT DISTINCT tbl_name FROM sqlite_master WHERE tbl_name<>'sqlite_sequence' ORDER BY tbl_name")
	} else {
		rows, err = db.Query(fmt.Sprintf("SELECT DISTINCT tbl_name FROM sqlite_master WHERE tbl_name<>'sqlite_sequence' ORDER BY tbl_name AND tbl_name LIKE '%s'", pattern))
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

func NewDb(opt Options) (sqlite *Db, err error) {
	db, err := sql.Open("sqlite3", opt.Dsn())
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

	return &Db{DB: db, opts: opt}, nil
}
