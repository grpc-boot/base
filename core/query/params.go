package query

import (
	"strings"
)

const (
	And = `AND`
	Or  = `OR`
)

type Params []interface{}

type AndParams Params

type OrParams Params

func (ap AndParams) Build() (sql string, err error) {
	return BuildParams(Params(ap), And)
}

func (op OrParams) Build() (sql string, err error) {
	return BuildParams(Params(op), Or)
}

func BuildParams(params Params, opt string) (sql string, err error) {
	switch len(params) {
	case 0:
		return
	case 1:
		return params[0].Build()
	}

	sqlStr, err := params[0].Build()
	if err != nil {
		return
	}

	var buffer strings.Builder
	buffer.WriteString(sqlStr)
	for _, conf := range params[1:] {
		if sqlStr, err = conf.Build(); err != nil {
			return
		}
		buffer.WriteByte(' ')
		buffer.WriteString(opt)
		buffer.WriteByte(' ')
		buffer.WriteString(sql)
	}

	sql = buffer.String()
	return
}
