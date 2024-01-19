package sqlite

import (
	"fmt"
	"time"

	"github.com/grpc-boot/base/v2/orm/basis"
	"github.com/grpc-boot/base/v2/utils"
)

func DefaultOption() Options {
	return Options{
		MaxOpenConns:          8,
		MaxIdleConns:          2,
		ConnMaxLifetimeSecond: 600,
		ConnMaxIdleTimeSecond: 60,
	}
}

func Flag2Options(f *basis.Flag) Options {
	opt := DefaultOption()
	opt.DbName = f.DbName()
	return opt
}

func Flag2Generator(f *basis.Flag) (basis.Generator, error) {
	return NewDb(Flag2Options(f))
}

func OptionsWithJson(conf string) (opt Options, err error) {
	opt = DefaultOption()
	err = utils.JsonDecode(conf, &opt)
	return
}

func OptionsWithYaml(conf string) (opt Options, err error) {
	opt = DefaultOption()
	err = utils.YamlDecode(conf, &opt)
	return
}

type Options struct {
	DbName                string `json:"dbName" yaml:"dbName"`
	MaxIdleConns          int    `json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns          int    `json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxIdleTimeSecond int64  `json:"connMaxIdleTimeSecond" yaml:"connMaxIdleTimeSecond"`
	ConnMaxLifetimeSecond int64  `json:"connMaxLifetimeSecond" yaml:"connMaxLifetimeSecond"`
}

func (o *Options) ConnMaxIdleTime() time.Duration {
	return time.Duration(o.ConnMaxIdleTimeSecond) * time.Second
}

func (o *Options) ConnMaxLifetime() time.Duration {
	return time.Duration(o.ConnMaxLifetimeSecond) * time.Second
}

func (o *Options) Dsn() string {
	return fmt.Sprintf("%s?parseTime=True",
		o.DbName,
	)
}
