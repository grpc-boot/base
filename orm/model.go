package orm

import "github.com/grpc-boot/base/v2/orm/base"

type Model interface {
	TableName() string
	PrimaryField() string
	PrimaryValue() any
	GetDefault(name string) string
	GetLabel(name string) string
	GetEnums(name string) []string
	GetSize(name string) int
	Query() base.Query
}
