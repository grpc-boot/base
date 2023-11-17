package components

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/grpc-boot/base/v2/utils"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlOption struct {
	DbName                string `json:"dbName" yaml:"dbName"`
	Host                  string `json:"host" yaml:"host"`
	Port                  uint32 `json:"port" yaml:"port"`
	UserName              string `json:"userName" yaml:"userName"`
	Password              string `json:"password" yaml:"password"`
	CharSet               string `json:"charSet" yaml:"charSet"`
	MaxIdleConns          int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns          int    `json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxIdleTimeSecond int64  `json:"connMaxIdleTimeSecond" yaml:"connMaxIdleTimeSecond"`
	ConnMaxLifetimeSecond int64  `json:"connMaxLifetimeSecond" yaml:"connMaxLifetimeSecond"`
}

type Mysql struct {
	*sql.DB
	Options MysqlOption
}

func (mo *MysqlOption) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		mo.UserName,
		mo.Password,
		mo.Host,
		mo.Port,
		mo.DbName,
		mo.CharSet,
	)
}

func DefaultMysqlOption() MysqlOption {
	return MysqlOption{
		CharSet:               "utf8",
		MaxOpenConns:          8,
		MaxIdleConns:          2,
		ConnMaxLifetimeSecond: 600,
		ConnMaxIdleTimeSecond: 60,
	}
}

func MysqlOptionWithJson(conf string) (opt MysqlOption, err error) {
	opt = DefaultMysqlOption()
	err = utils.JsonDecode(conf, &opt)
	return
}

func MysqlOptionWithYaml(conf string) (opt MysqlOption, err error) {
	opt = DefaultMysqlOption()
	err = utils.YamlDecode(conf, &opt)
	return
}

func NewMysql(opt MysqlOption) (mysql *Mysql, err error) {
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

	return &Mysql{DB: db, Options: opt}, nil
}
