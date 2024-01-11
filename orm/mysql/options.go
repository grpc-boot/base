package mysql

import (
	"fmt"

	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/utils"
)

func DefaultMysqlOption() Options {
	return Options{
		CharSet:               "utf8",
		MaxOpenConns:          8,
		MaxIdleConns:          2,
		ConnMaxLifetimeSecond: 600,
		ConnMaxIdleTimeSecond: 60,
	}
}

func Flag2Options(f *basis.Flag) Options {
	opt := DefaultMysqlOption()

	opt.UserName = f.UserName()
	opt.Password = f.Password()
	opt.Host = f.Host()
	opt.Port = f.Port()
	opt.DbName = f.DbName()
	opt.CharSet = f.Charset()

	return opt
}

func OptionsWithJson(conf string) (opt Options, err error) {
	opt = DefaultMysqlOption()
	err = utils.JsonDecode(conf, &opt)
	return
}

func OptionsWithYaml(conf string) (opt Options, err error) {
	opt = DefaultMysqlOption()
	err = utils.YamlDecode(conf, &opt)
	return
}

type Options struct {
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

func (o *Options) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		o.UserName,
		o.Password,
		o.Host,
		o.Port,
		o.DbName,
		o.CharSet,
	)
}
