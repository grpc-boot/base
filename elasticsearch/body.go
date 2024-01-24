package elasticsearch

import (
	"github.com/grpc-boot/base/v2/elasticsearch/result"
	"github.com/grpc-boot/base/v2/kind"
	"github.com/grpc-boot/base/v2/utils"
)

type Body kind.JsonParam

func (b Body) WithProperties(properties result.MappingProperties) {
	if len(properties) == 0 {
		return
	}

	b["properties"] = properties
}

func (b Body) WithQuery(query any) {
	b["query"] = query
}

func (b Body) WithSql(sqlStr string, args ...any) {
	b["query"] = sqlStr

	if len(args) > 0 {
		b["params"] = args
	}
}

func (b Body) WithFetchSize(size int64) {
	b["fetch_size"] = size
}

func (b Body) WithCursor(cursor string) {
	b["cursor"] = cursor
}

func (b Body) WithArgs(args ...Arg) {
	if len(args) == 0 {
		return
	}

	for _, arg := range args {
		b[arg.key] = arg.value
	}
}

func (b Body) Marshal() (body []byte, err error) {
	return utils.JsonMarshal(b)
}
