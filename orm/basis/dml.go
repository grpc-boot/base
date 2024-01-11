package basis

import "github.com/grpc-boot/base/v2/orm/condition"

type Row []any

type Columns []string

type Delete func(table string, condition condition.Condition) (sql string, args []any)
type Update func(table, setters string, condition condition.Condition) (sql string, args []any)
type Insert func(table string, columns Columns, rows []Row, ignore bool) (sql string, args []any)
